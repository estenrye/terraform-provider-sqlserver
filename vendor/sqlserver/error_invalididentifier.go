package sqlserver

import (
	"fmt"
)

// InvalidIdentifierError for reporting invalid characters in an identifier.
type InvalidIdentifierError struct {
	IdentifierType string
	FilterSet      string
}

func (e InvalidIdentifierError) Error() string {
	return fmt.Sprintf("Invalid characters in identifier. %v identifiers are not allowed to contain the following character: %v", e.IdentifierType, e.FilterSet)
}
