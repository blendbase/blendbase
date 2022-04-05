package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"blendbase/graph/generated"
	"blendbase/graph/model"
	"context"
	"errors"

	"github.com/google/uuid"
)

func (r *connectResolver) Integrations(ctx context.Context, obj *model.Connect) ([]*model.ConsumerIntegration, error) {
	connectClient, err := r.getConnectClient(ctx)
	if err != nil {
		return nil, err
	}

	if connectClient.ConsumerID == uuid.Nil {
		return nil, errors.New(MISSING_CONSUMER_ID_ERROR)
	}

	integrations, err := connectClient.ListIntegrations()
	if err != nil {
		return nil, err
	}

	return integrations, nil
}

func (r *mutationResolver) CreateConsumer(ctx context.Context) (string, error) {
	connectClient, err := r.getConnectClient(ctx)
	if err != nil {
		return "", err
	}

	consumerID, err := connectClient.CreateConsumer()
	if err != nil {
		return "", err
	}

	return consumerID.String(), nil
}

func (r *mutationResolver) EnableConsumerIntegration(ctx context.Context, serviceCode string, enabled bool) (bool, error) {
	connectClient, err := r.getConnectClient(ctx)
	if err != nil {
		return false, err
	}

	if connectClient.ConsumerID == uuid.Nil {
		return false, errors.New(MISSING_CONSUMER_ID_ERROR)
	}

	success, err := connectClient.EnableIntegration(serviceCode, enabled)
	if err != nil || !success {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SetConsumerIntegrationSecret(ctx context.Context, consumerIntegrationID string, secret string) (bool, error) {
	connectClient, err := r.getConnectClient(ctx)
	if err != nil {
		return false, err
	}

	if connectClient.ConsumerID == uuid.Nil {
		return false, errors.New(MISSING_CONSUMER_ID_ERROR)
	}

	id, err := uuid.Parse(consumerIntegrationID)
	if err != nil {
		return false, errors.New("invalid consumer integration id. must be a valid uuid")
	}

	err = connectClient.SetConsumerIntegrationSecret(id, secret)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ConfigureConsumerIntegrationOAuth(ctx context.Context, consumerIntegrationID string, input *model.OAuth2ConfigurationInput) (bool, error) {
	connectClient, err := r.getConnectClient(ctx)
	if err != nil {
		return false, err
	}

	if connectClient.ConsumerID == uuid.Nil {
		return false, errors.New(MISSING_CONSUMER_ID_ERROR)
	}

	id, err := uuid.Parse(consumerIntegrationID)
	if err != nil {
		return false, errors.New("invalid consumer integration id. must be a valid uuid")
	}

	success, err := connectClient.ConfigureOAuth2(id, input)
	if err != nil || !success {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Connect(ctx context.Context) (*model.Connect, error) {
	return &model.Connect{}, nil
}

// Connect returns generated.ConnectResolver implementation.
func (r *Resolver) Connect() generated.ConnectResolver { return &connectResolver{r} }

type connectResolver struct{ *Resolver }
