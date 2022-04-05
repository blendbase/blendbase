package connectors

import (
	"blendbase/graph/model"
	"context"
	"encoding/base64"
	"reflect"
)

const (
	CONNECTOR_TYPE_CRM       = "crm"
	CONNECTOR_CRM_SALESFORCE = "crm_salesforce"
	CONNECTOR_CRM_HUBSPOT    = "crm_hubspot"

	AUTH_TYPE_OAUTH2 = "oauth2"
	AUTH_TYPE_SECRET = "secret"
)

type Connector struct {
	ServiceCode string // e.g. "crm_salesforce", will be unique across all connectors
	Type        string // e.g. "crm"
	Name        string // e.g "Salesforce"
	Description string
	AuthType    string // e.g. "oauth2" or "secret"
}

var AvailableConnectors = [...]Connector{
	{
		ServiceCode: CONNECTOR_CRM_SALESFORCE,
		Type:        CONNECTOR_TYPE_CRM,
		Name:        "Salesforce",
		Description: "Salesforce is the world’s #1 customer relationship management (CRM) platform.",
		AuthType:    AUTH_TYPE_OAUTH2,
	},
	{
		ServiceCode: CONNECTOR_CRM_HUBSPOT,
		Type:        CONNECTOR_TYPE_CRM,
		Name:        "Hubspot",
		Description: "HubSpot’s CRM platform also offers enterprise software for marketing, sales, customer service, content management, and operations.",
		AuthType:    AUTH_TYPE_SECRET,
	},
}

// Takes a struct and returns a slice of its field names
func StructFieldNames(iface interface{}) []string {
	fields := make([]string, 0)
	ifv := reflect.Indirect(reflect.ValueOf(iface))
	ift := reflect.TypeOf(iface)

	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Type().Field(i)
		fields = append(fields, v.Name)
	}

	return fields
}

type CrmConnector interface {
	ListContacts(ctx context.Context, first int, after *string) (*model.ContactConnection, error)
	GetContact(ctx context.Context, contactId string) (*model.Contact, error)
	CreateContact(ctx context.Context, input *model.ContactInput) (*model.Contact, error)
	UpdateContact(ctx context.Context, contactId string, input *model.ContactInput) (bool, error)
	DeleteContact(ctx context.Context, contactId string) (bool, error)
	ListContactNotes(ctx context.Context, contactId string) ([]*model.Note, error)
	CreateContactNote(ctx context.Context, contactId string, input *model.NoteInput) (*model.Note, error)

	ListOpportunities(ctx context.Context, first int, after *string) (*model.OpportunityConnection, error)
	GetOpportunity(ctx context.Context, opportunityId string) (*model.Opportunity, error)
	CreateOpportunity(ctx context.Context, input *model.OpportunityInput) (*model.Opportunity, error)
	UpdateOpportunity(ctx context.Context, opportunityId string, input *model.OpportunityInput) (bool, error)
	DeleteOpportunity(ctx context.Context, opportunityId string) (bool, error)
	ListOpportunityNotes(ctx context.Context, opportunityId string) ([]*model.Note, error)
	CreateOpportunityNote(ctx context.Context, opportunityId string, input *model.NoteInput) (*model.Note, error)
}

func EncodeCursor(cursor string) string {
	return base64.StdEncoding.EncodeToString([]byte(cursor))
}

func DecodeCursor(cursor string) string {
	decoded, _ := base64.StdEncoding.DecodeString(cursor)
	return string(decoded)
}
