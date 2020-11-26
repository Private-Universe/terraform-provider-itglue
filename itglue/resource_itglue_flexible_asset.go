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
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
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

	a := &itglueRest.FlexibleAsset{}
	a.Data = asset.Data

	d.Set("traits", a.Data.Attributes.Traits)
	d.Set("organization_id", a.Data.Attributes.OrganizationID)
	d.Set("flexible_asset_type_id", a.Data.Attributes.FlexibleAssetTypeID)

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
	traits := d.Get("traits").(map[string]interface{})
	organizationID := d.Get("organization_id").(int)
	flexibleAssetTypeID := d.Get("flexible_asset_type_id").(int)

	if d.HasChanges("traits", "organization_id", "flexible_asset_type_id") {
		a := &itglueRest.FlexibleAsset{}
		a.Data.Attributes.Traits = traits
		a.Data.Attributes.OrganizationID = organizationID
		a.Data.Attributes.FlexibleAssetTypeID = flexibleAssetTypeID

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
