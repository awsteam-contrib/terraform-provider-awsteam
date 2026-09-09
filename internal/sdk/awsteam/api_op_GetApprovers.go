package awsteam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/smithy-go/ptr"
)

type GetApproversInput struct {
	Id *string
}

type GetApproversOutput struct {
	Approvers *Approvers `json:"getApprovers"`
}

func (client *Client) GetApprovers(ctx context.Context, in *GetApproversInput) (*GetApproversOutput, error) {
	out := &GetApproversOutput{}

	if in.Id == nil {
		return nil, errors.New("id is required to get Approvers")
	}

	q := fmt.Sprintf(`query GetApprovers {
		getApprovers(id: "%s") {
			id
			name
			type
			approvers
			groupIds
			ticketNo
			modifiedBy
			createdAt
			updatedAt
		}
	}	
	`, ptr.ToString(in.Id))

	raw, err := client.GraphClient.ExecRaw(ctx, q, nil)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw, out)

	if err != nil {
		return nil, err
	}

	return out, nil
}
