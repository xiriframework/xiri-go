# Enums (component/core/enums.go)

## ButtonAction

```go
ButtonActionLink     ButtonAction = "link"      // Navigate to URL
ButtonActionDialog   ButtonAction = "dialog"    // Open dialog
ButtonActionApi      ButtonAction = "api"       // API call (fire & forget)
ButtonActionDownload ButtonAction = "download"  // Download file
ButtonActionForm     ButtonAction = "form"      // Submit form
ButtonActionBack     ButtonAction = "back"      // Go back
ButtonActionClose    ButtonAction = "close"     // Close dialog
ButtonActionSave     ButtonAction = "save"      // Save form
ButtonActionHref     ButtonAction = "href"      // External link
ButtonActionGet      ButtonAction = "get"       // HTTP GET
ButtonActionPost     ButtonAction = "post"      // HTTP POST
ButtonActionPut      ButtonAction = "put"       // HTTP PUT
ButtonActionDelete   ButtonAction = "delete"    // HTTP DELETE
ButtonActionPrev     ButtonAction = "prev"      // Previous step
ButtonActionNext     ButtonAction = "next"      // Next step
```

## ButtonType

```go
ButtonTypeRaised   ButtonType = "raised"     // Filled button (primary action)
ButtonTypeBasic    ButtonType = "basic"       // Text-only
ButtonTypeStroked  ButtonType = "stroked"     // Outlined
ButtonTypeFlat     ButtonType = "flat"        // Flat filled
ButtonTypeMiniFab  ButtonType = "minifab"     // Small FAB
ButtonTypeFab      ButtonType = "fab"         // Floating Action Button
ButtonTypeIcon     ButtonType = "icon"        // Icon-only
ButtonTypeIconText ButtonType = "icontext"    // Icon + text
```

## Color

```go
ColorPrimary   Color = "primary"
ColorSecondary Color = "secondary"
ColorTertiary  Color = "tertiary"
ColorAccent    Color = "accent"
ColorWarning   Color = "warn"
ColorError     Color = "error"
ColorSuccess   Color = "success"
ColorEmerald   Color = "emerald"
ColorRed       Color = "red"
ColorYellow    Color = "yellow"
ColorGreen     Color = "green"
ColorPurple    Color = "purple"
ColorBlue      Color = "blue"
ColorOrange    Color = "orange"
ColorGray      Color = "gray"
ColorLightGray Color = "lightgray"
ColorDarkGray  Color = "darkgray"
ColorWhite     Color = "white"
ColorBlack     Color = "black"
ColorInherit   Color = "inherit"
```

## CardType

```go
CardTypeTable CardType = "table"
```

## DialogType

```go
DialogTypeForm      DialogType = "form"
DialogTypeQuestion  DialogType = "question"
DialogTypeWaiting   DialogType = "waiting"
DialogTypeTable     DialogType = "table"
DialogTypeComponent DialogType = "component"   // beliebige core.Component als Dialog-Inhalt
```

## TabHeaderPosition

```go
TabHeaderPositionAbove TabHeaderPosition = "above"
TabHeaderPositionBelow TabHeaderPosition = "below"
```

## TabAlignment

```go
TabAlignmentStart  TabAlignment = "start"
TabAlignmentCenter TabAlignment = "center"
TabAlignmentEnd    TabAlignment = "end"
```

## ExpansionDisplayMode

```go
ExpansionDisplayModeDefault ExpansionDisplayMode = "default"
ExpansionDisplayModeFlat    ExpansionDisplayMode = "flat"
```

## ExpansionTogglePosition

```go
ExpansionTogglePositionBefore ExpansionTogglePosition = "before"
ExpansionTogglePositionAfter  ExpansionTogglePosition = "after"
```

## TimelineOrientation

```go
TimelineOrientationVertical   TimelineOrientation = "vertical"    // Default
TimelineOrientationHorizontal TimelineOrientation = "horizontal"
```

## MessageType (response package)

```go
MessageSuccess MessageType = "success"
MessageError   MessageType = "error"
MessageInfo    MessageType = "info"
MessageWarning MessageType = "warning"
```

## ConditionOperator (form/field package)

```go
CondEquals    ConditionOperator = "equals"
CondNotEquals ConditionOperator = "notEquals"
CondContains  ConditionOperator = "contains"
CondGreater   ConditionOperator = "greaterThan"
CondLess      ConditionOperator = "lessThan"
CondIn        ConditionOperator = "in"
CondNotEmpty  ConditionOperator = "notEmpty"
```

## FieldType (form/field package)

```go
FieldTypeText       FieldType = "text"
FieldTypeInt        FieldType = "number"
FieldTypeBool       FieldType = "bool"
FieldTypeSelect     FieldType = "select"
FieldTypeModel      FieldType = "model"
FieldTypeModelList  FieldType = "modellist"
FieldTypeDeviceList FieldType = "devicelist"
FieldTypeDriverList FieldType = "driverlist"
FieldTypeTime       FieldType = "time"
FieldTypeTimeRange  FieldType = "timerange"
FieldTypeTimelimit  FieldType = "timelimit"
FieldTypeFile       FieldType = "file"
FieldTypeArray      FieldType = "array"
FieldTypeJson       FieldType = "json"
FieldTypeChips      FieldType = "chips"
FieldTypeGeoform    FieldType = "geoform"
FieldTypeHeader     FieldType = "header"
FieldTypeHtml       FieldType = "html"
FieldTypeInfo       FieldType = "info"
FieldTypeSerial     FieldType = "serial"
```

## Table OutputType

```go
OutputWeb   OutputType = 0  // HTML table for Angular
OutputCSV   OutputType = 1  // CSV export
OutputPDF   OutputType = 2  // PDF export
OutputExcel OutputType = 3  // Excel export
```

## Table FieldAlign

```go
AlignLeft   FieldAlign = "left"
AlignCenter FieldAlign = "center"
AlignRight  FieldAlign = "right"
```

## Table FieldFooter

```go
FooterNo     FieldFooter = "no"
FooterSum    FieldFooter = "sum"
FooterCount  FieldFooter = "count"
FooterStatic FieldFooter = "static"
```
