package chat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"
)

// SendInput ist der Body des Sende-POST: {"text": "...", "clientId": "..."}.
type SendInput struct {
	Text     string `json:"text"`
	ClientID string `json:"clientId"`
}

var (
	ErrEmpty    = errors.New("chat: message is empty")
	ErrTooLong  = errors.New("chat: message too long")
	ErrClientID = errors.New("chat: clientId too long")
	ErrCursor   = errors.New("chat: cursor must not contain line breaks")
)

const (
	maxClientID = 64
	maxBody     = 1 << 20
)

// ParseSend liest und prüft den Sende-Body: Text getrimmt und nicht leer, höchstens maxLength Zeichen
// (Runen; maxLength <= 0 = ohne Grenze), clientId höchstens 64 Byte. Der Body wird auf 1 MiB begrenzt.
func ParseSend(r io.Reader, maxLength int) (SendInput, error) {
	var in SendInput
	if err := json.NewDecoder(io.LimitReader(r, maxBody)).Decode(&in); err != nil {
		return SendInput{}, fmt.Errorf("chat: invalid body: %w", err)
	}
	in.Text = strings.TrimSpace(in.Text)
	switch {
	case in.Text == "":
		return SendInput{}, ErrEmpty
	case maxLength > 0 && utf8.RuneCountInString(in.Text) > maxLength:
		return SendInput{}, ErrTooLong
	case len(in.ClientID) > maxClientID:
		return SendInput{}, ErrClientID
	}
	return in, nil
}

// Cursor liefert den Stand, ab dem Änderungen gefragt sind: Last-Event-ID (SSE-Reconnect) vor ?cursor=.
func Cursor(r *http.Request) string {
	if id := r.Header.Get("Last-Event-ID"); id != "" {
		return id
	}
	return r.URL.Query().Get("cursor")
}

// Stream schreibt Seiten als Server-Sent Events. Woher neue Nachrichten kommen (Channel, DB-Notify, …),
// entscheidet die App; Stream kennt weder Pub/Sub noch Persistenz.
type Stream struct {
	w http.ResponseWriter
	f http.Flusher
}

var ErrNoFlush = errors.New("chat: stream needs a flushable ResponseWriter")

// NewStream prüft zuerst, ob der Writer flushen kann (echo.Response kann es), damit der Handler bei
// ErrNoFlush noch einen Fehlerstatus senden kann. Dann setzt er die SSE-Header und schreibt sie sofort.
func NewStream(w http.ResponseWriter) (*Stream, error) {
	f, ok := w.(http.Flusher)
	if !ok {
		return nil, ErrNoFlush
	}
	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("X-Accel-Buffering", "no") // nginx und Co.: nicht puffern
	w.WriteHeader(http.StatusOK)
	f.Flush()
	return &Stream{w: w, f: f}, nil
}

// Send schreibt eine Seite als Event; der Cursor wird zur Event-ID (Browser schickt ihn beim Reconnect
// als Last-Event-ID).
func (s *Stream) Send(p *Page) error {
	if strings.ContainsAny(p.cursor, "\r\n") {
		return ErrCursor
	}
	payload, err := json.Marshal(p.Print())
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if p.cursor != "" {
		fmt.Fprintf(&buf, "id: %s\n", p.cursor)
	}
	fmt.Fprintf(&buf, "data: %s\n\n", payload)
	return s.write(buf.Bytes())
}

// Ping schreibt einen Kommentar, damit Proxys die Verbindung nicht wegen Leerlauf schließen.
func (s *Stream) Ping() error { return s.write([]byte(": ping\n\n")) }

func (s *Stream) write(b []byte) error {
	if _, err := s.w.Write(b); err != nil {
		return err
	}
	s.f.Flush()
	return nil
}
