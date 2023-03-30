package cachet

import (
	"context"
	"strconv"

	"github.com/andygrunwald/cachet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	unit            = "unit"
	defaultValue    = "default_value"
	displayChart    = "display_chart"
	calculationType = "calculation_type"
	sum             = "sum"
	average         = "average"
	decimalPlaces   = "decimal_places"
	defaultView     = "default_view"
	threshold       = "mins_between_datapoints"
	visibility      = "visibility"

	visibilityPublic  = "public"
	visibilityPrivate = "private"
	visibilityHidden  = "hidden"

	viewHour    = "HOUR"
	view12Hours = "12_HOURS"
	viewWeek    = "WEEK"
	viewMonth   = "MONTH"
)

var (
	metricViewMap = map[string]interface{}{
		viewHour:    cachet.MetricsViewLastHour,
		view12Hours: cachet.MetricsViewLast12Hours,
		viewWeek:    cachet.MetricsViewLastWeek,
		viewMonth:   cachet.MetricsViewLastMonth,
	}
	metricVisibilityMap = map[string]interface{}{
		visibilityPublic:  cachet.MetricsVisibilityPublic,
		visibilityPrivate: cachet.MetricsVisibilityLoggedIn,
		visibilityHidden:  cachet.MetricsVisibilityHidden,
	}
	metricVisibilities = []string{visibilityPublic, visibilityPrivate, visibilityHidden}
	metricViews        = []string{viewHour, view12Hours, viewWeek, viewMonth}
)

func resourceCachetMetric() *schema.Resource {
	return &schema.Resource{
		Schema:        getMetricSchema(false),
		CreateContext: resourceCachetMetricCreate,
		ReadContext:   resourceCachetMetricRead,
		//UpdateContext:      resourceCachetMetricUpdate,
		DeleteContext:      resourceCachetMetricDelete,
		Importer:           nil,
		DeprecationMessage: "",
		Description:        "A component group is a resource that defines a group of components",
	}
}

func resourceCachetMetricDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Metrics.Delete(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// TODO: this needs update support in the SDK first!
//func resourceCachetMetricUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	client := meta.(*Config).Client
//
//	metirc := buildMetric(d)
//	idInt, err := strconv.Atoi(d.Id())
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	metirc.ID = idInt
//
//	updatedMetric, _, err := client.Metrics(idInt, metirc)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	return setMetric(d, updatedMetric)
//}

func resourceCachetMetricCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client

	metirc := buildMetric(d)

	createdMetric, _, err := client.Metrics.Create(metirc)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(createdMetric.ID))
	return nil
}

func resourceCachetMetricRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	metirc, _, err := client.Metrics.Get(id)
	if err != nil {
		return diag.FromErr(err)
	}

	return setMetric(d, metirc)
}

func buildMetric(d *schema.ResourceData) *cachet.Metric {
	metirc := &cachet.Metric{
		Name:         d.Get(name).(string),
		Suffix:       d.Get(unit).(string),
		Description:  d.Get(description).(string),
		DefaultValue: d.Get(defaultValue).(int),
		DisplayChart: d.Get(displayChart).(bool),
		Places:       d.Get(decimalPlaces).(int),
		DefaultView:  metricViewMap[d.Get(defaultView).(string)].(int),
		Threshold:    d.Get(threshold).(int),
		Visible:      metricVisibilityMap[d.Get(visibility).(string)].(int),
	}

	if d.Get(calculationType).(string) == sum {
		metirc.CalcType = cachet.MetricsCalculationSum
	} else {
		metirc.CalcType = cachet.MetricsCalculationAverage
	}

	return metirc

}

func setMetric(d *schema.ResourceData, metric *cachet.Metric) diag.Diagnostics {
	d.SetId(strconv.Itoa(metric.ID))
	d.Set(name, metric.Name)
	d.Set(unit, metric.Suffix)
	d.Set(description, metric.Description)
	d.Set(defaultValue, metric.DefaultValue)
	d.Set(displayChart, metric.DisplayChart)
	d.Set(decimalPlaces, metric.Places)

	d.Set(defaultView, FindIntInMap(metricViewMap, metric.DefaultView, "12_HOURS"))
	d.Set(threshold, metric.Threshold)
	d.Set(visibility, FindIntInMap(metricVisibilityMap, metric.Visible, "private"))

	if metric.CalcType == cachet.MetricsCalculationSum {
	} else {
		d.Set(calculationType, average)
	}

	return nil
}
