package field

import "testing"

func validTimeLimit() TimeLimitValue {
	return TimeLimitValue{Check: true, FromHour: "07", FromMin: "30", ToHour: "24", ToMin: "00", In: true}
}

func TestTimeLimitField_Validate_RejectsNonDigits(t *testing.T) {
	f := NewTimeLimitField("tl", "TL", false)
	mutations := map[string]func(*TimeLimitValue){
		"garbage":        func(v *TimeLimitValue) { v.FromHour = "x" },
		"trailing":       func(v *TimeLimitValue) { v.FromMin = "30'); DROP" },
		"sign":           func(v *TimeLimitValue) { v.ToHour = "+5" },
		"too long":       func(v *TimeLimitValue) { v.ToMin = "000" },
		"unchecked junk": func(v *TimeLimitValue) { v.Check = false; v.FromHour = "x" },
	}
	for name, mut := range mutations {
		v := validTimeLimit()
		mut(&v)
		if err := f.Validate(v); err == nil {
			t.Errorf("%s: expected error for %+v", name, v)
		}
	}
}

func TestTimeLimitField_Validate_AcceptsValidAndDefault(t *testing.T) {
	f := NewTimeLimitField("tl", "TL", false)
	if err := f.Validate(validTimeLimit()); err != nil {
		t.Errorf("valid value rejected: %v", err)
	}
	if err := f.Validate(f.GetDefault()); err != nil {
		t.Errorf("default rejected: %v", err)
	}
	// Zahlen vom Frontend: float64(7) wird per %v zu "7"
	v, err := f.Parse(map[string]interface{}{"check": true, "fromhour": float64(7), "frommin": float64(5),
		"tohour": float64(18), "tomin": float64(0)})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := f.Validate(v); err != nil {
		t.Errorf("numeric frontend value rejected: %v", err)
	}
}
