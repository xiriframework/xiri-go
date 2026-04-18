# TachoTime — Fahrtenschreiber-Diagramm

`component/tachotime` rendert Tacho-/Lenkzeit-Daten (Compliance-Dashboards, LKW-Telemetrie). Selten gebraucht, hier als Referenz.

## Konstruktoren

Quelle: `component/tachotime/tachotime.go`.

```go
func NewTachoTime(header string, data []TachoTimeDay, display *string) *TachoTime

func NewTachoTimeDay(
    date string,
    minDate, maxDate int64,
    data []TachoTimeData,
    drive []TachoTimeDriveBlock,
    driveDay []TachoTimeDriveDay,
    driving, working, breakTime, available, total int32,
) TachoTimeDay

func NewTachoTimeData(start, end int64, activity int32, duration string) TachoTimeData
func NewTachoTimeDriveBlock(start, end int64, length int32, data TachoTimeDriveBlockData) TachoTimeDriveBlock
func NewTachoTimeDriveBlockData(driving int32, duration, start, end string) TachoTimeDriveBlockData
func NewTachoTimeDriveDay(start, end int64, unknown int32, data TachoTimeDriveDayData) TachoTimeDriveDay
func NewTachoTimeDriveDayData(duration, start, end string) TachoTimeDriveDayData
```

## Beispiel

```go
day := tachotime.NewTachoTimeDay(
    "2024-04-17",
    1713312000, 1713398400,
    []tachotime.TachoTimeData{
        tachotime.NewTachoTimeData(1713312000, 1713315600, 1, "01:00"), // 1 = driving
        tachotime.NewTachoTimeData(1713315600, 1713317400, 2, "00:30"), // 2 = break
    },
    []tachotime.TachoTimeDriveBlock{
        tachotime.NewTachoTimeDriveBlock(
            1713312000, 1713322800, 25000,
            tachotime.NewTachoTimeDriveBlockData(240, "04:00", "08:00", "12:00"),
        ),
    },
    []tachotime.TachoTimeDriveDay{
        tachotime.NewTachoTimeDriveDay(
            1713312000, 1713398400, 0,
            tachotime.NewTachoTimeDriveDayData("09:00", "08:00", "17:00"),
        ),
    },
    480, 120, 45, 195, 840,
)

tt := tachotime.NewTachoTime("Tachograph Fahrer 42", []tachotime.TachoTimeDay{day}, nil)
return wc.Component(tt)
```

## Activity-Codes

Der `activity`-Parameter in `TachoTimeData` ist ein kodierter Zustand — welche Werte dein Frontend rendert, hängt vom Projekt-Mapping ab. Übliche Konvention:

| Code | Bedeutung       |
| ---- | --------------- |
| 1    | Lenken          |
| 2    | Ruhezeit        |
| 3    | Bereitschaft    |
| 4    | Sonstige Arbeit |

**Wichtig:** Das ist **kein** im xiri-go definierter Enum — prüfe das Frontend-Rendering.

## Duration-Format

Alle String-Durations (`"01:00"`, `"08:00"`) sind `HH:MM`. Die numerischen Minuten-Werte am Ende von `NewTachoTimeDay` (driving/working/breakTime/available/total) werden separat geliefert — das Frontend nutzt sie für Summary-Zeilen, die Strings für Balken-Labels.
