package connect

import (
	"os"
	"testing"

	"blendbase/config"
	"blendbase/connectors"
	"blendbase/graph/model"
	"blendbase/integrations"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"blendbase/misc/db_utils"

	log "github.com/sirupsen/logrus"
)

var (
	app           *config.App
	consumer      integrations.Consumer
	connectClient *ConnectClient
)

func TestMain(m *testing.M) {
	log.Println("Running tests!")

	godotenv.Load("../.env.test")

	app, _ = config.NewApp()
	db_utils.Migrate(app)

	if err := app.DB.Create(&consumer).Error; err != nil {
		app.Logger.Fatalf("Could not create create a test consumer: %s", err)
		return
	}

	connectClient = NewConnectClient(app, consumer.ID)

	exitVal := m.Run()

	log.Println("Cleaning up after tests!")
	app.DB.Where("1 = 1").Delete(&integrations.ConsumerOauth2Configuration{})
	app.DB.Where("1 = 1").Delete(&integrations.ConsumerIntegration{})
	// app.DB.Where("1 = 1").Delete(&integrations.Consumer{})

	os.Exit(exitVal)
}

func addCrmIntegrations(t *testing.T, consumerID uuid.UUID) []*integrations.ConsumerIntegration {
	hubspotIntegration := integrations.ConsumerIntegration{
		ConsumerID:  consumerID,
		Type:        "crm",
		ServiceCode: connectors.CONNECTOR_CRM_HUBSPOT,
	}
	app.DB.Create(&hubspotIntegration)

	salesforceIntegration := integrations.ConsumerIntegration{
		ConsumerID:  consumerID,
		Type:        "crm",
		ServiceCode: connectors.CONNECTOR_CRM_SALESFORCE,
	}
	app.DB.Create(&salesforceIntegration)

	t.Cleanup(func() {
		app.DB.Delete(&hubspotIntegration)
		app.DB.Delete(&salesforceIntegration)
	})

	return []*integrations.ConsumerIntegration{&hubspotIntegration, &salesforceIntegration}
}

// Test list integrations when non of the integrations are available
func TestListIngrationsWithNonAdded(t *testing.T) {
	integrations, err := connectClient.ListIntegrations()

	if err != nil {
		t.Errorf("Error listing integrations: %s", err)
	}

	integration := integrations[0]

	assert.Greater(t, len(integrations), 0, "There should be at least one integration")
	assert.False(t, *integration.Enabled, "The integration should not be enabled")
	assert.Nil(t, integration.ID, "The integration should not have an ID assigned")
}

func TestListIntegrationsWithSeededDB(t *testing.T) {
	addCrmIntegrations(t, consumer.ID)

	integrations, _ := connectClient.ListIntegrations()

	integration := integrations[0]

	assert.Greater(t, len(integrations), 0, "There should be at least one integration")
	assert.False(t, *integration.Enabled, "The integration should be enabled")
	assert.NotEmpty(t, *integration.ID, "The integration should have an ID assigned")
}

func TestEnableIntegrationWhenNoIntegrationsExistInTheDB(t *testing.T) {
	_, err := connectClient.EnableIntegration(connectors.CONNECTOR_CRM_HUBSPOT, false)

	assert.Nil(t, err, "There should be no error")

	var count int64
	app.DB.Model(&integrations.ConsumerIntegration{}).Where("consumer_id = ?", consumer.ID).Count(&count)

	assert.Equal(t, int64(1), count, "There should be one integration")
}

