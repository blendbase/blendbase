package hubspot

import (
	"blendbase/connectors"
	"blendbase/graph/model"
	"context"
	"time"
)

type HSDeal struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Archived  bool   `json:"archived"`

	Properties struct {
		Amount         *string `json:"amount"`
		CloseDate      *string `json:"closedate"`
		DealName       string  `json:"dealname"`
		DealStage      *string `json:"dealstage"`
		HubspotOwnerId *string `json:"hubspot_owner_id"`
		Pipeline       *string `json:"pipeline"`
	} `json:"properties"`
}

type HSDealsListSuccessResponse struct {
	Results []HSDeal `json:"results"`
}

func (hsDeal *HSDeal) mapOpportunityProperties() *model.Opportunity {
	var closeDatePtr *time.Time = nil

	if hsDeal.Properties.CloseDate != nil {
		closeDate, _ := time.Parse(time.RFC3339, *hsDeal.Properties.CloseDate)
		closeDatePtr = &closeDate
	}

	return &model.Opportunity{
		ID:        hsDeal.Id,
		Name:      hsDeal.Properties.DealName,
		StageName: hsDeal.Properties.DealStage,
		CloseDate: closeDatePtr,
		Amount:    hsDeal.Properties.Amount,
	}
}

type HSDealCreateUpdatePayload struct {
	Properties struct {
		Amount         *string `json:"amount,omitempty"`
		CloseDate      *string `json:"closedate"`
		DealName       string  `json:"dealname"`
		DealStage      *string `json:"dealstage"`
		HubspotOwnerId *string `json:"hubspot_owner_id"`
		Pipeline       *string `json:"pipeline"`
	} `json:"properties"`
}

func (client *Client) ListOpportunities(ctx context.Context, first int, after *string) (*model.OpportunityConnection, error) {
	response := HSDealsListSuccessResponse{}
	err := client.list(ctx, "deals", first, after, connectors.StructFieldNames(HSContact{}.Properties), &response)
	if err != nil {
		return nil, err
	}

	opportunityEdges := make([]*model.OpportunityEdge, len(response.Results))

	var opportunity *model.Opportunity
	for i, hsOpportunity := range response.Results {
		opportunity = hsOpportunity.mapOpportunityProperties()
		opportunityEdges[i] = &model.OpportunityEdge{
			Node:   opportunity,
			Cursor: connectors.EncodeCursor(opportunity.ID),
		}
	}

	recordsValue, pageInfo := client.prepareListResults(first, after, &opportunityEdges)

	return &model.OpportunityConnection{
		Edges:    recordsValue.Interface().([]*model.OpportunityEdge),
		PageInfo: pageInfo,
	}, nil
}

func (client *Client) GetOpportunity(ctx context.Context, opportunityId string) (*model.Opportunity, error) {
	deal := &HSDeal{}
	if err := client.get(ctx, "deals", opportunityId, connectors.StructFieldNames(HSDeal{}.Properties), deal); err != nil {
		return nil, err
	}

	opportunity := deal.mapOpportunityProperties()
	return opportunity, nil
}

func (client *Client) CreateOpportunity(ctx context.Context, input *model.OpportunityInput) (*model.Opportunity, error) {
	payload := createHSDealPayload(input)

	deal := HSDeal{}
	if err := client.create(ctx, "deals", payload, &deal); err != nil {
		return nil, err
	}

	return deal.mapOpportunityProperties(), nil
}

func (client *Client) UpdateOpportunity(ctx context.Context, opportunityId string, input *model.OpportunityInput) (bool, error) {
	payload := createHSDealPayload(input)

	deal := HSDeal{}
	if err := client.update(ctx, "deals", opportunityId, payload, &deal); err != nil {
		return false, err
	}

	return true, nil
}

func (client *Client) DeleteOpportunity(ctx context.Context, contactId string) (bool, error) {
	response := HSDeal{}
	if err := client.delete(ctx, "deals", contactId, &response); err != nil {
		return false, err
	}

	return true, nil
}

func createHSDealPayload(input *model.OpportunityInput) *HSDealCreateUpdatePayload {
	payload := HSDealCreateUpdatePayload{}

	payload.Properties.DealName = input.Name
	payload.Properties.DealStage = &input.StageName
	formattedDate := input.CloseDate.Format(time.RFC3339)
	payload.Properties.CloseDate = &formattedDate
	payload.Properties.Amount = input.Amount

	return &payload
}
