provider "sqlserver" {
  server = "localhost"
  user = "sa"
  password = "ownerTest1234"
}

resource "sqlserver_database" "hello" {
  database_name = "hello"
}