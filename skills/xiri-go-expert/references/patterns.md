# End-to-End Patterns

Die kompletten Workflows, die man 80% der Zeit braucht. Alle Beispiele gehen davon aus, dass der Controller `pageUrl(...)` und `apiUrl(...)` Helper hat (siehe `url-routing.md`).

## 1. CRUD-Vollpacker (Page + Data + Add + Edit + Delete + Inline-Save)

```go
type Controller struct {
    svc *service.Service
}

// ---------- Page ----------
func (c *Controller) Page(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    uc := wc.UiContext()

    fg := c.buildFilterGroup(uc)
    tbl := c.buildTable(uc)
    tbl.SetURL(c.apiUrl("data"))
    tbl.SetFilter(fg)

    q := query.NewQueryWithFormGroup(fg, nil, c.apiUrl("data"), nil, nil, nil).
        WithSaveStateId("devices.filter").
        Collapsed(false)
    q.Add(tbl, uc)

    p := page.NewPage()
    p.Bread("Devices", c.pageUrl(), false)

    addBtn := button.NewLinkButton("Neu", c.pageUrl("add"),
        core.ColorPrimary, core.ButtonTypeRaised, "add", false, nil, nil)
    buttons := button.NewButtonLine("", nil)
    buttons.Add(addBtn)
    p.Add(pageheader.New("Geräte").Buttons(buttons))
    p.Add(q)

    return wc.Page(p)
}

// ---------- Data ----------
func (c *Controller) Data(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    tbl := c.buildTable(wc.UiContext())
    tbl.SetFilter(c.buildFilterGroup(wc.UiContext()))

    filters, err := tbl.LoadFilterData(ctx)
    if err != nil { return wc.BadRequest(err.Error()) }
    pg := tbl.LoadPaginationParams()

    rows, total, err := c.svc.FindDevices(filters, pg)
    if err != nil { return wc.InternalServerError(err.Error()) }

    tbl.SetData(rows)
    tbl.SetTotal(int(total))      // für Server-Side-Pagination
    return wc.Data(tbl)
}

// ---------- Add ----------
func (c *Controller) Add(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    uc := wc.UiContext()
    nameF, activeF := buildFormFields(nil)
    fb := formbuilder.NewFormBuilder(uc).AddField(nameF).AddField(activeF)

    if ctx.Request().Method == "POST" {
        fg, _, _ := fb.BuildAdd()
        if err := formbuilder.BindAndValidate(ctx, fg); err != nil {
            return wc.BadRequest(err.Error())
        }
        d := &entities.Device{}
        if nameF.Value   != nil { d.Name   = *nameF.Value }
        if activeF.Value != nil { d.Active = *activeF.Value }
        if err := c.svc.DB.Device.Create(ctx.Request().Context(), d); err != nil {
            return wc.InternalServerError(err.Error())
        }
        return wc.Goto(c.pageUrl().Print())
    }

    fields, _ := fb.BuildAddForDisplay()
    f := form.NewForm(fields, c.apiUrl("add"), nil, nil, nil, uc)
    p := page.NewPage()
    p.Bread("Devices", c.pageUrl(), false)
    p.Bread("Neu", nil, false)
    p.Add(pageheader.New("Neues Gerät"))
    p.Add(f)
    return wc.Page(p)
}

// ---------- Edit ----------
func (c *Controller) Edit(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    uc := wc.UiContext()
    id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

    d, err := c.svc.DB.Device.GetByID(id)
    if err != nil { return wc.NotFound(err.Error()) }

    nameF, activeF := buildFormFields(d)
    fb := formbuilder.NewFormBuilder(uc).AddField(nameF).AddField(activeF)

    if ctx.Request().Method == "POST" {
        fg, _, _ := fb.BuildEdit()
        if err := formbuilder.BindAndValidate(ctx, fg); err != nil {
            return wc.BadRequest(err.Error())
        }
        if nameF.Value   != nil { d.Name   = *nameF.Value }
        if activeF.Value != nil { d.Active = *activeF.Value }
        if err := c.svc.DB.Device.Update(ctx.Request().Context(), d); err != nil {
            return wc.InternalServerError(err.Error())
        }
        return wc.Goto(c.pageUrl().Print())
    }

    fields, _ := fb.BuildEditForDisplay()
    f := form.NewForm(fields, c.apiUrl("edit", strconv.FormatInt(id, 10)), nil, nil, nil, uc)
    p := page.NewPage()
    p.Bread("Devices", c.pageUrl(), false)
    p.Bread(d.Name, nil, false)
    p.Add(pageheader.New("Bearbeiten"))
    p.Add(f)
    return wc.Page(p)
}

// ---------- Delete ----------
func (c *Controller) Delete(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

    if ctx.Request().Method == "GET" {
        d, _ := c.svc.DB.Device.GetByID(id)
        return wc.DeleteDialog(d.Name)
    }

    if err := c.svc.DB.Device.Delete(ctx.Request().Context(), id); err != nil {
        return wc.InternalServerError(err.Error())
    }
    return wc.RefreshTable()
}

// ---------- Shared field builder ----------
func buildFormFields(d *entities.Device) (*field.TextField, *field.BoolField) {
    var name string
    var active bool
    if d != nil { name = d.Name; active = d.Active }
    return field.NewTextField("name", "Name", true, name),
           field.NewBoolField("active", "Aktiv", false, active)
}
```

