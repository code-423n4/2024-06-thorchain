package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

// -------------------------------------------------------------------------------------
// Main
// -------------------------------------------------------------------------------------

func main() {
	filename := "mimir/id.go"
	fset := token.NewFileSet()
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	node, err := parser.ParseFile(fset, filename, file, parser.ParseComments)
	if err != nil {
		// trunk-ignore(golangci-lint/gocritic)
		log.Fatal(err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		// look for variable declarations
		genDecl, ok := n.(*ast.GenDecl)
		if !ok {
			return true
		}

		// track the type to detect implicit continuation of the same type
		currentType := ""

		for _, spec := range genDecl.Specs {
			var valueSpec *ast.ValueSpec
			valueSpec, ok = spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			// if the type is not explicitly stated, use the last type detected
			if valueSpec.Type != nil {
				var typeIdent *ast.Ident
				typeIdent, ok = valueSpec.Type.(*ast.Ident)
				if ok {
					currentType = typeIdent.Name
				}
			}

			if currentType == "Id" {
				for _, name := range valueSpec.Names {
					fmt.Println(name.Name)
				}
			}
		}

		return true
	})
}
