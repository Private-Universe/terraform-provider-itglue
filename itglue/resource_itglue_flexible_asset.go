package itglue

import (
	"context"
	"strconv"

	itglueRest "github.com/Private-Universe/itglue"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFlexibleAsset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlexibleAssetCreate,
		ReadContext:   resourceFlexibleAssetRead,
		UpdateContext: resourceFlexibleAssetUpdate,
		DeleteContext: resourceFlexibleAssetDelete,
		Schema: map[string]*schema.Schema{
			"traits": {
				Type:     schema.TypeMap,
				Required: true,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"flexible_asset_type_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceFlexibleAssetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*itglueRest.ITGAPI)

	var diags diag.Diagnostics

	traits := d.Get("traits").(map[string]interface{})
	organizationID := d.Get("organization_id").(int)
	flexibleAssetTypeID := d.Get("flexible_asset_type_id").(int)

	a := &itglueRest.FlexibleAsset{}
	a.Data.Type = "flexible-assets"
	a.Data.Attributes.Traits = traits
	a.Data.Attributes.OrganizationID = organizationID
	a.Data.Attributes.FlexibleAssetTypeID = flexibleAssetTypeID

	asset, err := client.PostFlexibleAsset(a)
	if err != nil {
		return diag.FromErr(err)
	}

	newID := asset.Data.ID
	d.SetId(newID)
	resourceFlexibleAssetRead(ctx, d, meta)

	return diags
}

func resourceFlexibleAssetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*itglueRest.ITGAPI)

	var diags diag.Diagnostics

	s := d.Id()
	id, err := strconv.Atoi(s)
	if err != nil {
		return diag.FromErr(err)

	}
	asset, err := client.GetFlexibleAssetsByID(id)
	if err != nil {
		return diag.FromErr(err)

	}

	a := flattenFlexibleAsset(asset)
	if err := d.Set("traits", a.Data.Attributes.Traits); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("organization_id", a.Data.Attributes.OrganizationID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("flexible_asset_type_id", a.Data.Attributes.FlexibleAssetTypeID); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceFlexibleAssetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*itglueRest.ITGAPI)

	s := d.Id()
	id, err := strconv.Atoi(s)
	if err != nil {
		return diag.FromErr(err)

	}

	// This is done together as any traits not specified in the patch request will be deleted in IT Glue.
	if d.HasChanges("traits", "organization_id", "flexible_asset_type_id") {
		traits := d.Get("traits").(map[string]interface{})
		organizationID := d.Get("organization_id").(int)
		flexibleAssetTypeID := d.Get("flexible_asset_type_id").(int)

		asset := &itglueRest.FlexibleAsset{}
		asset.Data.Attributes.Traits = traits
		asset.Data.Attributes.OrganizationID = organizationID
		asset.Data.Attributes.FlexibleAssetTypeID = flexibleAssetTypeID
		a := flattenFlexibleAsset(asset)
		a.Data.Type = "flexible-assets"

		_, err = client.PatchFlexibleAsset(id, a)
		if err != nil {
			return diag.FromErr(err)

		}
	}

	return resourceFlexibleAssetRead(ctx, d, meta)
}

func resourceFlexibleAssetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*itglueRest.ITGAPI)

	var diags diag.Diagnostics

	s := d.Id()
	id, err := strconv.Atoi(s)
	if err != nil {
		return diag.FromErr(err)

	}
	_, err = client.DeleteFlexibleAsset(id)
	if err != nil {
		return diag.FromErr(err)

	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

// This function flattens out flexible resources.
// Instead of having a list of structs which contain IDs somewhere, you get a list of just the IDs.
// This makes it achievable to compare/detect changes with Terraform.
// In the future, this will hopefully enable specifying multiple IDs in Tag Traits.
func flattenFlexibleAsset(fa *itglueRest.FlexibleAsset) *itglueRest.FlexibleAsset {
	nfa := &itglueRest.FlexibleAsset{}
	nmap := make(map[string]interface{})

	for i, asset := range fa.Data.Attributes.Traits {
		switch asset.(type) {
		case map[string]interface{}:
			tidl := strconv.FormatFloat(getTraitTagIDList(asset.(map[string]interface{}))[0], 'f', 0, 64)
			nmap[i] = tidl
		default:
			nmap[i] = asset
		}
	}

	nfa.Data.Attributes.OrganizationID = fa.Data.Attributes.OrganizationID
	nfa.Data.Attributes.FlexibleAssetTypeID = fa.Data.Attributes.FlexibleAssetTypeID
	nfa.Data.Attributes.Traits = nmap
	return nfa
}

// Will work better in the future for type map containing a list of integers
func getTraitTagIDList(tagList map[string]interface{}) []float64 {
	var list []float64
	for _, trait := range tagList["values"].([]interface{}) {
		l := trait.(map[string]interface{})
		for k, t := range l {
			if k == "id" {
				list = append(list, t.(float64))
			}
		}
	}
	return list
}

// Currently unused function which combines separate Terraform variables into one traits variable.
// E.g having separate Terraform variables for tags, text fields, numbers, etc.
func combineTraits(traitsList []interface{}) map[string]interface{} {
	saves := make(map[string]interface{})

	for _, traitGroup := range traitsList {
		l := traitGroup.(map[string]interface{})
		for _, tg := range l {
			t := tg.(map[string]interface{})
			if len(t) != 0 {
				for tkey, value := range t {
					saves[tkey] = value
				}
			}
		}
	}

	return saves
}
