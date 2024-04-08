package cachet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func Test_CachetComponentGroupDatasource(t *testing.T) {
	groupName := "group-name"
	resource.Test(t, resource.TestCase{
		IsUnitTest:        false,
		PreCheck:          func() { testPreCheck(t) },
		ProviderFactories: testProviderFactory(),
		CheckDestroy:      testCheckCachetComponentGroupDestroy,
		Steps: []resource.TestStep{
			testStepComponentGroup("service", "description", groupName, true),
			testStepComponentGroupData(groupName),
		},
	})
}

func testStepComponentGroupData(groupName string) resource.TestStep {
	return resource.TestStep{
		Config: testPlanSimpleComponentGroupData(groupName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cachet_component_group.this", "name", groupName),
			testDatasourceComponentGroup("cachet_component_group.this", "data.cachet_component_group.this"),
		),
	}
}

func testDatasourceComponentGroup(src, dest string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[dest]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get a component ID from Cachet")
		}

		testAtts := []string{"id", "name"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("expected the component %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}

}

func testPlanSimpleComponentGroupData(name string) string {
	return fmt.Sprintf(`
resource "cachet_component_group" "this" {
  name = "%s"
  public = true
}

data "cachet_component_group" "this" {
  name = "%s"
}
`, name, name)
}
