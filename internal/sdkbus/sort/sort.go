// Package sorting provides support for data sorting.
package sort

import (
	"fmt"
	"strings"
)

// Set of directions for data sorting.
const (
	ASC  = "ASC"
	DESC = "DESC"
)

var directions = map[string]string{
	ASC:  "ASC",
	DESC: "DESC",
}

// By represents a field used to sort by and direction.
type By struct {
	Field     string
	Direction string
}

// NewBy constructs a new By value with no checks.
func NewBy(field string, direction string) By {
	if _, exists := directions[direction]; !exists {
		return By{
			Field:     field,
			Direction: ASC,
		}
	}

	return By{
		Field:     field,
		Direction: direction,
	}
}

// Parse constructs a By value by parsing a string in the form of "field,direction" ie "user_id,ASC".
func Parse(fieldMappings map[string]string, sortBy string, defaultSortBy By) (By, error) {
	if sortBy == "" {
		return defaultSortBy, nil
	}

	sortParts := strings.Split(sortBy, ",")

	orgFieldName := strings.TrimSpace(sortParts[0])
	fieldName, exists := fieldMappings[orgFieldName]
	if !exists {
		return By{}, fmt.Errorf("unknown: %s", orgFieldName)
	}

	switch len(sortParts) {
	case 1:
		return NewBy(fieldName, ASC), nil

	case 2:
		direction := strings.ToUpper(strings.TrimSpace(sortParts[1]))
		if _, exists := directions[direction]; !exists {
			return By{}, fmt.Errorf("unknown direction: %s", direction)
		}

		return NewBy(fieldName, direction), nil

	default:
		return By{}, fmt.Errorf("unknown: %s", sortBy)
	}
}
