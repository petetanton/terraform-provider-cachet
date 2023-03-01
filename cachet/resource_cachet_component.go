package cachet

import (
	"context"
	"strconv"

	"github.com/andygrunwald/cachet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	name        = "name"
	description = "description"
	enabled     = "enabled"
	link        = "link"
	status      = "status"
	groupId     = "group_id"
)

var (
	componentStatuses = map[string]interface{}{
		"unknown":            cachet.ComponentStatusUnknown,
		"operational":        cachet.ComponentStatusOperational,
		"performance_issues": cachet.ComponentStatusPerformanceIssues,
		"partial_outage":     cachet.ComponentStatusPartialOutage,
		"major_outage":       cachet.ComponentStatusMajorOutage,
	}
)

func resourceCachetComponent() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			name: {
				Type:     schema.TypeString,
				Required: true,
			},
			description: {
				Type:     schema.TypeString,
				Required: true,
			},
			link: {
				Type:     schema.TypeString,
				Optional: true,
			},
			status: {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: StringInMapKeys(componentStatuses, false),
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
		},
		CreateContext:      resourceCachetComponentCreate,
		ReadContext:        resourceCachetComponentRead,
		UpdateContext:      resourceCachetComponentUpdate,
		DeleteContext:      resourceCachetComponentDelete,
		Importer:           nil,
		DeprecationMessage: "",
		Timeouts:           getDefaultTimeout(),
		Description:        "A component is a resource that defines a deployable thing that can be operational or degraded",
	}
}

func resourceCachetComponentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	_, err = client.Components.Delete(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceCachetComponentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client

	component := buildComponent(d)
	idInt, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	component.ID = idInt

	updatedComponent, _, err := client.Components.Update(idInt, component)
	if err != nil {
		return diag.FromErr(err)
	}

	return setComponent(d, updatedComponent)
}

func resourceCachetComponentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client

	component := buildComponent(d)

	createdComponent, _, err := client.Components.Create(component)
	if err != nil {
		diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(createdComponent.ID))
	return nil
}

func resourceCachetComponentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	component, _, err := client.Components.Get(id)
	if err != nil {
		return diag.FromErr(err)
	}

	return setComponent(d, component)
}

func buildComponent(d *schema.ResourceData) *cachet.Component {
	component := &cachet.Component{
		Name:        d.Get(name).(string),
		Description: d.Get(description).(string),
		Enabled:     d.Get(enabled).(bool),
	}

	if attr, ok := d.GetOk(link); ok {
		component.Link = attr.(string)
	}

	if attr, ok := d.GetOk(status); ok {
		component.Status = componentStatuses[attr.(string)].(int)
	} else {
		component.Status = cachet.ComponentStatusUnknown
	}

	if attr, ok := d.GetOk(groupId); ok {
		component.GroupID = attr.(int)
	}

	return component

}

func setComponent(d *schema.ResourceData, component *cachet.Component) diag.Diagnostics {
	d.SetId(strconv.Itoa(component.ID))
	d.Set(name, component.Name)

	d.Set(description, component.Description)
	d.Set(link, component.Link)
	d.Set(status, component.Status)
	d.Set(enabled, component.Enabled)
	d.Set(groupId, component.GroupID)

	return nil
}
