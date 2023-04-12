package cachet

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/petetanton/cachet-sdk"

	"github.com/petetanton/terraform-provider-cachet/pkg/cachet2"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CACHET_TOKEN", nil),
				Description: "The API token to use for Cachet. Generate it from within your account on the Cachet server.",
			},
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CACHET_URL", nil),
				Description: "The endpoint where Cachet is running.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cachet_component":       resourceCachetComponent(),
			"cachet_component_group": resourceCachetComponentGroup(),
			"cachet_metric":          resourceCachetMetric(),
			"cachet_metric_group":    resourceCachetMetricGroup(),
			"cachet_subscriber":      resourceCachetSubscriber(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cachet_component":       dataSourceCachetComponent(),
			"cachet_component_group": dataSourceCachetComponentGroup(),
			"cachet_metric":          dataSourceCachetMetric(),
			"cachet_metric_group":    dataSourceCachetMetricGroup(),
		},
	}

	p.ConfigureContextFunc = func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
		config, err := providerConfig(data)
		if err != nil {
			return nil, diag.Diagnostics{diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
				Detail:   err.Error(),
			}}
		}

		return config, nil
	}

	return p
}

func providerConfig(data *schema.ResourceData) (interface{}, error) {
	url := data.Get("api_url").(string)
	token := data.Get("token").(string)

	client, err := cachet.NewClient(url, nil)
	if err != nil {
		return nil, err
	}

	client.Authentication.SetTokenAuth(token)

	return &Config{Client: client, Client2: cachet2.New(url, token)}, nil
}

type Config struct {
	Client  *cachet.Client
	Client2 *cachet2.Client
}
