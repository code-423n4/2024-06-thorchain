package main

// trunk-ignore-all(golangci-lint/govet): skip shadowing noise on "ok" for ast inspect

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var reVersionedName = regexp.MustCompile(`.*V([0-9]+)$`)

// -------------------------------------------------------------------------------------
// MapIteration
// -------------------------------------------------------------------------------------

func MapIteration(pass *analysis.Pass) (interface{}, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errors.New("analyzer is not type *inspector.Inspector")
	}

	// track lines with the expected ignore comment
	type ignorePos struct {
		file string
		line int
	}
	ignore := map[ignorePos]bool{}

	// one pass to find all comments
	inspect.Preorder([]ast.Node{(*ast.File)(nil)}, func(node ast.Node) {
		n, ok := node.(*ast.File)
		if !ok {
			panic("node was not *ast.File")
		}
		for _, c := range n.Comments {
			if strings.Contains(c.Text(), "analyze-ignore(map-iteration)") {
				p := pass.Fset.Position(c.Pos())
				ignore[ignorePos{p.Filename, p.Line + strings.Count(c.Text(), "\n")}] = true
			}
		}
	})

	inspect.Preorder([]ast.Node{(*ast.RangeStmt)(nil)}, func(node ast.Node) {
		n, ok := node.(*ast.RangeStmt)
		if !ok {
			panic("node was not *ast.RangeStmt")
		}
		// skip if this is not a range over a map
		if !strings.HasPrefix(pass.TypesInfo.TypeOf(n.X).String(), "map") {
			return
		}

		// skip if this is a test file
		p := pass.Fset.Position(n.Pos())
		if strings.HasSuffix(p.Filename, "_test.go") {
			return
		}

		// skip if the previous line contained the ignore comment
		if ignore[ignorePos{p.Filename, p.Line}] {
			return
		}

		pass.Reportf(node.Pos(), "found map iteration")
	})

	return nil, nil
}

// -------------------------------------------------------------------------------------
// VersionSwitch
// -------------------------------------------------------------------------------------

func VersionSwitch(pass *analysis.Pass) (interface{}, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errors.New("analyzer is not type *inspector.Inspector")
	}

	// find all switch cases
	inspect.Preorder([]ast.Node{(*ast.CaseClause)(nil)}, func(node ast.Node) {
		n, ok := node.(*ast.CaseClause)
		if !ok {
			panic("node was not *ast.CaseClause")
		}

		for _, e := range n.List {
			version := ""
			cmpFn := ""
			var parent ast.Node

			ast.Inspect(e, func(n ast.Node) bool {
				if n == nil || version != "" {
					return false
				}
				if c, ok := n.(*ast.CallExpr); ok {

					// extract the version from semver.MustParse argument
					if s, ok := c.Fun.(*ast.SelectorExpr); ok {
						if x, ok := s.X.(*ast.Ident); ok {
							if x.Name == "semver" && s.Sel.Name == "MustParse" {
								if l, ok := c.Args[0].(*ast.BasicLit); ok {
									version = l.Value
									return false
								}
							}
						}
					}

					// extract the comparison function name
					fnType := pass.TypesInfo.TypeOf(c.Fun)
					if fnType.String() == "func(o github.com/blang/semver.Version) bool" {
						switch ft := c.Fun.(type) {
						case *ast.Ident:
							cmpFn = ft.Name
						case *ast.SelectorExpr:
							cmpFn = ft.Sel.Name
						}
						parent = n
					}
				}

				return true
			})
			if version == "" {
				continue
			}

			// ensure version switch is using GTE
			if cmpFn != "GTE" {
				pass.Reportf(parent.Pos(), "must use GTE in version switch")
			}

			// extract the minor version
			minor := strings.Split(version, ".")[1]

			// extract versioned functions called in the case body
			for _, s := range n.Body {
				ast.Inspect(s, func(n ast.Node) bool {
					if n == nil {
						return false
					}
					if c, ok := n.(*ast.CallExpr); ok {
						// extract function names from the body
						vFn := ""
						switch ft := c.Fun.(type) {
						case *ast.Ident:
							vFn = ft.Name
						case *ast.SelectorExpr:
							vFn = ft.Sel.Name
						default:
							return true
						}

						// verify function versions match
						v := reVersionedName.FindStringSubmatch(vFn)
						if len(v) == 2 && v[1] != minor {
							pass.Reportf(e.Pos(), "bad version switch body: %s != %s", v[1], minor)
						}
					}
					return true
				})
			}
		}
	})

	return nil, nil
}

