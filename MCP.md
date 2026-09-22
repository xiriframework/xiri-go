# MCP: die App für Agenten öffnen

Das Paket `mcp` stellt eine Xiri-App als [MCP](https://modelcontextprotocol.io)-Server bereit
(Streamable HTTP, offizielles [Go-SDK](https://github.com/modelcontextprotocol/go-sdk)). Ein Agent
liest damit Seiten und führt Aktionen aus — ohne dass Handler, Komponenten oder das Seiten-JSON
geändert werden.

Eine Zeile mountet ihn:

```go
import xirimcp "github.com/xiriframework/xiri-go/mcp"

e.Any("/mcp", echo.WrapHandler(xirimcp.Handler(e, xirimcp.Options{Name: "devices", Version: "1.0"})))
```

---

## Warum nur zwei Tools

Der übliche Weg wäre, jede Aktion einzeln als MCP-Tool zu registrieren: `create_device`,
`delete_device`, `list_devices`. Bei 200 Endpunkten sind das 200 Tool-Definitionen, die mit jeder
Änderung veralten.

Xiri braucht das nicht, weil das Seiten-JSON **sich selbst beschreibt**. Was
`GET /api/Portal/Devices` liefert, enthält bereits alles, was ein Agent wissen muss:

```json
{ "type": "form",  "data": { "url": "/api/Portal/Devices/Save",
                             "fields": [{ "id": "name", "type": "text", "required": true }] } }
{ "action": "api", "url": "/api/Portal/Devices/DeactivateAll", "text": "Alle deaktivieren" }
```

Ein Agent liest das, wie ein Mensch die Oberfläche liest. Also genügen zwei generische Tools:

| Tool | Entspricht | Macht |
| --- | --- | --- |
| `read_page(route)` | eine Seite ansehen | `GET <ApiPrefix><route>`, liefert das Seiten-JSON |
| `act(url, method, data)` | Button klicken, Formular absenden | `POST` (Default) oder `GET` an `url` |

Die Seite ist der **Katalog**, aus dem der Agent URLs abliest — sie ist aber **keine Schranke**.
Siehe [Umfang und Rechte](#umfang-und-rechte).

---

## Wie es funktioniert

Der ganze Ansatz steht auf einem Punkt: ein Tool-Aufruf baut einen `http.Request` im Speicher und
schickt ihn per `e.ServeHTTP` gegen **dasselbe `*echo.Echo`**, an dem auch die App hängt.

```
MCP-Client ──HTTP──► /mcp ──► Tool-Handler ──► dispatch()
                                                  │
                                        baut http.Request im Speicher
                                                  │
                                       e.ServeHTTP(recorder, request)
                                                  │
                              ┌───────────────────▼────────────────────┐
                              │  dasselbe Echo: Middleware → Routing    │
                              │  → derselbe Handler wie für den Browser │
                              └───────────────────┬────────────────────┘
                                                  │
                                       Antwort landet im recorder
                                                  ▼
                                      Body als Text zurück an den Client
```

Es geht **kein Netzwerk-Request** hinaus — `ServeHTTP` ist ein gewöhnlicher Methodenaufruf. Auch
eine URL wie `https://evil.example/x` landet deshalb nur als harmloser interner Pfad und erreicht
nichts Externes.

Der Gewinn: Auth-Middleware, Routing und Handler greifen unverändert, und es gibt keinen zweiten
Codepfad, der vom Browser-Pfad abweichen könnte. `read_page` liefert byte-identisch dasselbe wie
ein `curl` auf dieselbe Route.

### Was `dispatch` im Einzelnen tut

1. **URL bauen und prüfen.** Der `ApiPrefix` wird nur ergänzt, wenn er fehlt (Formulare und
   `api`-Buttons exportieren ihn bereits, siehe unten). Danach `path.Clean` plus die Prüfung, dass
   das Ziel unterhalb des Prefix bleibt — das fängt `../`-Ausbrüche. Geroutet wird anschließend
   exakt der geprüfte Pfad, damit Echo nichts anderes sieht als die Prüfung.
2. **Header übernehmen.** Nur was in `ForwardHeaders` steht. Transport- und Body-Header (`Host`,
   `Content-Type`, `Content-Length`, `Transfer-Encoding`, `Connection`) nie — sonst könnte ein
   Client den internen Request umbiegen. Werte werden kopiert, mehrere `Cookie`-Zeilen bleiben
   erhalten.
3. **Antwort einsammeln.** Ein eigener `recorder` statt `httptest.NewRecorder`: er puffert nur bis
   `MaxResponseBytes` und gibt dem Handler darüber einen Schreibfehler. Sonst könnte ein Endpunkt,
   der 500 MB schreibt, den Prozess aus dem Speicher drücken.
4. **Ergebnis verpacken.** Zurück geht Text. Binärantworten (z. B. Excel-Export) und ungültiges
   UTF-8 werden mit einer Meldung abgelehnt statt zerstört. Status ≠ 200 wird als
   `HTTP 302 Location: /Login` vorangestellt — ohne das sähe ein Login-Redirect wie ein leerer
   Erfolg aus. Ab Status 300 ist das Ergebnis `IsError`.

---

## Mounten und verbinden

```go
e := echo.New()
e.Use(meineAuthMiddleware)          // schützt /api UND /mcp
registriereApiRouten(e)

e.Any("/mcp", echo.WrapHandler(xirimcp.Handler(e, xirimcp.Options{
    Name: "devices", Version: "1.0",
})))
```

Der Mount-Pfad gehört **hinter dieselbe Auth-Middleware wie `/api`**. Das Paket authentifiziert
nichts, es reicht nur Header weiter.

Client-Seite, z. B. Claude Code:

```bash
claude mcp add --transport http devices https://host/mcp
```

Prüfen ohne Client:

```bash
curl -s -X POST localhost:8080/mcp \
  -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'
```

Die Antwort kommt als Server-Sent Event (`data: {…}`), nicht als nacktes JSON.

---

## Options

| Feld | Default | Bedeutung |
| --- | --- | --- |
| `Name` | `"xiri"` | Servername im `initialize`-Handshake |
| `Version` | `"0"` | dito |
| `Instructions` | Xiri-Protokollbeschreibung | Hinweise an den Client; beschreibt das Antwortformat und welche Button-Actions `act` überhaupt ausführt |
| `ApiPrefix` | `"/api/"` | wo die Xiri-Handler am Echo hängen; fehlende Slashes werden ergänzt, Ziele ausserhalb abgelehnt |
| `ForwardHeaders` | `["Cookie", "Authorization"]` | was aus dem MCP-Request in den internen Request kopiert wird |
| `MaxResponseBytes` | 1 MiB | darüber bekommt der Handler einen Schreibfehler und der Aufruf wird abgelehnt |

### URLs mit und ohne Prefix

`act` und `read_page` akzeptieren die URL **mit oder ohne** `ApiPrefix`. Das ist kein Komfort,
sondern notwendig: echte Komponenten exportieren API-URLs samt Prefix — `form.NewForm` und
`button.NewApiButton` bauen sie über `xurl.NewUrlPrefix` und geben `"/api/Portal/Devices/Save"`
aus. Routen aus `goto` und Breadcrumbs kommen dagegen ohne (`"/Portal/Devices"`). Ein Agent reicht
schlicht weiter, was im Seiten-JSON steht, und beide Formen müssen funktionieren.

---

## Der typische Ablauf

```
read_page("Portal/Devices")
  → table mit Zeilen
  → button {action:"api",  url:"/api/Portal/Devices/DeactivateAll"}
  → button {action:"link", url:"/Portal/Devices/New"}      ← Navigation, kein act

read_page("Portal/Devices/New")
  → form {url:"/api/Portal/Devices/Save", fields:[name, serial, active]}

act("/api/Portal/Devices/Save", data:{name:"Sensor West", serial:"SN-009"})
  → {"done":true, "goto":"/Portal/Devices"}   → read_page auf diese Route
  → oder HTTP 400 {"error":"…","fields":{"serial":"existiert bereits"}}
```

Welche Button-Action wie behandelt wird:

| `action` | Was der Agent tut |
| --- | --- |
| `api` | `act(url, data)` |
| `link`, `href`, `back` | Navigation → `read_page` auf die Ziel-Route, **nicht** `act` |
| `dialog` | die `url` liefert den Dialog-Inhalt: ohne `data`/`filter` per `read_page`, mit Filter- oder Initialdaten per `act` |
| `download` | nicht unterstützt (Binärantwort) |
| `close`, `return`, `debug`, `simulate` | reine Client-Aktionen, serverseitig ohne Bedeutung |

Antworttypen: `{"done":true}`, `{"goto":"/Route"}`, `{"refresh":"page|table|panel"}`, jeweils
optional mit `message` und `messageType`.

Der `fields`-Teil einer 400-Antwort ist der Grund, warum das praktisch funktioniert: der Agent
erfährt nicht bloß „ging nicht", sondern **welches Feld** warum — derselbe Mechanismus, mit dem
xiri-ng den Fehler am Eingabefeld anzeigt. Damit kann er gezielt korrigieren und `act` wiederholen.

---

## Umfang und Rechte

Der Agent handelt unter der Session, deren `Cookie` bzw. `Authorization` weitergereicht wird. Er
kann damit alles, was dieser Nutzer im Browser auch könnte — **einschließlich Endpunkten, die auf
der gerade gelesenen Seite nicht vorkamen.** `act` erreicht jeden Endpunkt unterhalb von
`ApiPrefix`. Das Seiten-JSON ist der Katalog, nicht die Schranke.

Das ist Absicht. Eine zweite Erlaubnisliste im Framework liefe neben der Endpunktautorisierung des
Hosts her und würde irgendwann von ihr abweichen — genau dort entstehen Lücken. Daraus folgt:

- Der Mount-Pfad gehört hinter dieselbe Auth-Middleware wie `/api`.
- Wer weniger freigeben will, mountet ein Echo mit weniger Routen.
- `ForwardHeaders` bewusst klein halten: jeder Eintrag erweitert, was ein Client in die App
  hineinreicht.
- Ist eine POST-Route CSRF-geschützt, muss der CSRF-Header in `ForwardHeaders` stehen, sonst
  scheitert `act` dort mit 403.

Die Gleichung „dieselben Rechte wie der Browser" gilt außerdem nur, solange der Host jeden Endpunkt
anhand dieser Credentials **im selben Echo** autorisiert. Prüfungen am Reverse-Proxy oder an
Request-Metadaten greifen beim internen Request nicht.

---

## Grenzen

- **Host-basiertes Routing und `HTTPSRedirect`-Middleware** am selben Echo funktionieren nicht: der
  interne Request hat weder `Host` noch TLS, und das SDK liefert den äußeren Request nicht aus
  (`RequestExtra` kennt nur `TokenInfo`, `Header`, `CloseSSEStream`).
- **Globale Echo-Middleware läuft zweimal**, für den MCP-Request und den internen Request. Ein
  globales Parallelitätslimit mit Kapazität 1 blockiert sich damit selbst.
- **Middleware am vorgelagerten Reverse-Proxy** läuft für den internen Request nicht mit.
- **Antwort-Header gehen verloren, auch `Set-Cookie`** — eine Session-Rotation des Hosts erreicht
  den Client nicht. Gemeldet werden nur Status und `Location`.
- **Kein eigener Timeout:** `e.ServeHTTP` ist synchron. Der MCP-Kontext wird durchgereicht,
  kooperative Handler brechen ab; ein unkooperativer hängt den Tool-Call.
- **Kein Streaming:** die Antwort wird gepuffert und erst am Ende zurückgegeben.
- **Keine Binärantworten:** Downloads bleiben Sache der UI.

---

## Testen

Gegen einen echten Client testen, nicht gegen `dispatch` — nur so läuft der Aufruf wirklich durch
Transport, Routing und Middleware:

```go
ts := httptest.NewServer(e)   // e enthält die App UND den /mcp-Mount
client := sdk.NewClient(&sdk.Implementation{Name: "test", Version: "0"}, nil)
sess, err := client.Connect(ctx, &sdk.StreamableClientTransport{
    Endpoint: ts.URL + "/mcp", HTTPClient: hc}, nil)
res, err := sess.CallTool(ctx, &sdk.CallToolParams{
    Name: "act", Arguments: map[string]any{"url": "Test/Save", "data": …}})
```

Header setzt man über einen `http.RoundTripper` am `HTTPClient` des Transports. Der wichtigste Test
ist der, dass die Middleware des Hosts greift: geschützte Route ohne Header → `HTTP 401` und
`IsError`, mit Header → Inhalt.

Drei Fallstricke:

- Der MCP-Transport erzwingt sein eigenes `Content-Type`; ein Test, der einen abweichenden Wert
  schickt, scheitert schon an der Verbindung (`Unsupported Media Type`). Die Header-Blockliste ist
  deshalb nur direkt über `dispatch` prüfbar.
- Jeder Aufruf braucht ein `context.WithTimeout`. `t.Context()` bricht erst beim Cleanup ab, ein
  hängender Aufruf fiele sonst erst am globalen Test-Timeout auf.
- **Handgeschriebene Test-JSONs bilden die Realität nicht ab.** Ein Testformular mit
  `"url": "Test/Save"` kommt in echtem Xiri-JSON nie vor — echte Komponenten schreiben
  `"/api/Test/Save"`. Wer nur gegen selbstgebaute Antworten testet, übersieht genau diese Klasse
  von Fehlern; ein Durchlauf gegen eine App aus echten Buildern findet sie.

---

## Troubleshooting

| Symptom | Ursache |
| --- | --- |
| `HTTP 404` bei `act` | Route stimmt nicht — die `url` unverändert aus dem Seiten-JSON übernehmen |
| `url must stay below /api/` | Ziel liegt ausserhalb von `ApiPrefix` (oder ein `../` darin) |
| `HTTP 401`/`403` | `Cookie`/`Authorization` fehlt im MCP-Request oder nicht in `ForwardHeaders`; bei 403 an CSRF denken |
| `HTTP 302 Location: /Login` | Auth-Middleware leitet um — der Agent ist nicht angemeldet |
| `binary response (…) is not supported` | Download-Endpunkt; gehört in die UI, nicht in `act` |
| `response exceeds … bytes` | `MaxResponseBytes` erhöhen oder den Endpunkt paginieren |
| `data is only supported with POST` | bei `GET` gehören Parameter in die `url` |

---

## Siehe auch

- Kurzeinstieg im [README](README.md#mcp-die-app-für-agenten-öffnen)
- Referenz für den `xiri-go-expert`-Skill: `skills/xiri-go-expert/references/mcp.md`
- Doc-Kommentar des Pakets: `go doc github.com/xiriframework/xiri-go/mcp`
