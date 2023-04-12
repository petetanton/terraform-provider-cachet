package cachet

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/petetanton/cachet-sdk"
)

func resourceCachetMetricGroup() *schema.Resource {
	return &schema.Resource{
		Schema:             getMetricGroupSchema(false),
		CreateContext:      resourceCachetMetricGroupCreate,
		ReadContext:        resourceCachetMetricGroupRead,
		UpdateContext:      resourceCachetMetricGroupUpdate,
		DeleteContext:      resourceCachetMetricGroupDelete,
		Importer:           nil,
		DeprecationMessage: "",
		Description:        "A metric group is a resource that defines a group of metrics",
	}
}

func resourceCachetMetricGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.MetricGroups.Delete(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceCachetMetricGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client

	componentGroup := buildMetricGroup(d)
	idInt, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	componentGroup.ID = idInt

	updatedMetricGroup, _, err := client.MetricGroups.Update(idInt, componentGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	return setMetricGroup(d, updatedMetricGroup)
}

func resourceCachetMetricGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client

	metricGroup := buildMetricGroup(d)

	createdMetricGroup, _, err := client.MetricGroups.Create(metricGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(createdMetricGroup.ID))
	return nil
}

func resourceCachetMetricGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	metricGroup, _, err := client.MetricGroups.Get(id)
	if err != nil {
		return diag.FromErr(err)
	}

	return setMetricGroup(d, metricGroup)
}

func buildMetricGroup(d *schema.ResourceData) *cachet.MetricGroup {
	metricGroup := &cachet.MetricGroup{
		Name: d.Get(name).(string),
	}

	if d.Get(public).(bool) {
		metricGroup.Visible = cachet.MetricGroupVisibilityPublic
	} else {
		metricGroup.Visible = cachet.MetricGroupVisibilityLoggedIn
	}

	return metricGroup

}

func setMetricGroup(d *schema.ResourceData, metricGroup *cachet.MetricGroup) diag.Diagnostics {
	d.SetId(strconv.Itoa(metricGroup.ID))
	d.Set(name, metricGroup.Name)
	if metricGroup.Visible == cachet.MetricGroupVisibilityPublic {
		d.Set(public, true)
	} else {
		d.Set(public, false)
	}

	return nil
}
