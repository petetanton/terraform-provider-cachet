package cachet

import (
	"context"
	"time"

	"github.com/andygrunwald/cachet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CACHET_TOKEN", nil),
			},
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CACHET_URL", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cachet_component":       resourceCachetComponent(),
			"cachet_component_group": resourceCachetComponentGroup(),
			"cachet_metric":          resourceCachetMetric(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cachet_component":       dataSourceCachetComponent(),
			"cachet_component_group": dataSourceCachetComponentGroup(),
			"cachet_metric":          dataSourceCachetMetric(),
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
	client, err := cachet.NewClient(data.Get("api_url").(string), nil)
	if err != nil {
		return nil, err
	}

	client.Authentication.SetTokenAuth(data.Get("token").(string))

	return &Config{Client: client}, nil
}

type Config struct {
	Client *cachet.Client
}

func getDefaultTimeout() *schema.ResourceTimeout {
	timeout := time.Second * 30
	return &schema.ResourceTimeout{
		Create:  &timeout,
		Read:    &timeout,
		Update:  &timeout,
		Delete:  &timeout,
		Default: &timeout,
	}
}
