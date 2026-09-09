package awsteam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/smithy-go/ptr"
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
		return nil, errors.New("id is required to delete Eligibility")
	}

	q := fmt.Sprintf(`mutation DeleteEligibility {
		deleteEligibility(input: { id: "%s" }) {
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
