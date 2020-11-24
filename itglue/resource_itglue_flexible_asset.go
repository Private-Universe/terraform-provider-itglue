package itglue

import (
	"encoding/json"
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

	flexibleAsset := &itglueRest.FlexibleAsset{}
	flexibleAsset.Data.Attributes.Traits = traits
	flexibleAsset.Data.Attributes.OrganizationID = organizationID
	flexibleAsset.Data.Attributes.FlexibleAssetTypeID = flexibleAssetTypeID

	b, err := json.Marshal(flexibleAsset)
	if err != nil {
		return err
	}

	asset, err := client.PostFlexibleAsset(b)
	if err != nil {
		return err
	}

	newFlexibleAsset := &itglueRest.FlexibleAsset{}
	err = json.Unmarshal(asset, newFlexibleAsset)
	if err != nil {
		return err
	}

	d.SetId(newFlexibleAsset.Data.ID)
	return resourceITGlueFlexibleAssetRead(d, meta)
}

func resourceITGlueFlexibleAssetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*itglueRest.ITGAPI)
	id := d.Id()
	sid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	var a itglueRest.FlexibleAsset
	req, err := client.GetFlexibleAssetsJSONByID(sid)
	if err != nil {
		return err
	}
	asset := &itglueRest.FlexibleAsset{}
	err = json.Unmarshal(req, asset)
	if err != nil {
		return err
	}
	a.Data = asset.Data

	d.Set("traits", a.Data.Attributes.Traits)

	return nil
}

func resourceITGlueFlexibleAssetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*itglueRest.ITGAPI)
	id := d.Id()
	sid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	traits := d.Get("traits").(map[string]interface{})

	if d.HasChange("traits") {
		a := &itglueRest.FlexibleAsset{}
		a.Data.Attributes.Traits = traits
		b, err := json.Marshal(a)
		if err != nil {
			return err
		}
		_, err = client.PatchFlexibleAsset(sid, b)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceITGlueFlexibleAssetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*itglueRest.ITGAPI)
	id := d.Id()
	sid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	_, err = client.DeleteFlexibleAsset(sid)
	if err != nil {
		return err
	}

	return nil
}
