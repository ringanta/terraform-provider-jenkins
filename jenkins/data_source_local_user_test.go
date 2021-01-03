package jenkins

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLocalUserDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "jenkins_local_user" "admin" {
					username = "admin"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.jenkins_local_user.admin", "username", "admin"),
					resource.TestCheckResourceAttr("data.jenkins_local_user.admin", "email", ""),
					resource.TestCheckResourceAttr("data.jenkins_local_user.admin", "fullname", "admin"),
				),
			},
		},
	})
}
