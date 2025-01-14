package _12284

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = analysis.Analyzer{
	Name: "test_naming",
	Doc:  "Checks test function naming conventions",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		if !isTestFile(f) {
			continue
		}
		checkTestFunctions(pass, f)
	}
	return nil, nil
}

func isTestFile(f *ast.File) bool {
	return f.Name.String() == "_test.go"
}

func checkTestFunctions(pass *analysis.Pass, f *ast.File) {
	for _, decl := range f.Decls {
		if funcdecl, ok := decl.(*ast.FuncDecl); ok {
			if !isTestFunction(funcdecl) {
				continue
			}
			if !matchesConvention(funcdecl.Name.String()) {
				pass.Reportf(funcdecl.Name.Pos(), "test function name does not match convention")
			}
		}
	}
}

func isTestFunction(funcdecl *ast.FuncDecl) bool {
	return funcdecl.Recv == nil && funcdecl.Name.Name == "Test"
}

func matchesConvention(name string) bool {
	// Implement logic to check if the name matches the convention
	// e.g., check if it starts with "Test" and has an appropriate format
	return true
}
