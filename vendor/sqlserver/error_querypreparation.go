package sqlserver

import (
	"fmt"
)

// Custom error class for reporting Query Preparation Errors
type QueryPreparationError struct {
	Query      string
	InnerError error
}

func (e QueryPreparationError) Error() string {
	return fmt.Sprintf("an error has occured in query preparation:\n%v\n%v", e.Query, e.InnerError)
}
