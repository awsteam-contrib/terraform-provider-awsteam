package provider

import (
	"context"
	"fmt"

	"github.com/aws/smithy-go/ptr"
	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/acctest"
	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/sdk/awsteam"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccEligibilityResourceExists(ctx context.Context, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource (%s) ID not set", resourceName)
		}

		client := acctest.NewAWSTeamClient(ctx)
		out, err := client.GetEligibility(ctx, &awsteam.GetEligibilityInput{Id: ptr.String(rs.Primary.ID)})

		if err != nil {
			return err
		}

		if out == nil {
			return fmt.Errorf("Eligibility %q does not exist", rs.Primary.ID)
		}

		if out.Eligibility == nil {
			return fmt.Errorf("Eligibility %q does not exist", rs.Primary.ID)
		}

		return nil
	}
}

func testAccEligibilityResourceDisappears(ctx context.Context, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource (%s) ID not set", resourceName)
		}

		client := acctest.NewAWSTeamClient(ctx)
		_, err := client.DeleteEligibility(ctx, &awsteam.DeleteEligibilityInput{Id: ptr.String(rs.Primary.ID)})

		if err != nil {
			return err
		}

		return nil
	}
}
