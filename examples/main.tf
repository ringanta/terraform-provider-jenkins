data "jenkins_local_user" "admin" {
  username = "admin"
}

resource "jenkins_local_user" "test" {
  username = "test"
  password = "testpwd"
  email    = "test@localhost"
  fullname = "test"
}

resource "jenkins_local_user" "test2" {
  username = "test2"
  password = "test2pwd"
  email    = "test2@localhost"
  fullname = "Test 2"
}

output "admin_email" {
  description = "Email of admin user"
  value       = data.jenkins_local_user.admin.email
}

output "test_fullname" {
  description = "Full name of the test user"
  value       = jenkins_local_user.test.fullname
}
