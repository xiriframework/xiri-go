# FormBuilder & FormGroup (form/builder, form/group)

## FormBuilder

Import: `"github.com/xiriframework/xiri-go/form/builder"`

Fluent API für Formular-Erstellung mit Type-Safe Value Binding.

### Konstruktor

```go
fb := builder.NewFormBuilder(ctx) // ctx = *core.UiContext
```

### Fields hinzufügen

```go
fb.AddField(field.NewTextField("name", "vehicle.name", true, ""))
fb.AddField(field.NewModelField("group_id", "vehicle.group", true, "group", 0))
fb.AddField(field.NewTextField("plate", "vehicle.plate", false, ""))
fb.AddField(field.NewBoolField("active", "vehicle.active", false, false))
```

Das letzte Argument ist der Default als **Wert**, nicht als Pointer — `nil` kompiliert nicht.

### Build-Methoden

#### BuildAdd — Neues Objekt erstellen

```go
fg, defaults, err := fb.BuildAdd()
// fg: *FormGroup mit allen Fields
// defaults: map[string]interface{} mit Default-Werten
// Verwendet für: Formular-Anzeige (leeres Formular)
```

#### BuildEdit — Bestehendes Objekt bearbeiten

```go
// Zuerst Fields mit aktuellen Werten erstellen:
fb.AddField(field.NewTextField("name", "vehicle.name", true, vehicle.Name))
fb.AddField(field.NewBoolField("active", "vehicle.active", false, vehicle.Active))

fg, values, err := fb.BuildEdit()
// fg: *FormGroup mit allen Fields
// values: map[string]interface{} mit aktuellen Werten
```

#### BuildAddForDisplay / BuildEditForDisplay — Direkt für Frontend

```go
fields, err := fb.BuildAddForDisplay()
// fields: []map[string]interface{} — direkt für form.NewForm()
```

### OnEditValueCheck Hook

```go
fb.OnEditValueCheck = func(fg *group.FormGroup, values map[string]interface{}) error {
    // Custom Validierung nach dem Laden der Edit-Werte
    return nil
}
```

## FormGroup

Import: `"github.com/xiriframework/xiri-go/form/group"`

Container für FormFields mit Context-aware Options Loading.

### Konstruktoren

```go
fg := group.NewFormGroup(fields)                  // Ohne Context
fg := group.NewFormGroupWithContext(fields, ctx)   // Mit Context (lädt automatisch Options)
```

### Kernmethoden

```go
fg.GetFields()                    // []FormField — alle Fields
fg.GetField("name")               // (FormField, bool) — einzelnes Field
fg.GetDefaults()                   // map[string]interface{} — Default-Werte
fg.GetRequiredFields()             // []FormField — nur Required
fg.GetOptionalFields()             // []FormField — nur Optional
fg.GetFieldIDs()                   // []string — alle Field-IDs
fg.SetContext(ctx)                 // Context setzen + Options laden
```

### Export für Frontend

```go
fields := fg.ExportForFrontend()
// []map[string]interface{} — Fields ohne Werte

fields := fg.ExportForFrontendWithValues(values)
// []map[string]interface{} — Fields mit Werten
```

### Parsing & Validierung

```go
values, err := fg.ParseValues(rawData)          // Werte parsen
err := fg.ValidateValues(values)                 // Werte validieren
values, err := fg.ParseAndValidate(rawData)      // Beides
```

## BindAndValidate (Request Binding)

Import: `"github.com/xiriframework/xiri-go/form/builder"`

Extrahiert Form-Daten aus Echo-Request und bindet sie an FormGroup-Fields.

```go
if err := builder.BindAndValidate(c, fg); err != nil {
    return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
}
```

Nach dem Binding sind die Werte direkt auf den Fields verfügbar:

