// Package mcp stellt eine Xiri-App als MCP-Server (Streamable HTTP) bereit.
//
// Es gibt bewusst nur zwei generische Tools, weil das Xiri-Seiten-JSON bereits beschreibt, was ein
// Agent tun kann: read_page liefert die Seite samt Tabellen, Formularen (fields, url) und Buttons
// (action, url, data); act führt eine API-Aktion daraus aus und gibt die Response
// (done/goto/refresh/message oder 400 mit error/fields) zurück. Kein Action-Registry, keine
// Protokolländerung.
//
// Mount (der Host sichert den Pfad mit derselben Auth-Middleware wie /api):
//
//	e.Any("/mcp", echo.WrapHandler(mcp.Handler(e, mcp.Options{Name: "meine-app", Version: "1.0"})))
//
// # Umfang und Rechte
//
// Tool-Aufrufe laufen als interner http.Request durch dasselbe Echo, also durch Routing und
// Middleware des Hosts. Der Agent handelt damit unter der Session, deren Cookie bzw.
// Authorization-Header weitergereicht wird, und erreicht jeden Endpunkt unterhalb von ApiPrefix —
// nicht nur die Aktionen der zuletzt gelesenen Seite. Das Seiten-JSON ist der Katalog, aus dem der
// Agent URLs liest, es ist keine Schranke. Wer weniger freigeben will, mountet einen eigenen Echo
// mit weniger Routen.
//
// Das entspricht den Rechten des Browsers nur, solange der Host jeden Endpunkt anhand der
// weitergereichten Credentials im selben Echo autorisiert. Prüfungen, die am Reverse-Proxy oder an
// Request-Metadaten (Host, TLS, RemoteAddr) hängen, greifen hier nicht — siehe Grenzen.
//
// ReadOnlyHint an read_page setzt voraus, was in Xiri gilt: Seitenrouten sind GET und
// nebenwirkungsfrei. Ein Host mit schreibendem GET bricht dieselbe Annahme auch im Browser.
//
// # Grenzen
//
//   - Host-basiertes Routing und HTTPSRedirect-Middleware am selben Echo funktionieren nicht: der
//     interne Request hat weder Host noch TLS, und das SDK liefert den äußeren Request nicht aus
//     (RequestExtra kennt nur TokenInfo, Header, CloseSSEStream).
//   - Globale Echo-Middleware läuft zweimal, für den MCP-Request und den internen Request. Ein
//     globales Parallelitätslimit mit Kapazität 1 blockiert sich damit selbst.
//   - Middleware am vorgelagerten Reverse-Proxy läuft für den internen Request nicht mit.
//   - Antwort-Header gehen verloren, auch Set-Cookie: eine Session-Rotation des Hosts erreicht den
//     MCP-Client nicht. Nur Status und Location werden gemeldet.
//   - Es gibt keinen eigenen Timeout: e.ServeHTTP ist synchron. Der MCP-Kontext wird
//     durchgereicht, kooperative Handler brechen ab; ein unkooperativer hängt den Tool-Call.
//   - Streaming wird nicht unterstützt: die Antwort wird gepuffert (begrenzt durch
//     MaxResponseBytes) und erst am Ende zurückgegeben.
//   - Binäre Antworten (z. B. Excel-Export) werden abgelehnt statt zerstört; dafür bleibt die UI
//     zuständig.
package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Options konfiguriert den MCP-Server. Leere Felder bekommen Defaults (siehe withDefaults).
type Options struct {
	Name         string // Servername im initialize-Handshake, Default "xiri"
	Version      string // Default "0"
	Instructions string // Hinweise an den Client; Default beschreibt das Xiri-Protokoll
	// ApiPrefix ist der Pfad, unter dem die Xiri-Handler am Echo hängen. Default "/api/";
	// fehlende Slashes werden ergänzt. Ziele ausserhalb werden abgelehnt.
	ApiPrefix string
	// ForwardHeaders werden vom MCP-Request in den internen Request kopiert. Default Cookie,
	// Authorization. Transport- und Body-Header (siehe blockedHeaders) werden nie kopiert.
	// Bewusst klein halten: jeder Eintrag erweitert, was ein MCP-Client in die App hineinreicht.
	ForwardHeaders []string
	// MaxResponseBytes begrenzt, wieviel vom Response-Body gepuffert wird. Default 1 MiB.
	// Schreibt ein Handler mehr, bekommt er einen Schreibfehler und der Aufruf wird abgelehnt.
	MaxResponseBytes int
}

