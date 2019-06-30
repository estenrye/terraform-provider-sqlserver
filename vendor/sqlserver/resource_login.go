package sqlserver

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSQLLogin() *schema.Resource {
	return &schema.Resource{
		Create: resourceSQLLoginCreate,
		Exists: resourceSQLLoginExists,
		Read:   resourceSQLLoginRead,
		Delete: resourceSQLLoginDelete,
		Update: resourceSQLLoginUpdate,
		//TODO:  Implement Importer
		// See https://www.terraform.io/docs/plugins/provider.html#resources

		Schema: map[string]*schema.Schema{
			"login_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"login_sid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"login_password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: false,
			},
		},
	}
}

func resourceSQLLoginCreate(d *schema.ResourceData, m interface{}) error {
	data := login{
		Name:     d.Get("login_name").(string),
		Password: d.Get("login_password").(string),
		SID:      cleanString(d.Get("login_sid").(string)),
	}

	// createTemplate := `USE master
	// DECLARE @login NVARCHAR(128) = $1
	// DECLARE @pass NVARCHAR(128) = $2
	// `

	// if data.SID != "" {
	// 	createTemplate += `DECLARE @sid VARBINARY(16) = $3
	// 	`
	// }

	createTemplate := `DECLARE @sql nvarchar(400) = 'CREATE LOGIN ' + QUOTENAME($1) + ' WITH PASSWORD = ' + QUOTENAME($2, '''')`

	if data.SID != "" {
		createTemplate += `+ ', SID = ' + $3
		`
	}

	createTemplate += `EXEC sp_sqlexec @sql
	SELECT ISNULL(SUSER_ID($1), -1)
	`

	var dbID string
	var err error
	if data.SID == "" {
		dbID, err = executeQuery(m, "Login", createTemplate, data.Name, data.Password)
	} else {
		dbID, err = executeQuery(m, "Login", createTemplate, data.Name, data.Password, data.SID)
	}

	d.SetId(dbID)
	return err
}

func resourceSQLLoginRead(d *schema.ResourceData, m interface{}) error {
	data := login{
		Name:     d.Get("login_name").(string),
		Password: d.Get("login_password").(string),
		SID:      d.Get("login_sid").(string),
	}

	readTemplate := `USE master
	SELECT ISNULL(SUSER_ID($1), -1)
	`

	dbID, err := executeQuery(m, "Login", readTemplate, data.Name)

	d.SetId(dbID)
	return err
}

func resourceSQLLoginExists(d *schema.ResourceData, m interface{}) (bool, error) {
	data := login{
		Name:     d.Get("login_name").(string),
		Password: d.Get("login_password").(string),
		SID:      d.Get("login_sid").(string),
	}

	readTemplate := `USE master
	SELECT ISNULL(SUSER_ID($1), -1)
	`

	dbID, err := executeQuery(m, "Login", readTemplate, data.Name)

	if err != nil {
		return false, err
	}

	if dbID != "-1" {
		return false, nil
	}

	return true, nil
}

func resourceSQLLoginDelete(d *schema.ResourceData, m interface{}) error {
	data := login{
		Name:     d.Get("login_name").(string),
		Password: d.Get("login_password").(string),
		SID:      d.Get("login_sid").(string),
	}

	deleteTemplate := `USE master
	DECLARE @login NVARCHAR(128) = $1
	DECLARE @sql nvarchar(200) = 'DROP LOGIN ' + QUOTENAME(@login)
	EXEC sp_sqlexec @sql`

	dbID, err := executeQuery(m, "Login", deleteTemplate, data.Name)

	d.SetId(dbID)
	return err
}

type login struct {
	Name     string
	Password string
	SID      string
}

func resourceSQLLoginUpdate(d *schema.ResourceData, m interface{}) error {
	data := login{
		Name:     d.Get("login_name").(string),
		Password: d.Get("login_password").(string),
		SID:      cleanString(d.Get("login_sid").(string)),
	}

	updateTemplate := `DECLARE @login NVARCHAR(128) = 'helloUser43'
	DECLARE @pass NVARCHAR(128) = 'passworD''1235'
	DECLARE @sql nvarchar(300) = 'ALTER LOGIN ' + QUOTENAME(@login) + ' WITH PASSWORD = ' + QUOTENAME(@pass, '''')
	EXEC sp_sqlexec @sql
	`

	dbID, err := executeQuery(m, "Login", updateTemplate, data.Name, data.Password)

	d.SetId(dbID)
	return err
}

// TODO: Implement func resourceDatabaseImporter
//  See https://www.terraform.io/docs/plugins/provider.html#resources
