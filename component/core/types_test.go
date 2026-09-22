package core

import (
	"bytes"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

var (
	typeLiteral     = regexp.MustCompile("\"type\":\\s*[\"`]([a-z0-9_-]+)[\"`]")
	envelopeLiteral = regexp.MustCompile("Envelope\\(\\s*[\"`]([a-z0-9_-]+)[\"`]")
)

// scanComponentTypes sammelt alle statischen "type"-Literale unter component/ (ohne Tests).
// Die Quelle wird über go/parser ohne Kommentare neu formatiert, damit Kommentare nicht zählen,
// Strings mit "//" aber intakt bleiben.
// ponytail: Regex auf kommentarfreier Quelle; AST-Walk erst, wenn ein Builder den Typ nicht mehr als Literal schreibt.
func scanComponentTypes(t *testing.T) map[string][]string {
	t.Helper()
	found := map[string][]string{}
	err := filepath.WalkDir("..", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		var src bytes.Buffer
		if err := format.Node(&src, fset, file); err != nil {
			return err
		}
		for _, re := range []*regexp.Regexp{typeLiteral, envelopeLiteral} {
			for _, m := range re.FindAllStringSubmatch(src.String(), -1) {
				found[m[1]] = append(found[m[1]], path)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return found
}

func TestComponentTypes_CoverAllLiterals(t *testing.T) {
	types := ComponentTypes()
	if len(types) == 0 {
		t.Fatal("core.ComponentTypes() is empty")
	}
	if !slices.IsSorted(types) {
		t.Error("core.ComponentTypes() must be sorted")
	}
	if len(slices.Compact(slices.Clone(types))) != len(types) {
		t.Error("core.ComponentTypes() contains duplicates")
	}
	found := scanComponentTypes(t)
	for typ, files := range found {
		if !slices.Contains(types, typ) {
			t.Errorf("type %q is printed in %v but missing in core.ComponentTypes()", typ, files)
		}
	}
	for _, typ := range types {
		if _, ok := found[typ]; !ok {
			t.Errorf("core.ComponentTypes() lists %q but no component prints it", typ)
		}
	}
}

func TestComponentTypes_ReturnsCopy(t *testing.T) {
	a := ComponentTypes()
	a[0] = "corrupted"
	if b := ComponentTypes(); b[0] == "corrupted" {
		t.Error("ComponentTypes() must return a fresh copy")
	}
}

// TestComponentTypes_MatchXiriNg vergleicht mit dem xiri-ng-Katalog, wenn das Schwester-Repo im
// Workspace liegt (../../../xiri-ng). Nur ein fehlendes Schwester-Repo führt zum Skip, damit
// go test in CI ohne Workspace grün bleibt; ein fehlender oder unlesbarer Katalog ist ein Fehler.
func TestComponentTypes_MatchXiriNg(t *testing.T) {
	ngRoot := filepath.Join("..", "..", "..", "xiri-ng")
	if _, err := os.Stat(ngRoot); err != nil {
		if os.IsNotExist(err) {
			t.Skipf("xiri-ng not found at %s", ngRoot)
		}
		t.Fatal(err)
	}
	src, err := os.ReadFile(filepath.Join(ngRoot, "projects", "xiri-ng", "src", "lib", "dyncomponent", "component-catalog.ts"))
	if err != nil {
		t.Fatal(err)
	}
	// Eintragsweise: jede Objektzeile im Katalog-Array muss type und goBuilder tragen.
	// ponytail: ein Eintrag pro Zeile ist Katalog-Konvention; fremde Schreibweisen scheitern laut statt still.
	typeRe := regexp.MustCompile(`type:\s*'([a-z0-9_-]+)'`)
	goRe := regexp.MustCompile(`goBuilder:\s*(true|false)`)
	commentRe := regexp.MustCompile(`/\*.*?\*/|//.*$`)
	goBuilder := map[string]bool{}
	inArray := false
	for i, line := range strings.Split(string(src), "\n") {
		trimmed := strings.TrimSpace(commentRe.ReplaceAllString(line, ""))
		switch {
		case strings.HasPrefix(trimmed, "export const COMPONENT_CATALOG"):
			inArray = true
			continue
		case inArray && strings.HasPrefix(trimmed, "];"):
			inArray = false
		}
		if !inArray || !strings.HasPrefix(trimmed, "{") {
			continue
		}
		tms, gms := typeRe.FindAllStringSubmatch(trimmed, -1), goRe.FindAllStringSubmatch(trimmed, -1)
		if len(tms) != 1 || len(gms) != 1 {
			t.Fatalf("component-catalog.ts:%d: entry must have exactly one type and one goBuilder: %s", i+1, trimmed)
		}
		typ := tms[0][1]
		if _, dup := goBuilder[typ]; dup {
			t.Fatalf("component-catalog.ts:%d: duplicate type %q", i+1, typ)
		}
		goBuilder[typ] = gms[0][1] == "true"
	}
	if len(goBuilder) == 0 {
		t.Fatal("no catalog entries parsed; catalog format changed")
	}
	types := ComponentTypes()
	for typ, isGo := range goBuilder {
		if isGo && !slices.Contains(types, typ) {
			t.Errorf("xiri-ng catalog claims goBuilder for %q, but core.ComponentTypes() lacks it", typ)
		}
	}
	rendererless := []string{"tachotime"} // Go-Builder ohne @case: Host-Apps rendern per Template.
	for _, typ := range types {
		if slices.Contains(rendererless, typ) {
			continue
		}
		if isGo, ok := goBuilder[typ]; !ok || !isGo {
			t.Errorf("core.ComponentTypes() has %q, but xiri-ng catalog has no goBuilder:true entry", typ)
		}
	}
}
