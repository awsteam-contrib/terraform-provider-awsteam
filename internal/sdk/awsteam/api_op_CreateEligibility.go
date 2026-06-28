package awsteam

import (
	"context"
	"errors"
)

type CreateEligibilityInput struct {
	Id               *string                  `json:"id"`
	Name             *string                  `json:"name"`
	Type             *string                  `json:"type"`
	Accounts         []*EligibilityAccount    `json:"accounts"`
	OUs              []*EligibilityOU         `json:"ous"`
	Permissions      []*EligibilityPermission `json:"permissions"`
	TicketNo         *string                  `json:"ticketNo"`
	ApprovalRequired *bool                    `json:"approvalRequired"`
	Duration         *int64                   `json:"duration,string"`
	ModifiedBy       *string                  `json:"modifiedBy"`
}

type CreateEligibilityOutput struct {
	Eligibility *Eligibility `json:"createEligibility"`
}

func (client *Client) CreateEligibility(ctx context.Context, in *CreateEligibilityInput) (*CreateEligibilityOutput, error) {
	out := &CreateEligibilityOutput{}

	if in.Id == nil {
		return nil, errors.New("Id is required to create Eligibility.")
	}

	if in.Accounts == nil {
		in.Accounts = []*EligibilityAccount{}
	}

	if in.OUs == nil {
		in.OUs = []*EligibilityOU{}
	}

	variables := map[string]interface{}{
		"input": *in,
	}

	query := `mutation CreateEligibility($input: CreateEligibilityInput!) {
		createEligibility(input: $input) {
		id
		name
		type
		accounts {
		  name
		  id
		}
		ous {
		  name
		  id
		}
		permissions {
		  name
		  id
		}
		ticketNo
		approvalRequired
		duration
		modifiedBy
		createdAt
		updatedAt
	  }
	}`

	if err := client.executeGraphQL(ctx, query, variables, out); err != nil {
		return nil, err
	}

	return out, nil
}
