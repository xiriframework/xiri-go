package response

import (
	"encoding/json"
	"testing"
)

// ReturnFields must be usable everywhere a success response is expected.
var _ SuccessResponse = ReturnFields{}

func TestNewReturnFields(t *testing.T) {
	r := NewReturnFields(map[string]interface{}{
		"tags": map[string]interface{}{"required": true},
	})

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `{"fields":{"tags":{"required":true}}}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

// A field patch is not a completed action - a "done" would be read as one by the frontend's
// form service, which turns any done-shaped response into a finished-form signal.
func TestNewReturnFields_HasNoDone(t *testing.T) {
	data, err := json.Marshal(NewReturnFields(map[string]interface{}{}))
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if _, present := decoded["done"]; present {
		t.Errorf("expected no done key, got %s", string(data))
	}
}

func TestReturnFields_WithMessage(t *testing.T) {
	r := NewReturnFields(map[string]interface{}{}).WithMessage("Liste aktualisiert", MessageInfo)

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `{"fields":{},"message":"Liste aktualisiert","messageType":"info"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}
