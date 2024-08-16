package orderbus

import "fmt"

type statusSet struct {
	Created  Status
	Finished Status
}

var Statuses = statusSet{
	Created:  newStatus("CREATED"),
	Finished: newStatus("FINISHED"),
}

// =============================================================================

var statuses = make(map[string]Status)

type Status struct {
	name string
}

func newStatus(role string) Status {
	r := Status{role}
	statuses[role] = r
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
	role, exists := statuses[value]
	if !exists {
		return Status{}, fmt.Errorf("invalid role %q", value)
	}

	return role, nil
}

func MustParseStatus(value string) Status {
	role, err := ParseStatus(value)
	if err != nil {
		panic(err)
	}

	return role
}
