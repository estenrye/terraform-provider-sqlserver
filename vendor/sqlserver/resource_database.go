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
	dbName, err := cleanDatabaseName(d)
	if err != nil {
		return "", err
	}

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

func cleanDatabaseName(d *schema.ResourceData) (string, error) {
	specifiedDatabaseName := d.Get("database_name").(string)
	cleanDatabaseName := strings.TrimSpace(specifiedDatabaseName)

	cleanDatabaseName = strings.ReplaceAll(cleanDatabaseName, "\n", "")
	cleanDatabaseName = strings.ReplaceAll(cleanDatabaseName, "\r", "")
	cleanDatabaseName = strings.ReplaceAll(cleanDatabaseName, "[", "")
	cleanDatabaseName = strings.ReplaceAll(cleanDatabaseName, "]", "")

	if specifiedDatabaseName != cleanDatabaseName {
		return "", InvalidIdentifierError{
			IdentifierType: "Database Name",
			FilterSet:      `Leading or Trailing Whitespace, Newlines, and/or '[' or ']' characters.`,
		}
	}

	return specifiedDatabaseName, nil
}

func resourceDatabaseCreate(d *schema.ResourceData, m interface{}) error {
	dbName, err := cleanDatabaseName(d)
	if err != nil {
		d.SetId("")
		return err
	}

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
	_, err := cleanDatabaseName(d)
	if err != nil {
		d.SetId("")
		return err
	}

	readTemplate := `USE master
	SELECT ISNULL(DB_ID($1), -1)
	`

	dbID, err := resourceDatabaseExecuteQuery(d, m, readTemplate)

	d.SetId(dbID)
	return err
}

func resourceDatabaseDelete(d *schema.ResourceData, m interface{}) error {
	dbName, err := cleanDatabaseName(d)
	if err != nil {
		d.SetId("")
		return err
	}

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