const defaultInstructions = `Xiri-App. read_page(route) liefert das Seiten-JSON: Komponenten mit "type" ` +
	`(table, form, card, buttonline, …). Formulare haben data.url und data.fields (id, type, required, showWhen); ` +
	`Buttons haben action, url und optional data. act(url, method, data) führt aus, was serverseitig etwas tut: ` +
	`Formular-Submit (POST der Feldwerte an data.url) und Buttons mit action "api". Buttons mit action "link", ` +
	`"href" oder "back" sind Navigation — dafür read_page auf die Ziel-Route statt act; "download" liefert eine ` +
	`Datei und wird nicht unterstützt; "close", "return", "debug", "simulate" sind reine Client-Aktionen. ` +
	`Bei action "dialog" liefert die url den Dialog-Inhalt: ohne data und ohne filter per read_page, ` +
	`mit Filter- oder Initialdaten per act(url, data). ` +
	`Antworten: {"done":true}, {"goto":"/Route"} (danach read_page auf die Route ohne führenden Slash), ` +
	`{"refresh":"page|table|panel"}, optional message; bei 400 {"error":"…","fields":{"<fieldId>":"…"}}. ` +
	`Status ≠ 200 steht als "HTTP <code>" in der ersten Zeile.`

// blockedHeaders werden nie aus dem MCP-Request übernommen: sie beschreiben den Transport bzw. den
// Body des internen Requests, den dieses Paket selbst baut.
var blockedHeaders = map[string]bool{
	"Host": true, "Content-Length": true, "Content-Type": true,
	"Transfer-Encoding": true, "Connection": true,
}

func withDefaults(o Options) Options {
	if o.Name == "" {
		o.Name = "xiri"
	}
	if o.Version == "" {
		o.Version = "0"
	}
	if o.Instructions == "" {
		o.Instructions = defaultInstructions
	}
	if o.ApiPrefix == "" {
		o.ApiPrefix = "/api/"
	}
	if !strings.HasPrefix(o.ApiPrefix, "/") {
		o.ApiPrefix = "/" + o.ApiPrefix
	}
	if !strings.HasSuffix(o.ApiPrefix, "/") {
		o.ApiPrefix += "/"
	}
	if o.ForwardHeaders == nil {
		o.ForwardHeaders = []string{"Cookie", "Authorization"}
	}
	if o.MaxResponseBytes <= 0 {
		o.MaxResponseBytes = 1 << 20
	}
	return o
}

// ReadPageInput ist das Argument von read_page.
type ReadPageInput struct {
	Route string `json:"route" jsonschema:"Xiri-Route ohne führenden Slash, z. B. Portal/Devices"`
}

// ActInput ist das Argument von act.
type ActInput struct {
	URL    string         `json:"url" jsonschema:"url eines api-Buttons oder Formulars aus read_page"`
	Method string         `json:"method,omitempty" jsonschema:"GET oder POST, Default POST"`
	Data   map[string]any `json:"data,omitempty" jsonschema:"Formularwerte bzw. button.data; nur bei POST erlaubt"`
}

