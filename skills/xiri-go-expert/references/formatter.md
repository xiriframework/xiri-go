# Formatter Package — Datum, Zahlen, Einheiten

Das `formatter`-Package liefert Locale-aware String-Formatierungen. In Komponenten werden Felder **meistens automatisch** formatiert (Table `FloatField`, `DateTimeField`, `DistanceField` etc. rufen intern die Formatter). Du brauchst sie v.a. für **Custom-Formatter**, **Card-Content-Texte** und **Page-Header-Subtitles**.

Import:

```go
import (
    "github.com/xiriframework/xiri-go/formatter"
    "github.com/xiriframework/xiri-go/types/distance"
    "github.com/xiriframework/xiri-go/types/locale"
    "github.com/xiriframework/xiri-go/types/pressure"
)
```

Alle Funktionen sind nil-safe bzgl. `*core.UiContext` (nutzen intern `ctx.SafeLocale()` / `ctx.SafeTimezone()`).

## Datum / Zeit (`formatter/datetime.go`)

### Timestamp ↔ `time.Time`

```go
formatter.ToUnixTimestamp(t time.Time) int64          // Seconds since epoch
formatter.ToUnixTimestampBigInt(t time.Time) int64    // Alias für int64-Emphasis
formatter.FromUnixTimestamp(ts int64) time.Time
formatter.FromUnixTimestampBigInt(ts int64) time.Time
```

### Unix-Timestamp → Locale-String

```go
formatter.FormatTimestampDateTime(ts int64, ctx *core.UiContext) string
  // "24.02.2024 18:30" (DE) / "02/24/2024 6:30 PM" (EnUS) / "24/02/2024 18:30" (EnGB)

formatter.FormatTimestampDate(ts int64, ctx *core.UiContext) string
  // "24.02.2024"

formatter.FormatTimestampFullDate(ts int64, ctx *core.UiContext) string
  // lang, mit Wochentag — Locale-abhängig

formatter.FormatTimestampToTextRange(ts int64, includeTime bool,
    timezone string, translate ...func(string) string) string
  // "heute um 14:00" / "vor 3 Stunden" / "gestern" — relative Text-Darstellung
```

### `time.Time` → Locale-String

```go
formatter.FormatDate(t time.Time, ctx *core.UiContext) string
formatter.FormatDateTime(t time.Time, ctx *core.UiContext) string
formatter.FormatTime(t time.Time, ctx *core.UiContext) string
```

### Spezial: Minutes-After-Midnight

```go
formatter.FormatMinutesAfterMidnight(
    dayTimestamp        int32,   // Unix-Timestamp (Sekunden) des Tages-Anfangs
    minutesAfterMidnight int16,  // 0..1440
    timezone            string,  // IANA-Name, z.B. "Europe/Vienna"
) string
  // "08:30" — für Fahrplan-artige Daten
```

## Zahlen (`formatter/numbers.go`)

```go
formatter.FormatInteger(value int64, ctx *core.UiContext) string
  // "1.234.567" (DE) / "1,234,567" (EnUS)

formatter.FormatDouble2(value float64, ctx *core.UiContext) string
  // 2 Nachkommastellen, Locale-Separatoren: "1.234,56" (DE) / "1,234.56" (EnUS)

formatter.FormatBigNumber(value float64, ctx *core.UiContext) string
  // "1,2M" / "3,5k" — kompakte Darstellung für KPI-Werte
```

## Locale-aware mit expliziter Locale (`formatter/locale.go`)

Für Fälle, wo du eine Locale hast, aber keinen `*UiContext` (z.B. in Tools oder Bulk-Exports):

```go
formatter.FormatNumberLocale(value float64, decimals int, loc locale.Locale) string
  // z.B. FormatNumberLocale(3.14159, 3, locale.De) → "3,142"

formatter.FormatDistanceLocaleWithDecimals(
    km        float64,
    distUnit  distance.Distance,  // Kilometer/Miles/Seemiles
    loc       locale.Locale,
    decimals  int,
) string
  // (42.5, Kilometer, De, 1) → "42,5 km"
  // (42.5, Miles,     De, 1) → "26,4 mi"

formatter.ConvertDistanceToKm(value float64, distUnit distance.Distance) float64
  // Rohwert-Konvertierung ohne Formatierung

formatter.FormatPressureLocale(bar float64, pressUnit pressure.Pressure, loc locale.Locale) string
  // (2.5, Bar, De)  → "2,5 bar"
  // (2.5, Psi, De)  → "36,3 psi"

formatter.FormatSpeedLocale(kmh float64, distUnit distance.Distance, loc locale.Locale) string
  // (50, Kilometer, De) → "50 km/h"
  // (50, Miles,     De) → "31 mph"
```

