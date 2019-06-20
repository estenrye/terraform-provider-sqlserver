# Terraform SQL Server Provider

This repository implements a Terraform Provider for SQL Server which you can 
use to create objects in SQL Server using Terraform.

The provider currently provides support for maintaining the following objects:

* Databases (Create/Drop)

# Using the Provider

## Example

```terraform
provider "sqlserver" {
  server = "localhost"
  user = "sa"
  password = "ownerTest1234"
}

resource "sqlserver_database" "hello" {
  database_name = "hello"
}
```

# Building the Provider for local development

1. [Download](https://golang.org/) and install Go
2. Clone the repository to your src directory in your gopath.
3. `./install.ps1 -skipGet`

# Installing the Provider

1. [Download](https://golang.org/) and install Go
2. `wget https://raw.githubusercontent.com/estenrye/terraform-provider-sqlserver/master/install.ps1`
3. `./install.ps1`

