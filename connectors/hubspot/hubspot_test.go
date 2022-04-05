package hubspot

import (
	"blendbase/misc/test_utils"
	"context"
	"encoding/base64"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestListContacts(t *testing.T) {
	godotenv.Load("../../.env")
	c := HubspotClient(os.Getenv("HUBSPOT_ACCESS_TOKEN"))

	ctx := context.Background()
	contactConnection, err := c.ListContacts(ctx, 10, nil)

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, contactConnection.Edges, "expecting non-nil contacts")
	assert.Greater(t, len(contactConnection.Edges), 0, "expecting more than zero contacts")
	assert.NotEmpty(t, contactConnection.Edges[0].Node.ID, 0, "expecting a non-empty ID for the first contact")
}

func TestListContactsPagination(t *testing.T) {
	godotenv.Load("../../.env")
	c := HubspotClient(os.Getenv("HUBSPOT_ACCESS_TOKEN"))

	ctx := context.Background()
	contactConnection, err := c.ListContacts(ctx, 1, nil)

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, contactConnection.Edges, "expecting non-nil result")
	assert.Equal(t, len(contactConnection.Edges), 1, "expecting single result")
	assert.NotEmpty(t, contactConnection.Edges[0].Node.ID, 0, "expecting a non-empty ID for the first contact")
	assert.Equal(t,
		contactConnection.Edges[0].Cursor,
		base64.StdEncoding.EncodeToString([]byte(contactConnection.Edges[0].Node.ID)),
		"expecting cursor to be equal to the encoded ID of the first contact")

	afterParam := contactConnection.PageInfo.EndCursor
	decodedAfterParam, _ := base64.StdEncoding.DecodeString(*afterParam)
	assert.Equal(t,
		string(decodedAfterParam),
		contactConnection.Edges[0].Node.ID,
		"expecting EndCursor to be equal to the last contact ID")

	contactConnection, err = c.ListContacts(ctx, 1, afterParam)

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, contactConnection.Edges, "expecting non-nil result")
	assert.Equal(t, len(contactConnection.Edges), 1, "expecting single result")
	assert.NotEqual(t, contactConnection.Edges[0].Node.ID, afterParam, "expecting a different ID")
}

func TestGetContact(t *testing.T) {
	godotenv.Load("../../.env")
	c := HubspotClient(os.Getenv("HUBSPOT_ACCESS_TOKEN"))

	ctx := context.Background()

	const contactId = "1"
	contact, err := c.GetContact(ctx, contactId)

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, contact, "expecting non-nil result")
	assert.Equal(t, contact.ID, contactId, "expecting a ID for the contact equal to the ID requested")
	assert.NotEmpty(t, contact.FirstName, "expecting non-nil first name")
	assert.NotEmpty(t, contact.LastName, "expecting non-nil last name")
	assert.NotEmpty(t, contact.Email, "expecting non-nil email")
	assert.NotEmpty(t, contact.Phone, "expecting non-nil phone")
	assert.NotEmpty(t, contact.Website, "expecting non-nil website")
	assert.NotEmpty(t, contact.CompanyName, "expecting non-nil company name")
}

func TestCreateContact(t *testing.T) {
	godotenv.Load("../../.env")
	c := HubspotClient(os.Getenv("HUBSPOT_ACCESS_TOKEN"))

	ctx := context.Background()
	input := test_utils.GenerateContactInput()

	contact, err := c.CreateContact(ctx, input)
	assert.Nil(t, err, "expecting nil error")
	assert.NotEmpty(t, contact.ID, "expecting a non-empty ID for the contact")
}

func TestUpdateContact(t *testing.T) {
	godotenv.Load("../../.env")
	c := HubspotClient(os.Getenv("HUBSPOT_ACCESS_TOKEN"))

	ctx := context.Background()
	const contactId = "1"
	input := test_utils.GenerateContactInput()

	success, err := c.UpdateContact(ctx, contactId, input)
	assert.Nil(t, err, "expecting nil error")
	assert.True(t, success, "expecting true update result")

	contact, err := c.GetContact(ctx, contactId)
	assert.Nil(t, err, "expecting nil error")
	assert.Equal(t, contact.CompanyName, input.CompanyName, "expecting a company name for the contact equal to the company name requested")
	assert.Equal(t, contact.FirstName, input.FirstName, "expecting a first name for the contact equal to the first name requested")
	assert.Equal(t, contact.LastName, input.LastName, "expecting a last name for the contact equal to the last name requested")
	assert.Equal(t, contact.Email, input.Email, "expecting a email for the contact equal to the email requested")
	assert.Equal(t, contact.Phone, input.Phone, "expecting a phone for the contact equal to the phone requested")
	assert.Contains(t, *contact.Website, *input.Website, "expecting a website for the contact equal to the website requested")
}

