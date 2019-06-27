package sqlserver

import (
	"database/sql"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"

	// driver for database/sql
	_ "github.com/denisenkom/go-mssqldb"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseCreate,
		Read:   resourceDatabaseRead,
		Delete: resourceDatabaseDelete,
		//TODO:  Implement Importer
		// See https://www.terraform.io/docs/plugins/provider.html#resources

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
	dbName := cleanDatabaseName(d)

	conn, err := sql.Open("mssql", client.connectionString)
	if err != nil {
		return "", ConnectionError{
			ConnectionString: client.connectionString,
			InnerError:       err,
		}
	}
	defer conn.Close()

	stmt, err := conn.Prepare(queryTemplate)
	if err != nil {
		return "", QueryPreparationError{
			Query:      queryTemplate,
			InnerError: err,
		}
	}

	row := stmt.QueryRow(dbName)
	var dbID string
	err = row.Scan(&dbID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	if dbID == "-1" {
		dbID = ""
	}

	return dbID, nil
}

func cleanDatabaseName(d *schema.ResourceData) string {
	databaseName := d.Get("database_name").(string)
	databaseName = strings.TrimSpace(databaseName)
	databaseName = strings.ReplaceAll(databaseName, "\n", "")
	databaseName = strings.ReplaceAll(databaseName, "\r", "")
	databaseName = strings.ReplaceAll(databaseName, "[", "")
	databaseName = strings.ReplaceAll(databaseName, "]", "")

	return databaseName
}

func resourceDatabaseCreate(d *schema.ResourceData, m interface{}) error {
	dbName := cleanDatabaseName(d)

	createTemplate := `USE master
	IF (ISNULL(DB_ID($1), -1) = -1)
	BEGIN
		CREATE DATABASE [` + dbName + `]
	END
	SELECT DB_ID($1)
	`

	dbID, err := resourceDatabaseExecuteQuery(d, m, createTemplate)

	d.SetId(dbID)
	return err
}

func resourceDatabaseRead(d *schema.ResourceData, m interface{}) error {
	readTemplate := `USE master
	SELECT ISNULL(DB_ID($1), -1)
	`

	dbID, err := resourceDatabaseExecuteQuery(d, m, readTemplate)

	d.SetId(dbID)
	return err
}

func resourceDatabaseDelete(d *schema.ResourceData, m interface{}) error {
	dbName := cleanDatabaseName(d)

	deleteTemplate := `USE master
	IF (ISNULL(DB_ID($1), -1) <> -1)
	BEGIN
		ALTER DATABASE [` + dbName + `]  SET OFFLINE WITH ROLLBACK IMMEDIATE
		DROP DATABASE [` + dbName + `] 
	END
	`

	dbID, err := resourceDatabaseExecuteQuery(d, m, deleteTemplate)

	d.SetId(dbID)
	return err
}

type database struct {
	Name string
}

// TODO: Implement func resourceDatabaseImporter
//  See https://www.terraform.io/docs/plugins/provider.html#resources
