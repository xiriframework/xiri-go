package field

import (
	"strings"
	"testing"
)

func TestGeoformField_Validate_Required(t *testing.T) {
	f := NewGeoformField("geo", "GEO", true)
	err := f.Validate(nil)
	if err == nil {
		t.Fatal("expected error for nil on required geoform field")
	}
	if !strings.Contains(err.Error(), "required") {
		t.Errorf("expected 'required' in error, got: %v", err)
	}
}

func TestGeoformField_Validate_Optional_Nil(t *testing.T) {
	f := NewGeoformField("geo", "GEO", false)
	if err := f.Validate(nil); err != nil {
		t.Fatalf("unexpected error for nil on optional field: %v", err)
	}
}

func TestGeoformField_Validate_InvalidType(t *testing.T) {
	f := NewGeoformField("geo", "GEO", true)
	err := f.Validate("not a geoform value")
	if err == nil {
		t.Fatal("expected error for non-GeoformValue input")
	}
	if !strings.Contains(err.Error(), "invalid geoform value type") {
		t.Errorf("expected 'invalid geoform value type' in error, got: %v", err)
	}
}

func TestGeoformField_Validate_InvalidGeoType(t *testing.T) {
	f := NewGeoformField("geo", "GEO", true)
	gv := &GeoformValue{Type: 3, Path: nil}
	err := f.Validate(gv)
	if err == nil {
		t.Fatal("expected error for invalid geometry type 3")
	}
	if !strings.Contains(err.Error(), "invalid geometry type") {
		t.Errorf("expected 'invalid geometry type' in error, got: %v", err)
	}
}

func TestGeoformField_Validate_Polygon_Valid(t *testing.T) {
	f := NewGeoformField("geo", "GEO", true)
	gv := &GeoformValue{
		Type: 1,
		Path: []map[string]string{
			{"lat": "48.2082", "lng": "16.3738"},
			{"lat": "48.2100", "lng": "16.3800"},
			{"lat": "48.2050", "lng": "16.3700"},
			{"lat": "48.2000", "lng": "16.3650"},
		},
	}
	if err := f.Validate(gv); err != nil {
		t.Fatalf("unexpected error for valid polygon: %v", err)
	}
}

func TestGeoformField_Validate_Polygon_TooFewPoints(t *testing.T) {
	f := NewGeoformField("geo", "GEO", true)
	gv := &GeoformValue{
		Type: 1,
		Path: []map[string]string{
			{"lat": "48.2082", "lng": "16.3738"},
			{"lat": "48.2100", "lng": "16.3800"},
		},
	}
	err := f.Validate(gv)
	if err == nil {
		t.Fatal("expected error for polygon with fewer than 3 points")
	}
	if !strings.Contains(err.Error(), "at least 3 points") {
		t.Errorf("expected 'at least 3 points' in error, got: %v", err)
	}
}

func TestGeoformField_Validate_Polygon_InvalidLat(t *testing.T) {
	f := NewGeoformField("geo", "GEO", true)
	gv := &GeoformValue{
		Type: 1,
		Path: []map[string]string{
			{"lat": "91.0", "lng": "16.3738"},
			{"lat": "48.2100", "lng": "16.3800"},
			{"lat": "48.2050", "lng": "16.3700"},
		},
	}
	err := f.Validate(gv)
	if err == nil {
		t.Fatal("expected error for latitude > 90")
	}
	if !strings.Contains(err.Error(), "latitude must be between") {
		t.Errorf("expected 'latitude must be between' in error, got: %v", err)
	}
}

func TestGeoformField_Validate_Polygon_InvalidLng(t *testing.T) {
	f := NewGeoformField("geo", "GEO", true)
	gv := &GeoformValue{
		Type: 1,
		Path: []map[string]string{
			{"lat": "48.2082", "lng": "181.0"},
			{"lat": "48.2100", "lng": "16.3800"},
			{"lat": "48.2050", "lng": "16.3700"},
		},
	}
	err := f.Validate(gv)
	if err == nil {
		t.Fatal("expected error for longitude > 180")
	}
	if !strings.Contains(err.Error(), "longitude must be between") {
		t.Errorf("expected 'longitude must be between' in error, got: %v", err)
	}
}

