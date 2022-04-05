package connect

import (
	"blendbase/connectors"
	"blendbase/graph/model"
	"blendbase/integrations"
	"blendbase/misc/gormext"
	"encoding/json"
	"fmt"
	"os"

	"blendbase/config"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

type ConnectClient struct {
	ConsumerID uuid.UUID
	App        *config.App
}

type OAuth2Settings struct {
	ClientID     string
	ClientSecret string
}

func NewConnectClient(app *config.App, consumerID uuid.UUID) *ConnectClient {
	return &ConnectClient{
		App:        app,
		ConsumerID: consumerID,
	}
}

func (client *ConnectClient) loadOAuth2Configuration(consumerIntegrationID uuid.UUID) (*integrations.ConsumerOauth2Configuration, error) {
	oauth2Configuration := integrations.ConsumerOauth2Configuration{}
	query := client.App.DB.Where("consumer_integration_id = ?", consumerIntegrationID).First(&oauth2Configuration)
	if err := query.Error; err != nil {
		return nil, nil
	}

	return &oauth2Configuration, nil
}

func (client *ConnectClient) ListIntegrations() ([]*model.ConsumerIntegration, error) {
	consumerIntegrations := []integrations.ConsumerIntegration{}
	if err := client.App.DB.Where("consumer_id = ?", client.ConsumerID).Find(&consumerIntegrations).Error; err != nil {
		log.Errorf("Error finding integrations for %s: %s", client.ConsumerID, err)
	}

	outputIntegrations := []*model.ConsumerIntegration{}

	// Loop through all available connectors
	for _, availableConnector := range connectors.AvailableConnectors {
		connector := availableConnector
		enabled := false // disabled by default

		outputIntegration := client.createOutputIntegrationFromConnector(&connector)

		// find the match from the DB consumer integration by service code
		consumerIntegration := findConsumerIntegrationByServiceCode(&consumerIntegrations, connector.ServiceCode)

		// if there is a match, assing values from the exiting integration
		if consumerIntegration != nil {
			enabled = consumerIntegration.Enabled
			ID := consumerIntegration.ID.String()
			outputIntegration.ID = &ID

			clientCredentialsSet := false
			tokensSet := false
			oauth2Config, _ := client.loadOAuth2Configuration(consumerIntegration.ID)

			if oauth2Config != nil {
				clientCredentialsSet = (oauth2Config.ClientID.Raw != "" && oauth2Config.ClientSecret.Raw != "")
				tokensSet = (oauth2Config.AccessToken.Raw != "" && oauth2Config.RefreshToken.Raw != "")
			}

			outputIntegration.Oauth2Metadata = &model.OAuth2Metadata{
				ClientCredentialsSet: clientCredentialsSet,
				TokensSet:            tokensSet,
			}
		}
		outputIntegration.Enabled = &enabled

		outputIntegrations = append(outputIntegrations, outputIntegration)
	}

	return outputIntegrations, nil
}

// Enables or disables an integration
// Adds or removes the integration to the DB if it doesn't exist
func (client *ConnectClient) EnableIntegration(serviceCode string, enabled bool) (bool, error) {
	connector := findConnectorByServiceCode(serviceCode)
	if connector == nil {
		return false, fmt.Errorf("cannot add %s integration. %s is not in the list of available integrations for this client", serviceCode, serviceCode)
	}

	consumerIntegration := integrations.ConsumerIntegration{}
	query := client.App.DB.FirstOrCreate(&consumerIntegration, integrations.ConsumerIntegration{
		ConsumerID:  client.ConsumerID,
		ServiceCode: connector.ServiceCode,
		Type:        connector.Type,
	})
	if err := query.Error; err != nil {
		return false, fmt.Errorf("error creating integration for %s: %s", serviceCode, err)
	}

	// updated enabled flag for the integration
	if err := client.App.DB.Model(&consumerIntegration).Updates(map[string]interface{}{"enabled": &enabled}).Error; err != nil {
		return false, fmt.Errorf("error enabling integrations #%s: %s", consumerIntegration.ID, err)
	}

	//disable remaining integrations of the same type
	if enabled {
		client.App.DB.Model(
			&integrations.ConsumerIntegration{},
		).Where(
			"consumer_id = ?", client.ConsumerID,
		).Where(
			"id <> ?", consumerIntegration.ID,
		).Where(
			"type = ?", consumerIntegration.Type,
		).Update("enabled", false)
	}

	return true, nil
}

// Configure OAuth2 settings for the existing consumer integration
// Arguments:
//   consumerIntegrationID: the ID of the consumer integration
//   oauth2Settings: the OAuth2 settings to configure (clientID, clientSecret)
// Returns:
//   true if the settings were successfully configured
func (client *ConnectClient) ConfigureOAuth2(consumerIntegrationID uuid.UUID, oauth2Settings *model.OAuth2ConfigurationInput) (bool, error) {
	consumerIntegration := integrations.ConsumerIntegration{}
	if err := client.App.DB.Where("consumer_id = ?", client.ConsumerID).Where("id = ?", consumerIntegrationID).First(&consumerIntegration).Error; err != nil {
		return false, fmt.Errorf("error finding integration #%s: %s", consumerIntegrationID, err)
	}

	oauth2Configuration := integrations.ConsumerOauth2Configuration{}
	query := client.App.DB.FirstOrCreate(&oauth2Configuration, integrations.ConsumerOauth2Configuration{
		ConsumerIntegrationID: consumerIntegration.ID,
	})
	if err := query.Error; err != nil {
		return false, fmt.Errorf("error creating/finding oauth2 configuration for consumer integration  #%s: %s", consumerIntegration.ID.String(), err)
	}

	if oauth2Settings.ClientID == nil || oauth2Settings.ClientSecret == nil {
		return false, fmt.Errorf("clientID and clientSecret must be provided")
	}

	if *oauth2Settings.ClientID == "" || *oauth2Settings.ClientSecret == "" {
		return false, fmt.Errorf("clientID and clientSecret must be non-empty")
	}

	if consumerIntegration.ServiceCode == connectors.CONNECTOR_CRM_SALESFORCE && (oauth2Settings.SalesforceInstanceSubdomain == nil || *oauth2Settings.SalesforceInstanceSubdomain == "") {
		return false, fmt.Errorf("salesforceInstanceSubdomain must be provided")
	}

	oauth2Configuration.ClientID = gormext.EncryptedValue{Raw: *oauth2Settings.ClientID}
	oauth2Configuration.ClientSecret = gormext.EncryptedValue{Raw: *oauth2Settings.ClientSecret}
	oauth2Configuration.RedirectURL = client.getCallbackUrl(consumerIntegration.ServiceCode)

	if consumerIntegration.ServiceCode == connectors.CONNECTOR_CRM_SALESFORCE {
		customSettings := integrations.ConsumerOauth2ConfigurationCustomSettings{
			SalesforceInstanceSubdomain: *oauth2Settings.SalesforceInstanceSubdomain,
		}
		customSettingsJson, _ := json.Marshal(customSettings)
		oauth2Configuration.CustomSettings = datatypes.JSON(customSettingsJson)
	}

	if err := client.App.DB.Save(&oauth2Configuration).Error; err != nil {
		return false, fmt.Errorf("error saving oauth2 configuration for consumer integration #%s: %s", consumerIntegration.ID.String(), err)
	}

	return true, nil
}

// Creates new consumer and returns the ID of the new consumer
func (client *ConnectClient) CreateConsumer() (uuid.UUID, error) {
	consumer := integrations.Consumer{}
	if err := client.App.DB.Create(&consumer).Error; err != nil {
		log.Infof("Created consumer error: %v", err)
		return uuid.Nil, fmt.Errorf("error creating consumer: %s", err)
	}

	return consumer.ID, nil
}

// Sets Consumer Integration Secret
func (client *ConnectClient) SetConsumerIntegrationSecret(consumerIntegrationID uuid.UUID, secret string) error {
	consumerIntegration := integrations.ConsumerIntegration{}
	if err := client.App.DB.Where("consumer_id = ?", client.ConsumerID).Where("id = ?", consumerIntegrationID).First(&consumerIntegration).Error; err != nil {
		return fmt.Errorf("error finding integration #%s: %s", consumerIntegrationID, err)
	}

	if secret == "" {
		return fmt.Errorf("secret key must be provided")
	}

	consumerIntegration.Secret = gormext.EncryptedValue{Raw: secret}

	if err := client.App.DB.Save(&consumerIntegration).Error; err != nil {
		return fmt.Errorf("error saving secret for consumer integration #%s: %s", consumerIntegration.ID.String(), err)
	}

	return nil
}

// -------- Private --------
func findConsumerIntegrationByServiceCode(integrations *[]integrations.ConsumerIntegration, serviceCode string) *integrations.ConsumerIntegration {
	for _, integration := range *integrations {
		if integration.ServiceCode == serviceCode {
			return &integration
		}
	}

	return nil
}

func findConnectorByServiceCode(serviceCode string) *connectors.Connector {
	for _, connector := range connectors.AvailableConnectors {
		if connector.ServiceCode == serviceCode {
			return &connector
		}
	}

	return nil
}

func (client *ConnectClient) createOutputIntegrationFromConnector(connector *connectors.Connector) *model.ConsumerIntegration {
	loginUrl := client.getLoginUrl(connector.ServiceCode)
	callbackUrl := client.getCallbackUrl(connector.ServiceCode)
	var authType model.AuthType

	switch connector.AuthType {
	case connectors.AUTH_TYPE_OAUTH2:
		authType = model.AuthTypeOauth2
	default:
		authType = model.AuthTypeSecret
	}

	return &model.ConsumerIntegration{
		Type:        &connector.Type,
		ServiceCode: &connector.ServiceCode,
		ServiceName: &connector.Name,
		Description: &connector.Description,
		LoginURL:    &loginUrl,
		CallbackURL: &callbackUrl,
		AuthType:    authType,
	}
}

func (client *ConnectClient) getLoginUrl(serviceCode string) string {
	baseUrl := os.Getenv("BASE_SERVICE_URL")
	return fmt.Sprintf("%s/connect/%s/integrations/%s/oauth2/login", baseUrl, client.ConsumerID, serviceCode)
}

func (client *ConnectClient) getCallbackUrl(serviceCode string) string {
	baseUrl := os.Getenv("BASE_SERVICE_URL")
	return fmt.Sprintf("%s/connect/%s/integrations/%s/oauth2/callback", baseUrl, client.ConsumerID, serviceCode)
}
