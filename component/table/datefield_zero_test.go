package table_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/xiriframework/xiri-go/component/table"
)

// #12: a zero time.Time must render as empty, not as year 0001.
func TestDateFieldZeroTimeEmpty(t *testing.T) {
	type row struct{ When time.Time }
	ctx := exampleContext()

	builder := table.NewBuilder[row]()
	builder.DateTimeField("dt", "dt", func(r row) time.Time { return r.When })
	builder.DateField("d", "d", func(r row) time.Time { return r.When })

	tbl := builder.Build()
	tbl.SetData([]row{{When: time.Time{}}}) // zero time
	data := tbl.GetData(ctx, table.OutputWeb)

	want := map[string]any{"d": "", "v": nil}
	if got := data[0]["dt"]; !reflect.DeepEqual(got, want) {
		t.Errorf("DateTimeField zero time = %#v, want %#v", got, want)
	}
	if got := data[0]["d"]; !reflect.DeepEqual(got, want) {
		t.Errorf("DateField zero time = %#v, want %#v", got, want)
	}
}
