package jenkins

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAuthorizationGlobalMatrix() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthorizationGlobalMatrixCreate,
		ReadContext:   resourceAuthorizationGlobalMatrixRead,
		UpdateContext: resourceAuthorizationGlobalMatrixUpdate,
		DeleteContext: resourceAuthorizationGlobalMatrixDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceAuthorizationGlobalMatrixSchema,
	}
}

func resourceAuthorizationGlobalMatrixCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(jenkinsClient)

	username := d.Get("username").(string)
	permissions := d.Get("permissions").([]interface{})
	permsStr := make([]string, len(permissions))

	for i, v := range permissions {
		permsStr[i] = fmt.Sprint(v)
	}

	err := client.CreateUserPermissions(username, permsStr)
	if err != nil {
		diag.FromErr(err)
	}

	d.SetId(username)
	return resourceAuthorizationGlobalMatrixRead(ctx, d, m)
}

func resourceAuthorizationGlobalMatrixRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(jenkinsClient)

	username := d.Id()
	userPermission, err := client.GetUserPermissions(username)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("username", userPermission.Username); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("permissions", userPermission.Permissions); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAuthorizationGlobalMatrixUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func resourceAuthorizationGlobalMatrixDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

var resourceAuthorizationGlobalMatrixSchema = map[string]*schema.Schema{
	"username": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Username of the Jenkins local user",
	},
	"permissions": {
		Type:        schema.TypeSet,
		Required:    true,
		Description: "Global matrix permission set for the local user",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
}
