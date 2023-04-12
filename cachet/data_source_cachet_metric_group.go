package cachet

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/petetanton/cachet-sdk"
)

func dataSourceCachetMetricGroup() *schema.Resource {
	return &schema.Resource{
		Schema:      getMetricGroupSchema(true),
		ReadContext: dataSourceCachetMetricGroupRead,
	}
}

func dataSourceCachetMetricGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client
	metricGroupsResponse, _, err := client.MetricGroups.GetAll(&cachet.MetricGroupsQueryParams{
		Name:         d.Get("name").(string),
		QueryOptions: cachet.QueryOptions{},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	for _, metricGroup := range metricGroupsResponse.MetricGroups {
		if metricGroup.Name == d.Get("name").(string) {
			d.SetId(strconv.Itoa(metricGroup.ID))
			setMetricGroup(d, &metricGroup)
		}
	}

	return nil
}
