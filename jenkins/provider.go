package jenkins

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider - create a new Jenkins provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JENKINS_URL", nil),
				Description: "The URL of the Jenkins server to connect to.",
			},
			"ca_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JENKINS_CA_CERT", nil),
				Description: "The path to the Jenkins self-signed certificate.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JENKINS_USERNAME", nil),
				Description: "Username to authenticate to Jenkins.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JENKINS_PASSWORD", nil),
				Description: "Password to authenticate to Jenkins. User API Token can be used to replace user password.",
			},
			"verify_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JENKINS_VERIFY_SSL", true),
				Description: "Flag to turn off ssl verification",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"jenkins_local_user":                  resourceLocalUser(),
			"jenkins_authorization_global_matrix": resourceAuthorizationGlobalMatrix(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"jenkins_local_user": dataSourceLocalUser(),
		},

		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		ServerURL: d.Get("server_url").(string),
		Username:  d.Get("username").(string),
		Password:  d.Get("password").(string),
		VerifySSL: d.Get("verify_ssl").(bool),
	}

	var err error
	if d.Get("ca_cert").(string) != "" {
		config.CACert, err = os.Open(d.Get("ca_cert").(string))
		if err != nil {
			return nil, diag.Errorf("Unable to open certificate file %s: %s", d.Get("ca_cert").(string), err.Error())
		}
	}

	client := newJenkinsClient(&config)
	if _, err = client.Init(); err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}
