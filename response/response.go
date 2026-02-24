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
	Done bool `json:"done"` // Always true
	Message
}

func (r ReturnDone) isSuccessResponse() {}

// WithMessage returns a copy with the given message and type.
func (r ReturnDone) WithMessage(text string, msgType MessageType) ReturnDone {
	r.MessageText = text
	r.MessageType = msgType
	return r
}

// ReturnMessage represents a standalone message response (snackbar only, no navigation).
//
// JSON output: {"done": true, "message": "Settings saved", "messageType": "success"}
//
// Use case: Show a message to the user without any page action
type ReturnMessage struct {
	Done bool `json:"done"` // Always true
	Message
}

func (r ReturnMessage) isSuccessResponse() {}

// Constructor functions

// NewReturnRefreshPage creates a refresh page response.
//
// Returns: {"done": true, "refresh": "page"}
func NewReturnRefreshPage() ReturnRefreshPage {
	return ReturnRefreshPage{Done: true, Refresh: "page"}
}

// NewReturnRefreshTable creates a refresh table response.
//
// Returns: {"done": true, "refresh": "table"}
func NewReturnRefreshTable() ReturnRefreshTable {
	return ReturnRefreshTable{Done: true, Refresh: "table"}
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
