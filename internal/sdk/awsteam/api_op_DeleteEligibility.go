package awsteam

import (
	"context"
	"errors"
	"fmt"
)

type DeleteEligibilityInput struct {
	Id *string
}

type DeleteEligibilityOutput struct {
	Eligibility *Eligibility `json:"deleteEligibility"`
}

func (client *Client) DeleteEligibility(ctx context.Context, in *DeleteEligibilityInput) (*DeleteEligibilityOutput, error) {
	out := &DeleteEligibilityOutput{}

	if in.Id == nil {
		return nil, errors.New("Id is required to delete Eligibility.")
	}

	query := fmt.Sprintf(`mutation DeleteEligibility {
		deleteEligibility(input: { id: %s }) {
			id
		}
	}
	`, graphqlStringPtr(in.Id))

	if err := client.executeGraphQL(ctx, query, nil, out); err != nil {
		return nil, err
	}

	return out, nil
}