func TestEnableIntegrationWhenIntegrationsExist(t *testing.T) {
	addedIntegrations := addCrmIntegrations(t, consumer.ID) // add two integrations

	firstIntegration := addedIntegrations[0]
	// Both integrations should be disabled after creation
	_, err := connectClient.EnableIntegration(firstIntegration.ServiceCode, true)

	assert.Nil(t, err, "There should be no error")

	var count int64
	app.DB.Model(&integrations.ConsumerIntegration{}).Where("consumer_id = ?", consumer.ID).Where("enabled = ?", true).Count(&count)
	assert.Equal(t, int64(1), count, "There should be one enabled integration in the database")

	// Now enabling the second integration of the same type
	secondIntegration := addedIntegrations[1]

	_, err = connectClient.EnableIntegration(secondIntegration.ServiceCode, true)
	assert.Nil(t, err, "There should be no error")

	app.DB.Model(&integrations.ConsumerIntegration{}).Where("consumer_id = ?", consumer.ID).Where("enabled = ?", true).Count(&count)
	assert.Equal(t, int64(1), count, "There should be one enabled integration in the database")

	// The first integration should still be DISABLED
	firstIntegrationReloaded := integrations.ConsumerIntegration{}
	app.DB.Where("id = ?", firstIntegration.ID).First(&firstIntegrationReloaded)
	assert.False(t, firstIntegration.Enabled, "The first integration should not be disabled")
}

func TestConfigureOAuth2WhenProvided(t *testing.T) {
	testClientID := "test_client_id"
	testClientSecret := "test_client_secret"
	testInstanceSubdomain := "test-domain"

	addedIntegrations := addCrmIntegrations(t, consumer.ID)
	firstIntegration := addedIntegrations[0]

	oauth2Settings := model.OAuth2ConfigurationInput{
		ClientID:                    &testClientID,
		ClientSecret:                &testClientSecret,
		SalesforceInstanceSubdomain: &testInstanceSubdomain,
	}

	success, err := connectClient.ConfigureOAuth2(firstIntegration.ID, &oauth2Settings)

	assert.Nil(t, err, "There should be no error")
	assert.True(t, success, "The configuration should be successful")

	var count int64
	app.DB.Model(&integrations.ConsumerOauth2Configuration{}).Count(&count)
	assert.Equal(t, count, int64(1), "There should be one oauth2 configuration in the database")

	consumerOauth2Configuration := integrations.ConsumerOauth2Configuration{}
	app.DB.First(&consumerOauth2Configuration)

	assert.Equal(t, testClientID, consumerOauth2Configuration.ClientID.Raw, "The client ID should be the same")
	assert.Equal(t, testClientSecret, consumerOauth2Configuration.ClientSecret.Raw, "The client secret should be the same")

	app.DB.Where("1 = 1").Delete(&integrations.ConsumerOauth2Configuration{})
}

func TestConfigureOAuth2WhenMissingInput(t *testing.T) {
	testClientSecret := "test_client_secret"

	addedIntegrations := addCrmIntegrations(t, consumer.ID)
	firstIntegration := addedIntegrations[0]

	clientID := ""
	oauth2Settings := model.OAuth2ConfigurationInput{
		ClientID:     &clientID,
		ClientSecret: &testClientSecret,
	}

	success, err := connectClient.ConfigureOAuth2(firstIntegration.ID, &oauth2Settings)
	assert.NotNil(t, err, "There should be an error")
	assert.False(t, success, "The configuration should not be successful")

	oauth2Settings = model.OAuth2ConfigurationInput{}

	success, err = connectClient.ConfigureOAuth2(firstIntegration.ID, &oauth2Settings)
	assert.NotNil(t, err, "There should be an error")
	assert.False(t, success, "The configuration should not be successful")

	app.DB.Where("1 = 1").Delete(&integrations.ConsumerOauth2Configuration{})
}

func TestCreateConsumer(t *testing.T) {
	consumerID, err := connectClient.CreateConsumer()

	assert.Nil(t, err, "There should be no error")
	assert.IsType(t, consumerID, uuid.UUID{}, "The consumer ID should not be empty")
}

func TestSetConsumerIntegrationSecret(t *testing.T) {
	const secret = "super_secret"
	addedIntegrations := addCrmIntegrations(t, consumer.ID) // add two integrations
	firstIntegration := addedIntegrations[0]

	err := connectClient.SetConsumerIntegrationSecret(firstIntegration.ID, secret)
	assert.Nil(t, err, "There should be no error")

	integration := integrations.ConsumerIntegration{}
	app.DB.Where("id = ?", firstIntegration.ID).First(&integration)

	assert.Equal(t, integration.Secret.Raw, secret, "The secret in the DB should match the one provided")
}
