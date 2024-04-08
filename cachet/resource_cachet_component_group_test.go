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

func Test_CachetComponentGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:        false,
		PreCheck:          func() { testPreCheck(t) },
		ProviderFactories: testProviderFactory(),
		CheckDestroy:      testCheckCachetComponentGroupDestroy,
		Steps: []resource.TestStep{
			testStepComponent("service-name", "service-description"),
			testStepComponentGroup("service-name", "service-description", "group-name", true),
			testStepComponentGroup("service-name", "service-description", "group-name", false),
		},
	})
}

func testStepComponentGroup(serviceName, serviceDescription, groupName string, public bool) resource.TestStep {
	pubStr := "false"
	if public {
		pubStr = "true"
	}
	return resource.TestStep{
		Config: testPlanSimpleServiceWithGroup(serviceName, serviceDescription, groupName, pubStr),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cachet_component.this", "name", serviceName),
			resource.TestCheckResourceAttr("cachet_component.this", "description", serviceDescription),
			resource.TestCheckResourceAttr("cachet_component_group.this", "name", groupName),
			resource.TestCheckResourceAttr("cachet_component_group.this", "public", pubStr),
		),
	}
}

func testPlanSimpleServiceWithGroup(serviceName, serviceDescription, groupName, pubStr string) string {
	return fmt.Sprintf(`
resource "cachet_component" "this" {
  name = "%s"
  description = "%s"
}

resource "cachet_component_group" "this" {
  name = "%s"
  public = %s
}
`, serviceName, serviceDescription, groupName, pubStr)
}

func testCheckCachetComponentGroupDestroy(s *terraform.State) error {
	client := testProvider.Meta().(*Config).Client
	for _, r := range s.RootModule().Resources {
		if r.Type == "cachet_component_group" {
			id, err := strconv.Atoi(r.Primary.ID)
			if err != nil {
				return err
			}
			_, response, err := client.ComponentGroups.Get(id)
			if err != nil && !strings.Contains(err.Error(), "404 Not Found") {
				return err
			}
			if response.StatusCode != http.StatusNotFound {
				return fmt.Errorf("component group %d has HTTP status %d", id, response.StatusCode)
			}
		}
	}

	return nil
}
