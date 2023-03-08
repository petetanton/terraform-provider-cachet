package cachet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func Test_CachetComponentDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:        false,
		PreCheck:          func() { testPreCheck(t) },
		ProviderFactories: testProviderFactory(),
		CheckDestroy:      testCheckCachetComponentDestroy,
		Steps: []resource.TestStep{
			testStepComponent("name-1", "description-1"),
			testStepComponentData("name-1", "description-1"),
		},
	})
}

func testStepComponentData(name, description string) resource.TestStep {
	return resource.TestStep{
		Config: testPlanSimpleServiceData(name, description),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cachet_component.this", "name", name),
			resource.TestCheckResourceAttr("cachet_component.this", "description", description),
			testDatasourceComponent("cachet_component.this", "data.cachet_component.this"),
		),
	}
}

func testDatasourceComponent(src, dest string) resource.TestCheckFunc {
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

func testPlanSimpleServiceData(name, description string) string {
	return fmt.Sprintf(`
resource "cachet_component" "this" {
  name = "%s"
  description = "%s"
}

data "cachet_component" "this" {
  name = "%s"
}
`, name, description, name)
}
