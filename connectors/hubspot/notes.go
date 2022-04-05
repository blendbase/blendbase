package hubspot

import (
	"blendbase/graph/model"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type HSNote struct {
	Id        string  `json:"id"`
	CreatedAt *string `json:"createdAt"`
	UpdatedAt *string `json:"updatedAt"`
	Archived  bool    `json:"archived"`

	Properties struct {
		HsCreatedDate  *string `json:"hs_createdate,omitempty"`
		HsNoteBody     *string `json:"hs_note_body,omitempty"` // can contain HTML
		HubspotOwnerId *string `json:"hubspot_owner_id,omitempty"`
	} `json:"properties"`
}

type HSNoteCreateUpdatePayload struct {
	Properties struct {
		HsNoteBody  *string `json:"hs_note_body,omitempty"`
		HsTimestamp *string `json:"hs_timestamp,omitempty"`
	} `json:"properties"`
}

func (client *Client) ListContactNotes(ctx context.Context, contactId string) ([]*model.Note, error) {
	return client.listNotes(ctx, "contacts", contactId)
}

func (client *Client) CreateContactNote(ctx context.Context, contactId string, input *model.NoteInput) (*model.Note, error) {
	return client.createNoteAndAssociate(ctx, "contacts", contactId, "note_to_contact", input)
}

func (client *Client) ListOpportunityNotes(ctx context.Context, opportunityId string) ([]*model.Note, error) {
	return client.listNotes(ctx, "deals", opportunityId)
}

func (client *Client) CreateOpportunityNote(ctx context.Context, opportunityId string, input *model.NoteInput) (*model.Note, error) {
	return client.createNoteAndAssociate(ctx, "deals", opportunityId, "note_to_deal", input)
}

func (client *Client) listNotes(ctx context.Context, objectPath string, objectId string) ([]*model.Note, error) {
	query := url.Values{}
	url := fmt.Sprintf("%s/%s/%s/associations/notes?%s", client.BaseURL, objectPath, objectId, query.Encode())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	associationResponse := HSAssociationListSuccessResponse{}

	req = req.WithContext(ctx)
	if err := client.sendRequest(req, &associationResponse); err != nil {
		return nil, err
	}

	hsNote := HSNote{}
	notes := make([]*model.Note, len(associationResponse.Results))
	for i, hsAssociation := range associationResponse.Results {
		if err := client.get(ctx, "notes", hsAssociation.Id, []string{"hs_note_body"}, &hsNote); err != nil {
			return nil, err
		}

		notes[i] = client.mapNoteProperties(&hsNote)
	}

	return notes, nil
}

func (client *Client) createNoteAndAssociate(ctx context.Context, relatedObjectName string, relatedObjectId string, associationName string, input *model.NoteInput) (*model.Note, error) {
	notePayload := HSNoteCreateUpdatePayload{}
	notePayload.Properties.HsNoteBody = &input.Content
	hsTimestamp := formatHSDateTime(time.Now().UTC())
	notePayload.Properties.HsTimestamp = &hsTimestamp

	hsNote := HSNote{}
	if err := client.create(ctx, "notes", &notePayload, &hsNote); err != nil {
		return nil, err
	}

	query := url.Values{}
	url := fmt.Sprintf("%s/notes/%s/associations/%s/%s/%s?%s",
		client.BaseURL, hsNote.Id, relatedObjectName, relatedObjectId, associationName, query.Encode())
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return nil, err
	}
	if err := client.sendRequest(req, nil); err != nil {
		return nil, err
	}

	note := client.mapNoteProperties(&hsNote)
	return note, nil
}

func (client *Client) mapNoteProperties(hsNote *HSNote) *model.Note {
	return &model.Note{
		ID:        hsNote.Id,
		CreatedAt: parseHSDateTime(hsNote.Properties.HsCreatedDate),
		UpdatedAt: parseHSDateTime(hsNote.UpdatedAt),
		Content:   *hsNote.Properties.HsNoteBody,
	}
}
