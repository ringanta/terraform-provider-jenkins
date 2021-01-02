data "jenkins_local_user" "admin" {
  username = "admin"
}

output "admin_email" {
  description = "Email of admin user"
  value       = data.jenkins_local_user.admin.email
}
