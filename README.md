# xiri-go

Go library for type-safe UI components, form builders, and API response types. Generates JSON structures for the [xiri-ng](https://github.com/xiriframework/xiri-ng) Angular frontend.

## Installation

```bash
go get github.com/xiriframework/xiri-go
```

## Overview

xiri-go provides Go types and builders for generating JSON that drives the xiri-ng Angular UI library. Every UI component implements the `Component` interface:

```go
type Component interface {
    Print(TranslateFunc) map[string]any
}
```

### Packages

- **component/** - UI component builders (table, form, card, dialog, stepper, tabs, etc.)
- **form/** - Form field and group builders with struct binding
- **formatter/** - Number, date, and time formatting utilities
- **response/** - HTTP response helpers for Echo framework
- **types/** - Shared type definitions
- **uicontext/** - Request context with locale and timezone

## Quick Start

```go
import (
    "github.com/xiriframework/xiri-go/component/table"
    "github.com/xiriframework/xiri-go/response"
    "github.com/xiriframework/xiri-go/uicontext"
)

type Device struct {
    ID   int64
    Name string
}

func handleDeviceTable(ctx *uicontext.UiContext, translate func(string) string) map[string]any {
    builder := table.NewBuilder[Device](ctx, translate)

    builder.IdField("id", "device.id", func(r Device) int64 { return r.ID })
    builder.TextField("name", "device.name", func(r Device) string { return r.Name })

    tbl := builder.Build()
    tbl.SetData(devices)

    return response.NewDataResponse(tbl.Print(translate))
}
```

## Requirements

- Go 1.25+
- [Echo v4](https://github.com/labstack/echo) (for HTTP response helpers)

## License

Apache-2.0
