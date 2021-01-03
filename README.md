# Terraform Provider Jenkins

Run the following command to build the provider

```shell
go build -o terraform-provider-jenkins
```

## Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, run the following command to apply configuration

```shell
terraform init
terraform plan -out=tfplan.out
terraform apply tfplan.out
```

## Attribution

This provider is inspired from the work of https://github.com/taiidani/terraform-provider-jenkins.
Some codes and documentations are copy pasted from there.
