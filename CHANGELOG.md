# Changelog

Alle nennenswerten Änderungen an `github.com/xiriframework/xiri-go` werden hier festgehalten.

Format nach [Keep a Changelog](https://keepachangelog.com/de/1.1.0/),
Versionierung nach [Semantic Versioning](https://semver.org/lang/de/).

## [Unreleased]

Änderungen seit `v0.2.30`, noch nicht getaggt.

### Added

- **`component/progress`**: `Progress` — einzelne determinate Fortschrittsanzeige („current of total"), inkl. `Indeterminate()`-Modus. Für Share-of-Sum weiterhin `MultiProgress`.
- **`component/bulletchart`**: Chart-Builder als kompakte Gauge-Alternative (`value` / `target` / `max` / `label`).
- **`component/stat`**: `Stat.Reference(text)` — muted Benchmark-/Anker-Zeile neben dem Wert (z. B. „Gate ≥ 1,1"), macht die Zahl auf einen Blick einordbar.
- **`form/field`**: `radio`-Feldtyp — Single-Select über `SelectField` mit `type=radio`, für kleine Optionsmengen.

[Unreleased]: https://github.com/xiriframework/xiri-go/compare/v0.2.30...HEAD