## 2. Inline-Edit in Table

```go
// Builder
tbl := table.NewBuilder[Device]()
tbl.TextField("name", "Name", func(r Device) string { return r.Name }).
    WithEditable(true).
    WithInputType("text")
tbl.TextField("priority", "Priorität", func(r Device) string { return r.Priority }).
    WithEditableOptions(map[string]string{"low":"Niedrig","high":"Hoch"})
tbl.SetEditUrl(c.apiUrl("inline").PrintPrefix())
t := tbl.Build()

// Controller-Handler
func (c *Controller) InlineSave(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    t := c.buildTable(wc.UiContext())   // derselbe Builder wie in Data

    req, err := t.ParseInlineEdit(ctx)
    if err != nil { return wc.BadRequest(err.Error()) }

    // req.ID, req.Field, req.Value
    d, err := c.svc.DB.Device.GetByID(req.ID)
    if err != nil { return wc.NotFound(err.Error()) }
    switch req.Field {
    case "name":     d.Name = req.Value.(string)
    case "priority": d.Priority = req.Value.(string)
    default:
        return wc.BadRequest("unknown field")
    }
    if err := c.svc.DB.Device.Update(ctx.Request().Context(), d); err != nil {
        return wc.InternalServerError(err.Error())
    }

    return wc.Component(response.NewReturnInlineEdit().
        WithUpdates(map[string]any{req.Field: req.Value}).
        WithMessage("Gespeichert", response.MessageSuccess))
}
```

## 3. Dashboard mit StatGrid + Table

```go
func (c *Controller) Dashboard(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    uc := wc.UiContext()

    // KPI-Stats
    sg := statgrid.New().Title("Übersicht").Columns(4)
    sg.Add(stat.New(strconv.Itoa(c.svc.CountDevices()), "Geräte").Icon("router"))
    sg.Add(stat.New(strconv.Itoa(c.svc.CountActive()), "Aktiv").
        Icon("check_circle").IconColor(core.ColorSuccess))
    sg.Add(stat.New(strconv.Itoa(c.svc.CountErrors()), "Fehler").
        Icon("error").IconColor(core.ColorError))

    // Jüngste Geräte als Card-Tabelle
    recent := c.buildTable(uc)
    recent.SetData(c.svc.RecentDevices(10))
    recent.SetTitle("Zuletzt geändert")

    p := page.NewPage()
    p.Add(pageheader.New("Dashboard").Icon("dashboard", core.ColorPrimary))
    p.Add(sg)
    p.Add(recent)
    return wc.Page(p)
}
```

## 4. Stepper-Wizard (Multi-Step Form)

```go
func (c *Controller) WizardPage(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    uc := wc.UiContext()

    fb1 := formbuilder.NewFormBuilder(uc).
        AddField(field.NewTextField("firstName", "Vorname", true, ""))
    fb2 := formbuilder.NewFormBuilder(uc).
        AddField(field.NewTextField("street", "Straße", true, ""))
    fb3 := formbuilder.NewFormBuilder(uc).
        AddField(field.NewBoolField("tosAccepted", "AGB akzeptiert", true, false))

    step1Fields, _ := fb1.BuildAddForDisplay()
    step2Fields, _ := fb2.BuildAddForDisplay()
    step3Fields, _ := fb3.BuildAddForDisplay()

    s, err := stepper.NewStepper(
        c.apiUrl("wizard", "save").PrintPrefix(),
        3,
        []string{"Person", "Adresse", "Bestätigung"},
        []stepper.StepFields{step1Fields, step2Fields, step3Fields},
        "Zurück", "Weiter", "Fertig", "",
    )
    if err != nil { return wc.InternalServerError(err.Error()) }

    p := page.NewPage()
    p.Add(pageheader.New("Registrierung"))
    p.Add(s)
    return wc.Page(p)
}
```

