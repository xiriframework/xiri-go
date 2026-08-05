package builder

import (
	"strings"

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
//	    builder := NewFormBuilder(ctx).
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
	formData, err := extractFormData(c, fg)
	if err != nil {
		return err
	}

	return BindFromMap(formData, fg)
}

// extractFormData pulls the declared fields out of the request body.
//
// Only Form=true and non-disabled fields are read; any extra key sent by the client is
// dropped, which is what prevents over-posting. Supports JSON and form-urlencoded/multipart.
func extractFormData(c echo.Context, fg *group.FormGroup) (map[string]interface{}, error) {
	allFields := fg.GetFields()
	fieldIDs := make([]string, 0, len(allFields))
	for _, f := range allFields {
		if f.GetForm() && !f.IsDisabled() {
			fieldIDs = append(fieldIDs, f.GetID())
		}
	}

	// Extract ONLY declared fields from request
	formData := make(map[string]interface{})

	// Check content type
	contentType := c.Request().Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		// Handle JSON body
		var rawData map[string]interface{}
		if err := c.Bind(&rawData); err != nil {
			return nil, err
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
			return nil, err
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

	return formData, nil
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
		var rawValue interface{}
		if f.GetForm() && !f.IsDisabled() {
			rawValue = resolveFieldValue(f, formData)
		} else {
			rawValue = f.GetDefault()
		}
		if err := bindFieldValue(f, rawValue); err != nil {
			return err
		}
	}
	return nil
}

// BindReload binds a reload request to the field instances.
//
// A reload fires while the user is still editing, so an empty required field or a value the
// user has not finished typing is normal input, not an error. Unlike BindAndValidate this
// therefore never fails on a single field - only an unreadable request body is an error.
//
// Use this in the handler behind a field's SetReloadOn URL:
//
//	func (ctrl *Controller) FormReload(c echo.Context) error {
//	    status, tags, fg := ctrl.buildThingForm(wc.UiContext())
//	    if err := builder.BindReload(c, fg); err != nil {
//	        return wc.BadRequest(err.Error())
//	    }
//	    tags.Options = ctrl.tagsForStatus(status.Value)
//	    return c.JSON(http.StatusOK, response.NewReturnFields(fg.ExportPatch()))
//	}
func BindReload(c echo.Context, fg *group.FormGroup) error {
	formData, err := extractFormData(c, fg)
	if err != nil {
		return err
	}

	return BindReloadFromMap(formData, fg)
}

// BindReloadFromMap is the lenient counterpart to BindFromMap.
//
// Every field is bound to its default first, so a field whose request value turns out to be
// unusable still holds the default afterwards instead of a nil/zero value. A failure on one
// field does not stop the remaining ones - any of them could be the trigger the server needs.
func BindReloadFromMap(formData map[string]interface{}, fg *group.FormGroup) error {
	for _, f := range fg.GetFields() {
		// Default first: a later failure leaves the field on a usable value.
		_ = bindFieldValue(f, f.GetDefault())

		if !f.GetForm() || f.IsDisabled() {
			continue
		}
		if rawValue, exists := formData[f.GetID()]; exists {
			_ = bindFieldValue(f, rawValue)
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
