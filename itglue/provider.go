package itglue

import (
	itglueRest "github.com/Private-Universe/itglue"
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
			"itglue_flexible_asset": resourceITGlueFlexibleAsset(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureFunc:  configureFunc(),
	}
}

func configureFunc() func(*schema.ResourceData) (interface{}, error) {
	return func(d *schema.ResourceData) (interface{}, error) {
		client := itglueRest.NewITGAPI(d.Get("api_key").(string))
		return client, nil
	}
}
