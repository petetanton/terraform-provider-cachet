package cachet

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getComponentSchema(dataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		name: {
			Type:     schema.TypeString,
			Required: true,
		},
		description: {
			Type:     schema.TypeString,
			Required: !dataSource,
			Computed: dataSource,
		},
		link: {
			Type:     schema.TypeString,
			Optional: true,
		},
		status: {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "unknown",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"unknown", "operational", "performance_issues", "partial_outage", "major_outage"}, false)),
		},
		enabled: {
			Type:     schema.TypeBool,
			Default:  true,
			Optional: true,
		},
		groupId: {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
}

func getComponentGroupSchema(dataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		name: {
			Type:     schema.TypeString,
			Required: true,
		},
		public: {
			Type:     schema.TypeBool,
			Required: !dataSource,
			Optional: dataSource,
		},
	}
}
