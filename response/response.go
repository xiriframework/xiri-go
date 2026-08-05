// Package response provides type-safe API response structs that serialize directly to JSON.
package response

// SuccessResponse is a marker interface for all success response types.
// Use this for type-safe function parameters that accept any success response.
//
// Example:
//
//	func HandleRequest(response SuccessResponse) {
//	    // Accepts: ReturnRefreshTable, ReturnRefreshPage, ReturnGoto, ReturnDone, ReturnMessage
//	}
type SuccessResponse interface {
	isSuccessResponse()
}

// ButtonPatch describes optional overrides applied to the initiating button after a response
// (e.g. when a worker finishes): change its label, color, icon, type, hint or disabled state.
// Attached via WithButton on a response and merged onto the button in the frontend.
//
// JSON output (only set fields): {"text":"Erledigt","color":"success","disabled":true}
type ButtonPatch struct {
	Text     string `json:"text,omitempty"`
	Color    string `json:"color,omitempty"`
	Icon     string `json:"icon,omitempty"`
	Type     string `json:"type,omitempty"`
	Hint     string `json:"hint,omitempty"`
	Disabled *bool  `json:"disabled,omitempty"`
}

// NewButtonPatch creates an empty button patch; chain the With* helpers to set fields.
func NewButtonPatch() ButtonPatch { return ButtonPatch{} }

// WithText sets the new button label.
func (b ButtonPatch) WithText(text string) ButtonPatch { b.Text = text; return b }

// WithColor sets the new button color (e.g. "success", "warn", "primary").
func (b ButtonPatch) WithColor(color string) ButtonPatch { b.Color = color; return b }

// WithIcon sets the new button icon.
func (b ButtonPatch) WithIcon(icon string) ButtonPatch { b.Icon = icon; return b }

// WithType sets the new button type (e.g. "raised", "flat").
func (b ButtonPatch) WithType(t string) ButtonPatch { b.Type = t; return b }

// WithHint sets the new button tooltip/hint.
func (b ButtonPatch) WithHint(hint string) ButtonPatch { b.Hint = hint; return b }

// Disable marks the button disabled.
func (b ButtonPatch) Disable() ButtonPatch { v := true; b.Disabled = &v; return b }

// Enable marks the button enabled.
func (b ButtonPatch) Enable() ButtonPatch { v := false; b.Disabled = &v; return b }

// MessageType defines the type of snackbar/toast message shown to the user.
type MessageType string

const (
	MessageSuccess MessageType = "success"
	MessageError   MessageType = "error"
	MessageInfo    MessageType = "info"
	MessageWarning MessageType = "warning"
)

// Message is an embeddable struct that adds optional snackbar/toast messages to any response.
// When embedded, the message fields are omitted from JSON if empty (backward-compatible).
type Message struct {
	MessageText string      `json:"message,omitempty"`
	MessageType MessageType `json:"messageType,omitempty"`
}

// ReturnRefreshPage represents a refresh page response.
//
// JSON output: {"done": true, "refresh": "page"}
// With message: {"done": true, "refresh": "page", "message": "Saved", "messageType": "success"}
//
// Use case: Operation completed, reload the current page
type ReturnRefreshPage struct {
	Done    bool   `json:"done"`    // Always true
	Refresh string `json:"refresh"` // Always "page"
	Message
}

func (r ReturnRefreshPage) isSuccessResponse() {}

// WithMessage returns a copy with the given message and type.
func (r ReturnRefreshPage) WithMessage(text string, msgType MessageType) ReturnRefreshPage {
	r.MessageText = text
	r.MessageType = msgType
	return r
}

// ReturnFields represents a field patch response for a form reload.
//
// JSON output: {"fields": {"tags": {"list": [...], "required": true}}}
// With message: {"fields": {...}, "message": "Liste aktualisiert", "messageType": "info"}
//
// Use case: a field declared SetReloadOn(...) and one of its trigger values changed. The
// frontend merges each entry into the matching field definition. Deliberately carries no
// "done" key - a patch is not a completed action.
//
// Build the map with FormGroup.ExportPatch().
type ReturnFields struct {
	Fields map[string]interface{} `json:"fields"`
	Message
}

