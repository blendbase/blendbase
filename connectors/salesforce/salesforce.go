package salesforce

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"blendbase/config"
	"blendbase/connectors"
	"blendbase/graph/model"
	"blendbase/integrations"
)

const (
	BaseUrlTemplate        = "https://%s.my.salesforce.com/services/data/v53.0"
	SALESFORCE_TIME_FORMAT = "2006-01-02T15:04:05.000+0000"
)

type Client struct {
	SalesforceInstanceSubdomain string
	OAuthStateString            string
	HTTPClient                  *http.Client

	app                 *config.App
	consumerOAuthConfig *integrations.ConsumerOauth2Configuration
}

type SFErrorResponse = []SFError
type SFError struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type SFListQuerySuccessResponseBase struct {
	TotalSize int  `json:"totalSize"`
	Done      bool `json:"done"`
}

type SFCreateObjectResponse struct {
	ID      string    `json:"id"`
	Success bool      `json:"success"`
	Errors  []SFError `json:"errors"`
}

func SaleforceClient(app *config.App, consumerOAuthConfig *integrations.ConsumerOauth2Configuration) *Client {
	context := context.Background()

	customSettings, _ := consumerOAuthConfig.GetCustomSettings()

	return &Client{
		SalesforceInstanceSubdomain: customSettings.SalesforceInstanceSubdomain,
		OAuthStateString:            os.Getenv("OAUTH_STATE_STRING"),
		HTTPClient:                  getOAuthConfig(consumerOAuthConfig).Client(context, getOAuthToken(consumerOAuthConfig)),
		app:                         app,
		consumerOAuthConfig:         consumerOAuthConfig,
	}
}

func LoadClientFromDB(app *config.App, consumer *integrations.Consumer) (*Client, error) {
	var consumerIntegration integrations.ConsumerIntegration
	if err := app.DB.Where("consumer_id = ?", consumer.ID).Where("service_code = ?", "crm_salesforce").Order("created_at DESC").First(&consumerIntegration).Error; err != nil {
		return nil, err
	}

	var consumerOAuthConfig integrations.ConsumerOauth2Configuration
	if err := app.DB.Where("consumer_integration_id = ?", consumerIntegration.ID).First(&consumerOAuthConfig).Error; err != nil {
		return nil, err
	}

	return SaleforceClient(app, &consumerOAuthConfig), nil
}

func (client *Client) baseUrl() string {
	return fmt.Sprintf(BaseUrlTemplate, client.SalesforceInstanceSubdomain)
}

func (client *Client) sendAPIRequest(req *http.Request, response interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 401 {
		// SF doesn't provide an expiration date for tokens, so we need to refresh the token
		// if the call to the API fails with a 401 error
		if err := client.refreshToken(); err != nil {
			return err
		}

		// retry
		res, err = client.HTTPClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
	}

	log.WithFields(log.Fields{
		"status_code": res.StatusCode,
	}).Info("Salesforce request")

	if res.StatusCode >= http.StatusBadRequest {
		var errRes SFErrorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes[0].Message)
		}

		return fmt.Errorf("unexpected error response message structure: %s", err)
	}
	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Errorf("Error decoding response body: %s", err)
		return err
	}

	return nil
}

// === API Specific Functions ===
// Generalized list objects request
func (client *Client) list(objectName string, fields []string, first int, after *string, response interface{}) error {
	query := url.Values{}

	// +1 to see if there are more pages
	first += 1

	var selectQuery string
	if after != nil {
		decodedAfter := connectors.DecodeCursor(*after)
		selectQuery = fmt.Sprintf(`SELECT %s from %s where Id > '%s' order by Id limit %d`, strings.Join(fields, ","), objectName, decodedAfter, first)
	} else {
		selectQuery = fmt.Sprintf("SELECT %s from %s order by Id limit %d", strings.Join(fields, ","), objectName, first)
	}

	log.Debugf("Listing objects of %s type: %s", objectName, selectQuery)

	query.Set("q", selectQuery)
	url := fmt.Sprintf("%s/query?%s", client.baseUrl(), query.Encode())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Warn(err)
		return err
	}

	if err := client.sendAPIRequest(req, &response); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (client *Client) prepareListResults(first int, after *string, edgesPtr interface{}) (reflect.Value, *model.PageInfo) {
	recordsValue := reflect.ValueOf(edgesPtr).Elem()
	recordsLenght := recordsValue.Len()

	start := 0
	end := recordsLenght
	hasNextPage := false

	if (recordsLenght - start) > first {
		// remove the item that was added to see if there are more pages
		end = first + start
		hasNextPage = true
	}

	ret := recordsValue.Slice(start, end)
	pageInfo := &model.PageInfo{
		HasNextPage: hasNextPage,
	}

	if recordsLenght > 0 {
		startCursor := ret.Index(0).Elem().FieldByName("Cursor").String()
		endCursor := ret.Index(ret.Len() - 1).Elem().FieldByName("Cursor").String()

		pageInfo.StartCursor = &startCursor
		pageInfo.EndCursor = &endCursor
	}

	return ret, pageInfo
}

