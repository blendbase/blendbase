package salesforce

import (
	"context"
	"encoding/base64"
	"testing"

	"blendbase/config"
	"blendbase/integrations"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"blendbase/misc/db_utils"
	"blendbase/misc/test_utils"

	log "github.com/sirupsen/logrus"
)

var (
	client *Client
)

const (
	contactId = "0035f00000AHo1uAAD"
)

func init() {
	godotenv.Load("../../.env.test")

	app, _ := config.NewApp()

	db_utils.Migrate(app)
	// THIS IS A TEMPORARY FIX FOR THE TESTING PURPOSES
	// you would need to load the app and get access token for the first time only,
	// after that the app will be able to refresh the token with the refresh token from the database

	consumer := integrations.Consumer{}
	if err := app.DB.Where("id = ?", db_utils.TEST_CONSUMER_ID).Order("created_at asc").First(&consumer).Error; err != nil {
		log.Fatalf("Could not find at lest one consumer after seeding the DB: %s", err)
		return
	}
	log.Debugf("Consumer %+#v", consumer)

	var err error
	client, err = LoadClientFromDB(app, &consumer)
	if err != nil {
		log.Fatalf("Could not load client from DB: %s", err)
	}

	gofakeit.Seed(0)
}

func TestListContacts(t *testing.T) {
	ctx := context.Background()
	contactConnection, err := client.ListContacts(ctx, 10, nil)

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, contactConnection, "expecting non-nil result")
	assert.Greater(t, len(contactConnection.Edges), 0, "expecting more than zero contacts")

	contact := *contactConnection.Edges[0].Node

	assert.NotEmpty(t, contact.ID, 0, "expecting a non-empty ID for the first contact")
}

func TestListContactsPagination(t *testing.T) {
	ctx := context.Background()
	contactConnection, err := client.ListContacts(ctx, 1, nil)

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

	contactConnection, err = client.ListContacts(ctx, 1, afterParam)

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, contactConnection.Edges, "expecting non-nil result")
	assert.Equal(t, len(contactConnection.Edges), 1, "expecting single result")
	assert.NotEqual(t, contactConnection.Edges[0].Node.ID, afterParam, "expecting a different ID")
}

func TestGetContact(t *testing.T) {
	ctx := context.Background()
	const contactId = "0035f00000AHo1uAAD"
	contact, err := client.GetContact(ctx, contactId)

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, contact, "expecting non-nil result")
	assert.Equal(t, contact.ID, contactId, "expecting a ID for the contact equal to the ID requested")
	assert.NotEmpty(t, *contact.Email, "expecting non-nil email")
	assert.NotEmpty(t, *contact.Phone, "expecting non-nil phone")
}

func TestCreateContact(t *testing.T) {
	ctx := context.Background()
	input := test_utils.GenerateContactInput()

	contact, err := client.CreateContact(ctx, input)
	log.Debugf("Contact created: %+v\n", contact)
	assert.Nil(t, err, "expecting nil error")
	assert.NotEmpty(t, contact.ID, "expecting a non-empty ID for the contact")
}

func TestUpdateContact(t *testing.T) {
	ctx := context.Background()
	input := test_utils.GenerateContactInput()

	success, err := client.UpdateContact(ctx, contactId, input)
	assert.Nil(t, err, "expecting nil error")
	assert.True(t, success, "expecting success as a result of the update process")

	contact, err := client.GetContact(ctx, contactId)
	assert.Nil(t, err, "expecting nil error while retreiving the contact")

	assert.Equal(t, contactId, contact.ID, "expecting a ID for the contact equal to the ID requested")
	assert.Equal(t, *input.FirstName, *contact.FirstName, "expecting a last name for the contact equal to the last name requested")
	assert.Equal(t, *input.LastName, *contact.LastName, "expecting a last name for the contact equal to the last name requested")
	assert.Equal(t, *input.Email, *contact.Email, "expecting a email for the contact equal to the email requested")
	assert.Equal(t, *input.Phone, *contact.Phone, "expecting a phone for the contact equal to the phone requested")
}

func TestDeleteContact(t *testing.T) {
	ctx := context.Background()
	contact, err := client.CreateContact(ctx, test_utils.GenerateContactInput())
	contactId := contact.ID

	assert.Nil(t, err, "expecting nil error")
	assert.NotEmpty(t, contactId, "expecting a non-empty ID for the contact")

	success, err := client.DeleteContact(ctx, contactId)
	assert.Nil(t, err, "expecting nil error")
	assert.True(t, success, "expecting true archive result")
}

func TestListOpportunitiesPagination(t *testing.T) {
	ctx := context.Background()
	connection, err := client.ListOpportunities(ctx, 1, nil)

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, connection.Edges, "expecting non-nil result")
	assert.Equal(t, len(connection.Edges), 1, "expecting single result")
	assert.NotEmpty(t, connection.Edges[0].Node.ID, 0, "expecting a non-empty ID for the first record")
	assert.Equal(t,
		connection.Edges[0].Cursor,
		base64.StdEncoding.EncodeToString([]byte(connection.Edges[0].Node.ID)),
		"expecting cursor to be equal to the encoded ID of the first opportunity")

	afterParam := connection.PageInfo.EndCursor
	decodedAfterParam, _ := base64.StdEncoding.DecodeString(*afterParam)
	assert.Equal(t,
		string(decodedAfterParam),
		connection.Edges[0].Node.ID,
		"expecting EndCursor to be equal to the last ID")

	connection, err = client.ListOpportunities(ctx, 1, afterParam)

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, connection.Edges, "expecting non-nil result")
	assert.Equal(t, len(connection.Edges), 1, "expecting single result")
	assert.NotEqual(t, connection.Edges[0].Node.ID, afterParam, "expecting a different ID")
}

func TestOpportunityCRUD(t *testing.T) {
	ctx := context.Background()
	input := test_utils.GenerateOpportunityInput()

	opportunity, err := client.CreateOpportunity(ctx, input)
	assert.Nil(t, err, "expecting nil error")
	assert.NotEmpty(t, opportunity.ID, "expecting a non-empty ID for the opportunity")

	foundOpportunity, err := client.GetOpportunity(ctx, opportunity.ID)
	assert.Nil(t, err, "expecting nil error")
	assert.Equal(t, opportunity.ID, foundOpportunity.ID, "expecting a ID for the opportunity equal to the ID requested")
	assert.Equal(t, opportunity.Name, foundOpportunity.Name, "expecting a name for the opportunity equal to the created one")

	input = test_utils.GenerateOpportunityInput()

	success, err := client.UpdateOpportunity(ctx, opportunity.ID, input)
	assert.Nil(t, err, "expecting nil error")
	assert.True(t, success, "expecting true update result")

	success, err = client.DeleteOpportunity(ctx, opportunity.ID)
	assert.Nil(t, err, "expecting nil error")
	assert.True(t, success, "expecting true archive result")
}

func TestCreateOpportunityAndAddNote(t *testing.T) {
	ctx := context.Background()
	input := test_utils.GenerateOpportunityInput()

	opportunity, err := client.CreateOpportunity(ctx, input)
	assert.Nil(t, err, "expecting nil error")
	assert.NotEmpty(t, opportunity.ID, "expecting a non-empty ID for the opportunity")

	noteInput := test_utils.GenerateNoteInput()
	note, err := client.CreateOpportunityNote(ctx, opportunity.ID, noteInput)

	assert.Nil(t, err, "expecting nil error")
	assert.NotEmpty(t, note.ID, "expecting a non-empty ID for the note")
	assert.Equal(t, note.Content, noteInput.Content, "expecting a content for the note equal to the content requested")
}