func (r ReturnFields) isSuccessResponse() {}

// WithMessage returns a copy with the given message and type.
func (r ReturnFields) WithMessage(text string, msgType MessageType) ReturnFields {
	r.MessageText = text
	r.MessageType = msgType
	return r
}

// ReturnRefreshTable represents a refresh table response.
//
// JSON output: {"done": true, "refresh": "table"}
// With message: {"done": true, "refresh": "table", "message": "Row deleted", "messageType": "success"}
//
// Use case: Operation completed on table row, reload the table data
type ReturnRefreshTable struct {
	Done    bool   `json:"done"`    // Always true
	Refresh string `json:"refresh"` // Always "table"
	Message
}

func (r ReturnRefreshTable) isSuccessResponse() {}

// WithMessage returns a copy with the given message and type.
func (r ReturnRefreshTable) WithMessage(text string, msgType MessageType) ReturnRefreshTable {
	r.MessageText = text
	r.MessageType = msgType
	return r
}

// ReturnInlineEdit represents an inline edit save response.
// The frontend patches row fields from Updates, then handles refresh/goto via callReturn.
//
// JSON output examples:
//
//	Updates only:    {"done": true, "updates": {"price": ["1.299,00", 1299], "lastModified": "11.03.2026"}}
//	With message:    {"done": true, "updates": {"status": "Active"}, "message": "Saved", "messageType": "success"}
//	Refresh table:   {"done": true, "refresh": "table"}
//	Goto:            {"done": true, "goto": "/other/page"}
//	Combined:        {"done": true, "updates": {"status": "Active"}, "refresh": "table", "message": "Saved", "messageType": "success"}
//
// Use case: Inline edit completed — patch row fields, optionally trigger navigation/reload
type ReturnInlineEdit struct {
	Done    bool           `json:"done"`              // Always true
	Updates map[string]any `json:"updates,omitempty"`  // Fields to patch in the row
	Refresh string         `json:"refresh,omitempty"`  // "table" or "page"
	Goto    string         `json:"goto,omitempty"`     // Redirect URL
	Message
}

func (r ReturnInlineEdit) isSuccessResponse() {}

// WithUpdates sets the fields to patch in the row.
func (r ReturnInlineEdit) WithUpdates(updates map[string]any) ReturnInlineEdit {
	r.Updates = updates
	return r
}

// WithRefreshTable triggers a table reload after the update.
func (r ReturnInlineEdit) WithRefreshTable() ReturnInlineEdit {
	r.Refresh = "table"
	return r
}

// WithRefreshPage triggers a full page reload after the update.
func (r ReturnInlineEdit) WithRefreshPage() ReturnInlineEdit {
	r.Refresh = "page"
	return r
}

// WithGoto triggers a navigation after the update.
func (r ReturnInlineEdit) WithGoto(url string) ReturnInlineEdit {
	r.Goto = url
	return r
}

// WithMessage returns a copy with the given message and type.
func (r ReturnInlineEdit) WithMessage(text string, msgType MessageType) ReturnInlineEdit {
	r.MessageText = text
	r.MessageType = msgType
	return r
}

// NewReturnInlineEdit creates a new inline edit response.
//
// Example:
//
//	response.NewReturnInlineEdit().
//	    WithUpdates(map[string]any{"price": []any{"1.299,00", 1299}, "lastModified": "11.03.2026"}).
//	    WithMessage("Saved", response.MessageSuccess)
func NewReturnInlineEdit() ReturnInlineEdit {
	return ReturnInlineEdit{Done: true}
}

// ReturnGoto represents a redirect/goto response.
//
// JSON output: {"done": true, "goto": "/Portal/User/Page/7"}
// With message: {"done": true, "goto": "/Portal/User/Page/7", "message": "Saved", "messageType": "success"}
//
// Use case: Operation completed, navigate to different URL
type ReturnGoto struct {
	Done bool   `json:"done"` // Always true
	Goto string `json:"goto"` // Redirect URL
	Message
}

