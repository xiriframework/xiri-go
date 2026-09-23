package builder

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/xiriframework/xiri-go/form/field"
)

// BindSuggest handles a suggestion request (TextField.SetSuggestionsURL).
//
// It reads the JSON body once, returns the typed search text and binds ONLY the given context
// fields leniently: default first, then the posted value if present and usable. Disabled fields
// and a field with the reserved id "search" keep their default. Nothing else in the request is
// looked at - the context fields are the allowlist.
//
// Example:
//
//	func (ctrl *Controller) CitySearch(c echo.Context) error {
//	    dept := field.NewIntField("dept", "DEPT", false, 0)
//	    search, err := builder.BindSuggest(c, dept)
//	    if err != nil {
//	        return wc.BadRequest(err.Error())
//	    }
//	    return c.JSON(http.StatusOK, ctrl.citiesFor(*dept.Value, search))
//	}
func BindSuggest(c echo.Context, contextFields ...field.FormField) (string, error) {
	var raw map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&raw); err != nil {
		return "", err
	}
	search, _ := raw["search"].(string)
	for _, f := range contextFields {
		// Default first: a failed bind leaves the field on a usable value (like BindReloadFromMap).
		_ = bindFieldValue(f, nil)
		if f.IsDisabled() || f.GetID() == "search" {
			continue
		}
		if v, ok := raw[f.GetID()]; ok {
			_ = bindFieldValue(f, v)
		}
	}
	return search, nil
}
