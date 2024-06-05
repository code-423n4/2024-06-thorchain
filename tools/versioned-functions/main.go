package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"hash/fnv"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// -------------------------------------------------------------------------------------
// Flags
// -------------------------------------------------------------------------------------

var flagVersion *int

func init() {
	flagVersion = flag.Int("version", 0, "current version allowing changes")
}

// -------------------------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------------------------

var (
	reCurrentVersionName   = regexp.MustCompile(`.*VCUR$`)
	reVersionedName        = regexp.MustCompile(`.*V([0-9]+)$`)
	reNonMainnetBuildFlags = regexp.MustCompile(`([^!](test|mock|stage)net|[^!]regtest)`)
	currentManagerVersions = map[string]string{}
	skipRootPackages       = map[string]bool{
		"bifrost": true,
		"chain":   true,
		"docs":    true,
		"openapi": true,
		"scripts": true,
		"test":    true,
		"tools":   true,
	}
)

func isVersionedFunction(node ast.Node, fset *token.FileSet) (bool, int) {
	n, ok := node.(*ast.FuncDecl)
	var version string

	switch {
	case !ok:
		return false, 0

	case reVersionedName.MatchString(n.Name.Name):
		// extract the version from the function name
		version = reVersionedName.FindStringSubmatch(n.Name.Name)[1]

	case reCurrentVersionName.MatchString(n.Name.Name):
		// if this is a current version get the mapping from managers.go
		manager := n.Name.Name
		if strings.HasPrefix(manager, "new") {
			manager = strings.TrimPrefix(n.Name.Name, "new")
		}
		version, ok = currentManagerVersions[manager]
		if !ok {
			fmt.Println("Error: could not find current version for", n.Name.Name)
			os.Exit(1)
		}

	case !reVersionedName.MatchString(n.Name.Name) && !reCurrentVersionName.MatchString(n.Name.Name):
		// search receiver for a versioned struct name
		if n.Recv != nil {
			for _, r := range n.Recv.List {
				buf := new(bytes.Buffer)
				printer.Fprint(buf, fset, r.Type)

				// extract the version from the struct type
				if reVersionedName.MatchString(buf.String()) {
					version = reVersionedName.FindStringSubmatch(buf.String())[1]
					break
				}

				// if this is a current version get the mapping from managers.go
				if reCurrentVersionName.MatchString(buf.String()) {
					version, ok = currentManagerVersions[strings.TrimPrefix(buf.String(), "*")]
					if !ok {
						fmt.Println("Error: could not find current version for", buf.String())
						os.Exit(1)
					}
					break
				}
			}
		}

		// if version was not found in receivers it is not versioned
		if version == "" {
			return false, 0
		}
	}

	fnVersion, _ := strconv.Atoi(version)
	return true, fnVersion
}

func hasBuildFlags(file *ast.File) bool {
	for _, comment := range file.Comments {
		if strings.Contains(comment.Text(), "+build") {
			return true
		}
	}
	return false
}

func hasMainnetBuildFlags(file *ast.File) bool {
	for _, comment := range file.Comments {
		if strings.Contains(comment.Text(), "+build") {
			if reNonMainnetBuildFlags.MatchString(comment.Text()) {
				return false
			}
		}
	}
	return true
}

// Returns true if the function is just a single return statement
func skipDedupe(fn *ast.FuncDecl) bool {
	if fn.Body == nil || len(fn.Body.List) != 1 {
		return false
	}

	_, ok := fn.Body.List[0].(*ast.ReturnStmt)
	return ok
}

// -------------------------------------------------------------------------------------
// Main
// -------------------------------------------------------------------------------------

