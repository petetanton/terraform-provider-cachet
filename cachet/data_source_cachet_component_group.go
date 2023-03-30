package cachet

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/petetanton/cachet-sdk"
)

func dataSourceCachetComponentGroup() *schema.Resource {
	return &schema.Resource{
		Schema:      getComponentGroupSchema(true),
		ReadContext: dataSourceCachetComponentGroupRead,
	}
}

func dataSourceCachetComponentGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client
	componentGroups, _, err := client.ComponentGroups.GetAll(&cachet.ComponentGroupsQueryParams{
		Name:         d.Get("name").(string),
		QueryOptions: cachet.QueryOptions{},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	for _, component := range componentGroups.ComponentGroups {
		if component.Name == d.Get("name").(string) {
			d.SetId(strconv.Itoa(component.ID))
			setComponentGroup(d, &component)
		}
	}

	return nil
}
