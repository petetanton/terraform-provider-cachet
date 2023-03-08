package cachet

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/andygrunwald/cachet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func Test_CachetSubscriberResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:        false,
		PreCheck:          func() { testPreCheck(t) },
		ProviderFactories: testProviderFactory(),
		CheckDestroy:      testCheckCachetSubscriberDestroy,
		Steps: []resource.TestStep{{
			Config: `
resource "cachet_subscriber" "this" {
  email = "admin@admin.com"
}
`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("cachet_subscriber.this", "email", "admin@admin.com"),
			),
		}},
	})
}

func testCheckCachetSubscriberDestroy(s *terraform.State) error {
	client := testProvider.Meta().(*Config).Client
	for _, r := range s.RootModule().Resources {
		if r.Type == "cachet_subscriber" {
			id, err := strconv.Atoi(r.Primary.ID)
			if err != nil {
				return err
			}
			subscribers, _, err := client.Subscribers.GetAll(&cachet.SubscribersQueryParams{})
			if err != nil {
				return err
			}
			for _, subscriber := range subscribers.Subscribers {
				if subscriber.ID == id {
					return fmt.Errorf("found subscriber %d when we didn't expect to", id)
				}
			}
		}
	}

	return nil
}
