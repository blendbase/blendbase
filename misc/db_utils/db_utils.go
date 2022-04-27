package db_utils

import (
	"blendbase/config"
	"blendbase/connectors"
	"blendbase/integrations"
	"blendbase/misc/gormext"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	// this UUID is used in the Salesforce callback
	TEST_CONSUMER_ID = "c6a82fd9-7e22-40c2-8bf2-db58a40839a9"
)

// Migrates the database based on the lastest version of the models
// see https://gorm.io/docs/migration.html
func Migrate(app *config.App) error {
	err := app.DB.AutoMigrate(
		&integrations.Consumer{},
		&integrations.ConsumerIntegration{},
		&integrations.ConsumerOauth2Configuration{},
	)

	if err != nil {
		return err
	}

	app.Logger.Info("Database migrated")

	return nil
}

// Creates seed data for the database.
// This is only for development and testing
func Seed(app *config.App) error {
	consumer := integrations.Consumer{}

	result := app.DB.First(&consumer, "id = ?", TEST_CONSUMER_ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Fatalf("failed to seed data: %+v", result.Error)
		return result.Error
	}

	if result.Error == nil {
		log.Println("Seed data already exists")
		return nil
	}

	consumer = integrations.Consumer{Base: integrations.Base{ID: uuid.MustParse(TEST_CONSUMER_ID)}}
	app.DB.Where("1 = 1").Delete(&integrations.ConsumerOauth2Configuration{})

	if err := app.DB.Create(&consumer).Error; err != nil {
		log.Fatal("failed to create consumer")
		return err
	}

	hubspotIntegration := integrations.ConsumerIntegration{
		ConsumerID:  consumer.ID,
		Type:        connectors.CONNECTOR_TYPE_CRM,
		ServiceCode: connectors.CONNECTOR_CRM_HUBSPOT,
		Secret: gormext.EncryptedValue{
			Raw: os.Getenv("HUBSPOT_ACCESS_TOKEN"),
		},
	}
	if err := app.DB.Create(&hubspotIntegration).Error; err != nil {
		log.Fatalf("failed to create hubspot integration: %+v", err)
		return err
	}

	salesforceIntegration := integrations.ConsumerIntegration{
		ConsumerID:  consumer.ID,
		Type:        connectors.CONNECTOR_TYPE_CRM,
		ServiceCode: connectors.CONNECTOR_CRM_SALESFORCE,
	}

	if err := app.DB.Create(&salesforceIntegration).Error; err != nil {
		log.Fatalf("failed to create salesforce integration: %+v", err)
		return err
	}

	oauth2Configuration := integrations.ConsumerOauth2Configuration{
		ConsumerIntegrationID: salesforceIntegration.ID,
		ClientID:              gormext.EncryptedValue{Raw: os.Getenv("SALESFORCE_CLIENT_ID")},
		ClientSecret:          gormext.EncryptedValue{Raw: os.Getenv("SALESFORCE_CLIENT_SECRET")},
		RedirectURL:           fmt.Sprintf("%s/connect/%s/integrations/crm_salesforce/oauth2/callback", os.Getenv("BASE_SERVICE_URL"), consumer.ID),
	}

	if err := app.DB.Create(&oauth2Configuration).Error; err != nil {
		log.Fatalf("failed to create salesforce oauth2 configuration: %+v", err)
		return err
	}

	return nil
}
