# jenkins_local_user Data Source

Get attributes of a local user from Jenkins.
The target Jenkins ssytem must use Jenkin's own user database as its security realm.

## Example Usage

```hcl
data "jenkins_local_user" "admin" {
  username = "admin"
}
```

## Argument Reference

The following arguments are supported:
- `username` - (Required) Username of the user being read.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `description` - Description of the Jenkins local user.
- `email` - Email address of the Jenkins local user.
- `fullname` - Full name of the Jenkins local user.
- `password_hash` - Password hash of the Jenkins local user.
