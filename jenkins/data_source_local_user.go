package jenkins

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLocalUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLocalUserRead,
		Schema:      dataSourceLocalUserSchema,
	}
}

func dataSourceLocalUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(jenkinsClient)
	username := d.Get("username").(string)

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

	if err := d.Set("fullname", user.Fullname); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("password_hash", user.PasswordHash); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", user.Description); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(username)
	return nil
}

var dataSourceLocalUserSchema = map[string]*schema.Schema{
	"email": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Email address of the Jenkins local user",
	},
	"fullname": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Full name of the Jenkins local user",
	},
	"password_hash": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Password hash of the jenkins local user",
	},
	"username": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Username of the Jenkins local user",
	},
	"description": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Description of the Jenkins local user",
	},
}