// Handler baut den MCP-HTTP-Handler über dem gegebenen Echo. Tool-Aufrufe laufen intern durch
// e.ServeHTTP, also durch Middleware und Routing des Hosts.
func Handler(e *echo.Echo, opts Options) http.Handler {
	opts = withDefaults(opts)
	srv := sdk.NewServer(&sdk.Implementation{Name: opts.Name, Version: opts.Version},
		&sdk.ServerOptions{Instructions: opts.Instructions})

	sdk.AddTool(srv, &sdk.Tool{
		Name:        "read_page",
		Description: "Liest eine Xiri-Seite (GET " + opts.ApiPrefix + "<route>) und liefert ihr JSON: Komponenten, Formulare mit Feldern und URL, Buttons mit action/url/data.",
		Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, req *sdk.CallToolRequest, in ReadPageInput) (*sdk.CallToolResult, any, error) {
		return dispatch(ctx, e, opts, req, http.MethodGet, in.Route, nil)
	})

	sdk.AddTool(srv, &sdk.Tool{
		Name:        "act",
		Description: "Führt eine Xiri-Aktion aus: POST (Default) der Feldwerte an die url eines Formulars oder eines Buttons mit action \"api\", alternativ GET. Liefert die Xiri-Response (done/goto/refresh/message; bei 400 error und fields). Navigation (link/href/back) und Downloads gehören nicht hierher.",
	}, func(ctx context.Context, req *sdk.CallToolRequest, in ActInput) (*sdk.CallToolResult, any, error) {
		method := strings.ToUpper(in.Method)
		if method == "" {
			method = http.MethodPost
		}
		if method != http.MethodGet && method != http.MethodPost {
			return errorResult("method must be GET or POST"), nil, nil
		}
		return dispatch(ctx, e, opts, req, method, in.URL, in.Data)
	})

	return sdk.NewStreamableHTTPHandler(func(*http.Request) *sdk.Server { return srv },
		&sdk.StreamableHTTPOptions{Stateless: true})
}

// dispatch schickt einen internen Request durch das Echo und verpackt den Body als Tool-Ergebnis.
func dispatch(ctx context.Context, e *echo.Echo, opts Options, req *sdk.CallToolRequest, method, route string, data map[string]any) (*sdk.CallToolResult, any, error) {
	if method == http.MethodGet && len(data) > 0 {
		return errorResult("data is only supported with POST; put parameters into the url"), nil, nil
	}
	var body io.Reader
	if method == http.MethodPost {
		if data == nil {
			data = map[string]any{}
		}
		b, err := json.Marshal(data)
		if err != nil {
			return errorResult("data is not JSON-serialisable: " + err.Error()), nil, nil
		}
		body = bytes.NewReader(b)
	}
	r, err := http.NewRequestWithContext(ctx, method, opts.ApiPrefix+strings.TrimPrefix(route, "/"), body)
	if err != nil {
		return errorResult("invalid url: " + err.Error()), nil, nil
	}
	// Das Ziel muss unterhalb des Prefix bleiben: fängt "..". Absolute URLs (https://host/x) landen
	// als harmloser interner Pfad /api/https:/host/x, nicht als externer Aufruf — e.ServeHTTP
	// spricht nie mit dem Netz. Der Query-String bleibt in r.URL.RawQuery erhalten.
	clean := path.Clean(r.URL.Path)
	if r.URL.Host != "" || !strings.HasPrefix(clean+"/", opts.ApiPrefix) {
		return errorResult("url must stay below " + opts.ApiPrefix), nil, nil
	}
	// Echo unterscheidet /x und /x/, path.Clean nicht — der Slash gehört zur Route.
	if strings.HasSuffix(r.URL.Path, "/") && !strings.HasSuffix(clean, "/") {
		clean += "/"
	}
	// Geroutet wird genau der Pfad, der geprüft wurde — sonst sähe Echo eine andere Schreibweise.
	r.URL.Path, r.URL.RawPath = clean, ""
	if req.Extra != nil {
		for _, h := range opts.ForwardHeaders {
			key := http.CanonicalHeaderKey(h)
			if blockedHeaders[key] {
				continue
			}
			// Values statt Get: mehrere Cookie-Zeilen bleiben erhalten. Kopie, damit innere
			// Middleware nicht in den Header des äußeren Requests schreibt.
			if v := req.Extra.Header.Values(key); len(v) > 0 {
				r.Header[key] = append([]string(nil), v...)
			}
		}
	}
	if body != nil {
		r.Header.Set("Content-Type", "application/json") // nach dem Forward, damit nichts ihn überschreibt.
	}

	rec := &recorder{header: http.Header{}, limit: opts.MaxResponseBytes}
	e.ServeHTTP(rec, r)
	// Ein Handler, der nur nil zurückgibt, schreibt nichts — ohne das meldete das Tool "HTTP 0".
	rec.WriteHeader(http.StatusOK)

	raw := rec.buf.Bytes()
	if rec.truncated {
		return errorResult(fmt.Sprintf("response exceeds %d bytes; use the UI for this endpoint", opts.MaxResponseBytes)), nil, nil
	}
	ct := rec.header.Get("Content-Type")
	if ct == "" && len(raw) > 0 {
		ct = http.DetectContentType(raw) // wie ein echter ResponseWriter: Typ aus dem Inhalt raten.
	}
	if !isTextual(ct, raw) {
		return errorResult(fmt.Sprintf("binary response (%s, %d bytes) is not supported; use the UI for this endpoint",
			ct, len(raw))), nil, nil
	}

	text := strings.TrimSpace(string(raw))
	if rec.code != http.StatusOK {
		// Ohne den Status sieht ein 302 auf die Login-Seite wie ein leerer Erfolg aus.
		prefix := fmt.Sprintf("HTTP %d", rec.code)
		if loc := rec.header.Get("Location"); loc != "" {
			prefix += " Location: " + loc
		}
		text = strings.TrimSpace(prefix + "\n" + text)
	}
	res := &sdk.CallToolResult{Content: []sdk.Content{&sdk.TextContent{Text: text}}}
	if rec.code >= http.StatusMultipleChoices {
		res.IsError = true
	}
	return res, nil, nil
}

