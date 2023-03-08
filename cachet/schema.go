package cachet

import (
	"fmt"
	"strings"

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

func getMetricSchema(dataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		name: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		description: {
			Type:     schema.TypeString,
			Required: !dataSource,
			Optional: dataSource,
			ForceNew: true,
		},
		unit: {
			Type:     schema.TypeString,
			Required: !dataSource,
			Optional: dataSource,
			ForceNew: true,
		},
		defaultValue: {
			Type:     schema.TypeInt,
			Required: !dataSource,
			Optional: dataSource,
			ForceNew: true,
		},
		calculationType: {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{sum, average}, false)),
			Default:          sum,
			ForceNew:         true,
		},
		displayChart: {
			Type:     schema.TypeBool,
			Default:  true,
			Optional: true,
			ForceNew: true,
		},
		decimalPlaces: {
			Type:             schema.TypeInt,
			Required:         !dataSource,
			Optional:         dataSource,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			ForceNew:         true,
		},
		defaultView: {
			Type:             schema.TypeString,
			Required:         !dataSource,
			Optional:         dataSource,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(metricViews, false)),
			ForceNew:         true,
			Description:      fmt.Sprintf("default view of the metric, must be on of: %s", strings.Join(metricViews, ",")),
		},
		threshold: {
			Type:             schema.TypeInt,
			Required:         !dataSource,
			Optional:         dataSource,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
			ForceNew:         true,
		},
		visibility: {
			Type:             schema.TypeString,
			Required:         !dataSource,
			Optional:         dataSource,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(metricVisibilities, false)),
			ForceNew:         true,
			Description:      fmt.Sprintf("visibility of the metric, must be on of: %s", strings.Join(metricVisibilities, ",")),
		},
	}
}
