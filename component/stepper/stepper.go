// Package stepper provides multi-step wizard components for the Angular frontend.
package stepper

import (
	"fmt"

	"github.com/xiriframework/xiri-go/component/button"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
)

// StepFields represents the fields for a specific step
type StepFields interface {
	GetFormFields() []map[string]any
}

// Stepper represents a multi-step wizard component
type Stepper struct {
	url     *url.Url
	steps   int
	titles  []string
	fields  [][]map[string]any // Fields for each step
	back    string
	next    string
	done    string
	display *string
}

// NewStepper creates a new stepper component
func NewStepper(
	u *url.Url,
	steps int,
	titles []string,
	fields [][]map[string]any,
	textBack string,
	textNext string,
	textDone string,
	display *string,
) (*Stepper, error) {
	if len(titles) != steps {
		return nil, fmt.Errorf("XiriStepper: wrong step count != titles (got %d steps, %d titles)", steps, len(titles))
	}
	if len(fields) != steps {
		return nil, fmt.Errorf("XiriStepper: wrong step count != fields (got %d steps, %d field arrays)", steps, len(fields))
	}

	return &Stepper{
		url:     u,
		steps:   steps,
		titles:  titles,
		fields:  fields,
		back:    textBack,
		next:    textNext,
		done:    textDone,
		display: display,
	}, nil
}

// WithDisplay sets the display/layout class (optional)
func (s *Stepper) WithDisplay(display string) *Stepper {
	s.display = &display
	return s
}

// Print returns the JSON representation of the stepper
func (s *Stepper) Print(translator core.TranslateFunc) map[string]any {
	stepArray := make([]map[string]any, s.steps)

	for i := 1; i <= s.steps; i++ {
		firstAction := core.ButtonActionPrev
		if i == 1 {
			firstAction = core.ButtonActionBack
		}

		secondText := s.next
		secondAction := core.ButtonActionNext
		if i == s.steps {
			secondText = s.done
		}

		backBtn := button.NewButton(
			firstAction,
			s.back,
			url.NewUrl(""),
			core.ColorPrimary,
			core.ButtonTypeStroked,
			"",
			"",
			false,
			nil,
			false,
			"_self",
			nil,
		)

		nextBtn := button.NewButton(
			secondAction,
			secondText,
			s.url,
			core.ColorPrimary,
			core.ButtonTypeRaised,
			"",
			"",
			false,
			nil,
			true,
			"_self",
			nil,
		)

		stepArray[i-1] = map[string]any{
			"title": s.titles[i-1],
			"extra": map[string]int{
				"step": i,
			},
			"fields": s.fields[i-1],
			"buttons": []map[string]any{
				backBtn.Print(translator),
				nextBtn.Print(translator),
			},
		}
	}

	return map[string]any{
		"type":    "stepper",
		"display": s.display,
		"data": map[string]any{
			"url":   s.url.PrintPrefix(),
			"steps": stepArray,
		},
	}
}

// StepperStep represents a stepper step handler (server-side processing)
type StepperStep struct {
	step   int
	fields []map[string]any
}

// NewStepperStep creates a new stepper step handler
func NewStepperStep(
	step int,
	fields []map[string]any,
) *StepperStep {
	return &StepperStep{
		step:   step,
		fields: fields,
	}
}

// Print returns the next step's field configuration
func (ss *StepperStep) Print(translator core.TranslateFunc) map[string]any {
	return map[string]any{
		"fields": ss.fields,
		"step":   ss.step,
	}
}
