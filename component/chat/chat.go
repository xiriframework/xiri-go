// Package chat baut die Chat-Komponente (Typ "chat"): Nachrichtenverlauf mit Eingabe, neue Nachrichten unten,
// aktualisiert per Polling oder SSE. Protokoll: siehe skills/xiri-go-expert/references/components.md (Chat).
package chat

import (
	"maps"
	"time"

	"github.com/xiriframework/xiri-go/component/core"
)

// Status ist der Zustellstatus einer eigenen Nachricht. "sending" und "failed" vergibt nur der Client.
type Status string

const (
	Sent      Status = "sent"
	Delivered Status = "delivered"
	Read      Status = "read"
)

// Message ist eine Chat-Nachricht. Der Text ist Nutzerinhalt und wird nicht übersetzt.
type Message struct {
	id, text, from, avatar, clientID string
	time                             time.Time
	own                              bool
	status                           Status
}

// NewMessage erzeugt eine Nachricht, z. B. für SendResult.
func NewMessage(id, text string, t time.Time) *Message {
	return &Message{id: id, text: text, time: t}
}

// Own markiert die Nachricht als eigene (rechts, mit Status); from/avatar werden dann nicht ausgegeben.
func (m *Message) Own() *Message { m.own = true; return m }

// From setzt den Absendernamen (fremde Nachrichten).
func (m *Message) From(name string) *Message { m.from = name; return m }

// Avatar setzt die Bild-URL des Absenders; ohne Avatar zeigt der Client Initialen aus From.
func (m *Message) Avatar(url string) *Message { m.avatar = url; return m }

// Status setzt den Zustellstatus (nur bei Own ausgegeben).
func (m *Message) Status(s Status) *Message { m.status = s; return m }

// ClientID gibt die clientId des Sende-Requests zurück, damit der Client seine optimistische Nachricht ersetzt.
func (m *Message) ClientID(id string) *Message { m.clientID = id; return m }

func (m *Message) print() map[string]any {
	out := map[string]any{"id": m.id, "text": m.text, "time": m.time.Unix()}
	if m.own {
		out["own"] = true
		if m.status != "" {
			out["status"] = string(m.status)
		}
	} else {
		if m.from != "" {
			out["from"] = m.from
		}
		if m.avatar != "" {
			out["avatar"] = m.avatar
		}
	}
	if m.clientID != "" {
		out["clientId"] = m.clientID
	}
	return out
}

// SendResult ist die Antwort auf den Sende-POST: {"message": {...}}.
func SendResult(m *Message) map[string]any {
	return map[string]any{"message": m.print()}
}

// Page ist eine Antwort der Chat-Endpunkte (Startseite, ?cursor=, ?before=, SSE-Event).
// cursor ist für den Client undurchsichtig; er schickt ihn beim nächsten Abruf zurück.
type Page struct {
	messages []*Message
	cursor   string
	hasOlder bool
}

// NewPage erzeugt eine Seite mit dem Cursor, ab dem der nächste Abruf Änderungen liefern soll.
func NewPage(cursor string) *Page { return &Page{cursor: cursor} }

// HasOlder meldet, dass es vor der ältesten Nachricht dieser Seite weitere gibt.
func (p *Page) HasOlder(v bool) *Page { p.hasOlder = v; return p }

// Add hängt eine Nachricht an und liefert sie zum Weiterkonfigurieren.
func (p *Page) Add(id, text string, t time.Time) *Message {
	m := NewMessage(id, text, t)
	p.messages = append(p.messages, m)
	return m
}

// Print liefert {"messages": [...], "cursor": "...", "hasOlder": bool}.
func (p *Page) Print() map[string]any {
	msgs := make([]map[string]any, len(p.messages))
	for i, m := range p.messages {
		msgs[i] = m.print()
	}
	return map[string]any{"messages": msgs, "cursor": p.cursor, "hasOlder": p.hasOlder}
}

// Chat ist die Komponente. Ohne WithStream und WithPoll aktualisiert sie sich nur nach dem eigenen Senden.
type Chat struct {
	url, stream, placeholder, height, display string
	poll, maxLength                           int
	readonly                                  bool
	page                                      *Page
}

// New erzeugt einen Chat; url ist der Endpunkt für GET (Seiten) und POST (Senden).
func New(url string) *Chat { return &Chat{url: url} }

// WithStream schaltet auf SSE um (ersetzt Polling).
func (c *Chat) WithStream(url string) *Chat { c.stream = url; return c }

// WithPoll fragt alle seconds Sekunden nach Änderungen (ignoriert, wenn WithStream gesetzt ist).
func (c *Chat) WithPoll(seconds int) *Chat { c.poll = seconds; return c }

// WithPlaceholder setzt den Platzhalter der Eingabe (Übersetzungs-Key).
func (c *Chat) WithPlaceholder(key string) *Chat { c.placeholder = key; return c }

// WithMaxLength begrenzt die Eingabe; denselben Wert an ParseSend übergeben.
func (c *Chat) WithMaxLength(n int) *Chat { c.maxLength = n; return c }

// WithHeight setzt die Höhe (CSS), z. B. "60vh". Ohne füllt der Chat seinen Container.
func (c *Chat) WithHeight(css string) *Chat { c.height = css; return c }

// WithDisplay setzt die Layout-Klasse.
func (c *Chat) WithDisplay(display string) *Chat { c.display = display; return c }

// ReadOnly blendet die Eingabe aus (z. B. Archiv).
func (c *Chat) ReadOnly() *Chat { c.readonly = true; return c }

// WithPage liefert die Startseite gleich mit und spart den ersten GET.
func (c *Chat) WithPage(p *Page) *Chat { c.page = p; return c }

// Print liefert {"type": "chat", "data": {...}}.
func (c *Chat) Print(ctx *core.UiContext) map[string]any {
	data := map[string]any{"url": c.url}
	if c.stream != "" {
		data["stream"] = c.stream
	}
	if c.poll > 0 {
		data["poll"] = c.poll
	}
	if c.placeholder != "" {
		data["placeholder"] = core.Translate(ctx, c.placeholder)
	}
	if c.maxLength > 0 {
		data["maxLength"] = c.maxLength
	}
	if c.height != "" {
		data["height"] = c.height
	}
	if c.readonly {
		data["readonly"] = true
	}
	if c.page != nil {
		maps.Copy(data, c.page.Print())
	}
	out := map[string]any{"type": "chat", "data": data}
	if c.display != "" {
		out["display"] = c.display
	}
	return out
}
