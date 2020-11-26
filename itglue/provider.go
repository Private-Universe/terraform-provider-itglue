package itglue

import (
	"context"

	itglueRest "github.com/Private-Universe/itglue"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Description: "Your IT Glue API key",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ITGLUE_API_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"itglue_flexible_asset": resourceFlexibleAsset(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	key := d.Get("api_key").(string)
	client := itglueRest.NewITGAPI(key)

	return client, diags
}
