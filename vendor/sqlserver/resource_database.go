package sqlserver

import (
	"bytes"
	"database/sql"
	"github.com/hashicorp/terraform/helper/schema"
	"text/template"

	// driver for database/sql
	_ "github.com/denisenkom/go-mssqldb"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseCreate,
		Read:   resourceDatabaseRead,
		Delete: resourceDatabaseDelete,

		Schema: map[string]*schema.Schema{
			"database_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDatabaseExecuteQuery(d *schema.ResourceData, m interface{}, queryTemplate string) (string, error) {
	client := m.(*sqlServerClient)

	conn, err := sql.Open("mssql", client.connectionString)
	if err != nil {
		return "", ConnectionError{
			ConnectionString: client.connectionString,
			InnerError:       err,
		}
	}
	defer conn.Close()

	t := template.Must(template.New("template").Parse(queryTemplate))

	db := database{
		Name: d.Get("database_name").(string),
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, db)
	if err != nil {
		return "", err
	}
	query := tpl.String()

	stmt, err := conn.Prepare(query)
	if err != nil {
		return "", QueryPreparationError{
			Query:      query,
			InnerError: err,
		}
	}

	row := stmt.QueryRow()
	var dbID string
	err = row.Scan(&dbID)
	if err != nil {
		return "", err
	}

	if dbID == "-1" {
		dbID = ""
	}

	return dbID, nil
}

func resourceDatabaseCreate(d *schema.ResourceData, m interface{}) error {
	dbID, err := resourceDatabaseExecuteQuery(d, m, createTemplate)

	d.SetId(dbID)
	return err
}

func resourceDatabaseRead(d *schema.ResourceData, m interface{}) error {
	dbID, err := resourceDatabaseExecuteQuery(d, m, readTemplate)

	d.SetId(dbID)
	return err
}

func resourceDatabaseDelete(d *schema.ResourceData, m interface{}) error {
	dbID, err := resourceDatabaseExecuteQuery(d, m, readTemplate)

	d.SetId(dbID)
	return err
}

type database struct {
	Name string
}

var createTemplate = `USE master
IF (ISNULL(DB_ID('{{.Name}}'), -1) = -1)
BEGIN
	CREATE DATABASE [{{.Name}}]
END
SELECT DB_ID('{{.Name}}')
`

var readTemplate = `USE master
SELECT ISNULL(DB_ID('{{.Name}}'), -1)
GO
`

var deleteTemplate = `USE master
IF (ISNULL(DB_ID('{{.Name}}'), -1) <> -1)
BEGIN
	ALTER DATABASE [{{.Name}}] SET OFFLINE WITH ROLLBACK IMMEDIATE
	DROP DATABASE [{{.Name}}]
END
GO
`
