package field

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"
)

// fileValue builds the value xiri-ng submits for a file field: the JSON-decoded
// [{name, data}] list with each file as a base64 data URL.
func fileValue(name string, size int) []interface{} {
	data := base64.StdEncoding.EncodeToString([]byte(strings.Repeat("x", size)))
	return []interface{}{map[string]interface{}{"name": name, "data": "data:application/pdf;base64," + data}}
}

// xiri-ng reads field.max and field.accept (file.component.ts / .html);
// the former maxSize/allowedTypes/allowedExtensions keys were never read.
func TestFileField_ExportForFrontend_MaxAndAccept(t *testing.T) {
	f := NewFileField("doc", "DOC", false, 1024)
	f.AllowedTypes = []string{"application/pdf", "image/png"}
	f.AllowedExtensions = []string{".pdf", ".png"}

	result := f.ExportForFrontend(nil, nil)

	if result["max"] != int64(1024) {
		t.Errorf("max = %#v, want int64(1024)", result["max"])
	}
	if result["accept"] != "application/pdf,image/png,.pdf,.png" {
		t.Errorf("accept = %#v", result["accept"])
	}
	for _, key := range []string{"maxSize", "allowedTypes", "allowedExtensions"} {
		if _, ok := result[key]; ok {
			t.Errorf("unexpected legacy key %q", key)
		}
	}
}

func TestFileField_ExportForFrontend_NoLimits(t *testing.T) {
	result := NewFileField("doc", "DOC", false, 0).ExportForFrontend(nil, nil)

	if _, ok := result["max"]; ok {
		t.Errorf("unexpected max without MaxSize: %#v", result["max"])
	}
	if _, ok := result["accept"]; ok {
		t.Errorf("unexpected accept without types/extensions: %#v", result["accept"])
	}
}

// The browser limit is client-side only; Validate enforces MaxSize per file on the server.
func TestFileField_Validate_MaxSize(t *testing.T) {
	f := NewFileField("doc", "DOC", false, 10)

	if err := f.Validate(fileValue("a.pdf", 10)); err != nil {
		t.Errorf("file of exactly MaxSize bytes: unexpected error %v", err)
	}

	err := f.Validate(fileValue("a.pdf", 11))
	var ve *ValidationError
	if !errors.As(err, &ve) || ve.Code != "max_size" || ve.Params["max"] != int64(10) {
		t.Errorf("file one byte over MaxSize: got %v, want max_size error with max=10", err)
	}

	// Far over the limit: rejected before decoding, same error.
	if err := f.Validate(fileValue("a.pdf", 100)); !errors.As(err, &ve) || ve.Code != "max_size" {
		t.Errorf("file far over MaxSize: got %v, want max_size error", err)
	}
}

func TestFileField_Validate_MaxSizePerFile(t *testing.T) {
	f := NewFileField("doc", "DOC", false, 10)
	f.Multiple = true

	ok := append(fileValue("a.pdf", 10), fileValue("b.pdf", 10)...)
	if err := f.Validate(ok); err != nil {
		t.Errorf("two files each within MaxSize: unexpected error %v", err)
	}
	tooBig := append(fileValue("a.pdf", 10), fileValue("b.pdf", 11)...)
	if err := f.Validate(tooBig); err == nil {
		t.Error("second file over MaxSize: expected error")
	}
}

func TestFileField_Validate_NoMaxSize(t *testing.T) {
	if err := NewFileField("doc", "DOC", false, 0).Validate(fileValue("a.pdf", 5000)); err != nil {
		t.Errorf("MaxSize 0 means no limit, got %v", err)
	}
}

func TestFileField_Validate_Malformed(t *testing.T) {
	f := NewFileField("doc", "DOC", false, 10)
	cases := map[string]interface{}{
		"string":            "somefile.pdf",
		"no list entry map": []interface{}{"a.pdf"},
		"data not string":   []interface{}{map[string]interface{}{"name": "a.pdf", "data": 5}},
		"no base64 marker":  []interface{}{map[string]interface{}{"name": "a.pdf", "data": "data:text/plain,abc"}},
		"broken base64":     []interface{}{map[string]interface{}{"name": "a.pdf", "data": "data:x;base64,!!!"}},
		"no data: prefix":   []interface{}{map[string]interface{}{"name": "a.pdf", "data": "junk;base64,YQ=="}},
		"no name":           []interface{}{map[string]interface{}{"data": "data:x;base64,YQ=="}},
	}
	for name, value := range cases {
		if err := f.Validate(value); err == nil {
			t.Errorf("%s: expected error, got nil", name)
		}
	}
}

func TestFileField_Validate_RequiredEmptyList(t *testing.T) {
	err := NewFileField("doc", "DOC", true, 0).Validate([]interface{}{})
	var ve *ValidationError
	if !errors.As(err, &ve) || ve.Code != "required" {
		t.Errorf("required with empty list: got %v, want required error", err)
	}
	if err := NewFileField("doc", "DOC", false, 0).Validate([]interface{}{}); err != nil {
		t.Errorf("optional with empty list: unexpected error %v", err)
	}
}