func (r ReturnGoto) isSuccessResponse() {}

// WithMessage returns a copy with the given message and type.
func (r ReturnGoto) WithMessage(text string, msgType MessageType) ReturnGoto {
	r.MessageText = text
	r.MessageType = msgType
	return r
}

// ReturnDone represents a simple done response.
//
// JSON output: {"done": true}
// With message: {"done": true, "message": "Done", "messageType": "success"}
//
// Use case: Operation completed, no further action needed
type ReturnDone struct {
	Done   bool         `json:"done"` // Always true
	Button *ButtonPatch `json:"button,omitempty"`
	Message
}

func (r ReturnDone) isSuccessResponse() {}

// WithMessage returns a copy with the given message and type.
func (r ReturnDone) WithMessage(text string, msgType MessageType) ReturnDone {
	r.MessageText = text
	r.MessageType = msgType
	return r
}

// WithButton attaches button overrides applied after the action (text, color, disabled, …).
func (r ReturnDone) WithButton(button ButtonPatch) ReturnDone {
	r.Button = &button
	return r
}

// ReturnMessage represents a standalone message response (snackbar only, no navigation).
//
// JSON output: {"done": true, "message": "Settings saved", "messageType": "success"}
//
// Use case: Show a message to the user without any page action
type ReturnMessage struct {
	Done   bool         `json:"done"` // Always true
	Button *ButtonPatch `json:"button,omitempty"`
	Message
}

func (r ReturnMessage) isSuccessResponse() {}

// WithButton attaches button overrides applied after the action (text, color, disabled, …).
func (r ReturnMessage) WithButton(button ButtonPatch) ReturnMessage {
	r.Button = &button
	return r
}

// ReturnPoll tells a button (or other initiator) to keep polling a status endpoint
// at the given interval until a response without "poll" is returned.
//
// JSON output: {"done": true, "poll": 2000, "pollUrl": "Test/Worker/Status"}
// With message: {"done": true, "poll": 2000, "pollUrl": "...", "message": "Started", "messageType": "info"}
//
// Use case: A button triggers a background worker (via api action or a dialog). While the
// worker runs, return ReturnPoll so the frontend re-fetches PollUrl (GET) every Poll ms and
// shows a countdown in the button. When the worker finishes, return a normal response without
// poll (e.g. ReturnRefreshTable/ReturnDone) to stop the polling.
type ReturnPoll struct {
	Done    bool   `json:"done"`              // Always true
	Poll    int    `json:"poll"`              // Poll interval in milliseconds
	PollUrl string       `json:"pollUrl,omitempty"` // Status endpoint polled via GET (optional)
	Text    string       `json:"text,omitempty"`    // Optional label shown inside the button while polling (e.g. "läuft… 50 %")
	Button  *ButtonPatch `json:"button,omitempty"`  // Optional overrides applied to the initiating button
	Message
}

func (r ReturnPoll) isSuccessResponse() {}

// WithButton attaches button overrides applied while/after polling (text, color, disabled, …).
func (r ReturnPoll) WithButton(button ButtonPatch) ReturnPoll {
	r.Button = &button
	return r
}

// WithMessage returns a copy with the given message and type.
func (r ReturnPoll) WithMessage(text string, msgType MessageType) ReturnPoll {
	r.MessageText = text
	r.MessageType = msgType
	return r
}

// WithText returns a copy that shows the given label inside the polling button
// (replacing the default countdown). Can be updated on every poll tick to reflect
// real progress, e.g. "läuft… 50 %".
func (r ReturnPoll) WithText(text string) ReturnPoll {
	r.Text = text
	return r
}

// Constructor functions

// NewReturnRefreshPage creates a refresh page response.
//
// Returns: {"done": true, "refresh": "page"}
func NewReturnRefreshPage() ReturnRefreshPage {
	return ReturnRefreshPage{Done: true, Refresh: "page"}
}

