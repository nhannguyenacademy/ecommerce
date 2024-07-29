package auth

import (
	_ "embed"
	"fmt"
)

// Rule represents an authorize rule in the system.
type Rule struct {
	name string
}

// Equal provides support for the go-cmp package and testing.
func (r Rule) Equal(r2 Rule) bool {
	return r.name == r2.name
}

// Set of known rules.
var rules = make(map[string]Rule)

func newRule(rule string) Rule {
	r := Rule{rule}
	rules[rule] = r
	return r
}

type ruleSet struct {
	Any            Rule
	Admin          Rule
	User           Rule
	AdminOrSubject Rule
}

var Rules = ruleSet{
	Any:            newRule("rule_any"),
	Admin:          newRule("rule_admin_only"),
	User:           newRule("rule_user_only"),
	AdminOrSubject: newRule("rule_admin_or_subject"),
}

// =============================================================================

func ParseRule(value string) (Rule, error) {
	role, exists := rules[value]
	if !exists {
		return Rule{}, fmt.Errorf("invalid rule %q", value)
	}

	return role, nil
}

func MustParseRule(value string) Rule {
	role, err := ParseRule(value)
	if err != nil {
		panic(err)
	}

	return role
}
