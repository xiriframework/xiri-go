# Response-Typen (response package)

Import: `"github.com/xiriframework/xiri-go/response"`

## Success Responses

Alle implementieren `SuccessResponse` Interface und haben `Done: true`.

### ReturnRefreshPage

Seite neu laden nach erfolgreicher Operation.

```go
resp := response.NewReturnRefreshPage()
// → {"done": true, "refresh": "page"}

resp := response.NewReturnRefreshPage().WithMessage("Gespeichert", response.MessageSuccess)
// → {"done": true, "refresh": "page", "message": "Gespeichert", "messageType": "success"}
```

### ReturnRefreshTable

Tabelle neu laden (z.B. nach Zeile löschen).

```go
resp := response.NewReturnRefreshTable()
// → {"done": true, "refresh": "table"}

resp := response.NewReturnRefreshTable().WithMessage("Gelöscht", response.MessageSuccess)
```

### ReturnGoto

Weiterleitung zu anderer URL.

```go
resp := response.NewReturnGoto("/portal/user/page/7")
// → {"done": true, "goto": "/portal/user/page/7"}

resp := response.NewReturnGoto("/portal/devices").WithMessage("Erstellt", response.MessageSuccess)
```

### ReturnDone

Einfache Bestätigung ohne Aktion.

```go
resp := response.NewReturnDone()
// → {"done": true}

resp := response.NewReturnDone().WithMessage("OK", response.MessageInfo)
```

### ReturnMessage / NewReturnSuccess / NewReturnError

Nur Snackbar-Nachricht, keine Navigation.

```go
resp := response.NewReturnSuccess("Einstellungen gespeichert")
// → {"done": true, "message": "Einstellungen gespeichert", "messageType": "success"}

resp := response.NewReturnError("Fehler beim Speichern")
// → {"done": true, "message": "Fehler beim Speichern", "messageType": "error"}

resp := response.NewReturnMessage("Hinweis", response.MessageWarning)
```

### ReturnPoll — selbst-pollender Button (Background-Worker)

Sagt dem auslösenden **Button** (oder anderem Initiator), eine Status-URL im Intervall
weiter abzufragen, bis eine Antwort **ohne** `poll` kommt. Der Button zeigt währenddessen
Spinner + Countdown (oder freien Text, s.u.) und ist disabled. Verwendung: Ein Button (per
`action:'api'` oder Dialog) stößt einen Worker an; solange er läuft, liefert der
Status-Handler `ReturnPoll`, am Ende eine normale Response (Snackbar/Refresh/Goto/Done).

```go
// Start-Handler: Worker anstoßen + Polling aktivieren
return wc.Component(response.NewReturnPoll(statusUrl.PrintPrefix(), 2000).
    WithMessage("Worker gestartet", response.MessageInfo))
// → {"done":true,"poll":2000,"pollUrl":"/api/.../Status","message":"Worker gestartet","messageType":"info"}

// Status-Handler (wird vom Button per GET gepollt):
if worker.Running() {
    return wc.Component(response.NewReturnPoll(statusUrl.PrintPrefix(), 2000).
        WithText(fmt.Sprintf("läuft… %d %%", worker.Percent())))   // freier Button-Text pro Tick
}
// fertig → kein poll → Button stoppt; finale Response normal verarbeitet:
return wc.Component(response.NewReturnRefreshTable().WithMessage("Fertig", response.MessageSuccess))
```

- `poll` (ms): Intervall; fehlt es / 0 → Polling stoppt.
- `pollUrl`: Status-Endpoint (GET). Optional — fehlt er, nutzt der Button seine eigene `url`.
  Bei **Dialog**-Buttons ist `pollUrl` Pflicht (button.url ist dort die Dialog-URL).
- `WithText(...)`: optionaler Text **im Button** während des Pollings (überschreibt den
  Countdown); kann pro Tick aktualisiert werden → echter Fortschritt.
- Die finale Response (ohne poll) läuft durch den ResponseHandler → kann `RefreshPage`,
  `RefreshTable`, `Goto`, `Message` sein.

### ButtonPatch — Button am Ende/laufend ändern (`WithButton`)

Jede Antwort, die ein Button verarbeitet (Poll-Tick **oder** finale Antwort), darf den
Button selbst patchen: Text, Farbe, Icon, Type, Hint, disabled. Verfügbar via `WithButton`
auf `ReturnPoll`, `ReturnMessage`, `ReturnDone`. Der Patch bleibt bestehen, bis die Aktion
erneut ausgelöst wird.

```go
response.NewButtonPatch().
    WithText("Erledigt ✓").WithColor("success").WithIcon("check").Disable()
// → {"text":"Erledigt ✓","color":"success","icon":"check","disabled":true}

// am Ende: Snackbar + Button auf "Erledigt"/grün/disabled
return wc.Component(response.NewReturnSuccess("Worker fertig").
    WithButton(response.NewButtonPatch().WithText("Erledigt ✓").WithColor("success").Disable()))

// während des Pollings: Button-Farbe wechseln
return wc.Component(response.NewReturnPoll(statusUrl.PrintPrefix(), 2000).
    WithButton(response.NewButtonPatch().WithColor("warn")))
```

