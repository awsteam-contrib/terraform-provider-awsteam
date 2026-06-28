package awsteam

import (
	"context"
	"errors"
	"fmt"
)

type GetEligibilityInput struct {
	Id *string
}

type GetEligibilityOutput struct {
	Eligibility *Eligibility `json:"getEligibility"`
}

func (client *Client) GetEligibility(ctx context.Context, in *GetEligibilityInput) (*GetEligibilityOutput, error) {
	out := &GetEligibilityOutput{}

	if in.Id == nil {
		return nil, errors.New("Id is required to get Eligibility.")
	}

	query := fmt.Sprintf(`query GetEligibility {
		getEligibility(id: %s) {
			id
			name
			type
			ticketNo
			approvalRequired
			duration
			modifiedBy
			createdAt
			updatedAt
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
		}
	}
	`, graphqlStringPtr(in.Id))

	if err := client.executeGraphQL(ctx, query, nil, out); err != nil {
		return nil, err
	}

	return out, nil
}
