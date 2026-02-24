package field

import (
	"fmt"

	"github.com/xiriframework/xiri-go/uicontext"
)

// GeoformField represents a geometry form field for geofencing/map drawing
// This is the 16th field type in the formfield system
type GeoformField struct {
	*BaseField
}

// GeoformValue represents a parsed geometry value
// Type: 1 = POLYGON, 2 = KREIS (circle)
// Note: Path is used differently for polygon vs circle:
// - Polygon (type 1): Path is array of {lat,lng} objects
// - Circle (type 2): Path is single {lat,lng,radius} object
type GeoformValue struct {
	Type int         `json:"type"` // 1 = POLYGON, 2 = KREIS
	Path interface{} `json:"path"` // For polygon: []map[string]string, for circle: map[string]string
}

func (f *GeoformField) Validate(value interface{}) error {
	if value == nil {
		if f.Required {
			return fmt.Errorf("geoform field %s is required", f.ID)
		}
		return nil
	}

	gv, ok := value.(*GeoformValue)
	if !ok {
		return fmt.Errorf("invalid geoform value type for %s", f.ID)
	}

	// Validate based on type
	if gv.Type == 1 {
		// Polygon validation - Path should be array of objects
		pathArr, ok := gv.Path.([]map[string]string)
		if !ok {
			return fmt.Errorf("polygon path must be an array of {lat,lng} objects")
		}
		if len(pathArr) < 3 {
			return fmt.Errorf("polygon must have at least 3 points")
		}
		for i, point := range pathArr {
			latStr, latOk := point["lat"]
			lngStr, lngOk := point["lng"]
			if !latOk || !lngOk {
				return fmt.Errorf("polygon point %d must have lat and lng fields", i)
			}

			// Parse lat/lng from string
			var lat, lng float64
			if _, err := fmt.Sscanf(latStr, "%f", &lat); err != nil {
				return fmt.Errorf("polygon point %d: invalid latitude: %w", i, err)
			}
			if _, err := fmt.Sscanf(lngStr, "%f", &lng); err != nil {
				return fmt.Errorf("polygon point %d: invalid longitude: %w", i, err)
			}

			if lat < -90 || lat > 90 {
				return fmt.Errorf("polygon point %d: latitude must be between -90 and 90", i)
			}
			if lng < -180 || lng > 180 {
				return fmt.Errorf("polygon point %d: longitude must be between -180 and 180", i)
			}
		}
	} else if gv.Type == 2 {
		// Circle validation - Path should be single object with lat, lng, radius
		pathObj, ok := gv.Path.(map[string]string)
		if !ok {
			return fmt.Errorf("circle path must be an object with {lat,lng,radius}")
		}

		latStr, latOk := pathObj["lat"]
		lngStr, lngOk := pathObj["lng"]
		radiusStr, radiusOk := pathObj["radius"]
		if !latOk || !lngOk || !radiusOk {
			return fmt.Errorf("circle path must have lat, lng, and radius fields")
		}

		// Parse from strings
		var lat, lng, radius float64
		if _, err := fmt.Sscanf(latStr, "%f", &lat); err != nil {
			return fmt.Errorf("circle: invalid latitude: %w", err)
		}
		if _, err := fmt.Sscanf(lngStr, "%f", &lng); err != nil {
			return fmt.Errorf("circle: invalid longitude: %w", err)
		}
		if _, err := fmt.Sscanf(radiusStr, "%f", &radius); err != nil {
			return fmt.Errorf("circle: invalid radius: %w", err)
		}

		if lat < -90 || lat > 90 {
			return fmt.Errorf("circle center latitude must be between -90 and 90")
		}
		if lng < -180 || lng > 180 {
			return fmt.Errorf("circle center longitude must be between -180 and 180")
		}
		if radius < 1 || radius > 50000 {
			return fmt.Errorf("circle radius must be between 1 and 50000 meters")
		}
	} else {
		return fmt.Errorf("invalid geometry type: %d (must be 1=polygon or 2=circle)", gv.Type)
	}

	return nil
}

