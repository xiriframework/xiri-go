package builder

import (
	"github.com/labstack/echo/v4"
	"github.com/xiriframework/xiri-go/form/field"
	"github.com/xiriframework/xiri-go/form/group"
)

// BindAndValidate safely extracts form data for declared fields only,
// then parses, validates, and stores values directly in the field instances.
//
// This function prevents over-posting attacks by only processing fields
// that are explicitly declared in the form definition. Any extra fields
// sent by the client are ignored.
//
// After calling this function, access values via the field instances:
//
//	name.Value, groupID.Value, enabled.Value
//
// Supports both JSON and form-urlencoded/multipart request bodies.
//
// Example usage in controller:
//
//	func (ctrl *Controller) AddSave(c echo.Context) error {
//	    // Create field instances
//	    name := field.NewTextField("name", "NAME", true)
//	    groupID := field.NewModelField("grp", "GRUPPE", true, "Group", int32(ctx.GroupID))
//	    enabled := field.NewBoolField("enabled", "ENABLED", false, true)
//
//	    // Build form
//	    builder := NewFormBuilder(ctx, t).
//	        AddField(name).
//	        AddField(groupID).
//	        AddField(enabled)
//
//	    fg, _, _ := builder.BuildAdd()
//
//	    // Bind and validate - values stored in field instances
//	    if err := formhelper.BindAndValidate(c, fg); err != nil {
//	        return wc.BadRequest(err.Error())
//	    }
//
//	    // Type-safe access - NO type assertions!
//	    einsatz := &core.Einsatz{
//	        Name:    name.Value,      // *string (compiler-checked)
//	        GroupId: groupID.Value,   // *int32 (compiler-checked)
//	        Enabled: enabled.Value,   // *bool (compiler-checked)
//	    }
//
//	    dbm.Einsatz.Create(einsatz)
//	    return wc.Goto("/Portal/Usage/Table")
//	}
//
// Returns:
//   - error: Validation error if any field fails validation
func BindAndValidate(c echo.Context, fg *group.FormGroup) error {
	// Get all declared field IDs from FormGroup
	fieldIDs := fg.GetFieldIDs()

	// Extract ONLY declared fields from request
	formData := make(map[string]interface{})

	// Check content type
	contentType := c.Request().Header.Get("Content-Type")

	if contentType == "application/json" || contentType == "application/json; charset=UTF-8" {
		// Handle JSON body
		var rawData map[string]interface{}
		if err := c.Bind(&rawData); err != nil {
			return err
		}

		// Filter to declared fields only (prevents over-posting)
		for _, fieldID := range fieldIDs {
			if value, exists := rawData[fieldID]; exists {
				formData[fieldID] = value
			}
		}
	} else {
		// Handle form-urlencoded / multipart
		// Parse form first
		if err := c.Request().ParseForm(); err != nil {
			return err
		}

		for _, fieldID := range fieldIDs {
			// Check if field exists in form
			if values, exists := c.Request().Form[fieldID]; exists && len(values) > 0 {
				if len(values) == 1 {
					// Single value field
					formData[fieldID] = values[0]
				} else {
					// Array field (e.g., ModelListField)
					formData[fieldID] = values
				}
			} else if value := c.FormValue(fieldID); value != "" {
				// Fallback to FormValue for single values
				formData[fieldID] = value
			}
		}
	}

	return BindFromMap(formData, fg)
}

// BindFromMap binds values from an already-parsed map to field instances.
//
// Each field's BindValue() method is called, which:
//   - Parses the raw value (e.g., float64 -> int32)
//   - Validates against field rules
//   - Stores in field.Value property
//
// After this function returns, access values via field.Value (type-safe).
func BindFromMap(formData map[string]interface{}, fg *group.FormGroup) error {
	for _, f := range fg.GetFields() {
		rawValue := resolveFieldValue(f, formData)
		if err := bindFieldValue(f, rawValue); err != nil {
			return err
		}
	}
	return nil
}

// resolveFieldValue extracts the raw value for a field from form data,
// falling back to the field's default if the field is not present.
func resolveFieldValue(field field.FormField, formData map[string]interface{}) interface{} {
	rawValue, exists := formData[field.GetID()]
	if !exists {
		return field.GetDefault()
	}
	return rawValue
}

// bindFieldValue binds a single raw value to a field instance.
// It prefers the BindValue method if available, falling back to Parse+Validate.
func bindFieldValue(field field.FormField, rawValue interface{}) error {
	type ValueBinder interface {
		BindValue(raw interface{}) error
	}

	if binder, ok := field.(ValueBinder); ok {
		return binder.BindValue(rawValue)
	}

	// Fallback for fields that don't have BindValue yet
	parsed, err := field.Parse(rawValue)
	if err != nil {
		return err
	}
	return field.Validate(parsed)
}
