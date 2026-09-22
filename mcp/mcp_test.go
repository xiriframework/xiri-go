package mcp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// app ist eine Mini-Xiri-App: eine Seite mit Formular, ein Save-Endpoint mit Validierung, ein
// geschützter Endpoint, ein DELETE-Handler der mitzählt, und Endpoints für die Antwort-Grenzfälle.
type app struct {
	ts       *httptest.Server
	echo     *echo.Echo
	delCalls atomic.Int32
	blobSize int    // Body-Größe von /api/Test/Big
	seenAuth string // was der Handler im Authorization-Header gesehen hat
	seenHost string // was der Handler als Host-Header gesehen hat
	seenRaw  []string
}

func newApp(t *testing.T, opts Options) *app {
	t.Helper()
	a := &app{blobSize: 32}
	e := echo.New()
	e.GET("/api/Test/Page", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"bread":  []any{},
			"data":   []any{map[string]any{"type": "form", "data": map[string]any{"url": "Test/Save"}}},
			"cookie": c.Request().Header.Get("Cookie"),
		})
	})
	e.GET("/api/Test/Ping", func(c echo.Context) error {
		a.seenAuth = c.Request().Header.Get("Authorization")
		a.seenHost = c.Request().Header.Get("Host")
		a.seenRaw = c.Request().Header.Values("Cookie")
		return c.JSON(http.StatusOK, map[string]any{"pong": true, "ct": c.Request().Header.Get("Content-Type")})
	})
	e.POST("/api/Test/Save", func(c echo.Context) error {
		a.seenHost = c.Request().Header.Get("Host")
		var body map[string]any
		if err := c.Bind(&body); err != nil {
			return err
		}
		if body["name"] == nil || body["name"] == "" {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": "name fehlt", "fields": map[string]string{"name": "name fehlt"}})
		}
		return c.JSON(http.StatusOK, map[string]any{"done": true, "goto": "/Test/Page"})
	})
	// Registriert, damit die Methodenprüfung in act etwas zu verhindern hat: ohne sie käme DELETE
	// hier an und lieferte 200 statt eines Tool-Fehlers.
	e.DELETE("/api/Test/Save", func(c echo.Context) error {
		a.delCalls.Add(1)
		return c.JSON(http.StatusOK, map[string]any{"done": true})
	})
	e.GET("/api/Test/Big", func(c echo.Context) error {
		return c.Blob(http.StatusOK, echo.MIMEApplicationJSON, []byte(`"`+strings.Repeat("x", a.blobSize)+`"`))
	})
	// Gültiges UTF-8 mit binärem MIME-Typ: belegt den MIME-Zweig unabhängig von der utf8-Prüfung.
	e.GET("/api/Test/Xlsx", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", []byte("PK plain text"))
	})
	// Kein Content-Type, aber gültiges UTF-8 mit binärer Signatur: nur die Inhaltserkennung
	// (http.DetectContentType → application/pdf) kann das noch ablehnen.
	e.GET("/api/Test/Sniff", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "", []byte("%PDF-1.4 trailer"))
	})
	// application/problem+json ist Text und muss durchkommen.
	e.GET("/api/Test/Problem", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/problem+json", []byte(`{"title":"nope"}`))
	})
	// Handler ohne jede Ausgabe: der Status muss trotzdem 200 sein, nicht 0.
	e.GET("/api/Test/Empty", func(c echo.Context) error { return nil })
	// Route, die nur mit abschliessendem Slash registriert ist.
	e.GET("/api/Test/Slash/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{"slash": true, "q": c.QueryParam("a")})
	})
	// text/plain, aber kein gültiges UTF-8: als Text ginge der Inhalt verfälscht zurück.
	e.GET("/api/Test/Broken", func(c echo.Context) error {
		return c.Blob(http.StatusOK, echo.MIMETextPlain, []byte{0xff, 0xfe, 0x00})
	})
	e.GET("/api/Test/Login", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/Login")
	})
	// Geschützt: nur mit Authorization erreichbar — prüft, dass Middleware des Hosts greift.
	admin := e.Group("/api/Admin", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("Authorization") != "Bearer secret" {
				return c.JSON(http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
			}
			return next(c)
		}
	})
	admin.GET("/Secret", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{"secret": "42"})
	})
	e.Any("/mcp", echo.WrapHandler(Handler(e, opts)))
	a.echo = e
	a.ts = httptest.NewServer(e)
	t.Cleanup(a.ts.Close)
	return a
}

