provider "sqlserver" {
  server = "localhost"
  user = "sa"
  password = "ownerTest1234"
}

resource "sqlserver_database" "hello" {
  database_name = "hello"
}

resource "sqlserver_sqllogin" "myLogin5" {
  login_name = "myLogin5"
  login_password = "Mypassword1234"
}

resource "sqlserver_sqllogin" "myLogin6" {
  login_name = "myLogin6"
  login_password = "Mypassword1234"
  login_sid = "0xDEADBEEFDEADBEEFDEADBEEF00000060"
}