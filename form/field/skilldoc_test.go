package field

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Die ausgelieferte Skill-Doku (skills/xiri-go-expert/) ist die Vorlage, nach der Consumer Code
// schreiben — eine falsche Signatur dort kostet sie einen Compile-Fehler. Dieser Test prueft die
// Doku-Snippets gegen die echten Konstruktoren dieses Packages, damit Signatur-Aenderungen die
// Doku nicht still veralten lassen.
//
// ponytail: nur zwei Regeln (Konstruktor existiert, kein nil/& an einen Wert-Parameter) statt
// vollem Typecheck der Snippets. Echtes Kompilieren wuerde verlangen, jedes Snippet zu einer
// vollstaendigen Datei mit Imports auszubauen — dann lieber go/packages auf generierte Dateien.

var docCallRe = regexp.MustCompile(`field\.(New\w+)\s*\(`)

func TestSkillDocConstructorSignatures(t *testing.T) {
	skillsDir := filepath.Join("..", "..", "skills")
	if _, err := os.Stat(skillsDir); err != nil {
		t.Skip("keine skills/ neben dem Package — nichts zu pruefen")
	}

	lastParams := packageConstructors(t)

	err := filepath.WalkDir(skillsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		doc := string(src)

		for _, m := range docCallRe.FindAllStringSubmatchIndex(doc, -1) {
			name := doc[m[2]:m[3]]
			line := 1 + strings.Count(doc[:m[0]], "\n")

			last, ok := lastParams[name]
			if !ok {
				t.Errorf("%s:%d: field.%s existiert nicht im Package", path, line, name)
				continue
			}
			if last == nil {
				continue // parameterloser Konstruktor
			}

			args, ok := splitArgs(doc[m[1]-1:])
			if !ok || len(args) == 0 {
				continue // unvollstaendiges Snippet — kein Signatur-Befund
			}
			arg := args[len(args)-1]
			if arg != "nil" && !strings.HasPrefix(arg, "&") {
				continue
			}
			if nilable(last) {
				continue
			}
			t.Errorf("%s:%d: field.%s(…, %s) — letzter Parameter ist %s, kein Pointer. Nullwert verwenden.",
				path, line, name, arg, types.ExprString(last))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("skills/ durchlaufen: %v", err)
	}
}

// packageConstructors liefert pro exportiertem New*-Konstruktor den Typ des letzten Parameters
// (nil bei parameterlosen Konstruktoren).
func packageConstructors(t *testing.T) map[string]ast.Expr {
	t.Helper()

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("Package-Verzeichnis lesen: %v", err)
	}

	fset := token.NewFileSet()
	out := map[string]ast.Expr{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, e.Name(), nil, 0)
		if err != nil {
			t.Fatalf("%s parsen: %v", e.Name(), err)
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !strings.HasPrefix(fn.Name.Name, "New") {
				continue
			}
			params := fn.Type.Params.List
			if len(params) == 0 {
				out[fn.Name.Name] = nil
				continue
			}
			out[fn.Name.Name] = params[len(params)-1].Type
		}
	}
	if len(out) == 0 {
		t.Fatal("keine New*-Konstruktoren gefunden — Parser-Setup kaputt")
	}
	return out
}

// nilable sagt, ob nil bzw. eine Adresse ein gueltiges Argument fuer diesen Parametertyp ist.
func nilable(typ ast.Expr) bool {
	switch typ.(type) {
	case *ast.StarExpr, *ast.ArrayType, *ast.MapType, *ast.InterfaceType, *ast.ChanType, *ast.FuncType, *ast.Ellipsis:
		return true
	}
	return false
}

// splitArgs zerlegt die Argumentliste eines Calls. src beginnt mit der oeffnenden Klammer.
// ok ist false, wenn die Klammer im Text nicht geschlossen wird (abgeschnittenes Snippet).
func splitArgs(src string) (args []string, ok bool) {
	depth := 0
	start := 1
	var quote rune

	for i, r := range src {
		if quote != 0 {
			if r == quote && src[i-1] != '\\' {
				quote = 0
			}
			continue
		}
		switch r {
		case '"', '\'', '`':
			quote = r
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			depth--
			if depth == 0 {
				return append(args, strings.TrimSpace(src[start:i])), true
			}
		case ',':
			if depth == 1 {
				args = append(args, strings.TrimSpace(src[start:i]))
				start = i + 1
			}
		}
	}
	return nil, false
}
