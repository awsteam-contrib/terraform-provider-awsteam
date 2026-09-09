package awsteam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/smithy-go/ptr"
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
		return nil, errors.New("id is required to delete Approvers")
	}

	q := fmt.Sprintf(`mutation DeleteApprovers {
		deleteApprovers(input: { id: "%s" }) {
			id
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
