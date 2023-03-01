package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/petetanton/terraform-provider-cachet/cachet"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cachet.Provider})
}