type headerRoundTripper struct {
	header http.Header
	next   http.RoundTripper
}

func (h headerRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	for k, v := range h.header {
		r.Header[k] = v
	}
	return h.next.RoundTrip(r)
}

// ctx begrenzt jeden Aufruf: ohne Frist würde ein hängender Aufruf erst am globalen Test-Timeout
// auffallen, und t.Context() bricht erst beim Cleanup ab.
func ctx(t *testing.T) context.Context {
	t.Helper()
	c, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	t.Cleanup(cancel)
	return c
}

func connect(t *testing.T, a *app, header http.Header) *sdk.ClientSession {
	t.Helper()
	client := sdk.NewClient(&sdk.Implementation{Name: "test", Version: "0"}, nil)
	hc := &http.Client{Transport: headerRoundTripper{header: header, next: http.DefaultTransport}}
	sess, err := client.Connect(ctx(t), &sdk.StreamableClientTransport{Endpoint: a.ts.URL + "/mcp", HTTPClient: hc}, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sess.Close() })
	return sess
}

// call ruft ein Tool auf und liefert Ergebnistext und IsError.
func call(t *testing.T, sess *sdk.ClientSession, name string, args map[string]any) (string, bool) {
	t.Helper()
	res, err := sess.CallTool(ctx(t), &sdk.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Content) == 0 {
		t.Fatal("empty content")
	}
	tc, ok := res.Content[0].(*sdk.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", res.Content[0])
	}
	return tc.Text, res.IsError
}

func TestListTools(t *testing.T) {
	sess := connect(t, newApp(t, Options{}), nil)
	res, err := sess.ListTools(ctx(t), nil)
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	for _, tool := range res.Tools {
		names[tool.Name] = true
	}
	if len(names) != 2 || !names["act"] || !names["read_page"] {
		t.Fatalf("expected act and read_page, got %v", names)
	}
}

func TestReadPage_ForwardsCookieAndReturnsPage(t *testing.T) {
	sess := connect(t, newApp(t, Options{}), http.Header{"Cookie": {"session=abc"}})
	body, isErr := call(t, sess, "read_page", map[string]any{"route": "Test/Page"})
	if isErr || !strings.Contains(body, `"url":"Test/Save"`) || !strings.Contains(body, `"cookie":"session=abc"`) {
		t.Fatalf("unexpected result (isError=%v): %s", isErr, body)
	}
}

func TestAct_PostsDataAndReturnsResponse(t *testing.T) {
	sess := connect(t, newApp(t, Options{}), nil)
	body, isErr := call(t, sess, "act", map[string]any{"url": "Test/Save", "data": map[string]any{"name": "Alpha"}})
	if isErr || !strings.Contains(body, `"goto":"/Test/Page"`) {
		t.Fatalf("unexpected result (isError=%v): %s", isErr, body)
	}
}

func TestAct_ValidationErrorIsToolError(t *testing.T) {
	sess := connect(t, newApp(t, Options{}), nil)
	body, isErr := call(t, sess, "act", map[string]any{"url": "Test/Save"})
	if !isErr || !strings.Contains(body, `"fields":{"name":"name fehlt"}`) || !strings.Contains(body, "HTTP 400") {
		t.Fatalf("expected isError with status and fields, got (isError=%v) %s", isErr, body)
	}
}

func TestAct_GetReachesHandler(t *testing.T) {
	sess := connect(t, newApp(t, Options{}), nil)
	body, isErr := call(t, sess, "act", map[string]any{"url": "Test/Ping", "method": "GET"})
	if isErr || !strings.Contains(body, `"pong":true`) {
		t.Fatalf("unexpected result (isError=%v): %s", isErr, body)
	}
}

