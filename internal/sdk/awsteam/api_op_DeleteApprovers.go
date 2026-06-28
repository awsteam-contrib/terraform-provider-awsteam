package awsteam

import (
	"context"
	"errors"
	"fmt"
)

type DeleteApproversInput struct {
	Id *string
}

type DeleteApproversOutput struct {
	Approvers *Approvers `json:"deleteApprovers"`
}

func (client *Client) DeleteApprovers(ctx context.Context, in *DeleteApproversInput) (*DeleteApproversOutput, error) {
	out := &DeleteApproversOutput{}

	if in.Id == nil {
		return nil, errors.New("Id is required to delete Approvers.")
	}

	query := fmt.Sprintf(`mutation DeleteApprovers {
		deleteApprovers(input: { id: %s }) {
			id
		}
	}
	`, graphqlStringPtr(in.Id))

	if err := client.executeGraphQL(ctx, query, nil, out); err != nil {
		return nil, err
	}

	return out, nil
}