func main() {
	// parse flags
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
		fmt.Println("Error parsing files:", err)
		os.Exit(1)
	}

	// extract current version for VCUR managers from managers.go
	for _, pkg := range pkgs {
		file, ok := pkg.Files["x/thorchain/managers.go"]
		if !ok {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			var n *ast.CaseClause
			n, ok = node.(*ast.CaseClause)
			if !ok {
				return true
			}

			// iterate switch cases
			for _, e := range n.List {
				version := ""

				ast.Inspect(e, func(n ast.Node) bool {
					if n == nil || version != "" {
						return false
					}
					var c *ast.CallExpr
					if c, ok = n.(*ast.CallExpr); ok {
						// extract the version from semver.MustParse argument
						var s *ast.SelectorExpr
						if s, ok = c.Fun.(*ast.SelectorExpr); ok {
							var x *ast.Ident
							if x, ok = s.X.(*ast.Ident); ok {
								if x.Name == "semver" && s.Sel.Name == "MustParse" {
									var l *ast.BasicLit
									if l, ok = c.Args[0].(*ast.BasicLit); ok {
										version = l.Value
										return false
									}
								}
							}
						}
					}

					return true
				})
				if version == "" {
					continue
				}

				// extract the minor version
				minor := strings.Split(version, ".")[1]

				// extract versioned functions called in the case body
				for _, s := range n.Body {
					ast.Inspect(s, func(n ast.Node) bool {
						if n == nil {
							return false
						}
						var c *ast.CallExpr
						if c, ok = n.(*ast.CallExpr); ok {
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
							if !strings.HasPrefix(vFn, "new") {
								return true
							}
							manager := strings.TrimPrefix(vFn, "new")

							// store the manager version and type mapping
							var v string
							if v, ok = currentManagerVersions[manager]; ok {
								if v != minor {
									fmt.Printf("Error: function version mismatch (%s): %s != %s\n", manager, v, minor)
								}
							}
							currentManagerVersions[manager] = minor
						}
						return true
					})
				}
			}

			return true
		})
	}

	// walk the ast and record all versioned functions
	fnsMap := map[token.Pos]ast.Node{}
	fnsDedupe := map[uint64][]string{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {

			// skip packages with no versioned functions relevant for consensus
			rootPackage := strings.Split(fset.File(file.Pos()).Name(), "/")[0]
			if skipRootPackages[rootPackage] {
				continue
			}

			if hasBuildFlags(file) {
				// explicitly disallow build flags on handler files
				if strings.HasPrefix(file.Name.Name, "handler") {
					fmt.Println("Error: build flags are not allowed on handler files")
					os.Exit(1)
				}

				// skip files with non mainnet build flags
				if !hasMainnetBuildFlags(file) {
					continue
				}
			}

			// skip test files
			if strings.HasSuffix(fset.File(file.Pos()).Name(), "_test.go") {
				continue
			}

			// skip generated files
			if strings.HasSuffix(fset.File(file.Pos()).Name(), ".pb.go") {
				continue
			}

			ast.Inspect(file, func(node ast.Node) bool {
				// remove all comments
				v := reflect.ValueOf(node)
				if v.Kind() == reflect.Ptr && !v.IsNil() {
					v = v.Elem()
				}
				if v.IsValid() {
					for _, field := range []string{"Doc", "Comments"} {
						field := v.FieldByName(field)
						if field.IsValid() {
							field.Set(reflect.Zero(field.Type()))
						}
					}
				}

				isVersioned, version := isVersionedFunction(node, fset)
				if isVersioned {
					fn, ok := node.(*ast.FuncDecl)
					if !ok {
						panic("unreachable")
					}

					fnHash := fnv.New64a()
					name := []string{}
					if fn.Recv != nil {
						for _, r := range fn.Recv.List {
							buf := new(bytes.Buffer)
							printer.Fprint(buf, fset, r.Type)
							name = append(name, buf.String())
							printer.Fprint(fnHash, fset, r.Type)
						}
					}
					printer.Fprint(fnHash, fset, fn.Body)

					// record all versioned functions outside current version
					if version != *flagVersion {
						fnsMap[node.Pos()] = node
					}

					// skip empty functions
					if skipDedupe(fn) {
						return true
					}

					// dedupe by function receiver and body
					if len(fn.Body.List) > 0 {
						name = append(name, fn.Name.Name)
						fnsDedupe[fnHash.Sum64()] = append(fnsDedupe[fnHash.Sum64()], strings.Join(name, "."))
					}
				}
				return true
			})
		}
	}

	// explicitly disallow duplicate versioned functions
	for _, fns := range fnsDedupe {
		if len(fns) > 1 {
			fmt.Fprintf(os.Stderr, "Error: duplicate versioned functions: %s\n", strings.Join(fns, ", "))
			os.Exit(1)
		}
	}

	// convert to slice for sorting
	fns := []ast.Node{}
	for _, fn := range fnsMap {
		fns = append(fns, fn)
	}

	// sort by function name in case filestructure changes
	sort.SliceStable(fns, func(i, j int) bool {
		fi, ok := fns[i].(*ast.FuncDecl)
		if !ok {
			panic("unreachable")
		}
		fj, ok := fns[j].(*ast.FuncDecl)
		if !ok {
			panic("unreachable")
		}

		ii := new(bytes.Buffer)
		if fi.Recv != nil {
			for _, f := range fi.Recv.List {
				printer.Fprint(ii, fset, f.Type)
			}
		}
		ii.WriteString(fi.Name.Name)
		printer.Fprint(ii, fset, fi.Type)

		jj := new(bytes.Buffer)
		if fj.Recv != nil {
			for _, f := range fj.Recv.List {
				printer.Fprint(jj, fset, f.Type)
			}
		}
		jj.WriteString(fj.Name.Name)
		printer.Fprint(jj, fset, fj.Type)

		// replace VCUR functions with version extracted from managers.go
		iis := ii.String()
		jjs := jj.String()
		for k, v := range currentManagerVersions {
			iis = strings.ReplaceAll(iis, k, strings.ReplaceAll(k, "CUR", v))
			jjs = strings.ReplaceAll(jjs, k, strings.ReplaceAll(k, "CUR", v))
		}

		return strings.ToLower(iis) < strings.ToLower(jjs)
	})

	// print package so gofumpt can format
	fmt.Println("package main")

	// print the versioned functions to buffer
	buf := new(bytes.Buffer)
	for _, fn := range fns {
		pos := fset.Position(fn.Pos())
		buf.WriteString(fmt.Sprintf("// %s:%d\n", pos.Filename, pos.Line))
		printer.Fprint(buf, fset, fn)
		buf.WriteString("\n\n")
	}

	// replace VCUR functions with version extracted from managers.go
	out := buf.String()
	for k, v := range currentManagerVersions {
		out = strings.ReplaceAll(out, k, strings.ReplaceAll(k, "CUR", v))
	}

	// print the versioned functions
	fmt.Print(out)
}