## 5a. MassDelete — Bulk-Delete via Confirm-Dialog

Der klassische Flow: User selektiert Zeilen → klickt "Löschen" → Confirm-Dialog → POST löscht alle. Zwei Handler (GET = Dialog, POST = Delete), **gleiche URL**.

```go
// Table-Setup: Select-Button mit action: dialog
b.SetSelectButtons([]*button.TableButton{
    button.NewTableButton(
        core.ButtonActionDialog,
        "delete",
        c.apiUrl("mass-delete"),
        "Löschen",
        core.ColorWarning,
        false,
        nil,
    ),
})

// Handler
func (c *Controller) MassDelete(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)

    // IDs aus Request auslesen (unterstützt verschiedene Payload-Formate)
    ids, _, _, err := dialog.ExtractMultiSelectRequest(ctx)
    if err != nil { return wc.BadRequest(err.Error()) }
    if len(ids) == 0 { return wc.BadRequest("keine Zeilen gewählt") }

    if ctx.Request().Method == "GET" {
        dlg := dialog.NewDialogFormMultiDelete(
            c.apiUrl("mass-delete"),    // dieselbe URL — POST löscht
            ids,
            fmt.Sprintf("%d Geräte wirklich löschen?", len(ids)),
            nil, nil, nil,
        )
        return wc.Component(dlg)
    }

    if err := c.svc.DB.Device.DeleteMany(ctx.Request().Context(), ids); err != nil {
        return wc.InternalServerError(err.Error())
    }
    return wc.Component(response.NewReturnRefreshTable().
        WithMessage(fmt.Sprintf("%d gelöscht", len(ids)), response.MessageSuccess))
}
```

## 5b. MassEdit — gleichzeitige Änderung mehrerer Zeilen

Etwas komplexer: nur Felder ändern, die der User aktiv im MassEdit-Form setzt — leere Felder bleiben in der DB unangetastet. Typische UX: Dialog mit wenigen Feldern (z.B. nur "Status", "Owner"), nicht alle Edit-Felder wie bei Single-Edit.

```go
// Table-Setup
b.SetSelectButtons([]*button.TableButton{
    button.NewTableButton(
        core.ButtonActionDialog,
        "edit",
        c.apiUrl("mass-edit"),
        "Gemeinsam bearbeiten",
        core.ColorPrimary,
        false,
        nil,
    ),
})

// Handler: GET = Form-Dialog, POST = Bulk-Update
func (c *Controller) MassEdit(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    uc := wc.UiContext()

    ids, _, _, err := dialog.ExtractMultiSelectRequest(ctx)
    if err != nil { return wc.BadRequest(err.Error()) }
    if len(ids) == 0 { return wc.BadRequest("keine Zeilen gewählt") }

    // Felder für MassEdit (alle optional — User lässt leer = keine Änderung)
    statusF := field.NewSelectField("status", "Status", false, statusOptions)
    ownerF  := field.NewModelField("ownerId", "Owner", false, "user", 0)
    noteF   := field.NewTextField("note", "Notiz anhängen", false, "")

    fb := formbuilder.NewFormBuilder(uc).
        AddField(statusF).AddField(ownerF).AddField(noteF)

    if ctx.Request().Method == "GET" {
        fields, _ := fb.BuildAddForDisplay()
        header := fmt.Sprintf("%d Geräte bearbeiten", len(ids))
        dlg := dialog.NewDialogForm(
            fields,
            c.apiUrl("mass-edit"),   // POST: dieselbe URL
            &header,
            map[string]any{"ids": ids},   // IDs durchreichen!
            nil, nil,
        ).WithOption("size", "md")
        return wc.Component(dlg)
    }

    // POST: Bulk-Update
    fg, _, _ := fb.BuildAdd()
    if err := formbuilder.BindAndValidate(ctx, fg); err != nil {
        return wc.BadRequest(err.Error())
    }

    updates := map[string]any{}
    if statusF.Value != 0 { updates["status"] = statusF.Value }
    if ownerF.Value  != 0 { updates["owner_id"] = ownerF.Value }
    if noteF.Value != nil && *noteF.Value != "" {
        updates["note"] = *noteF.Value
    }
    if len(updates) == 0 {
        return wc.BadRequest("keine Änderung angegeben")
    }

    if err := c.svc.DB.Device.
        Where("id IN ?", ids).
        Updates(updates); err != nil {
        return wc.InternalServerError(err.Error())
    }

    return wc.Component(response.NewReturnRefreshTable().
        WithMessage(fmt.Sprintf("%d aktualisiert", len(ids)), response.MessageSuccess))
}
```

