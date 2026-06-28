package awsteam

import (
	"context"
	"errors"
)

type UpdateEligibilityInput struct {
	Id               *string                  `json:"id"`
	Name             *string                  `json:"name"`
	Type             *string                  `json:"type"`
	Accounts         []*EligibilityAccount    `json:"accounts"`
	OUs              []*EligibilityOU         `json:"ous"`
	Permissions      []*EligibilityPermission `json:"permissions"`
	TicketNo         *string                  `json:"ticketNo"`
	ApprovalRequired *bool                    `json:"approvalRequired"`
	Duration         *int64                   `json:"duration"`
	ModifiedBy       *string                  `json:"modifiedBy"`
}

type UpdateEligibilityOutput struct {
	Eligibility *Eligibility `json:"updateEligibility"`
}

func (client *Client) UpdateEligibility(ctx context.Context, in *UpdateEligibilityInput) (*UpdateEligibilityOutput, error) {
	out := &UpdateEligibilityOutput{}

	if in.Id == nil {
		return nil, errors.New("Id is required to update Eligibility.")
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

	query := `mutation UpdateEligibility($input: UpdateEligibilityInput!) {
		updateEligibility(input: $input) {
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
