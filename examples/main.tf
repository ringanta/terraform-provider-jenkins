data "jenkins_local_user" "admin" {
  username = "admin"
}

resource "jenkins_local_user" "test" {
  username = "test"
  password = "testpwd"
  email    = "test@localhost"
  fullname = "Test user"
}

output "admin_fullname" {
  description = "Full name of the admin user"
  value       = data.jenkins_local_user.admin.fullname
}

output "test_fullname" {
  description = "Full name of the test user"
  value       = jenkins_local_user.test.fullname
}
