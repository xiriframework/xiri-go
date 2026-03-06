package dialog

import "github.com/xiriframework/xiri-go/component/core"

// DialogContent is an interface for content that needs custom serialization
type DialogContent interface {
	Print(ctx *core.UiContext) map[string]any
}

// DialogQuestionContent represents content for question-type dialogs (delete, warning)
type DialogQuestionContent struct {
	Icon     string `json:"icon"`
	Question string `json:"question"`
}

// Print returns the JSON representation of question content
func (d DialogQuestionContent) Print(ctx *core.UiContext) map[string]any {
	return map[string]any{
		"icon":     d.Icon,
		"question": d.Question,
	}
}

// DialogWaitingContent represents content for waiting dialogs
type DialogWaitingContent struct {
	Icon string `json:"icon"`
	Text string `json:"text"`
}

// Print returns the JSON representation of waiting content
func (d DialogWaitingContent) Print(ctx *core.UiContext) map[string]any {
	return map[string]any{
		"icon": d.Icon,
		"text": d.Text,
	}
}