// NewReturnFields creates a field patch response for a form reload.
//
// Returns: {"fields": {"<fieldID>": {...}}}
func NewReturnFields(fields map[string]interface{}) ReturnFields {
	return ReturnFields{Fields: fields}
}

// NewReturnRefreshTable creates a refresh table response.
//
// Returns: {"done": true, "refresh": "table"}
func NewReturnRefreshTable() ReturnRefreshTable {
	return ReturnRefreshTable{Done: true, Refresh: "table"}
}

// NewReturnPoll creates a poll response that asks the frontend to keep polling pollUrl
// every intervalMs milliseconds until a response without "poll" is returned.
//
// Parameters:
//   - pollUrl: the status endpoint to GET each tick (e.g. u.PrintPrefix()); may be empty
//     to let the initiator fall back to its own url
//   - intervalMs: poll interval in milliseconds
//
// Returns: {"done": true, "poll": intervalMs, "pollUrl": pollUrl}
func NewReturnPoll(pollUrl string, intervalMs int) ReturnPoll {
	return ReturnPoll{Done: true, Poll: intervalMs, PollUrl: pollUrl}
}

// NewReturnGoto creates a redirect/goto response.
//
// Parameters:
//   - url: The URL to navigate to (e.g., "/Portal/User/Page/7")
//
// Returns: {"done": true, "goto": url}
func NewReturnGoto(url string) ReturnGoto {
	return ReturnGoto{Done: true, Goto: url}
}

// NewReturnDone creates a simple done response.
//
// Returns: {"done": true}
func NewReturnDone() ReturnDone {
	return ReturnDone{Done: true}
}

// NewReturnMessage creates a standalone message response.
//
// Returns: {"done": true, "message": text, "messageType": msgType}
func NewReturnMessage(text string, msgType MessageType) ReturnMessage {
	return ReturnMessage{Done: true, Message: Message{MessageText: text, MessageType: msgType}}
}

// NewReturnSuccess creates a standalone success message response.
//
// Returns: {"done": true, "message": text, "messageType": "success"}
func NewReturnSuccess(text string) ReturnMessage {
	return NewReturnMessage(text, MessageSuccess)
}

// NewReturnError creates a standalone error message response.
//
// Returns: {"done": true, "message": text, "messageType": "error"}
func NewReturnError(text string) ReturnMessage {
	return NewReturnMessage(text, MessageError)
}

// ErrorResponse represents an HTTP error response body.
// Used by framework-level error helpers (BadRequest, NotFound, etc.)
// that set the appropriate HTTP status code.
//
// JSON output: {"error": "message"}
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewErrorResponse creates an error response with the given message.
//
// Returns: ErrorResponse{"error": message}
func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{Error: message}
}

// NewDataResponse wraps any value in the standard Ajax data envelope.
// The frontend expects {"data": ...} for component data endpoints
// (e.g. Card, Stat, StatGrid, List, MultiProgress).
//
// Returns: map[string]any{"data": data}
func NewDataResponse(data any) map[string]any {
	return map[string]any{"data": data}
}

// ResponseType indicates the format of the response body.
type ResponseType int

const (
	ResponseJSON  ResponseType = iota // Body is map[string]any
	ResponseCSV                       // Body is string
	ResponseExcel                     // Body is []byte
)

// DataResult represents a formatted response with type metadata.
// Components return this to indicate both the data and its format,
// without knowing about HTTP headers or framework specifics.
type DataResult struct {
	Type ResponseType
	Body any
}

// NewJSONDataResult wraps data in the standard {"data": ...} envelope.
func NewJSONDataResult(data any) DataResult {
	return DataResult{Type: ResponseJSON, Body: map[string]any{"data": data}}
}

// NewCSVDataResult creates a CSV response.
func NewCSVDataResult(csv string) DataResult {
	return DataResult{Type: ResponseCSV, Body: csv}
}

// NewExcelDataResult creates an Excel response.
func NewExcelDataResult(excel []byte) DataResult {
	return DataResult{Type: ResponseExcel, Body: excel}
}
