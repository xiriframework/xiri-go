package field

import "fmt"

// ValidationError is a user-facing validation failure with a stable code, so FormGroup can
// translate it via the TranslateFunc key "validation.<Code>". Placeholders in the translation
// are {field} and the Params keys. Error() stays the English text for callers without translation.
//
// Codes and params:
//
//	required                   —
//	min_length / max_length    {min} / {max}
//	pattern                    —
//	min / max                  {min} / {max}   (number fields)
//	min_items / max_items      {min} / {max}
//	not_allowed                —               (value is no offered option)
//	not_past / not_future      —               (time fields)
//	min_date / max_date        —               (time fields)
//	range_order                —               (time range start after end)
type ValidationError struct {
	Code   string
	Params map[string]any
	Msg    string
}

func (e *ValidationError) Error() string { return e.Msg }

func invalid(code string, params map[string]any, format string, args ...any) error {
	return &ValidationError{Code: code, Params: params, Msg: fmt.Sprintf(format, args...)}
}
