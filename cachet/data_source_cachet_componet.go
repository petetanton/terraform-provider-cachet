package cachet

import (
	"context"
	"strconv"

	"github.com/andygrunwald/cachet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCachetComponent() *schema.Resource {
	return &schema.Resource{
		Schema:      getComponentSchema(true),
		ReadContext: dataSourceCachetComponentRead,
	}
}

func dataSourceCachetComponentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client
	components, _, err := client.Components.GetAll(&cachet.ComponentsQueryParams{
		Name:         d.Get("name").(string),
		QueryOptions: cachet.QueryOptions{},
	})
	if err != nil {
		diag.FromErr(err)
	}

	for _, component := range components.Components {
		if component.Name == d.Get("name").(string) {
			d.SetId(strconv.Itoa(component.ID))
			setComponent(d, &component)
		}
	}

	return nil
}
