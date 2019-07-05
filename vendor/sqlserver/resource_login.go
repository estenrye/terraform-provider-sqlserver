package sqlserver

import (
	"database/sql"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	// driver for database/sql
	_ "github.com/denisenkom/go-mssqldb"
)

func resourceSQLLogin() *schema.Resource {
	return &schema.Resource{
		Create: resourceSQLLoginCreate,
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
			"login_password_hash": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceSQLLoginCreate(d *schema.ResourceData, m interface{}) error {
	data := login{
		Name:         d.Get("login_name").(string),
		PasswordHash: d.Get("login_password_hash").(string),
		SID:          cleanString(d.Get("login_sid").(string)),
	}

	quotedLogin, err := cleanIdentifier(data.Name, "Login")
	template := `CREATE LOGIN [` + quotedLogin + `] WITH PASSWORD = ` + data.PasswordHash + ` HASHED`

	if data.SID != "" {
		template += `, SID = ` + data.SID
	}

	template += `; SELECT ISNULL(SUSER_ID('` + cleanString(data.Name) + `'), -1)`

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

func resourceSQLLoginRead(d *schema.ResourceData, m interface{}) error {
	template := `USE master
	SELECT name, principal_id, sid, CONVERT(VARCHAR(512), password_hash, 1) FROM sys.sql_logins WHERE name = $1
	`

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

	row := stmt.QueryRow(d.Get("login_name").(string))

	var result login
	err = row.Scan(&result.Name, &result.ID, &result.SID, &result.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	d.SetId(result.ID)
	d.Set("login_name", result.Name)
	sID := d.Get("login_name").(string)
	if sID == "" {
		d.Set("login_sid", result.SID)
	}
	d.Set("login_password_hash", result.PasswordHash)

	return err
}

func resourceSQLLoginDelete(d *schema.ResourceData, m interface{}) error {
	data := login{
		Name:         d.Get("login_name").(string),
		PasswordHash: d.Get("login_password_hash").(string),
		SID:          d.Get("login_sid").(string),
	}

	deleteTemplate := `USE master
	DECLARE @login NVARCHAR(128) = $1
	DECLARE @sql nvarchar(200) = 'DROP LOGIN ' + QUOTENAME(@login)
	EXEC sp_sqlexec @sql`

	log.Println(`[Trace] ` + deleteTemplate)
	dbID, err := executeQuery(m, "Login", deleteTemplate, data.Name)

	d.SetId(dbID)
	return err
}

type login struct {
	Name         string
	PasswordHash string
	SID          string
	ID           string
}

func resourceSQLLoginUpdate(d *schema.ResourceData, m interface{}) error {
	data := login{
		Name:         d.Get("login_name").(string),
		PasswordHash: d.Get("login_password_hash").(string),
		SID:          cleanString(d.Get("login_sid").(string)),
	}

	updateTemplate := `DECLARE @login NVARCHAR(128) = 'helloUser43'
	DECLARE @pass NVARCHAR(128) = 'passworD''1235'
	DECLARE @sql nvarchar(300) = 'ALTER LOGIN ' + QUOTENAME(@login) + ' WITH PASSWORD = ' + QUOTENAME(@pass, '''')
	EXEC sp_sqlexec @sql
	`

	log.Println(`[Trace] ` + updateTemplate)

	dbID, err := executeQuery(m, "Login", updateTemplate, data.Name, data.PasswordHash)

	d.SetId(dbID)
	return err
}

// TODO: Implement func resourceDatabaseImporter
//  See https://www.terraform.io/docs/plugins/provider.html#resources
