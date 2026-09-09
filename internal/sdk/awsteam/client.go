package awsteam

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hasura/go-graphql-client"
)

// SettingsCapabilities records which optional Settings API fields are available
// in the deployed TEAM version. It is populated once at client initialization
// via schema introspection so that mutations can omit fields the API doesn't
// know about.
type SettingsCapabilities struct {
	UseOUCacheSupported bool
}

type Client struct {
	GraphEndpoint        string
	GraphClient          *graphql.Client
	Config               *Config
	SettingsCapabilities *SettingsCapabilities
}

func (client *Client) detectSettingsCapabilities(ctx context.Context) (*SettingsCapabilities, error) {
	caps := &SettingsCapabilities{}

	q := `query { __type(name: "Settings") { fields { name } } }`

	raw, err := client.GraphClient.ExecRaw(ctx, q, nil)
	if err != nil {
		return nil, fmt.Errorf("introspecting Settings API capabilities: %w", err)
	}

	var result struct {
		Type *struct {
			Fields []struct {
				Name string `json:"name"`
			} `json:"fields"`
		} `json:"__type"`
	}

	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parsing Settings API capabilities response: %w", err)
	}

	if result.Type == nil {
		return nil, fmt.Errorf("introspecting Settings API capabilities: Settings type not found in schema")
	}

	for _, f := range result.Type.Fields {
		if f.Name == "useOUCache" {
			caps.UseOUCacheSupported = true
			break
		}
	}

	return caps, nil
}
