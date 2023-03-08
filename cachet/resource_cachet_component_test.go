package cachet

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func Test_CachetComponentResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:        false,
		PreCheck:          func() { testPreCheck(t) },
		ProviderFactories: testProviderFactory(),
		CheckDestroy:      testCheckCachetComponentDestroy,
		Steps: []resource.TestStep{
			testStepComponent("name-1", "description-1"),
			testStepComponent("name-1", "description-2"),
		},
	})
}

func testStepComponent(name, description string) resource.TestStep {
	return resource.TestStep{
		Config: testPlanSimpleService(name, description),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cachet_component.this", "name", name),
			resource.TestCheckResourceAttr("cachet_component.this", "description", description),
		),
	}
}

func testPlanSimpleService(name, description string) string {
	return fmt.Sprintf(`
resource "cachet_component" "this" {
  name = "%s"
  description = "%s"
}
`, name, description)
}

func testCheckCachetComponentDestroy(s *terraform.State) error {
	client := testProvider.Meta().(*Config).Client
	for _, r := range s.RootModule().Resources {
		if r.Type == "cachet_component" {
			id, err := strconv.Atoi(r.Primary.ID)
			if err != nil {
				return err
			}
			_, response, err := client.Components.Get(id)
			if err != nil && !strings.Contains(err.Error(), "404 Not Found") {
				return err
			}
			if response.StatusCode != http.StatusNotFound {
				return fmt.Errorf("service %d has HTTP status %d", id, response.StatusCode)
			}
		}
	}

	return nil
}
