// Package stat provides a KPI/statistics card component for the Angular frontend.
package stat

import (
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// TrendDirection indicates the direction of a trend indicator.
type TrendDirection string

const (
	TrendUp      TrendDirection = "up"
	TrendDown    TrendDirection = "down"
	TrendNeutral TrendDirection = "neutral"
)

// Trend represents a trend indicator with value and direction.
type Trend struct {
	Value     float64        `json:"value"`
	Direction TrendDirection `json:"direction"`
}

// Stat represents a KPI/statistics card component.
type Stat struct {
	value     interface{} // string or number
	label     string
	icon      *string
	iconColor *string
	trend     *Trend
	prefix    *string
	suffix    *string
	color     *string
	display   *string
	url       *url.Url
	reload    *bool
}

// New creates a new Stat component with value and label.
func New(value interface{}, label string) *Stat {
	return &Stat{
		value: value,
		label: label,
	}
}

// Icon sets the icon for the stat card.
func (s *Stat) Icon(icon string) *Stat {
	s.icon = &icon
	return s
}

// IconColor sets the icon color.
func (s *Stat) IconColor(color string) *Stat {
	s.iconColor = &color
	return s
}

// SetTrend sets the trend indicator.
func (s *Stat) SetTrend(value float64, direction TrendDirection) *Stat {
	s.trend = &Trend{Value: value, Direction: direction}
	return s
}

// Prefix sets the value prefix (e.g. currency symbol).
func (s *Stat) Prefix(prefix string) *Stat {
	s.prefix = &prefix
	return s
}

// Suffix sets the value suffix (e.g. unit).
func (s *Stat) Suffix(suffix string) *Stat {
	s.suffix = &suffix
	return s
}

// Color sets the value color.
func (s *Stat) Color(color string) *Stat {
	s.color = &color
	return s
}

// WithDisplay sets the display/layout class.
func (s *Stat) WithDisplay(display string) *Stat {
	s.display = &display
	return s
}

// SetURL sets the AJAX data URL. When set, the frontend loads stat data dynamically.
func (s *Stat) SetURL(url *url.Url) *Stat {
	s.url = url
	return s
}

// WithReload enables periodic reload of the stat data when using AJAX mode.
func (s *Stat) WithReload(reload bool) *Stat {
	s.reload = &reload
	return s
}

// Print returns the JSON representation of the stat component.
func (s *Stat) Print(translator core.TranslateFunc) map[string]any {
	var data map[string]any

	if s.url != nil {
		data = map[string]any{
			"url": s.url.PrintPrefix(),
		}
		if s.reload != nil {
			data["reload"] = *s.reload
		}
	} else {
		data = s.printData(translator)
	}

	result := map[string]any{
		"type": "stat",
		"data": data,
	}

	if s.display != nil {
		result["display"] = *s.display
	}

	return result
}

// PrintData returns only the data portion of the stat (for use in data endpoints and StatGrid).
func (s *Stat) PrintData(translator core.TranslateFunc) map[string]any {
	return s.printData(translator)
}

// DataResponse returns a DataResult wrapping the stat data in {"data": ...} envelope.
func (s *Stat) DataResponse(translator core.TranslateFunc) response.DataResult {
	return response.NewJSONDataResult(s.PrintData(translator))
}

// printData builds the data map used by both Print and PrintData.
func (s *Stat) printData(translator core.TranslateFunc) map[string]any {
	data := map[string]any{
		"value": s.value,
		"label": core.Translate(translator, s.label),
	}

	if s.icon != nil {
		data["icon"] = *s.icon
	}
	if s.iconColor != nil {
		data["iconColor"] = *s.iconColor
	}
	if s.trend != nil {
		data["trend"] = map[string]any{
			"value":     s.trend.Value,
			"direction": string(s.trend.Direction),
		}
	}
	if s.prefix != nil {
		data["prefix"] = *s.prefix
	}
	if s.suffix != nil {
		data["suffix"] = *s.suffix
	}
	if s.color != nil {
		data["color"] = *s.color
	}

	return data
}
