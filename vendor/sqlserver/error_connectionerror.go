package sqlserver

import (
	"fmt"
)

// ConnectionError Custom error class for reporting Query Preparation Errors
type ConnectionError struct {
	ConnectionString string
	InnerError       error
}

func (e ConnectionError) Error() string {
	return fmt.Sprintf("an error has occured while connecting to the sql server database:\n%v\n%v", e.ConnectionString, e.InnerError)
}
