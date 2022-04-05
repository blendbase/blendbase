package salesforce

import (
	"blendbase/connectors"
	"blendbase/graph/model"
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	NOTE_OBJECT          = "Note"
	SF_NOTE_TITLE_LENGTH = 30
)

type SFNotesListSuccessResponse struct {
	SFListQuerySuccessResponseBase
	Records []SFNote `json:"records"`
}

// https://developer.salesforce.com/docs/atlas.en-us.object_reference.meta/object_reference/sforce_api_objects_note.htm
type SFNote struct {
	ID               string `json:"Id"`
	Body             string `json:"Body"`
	IsDeleted        bool   `json:"IsDeleted"`
	IsPrivate        bool   `json:"IsPrivate"`
	OwnerId          string `json:"OwnerId"`
	ParentId         string `json:"ParentId"`
	Title            string `json:"Title"`
	CreatedDate      string `json:"CreatedDate"`
	LastModifiedDate string `json:"LastModifiedDate"`
}

type SFNoteCreateUpdatePayload struct {
	Body     string `json:"Body"`
	Title    string `json:"Title"`
	ParentId string `json:"ParentId"`
}

func (client *Client) ListContactNotes(ctx context.Context, contactId string) ([]*model.Note, error) {
	return client.listNotes(contactId)
}

func (client *Client) CreateContactNote(ctx context.Context, contactId string, input *model.NoteInput) (*model.Note, error) {
	return client.createNote(contactId, input)
}

func (client *Client) ListOpportunityNotes(ctx context.Context, opportunityId string) ([]*model.Note, error) {
	return client.listNotes(opportunityId)
}

func (client *Client) CreateOpportunityNote(ctx context.Context, opportunityId string, input *model.NoteInput) (*model.Note, error) {
	return client.createNote(opportunityId, input)
}

func (client *Client) listNotes(parentId string) ([]*model.Note, error) {
	response := SFNotesListSuccessResponse{}
	err := client.listWithWhere("Note", connectors.StructFieldNames(SFNote{}),
		fmt.Sprintf("ParentId = '%s'", parentId), &response)

	if err != nil {
		return nil, err
	}

	notes := make([]*model.Note, len(response.Records))
	for i, sfNote := range response.Records {
		notes[i] = sfNote.mapNoteProperties()
	}

	return notes, err
}

func (client *Client) createNote(parentId string, input *model.NoteInput) (*model.Note, error) {
	payload := createSFNotePayload(parentId, input)

	// Create the contact
	objectId, err := client.create(NOTE_OBJECT, payload)
	if err != nil {
		log.Errorf("Error creating note: %s", err)
		return nil, err
	}

	// Fetch all the fields after the creation
	sfNote := SFNote{}
	err = client.get(NOTE_OBJECT, objectId, connectors.StructFieldNames(sfNote), &sfNote)
	if err != nil {
		log.Errorf("Error fetching note #%s after creation with :%s", objectId, err)
		return nil, err
	}

	note := sfNote.mapNoteProperties()
	return note, nil
}

func (sfNote *SFNote) mapNoteProperties() *model.Note {
	return &model.Note{
		ID:        sfNote.ID,
		CreatedAt: parseSFDateTime(&sfNote.CreatedDate),
		UpdatedAt: parseSFDateTime(&sfNote.LastModifiedDate),
		Content:   sfNote.Body,
	}
}

func createSFNotePayload(parentId string, input *model.NoteInput) *SFNoteCreateUpdatePayload {
	titleSliceLength := SF_NOTE_TITLE_LENGTH
	if len(input.Content) < titleSliceLength {
		titleSliceLength = len(input.Content)
	}

	payload := SFNoteCreateUpdatePayload{
		Body:     input.Content,
		Title:    input.Content[:titleSliceLength] + "...",
		ParentId: parentId,
	}

	return &payload
}
