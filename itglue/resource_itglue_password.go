package itglue

import (
	"context"
	"strconv"

	itglueRest "github.com/Private-Universe/itglue"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePassword() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePasswordCreate,
		ReadContext:   resourcePasswordRead,
		UpdateContext: resourcePasswordUpdate,
		DeleteContext: resourcePasswordDelete,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"organization_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"resource_id": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"password_category_id": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
			"password_category_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password_folder_id": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
		},
	}
}

func resourcePasswordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*itglueRest.ITGAPI)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	password := d.Get("password").(string)
	username := d.Get("username").(string)
	organizationID := d.Get("organization_id").(int)
	url := d.Get("url").(string)
	notes := d.Get("notes").(string)
	resourceID := d.Get("resource_id").(int)
	resourceType := d.Get("resource_type").(string)
	passwordCategoryID := d.Get("password_category_id").(int)
	passwordFolderID := d.Get("password_folder_id").(int)

	passwordType := "passwords"

	p := &itglueRest.Password{}
	p.Data.Type = passwordType
	p.Data.Attributes.Name = name
	p.Data.Attributes.Password = password
	p.Data.Attributes.Username = username
	p.Data.Attributes.OrganizationID = organizationID
	p.Data.Attributes.URL = url
	p.Data.Attributes.Notes = notes
	p.Data.Attributes.ResourceID = resourceID
	p.Data.Attributes.ResourceType = resourceType
	p.Data.Attributes.PasswordCategoryID = passwordCategoryID
	p.Data.Attributes.PasswordFolderID = passwordFolderID

	pass, err := client.PostPassword(p)
	if err != nil {
		return diag.FromErr(err)
	}

	newID := pass.Data.ID
	d.SetId(newID)
	resourcePasswordRead(ctx, d, meta)

	return diags
}

func resourcePasswordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*itglueRest.ITGAPI)

	var diags diag.Diagnostics

	s := d.Id()
	id, err := strconv.Atoi(s)
	if err != nil {
		return diag.FromErr(err)
	}

	p, err := client.GetPasswordsByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", p.Data.Attributes.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("password", p.Data.Attributes.Password); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("username", p.Data.Attributes.Username); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("organization_id", p.Data.Attributes.OrganizationID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("organization_name", p.Data.Attributes.OrganizationName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("url", p.Data.Attributes.URL); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("notes", p.Data.Attributes.Notes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("resource_id", p.Data.Attributes.ResourceID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("password_category_id", p.Data.Attributes.PasswordCategoryID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("password_category_name", p.Data.Attributes.PasswordCategoryName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("password_folder_id", p.Data.Attributes.PasswordFolderID); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePasswordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*itglueRest.ITGAPI)

	s := d.Id()
	id, err := strconv.Atoi(s)
	if err != nil {
		return diag.FromErr(err)

	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		p := &itglueRest.Password{}
		passwordType := "passwords"
		p.Data.Type = passwordType
		p.Data.Attributes.Name = name

		_, err = client.PatchPassword(id, p)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("password") {
		password := d.Get("password").(string)
		p := &itglueRest.Password{}
		passwordType := "passwords"
		p.Data.Type = passwordType
		p.Data.Attributes.Password = password

		_, err = client.PatchPassword(id, p)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("username") {
		username := d.Get("username").(string)
		p := &itglueRest.Password{}
		passwordType := "passwords"
		p.Data.Type = passwordType
		p.Data.Attributes.Username = username
		_, err = client.PatchPassword(id, p)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("organization_id") {
		organizationID := d.Get("organization_id").(int)
		p := &itglueRest.Password{}
		passwordType := "passwords"
		p.Data.Type = passwordType
		p.Data.Attributes.OrganizationID = organizationID

		_, err = client.PatchPassword(id, p)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("url") {
		url := d.Get("url").(string)
		p := &itglueRest.Password{}
		passwordType := "passwords"
		p.Data.Type = passwordType
		p.Data.Attributes.URL = url

		_, err = client.PatchPassword(id, p)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("notes") {
		notes := d.Get("notes").(string)
		p := &itglueRest.Password{}
		passwordType := "passwords"
		p.Data.Type = passwordType
		p.Data.Attributes.Notes = notes

		_, err = client.PatchPassword(id, p)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("resource_id") {
		resourceID := d.Get("resource_id").(int)
		p := &itglueRest.Password{}
		passwordType := "passwords"
		p.Data.Type = passwordType
		p.Data.Attributes.ResourceID = resourceID

		_, err = client.PatchPassword(id, p)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("resource_type") {
		resourceType := d.Get("resource_type").(string)
		p := &itglueRest.Password{}
		passwordType := "passwords"
		p.Data.Type = passwordType
		p.Data.Attributes.ResourceType = resourceType

		_, err = client.PatchPassword(id, p)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("password_category_id") {
		passwordCategoryID := d.Get("password_category_id").(int)
		p := &itglueRest.Password{}
		passwordType := "passwords"
		p.Data.Type = passwordType
		p.Data.Attributes.PasswordCategoryID = passwordCategoryID

		_, err = client.PatchPassword(id, p)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("password_folder_id") {
		passwordFolderID := d.Get("password_folder_id").(int)
		p := &itglueRest.Password{}
		passwordType := "passwords"
		p.Data.Type = passwordType
		p.Data.Attributes.PasswordFolderID = passwordFolderID

		_, err = client.PatchPassword(id, p)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourcePasswordRead(ctx, d, meta)
}

func resourcePasswordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*itglueRest.ITGAPI)

	var diags diag.Diagnostics

	s := d.Id()
	id, err := strconv.Atoi(s)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.DeletePassword(id)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
