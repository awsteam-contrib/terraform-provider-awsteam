package awsteam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/smithy-go/ptr"
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
		return nil, errors.New("id is required to get Eligibility")
	}

	q := fmt.Sprintf(`query GetEligibility {
		getEligibility(id: "%s") {
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
