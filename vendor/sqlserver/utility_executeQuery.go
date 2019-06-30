package sqlserver

import (
	"database/sql"
	"strings"
	// driver for database/sql
	_ "github.com/denisenkom/go-mssqldb"
)

func executeQuery(m interface{}, resourceType string, queryTemplate string, args ...interface{}) (string, error) {
	_, err := cleanIdentifier(args[0].(string), resourceType)
	if err != nil {
		return "", err
	}

	conn, err := getSQLConnection(m)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	stmt, err := prepareSQLStatement(conn, queryTemplate)
	if err != nil {
		return "", err
	}

	row := stmt.QueryRow(args...)
	return getIDFromRow(row)

}

func getSQLConnection(m interface{}) (*sql.DB, error) {
	client := m.(*sqlServerClient)
	conn, err := sql.Open("mssql", client.connectionString)
	if err != nil {
		return nil, ConnectionError{
			ConnectionString: client.connectionString,
			InnerError:       err,
		}
	}
	return conn, nil
}

func prepareSQLStatement(conn *sql.DB, template string) (*sql.Stmt, error) {
	stmt, err := conn.Prepare(template)
	if err != nil {
		return nil, QueryPreparationError{
			Query:      template,
			InnerError: err,
		}
	}
	return stmt, nil
}

func getIDFromRow(row *sql.Row) (string, error) {
	var dbID string
	err := row.Scan(&dbID)
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

func cleanIdentifier(identifier string, identifierType string) (string, error) {

	cleanIdentifier := strings.TrimSpace(identifier)

	cleanIdentifier = strings.ReplaceAll(cleanIdentifier, "\n", "")
	cleanIdentifier = strings.ReplaceAll(cleanIdentifier, "\r", "")
	cleanIdentifier = strings.ReplaceAll(cleanIdentifier, "[", "")
	cleanIdentifier = strings.ReplaceAll(cleanIdentifier, "]", "")

	if identifier != cleanIdentifier {
		return "", InvalidIdentifierError{
			IdentifierType: identifierType,
			FilterSet:      `Leading or Trailing Whitespace, Newlines, and/or '[' or ']' characters.`,
		}
	}

	return identifier, nil
}

func cleanString(value string) string {
	cleanValue := strings.ReplaceAll(value, "'", "''")
	return cleanValue
}