func TestGeoformField_Validate_Circle_Valid(t *testing.T) {
	f := NewGeoformField("geo", "GEO", true)
	gv := &GeoformValue{
		Type: 2,
		Path: map[string]string{
			"lat":    "48.2082",
			"lng":    "16.3738",
			"radius": "500",
		},
	}
	if err := f.Validate(gv); err != nil {
		t.Fatalf("unexpected error for valid circle: %v", err)
	}
}

func TestGeoformField_Validate_Circle_InvalidRadius(t *testing.T) {
	f := NewGeoformField("geo", "GEO", true)
	gv := &GeoformValue{
		Type: 2,
		Path: map[string]string{
			"lat":    "48.2082",
			"lng":    "16.3738",
			"radius": "60000",
		},
	}
	err := f.Validate(gv)
	if err == nil {
		t.Fatal("expected error for radius > 50000")
	}
	if !strings.Contains(err.Error(), "radius must be between") {
		t.Errorf("expected 'radius must be between' in error, got: %v", err)
	}
}

func TestGeoformField_Validate_Circle_MissingRadius(t *testing.T) {
	f := NewGeoformField("geo", "GEO", true)
	gv := &GeoformValue{
		Type: 2,
		Path: map[string]string{
			"lat": "48.2082",
			"lng": "16.3738",
		},
	}
	err := f.Validate(gv)
	if err == nil {
		t.Fatal("expected error for missing radius key")
	}
	if !strings.Contains(err.Error(), "radius") {
		t.Errorf("expected 'radius' in error, got: %v", err)
	}
}

func TestGeoformField_Parse_Polygon(t *testing.T) {
	f := NewGeoformField("geo", "GEO", true)
	input := map[string]interface{}{
		"type": 1,
		"path": []interface{}{
			map[string]interface{}{"lat": 48.2082, "lng": 16.3738},
			map[string]interface{}{"lat": 48.2100, "lng": 16.3800},
			map[string]interface{}{"lat": 48.2050, "lng": 16.3700},
		},
	}

	result, err := f.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gv, ok := result.(*GeoformValue)
	if !ok {
		t.Fatalf("expected *GeoformValue, got %T", result)
	}
	if gv.Type != 1 {
		t.Errorf("expected type 1, got %d", gv.Type)
	}

	points, ok := gv.Path.([]map[string]string)
	if !ok {
		t.Fatalf("expected []map[string]string path, got %T", gv.Path)
	}
	if len(points) != 3 {
		t.Errorf("expected 3 points, got %d", len(points))
	}
	if points[0]["lat"] == "" || points[0]["lng"] == "" {
		t.Error("expected non-empty lat/lng in first point")
	}
}

func TestGeoformField_Parse_Circle(t *testing.T) {
	f := NewGeoformField("geo", "GEO", true)
	input := map[string]interface{}{
		"type": 2,
		"path": map[string]interface{}{"lat": 48.2082, "lng": 16.3738, "radius": 500.0},
	}

	result, err := f.Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gv, ok := result.(*GeoformValue)
	if !ok {
		t.Fatalf("expected *GeoformValue, got %T", result)
	}
	if gv.Type != 2 {
		t.Errorf("expected type 2, got %d", gv.Type)
	}

	circle, ok := gv.Path.(map[string]string)
	if !ok {
		t.Fatalf("expected map[string]string path, got %T", gv.Path)
	}
	if circle["lat"] == "" || circle["lng"] == "" || circle["radius"] == "" {
		t.Error("expected non-empty lat, lng, radius in circle")
	}
}

func TestGeoformField_Parse_InvalidInput(t *testing.T) {
	f := NewGeoformField("geo", "GEO", true)
	_, err := f.Parse("not a map")
	if err == nil {
		t.Fatal("expected error for non-map input")
	}
	if !strings.Contains(err.Error(), "expects map") {
		t.Errorf("expected 'expects map' in error, got: %v", err)
	}
}
