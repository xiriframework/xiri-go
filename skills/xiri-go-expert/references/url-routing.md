# URL-Routing, Sidebar-Prefixes & Breadcrumbs

`component/url` + `component/page` bilden die URL-Infrastruktur. **Niemals** Strings manuell zusammenbauen — `*xurl.Url` ist klein und stringfrei-safe.

## `xurl.Url` — die einzige URL-Abstraktion

Quelle: `component/url/url.go`

```go
type Url struct {
    url    string
    prefix string
}

func NewUrl(url string) *Url                            // ohne Prefix
func NewUrlPrefix(url string, prefix string) *Url        // mit Prefix

func (u *Url) Add(url string) *Url                       // append-Chain
func (u *Url) AddPrefix(prefix string) *Url              // prepend-Chain

func (u *Url) Print() string                             // ohne Prefix (Frontend-Link)
func (u *Url) PrintPrefix() string                       // mit Prefix (API-POST-Target)
```

Beispiel:

```go
u := xurl.NewUrlPrefix("/Portal/Device/Table", "/api")
u.Print()        // "/Portal/Device/Table"
u.PrintPrefix()  // "/api/Portal/Device/Table"

u2 := xurl.NewUrl("/Portal").Add("Device").Add("Edit").Add("7")
u2.Print()       // "/Portal/Device/Edit/7"
```

## Zwei Welten: Frontend-Routes vs. API-Endpoints

| URL-Kategorie           | Prefix? | Wer konsumiert              | Beispielwerte                               |
| ----------------------- | ------- | --------------------------- | ------------------------------------------- |
| Angular-Route (Link)    | **nein** | `<a routerLink>`            | `/devices`, `/devices/edit/42`              |
| API-Endpoint (POST/GET) | **ja**   | HttpClient, `form.NewForm`  | `/api/v1/devices/save`                      |
| Breadcrumb              | **nein** | Angular Router              | `/devices`                                  |
| Button `action: link`   | **nein** | Router                      | `/devices/edit/42`                          |
| Button `action: api`    | **ja**   | HttpClient                  | `/api/v1/devices/save`                      |
| Button `action: dialog` | **ja**   | fetch + MatDialog           | `/api/v1/devices/42/delete`                 |
| `table.SetURL(...)`     | **ja**   | HttpClient                  | `/api/v1/devices/data`                      |
| `form.NewForm(_, url, …)` | **ja** | HttpClient                  | `/api/v1/devices/save`                      |

Deshalb: **`Print()` in Links/Breadcrumbs, `PrintPrefix()` in API-Zielen**. Builder wie `button.NewLinkButton` und `table.SetURL` akzeptieren `*xurl.Url` direkt — sie picken intern die richtige Variante.

## Controller-Helper für saubere Prefix-Nutzung

Wenn jeder Controller zwei Helpers hat, braucht man die Prefix-Logik nie wieder händisch:

```go
type Controller struct {
    svc        *service.Service
    pagePrefix string   // z.B. "/portal"
    apiPrefix  string   // z.B. "/api/v1/portal"
}

func (c *Controller) pageUrl(seg ...string) *xurl.Url {
    u := xurl.NewUrl(c.pagePrefix + "/devices")
    for _, s := range seg { u.Add(s) }
    return u
}

func (c *Controller) apiUrl(seg ...string) *xurl.Url {
    u := xurl.NewUrlPrefix(c.pagePrefix+"/devices", c.apiPrefix)
    for _, s := range seg { u.Add(s) }
    return u
}

// Benutzung im Handler
tbl.SetURL(c.apiUrl("data"))                              // POST-Target
buttons.Add(button.NewLinkButton("Neu", c.pageUrl("add"), // Angular-Route
    core.ColorPrimary, core.ButtonTypeRaised, "", false, nil, nil))
```

Der Prefix wandert damit in **einem** Konfigurations-Punkt (Service-Bootstrap) und nicht in 50 String-Konkatenationen.

## Breadcrumbs

```go
p := page.NewPage()
p.Bread("Start", xurl.NewUrl("/"), false)
p.Bread("Devices", c.pageUrl(), false)
p.Bread("Detail", nil, false)     // letztes Segment: kein Link
```

Das dritte Argument (`bool`) ist `extern` — `true` öffnet in neuem Tab.

## Sidebar-Navigation (Projekt-Seite)

xiri-go liefert **keine** eigene Sidebar-Komponente — Sidebars werden im Frontend (xiri-ng) gerendert. Aus Go-Sicht bedeutet das: Der Controller gibt einfach Page-URLs zurück, **ohne** API-Prefix. Für die Sidebar-Konfiguration konsumiert xiri-ng `XiriNavigationField[]`:

```json
{
  "prefix": "/portal/",
  "fields": [
    { "name": "Devices", "link": "/portal/devices", "icon": "router" },
    { "name": "Users",   "link": "/portal/users",   "icon": "group" }
  ]
}
```

Der Go-Endpoint der diese Struktur liefert, baut sie mit `xurl.NewUrl(pagePrefix + "/devices").Print()` — d.h. **dieselben** URLs, die auch in Breadcrumbs und Link-Buttons landen. **`prefix`** in der Sidebar-Config ist der Route-Prefix, der im Active-Match **entfernt** werden soll (z.B. `/portal/` damit `/portal/devices` als `devices` hervorgehoben wird).

**Regel:** `pagePrefix` in Go = `prefix` in Sidebar-Config = gemeinsamer Anfang aller Page-URLs.

### Aufklappbare Menüs (bis zu 3 Ebenen)

Ein `XiriNavigationField` ist entweder ein **Link** (`link`), ein **externer Link** (`extern`) oder ein **aufklappbares Menü** (`menu: true` + `sub: XiriNavigationField[]`). Seit der Sidebar-Erweiterung können Sub-Einträge ihrerseits `menu: true` mit eigenem `sub` sein — damit sind **drei Ebenen** möglich (Top → Sub → Sub-Sub). Tiefer rendert das Frontend bewusst nicht; die dritte Ebene sind Blätter (`link`/`extern`).

```json
{
  "prefix": "/portal/",
  "fields": [
    {
      "name": "Forms", "icon": "edit_note", "menu": true,
      "sub": [
        { "name": "Basic", "link": "/portal/forms/basic", "icon": "text_fields" },
        {
          "name": "Advanced", "icon": "tune", "menu": true,
          "sub": [
            { "name": "Select",  "link": "/portal/forms/advanced/select" },
            { "name": "Special", "link": "/portal/forms/advanced/special" }
          ]
        }
      ]
    }
  ]
}
```

`menu`-Knoten brauchen kein `link` — der Klick togglet nur das Auf-/Zuklappen (`showSubmenu` wird vom Frontend verwaltet, nicht setzen). Alle URLs in `link` bleiben **ohne** API-Prefix, auf jeder Ebene.

### Einträge weglassen (`hide`)

Soll ein Nutzer eine Route nicht sehen, setzt der Sidebar-Endpoint `"hide": true` — das Frontend
rendert den Eintrag dann **gar nicht** (kein `routerLink` im DOM), auf allen drei Ebenen. Ein
verstecktes Kind aktiviert und öffnet auch seinen Parent nicht.

```json
{ "name": "Admin", "link": "/portal/admin", "icon": "shield", "hide": true }
```

**Kein Berechtigungs-Contract.** Client-Filtern ist keine Autorisierung — der Server muss die
Routes ohnehin absichern. `hide` erspart nur den Klick ins Leere. Genauso gut (und einfacher):
den Eintrag serverseitig weglassen. `hide` ist dann sinnvoll, wenn dieselbe Struktur für alle
gebaut und pro Nutzer nur maskiert wird.

## Active-State Matching

Wenn die Sidebar eine bestimmte Section hervorheben soll, bei mehreren URL-Varianten (`/devices`, `/devices/edit/42`, `/devices/add`), kann man im `XiriNavigationField.path` ein Regex setzen — das funktioniert auf **jeder** Ebene (auch auf einem verschachtelten `menu`-Knoten) und ist reine Frontend-Sache in der JSON-Response des Sidebar-Endpoints.

Bei einem aktiven (Tief-)Link klappt das Frontend automatisch **alle Vorfahren** auf und markiert den aktiven Eintrag — man muss also keinen Expand-Zustand vorberechnen. `active`/`showSubmenu` nicht im Backend setzen; sie werden aus der aktuellen Route abgeleitet.

## Häufige Fallen

- **Falsche Mischung**: `tbl.SetURL(xurl.NewUrl("/data"))` ohne Prefix → xiri-ng POSTed an `/data` und bekommt 404. Immer `NewUrlPrefix` für API-URLs.
- **String-Konkatenation**: `fmt.Sprintf("%s/edit/%d", prefix, id)` in Handlern — gerne schleicht sich hier der API-Prefix in einen Link ein. Stattdessen `c.pageUrl("edit", strconv.FormatInt(id, 10)).Print()`.
- **Doppelte Slashes**: `NewUrl("/devices/").Add("/edit")` → `/devices//edit`. `Add` prefixt selbst einen Slash; `Add` **nicht** mit führendem Slash füttern.
- **Falscher Prefix beim Weiterdelegieren**: `u.AddPrefix("/tenant-x")` ist **cumulativ** — mehrfaches `AddPrefix` stackt. Wenn du nur einen Prefix setzen willst, verwende `NewUrlPrefix` beim Erzeugen.
