package jenkins

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLocalUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLocalUserCreate,
		ReadContext:   resourceLocalUserRead,
		UpdateContext: resourceLocalUserUpdate,
		DeleteContext: resourceLocalUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceLocalUserSchema,
	}
}

func resourceLocalUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(jenkinsClient)

	username := d.Get("username").(string)
	password := d.Get("password").(string)
	email := d.Get("email").(string)
	fullname := d.Get("fullname").(string)
	description := d.Get("description").(string)

	user, err := client.GetLocalUser(username)
	if err != nil {
		return diag.FromErr(err)
	}

	if user.Username != "" {
		return diag.Errorf("Local user %s is already existing in the Jenkins system", username)
	}

	err = client.CreateLocalUser(username, password, fullname, email, description)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(username)
	return resourceLocalUserRead(ctx, d, m)
}

func resourceLocalUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(jenkinsClient)

	username := d.Id()

	user, err := client.GetLocalUser(username)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("username", user.Username); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("email", user.Email); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("password_hash", user.PasswordHash); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("fullname", user.Fullname); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", user.Description); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceLocalUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(jenkinsClient)

	username := d.Get("username").(string)
	password := d.Get("password").(string)
	email := d.Get("email").(string)
	fullname := d.Get("fullname").(string)
	description := d.Get("description").(string)

	err := client.CreateLocalUser(username, password, fullname, email, description)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLocalUserRead(ctx, d, m)
}

func resourceLocalUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(jenkinsClient)
	var diags diag.Diagnostics

	username := d.Id()
	err := client.DeleteLocalUser(username)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

var resourceLocalUserSchema = map[string]*schema.Schema{
	"email": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Email address of the Jenkins local user",
	},
	"fullname": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Full name of the Jenkins local user",
	},
	"password": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Password of the jenkins local user",
		Sensitive:   true,
	},
	"password_hash": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Password hash of the jenkins local user",
	},
	"username": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Username of the Jenkins local user",
	},
	"description": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "Managed by Terraform",
		Description: "Description of the Jenkins local user",
	},
}