func (f *GeoformField) Parse(raw interface{}) (interface{}, error) {
	if raw == nil {
		return f.GetDefault(), nil
	}

	// Expect map with "type" and geometry data
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("geoform field expects map with geometry data")
	}

	// Parse type field
	typeRaw, ok := m["type"]
	if !ok {
		return nil, fmt.Errorf("geoform field missing 'type' field")
	}

	var geoType int
	switch v := typeRaw.(type) {
	case int:
		geoType = v
	case float64:
		geoType = int(v)
	case string:
		// Handle string type conversion
		if v == "1" {
			geoType = 1
		} else if v == "2" {
			geoType = 2
		} else {
			return nil, fmt.Errorf("invalid type string: %s", v)
		}
	default:
		return nil, fmt.Errorf("invalid type field: %T", typeRaw)
	}

	gv := &GeoformValue{Type: geoType}

	if geoType == 1 {
		// Parse polygon - expect "path" array of objects
		pathRaw, ok := m["path"]
		if !ok {
			return nil, fmt.Errorf("polygon geometry missing 'path' field")
		}

		pathArr, ok := pathRaw.([]interface{})
		if !ok {
			return nil, fmt.Errorf("polygon 'path' must be an array")
		}

		// Create typed slice for polygon points
		points := make([]map[string]string, len(pathArr))
		for i, pointRaw := range pathArr {
			pointMap, ok := pointRaw.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("polygon point %d must be an object", i)
			}

			// Extract lat and lng
			latRaw, latOk := pointMap["lat"]
			lngRaw, lngOk := pointMap["lng"]
			if !latOk || !lngOk {
				return nil, fmt.Errorf("polygon point %d must have lat and lng fields", i)
			}

			// Convert to strings for precision
			latStr := fmt.Sprintf("%v", latRaw)
			lngStr := fmt.Sprintf("%v", lngRaw)

			points[i] = map[string]string{
				"lat": latStr,
				"lng": lngStr,
			}
		}
		gv.Path = points
	} else if geoType == 2 {
		// Parse circle - expect "path" object with {lat, lng, radius}
		pathRaw, ok := m["path"]
		if !ok {
			return nil, fmt.Errorf("circle geometry missing 'path' field")
		}

		pathMap, ok := pathRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("circle 'path' must be an object")
		}

		// Extract lat, lng, radius and convert to strings
		latRaw, latOk := pathMap["lat"]
		lngRaw, lngOk := pathMap["lng"]
		radiusRaw, radiusOk := pathMap["radius"]
		if !latOk || !lngOk || !radiusOk {
			return nil, fmt.Errorf("circle path must have lat, lng, and radius fields")
		}

		// Convert to strings
		latStr := fmt.Sprintf("%v", latRaw)
		lngStr := fmt.Sprintf("%v", lngRaw)
		radiusStr := fmt.Sprintf("%v", radiusRaw)

		gv.Path = map[string]string{
			"lat":    latStr,
			"lng":    lngStr,
			"radius": radiusStr,
		}
	} else {
		return nil, fmt.Errorf("invalid geometry type: %d", geoType)
	}

	return gv, nil
}

// ExportForFrontend exports the field for frontend rendering
func (f *GeoformField) ExportForFrontend(ctx *uicontext.UiContext, value interface{}) map[string]interface{} {
	// Get base export fields
	result := f.BaseField.GetBaseExport(ctx, nil)

	// Set type and subtype
	result["type"] = "geoform"
	result["subtype"] = "geoform"

	// Handle value - expect GeoformValue
	if value != nil {
		if gv, ok := value.(*GeoformValue); ok {
			if gv.Type == 1 {
				// Polygon - path is array of {lat,lng} objects
				result["value"] = map[string]interface{}{
					"type": 1,
					"path": gv.Path,
				}
			} else if gv.Type == 2 {
				// Circle - path is single {lat,lng,radius} object
				result["value"] = map[string]interface{}{
					"type": 2,
					"path": gv.Path,
				}
			}
		}
	}

	// If no value, provide empty structure (polygon with empty path)
	if result["value"] == nil {
		result["value"] = map[string]interface{}{
			"type": 1,
			"path": []interface{}{},
		}
	}

	return result
}

// ============================================================================
// Builder Functions
// ============================================================================

// NewGeoformField creates a new geometry form field for geofencing/map drawing
func NewGeoformField(id, name string, required bool) *GeoformField {
	return &GeoformField{
		BaseField: &BaseField{
			ID:       id,
			Type:     FieldTypeGeoform,
			Name:     name,
			Required: required,
			Form:     true,
		},
	}
}

// parseFloatFromMap parses a float value from a map key
func parseFloatFromMap(m map[string]interface{}, key string) (float64, error) {
	raw, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("missing key %s", key)
	}
	return parseFloat(raw)
}

// parseFloat parses a float value from interface{}
func parseFloat(raw interface{}) (float64, error) {
	switch v := raw.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		var f float64
		_, err := fmt.Sscanf(v, "%f", &f)
		if err != nil {
			return 0, fmt.Errorf("invalid float string: %s", v)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float", raw)
	}
}

// ============================================================================
// Chainable Setter Methods
// ============================================================================

// SetClass sets the CSS class for frontend styling
func (f *GeoformField) SetClass(class string) *GeoformField {
	f.BaseField.SetClass(class)
	return f
}

// SetHint sets the tooltip/help text for the field
func (f *GeoformField) SetHint(hint string) *GeoformField {
	f.BaseField.SetHint(hint)
	return f
}

// SetStep sets the step indicator for multi-step forms
func (f *GeoformField) SetStep(step int) *GeoformField {
	f.BaseField.SetStep(step)
	return f
}

// SetDisabled sets whether the field is disabled
func (f *GeoformField) SetDisabled(disabled bool) *GeoformField {
	f.BaseField.SetDisabled(disabled)
	return f
}

// SetAccess sets the access control permissions
func (f *GeoformField) SetAccess(access []string) *GeoformField {
	f.BaseField.SetAccess(access)
	return f
}

// SetScenario sets which scenarios this field applies to
func (f *GeoformField) SetScenario(scenario []string) *GeoformField {
	f.BaseField.SetScenario(scenario)
	return f
}

// SetForm sets whether to show in form
func (f *GeoformField) SetForm(form bool) *GeoformField {
	f.BaseField.SetForm(form)
	return f
}