```go
nameField, _ := fg.GetField("name")
name := *nameField.(*field.TextField).Value          // string

groupField, _ := fg.GetField("group_id")
groupID := groupField.(*field.ModelField).Value       // int32

activeField, _ := fg.GetField("active")
active := *activeField.(*field.BoolField).Value       // bool

statusField, _ := fg.GetField("status")
status := statusField.(*field.SelectField).Value      // int32
```

### BindFromMap (Alternative)

```go
err := builder.BindFromMap(formData, fg)
// formData: map[string]interface{} — bereits geparstes Formular
```

## Vollständiges Beispiel: Add + Edit + Save

```go
// fields.go — Shared Field-Definition
func vehicleFields(ctx *core.UiContext, v *Vehicle) *builder.FormBuilder {
    fb := builder.NewFormBuilder(ctx)

    var name, plate string
    var groupID int32
    var active bool
    if v != nil {
        name = v.Name
        groupID = v.GroupID
        plate = v.Plate
        active = v.Active
    }

    fb.AddField(field.NewTextField("name", "vehicle.name", true, name))

    mf := field.NewModelField("group_id", "vehicle.group", true, "group", groupID)
    mf.SetLoaderFunc(func(ctx *core.UiContext) ([]field.ModelOption, error) {
        return loadGroups(ctx)
    })
    fb.AddField(mf)

    fb.AddField(field.NewTextField("plate", "vehicle.plate", false, plate))
    fb.AddField(field.NewBoolField("active", "vehicle.active", false, active))

    return fb
}

// handler_add.go
func HandleVehicleAdd(c echo.Context) error {
    ctx := getUiContext(c)
    fb := vehicleFields(ctx, nil)

    fg, defaults, err := fb.BuildAdd()
    if err != nil {
        return err
    }

    f := form.NewForm(
        fg.ExportForFrontendWithValues(defaults),
        "/api/vehicle/add",
        nil, nil, "", ctx,
    )

    p := page.NewPage()
    p.Bread("Fahrzeuge", "/vehicles", false)
    p.Bread("Hinzufügen", "", false)
    p.Add(f)
    return c.JSON(http.StatusOK, p.Print(ctx))
}

// handler_edit.go
func HandleVehicleEdit(c echo.Context) error {
    ctx := getUiContext(c)
    id := getParamID(c)

    vehicle, err := db.GetVehicle(id)
    if err != nil {
        return c.JSON(http.StatusNotFound, response.NewErrorResponse("not found"))
    }

    fb := vehicleFields(ctx, vehicle)
    fg, values, err := fb.BuildEdit()
    if err != nil {
        return err
    }

    f := form.NewForm(
        fg.ExportForFrontendWithValues(values),
        fmt.Sprintf("/api/vehicle/%d/edit", id),
        nil, nil, "", ctx,
    )

    p := page.NewPage()
    p.Bread("Fahrzeuge", "/vehicles", false)
    p.Bread(vehicle.Name, "", false)
    p.Add(f)
    return c.JSON(http.StatusOK, p.Print(ctx))
}

// handler_save.go
func HandleVehicleSave(c echo.Context) error {
    ctx := getUiContext(c)
    fb := vehicleFields(ctx, nil)

    fg, _, err := fb.BuildAdd()
    if err != nil {
        return err
    }

    if err := builder.BindAndValidate(c, fg); err != nil {
        return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
    }

    nameF, _ := fg.GetField("name")
    groupF, _ := fg.GetField("group_id")
    plateF, _ := fg.GetField("plate")
    activeF, _ := fg.GetField("active")

    vehicle := &Vehicle{
        Name:    *nameF.(*field.TextField).Value,
        GroupID: groupF.(*field.ModelField).Value,
        Plate:   *plateF.(*field.TextField).Value,
        Active:  *activeF.(*field.BoolField).Value,
    }

    if err := db.Create(vehicle).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
    }

    return c.JSON(http.StatusOK,
        response.NewReturnGoto(fmt.Sprintf("/vehicles/%d", vehicle.ID)).
            WithMessage(ctx.SafeTranslate("saved"), response.MessageSuccess))
}
```
