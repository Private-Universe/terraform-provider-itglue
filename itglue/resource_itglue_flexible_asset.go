package itglue

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"checkboxes": {
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeBool,
							},
							Optional: true,
						},
						"dates": {
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"text": {
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"textboxes": {
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"numbers": {
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeFloat,
							},
							Optional: true,
						},
						"percents": {
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeFloat,
							},
							Optional: true,
						},
						"selects": {
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"tags": {
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeList,
								Elem: &schema.Schema{
									Type: schema.TypeInt,
								},
							},
							Optional: true,
						},
					},
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

	traits := d.Get("traits").(*schema.Set).List()
	organizationID := d.Get("organization_id").(int)
	flexibleAssetTypeID := d.Get("flexible_asset_type_id").(int)

	a := &itglueRest.FlexibleAsset{}
	a.Data.Type = "flexible-assets"
	combinedTraits := combineTraits(traits)
	a.Data.Attributes.Traits = combinedTraits
	a.Data.Attributes.OrganizationID = organizationID
	a.Data.Attributes.FlexibleAssetTypeID = flexibleAssetTypeID

	asset, err := client.PostFlexibleAsset(a)
	if err != nil {
		return diag.Errorf("%s %s", err, traits)
	}

	newID := fmt.Sprintf("fa-%s", asset.Data.ID)
	d.SetId(newID)
	resourceFlexibleAssetRead(ctx, d, meta)

	return diags
}

func resourceFlexibleAssetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*itglueRest.ITGAPI)

	var diags diag.Diagnostics

	sid := d.Id()
	s := strings.Split(sid, "-")
	id, err := strconv.Atoi(s[1])
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

	sid := d.Id()
	s := strings.Split(sid, "-")
	id, err := strconv.Atoi(s[1])
	if err != nil {
		return diag.FromErr(err)

	}

	if d.HasChanges("traits", "organization_id", "flexible_asset_type_id") {
		traits := d.Get("traits").(*schema.Set).List()
		organizationID := d.Get("organization_id").(int)
		flexibleAssetTypeID := d.Get("flexible_asset_type_id").(int)

		asset := &itglueRest.FlexibleAsset{}
		combinedTraits := combineTraits(traits)
		asset.Data.Attributes.Traits = combinedTraits
		asset.Data.Attributes.OrganizationID = organizationID
		asset.Data.Attributes.FlexibleAssetTypeID = flexibleAssetTypeID
		a := flattenFlexibleAsset(asset)

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

	sid := d.Id()
	s := strings.Split(sid, "-")
	id, err := strconv.Atoi(s[1])
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

func flattenFlexibleAsset(fa *itglueRest.FlexibleAsset) *itglueRest.FlexibleAsset {
	nfa := &itglueRest.FlexibleAsset{}
	nmap := make(map[string]interface{})

	for i, asset := range fa.Data.Attributes.Traits {
		switch asset.(type) {
		case map[string]interface{}:
			tidl := getTraitTagIDList(asset.(map[string]interface{}))
			nmap[i] = tidl
		case string:
			nmap[i] = asset
		case float64:
			nmap[i] = asset
		default:
			continue
		}
	}

	nfa.Data.Attributes.OrganizationID = fa.Data.Attributes.OrganizationID
	nfa.Data.Attributes.FlexibleAssetTypeID = fa.Data.Attributes.FlexibleAssetTypeID
	nfa.Data.Attributes.Traits = nmap
	return nfa
}

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