func TestAct_GetWithDataIsRejected(t *testing.T) {
	sess := connect(t, newApp(t, Options{}), nil)
	body, isErr := call(t, sess, "act", map[string]any{"url": "Test/Ping", "method": "GET", "data": map[string]any{"x": 1}})
	if !isErr {
		t.Fatalf("GET with data must be rejected, got %s", body)
	}
}

// Gegenprobe zu Step 2: ohne die Methodenprüfung in act erreicht DELETE den registrierten Handler
// und liefert 200 — der Test scheitert dann am ersten Fatalf, delCalls belegt zusätzlich, dass die
// Ablehnung vor dem Dispatch greift und nicht bloß ein 405 zurückkam.
func TestAct_RejectsOtherMethodsWithoutReachingHandler(t *testing.T) {
	a := newApp(t, Options{})
	sess := connect(t, a, nil)
	body, isErr := call(t, sess, "act", map[string]any{"url": "Test/Save", "method": "DELETE"})
	if !isErr {
		t.Fatalf("DELETE must be rejected, got %s", body)
	}
	if n := a.delCalls.Load(); n != 0 {
		t.Fatalf("DELETE handler must not be reached, got %d calls", n)
	}
}

// Echte Xiri-Komponenten exportieren API-URLs samt Prefix ("/api/Portal/Devices/Save", via
// xurl.NewUrlPrefix). Ein Agent reicht genau diese url an act weiter — ohne diesen Test hängt
// dispatch den Prefix ein zweites Mal davor und jede reale App antwortet mit 404.
func TestAct_AcceptsUrlWithAndWithoutApiPrefix(t *testing.T) {
	sess := connect(t, newApp(t, Options{}), nil)
	for _, url := range []string{"Test/Save", "/api/Test/Save", "api/Test/Save"} {
		body, isErr := call(t, sess, "act", map[string]any{"url": url, "data": map[string]any{"name": "Alpha"}})
		if isErr || !strings.Contains(body, `"done":true`) {
			t.Fatalf("%q must reach the handler, got (isError=%v) %s", url, isErr, body)
		}
	}
	// read_page bekommt Routen normalerweise ohne Prefix (aus goto), muss den Prefix aber
	// genauso schlucken.
	body, isErr := call(t, sess, "read_page", map[string]any{"route": "/api/Test/Page"})
	if isErr || !strings.Contains(body, `"url":"Test/Save"`) {
		t.Fatalf("read_page with prefix failed (isError=%v): %s", isErr, body)
	}
}

func TestAct_RejectsTargetsOutsideApiPrefix(t *testing.T) {
	sess := connect(t, newApp(t, Options{}), nil)
	for _, url := range []string{"../mcp", "/../mcp", "Test/../../mcp", "%2e%2e/mcp"} {
		body, isErr := call(t, sess, "act", map[string]any{"url": url})
		if !isErr || !strings.Contains(body, "must stay below") {
			t.Fatalf("%q must be rejected, got (isError=%v) %s", url, isErr, body)
		}
	}
}

// Die Kernbehauptung des Designs: der interne Request läuft durch die Middleware des Hosts.
func TestHostMiddlewareApplies(t *testing.T) {
	a := newApp(t, Options{})

	body, isErr := call(t, connect(t, a, nil), "read_page", map[string]any{"route": "Admin/Secret"})
	if !isErr || !strings.Contains(body, "HTTP 401") {
		t.Fatalf("expected 401 without Authorization, got (isError=%v) %s", isErr, body)
	}

	sess := connect(t, a, http.Header{"Authorization": {"Bearer secret"}})
	body, isErr = call(t, sess, "read_page", map[string]any{"route": "Admin/Secret"})
	if isErr || !strings.Contains(body, `"secret":"42"`) {
		t.Fatalf("expected secret with Authorization, got (isError=%v) %s", isErr, body)
	}
}

