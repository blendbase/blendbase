package salesforce

import (
	"blendbase/connectors"
	"blendbase/graph/model"
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	OPPORTUNITY_OBJECT = "Opportunity"
)

type SFOpportunity struct {
	Id        string   `json:"Id"`
	Name      string   `json:"Name"`
	StageName string   `json:"StageName"`
	CloseDate *string  `json:"CloseDate"`
	Amount    *float32 `json:"Amount"`
}

type SFOpportunityListSuccessResponse struct {
	SFListQuerySuccessResponseBase
	Records []SFOpportunity `json:"records"`
}

type SFOpportunityCreateUpdatePayload struct {
	Name      string  `json:"Name"`
	StageName string  `json:"StageName"`
	CloseDate string  `json:"CloseDate"`
	Amount    *string `json:"Amount"`
}

func (client *Client) ListOpportunities(ctx context.Context, first int, after *string) (*model.OpportunityConnection, error) {
	response := SFOpportunityListSuccessResponse{}
	err := client.list(
		OPPORTUNITY_OBJECT,
		connectors.StructFieldNames(SFOpportunity{}),
		first,
		after,
		&response,
	)

	if err != nil {
		log.Errorf("Error listing opportunities: %s", err)
		return nil, err
	}

	var opportunity *model.Opportunity
	edges := make([]*model.OpportunityEdge, len(response.Records))
	for i, sfOpportunity := range response.Records {
		opportunity = sfOpportunity.mapProperties()
		edges[i] = &model.OpportunityEdge{
			Cursor: connectors.EncodeCursor(opportunity.ID),
			Node:   opportunity,
		}
	}

	recordsValue, pageInfo := client.prepareListResults(first, after, &edges)

	return &model.OpportunityConnection{
		Edges:    recordsValue.Interface().([]*model.OpportunityEdge),
		PageInfo: pageInfo,
	}, nil
}

func (client *Client) GetOpportunity(ctx context.Context, opportunityId string) (*model.Opportunity, error) {
	response := SFOpportunity{}
	err := client.get(
		OPPORTUNITY_OBJECT,
		opportunityId,
		connectors.StructFieldNames(SFOpportunity{}),
		&response,
	)

	if err != nil {
		log.Errorf("Error getting opportunity: %s", err)
		return nil, err
	}

	obj := response.mapProperties()
	return obj, nil
}

func (client *Client) CreateOpportunity(ctx context.Context, input *model.OpportunityInput) (*model.Opportunity, error) {
	payload := createSFOpportunityPayload(input)

	objectId, err := client.create(OPPORTUNITY_OBJECT, payload)
	if err != nil {
		log.Errorf("Error creating opportunity: %s", err)
		return nil, err
	}

	// Fetch all the fields after the creation
	opportunity, err := client.GetOpportunity(ctx, objectId)
	if err != nil {
		log.Errorf("Error fetching opportunity #%s after creation with :%s", objectId, err)
		return nil, err
	}

	return opportunity, nil
}

func (client *Client) UpdateOpportunity(ctx context.Context, opportunityId string, input *model.OpportunityInput) (bool, error) {
	payload := createSFOpportunityPayload(input)

	success, err := client.update(OPPORTUNITY_OBJECT, opportunityId, payload)
	if !success || err != nil {
		log.Errorf("Error updating opportunity #%s: %s", opportunityId, err)
		return false, err
	}

	return true, nil
}

func (client *Client) DeleteOpportunity(ctx context.Context, opportunityId string) (bool, error) {
	return client.delete(OPPORTUNITY_OBJECT, opportunityId)
}

func (sfOpportunity *SFOpportunity) mapProperties() *model.Opportunity {
	opportunity := model.Opportunity{
		ID:        sfOpportunity.Id,
		Name:      sfOpportunity.Name,
		StageName: &sfOpportunity.StageName,
	}

	if sfOpportunity.CloseDate != nil {
		closeDate, _ := time.Parse(time.RFC3339, *sfOpportunity.CloseDate)
		opportunity.CloseDate = &closeDate
	}

	if sfOpportunity.Amount != nil {
		amount := fmt.Sprintf("%f", *sfOpportunity.Amount)
		opportunity.Amount = &amount
	}

	return &opportunity
}

func createSFOpportunityPayload(input *model.OpportunityInput) *SFOpportunityCreateUpdatePayload {
	payload := SFOpportunityCreateUpdatePayload{}

	payload.Name = input.Name
	payload.StageName = input.StageName
	payload.CloseDate = formatSFDateTime(input.CloseDate)
	if input.Amount != nil {
		payload.Amount = input.Amount
	}

	return &payload
}
