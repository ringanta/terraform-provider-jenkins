# jenkins_authorization_global_matrix Resource

Manage global matrix permission set for local user on the Jenkins system.
The target Jenkins system must use Jenkin's own user database as its security realm.

## Example Usage

```hcl
resource "jenkins_authorization_global_matrix" "example" {
  username = "example"
  permissions = [
    "Overall/Read",
    "Job/Build",
    "Job/Cancel",
    "Job/Read"
  ]
}
```

## Argument Reference

The following arguments are required:

- `username` - (Required) Username of the local user.
- `permissions` - (Required) Permission set of the local user.
  Permission format are `<group>/<action>`.
  They are similiar with the permission name on the Jenkins authorization dashboard.


## Import

Local user can be imported using the `username` field, e.g.

```hcl
terraform import jenkins_authorization_global_matrix.example example
```