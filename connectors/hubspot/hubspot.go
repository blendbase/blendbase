package hubspot

import (
	"blendbase/connectors"
	"blendbase/graph/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	HSBaseUrlV3 = "https://api.hubapi.com/crm/v3/objects"
)

type Client struct {
	BaseURL     string
	AccessToken string
	HTTPClient  *http.Client
}

type ErrorResponse struct {
	Status        string `json:"status"`
	Message       string `json:"message"`
	CorrelationId string `json:"correlationId"`
	Category      string `json:"category"`
}

type HubspotError struct {
	StatusCode int
	Err        error
}

type HSAssociationListSuccessResponse struct {
	Results []HSAssociationsListItem `json:"results"`
}

type HSAssociationsListItem struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

func (e *HubspotError) Error() string {
	return e.Err.Error()
}

func HubspotClient(accessToken string) *Client {
	return &Client{
		BaseURL:     HSBaseUrlV3,
		AccessToken: accessToken,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (client *Client) sendRequest(req *http.Request, response interface{}) *HubspotError {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("authorization", "Bearer "+client.AccessToken)

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return &HubspotError{
			Err: err,
		}
	}
	defer res.Body.Close()

	url := strings.Split(res.Request.URL.String(), "?")[0]

	if res.StatusCode >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			log.WithFields(log.Fields{
				"status_code": res.StatusCode,
				"url":         url,
				"method":      req.Method,
				"err":         errRes.Message,
			}).Info("HubSpot request")

			return &HubspotError{
				StatusCode: res.StatusCode,
				Err:        errors.New(errRes.Message),
			}
		}

		log.WithFields(log.Fields{
			"status_code": res.StatusCode,
			"url":         url,
			"method":      req.Method,
			"err":         "cannot parse error response",
		}).Info("HubSpot request")

		return &HubspotError{
			StatusCode: res.StatusCode,
			Err:        fmt.Errorf("unexpected status code %d", res.StatusCode),
		}
	}

	log.WithFields(log.Fields{
		"status_code": res.StatusCode,
		"url":         url,
		"method":      req.Method,
	}).Info("HubSpot request")

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return &HubspotError{
			StatusCode: res.StatusCode,
			Err:        err,
		}
	}

	return nil
}

func (client *Client) get(ctx context.Context, objectPath, objectId string, props []string, obj interface{}) error {
	query := url.Values{}
	query.Set("properties", strings.Join(props, ","))

	url := fmt.Sprintf("%s/%s/%s?%s", client.BaseURL, objectPath, objectId, query.Encode())

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	if err := client.sendRequest(req, obj); err != nil {
		if err.StatusCode == 404 {
			return errors.New("not found")
		}

		return err
	}
	return nil
}

func (client *Client) create(ctx context.Context, objectPath, payload interface{}, response interface{}) error {
	query := url.Values{}
	url := fmt.Sprintf("%s/%s?%s", client.BaseURL, objectPath, query.Encode())

	payloadString, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadString))
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	if err := client.sendRequest(req, &response); err != nil {
		return err
	}
	return nil
}

func (client *Client) update(ctx context.Context, objectPath string, objectId string, payload interface{}, response interface{}) error {
	query := url.Values{}
	url := fmt.Sprintf("%s/%s/%s?%s", client.BaseURL, objectPath, objectId, query.Encode())

	payloadString, _ := json.Marshal(payload)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payloadString))
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	if err := client.sendRequest(req, &response); err != nil {
		return err
	}
	return nil
}

func (client *Client) delete(ctx context.Context, objectPath string, objectId string, response interface{}) error {
	query := url.Values{}
	url := fmt.Sprintf("%s/%s/%s?%s", client.BaseURL, objectPath, objectId, query.Encode())

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	if err := client.sendRequest(req, &response); err != nil {
		return err
	}

	return nil
}

func (client *Client) list(ctx context.Context, objectPath string, first int, after *string, props []string, response interface{}) error {
	query := url.Values{}

	query.Set("properties", strings.Join(props, ","))

	limitBuffer := 1 // +1 to see if there are more pages
	if after != nil {
		query.Set("after", connectors.DecodeCursor(*after))
		// increasing the buffer because Hubspot API returns the "after" item, we're going to drop it
		limitBuffer += 1
	}
	query.Set("limit", fmt.Sprint(first+limitBuffer))

	url := fmt.Sprintf("%s/%s?%s", client.BaseURL, objectPath, query.Encode())
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	if err := client.sendRequest(req, &response); err != nil {
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

	// remove the "after" item
	if after != nil && recordsLenght > 0 {
		start = 1
	}

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

func parseHSDateTime(dateTime *string) *time.Time {
	if dateTime == nil {
		return nil
	}

	t, err := time.Parse(time.RFC3339, *dateTime)

	if err != nil {
		log.WithFields(log.Fields{
			"dateTime": dateTime,
		}).Error("Failed to parse dateTime")
	}

	return &t
}

func formatHSDateTime(dateTime time.Time) string {
	return dateTime.Format(time.RFC3339)
}
