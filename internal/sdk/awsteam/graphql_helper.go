package awsteam

import (
	"context"
	"encoding/json"
	"strconv"
)

func (client *Client) executeGraphQL(ctx context.Context, query string, variables map[string]interface{}, out interface{}) error {
	raw, err := client.GraphClient.ExecRaw(ctx, query, variables)
	if err != nil {
		return err
	}

	return json.Unmarshal(raw, out)
}

func defaultString(value *string, fallback string) string {
	if value == nil {
		return fallback
	}

	return *value
}

func graphqlString(value string) string {
	return strconv.Quote(value)
}

func graphqlStringPtr(value *string) string {
	return graphqlString(defaultString(value, ""))
}

func graphqlBoolPtr(value *bool) string {
	if value != nil && *value {
		return "true"
	}

	return "false"
}

func graphqlInt64Ptr(value *int64) string {
	if value == nil {
		return "0"
	}

	return strconv.FormatInt(*value, 10)
}

func graphqlQuotedInt64Ptr(value *int64) string {
	return graphqlString(graphqlInt64Ptr(value))
}

func graphqlJSON(value interface{}) (string, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}
