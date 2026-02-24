package response

import (
	"encoding/json"
	"testing"
)

func TestNewReturnRefreshPage(t *testing.T) {
	r := NewReturnRefreshPage()

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `{"done":true,"refresh":"page"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestNewReturnRefreshTable(t *testing.T) {
	r := NewReturnRefreshTable()

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `{"done":true,"refresh":"table"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestNewReturnGoto(t *testing.T) {
	r := NewReturnGoto("/Portal/User/Page/7")

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `{"done":true,"goto":"/Portal/User/Page/7"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestNewReturnDone(t *testing.T) {
	r := NewReturnDone()

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `{"done":true}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestSuccessResponseInterface(t *testing.T) {
	// Verify all types implement SuccessResponse
	var _ SuccessResponse = NewReturnRefreshPage()
	var _ SuccessResponse = NewReturnRefreshTable()
	var _ SuccessResponse = NewReturnGoto("/test")
	var _ SuccessResponse = NewReturnDone()
	var _ SuccessResponse = NewReturnMessage("test", MessageSuccess)
}

func TestReturnRefreshPageWithMessage(t *testing.T) {
	r := NewReturnRefreshPage().WithMessage("Page reloaded", MessageSuccess)

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `{"done":true,"refresh":"page","message":"Page reloaded","messageType":"success"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestReturnRefreshTableWithMessage(t *testing.T) {
	r := NewReturnRefreshTable().WithMessage("Row deleted", MessageWarning)

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `{"done":true,"refresh":"table","message":"Row deleted","messageType":"warning"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestReturnGotoWithMessage(t *testing.T) {
	r := NewReturnGoto("/Dashboard").WithMessage("Saved successfully", MessageInfo)

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `{"done":true,"goto":"/Dashboard","message":"Saved successfully","messageType":"info"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestReturnDoneWithMessage(t *testing.T) {
	r := NewReturnDone().WithMessage("Operation complete", MessageSuccess)

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `{"done":true,"message":"Operation complete","messageType":"success"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestNewReturnMessage(t *testing.T) {
	r := NewReturnMessage("Settings saved", MessageSuccess)

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `{"done":true,"message":"Settings saved","messageType":"success"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestNewReturnSuccess(t *testing.T) {
	r := NewReturnSuccess("All good")

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `{"done":true,"message":"All good","messageType":"success"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestNewReturnError(t *testing.T) {
	r := NewReturnError("Something failed")

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	expected := `{"done":true,"message":"Something failed","messageType":"error"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}
