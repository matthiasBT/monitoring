package main

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
)

var osExitMainAnalyzer = &analysis.Analyzer{
	Name: "osExitMainAnalyzer",
	Doc:  "Check for direct os.Exit calls in main functions",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		fmt.Printf("> %s\n", file.Name)
		ast.Inspect(file, func(node ast.Node) bool {
			fn, ok := node.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "main" {
				return true
			}
			if hasDirectOsExitCall(fn.Body) {
				pass.Reportf(fn.Pos(), "Avoid direct os.Exit calls in main function")
			}
			return false
		})
	}
	return nil, nil
}

func hasDirectOsExitCall(fnBody *ast.BlockStmt) bool {
	for _, stmt := range fnBody.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		callExpr, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		ident, ok := selExpr.X.(*ast.Ident)
		if !ok || ident.Name != "os" {
			continue
		}
		if selExpr.Sel.Name == "Exit" {
			return true
		}
	}
	return false
}

func main() {
	multichecker.Main(
		buildssa.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		ctrlflow.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shift.Analyzer,
		stdmethods.Analyzer,
		structtag.Analyzer,
		unsafeptr.Analyzer,
		unmarshal.Analyzer,
		unusedresult.Analyzer,

		osExitMainAnalyzer,
	)
}
