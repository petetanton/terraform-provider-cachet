package cachet

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func Test_CachetMetricResource(t *testing.T) {
	metricName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	metricDescription := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		IsUnitTest:        false,
		PreCheck:          func() { testPreCheck(t) },
		ProviderFactories: testProviderFactory(),
		CheckDestroy:      testCheckCachetMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPlanSimpleMetric(metricName, metricDescription, "this"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cachet_metric.this", "name", metricName),
					resource.TestCheckResourceAttr("cachet_metric.this", "description", metricDescription),
				),
			},
			{
				Config:  testPlanSimpleMetric(metricName, metricDescription, "this"),
				Destroy: true,
			},
		},
	})
}

func testPlanSimpleMetric(name, description, tfRef string) string {
	return fmt.Sprintf(`
resource "cachet_metric" "%s" {
  name = "%s"
  description = "%s"
  visibility = "public"
  unit = "percent"
  default_value = 0
  default_view = "12_HOURS"
  mins_between_datapoints = 5
  decimal_places = 4
}
`, tfRef, name, description)
}

func testCheckCachetMetricDestroy(s *terraform.State) error {
	client := testProvider.Meta().(*Config).Client
	for _, r := range s.RootModule().Resources {
		if r.Type == "cachet_metric" {
			id, err := strconv.Atoi(r.Primary.ID)
			if err != nil {
				return err
			}
			_, response, err := client.Metrics.Get(id)
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
