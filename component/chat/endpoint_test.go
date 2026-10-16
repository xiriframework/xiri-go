package chat

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestParseSend(t *testing.T) {
	tests := []struct {
		name, body string
		max        int
		want       SendInput
		err        error
	}{
		{"trimmed", `{"text":"  Hallo \n","clientId":"c-1"}`, 0, SendInput{Text: "Hallo", ClientID: "c-1"}, nil},
		{"no clientId ok", `{"text":"Hi"}`, 0, SendInput{Text: "Hi"}, nil},
		{"empty", `{"text":"   ","clientId":"c"}`, 0, SendInput{}, ErrEmpty},
		{"too long in runes", `{"text":"äöü"}`, 2, SendInput{}, ErrTooLong},
		{"exactly max", `{"text":"äöü"}`, 3, SendInput{Text: "äöü"}, nil},
		{"clientId too long", `{"text":"x","clientId":"` + strings.Repeat("a", 65) + `"}`, 0, SendInput{}, ErrClientID},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSend(strings.NewReader(tt.body), tt.max)
			if !errors.Is(err, tt.err) || got != tt.want {
				t.Errorf("got %+v, %v; want %+v, %v", got, err, tt.want, tt.err)
			}
		})
	}
	if _, err := ParseSend(strings.NewReader(`{kaputt`), 0); err == nil {
		t.Error("invalid JSON must fail")
	}
}

func TestCursor(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/chat?cursor=c-3", nil)
	if got := Cursor(r); got != "c-3" {
		t.Errorf("query cursor = %q", got)
	}
	r.Header.Set("Last-Event-ID", "c-9")
	if got := Cursor(r); got != "c-9" {
		t.Errorf("Last-Event-ID must win, got %q", got)
	}
}

func TestStreamSend(t *testing.T) {
	rec := httptest.NewRecorder()
	s, err := NewStream(rec)
	if err != nil {
		t.Fatal(err)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/event-stream" {
		t.Errorf("Content-Type = %q", ct)
	}
	if rec.Header().Get("Cache-Control") != "no-cache" {
		t.Error("Cache-Control no-cache missing")
	}

	rec.Flushed = false // NewStream flusht schon die Header; hier zählt nur der Flush von Send
	p := NewPage("c-2")
	p.Add("m1", "Hi\nzweite Zeile", time.Unix(1, 0)).From("Anna")
	if err := s.Send(p); err != nil {
		t.Fatal(err)
	}
	body := rec.Body.String()
	prefix := "id: c-2\ndata: "
	if !strings.HasPrefix(body, prefix) || !strings.HasSuffix(body, "\n\n") {
		t.Fatalf("event format wrong: %q", body)
	}
	payload := strings.TrimSuffix(strings.TrimPrefix(body, prefix), "\n\n")
	if strings.Contains(payload, "\n") {
		t.Fatalf("data must be one line: %q", payload)
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(payload), &got); err != nil || got["cursor"] != "c-2" {
		t.Errorf("payload = %v, err %v", got, err)
	}
	if !rec.Flushed {
		t.Error("Send must flush")
	}
}

func TestStreamPingAndNoCursor(t *testing.T) {
	rec := httptest.NewRecorder()
	s, _ := NewStream(rec)
	_ = s.Ping()
	_ = s.Send(NewPage(""))
	if got := rec.Body.String(); !strings.HasPrefix(got, ": ping\n\ndata: ") {
		t.Errorf("body = %q (no id line without cursor)", got)
	}
}

func TestStreamRejectsCursorWithNewline(t *testing.T) {
	s, _ := NewStream(httptest.NewRecorder())
	if err := s.Send(NewPage("a\nb")); !errors.Is(err, ErrCursor) {
		t.Errorf("err = %v, want ErrCursor", err)
	}
}

type plainWriter struct{ h http.Header }

func (w *plainWriter) Header() http.Header         { return w.h }
func (w *plainWriter) Write(b []byte) (int, error) { return len(b), nil }
func (w *plainWriter) WriteHeader(int)             {}

func TestNewStreamNeedsFlusher(t *testing.T) {
	w := &plainWriter{h: http.Header{}}
	if _, err := NewStream(w); !errors.Is(err, ErrNoFlush) {
		t.Errorf("err = %v, want ErrNoFlush", err)
	}
	if w.h.Get("Content-Type") != "" {
		t.Error("headers must stay untouched, so the handler can still send an error status")
	}
}
