package cachet

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/petetanton/cachet-sdk"
)

const (
	public                 = "public"
	collapsed              = "collapsed"
	collapsedNo            = "no"
	collapsedYes           = "yes"
	collapsedIfOperational = "if_operational"
)

func resourceCachetComponentGroup() *schema.Resource {
	return &schema.Resource{
		Schema:             getComponentGroupSchema(false),
		CreateContext:      resourceCachetComponentGroupCreate,
		ReadContext:        resourceCachetComponentGroupRead,
		UpdateContext:      resourceCachetComponentGroupUpdate,
		DeleteContext:      resourceCachetComponentGroupDelete,
		Importer:           nil,
		DeprecationMessage: "",
		Description:        "A component group is a resource that defines a group of components",
	}
}

func resourceCachetComponentGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.ComponentGroups.Delete(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceCachetComponentGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client

	componentGroup := buildComponentGroup(d)
	idInt, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	componentGroup.ID = idInt

	updatedComponentGroup, _, err := client.ComponentGroups.Update(idInt, componentGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	return setComponentGroup(d, updatedComponentGroup)
}

func resourceCachetComponentGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client

	componentGroup := buildComponentGroup(d)

	createdComponentGroup, _, err := client.ComponentGroups.Create(componentGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(createdComponentGroup.ID))
	return nil
}

func resourceCachetComponentGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	componentGroup, _, err := client.ComponentGroups.Get(id)
	if err != nil {
		return diag.FromErr(err)
	}

	return setComponentGroup(d, componentGroup)
}

func buildComponentGroup(d *schema.ResourceData) *cachet.ComponentGroup {
	componentGroup := &cachet.ComponentGroup{
		Name: d.Get(name).(string),
	}

	if d.Get(public).(bool) {
		componentGroup.Visible = cachet.ComponentGroupVisibilityPublic
	} else {
		componentGroup.Visible = cachet.ComponentGroupVisibilityLoggedIn
	}

	if attr, ok := d.GetOk(collapsed); ok {
		componentGroup.Collapsed = getCollapsedInt(attr.(string))
	}

	return componentGroup

}

func setComponentGroup(d *schema.ResourceData, componentGroup *cachet.ComponentGroup) diag.Diagnostics {
	d.SetId(strconv.Itoa(componentGroup.ID))
	d.Set(name, componentGroup.Name)
	d.Set(collapsed, getCollapsedString(componentGroup.Collapsed))

	if componentGroup.Visible == cachet.ComponentGroupVisibilityPublic {
		d.Set(public, true)
	} else {
		d.Set(public, false)
	}

	return nil
}

func getCollapsedInt(in string) int {
	if strings.EqualFold(in, collapsedIfOperational) {
		return 2
	}

	if strings.EqualFold(in, collapsedYes) {
		return 1
	}

	return 0
}

func getCollapsedString(in int) string {
	if in == 2 {
		return collapsedIfOperational
	}

	if in == 1 {
		return collapsedYes
	}

	return collapsedNo
}
