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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password_hash": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"check_policy": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: false,
			},
			"check_expiration": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: false,
			},
		},
	}
}

func resourceSQLLoginCreate(d *schema.ResourceData, m interface{}) error {
	data := login{
		Name:            d.Get("name").(string),
		PasswordHash:    d.Get("password_hash").(string),
		SID:             cleanString(d.Get("sid").(string)),
		CheckPolicy:     d.Get("check_policy").(bool),
		CheckExpiration: d.Get("check_expiration").(bool),
	}

	quotedLogin, err := cleanIdentifier(data.Name, "Login")
	checkPolicy := `OFF`
	if data.CheckPolicy {
		checkPolicy = `ON`
	}
	checkExpiration := `OFF`
	if data.CheckExpiration {
		checkExpiration = `ON`
	}
	template := `CREATE LOGIN [` + quotedLogin +
		`] WITH PASSWORD = ` + data.PasswordHash +
		` HASHED, CHECK_POLICY = ` + checkPolicy +
		`, CHECK_EXPIRATION = ` + checkExpiration

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
	SELECT name, principal_id, sid, CONVERT(VARCHAR(512), password_hash, 1), is_policy_checked, is_expiration_checked FROM sys.sql_logins WHERE name = $1
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

	row := stmt.QueryRow(d.Get("name").(string))

	var result login
	err = row.Scan(&result.Name, &result.ID, &result.SID, &result.PasswordHash,
		&result.CheckPolicy, &result.CheckExpiration)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	d.SetId(result.ID)
	d.Set("name", result.Name)
	sID := d.Get("name").(string)
	if sID == "" {
		d.Set("sid", result.SID)
	}
	d.Set("password_hash", result.PasswordHash)

	return err
}

func resourceSQLLoginDelete(d *schema.ResourceData, m interface{}) error {
	data := login{
		Name:            d.Get("name").(string),
		PasswordHash:    d.Get("password_hash").(string),
		SID:             d.Get("sid").(string),
		CheckPolicy:     d.Get("check_policy").(bool),
		CheckExpiration: d.Get("check_expiration").(bool),
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
	Name            string
	PasswordHash    string
	SID             string
	ID              string
	CheckPolicy     bool
	CheckExpiration bool
}

func resourceSQLLoginUpdate(d *schema.ResourceData, m interface{}) error {
	data := login{
		Name:            d.Get("name").(string),
		PasswordHash:    d.Get("password_hash").(string),
		SID:             cleanString(d.Get("sid").(string)),
		CheckPolicy:     d.Get("check_policy").(bool),
		CheckExpiration: d.Get("check_expiration").(bool),
	}
	quotedLogin, err := cleanIdentifier(data.Name, "Login")
	checkPolicy := `OFF`
	if data.CheckPolicy {
		checkPolicy = `ON`
	}
	checkExpiration := `OFF`
	if data.CheckExpiration {
		checkExpiration = `ON`
	}

	template := `ALTER LOGIN [` + quotedLogin + `] WITH PASSWORD = ` +
		data.PasswordHash + ` HASHED, CHECK_POLICY = ` + checkPolicy +
		`, CHECK_EXPIRATION = ` + checkExpiration + `;`
	template += `SELECT ISNULL(SUSER_ID('` + cleanString(data.Name) + `'), -1)`

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

// TODO: Implement func resourceDatabaseImporter
//  See https://www.terraform.io/docs/plugins/provider.html#resources
