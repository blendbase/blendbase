package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"blendbase/graph/generated"
	"blendbase/graph/model"
	"context"
)

func (r *contactResolver) Notes(ctx context.Context, obj *model.Contact) ([]*model.Note, error) {
	c, err := r.getCrmConnector(ctx)
	if err != nil {
		return nil, err
	}

	return c.ListContactNotes(ctx, obj.ID)
}

func (r *crmResolver) Contact(ctx context.Context, obj *model.Crm, id string) (*model.Contact, error) {
	c, err := r.getCrmConnector(ctx)
	if err != nil {
		return nil, err
	}

	return c.GetContact(ctx, id)
}

func (r *crmResolver) Contacts(ctx context.Context, obj *model.Crm, first *int, after *string) (*model.ContactConnection, error) {
	c, err := r.getCrmConnector(ctx)
	if err != nil {
		return nil, err
	}

	firstOption := 10
	if first != nil {
		firstOption = *first
	}

	return c.ListContacts(ctx, firstOption, after)
}

func (r *crmResolver) Opportunities(ctx context.Context, obj *model.Crm, first *int, after *string) (*model.OpportunityConnection, error) {
	c, err := r.getCrmConnector(ctx)
	if err != nil {
		return nil, err
	}

	firstOption := 10
	if first != nil {
		firstOption = *first
	}

	return c.ListOpportunities(ctx, firstOption, after)
}

func (r *crmResolver) Opportunity(ctx context.Context, obj *model.Crm, id string) (*model.Opportunity, error) {
	c, err := r.getCrmConnector(ctx)
	if err != nil {
		return nil, err
	}

	return c.GetOpportunity(ctx, id)
}

func (r *mutationResolver) CreateContact(ctx context.Context, input model.ContactInput) (*model.Contact, error) {
	c, err := r.getCrmConnector(ctx)
	if err != nil {
		return nil, err
	}

	return c.CreateContact(ctx, &input)
}

func (r *mutationResolver) UpdateContact(ctx context.Context, id string, input model.ContactInput) (*bool, error) {
	c, err := r.getCrmConnector(ctx)
	if err != nil {
		return nil, err
	}

	success, err := c.UpdateContact(ctx, id, &input)

	return &success, err
}

func (r *mutationResolver) DeleteContact(ctx context.Context, id string) (*bool, error) {
	c, err := r.getCrmConnector(ctx)
	if err != nil {
		return nil, err
	}

	success, err := c.DeleteContact(ctx, id)

	return &success, err
}

func (r *mutationResolver) CreateContactNote(ctx context.Context, contactID string, input model.NoteInput) (*model.Note, error) {
	c, err := r.getCrmConnector(ctx)
	if err != nil {
		return nil, err
	}

	return c.CreateContactNote(ctx, contactID, &input)
}

func (r *mutationResolver) CreateOpportunity(ctx context.Context, input model.OpportunityInput) (*model.Opportunity, error) {
	c, err := r.getCrmConnector(ctx)
	if err != nil {
		return nil, err
	}

	opportunity, err := c.CreateOpportunity(ctx, &input)
	return opportunity, err
}

func (r *mutationResolver) UpdateOpportunity(ctx context.Context, id string, input model.OpportunityInput) (*bool, error) {
	c, err := r.getCrmConnector(ctx)
	if err != nil {
		return nil, err
	}

	success, err := c.UpdateOpportunity(ctx, id, &input)
	return &success, err
}

func (r *mutationResolver) DeleteOpportunity(ctx context.Context, id string) (*bool, error) {
	c, err := r.getCrmConnector(ctx)
	if err != nil {
		return nil, err
	}

	success, err := c.DeleteOpportunity(ctx, id)
	return &success, err
}

func (r *mutationResolver) CreateOpportunityNote(ctx context.Context, opportunityID string, input model.NoteInput) (*model.Note, error) {
	c, err := r.getCrmConnector(ctx)
	if err != nil {
		return nil, err
	}

	return c.CreateOpportunityNote(ctx, opportunityID, &input)
}

func (r *opportunityResolver) Notes(ctx context.Context, obj *model.Opportunity) ([]*model.Note, error) {
	c, err := r.getCrmConnector(ctx)
	if err != nil {
		return nil, err
	}

	return c.ListOpportunityNotes(ctx, obj.ID)
}

func (r *queryResolver) Crm(ctx context.Context) (*model.Crm, error) {
	return &model.Crm{}, nil
}

// Contact returns generated.ContactResolver implementation.
func (r *Resolver) Contact() generated.ContactResolver { return &contactResolver{r} }

// Crm returns generated.CrmResolver implementation.
func (r *Resolver) Crm() generated.CrmResolver { return &crmResolver{r} }

// Opportunity returns generated.OpportunityResolver implementation.
func (r *Resolver) Opportunity() generated.OpportunityResolver { return &opportunityResolver{r} }

type contactResolver struct{ *Resolver }
type crmResolver struct{ *Resolver }
type opportunityResolver struct{ *Resolver }