// recorder puffert die Antwort bis zu limit Bytes. Darüber hinaus bekommt der Handler einen
// Schreibfehler — httptest.NewRecorder würde beliebig viel in den Speicher nehmen.
type recorder struct {
	header    http.Header
	code      int
	buf       bytes.Buffer
	limit     int
	truncated bool
}

var errResponseTooLarge = errors.New("mcp: response exceeds limit")

func (w *recorder) Header() http.Header { return w.header }

// WriteHeader hält den ersten finalen Status fest; 1xx sind informativ und überschreiben ihn nicht.
func (w *recorder) WriteHeader(code int) {
	if w.code == 0 && code >= 200 {
		w.code = code
	}
}

func (w *recorder) Write(b []byte) (int, error) {
	w.WriteHeader(http.StatusOK)
	room := w.limit - w.buf.Len()
	if room >= len(b) {
		return w.buf.Write(b)
	}
	if room > 0 {
		w.buf.Write(b[:room])
	}
	w.truncated = true
	// Der Writer-Vertrag verlangt die tatsächlich geschriebene Anzahl zum Fehler.
	return max(room, 0), errResponseTooLarge
}

// Flush committet die Antwort wie ein echter ResponseWriter, puffert aber weiter: Streaming
// unterstützt dieses Paket nicht. Ohne die Methode würde c.Response().Flush() panicken.
func (w *recorder) Flush() { w.WriteHeader(http.StatusOK) }

// isTextual sagt, ob der Body verlustfrei als Text zurückgehen kann. Xiri antwortet JSON; alles
// andere ist entweder ein Download oder ein Fehler des Hosts.
func isTextual(contentType string, body []byte) bool {
	if !utf8.Valid(body) {
		return false
	}
	if contentType == "" {
		return true
	}
	ct, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return strings.HasPrefix(ct, "text/") || ct == "application/json" || strings.HasSuffix(ct, "+json")
}

func errorResult(msg string) *sdk.CallToolResult {
	return &sdk.CallToolResult{IsError: true, Content: []sdk.Content{&sdk.TextContent{Text: msg}}}
}
