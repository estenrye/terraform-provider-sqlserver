package sqlserver

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	// Example Provider requires an API Token.
	// The Email is optional
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1433,
			},
			"user": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sqlserver_database":     resourceDatabase(),
			"sqlserver_sqllogin":     resourceSQLLogin(),
			"sqlserver_databaseuser": resourceSQLDBUser(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	connString := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(d.Get("user").(string), d.Get("password").(string)),
		Host:   fmt.Sprintf("%s:%d", d.Get("server").(string), d.Get("port").(int)),
		Path:   d.Get("instance").(string),
	}

	client := sqlServerClient{
		connectionString: connString.String(),
	}

	return &client, nil
}

type sqlServerClient struct {
	connectionString string
}
