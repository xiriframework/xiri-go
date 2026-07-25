# UiContext — Translator, Locale, Timezone, Einheiten

Der `*core.UiContext` ist der Per-Request-Kontext, der jede `Print(ctx)`-Methode erreicht. Er bündelt Sprache, Locale, Timezone, Einheiten und den Translator. Alle Formatierungen (Datum, Zahlen, Distanzen) und alle übersetzbaren Strings hängen daran.

## Struct

Quelle: `component/core/context.go`.

```go
type UiContext struct {
    Timezone  timezone.Timezone
    Lang      language.Language
    Locale    locale.Locale
    Distance  distance.Distance
    Pressure  pressure.Pressure
    Translate TranslateFunc
}

type TranslateFunc func(key string) string
```

Die Typen sind `int`-Aliase (Enums), keine Strings. Werte kommen aus:

- `types/language`: `language.Deutsch`, `language.Englisch`, `language.Kroatisch`, `language.Spanisch`, `language.Franzoesisch`, `language.Italienisch`, … (27 Werte)
- `types/locale`: `locale.De`, `locale.EnGB`, `locale.DeAT`, `locale.DeCH`, `locale.EnUS`, `locale.Fr`, `locale.It`, … (31 Werte — feiner als Language)
- `types/timezone`: `timezone.EuropeVienna`, `timezone.EuropeBerlin`, `timezone.EuropeLondon`, … (viele)
- `types/distance`: `distance.Kilometer`, `distance.Miles`, `distance.Seemiles`
- `types/pressure`: `pressure.Bar`, `pressure.Psi`, `pressure.Kpa`

## Nil-safe Accessors

Ein `*UiContext` darf `nil` sein — alle Komponenten nutzen Nil-safe Methoden:

```go
func (uc *UiContext) SafeTranslate(key string) string
  // gibt key zurück wenn uc==nil oder uc.Translate==nil

func (uc *UiContext) SafeLocale() locale.Locale
  // gibt locale.De (Default) zurück wenn uc==nil

func (uc *UiContext) SafeTimezone() timezone.Timezone
  // gibt timezone.EuropeVienna (Default) zurück wenn uc==nil

func (uc *UiContext) SafeDistance() distance.Distance
  // gibt distance.Kilometer (Default) zurück wenn uc==nil

func (uc *UiContext) SafePressure() pressure.Pressure
  // gibt pressure.Bar (Default) zurück wenn uc==nil
```

Es gibt außerdem eine freie Funktion:

```go
func core.Translate(ctx *UiContext, key string) string
// Äquivalent zu ctx.SafeTranslate(key) — für Fälle wo ctx sehr tief via pointer weitergereicht wird
```

## Wo wird `UiContext` gebaut?

xiri-go selbst liefert **keinen** Bootstrap-Constructor — das ist Projekt-spezifisch. Typische Varianten:

### Variante A: Middleware setzt UiContext auf Echo-Context

```go
// Projekt-Middleware
func UiContextMiddleware(translator i18n.Translator) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            user := getUserFromSession(c)

            uc := &core.UiContext{
                Timezone:  user.Timezone,    // timezone.Timezone
                Lang:      user.Language,
                Locale:    user.Locale,
                Distance:  user.DistanceUnit,
                Pressure:  user.PressureUnit,
                Translate: func(key string) string {
                    return translator.T(user.Language, key)
                },
            }
            c.Set("uicontext", uc)
            return next(c)
        }
    }
}

// In Handlern
func (ctrl *Ctrl) List(c echo.Context) error {
    uc := c.Get("uicontext").(*core.UiContext)
    // ... uc.SafeTranslate(...), uc.SafeLocale() ...
}
```

### Variante B: GEM WebContext Pattern (bevorzugt)

Das GEM-Projekt wrapped Echo-Context in einen `webcontext.WebContext`, der `UiContext()` liefert:

```go
wc := webcontext.GetWebContext(ctx)
uc := wc.UiContext()     // *core.UiContext
```

Hier ist die Middleware bereits gesetzt — im Controller braucht man nur `wc.UiContext()`.

### Anonymer UiContext (Tests, Public-Endpoints)

```go
uc := &core.UiContext{
    Locale:   locale.De,
    Timezone: timezone.EuropeVienna,
    Distance: distance.Kilometer,
    Pressure: pressure.Bar,
    // Translate = nil → SafeTranslate gibt keys 1:1 zurück
}
```

Das ist ein gültiger "Default"-Context für Rendering ohne User-Bindung.

## `TranslateFunc` — der Übersetzer

```go
type TranslateFunc func(key string) string
```

Der Function-Type ist bewusst minimal: **ein String rein, ein String raus**. Das i18n-Backend (json-Files, Datenbank, gotext, was auch immer) ist komplett in der Projekt-Implementierung. xiri-go ruft `ctx.SafeTranslate(key)` an den richtigen Stellen auf — wenn `Translate==nil`, kommt der Key selbst zurück, d.h. der Fallback-Text ist der Key.

