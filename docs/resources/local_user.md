# jenkins_local_user Resource

Manage local user on the Jenkins system.
The target Jenkins system must use Jenkin's own user database as its security realm.

## Example Usage

```hcl
resource "jenkins_local_user" "example" {
  username = "example"
  password = "examplepwd"
  email    = "example@example.com"
  fullname = "Example"
}
```

## Argument Reference

The following arguments are required:

- `username` - (Required) Username of the local user.
- `password` - (Required) Password of the local user.
- `email` - (Required) Email of the local user.
- `fullname` - (Required) Fullname of the local user.

The following arguments are optional:

- `description` - (Optional) key value. Defaults to `Managed by Terraform`.

## Attributes Reference

In addition to all arguments above, the following attribute are exported:

- `password_hash` - Hash of current password.

## Import

Local user can be imported using the `username` field, e.g.

```hcl
terraform import jenkins_local_user.example example
```