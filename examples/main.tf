provider "sqlserver" {
  server = "localhost"
  user = "sa"
  password = "ownerTest1234"
}

resource "sqlserver_database" "hello" {
  database_name = "hello"
}

resource "sqlserver_sqllogin" "myLogin5" {
  name = "myLogin5"
  password_hash = "0x020088A31FCA2925F288B2A76C83B9245ABD1A7E5781B2300788B5B5E019D8439EB79D6B2F77553181F4F80BC4617F7731EDE21D59F80CC641BA7D2A79680A555453D0AC9795"
}
# Password:  P@ssw0rd
resource "sqlserver_sqllogin" "myLogin6" {
  name = "myLogin6"
  password_hash = "0x020088A31FCA2925F288B2A76C83B9245ABD1A7E5781B2300788B5B5E019D8439EB79D6B2F77553181F4F80BC4617F7731EDE21D59F80CC641BA7D2A79680A555453D0AC9795"
  sid = "0xDEADBEEFDEADBEEFDEADBEEF00000060"
}

resource "sqlserver_databaseuser" "myLogin5" {
  name = sqlserver_sqllogin.myLogin5.name
  login = sqlserver_sqllogin.myLogin5.name
  database = "hello"
}

resource "sqlserver_databaseuser" "myLogin6" {
  name = sqlserver_sqllogin.myLogin6.name
  login = sqlserver_sqllogin.myLogin6.name
  database = "hello"
}