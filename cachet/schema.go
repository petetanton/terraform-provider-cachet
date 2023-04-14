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
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the component",
		},
		description: {
			Type:        schema.TypeString,
			Required:    !dataSource,
			Computed:    dataSource,
			Description: "Description of the component",
		},
		link: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A link to the component - can be used for Runbooks etc",
		},
		status: {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "unknown",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"unknown", "operational", "performance_issues", "partial_outage", "major_outage"}, false)),
			Description:      "The current status of thc component. If you have automation that updates the component status, please use lifecycle rules in terraform",
		},
		enabled: {
			Type:        schema.TypeBool,
			Default:     true,
			Optional:    true,
			Description: "Is the component enabled",
		},
		groupId: {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The group that the component is within",
		},
	}
}

func getComponentGroupSchema(dataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		name: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the component group",
		},
		public: {
			Type:        schema.TypeBool,
			Required:    !dataSource,
			Optional:    dataSource,
			Description: "Is the component group public?",
		},
		collapsed: {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{collapsedNo, collapsedYes, collapsedIfOperational}, false)),
			Default:          collapsedNo,
		},
	}
}

func getMetricSchema(dataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		name: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the metric",
		},
		description: {
			Type:        schema.TypeString,
			Required:    !dataSource,
			Optional:    dataSource,
			Description: "Descripion of the metric",
		},
		unit: {
			Type:        schema.TypeString,
			Required:    !dataSource,
			Optional:    dataSource,
			Description: "The unit suffix for the metrics",
		},
		defaultValue: {
			Type:        schema.TypeInt,
			Required:    !dataSource,
			Optional:    dataSource,
			Description: "The default value of the metric to be displayed when there is no data",
		},
		calculationType: {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{sum, average}, false)),
			Default:          sum,
			Description:      "Define how math should be performed on this metric",
		},
		displayChart: {
			Type:        schema.TypeBool,
			Default:     true,
			Optional:    true,
			Description: "Should thi metric be displayed on the status page?",
		},
		decimalPlaces: {
			Type:             schema.TypeInt,
			Required:         !dataSource,
			Optional:         dataSource,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			Description:      "The number of decimal places to use for this metric",
		},
		defaultView: {
			Type:             schema.TypeString,
			Required:         !dataSource,
			Optional:         dataSource,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(metricViews, false)),
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
			Description:      fmt.Sprintf("visibility of the metric, must be on of: %s", strings.Join(metricVisibilities, ",")),
		},
		groupId: {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The group that the metric is within",
		},
		order: {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "The order that the metric should appear in",
		},
	}
}

func getSubscriberSchema(dataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"email": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "the email address for the subscriber",
			ForceNew:    true,
		},
		"verify": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "whether to send a verification email",
			ForceNew:    true,
			Default:     false,
		},
	}
}

func getMetricGroupSchema(dataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		name: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the metric group",
		},
		public: {
			Type:        schema.TypeBool,
			Required:    !dataSource,
			Optional:    dataSource,
			Description: "Is the metric group public?",
		},
	}
}