// -------------------------------------------------------------------------------------
// GetMimirCheck
// -------------------------------------------------------------------------------------

func GetMimirCheck(pass *analysis.Pass) (interface{}, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errors.New("analyzer is not type *inspector.Inspector")
	}

	var (
		switchStmt *ast.SwitchStmt
		constIds   []string
	)

	// getmimir func and const ids for mimirv2
	inspect.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(node ast.Node) {
		n, ok := node.(*ast.FuncDecl)
		if !ok {
			return
		}
		if n.Name.Name != "GetMimir" {
			return
		}
		for _, stmt := range n.Body.List {
			if s, ok := stmt.(*ast.SwitchStmt); ok {
				switchStmt = s
				break
			}
		}
	})

	if switchStmt == nil {
		return nil, nil
	}

	inspect.Preorder([]ast.Node{(*ast.GenDecl)(nil)}, func(node ast.Node) {
		n, ok := node.(*ast.GenDecl)
		if !ok {
			return
		}
		if n.Tok != token.CONST {
			return
		}
		if len(n.Specs) == 0 {
			return
		}
		firstSpec, ok := n.Specs[0].(*ast.ValueSpec)
		if !ok || len(firstSpec.Names) == 0 || firstSpec.Names[0].Name != "Unknown" {
			return
		}
		for _, spec := range n.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				for _, name := range valueSpec.Names {
					constIds = append(constIds, name.Name)
				}
			}
		}
	})

	caseClauses := make(map[string]struct{})
	for _, stmt := range switchStmt.Body.List {
		caseClause, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}
		for _, clause := range caseClause.List {
			caseClauses[fmt.Sprintf("%s", clause)] = struct{}{}
		}
	}

	for _, id := range constIds {
		if id == "Unknown" {
			continue
		}
		if _, found := caseClauses[id]; !found {
			pass.Reportf(switchStmt.Pos(), "get mimir case not found for : %s", id)
		}
	}
	return nil, nil
}

// -------------------------------------------------------------------------------------
// Rand
// -------------------------------------------------------------------------------------

func Rand(pass *analysis.Pass) (interface{}, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errors.New("analyzer is not type *inspector.Inspector")
	}

	allow := func(pos token.Pos) bool {
		p := pass.Fset.Position(pos)
		if strings.HasSuffix(p.Filename, "_test.go") {
			return true
		}
		if strings.HasSuffix(p.Filename, "querier_quotes.go") {
			return true
		}
		if strings.HasSuffix(p.Filename, "test_common.go") {
			return true
		}
		return false
	}

	inspect.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(node ast.Node) {
		n, ok := node.(*ast.CallExpr)
		if !ok {
			return
		}
		if id, ok := n.Fun.(*ast.Ident); ok {
			if strings.Contains(id.Name, "Rand") && !allow(n.Pos()) {
				pass.Reportf(n.Pos(), "use of functions with \"Rand\" in name is prohibited")
			}
		}
	})

	return nil, nil
}

// -------------------------------------------------------------------------------------
// Main
// -------------------------------------------------------------------------------------

func main() {
	multichecker.Main(
		&analysis.Analyzer{
			Name:     "map_iteration",
			Doc:      "fails on uncommented map iterations",
			Requires: []*analysis.Analyzer{inspect.Analyzer},
			Run:      MapIteration,
		},
		&analysis.Analyzer{
			Name:     "switch_version",
			Doc:      "fails on bad version switches",
			Requires: []*analysis.Analyzer{inspect.Analyzer},
			Run:      VersionSwitch,
		},
		&analysis.Analyzer{
			Name:     "mimir_v2_getmimir",
			Doc:      "fails if no get mimir case is defined for mimirv2 id",
			Requires: []*analysis.Analyzer{inspect.Analyzer},
			Run:      GetMimirCheck,
		},
		&analysis.Analyzer{
			Name:     "rand",
			Doc:      "fails on use of functions with \"Rand\" in name",
			Requires: []*analysis.Analyzer{inspect.Analyzer},
			Run:      Rand,
		},
	)
}