func TestDeleteContact(t *testing.T) {
	godotenv.Load("../../.env")
	c := HubspotClient(os.Getenv("HUBSPOT_ACCESS_TOKEN"))

	ctx := context.Background()
	contact, err := c.CreateContact(ctx, test_utils.GenerateContactInput())
	contactId := contact.ID

	assert.Nil(t, err, "expecting nil error")
	assert.NotEmpty(t, contactId, "expecting a non-empty ID for the contact")

	success, err := c.DeleteContact(ctx, contactId)
	assert.Nil(t, err, "expecting nil error")
	assert.True(t, success, "expecting true archive result")
}

func TestCreateContactAndAddNote(t *testing.T) {
	godotenv.Load("../../.env")
	c := HubspotClient(os.Getenv("HUBSPOT_ACCESS_TOKEN"))

	ctx := context.Background()
	input := test_utils.GenerateContactInput()

	contact, err := c.CreateContact(ctx, input)
	assert.Nil(t, err, "expecting nil error")
	assert.NotEmpty(t, contact.ID, "expecting a non-empty ID for the contact")

	noteInput := test_utils.GenerateNoteInput()
	note, err := c.CreateContactNote(ctx, contact.ID, noteInput)

	assert.Nil(t, err, "expecting nil error")
	assert.NotEmpty(t, note.ID, "expecting a non-empty ID for the note")
	assert.Equal(t, note.Content, noteInput.Content, "expecting a content for the note equal to the content requested")
}

func TestListOpportunities(t *testing.T) {
	godotenv.Load("../../.env")
	c := HubspotClient(os.Getenv("HUBSPOT_ACCESS_TOKEN"))

	ctx := context.Background()
	opportunityConnection, err := c.ListOpportunities(ctx, 10, nil)

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, opportunityConnection.Edges, "expecting non-nil opportunities")
	assert.Greater(t, len(opportunityConnection.Edges), 0, "expecting more than zero opportunities")
	assert.NotEmpty(t, opportunityConnection.Edges[0].Node.ID, 0, "expecting a non-empty ID for the first opportunity")
}

func TestOpportunityCRUD(t *testing.T) {
	godotenv.Load("../../.env")
	c := HubspotClient(os.Getenv("HUBSPOT_ACCESS_TOKEN"))

	ctx := context.Background()
	input := test_utils.GenerateOpportunityInput()

	opportunity, err := c.CreateOpportunity(ctx, input)
	assert.Nil(t, err, "expecting nil error", err)
	assert.NotEmpty(t, opportunity.ID, "expecting a non-empty ID for the opportunity")

	foundOpportunity, err := c.GetOpportunity(ctx, opportunity.ID)
	assert.Nil(t, err, "expecting nil error")
	assert.Equal(t, opportunity.ID, foundOpportunity.ID, "expecting a ID for the opportunity equal to the ID requested")
	assert.Equal(t, opportunity.Name, foundOpportunity.Name, "expecting a name for the opportunity equal to the created one")

	input = test_utils.GenerateOpportunityInput()

	success, err := c.UpdateOpportunity(ctx, opportunity.ID, input)
	assert.Nil(t, err, "expecting nil error")
	assert.True(t, success, "expecting true update result")

	success, err = c.DeleteOpportunity(ctx, opportunity.ID)
	assert.Nil(t, err, "expecting nil error")
	assert.True(t, success, "expecting true archive result")
}

func TestCreateOpportunityAndAddNote(t *testing.T) {
	godotenv.Load("../../.env")
	c := HubspotClient(os.Getenv("HUBSPOT_ACCESS_TOKEN"))

	ctx := context.Background()
	input := test_utils.GenerateOpportunityInput()

	opportunity, err := c.CreateOpportunity(ctx, input)
	assert.Nil(t, err, "expecting nil error")
	assert.NotEmpty(t, opportunity.ID, "expecting a non-empty ID for the opportunity")

	noteInput := test_utils.GenerateNoteInput()
	note, err := c.CreateOpportunityNote(ctx, opportunity.ID, noteInput)

	assert.Nil(t, err, "expecting nil error")
	assert.NotEmpty(t, note.ID, "expecting a non-empty ID for the note")
	assert.Equal(t, note.Content, noteInput.Content, "expecting a content for the note equal to the content requested")
}
