package sqlserver

import (
	"database/sql"
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	// driver for database/sql
	_ "github.com/denisenkom/go-mssqldb"
)

func resourceSQLDBUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceSQLDBUserCreate,
		Read:   resourceSQLDBUserRead,
		Delete: resourceSQLDBUserDelete,
		//TODO:  Implement Importer
		// See https://www.terraform.io/docs/plugins/provider.html#resources

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"login": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"database": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSQLDBUserCreate(d *schema.ResourceData, m interface{}) error {
	name, err := cleanIdentifier(d.Get("name").(string), "User")
	if err != nil {
		return err
	}
	login, err := cleanIdentifier(d.Get("login").(string), "Login")
	if err != nil {
		return err
	}
	database, err := cleanIdentifier(d.Get("database").(string), "Database")
	if err != nil {
		return err
	}

	if login == "" {
		login = name
	}

	data := user{
		Name:     name,
		Login:    login,
		Database: database,
	}

	template := `USE [` + data.Database + `]; 
		IF ISNULL(USER_ID('` + data.Name + `'), -1) <> -1
		BEGIN
			DROP USER [` + data.Name + `]
		END
		CREATE USER [` + data.Name + `] FOR LOGIN [` + data.Login + `];
		SELECT USER_ID('` + data.Name + `')`

	log.Println(template)

	client := m.(*sqlServerClient)
	conn, err := sql.Open("mssql", client.connectionString)
	if err != nil {
		return ConnectionError{
			ConnectionString: client.connectionString,
			InnerError:       err,
		}
	}

	stmt, err := conn.Prepare(template)
	if err != nil {
		return QueryPreparationError{
			Query:      template,
			InnerError: err,
		}
	}

	row := stmt.QueryRow()

	var dbID string
	err = row.Scan(&dbID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	if dbID == "-1" {
		dbID = ""
	}

	log.Println(`Id: ` + dbID)
	d.SetId(dbID)
	return err
}

func resourceSQLDBUserRead(d *schema.ResourceData, m interface{}) error {
	database, err := cleanIdentifier(d.Get("database").(string), "Database")
	if err != nil {
		return err
	}

	template := `USE [` + database + `]; 
		SELECT u.name, u.principal_id, s.name 
		FROM sys.database_principals u 
		JOIN sys.server_principals s ON s.sid = u.sid 
		WHERE u.name = $1`

	log.Println(template)

	client := m.(*sqlServerClient)
	conn, err := sql.Open("mssql", client.connectionString)
	if err != nil {
		return ConnectionError{
			ConnectionString: client.connectionString,
			InnerError:       err,
		}
	}

	stmt, err := conn.Prepare(template)
	if err != nil {
		return QueryPreparationError{
			Query:      template,
			InnerError: err,
		}
	}

	row := stmt.QueryRow(d.Get("name").(string))

	var result user
	result.Database = database
	err = row.Scan(&result.Name, &result.ID, &result.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	d.Set("name", result.Name)
	d.SetId(result.ID)
	d.Set("login", result.Login)
	d.Set("database", database)

	return nil
}

func resourceSQLDBUserDelete(d *schema.ResourceData, m interface{}) error {
	database, err := cleanIdentifier(d.Get("database").(string), "Database")
	if err != nil {
		return err
	}

	name, err := cleanIdentifier(d.Get("name").(string), "User")
	if err != nil {
		return err
	}

	template := `USE [` + database + `[
		IF ISNULL(USER_ID('` + name + `'), -1) <> -1
		BEGIN
			DROP USER [` + name + `]
			SELECT -1
		END`
	log.Println(`[Trace] ` + template)
	_, err = executeQuery(m, "User", template)

	return err
}

type user struct {
	Name     string
	Login    string
	Database string
	ID       string
}

// TODO: Implement func resourceDatabaseImporter
//  See https://www.terraform.io/docs/plugins/provider.html#resources
