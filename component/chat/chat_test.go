package chat

import (
	"testing"
	"time"

	"github.com/xiriframework/xiri-go/component/core"
)

func data(t *testing.T, out map[string]any) map[string]any {
	t.Helper()
	if out["type"] != "chat" {
		t.Fatalf("type = %v, want chat", out["type"])
	}
	d, ok := out["data"].(map[string]any)
	if !ok {
		t.Fatalf("data missing: %v", out)
	}
	return d
}

func TestPrintMinimal(t *testing.T) {
	out := New("Ticket/42/Chat").Print(nil)
	d := data(t, out)
	if len(d) != 1 || d["url"] != "Ticket/42/Chat" {
		t.Errorf("data = %v, want only url", d)
	}
	if _, ok := out["display"]; ok {
		t.Error("display must be absent when unset")
	}
}

func TestPrintAllOptions(t *testing.T) {
	ctx := &core.UiContext{Translate: func(k string) string {
		if k == "chat.placeholder" {
			return "Nachricht schreiben"
		}
		return k
	}}
	out := New("C").WithStream("C/Stream").WithPoll(5).WithPlaceholder("chat.placeholder").
		WithMaxLength(4000).WithHeight("60vh").WithDisplay("xcol-6").ReadOnly().Print(ctx)
	d := data(t, out)
	want := map[string]any{
		"url": "C", "stream": "C/Stream", "poll": 5, "placeholder": "Nachricht schreiben",
		"maxLength": 4000, "height": "60vh", "readonly": true,
	}
	for k, v := range want {
		if d[k] != v {
			t.Errorf("data[%q] = %v, want %v", k, d[k], v)
		}
	}
	if out["display"] != "xcol-6" {
		t.Errorf("display = %v", out["display"])
	}
}

func TestPrintEmbedsPage(t *testing.T) {
	p := NewPage("c-7").HasOlder(true)
	p.Add("m1", "Hallo", time.Unix(1729070000, 0)).From("Anna Berger")
	d := data(t, New("C").WithPage(p).Print(nil))
	if d["cursor"] != "c-7" || d["hasOlder"] != true {
		t.Errorf("cursor/hasOlder = %v/%v", d["cursor"], d["hasOlder"])
	}
	if msgs := d["messages"].([]map[string]any); len(msgs) != 1 || msgs[0]["id"] != "m1" {
		t.Errorf("messages = %v", d["messages"])
	}
}

func TestMessageFields(t *testing.T) {
	p := NewPage("c")
	p.Add("m1", "Hi", time.Unix(100, 0)).From("Anna").Avatar("/a.png").Status(Read)
	p.Add("m2", "Yo", time.Unix(200, 0)).Own().From("Ich").Avatar("/me.png").Status(Delivered).ClientID("c-1")
	msgs := p.Print()["messages"].([]map[string]any)

	foreign := msgs[0]
	if foreign["text"] != "Hi" || foreign["time"] != int64(100) || foreign["from"] != "Anna" || foreign["avatar"] != "/a.png" {
		t.Errorf("foreign = %v", foreign)
	}
	if _, ok := foreign["status"]; ok {
		t.Error("status only for own messages")
	}
	if _, ok := foreign["own"]; ok {
		t.Error("own must be absent for foreign messages")
	}

	own := msgs[1]
	if own["own"] != true || own["status"] != "delivered" || own["clientId"] != "c-1" {
		t.Errorf("own = %v", own)
	}
	for _, k := range []string{"from", "avatar"} {
		if _, ok := own[k]; ok {
			t.Errorf("%q must be absent for own messages", k)
		}
	}
}

func TestEmptyPagePrintsEmptyList(t *testing.T) {
	msgs, ok := NewPage("").Print()["messages"].([]map[string]any)
	if !ok || msgs == nil || len(msgs) != 0 {
		t.Errorf("messages = %#v, want empty non-nil slice (JSON [])", msgs)
	}
}

func TestAddPointerStableAfterRealloc(t *testing.T) {
	p := NewPage("")
	first := p.Add("m0", "x", time.Unix(0, 0))
	for i := 1; i < 100; i++ {
		p.Add("mx", "x", time.Unix(0, 0))
	}
	first.Own()
	if p.Print()["messages"].([]map[string]any)[0]["own"] != true {
		t.Error("pointer returned by Add went stale")
	}
}

func TestSendResult(t *testing.T) {
	out := SendResult(NewMessage("m9", "Hi", time.Unix(5, 0)).Own().Status(Sent).ClientID("c-9"))
	m := out["message"].(map[string]any)
	if m["id"] != "m9" || m["status"] != "sent" || m["clientId"] != "c-9" {
		t.Errorf("message = %v", m)
	}
}

func TestPageActivity(t *testing.T) {
	started := time.Unix(1793000000, 0)
	got := NewPage("c").Activity("Agent arbeitet … (Bash)", started).Print()["activity"]
	want := map[string]any{"text": "Agent arbeitet … (Bash)", "since": int64(1793000000)}
	if a, ok := got.(map[string]any); !ok || len(a) != 2 || a["text"] != want["text"] || a["since"] != want["since"] {
		t.Errorf("activity = %v, want %v", got, want)
	}

	got = NewPage("c").Activity("Agent arbeitet …", time.Time{}).Print()["activity"]
	if a, ok := got.(map[string]any); !ok || len(a) != 1 || a["text"] != "Agent arbeitet …" {
		t.Errorf("zero since: activity = %v, want only text", got)
	}

	for name, p := range map[string]*Page{"ohne Aufruf": NewPage("c"), "leerer Text": NewPage("c").Activity("", started)} {
		if a, ok := p.Print()["activity"]; ok {
			t.Errorf("%s: activity = %v, want absent", name, a)
		}
	}

	d := data(t, New("C").WithPage(NewPage("c").Activity("läuft", started)).Print(nil))
	if a, ok := d["activity"].(map[string]any); !ok || a["text"] != "läuft" {
		t.Errorf("WithPage: activity = %v", d["activity"])
	}
}
