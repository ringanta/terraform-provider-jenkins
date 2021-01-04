data "jenkins_local_user" "admin" {
  username = "admin"
}

resource "jenkins_local_user" "test" {
  username = "test"
  password = "testpwd"
  email    = "test@localhost"
  fullname = "Test user"
}

resource "jenkins_authorization_global_matrix" "test" {
  username = jenkins_local_user.test.id
  permissions = [
    "Overall/Read",
    "Job/Build",
    "Job/Cancel",
    "Job/Read"
  ]
}

output "admin_fullname" {
  description = "Full name of the admin user"
  value       = data.jenkins_local_user.admin.fullname
}

output "test_fullname" {
  description = "Full name of the test user"
  value       = jenkins_local_user.test.fullname
}
