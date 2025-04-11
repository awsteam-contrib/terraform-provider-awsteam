package acctest

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/envvar"
	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/sdk/awsteam"
)

// These two functions come form the terraform-provider-aws code
// RunSerialTests1Level runs test cases in parallel, optionally sleeping between each.
func RunSerialTests1Level(t *testing.T, testCases map[string]func(t *testing.T), d time.Duration) {
	t.Helper()

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc(t)
			time.Sleep(d)
		})
	}
}

// RunSerialTests2Levels runs test cases in parallel, optionally sleeping between each.
func RunSerialTests2Levels(t *testing.T, testCases map[string]map[string]func(t *testing.T), d time.Duration) {
	t.Helper()

	for group, m := range testCases {
		t.Run(group, func(t *testing.T) {
			RunSerialTests1Level(t, m, d)
		})
	}
}

func NewAWSTeamClient(ctx context.Context) *awsteam.Client {
	clientId := os.Getenv(envvar.AWSTEAMClientId)
	clientSecret := os.Getenv(envvar.AWSTEAMClientSecret)
	graphEndpoint := os.Getenv(envvar.AWSTEAMGraphEndpoint)
	TokenEndpoint := os.Getenv(envvar.AWSTEAMTokenEndpoint)

	config := &awsteam.Config{
		ClientId:      clientId,
		ClientSecret:  clientSecret,
		GraphEndpoint: graphEndpoint,
		TokenEndpoint: TokenEndpoint,
	}

	config.Build(ctx)

	return config.NewClient(ctx)
}
