package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	reVersionedName  = regexp.MustCompile(`[Vv]([0-9]+)`)
	skipRootPackages = map[string]bool{
		"bifrost": true,
		"chain":   true,
		"docs":    true,
		"openapi": true,
		"scripts": true,
		"test":    true,
		"tools":   true,
	}
)

func main() {
	var managersOnly bool
	flag.BoolVar(&managersOnly, "managers", false, "only check for managers")
	flag.Parse()

	fset := token.NewFileSet()
	pkgs := []*ast.Package{}

	// parse all subdirectories with go files
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// trunk-ignore(golangci-lint/govet): shadow
			subPkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}
			for _, pkg := range subPkgs {
				pkgs = append(pkgs, pkg)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalln("Error parsing files:", err)
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			fn := fset.File(file.Pos()).Name()

			// skip packages with no versioned functions relevant for consensus
			rootPackage := strings.Split(fn, "/")[0]
			if skipRootPackages[rootPackage] {
				continue
			}
			// skip non-test files
			if !strings.HasSuffix(fn, "_test.go") {
				continue
			}
			// skip archive tests
			if strings.HasSuffix(fn, "_archive_test.go") {
				continue
			}
			// ignore versioned test files
			fileVers := reVersionedName.FindString(fn)
			if fileVers != "" {
				continue
			}

			ast.Inspect(file, func(n ast.Node) bool {
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

					// manager constructors start with "new"
					if managersOnly && !strings.HasPrefix(vFn, "new") {
						return true
					}

					funcVers := reVersionedName.FindString(vFn)
					if funcVers != "" {
						position := fset.Position(n.Pos())
						fmt.Printf("%s called on %s:%d\n", vFn, fn, position.Line)
					}
				}
				return true
			})
		}
	}
}