// Mehrere Cookie-Zeilen bleiben erhalten (Header.Values statt Get), und Header ausserhalb von
// ForwardHeaders kommen nicht durch.
func TestForwardHeaders_KeepsMultipleValuesAndDropsOthers(t *testing.T) {
	a := newApp(t, Options{ForwardHeaders: []string{"Cookie"}})
	sess := connect(t, a, http.Header{"Cookie": {"a=1", "b=2"}, "Authorization": {"Bearer secret"}})
	if _, isErr := call(t, sess, "act", map[string]any{"url": "Test/Ping", "method": "GET"}); isErr {
		t.Fatal("ping failed")
	}
	if len(a.seenRaw) != 2 || a.seenRaw[0] != "a=1" || a.seenRaw[1] != "b=2" {
		t.Fatalf("expected both cookie values, got %v", a.seenRaw)
	}
	if a.seenAuth != "" {
		t.Fatalf("Authorization is not in ForwardHeaders and must not arrive, got %q", a.seenAuth)
	}
}

// Transport-Header gehören dem internen Request; ein MCP-Client darf sie nicht umbiegen, auch
// wenn sie in ForwardHeaders stehen. Direkt auf dispatch, weil der MCP-Transport selbst ein
// Content-Type erzwingt und ein abweichender Wert schon die Verbindung scheitern ließe.
func TestForwardHeaders_BlocksTransportHeaders(t *testing.T) {
	a := newApp(t, Options{})
	opts := withDefaults(Options{ForwardHeaders: []string{"Content-Type", "Host", "Authorization"}})
	req := &sdk.CallToolRequest{Extra: &sdk.RequestExtra{Header: http.Header{
		"Content-Type":  {"text/plain"},
		"Host":          {"evil.example"},
		"Authorization": {"Bearer secret"},
	}}}

	res, _, err := dispatch(ctx(t), a.echo, opts, req, http.MethodPost, "Test/Save",
		map[string]any{"name": "Alpha"})
	if err != nil {
		t.Fatal(err)
	}
	body := res.Content[0].(*sdk.TextContent).Text
	if res.IsError || !strings.Contains(body, `"done":true`) {
		t.Fatalf("Bind must still see JSON, got (isError=%v) %s", res.IsError, body)
	}
	if a.seenHost != "" {
		t.Fatalf("Host must not be forwarded, handler saw %q", a.seenHost)
	}
	// Gegenprobe zur Blockliste: ein nicht gesperrter Header derselben Liste kommt an.
	if _, _, err := dispatch(ctx(t), a.echo, opts, req, http.MethodGet, "Test/Ping", nil); err != nil {
		t.Fatal(err)
	}
	if a.seenAuth != "Bearer secret" {
		t.Fatalf("Authorization must be forwarded, handler saw %q", a.seenAuth)
	}
}

func TestResponseLimit(t *testing.T) {
	a := newApp(t, Options{MaxResponseBytes: 64})

	a.blobSize = 10 // 12 Bytes mit Quotes — passt
	if body, isErr := call(t, connect(t, a, nil), "act", map[string]any{"url": "Test/Big", "method": "GET"}); isErr {
		t.Fatalf("small response must pass: %s", body)
	}

	a.blobSize = 500
	body, isErr := call(t, connect(t, a, nil), "act", map[string]any{"url": "Test/Big", "method": "GET"})
	if !isErr || !strings.Contains(body, "exceeds 64 bytes") {
		t.Fatalf("large response must be rejected, got (isError=%v) %s", isErr, body)
	}
}

func TestBinaryResponseIsRejected(t *testing.T) {
	sess := connect(t, newApp(t, Options{}), nil)
	body, isErr := call(t, sess, "act", map[string]any{"url": "Test/Xlsx", "method": "GET"})
	if !isErr || !strings.Contains(body, "binary response") {
		t.Fatalf("binary response must be rejected, got (isError=%v) %s", isErr, body)
	}
}

func TestSniffedBinaryIsRejected(t *testing.T) {
	sess := connect(t, newApp(t, Options{}), nil)
	body, isErr := call(t, sess, "act", map[string]any{"url": "Test/Sniff", "method": "GET"})
	if !isErr || !strings.Contains(body, "binary response") {
		t.Fatalf("binary body without content-type must be rejected, got (isError=%v) %s", isErr, body)
	}
}

