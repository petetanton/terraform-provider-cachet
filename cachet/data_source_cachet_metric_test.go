package cachet

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func Test_CachetMetricDatasourceSimple(t *testing.T) {
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
				Config: testPlanSimpleMetric(metricName, metricDescription, "this") + "\n" + testDatasourcePlanMetric(metricName, "this"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cachet_metric.this", "name", metricName),
					resource.TestCheckResourceAttr("cachet_metric.this", "description", metricDescription),
					testDatasourceMetric("cachet_metric.this", "data.cachet_metric.this"),
				),
			},
		},
	})
}

func Test_CachetMetricDatasourceComplex(t *testing.T) {
	if os.Getenv("TEST_SKIP_SLOW") == "true" {
		t.SkipNow()
	}

	noOfMetricsToTest := acctest.RandIntRange(25, 100)

	var plan1 string
	var plan2 string
	var tfunc1 []resource.TestCheckFunc
	var tfunc2 []resource.TestCheckFunc
	for i := 0; i < noOfMetricsToTest; i++ {
		metricName := fmt.Sprintf("name-%d", i)
		metricDescription := fmt.Sprintf("description-%d", i)
		metricTfRef := fmt.Sprintf("this-%d", i)
		plan1 += "\n" + testPlanSimpleMetric(metricName, metricDescription, metricTfRef)
		tfunc1 = append(tfunc1, resource.TestCheckResourceAttr(fmt.Sprintf("cachet_metric.%s", metricTfRef), name, metricName))
		tfunc1 = append(tfunc1, resource.TestCheckResourceAttr(fmt.Sprintf("cachet_metric.%s", metricTfRef), description, metricDescription))

		plan2 += "\n" + testDatasourcePlanMetric(metricName, metricTfRef)
		tfunc2 = append(tfunc2, testDatasourceMetric(fmt.Sprintf("cachet_metric.%s", metricTfRef), fmt.Sprintf("data.cachet_metric.%s", metricTfRef)))
	}

	plan2 += plan1

	resource.Test(t, resource.TestCase{
		IsUnitTest:        false,
		PreCheck:          func() { testPreCheck(t) },
		ProviderFactories: testProviderFactory(),
		CheckDestroy:      testCheckCachetMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: plan1,
				Check: resource.ComposeTestCheckFunc(
					tfunc1...,
				),
			},
			{
				Config: plan2,
				Check: resource.ComposeTestCheckFunc(
					tfunc2...,
				),
			},
		},
	})
}

func testDatasourcePlanMetric(metricName, tfRef string) string {
	return fmt.Sprintf(`
data "cachet_metric" "%s" {
  name = "%s"
}
`, tfRef, metricName)
}

func testDatasourceMetric(src, dest string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[dest]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get a metric ID from Cachet")
		}

		testAtts := []string{"id", "name"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("expected the metric %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}

}
