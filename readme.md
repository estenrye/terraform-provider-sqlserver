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

resource "sqlserver_sqllogin" "myLogin6" {
  name = "myLogin6"
  password_hash = "0x020088A31FCA2925F288B2A76C83B9245ABD1A7E5781B2300788B5B5E019D8439EB79D6B2F77553181F4F80BC4617F7731EDE21D59F80CC641BA7D2A79680A555453D0AC9795"
  sid = "0xDEADBEEFDEADBEEFDEADBEEF00000060"
}

resource "sqlserver_databaseuser" "myLogin6" {
  name = sqlserver_sqllogin.myLogin6.name
  login = sqlserver_sqllogin.myLogin6.name
  database = "hello"
}
```

# Building the Provider for local development

1. [Download](https://golang.org/) and install Go
2. Clone the repository to your src directory in your gopath.
3. `./initialize-devDependencies.ps1`
4. `./install.ps1 -skipGet`

# Installing the Provider

1. [Download](https://golang.org/) and install Go
2. `wget https://raw.githubusercontent.com/estenrye/terraform-provider-sqlserver/master/install.ps1`
3. `./install.ps1`