**Wichtig:**
- `extra.ids` im Dialog mitschicken → das Frontend hängt sie beim Submit wieder an den Request, sodass der POST-Handler wieder `ExtractMultiSelectRequest` nutzen kann.
- **Leere Felder = keine Änderung** — deshalb alle MassEdit-Felder `required: false`. Check vor dem `Updates(...)`-Call, dass mindestens ein Feld gesetzt ist.
- GORM `.Updates(map[string]any)` ignoriert zero-values nicht — deshalb die Map nur mit tatsächlichen Änderungen füllen.

## 5. Bulk-Actions (einfach, ohne Dialog)

```go
// Table mit Select-Buttons
tbl := c.buildTable(uc)
tbl.SetSelectButtons([]*button.TableButton{
    button.NewTableButton(core.ButtonActionApi, "delete",
        c.apiUrl("bulk-delete"), "Löschen", core.ColorError, false, nil),
    button.NewTableButton(core.ButtonActionApi, "export",
        c.apiUrl("bulk-export"), "Export", core.ColorPrimary, false, nil),
})

// Bulk-Handler
func (c *Controller) BulkDelete(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    var req struct{ IDs []int64 `json:"ids"` }
    if err := ctx.Bind(&req); err != nil { return wc.BadRequest(err.Error()) }
    if err := c.svc.DB.Device.DeleteMany(ctx.Request().Context(), req.IDs); err != nil {
        return wc.InternalServerError(err.Error())
    }
    return wc.Component(response.NewReturnRefreshTable().
        WithMessage(fmt.Sprintf("%d gelöscht", len(req.IDs)), response.MessageSuccess))
}
```

## 6. AJAX-Card (lazy Server-Content)

```go
// Page
card := card.NewCard(core.CardTypeTable, nil, "Live-Status", "", "", "", true, false, "")
card.SetURL(c.apiUrl("card", "status")).WithReload(true)
p.Add(card)

// Endpoint — gibt nur den Card-Content zurück
func (c *Controller) CardStatus(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    content := card.NewCardListContent([]card.CardListContentLine{
        {Name: "Online",  Content: strconv.Itoa(c.svc.CountActive())},
        {Name: "Offline", Content: strconv.Itoa(c.svc.CountOffline())},
    })
    inner := card.NewCardList("Live-Status", content)
    return wc.Component(inner)
}
```

## 7. Tabs mit Lazy-Loading

```go
t := tabs.NewTabs().WithLazy(true).WithUnload(true)

tab1 := tabs.NewTab("Übersicht").WithIcon("info")
tab1.AddContent(overviewComponent)
t.AddTab(tab1)

tab2 := tabs.NewTab("Historie").WithIcon("history")
tab2.AddContent(historyTable)
t.AddTab(tab2)

p.Add(t)
```

## 8. Delete-Dialog + Custom-Message

Der Standard-DeleteDialog via `wc.DeleteDialog(name)` reicht meistens. Wenn du mehr Kontrolle brauchst:

```go
func (c *Controller) ConfirmArchive(ctx echo.Context) error {
    wc := webcontext.GetWebContext(ctx)
    content := &dialog.DialogQuestionContent{
        Icon:     "archive",
        Question: "Diesen Eintrag wirklich archivieren?",
    }
    d := dialog.NewDialog(core.DialogTypeQuestion, "Archivieren", content,
        []*button.Button{
            button.NewSimpleCloseButton("Abbrechen"),
            button.NewSimpleApiButton("Archivieren",
                c.apiUrl("archive", ctx.Param("id")).PrintPrefix(),
                core.ColorWarning),
        }, nil, nil)
    return wc.Component(d)
}
```
