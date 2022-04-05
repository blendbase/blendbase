package graph

import (
	"blendbase/config"
	"blendbase/connect"
	"blendbase/connectors"
	"blendbase/connectors/hubspot"
	"blendbase/connectors/salesforce"
	"blendbase/graph/auth"
	"blendbase/integrations"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

const MISSING_CONSUMER_ID_ERROR = "missing consumer ID. please provide consumer ID in order to use this endpoint"

type Resolver struct {
	App       *config.App
	GraphAuth *auth.GraphAuth
}

func (r *Resolver) getCrmConsumerIntegration(ctx context.Context) (*integrations.ConsumerIntegration, error) {
	consumerID := r.GraphAuth.GetConsumerIDFromContext(ctx)
	if consumerID == nil {
		return nil, fmt.Errorf("missing consumer ID")
	}

	if err := r.checkConsumerExistence(consumerID); err != nil {
		return nil, err
	}

	integration := integrations.ConsumerIntegration{}
	// TODO: add support for users having multiple integrations
	if err := r.App.DB.Where("consumer_id = ?", *consumerID).Where("type = ?", "crm").Where("enabled = ?", true).First(&integration).Error; err != nil {
		return nil, fmt.Errorf("crm integration not found")
	}

	return &integration, nil
}

func (r *Resolver) getOAuthConfig(consumerIntegration *integrations.ConsumerIntegration) *integrations.ConsumerOauth2Configuration {
	var consumerOAuthConfig integrations.ConsumerOauth2Configuration
	if err := r.App.DB.Where("consumer_integration_id = ?", consumerIntegration.ID).First(&consumerOAuthConfig).Error; err != nil {
		return nil
	}
	return &consumerOAuthConfig
}

func (r *Resolver) getCrmConnector(ctx context.Context) (connectors.CrmConnector, error) {
	integration, err := r.getCrmConsumerIntegration(ctx)
	if err != nil {
		return nil, err
	}

	if integration.ServiceCode == connectors.CONNECTOR_CRM_HUBSPOT {
		client := hubspot.HubspotClient(integration.Secret.Raw)

		return client, nil
	} else if integration.ServiceCode == connectors.CONNECTOR_CRM_SALESFORCE {
		oauthConfig := r.getOAuthConfig(integration)

		return salesforce.SaleforceClient(r.App, oauthConfig), nil
	}

	return nil, fmt.Errorf("crm integration not found")
}

func (r *Resolver) getConnectClient(ctx context.Context) (*connect.ConnectClient, error) {
	consumerID := r.GraphAuth.GetConsumerIDFromContext(ctx)

	if consumerID == nil {
		return connect.NewConnectClient(r.App, uuid.Nil), nil
	} else {
		if err := r.checkConsumerExistence(consumerID); err != nil {
			return nil, errors.New("consumer does not exist")
		}

		return connect.NewConnectClient(r.App, *consumerID), nil
	}
}

func (r *Resolver) checkConsumerExistence(consumerID *uuid.UUID) error {
	consumer := integrations.Consumer{}
	if err := r.App.DB.Where("id = ?", *consumerID).First(&consumer).Error; err != nil {
		return fmt.Errorf("unable to find the consumer; please provide a valid consumer ID")
	}

	return nil
}
