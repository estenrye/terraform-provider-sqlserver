package sqlserver

import (
	"github.com/hashicorp/terraform/helper/schema"
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

func resourceDatabaseCreate(d *schema.ResourceData, m interface{}) error {
	dbName, err := cleanIdentifier(d.Get("database_name").(string), "Database Name")
	if err != nil {
		d.SetId("")
		return err
	}

	template := `USE master
	IF (ISNULL(DB_ID($1), -1) = -1)
	BEGIN
		CREATE DATABASE [` + dbName + `]
	END
	SELECT ISNULL(DB_ID($1), -1)
	`

	dbID, err := executeQuery(m, "Database", template, d.Get("database_name").(string))

	d.SetId(dbID)
	return err
}

func resourceDatabaseRead(d *schema.ResourceData, m interface{}) error {
	template := `USE master
	SELECT ISNULL(DB_ID($1), -1)
	`

	dbID, err := executeQuery(m, "Database", template, d.Get("database_name").(string))

	d.SetId(dbID)
	return err
}

func resourceDatabaseDelete(d *schema.ResourceData, m interface{}) error {
	dbName, err := cleanIdentifier(d.Get("database_name").(string), "Database Name")
	if err != nil {
		d.SetId("")
		return err
	}

	template := `USE master
	IF (ISNULL(DB_ID($1), -1) <> -1)
	BEGIN
		ALTER DATABASE [` + dbName + `]  SET OFFLINE WITH ROLLBACK IMMEDIATE
		DROP DATABASE [` + dbName + `] 
	END
	`

	dbID, err := executeQuery(m, "Database", template, d.Get("database_name").(string))

	d.SetId(dbID)
	return err
}

type database struct {
	Name string
}

// TODO: Implement func resourceDatabaseImporter
//  See https://www.terraform.io/docs/plugins/provider.html#resources
