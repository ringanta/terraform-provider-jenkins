data "jenkins_local_user" "admin" {
  username = "admin"
}

resource "jenkins_local_user" "test" {
  username = "test"
  password = ""
  email    = "test@localhost"
  fullname = "Test"

  lifecycle {
    ignore_changes = [password]
  }
}

output "admin_description" {
  description = "Description of admin user"
  value       = data.jenkins_local_user.admin.description
}

output "test_fullname" {
  description = "Full name of the test user"
  value       = jenkins_local_user.test.fullname
}
