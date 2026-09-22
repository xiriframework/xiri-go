# MCP: die App für Agenten öffnen

Das Paket `mcp` stellt eine Xiri-App als MCP-Server bereit (Streamable HTTP, offizielles
`github.com/modelcontextprotocol/go-sdk`). Es gibt **zwei generische Tools**, kein Action-Registry:

| Tool | Macht | Ergebnis |
| --- | --- | --- |
| `read_page(route)` | `GET <ApiPrefix><route>` | das Seiten-JSON, wie es xiri-ng auch bekommt |
| `act(url, method, data)` | `POST` (Default) oder `GET` an `url` | die Xiri-Response |

Das Seiten-JSON beschreibt bereits, was ein Agent tun kann: Formulare haben `data.url` und
`data.fields`, Buttons haben `action`, `url` und optional `data`. Die Seite ist damit der **Katalog**
— aber **keine Schranke**: `act` erreicht jeden Endpunkt unterhalb von `ApiPrefix`.

## Mount

```go
import xirimcp "github.com/xiriframework/xiri-go/mcp"

e.Any("/mcp", echo.WrapHandler(xirimcp.Handler(e, xirimcp.Options{Name: "devices", Version: "1.0"})))
```

Der Mount-Pfad gehört **hinter dieselbe Auth-Middleware wie `/api`** — das Paket forwardet nur
Header, es authentifiziert nichts.

## Options

| Feld | Default | Bedeutung |
| --- | --- | --- |
| `Name` | `"xiri"` | Servername im `initialize`-Handshake |
| `Version` | `"0"` | dito |
| `Instructions` | Xiri-Protokollbeschreibung | Hinweise an den Client |
| `ApiPrefix` | `"/api/"` | wo die Xiri-Handler am Echo hängen; fehlende Slashes werden ergänzt, Ziele ausserhalb abgelehnt |
| `ForwardHeaders` | `["Cookie", "Authorization"]` | was aus dem MCP-Request in den internen Request kopiert wird |
| `MaxResponseBytes` | 1 MiB | ab hier bekommt der Handler einen Schreibfehler und der Aufruf wird abgelehnt |

## Sicherheit

- **Der Host schützt den Mount-Pfad.** Gleiche Middleware wie `/api`, sonst ist die App offen.
- **`ForwardHeaders` klein halten.** Jeder Eintrag erweitert, was ein MCP-Client in die App
  hineinreicht. Transport- und Body-Header (`Host`, `Content-Type`, `Content-Length`,
  `Transfer-Encoding`, `Connection`) werden nie kopiert, auch wenn sie in der Liste stehen.
- **CSRF:** ist eine POST-Route CSRF-geschützt, muss der CSRF-Header in `ForwardHeaders` stehen —
  sonst scheitert `act` dort mit 403.
- **Rechteumfang:** der Agent handelt unter der weitergereichten Session. Das entspricht den Rechten
  des Browsers, **solange der Host jeden Endpunkt anhand dieser Credentials im selben Echo
  autorisiert**. Prüfungen am Reverse-Proxy oder an Request-Metadaten greifen nicht (siehe Grenzen).
  Wer weniger freigeben will, mountet einen eigenen Echo mit weniger Routen.

## Grenzen

- **Host-basiertes Routing und `HTTPSRedirect`-Middleware** am selben Echo funktionieren nicht: der
  interne Request hat weder `Host` noch TLS, und das SDK liefert den äußeren Request nicht aus.
- **Globale Echo-Middleware läuft zweimal** (MCP-Request + interner Request). Ein globales
  Parallelitätslimit mit Kapazität 1 blockiert sich selbst.
- **Reverse-Proxy-Middleware** läuft für den internen Request nicht mit.
- **Antwort-Header gehen verloren, auch `Set-Cookie`** — eine Session-Rotation erreicht den Client
  nicht. Gemeldet werden nur Status und `Location`.
- **Kein eigener Timeout**, kein Streaming, keine Binärantworten (Downloads bleiben Sache der UI).

## Typischer Agenten-Ablauf

1. `read_page("Portal/Devices")` → Komponenten durchgehen.
2. Formular (`type: "form"` → `data.url`, `data.fields`) oder Button mit `action: "api"` finden.
   `link`, `href`, `back` sind **Navigation** → dafür `read_page` auf die Ziel-Route, nicht `act`.
   `download` wird nicht unterstützt; `close`, `return`, `debug`, `simulate` sind Client-Aktionen.
   Bei `dialog` liefert die `url` den Dialog-Inhalt: ohne `data`/`filter` per `read_page`, mit
   Filter- oder Initialdaten per `act`.
3. `act(url, data: {<fieldId>: <wert>, …})`.
4. Antwort auswerten: `{"done":true}`, `{"goto":"/Route"}` → `read_page` auf die Route **ohne**
   führenden Slash, `{"refresh":"page|table|panel"}`, optional `message`.
5. Bei 400 liefert die Antwort `error` und `fields` (`{"<fieldId>": "<meldung>"}`) — damit kann der
   Agent gezielt das falsche Feld korrigieren und `act` wiederholen. Status ≠ 200 steht als
   `HTTP <code>` in der ersten Zeile, ab 300 ist das Tool-Ergebnis `IsError`.

## Testmuster

Gegen einen echten Client testen, nicht gegen `dispatch`:

```go
ts := httptest.NewServer(e) // e enthält die App UND e.Any("/mcp", echo.WrapHandler(Handler(e, opts)))
client := sdk.NewClient(&sdk.Implementation{Name: "test", Version: "0"}, nil)
sess, err := client.Connect(ctx, &sdk.StreamableClientTransport{Endpoint: ts.URL + "/mcp", HTTPClient: hc}, nil)
res, err := sess.CallTool(ctx, &sdk.CallToolParams{Name: "act", Arguments: map[string]any{"url": "Test/Save", "data": …}})
```

Header setzt man über einen `http.RoundTripper` am `HTTPClient` des Transports. Der wichtigste Test
ist der, dass die **Middleware des Hosts greift**: geschützte Route ohne Header → `HTTP 401` und
`IsError`, mit Header → Inhalt.

Zwei Fallstricke: der MCP-Transport erzwingt sein eigenes `Content-Type`, ein Test mit abweichendem
Wert scheitert schon an der Verbindung (`Unsupported Media Type`) — die Header-Blockliste ist nur
direkt über `dispatch` prüfbar. Und jeder Aufruf braucht ein `context.WithTimeout`: `t.Context()`
bricht erst beim Cleanup ab, ein hängender Aufruf fiele sonst erst am globalen Test-Timeout auf.
