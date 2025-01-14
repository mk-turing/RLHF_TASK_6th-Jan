package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"regexp"
	"strings"
)

var (
	testFunctionRegexp = regexp.MustCompile(`^Test[A-Z][a-zA-Z0-9_]*_Should_[A-Z][a-zA-Z0-9_]*$`)
)

func main() {
	// Replace "path/to/your/package" with the actual package path
	if err := checkTestFunctionNames("path/to/your/package"); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func checkTestFunctionNames(packagePath string) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, packagePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.IsExported() && strings.HasPrefix(fn.Name.Name, "Test") {
					if !testFunctionRegexp.MatchString(fn.Name.Name) {
						log.Printf("Invalid test function name: %s. Expected format: %s", fn.Name.Name, testFunctionRegexp)
					}
				}
				return true
			})
		}
	}

	return nil
}