Setter: `WithText`, `WithColor`, `WithIcon`, `WithType`, `WithHint`, `Disable()`, `Enable()`
(alle fluent, nur gesetzte Felder landen im JSON).

## Error Response

Für HTTP-Fehler (400, 404, 500, etc.).

```go
resp := response.NewErrorResponse("Ungültige Eingabe")
// → {"error": "Ungültige Eingabe"}

// Typischer Einsatz:
return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
return c.JSON(http.StatusNotFound, response.NewErrorResponse("Nicht gefunden"))
```

## Data Response

Standard-Envelope für Komponenten-Daten (Card, Stat, StatGrid, List, etc.).

```go
resp := response.NewDataResponse(myData)
// → {"data": myData}
```

## DataResult

Für Responses mit verschiedenen Formaten (JSON, CSV, Excel). Wird von Table-Komponenten verwendet — aber auch von allen Komponenten, die eine `DataResponse(ctx)`-Methode haben.

```go
result := response.NewJSONDataResult(data)   // {"data": data}
result := response.NewCSVDataResult(csvStr)  // CSV string
result := response.NewExcelDataResult(bytes) // Excel bytes
```

### `component.DataResponse(ctx)` — Partial-Refresh-Pattern

Viele Komponenten bieten eine `DataResponse(ctx) response.DataResult`-Methode, die **nur die data-Portion** ohne den `type`-Wrapper liefert. Gedacht für AJAX-Endpoints, die einer bereits gerenderten Komponente im Frontend neue Daten einspielen (z.B. `WithReload`-gepolltes Card/Stat).

Unterstützt:

| Komponente        | Methode                                  | Typischer Use-Case                       |
| ----------------- | ---------------------------------------- | ---------------------------------------- |
| `Card`            | `card.DataResponse(ctx)`                 | AJAX-Card lädt neue Inhalte              |
| `Stat`            | `stat.DataResponse(ctx)`                 | Polling-KPI                              |
| `StatGrid`        | `statgrid.DataResponse(ctx)`             | Polling-Dashboard                        |
| `MultiStat`       | `multistat.DataResponse(ctx)`            | Polling-KPI-Karte (mehrere Zahlen)       |
| `List`            | `list.DataResponse(ctx)`                 | Dynamische Listen-Updates                |
| `MultiProgress`   | `mp.DataResponse(ctx)`                   | Progress-Bar-Update                      |
| `EmptyState`      | `es.DataResponse(ctx)`                   | State-Wechsel                            |
| `ButtonLine`      | `bl.DataResponse(ctx)`                   | Button-Refresh nach State-Änderung       |

Beispiel — Card mit `SetURL` + AJAX-Reload:

```go
// Page: Card als Shell, lädt Daten lazy
card := card.NewCard(core.CardTypeTable, nil, "Live-Status", "", "", "", true, false, "")
card.SetURL(c.apiUrl("card", "status"))
card.WithReload(true)   // Frontend pollt

// AJAX-Endpoint liefert nur die Daten (ohne type-Wrapper)
func (c *Controller) CardStatus(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    content := card.NewCardListContent([]card.CardListContentLine{
        {Name: "Online",  Content: strconv.Itoa(c.svc.CountActive())},
        {Name: "Offline", Content: strconv.Itoa(c.svc.CountOffline())},
    })
    inner := card.NewCardList("Live-Status", content)
    return c.JSON(http.StatusOK, inner.DataResponse(wc.UiContext()))
}
```

**Key-Unterschied zu `Print()`**: `Print()` liefert das vollständige Komponenten-JSON inkl. `type`, `display`, `data` etc. — geeignet für erstes Rendern. `DataResponse()` liefert nur die `{data: ...}`-Portion — geeignet um eine bereits existierende Komponente im Frontend mit neuen Inhalten zu füttern, **ohne** sie neu zu mounten.

### `component.PrintData(ctx)` — das rohe data-Map

Parallel zu `DataResponse` gibt es auf vielen Komponenten noch `PrintData(ctx) map[string]any` — das ist die **unverpackte** Form (nicht in `{"data": …}` geschachtelt). Sinnvoll, wenn du mehrere Daten-Stücke in einer eigenen JSON-Antwort kombinierst:

```go
return c.JSON(http.StatusOK, map[string]any{
    "main":     mainCard.PrintData(uc),
    "sidebar":  sideCard.PrintData(uc),
    "buttons":  bl.PrintData(uc),
})
```

## MessageType Werte

```go
response.MessageSuccess  // "success" — grüne Snackbar
response.MessageError    // "error"   — rote Snackbar
response.MessageInfo     // "info"    — blaue Snackbar
response.MessageWarning  // "warning" — gelbe Snackbar
```

## Typischer Handler-Pattern

```go
func HandleSave(c echo.Context) error {
    // ... Validierung, DB-Operation ...

    if err != nil {
        return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
    }

    return c.JSON(http.StatusOK,
        response.NewReturnRefreshPage().WithMessage(
            ctx.SafeTranslate("saved"),
            response.MessageSuccess,
        ))
}
```
