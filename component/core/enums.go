package core

// ButtonAction represents the action type for a button
type ButtonAction string

const (
	ButtonActionLink     ButtonAction = "link"
	ButtonActionDialog   ButtonAction = "dialog"
	ButtonActionApi      ButtonAction = "api"
	ButtonActionDownload ButtonAction = "download"
	ButtonActionForm     ButtonAction = "form"
	ButtonActionBack     ButtonAction = "back"
	ButtonActionClose    ButtonAction = "close"
	ButtonActionSave     ButtonAction = "save"
	ButtonActionHref     ButtonAction = "href"
	ButtonActionGet      ButtonAction = "get"
	ButtonActionPost     ButtonAction = "post"
	ButtonActionPut      ButtonAction = "put"
	ButtonActionDelete   ButtonAction = "delete"
	ButtonActionPrev     ButtonAction = "prev"
	ButtonActionNext     ButtonAction = "next"
)

// ButtonType represents the visual style of a button
type ButtonType string

const (
	ButtonTypeRaised   ButtonType = "raised"
	ButtonTypeBasic    ButtonType = "basic"
	ButtonTypeStroked  ButtonType = "stroked"
	ButtonTypeFlat     ButtonType = "flat"
	ButtonTypeMiniFab  ButtonType = "minifab"
	ButtonTypeFab      ButtonType = "fab"
	ButtonTypeIcon     ButtonType = "icon"
	ButtonTypeIconText ButtonType = "icontext"
)

// Color represents theme colors for components
type Color string

const (
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
)

// CardType represents the type of content in a card
type CardType string

const (
	CardTypeTable CardType = "table"
)

// DialogType represents the type of dialog
type DialogType string

const (
	DialogTypeForm     DialogType = "form"
	DialogTypeQuestion DialogType = "question"
	DialogTypeWaiting  DialogType = "waiting"
	DialogTypeTable    DialogType = "table"
)

// ListItemColor represents colors for list items
type ListItemColor string

const (
	ListItemColorPrimary ListItemColor = "primary"
	ListItemColorWarn    ListItemColor = "warn"
	ListItemColorError   ListItemColor = "error"
	ListItemColorSuccess ListItemColor = "success"
	ListItemColorRed     ListItemColor = "red"
	ListItemColorGreen   ListItemColor = "green"
	ListItemColorYellow  ListItemColor = "yellow"
	ListItemColorGray    ListItemColor = "gray"
	ListItemColorOrange  ListItemColor = "orange"
)

// TabHeaderPosition represents where tab headers are positioned
// Angular: XiriTabsSettings.headerPosition
type TabHeaderPosition string

const (
	TabHeaderPositionAbove TabHeaderPosition = "above"
	TabHeaderPositionBelow TabHeaderPosition = "below"
)

// TabAlignment represents horizontal alignment of tab headers
// Angular: XiriTabsSettings.alignTabs
type TabAlignment string

const (
	TabAlignmentStart  TabAlignment = "start"
	TabAlignmentCenter TabAlignment = "center"
	TabAlignmentEnd    TabAlignment = "end"
)

// ExpansionDisplayMode represents the display mode of an expansion panel accordion
// Angular: XiriExpansionSettings.displayMode
type ExpansionDisplayMode string

const (
	ExpansionDisplayModeDefault ExpansionDisplayMode = "default"
	ExpansionDisplayModeFlat    ExpansionDisplayMode = "flat"
)

// ExpansionTogglePosition represents where the expand/collapse toggle is positioned
// Angular: XiriExpansionSettings.togglePosition
type ExpansionTogglePosition string

const (
	ExpansionTogglePositionBefore ExpansionTogglePosition = "before"
	ExpansionTogglePositionAfter  ExpansionTogglePosition = "after"
)
