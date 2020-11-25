package itglue

import (
	"strconv"

	itglueRest "github.com/Private-Universe/itglue"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceITGlueFlexibleAsset() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"traits": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"organization-id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"flexible-asset-type-id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
		Create: resourceITGlueFlexibleAssetCreate,
		Read:   resourceITGlueFlexibleAssetRead,
		Update: resourceITGlueFlexibleAssetUpdate,
		Delete: resourceITGlueFlexibleAssetDelete,
	}
}

func resourceITGlueFlexibleAssetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*itglueRest.ITGAPI)
	traits := d.Get("traits").(map[string]interface{})
	organizationID := d.Get("organization-id").(int)
	flexibleAssetTypeID := d.Get("flexible-asset-type-id").(int)

	a := &itglueRest.FlexibleAsset{}
	a.Data.Attributes.Traits = traits
	a.Data.Attributes.OrganizationID = organizationID
	a.Data.Attributes.FlexibleAssetTypeID = flexibleAssetTypeID

	asset, err := client.PostFlexibleAsset(a)
	if err != nil {
		return err
	}

	d.SetId(asset.Data.ID)
	return resourceITGlueFlexibleAssetRead(d, meta)
}

func resourceITGlueFlexibleAssetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*itglueRest.ITGAPI)
	sid := d.Id()
	id, err := strconv.Atoi(sid)
	if err != nil {
		return err
	}
	asset, err := client.GetFlexibleAssetsByID(id)
	if err != nil {
		return err
	}

	a := &itglueRest.FlexibleAsset{}
	a.Data = asset.Data

	d.Set("traits", a.Data.Attributes.Traits)
	d.Set("organization-id", a.Data.Attributes.OrganizationID)
	d.Set("flexible-asset-type-id", a.Data.Attributes.FlexibleAssetTypeID)

	return nil
}

func resourceITGlueFlexibleAssetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*itglueRest.ITGAPI)
	sid := d.Id()
	id, err := strconv.Atoi(sid)
	if err != nil {
		return err
	}
	traits := d.Get("traits").(map[string]interface{})
	organizationID := d.Get("organization-id").(int)
	flexibleAssetTypeID := d.Get("flexible-asset-type-id").(int)

	if d.HasChanges("traits", "organization-id", "flexible-asset-type-id") {
		a := &itglueRest.FlexibleAsset{}
		a.Data.Attributes.Traits = traits
		a.Data.Attributes.OrganizationID = organizationID
		a.Data.Attributes.FlexibleAssetTypeID = flexibleAssetTypeID

		_, err = client.PatchFlexibleAsset(id, a)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceITGlueFlexibleAssetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*itglueRest.ITGAPI)
	sid := d.Id()
	id, err := strconv.Atoi(sid)
	if err != nil {
		return err
	}
	_, err = client.DeleteFlexibleAsset(id)
	if err != nil {
		return err
	}

	return nil
}