### Konvention für Translation-Keys

xiri-go erwartet **keine** bestimmte Key-Struktur. Gängig sind:

- **UPPERCASE_SNAKE**: `DEVICES_TITLE`, `SAVE`, `DELETE_QUESTION`
- **Dotted-lowercase**: `devices.title`, `buttons.save`, `dialog.delete`

Bei Mixing von Backend-Translation und UI-Labels wird meist ein Prefix gewählt (`form.`, `table.`, `button.`, …). Für Labels, die Go-Code definiert (Field-Namen, Table-Headers), nutzt man Translation-Keys in den Konstruktoren:

```go
field.NewTextField("name", "form.device.name", true, "")   // "form.device.name" ist der Translation-Key
```

## Accessor-Verhalten im Detail

### `SafeLocale()`

Default = `locale.De` (Deutschland). Nutze das, wenn du einen Format-Helper direkt aufrufst:

```go
import "github.com/xiriframework/xiri-go/formatter"

s := formatter.FormatDateTime(time.Now(), uc)   // nil-safe
```

### `SafeTimezone()`

Default = `timezone.EuropeVienna`. Wird von `FormatTimestampDateTime` und verwandten Helpern gelesen, um Unix-Timestamps in lokale Zeit umzurechnen.

### `SafeDistance()` / `SafePressure()`

Defaults = `distance.Kilometer` bzw. `pressure.Bar`. Relevant für `FormatDistanceLocale*` /
`FormatPressureLocale` / `FormatSpeedLocale` und die Tabellen-Formatter, die diese Einheiten
umrechnen. Nutze immer die Safe-Variante — `ctx.Distance` direkt zu dereferenzieren paniced bei
`nil`-Context, und `core.Component.Print` erlaubt `ctx == nil` ausdrücklich.

## Enum-Werte aufzählen

Die Lookup-Tabellen der `types/*`-Packages sind **nicht** öffentlich. Für einzelne Werte gibt es
Accessoren, zum Aufzählen (z. B. für ein Dropdown) `All()`:

```go
for _, tz := range timezone.All() {          // sortiert nach numerischem Wert
    fmt.Println(tz.ToInt32(), timezone.GetName(tz), tz.GetIANA())
}
```

Gibt es in `timezone`, `locale`, `language`, `distance` und `pressure`. Die zurückgegebene Slice ist
pro Aufruf frisch alloziert und darf gefahrlos sortiert oder gefiltert werden.

**Achtung:** `GetName` ist eine **Package-Funktion** (`timezone.GetName(tz)`), die übrigen
Accessoren sind Methoden auf dem Wert.

| Package    | Name                  | Zweiter Accessor (Methode) |
| ---------- | --------------------- | -------------------------- |
| `timezone` | `timezone.GetName(v)` | `v.GetIANA()`              |
| `locale`   | `locale.GetName(v)`   | `v.GetLocaleString()`      |
| `language` | `language.GetName(v)` | `v.GetCode()`              |
| `distance` | `distance.GetName(v)` | `v.GetSymbol()`            |
| `pressure` | `pressure.GetName(v)` | `v.GetSymbol()`            |

Alternativ liefert `v.String()` denselben Namen wie `GetName` (bzw. `"Unknown"`).

## Minimale Fake-Implementation für Tests

```go
func testContext(t *testing.T) *core.UiContext {
    return &core.UiContext{
        Locale:   locale.De,
        Timezone: timezone.EuropeVienna,
        Distance: distance.Kilometer,
        Pressure: pressure.Bar,
        Translate: func(key string) string {
            return "[" + key + "]"   // sichtbar im Test, welche Keys übersetzt wurden
        },
    }
}
```

## Häufige Fehler

- **`uc.Translate(key)` direkt aufrufen** — crasht wenn `uc` oder `uc.Translate` nil. Immer `uc.SafeTranslate(key)` oder `core.Translate(uc, key)`.
- **Language statt Locale in Formattern** — `Language` ist die UI-Sprache, `Locale` die Format-Konvention (Datum, Zahlen). Für `De` gibt es `locale.De`, `locale.DeAT`, `locale.DeCH` — alle sind "Deutsch", unterscheiden sich aber in Zahlenformat.
- **Timezone falsch setzen** — Unix-Timestamps sind UTC. `FormatTimestampDateTime(ts, uc)` konvertiert in `uc.Timezone`. Wenn dein User in NY sitzt und du setzt `EuropeVienna`, kriegt er die MEZ-Zeit angezeigt.
- **Translate-Keys leak ins UI** — passiert wenn die Middleware fehlt oder das Translation-Backend einen Key nicht kennt. Logs um `SafeTranslate` in der Translation-Func helfen beim Diagnosieren.