func TestProblemJSONPasses(t *testing.T) {
	sess := connect(t, newApp(t, Options{}), nil)
	body, isErr := call(t, sess, "act", map[string]any{"url": "Test/Problem", "method": "GET"})
	if isErr || !strings.Contains(body, `"title":"nope"`) {
		t.Fatalf("+json must pass, got (isError=%v) %s", isErr, body)
	}
}

func TestEmptyResponseIsOK(t *testing.T) {
	sess := connect(t, newApp(t, Options{}), nil)
	body, isErr := call(t, sess, "act", map[string]any{"url": "Test/Empty", "method": "GET"})
	if isErr || body != "" {
		t.Fatalf("empty response must be a plain success, got (isError=%v) %q", isErr, body)
	}
}

// path.Clean darf den Slash nicht fressen, und der Query-String muss ihn überleben.
func TestTrailingSlashRouteAndQuery(t *testing.T) {
	sess := connect(t, newApp(t, Options{}), nil)
	body, isErr := call(t, sess, "act", map[string]any{"url": "Test/Slash/?a=7", "method": "GET"})
	if isErr || !strings.Contains(body, `"slash":true`) || !strings.Contains(body, `"q":"7"`) {
		t.Fatalf("unexpected result (isError=%v): %s", isErr, body)
	}
}

// Der Recorder direkt: exakte Limitgrenze, verteilte Writes und die gemeldete Bytezahl.
func TestRecorderLimit(t *testing.T) {
	w := &recorder{header: http.Header{}, limit: 4}
	if n, err := w.Write([]byte("abcd")); n != 4 || err != nil || w.truncated {
		t.Fatalf("exact limit must pass: n=%d err=%v truncated=%v", n, err, w.truncated)
	}
	if n, err := w.Write([]byte("e")); n != 0 || err == nil || !w.truncated {
		t.Fatalf("one byte over must fail: n=%d err=%v truncated=%v", n, err, w.truncated)
	}

	w = &recorder{header: http.Header{}, limit: 4}
	n, err := w.Write([]byte("abcdef"))
	if n != 4 || err == nil || !w.truncated {
		t.Fatalf("partial write must report stored bytes: n=%d err=%v truncated=%v", n, err, w.truncated)
	}
	if got := w.buf.String(); got != "abcd" {
		t.Fatalf("buffer must hold exactly the limit, got %q", got)
	}
	if w.code != http.StatusOK {
		t.Fatalf("Write must commit 200, got %d", w.code)
	}

	// 1xx ist informativ und darf den finalen Status nicht belegen.
	w = &recorder{header: http.Header{}, limit: 4}
	w.WriteHeader(http.StatusEarlyHints)
	w.WriteHeader(http.StatusCreated)
	if w.code != http.StatusCreated {
		t.Fatalf("1xx must not win over the final status, got %d", w.code)
	}
}

func TestInvalidUTF8IsRejected(t *testing.T) {
	sess := connect(t, newApp(t, Options{}), nil)
	body, isErr := call(t, sess, "act", map[string]any{"url": "Test/Broken", "method": "GET"})
	if !isErr || !strings.Contains(body, "binary response") {
		t.Fatalf("invalid utf-8 must be rejected, got (isError=%v) %s", isErr, body)
	}
}

// Ein Redirect auf die Login-Seite hat einen leeren Body — ohne Status sähe er wie ein Erfolg aus.
func TestRedirectIsErrorWithStatusAndLocation(t *testing.T) {
	sess := connect(t, newApp(t, Options{}), nil)
	body, isErr := call(t, sess, "act", map[string]any{"url": "Test/Login", "method": "GET"})
	if !isErr || !strings.Contains(body, "HTTP 302") || !strings.Contains(body, "Location: /Login") {
		t.Fatalf("expected 302 as error with location, got (isError=%v) %s", isErr, body)
	}
}