func (client *Client) listWithWhere(objectName string, fields []string, whereFilter string, response interface{}) error {
	query := url.Values{}

	selectQuery := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(fields, ","), objectName, whereFilter)

	log.Infof("Listing objects of %s type: %s", objectName, selectQuery)

	query.Set("q", selectQuery)
	url := fmt.Sprintf("%s/query?%s", client.baseUrl(), query.Encode())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Warn(err)
		return err
	}

	if err := client.sendAPIRequest(req, &response); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// Gets an object by ID
func (client *Client) get(objectName string, objectId string, fields []string, response interface{}) error {
	query := url.Values{}
	fieldsString := strings.Join(fields, ",")
	query.Set("q", fieldsString)
	log.Infof("Fetching %s %s and fields %s", objectName, objectId, fieldsString)

	url := fmt.Sprintf("%s/sobjects/%s/%s?%s", client.baseUrl(), objectName, objectId, query.Encode())
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	if err := client.sendAPIRequest(req, &response); err != nil {
		return err
	}

	return nil
}

// Create a new object
func (client *Client) create(objectName string, payload interface{}) (string, error) {
	url := fmt.Sprintf("%s/sobjects/%s", client.baseUrl(), objectName)
	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Error encoding payload %v: %s", payload, err)
		return "", err
	}
	log.Debugf("creating %s with %s payload", objectName, string(encodedPayload))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(encodedPayload)))

	if err != nil {
		return "", err
	}

	response := SFCreateObjectResponse{}
	if err := client.sendAPIRequest(req, &response); err != nil {
		return "", err
	}

	if !response.Success && len(response.Errors) > 0 {
		return "", errors.New(formatSFErrors(response.Errors))
	}

	return response.ID, nil
}

func (client *Client) update(objectName string, objectId string, payload interface{}) (bool, error) {
	url := fmt.Sprintf("%s/sobjects/%s/%s", client.baseUrl(), objectName, objectId)
	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Error encoding payload %v: %s", payload, err)
		return false, err
	}
	log.Debugf("updating %s %s with %s payload", objectName, objectId, string(encodedPayload))

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer([]byte(encodedPayload)))

	if err != nil {
		return false, err
	}

	if err := client.sendAPIRequest(req, nil); err != nil {
		return false, err
	}

	return true, nil
}

func (client *Client) delete(objectName string, objectId string) (bool, error) {
	url := fmt.Sprintf("%s/sobjects/%s/%s", client.baseUrl(), objectName, objectId)
	log.Debugf("archiving %s %s", objectName, objectId)

	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		return false, err
	}

	if err := client.sendAPIRequest(req, nil); err != nil {
		return false, err
	}

	return true, nil
}

func formatSFErrors(errors []SFError) string {
	messages := []string{}
	for _, error := range errors {
		messages = append(messages, fmt.Sprintf("%s: %s\n", error.ErrorCode, error.Message))
	}

	return strings.Join(messages, ", ")
}

func parseSFDateTime(dateTime *string) *time.Time {
	if dateTime == nil {
		return nil
	}

	t, err := time.Parse(SALESFORCE_TIME_FORMAT, *dateTime)

	if err != nil {
		log.WithFields(log.Fields{
			"dateTime": *dateTime,
		}).Error("Failed to parse dateTime")
	}

	return &t
}

func formatSFDateTime(dateTime time.Time) string {
	return dateTime.Format(SALESFORCE_TIME_FORMAT)
}
