package orderbus

import "fmt"

type statusSet struct {
	Created   Status
	Finished  Status
	Cancelled Status
}

var Statuses = statusSet{
	Created:   newStatus("CREATED"),
	Finished:  newStatus("FINISHED"),
	Cancelled: newStatus("CANCELLED"),
}

// =============================================================================

var statuses = make(map[string]Status)

type Status struct {
	name string
}

func newStatus(status string) Status {
	r := Status{status}
	statuses[status] = r
	return r
}

func (s Status) String() string {
	return s.name
}

func (s Status) Equal(r2 Status) bool {
	return s.name == r2.name
}

// =============================================================================

func ParseStatus(value string) (Status, error) {
	status, exists := statuses[value]
	if !exists {
		return Status{}, fmt.Errorf("invalid status %q", value)
	}

	return status, nil
}

func MustParseStatus(value string) Status {
	status, err := ParseStatus(value)
	if err != nil {
		panic(err)
	}

	return status
}
