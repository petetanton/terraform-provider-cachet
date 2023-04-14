package cachet

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/petetanton/terraform-provider-cachet/pkg/cachet2"
)

func resourceCachetSubscriber() *schema.Resource {
	return &schema.Resource{
		Schema:        getSubscriberSchema(false),
		CreateContext: resourceCachetSubscriberCreate,
		ReadContext:   resourceCachetSubscriberRead,
		DeleteContext: resourceCachetSubscriberDelete,
		Description:   "A subscriber is an email address that will receive updates about components",
	}
}

func resourceCachetSubscriberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client2

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.SubscriberDelete(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceCachetSubscriberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client2

	createdSubscriber, err := client.SubscriberCreate(&cachet2.Subscriber{
		Email:  d.Get("email").(string),
		Verify: true,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(createdSubscriber.ID))
	return nil
}

func resourceCachetSubscriberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).Client2
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	subscriber, err := client.SubscriberGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	return setSubscriber(d, subscriber)
}

func setSubscriber(d *schema.ResourceData, subscriber *cachet2.Subscriber) diag.Diagnostics {
	d.SetId(strconv.Itoa(subscriber.ID))
	d.Set("email", subscriber.Email)

	return nil
}