## Zeitdauer (`formatter/timeformat.go`)

Input ist **Sekunden** (int64).

```go
formatter.FormatTimeLengthHM(seconds int64, ctx *core.UiContext) string
  // 5430 → "01:30" (h:mm, bei < 1h: "30 min")

formatter.FormatTimeLengthHMS(seconds int64, ctx *core.UiContext) string
  // 5430 → "01:30:30"

formatter.FormatTimeLengthH(seconds int64, ctx *core.UiContext) string
  // 5430 → "1,5 h"

formatter.FormatTimeLengthMin(seconds int64, ctx *core.UiContext) string
  // 5430 → "90 min"
```

## TimeLimit (`formatter/timelimit.go`)

Spezial-Formatter für Zeitfenster-Definitionen (z.B. Öffnungszeiten, Geo-Fencing-Regeln).

```go
formatter.FormatTimeLimitFromDB(
    weekdaysStr string,      // Postgres-Array-Literal: "{t,t,t,t,f,f,f}"  (Mo-So)
    timeFrom    *string,     // "08:00:00" oder nil
    timeTo      *string,     // "17:00:00" oder nil
    timeIn      *bool,       // true = "im Zeitraum"/aktiv, false = "außerhalb"
    t           func(string) string,  // Translator für Wochentage-Namen
) string
```

Liefert einen lesbaren String, der die Regel zusammenfasst (z.B. "Mo-Fr 08:00-17:00").

## Wann welchen Formatter?

| Use-Case                                        | Funktion                                      |
| ----------------------------------------------- | --------------------------------------------- |
| Table-Spalte mit Datum (automatisch)            | `DateTimeField` / `DateField` — kein Call nötig |
| Custom-Formatter in Table-Spalte                | `FormatTimestampDateTime(ts, ctx)`             |
| PageHeader-Subtitle mit Datum                   | `FormatDate(someTime, uc)`                     |
| Stat-Value mit großer Zahl                      | `FormatBigNumber(v, uc)` → "1,2M"              |
| Card-Content-Line mit Distanz                   | `FormatDistanceLocaleWithDecimals(km, …)`      |
| Export-CSV mit Zahlen (keine `*UiContext` da)   | `FormatNumberLocale(v, 2, locale.De)`          |
| Zeitdauer (Dauer einer Tour, Downtime)          | `FormatTimeLengthHM(seconds, uc)`              |
| Unix-Timestamp von DB → Display                 | `FormatTimestampDate(ts, uc)`                  |
| `time.Time` von GORM → Display                  | `FormatDateTime(t, uc)`                        |

## Custom-Table-Formatter Beispiel

```go
b := table.NewBuilder[Device]()
b.FloatField("uptime", "Uptime", func(r Device) float64 { return r.Uptime }).
    WithFormatter(func(v any) string {
        seconds, ok := v.(float64)
        if !ok { return "" }
        return formatter.FormatTimeLengthHMS(int64(seconds), uc)
    })
```

Achtung: `WithFormatter` nimmt eine Func ohne `ctx`-Parameter — du schließt `uc` in der Closure (Build-Zeit) mit ein. Wenn der Table-Builder pro Request neu gebaut wird (empfohlen), ist das unproblematisch; bei wiederverwendeten Buildern würde das die falsche Locale einfrieren.

## Häufige Fehler

- **Seconds vs. Milliseconds**: xiri-go arbeitet mit **Sekunden**. Wenn dein Frontend / deine DB Millisekunden schickt, **vorher teilen** (`ts / 1000`).
- **Time vs. Timestamp**: `FormatDateTime` nimmt `time.Time`, `FormatTimestampDateTime` nimmt `int64`. Name beachten.
- **Pressure falsch**: `FormatPressureLocale(bar, unit, …)` — erster Parameter ist **immer in Bar** (Canonical-Unit). Konversion in Psi/Kpa passiert intern.
- **Distance in km**: `FormatDistanceLocaleWithDecimals(km, …)` — erster Parameter immer in Kilometern. Für Miles-in → `ConvertDistanceToKm(miles, distance.Miles)` davor.
- **Locale-without-Context**: Einige Funktionen brauchen `locale.Locale`, andere `*core.UiContext`. Wenn du nur den `uc` hast: `uc.SafeLocale()` holt den Wert raus.
