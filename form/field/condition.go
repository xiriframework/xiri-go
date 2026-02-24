package field

// ConditionOperator defines the comparison operators for field conditions
type ConditionOperator string

const (
	CondEquals    ConditionOperator = "equals"
	CondNotEquals ConditionOperator = "notEquals"
	CondContains  ConditionOperator = "contains"
	CondGreater   ConditionOperator = "greaterThan"
	CondLess      ConditionOperator = "lessThan"
	CondIn        ConditionOperator = "in"
	CondNotEmpty  ConditionOperator = "notEmpty"
)

// Condition represents a visibility condition for a form field.
// When set on a field, the frontend will show/hide the field based on the value
// of another field in the same form.
type Condition struct {
	Field    string            `json:"field"`
	Operator ConditionOperator `json:"operator"`
	Value    interface{}       `json:"value,omitempty"`
}

// NewCondition creates a new condition
func NewCondition(field string, operator ConditionOperator, value interface{}) Condition {
	return Condition{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

// NewConditionNotEmpty creates a "notEmpty" condition for the given field
func NewConditionNotEmpty(field string) Condition {
	return Condition{
		Field:    field,
		Operator: CondNotEmpty,
	}
}
