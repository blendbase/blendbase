package salesforce

import (
	"blendbase/connectors"
	"blendbase/graph/model"
	"context"

	log "github.com/sirupsen/logrus"
)

const (
	CONTACT_OBJECT = "Contact"
)

type SFContactsListSuccessResponse struct {
	SFListQuerySuccessResponseBase
	Records []SFContact `json:"records"`
}

type SFContactBase struct {
	ID          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	Title       string `json:"Title"`
	Department  string `json:"Department"`
	Email       string `json:"Email"`
	Phone       string `json:"Phone"`
	MobilePhone string `json:"MobilePhone"`

	IsDeleted        bool   `json:"IsDeleted"`
	AccountID        string `json:"AccountId"`
	CreatedDate      string `json:"CreatedDate"`
	LastModifiedDate string `json:"LastModifiedDate"`
}

type SFContact struct {
	SFContactBase
	Attributes struct {
		Type string `json:"type"`
		Url  string `json:"url"`
	} `json:"attributes"`
}

type SFContactCreateUpdatePayload struct {
	FirstName *string `json:"FirstName,omitempty"`
	LastName  *string `json:"LastName,omitempty"`
	Email     *string `json:"Email,omitempty"`
	Phone     *string `json:"Phone,omitempty"`
}

func (client *Client) ListContacts(ctx context.Context, first int, after *string) (*model.ContactConnection, error) {
	response := SFContactsListSuccessResponse{}
	err := client.list(
		CONTACT_OBJECT,
		connectors.StructFieldNames(SFContactBase{}),
		first,
		after,
		&response,
	)

	if err != nil {
		log.Errorf("Error listing contacts: %s", err)
		return nil, err
	}

	var contact *model.Contact
	contactEdges := make([]*model.ContactEdge, len(response.Records))
	for i, sfContact := range response.Records {
		contact = sfContact.mapContactProperties()
		contactEdges[i] = &model.ContactEdge{
			Cursor: connectors.EncodeCursor(contact.ID),
			Node:   contact,
		}
	}

	recordsValue, pageInfo := client.prepareListResults(first, after, &contactEdges)

	return &model.ContactConnection{
		Edges:    recordsValue.Interface().([]*model.ContactEdge),
		PageInfo: pageInfo,
	}, nil
}

func (client *Client) GetContact(ctx context.Context, contactId string) (*model.Contact, error) {
	response := SFContact{}
	err := client.get(
		CONTACT_OBJECT,
		contactId,
		connectors.StructFieldNames(SFContactBase{}),
		&response,
	)

	if err != nil {
		log.Errorf("Error getting contact: %s", err)
		return nil, err
	}

	contact := response.mapContactProperties()

	return contact, nil
}

// Create contact using GraphQL input
func (client *Client) CreateContact(ctx context.Context, input *model.ContactInput) (*model.Contact, error) {
	payload := createSFContactPayload(input)

	// Create the contact
	objectId, err := client.create(CONTACT_OBJECT, payload)
	if err != nil {
		log.Errorf("Error creating contact: %s", err)
		return nil, err
	}

	// Fetch all the fields after the creation
	contact, err := client.GetContact(ctx, objectId)
	if err != nil {
		log.Errorf("Error fetching contact #%s after creation with :%s", objectId, err)
		return nil, err
	}

	return contact, nil
}

func (client *Client) UpdateContact(ctx context.Context, contactId string, input *model.ContactInput) (bool, error) {
	payload := createSFContactPayload(input)

	// Update the contact
	success, err := client.update(CONTACT_OBJECT, contactId, payload)
	if !success || err != nil {
		log.Errorf("Error updating contact #%s: %s", contactId, err)
		return false, err
	}

	return true, nil
}

func (client *Client) DeleteContact(ctx context.Context, contactId string) (bool, error) {
	return client.delete(CONTACT_OBJECT, contactId)
}

// Creates Salesforce Contact Update/Create payload from GraphQL input
func createSFContactPayload(input *model.ContactInput) *SFContactCreateUpdatePayload {
	payload := SFContactCreateUpdatePayload{}

	if input.FirstName != nil {
		payload.FirstName = input.FirstName
	}

	if input.LastName != nil {
		payload.LastName = input.LastName
	}

	if input.Email != nil {
		payload.Email = input.Email
	}

	if input.Phone != nil {
		payload.Phone = input.Phone
	}

	return &payload
}

func (sfContact *SFContact) mapContactProperties() *model.Contact {
	name := sfContact.FirstName + " " + sfContact.LastName
	return &model.Contact{
		ID:        sfContact.ID,
		Name:      &name,
		FirstName: &sfContact.FirstName,
		LastName:  &sfContact.LastName,
		Email:     &sfContact.Email,
		Phone:     &sfContact.Phone,
		CreatedAt: parseSFDateTime(&sfContact.CreatedDate),
		UpdatedAt: parseSFDateTime(&sfContact.LastModifiedDate),
	}
}
