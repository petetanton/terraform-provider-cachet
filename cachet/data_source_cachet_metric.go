package cachet

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/petetanton/cachet-sdk"
)

func dataSourceCachetMetric() *schema.Resource {
	return &schema.Resource{
		Schema:      getMetricSchema(true),
		ReadContext: dataSourceCachetMetricRead,
	}
}

func dataSourceCachetMetricRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client
	metrics, _, err := client.Metrics.GetAll(&cachet.MetricQueryParams{
		QueryOptions: cachet.QueryOptions{},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	metricName := d.Get("name").(string)

	metric, err := getMetricsPaginated(client, 0, metricName)
	if err != nil {
		return diag.FromErr(err)
	}

	if metric == nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("could not find a metric called %s, next page: %s", metricName, metrics.Meta.Pagination.Links.NextPage),
			},
		}
	}

	return setMetric(d, metric)
}

func getMetricsPaginated(client *cachet.Client, page int, metricName string) (*cachet.Metric, error) {
	metrics, _, err := client.Metrics.GetAll(&cachet.MetricQueryParams{
		QueryOptions: cachet.QueryOptions{
			Page:      page,
			PerPage:   25,
			SortField: "id",
		},
	})
	if err != nil {
		return nil, err
	}

	for _, metric := range metrics.Metrics {
		if metric.Name == metricName {
			return &metric, nil
		}
	}

	if metrics.Meta.Pagination.TotalPages > page {
		return getMetricsPaginated(client, page+1, metricName)
	}

	return nil, nil
}
