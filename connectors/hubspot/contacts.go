package hubspot

import (
	"blendbase/connectors"
	"blendbase/graph/model"
	"context"
)

type HSContact struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Archived  bool   `json:"archived"`

	Properties struct {
		Company          string `json:"company"`
		Phone            string `json:"phone"`
		Website          string `json:"website"`
		CreatedDate      string `json:"createddate"`
		Email            string `json:"email"`
		FirstName        string `json:"firstname"`
		LastName         string `json:"lastname"`
		HSObjectID       string `json:"hs_object_id"`
		LastModifiedDate string `json:"lastmodifieddate"`
	} `json:"properties"`
}

type HSContactCreateUpdatePayload struct {
	Properties struct {
		Company   *string `json:"company,omitempty"`
		FirstName *string `json:"firstname,omitempty"`
		LastName  *string `json:"lastname,omitempty"`
		Email     *string `json:"email,omitempty"`
		Phone     *string `json:"phone,omitempty"`
		Website   *string `json:"website,omitempty"`
	} `json:"properties"`
}

type HSContactsListSuccessResponse struct {
	Results []HSContact `json:"results"`
}

type HSContactUpdateSuccessResponse struct {
	ID         string    `json:"id"`
	Properties HSContact `json:"properties"`
	CreatedAt  string    `json:"createdAt"`
	UpdatedAt  string    `json:"updatedAt"`
	Archived   bool      `json:"archived"`
}

// List contacts from Hubspot API
func (client *Client) ListContacts(ctx context.Context, first int, after *string) (*model.ContactConnection, error) {
	response := HSContactsListSuccessResponse{}
	err := client.list(ctx, "contacts", first, after, connectors.StructFieldNames(HSContact{}.Properties), &response)
	if err != nil {
		return nil, err
	}

	contactEdges := make([]*model.ContactEdge, len(response.Results))

	var contact *model.Contact
	for i, hsContact := range response.Results {
		contact = hsContact.mapContactProperties()
		contactEdges[i] = &model.ContactEdge{
			Node:   contact,
			Cursor: connectors.EncodeCursor(contact.ID),
		}
	}

	recordsValue, pageInfo := client.prepareListResults(first, after, &contactEdges)

	return &model.ContactConnection{
		Edges:    recordsValue.Interface().([]*model.ContactEdge),
		PageInfo: pageInfo,
	}, nil
}

// --------------------------------------------------
// Contact:Get contact by ID
func (client *Client) GetContact(ctx context.Context, contactId string) (*model.Contact, error) {
	response := HSContact{}
	if err := client.get(ctx, "contacts", contactId, connectors.StructFieldNames(HSContact{}.Properties), &response); err != nil {
		return nil, err
	}

	contact := response.mapContactProperties()

	return contact, nil
}

// --------------------------------------------------
// Contact:Create a contact in Hubspot API
func (client *Client) CreateContact(ctx context.Context, input *model.ContactInput) (*model.Contact, error) {
	payload := createHSContactPayload(input)

	response := HSContact{}
	if err := client.create(ctx, "contacts", payload, &response); err != nil {
		return nil, err
	}

	return response.mapContactProperties(), nil
}

// --------------------------------------------------
// Contact:Update a contact in Hubspot API
func (client *Client) UpdateContact(ctx context.Context, contactId string, input *model.ContactInput) (bool, error) {
	payload := createHSContactPayload(input)

	response := HSContact{}
	if err := client.update(ctx, "contacts", contactId, payload, &response); err != nil {
		return false, err
	}

	return true, nil
}

// --------------------------------------------------
// Contact:Archive a contact in Hubspot API
func (client *Client) DeleteContact(ctx context.Context, contactId string) (bool, error) {
	response := HSContact{}
	if err := client.delete(ctx, "contacts", contactId, &response); err != nil {
		return false, err
	}

	return true, nil
}

// Creates Hubspot Contact Update/Create payload from GraphQL input
func createHSContactPayload(input *model.ContactInput) *HSContactCreateUpdatePayload {
	hsContactPayload := HSContactCreateUpdatePayload{}

	if input.FirstName != nil {
		hsContactPayload.Properties.FirstName = input.FirstName
	}

	if input.LastName != nil {
		hsContactPayload.Properties.LastName = input.LastName
	}

	if input.Email != nil {
		hsContactPayload.Properties.Email = input.Email
	}

	if input.Phone != nil {
		hsContactPayload.Properties.Phone = input.Phone
	}

	if input.Website != nil {
		hsContactPayload.Properties.Website = input.Website
	}

	if input.CompanyName != nil {
		hsContactPayload.Properties.Company = input.CompanyName
	}

	return &hsContactPayload
}

func (hsContact HSContact) mapContactProperties() *model.Contact {
	name := hsContact.Properties.FirstName + " " + hsContact.Properties.LastName

	return &model.Contact{
		ID:          hsContact.Id,
		Name:        &name,
		FirstName:   &hsContact.Properties.FirstName,
		LastName:    &hsContact.Properties.LastName,
		Email:       &hsContact.Properties.Email,
		Phone:       &hsContact.Properties.Phone,
		CompanyName: &hsContact.Properties.Company,
		Website:     &hsContact.Properties.Website,
		CreatedAt:   parseHSDateTime(&hsContact.CreatedAt),
		UpdatedAt:   parseHSDateTime(&hsContact.UpdatedAt),
		Archived:    &hsContact.Archived,
	}
}
