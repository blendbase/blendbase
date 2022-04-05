package test_utils

import (
	"blendbase/graph/model"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func GenerateContactInput() *model.ContactInput {
	companyName := gofakeit.Company()
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	email := gofakeit.Email()
	phone := gofakeit.Phone()
	website := gofakeit.DomainName()

	input := model.ContactInput{
		CompanyName: &companyName,
		FirstName:   &firstName,
		LastName:    &lastName,
		Email:       &email,
		Phone:       &phone,
		Website:     &website,
	}

	return &input
}

func GenerateOpportunityInput() *model.OpportunityInput {
	name := gofakeit.AppName()
	closeDate, _ := time.Parse(time.RFC3339, "2019-10-30T03:30:17.883Z")

	input := model.OpportunityInput{
		Name:      name,
		StageName: "contractsent",
		CloseDate: closeDate,
	}

	return &input
}

func GenerateNoteInput() *model.NoteInput {
	note := gofakeit.Sentence(30)

	input := model.NoteInput{
		Content: note,
	}

	return &input
}
