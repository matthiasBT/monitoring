// Package main demonstrates the setup of a multichecker which integrates
// various static analysis tools for Go code. This multichecker combines
// standard analyzers from the golang.org/x/tools/go/analysis/passes package,
// Staticcheck, Simple, Stylecheck, Errcheck, and Go-Critic checkers, along
// with a custom analyzer.

// The custom analyzer, osExitMainAnalyzer, checks for direct calls to
// os.Exit within the main function, a practice generally discouraged as it
// bypasses deferred function calls and can make cleanup and error handling
// more complex.

// This package iterates over a list of analyzers, appending them to a
// multichecker instance. This approach allows for comprehensive static analysis
// covering a wide range of common Go programming mistakes, style issues, and
// potential bugs.
package main

import (
	"go/ast"
	"strings"

	gocritic "github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// osExitMainAnalyzer is a custom analyzer that checks for direct usages of
// os.Exit in the main function. This is generally discouraged in Go as it
// prevents deferred functions from running.
var osExitMainAnalyzer = &analysis.Analyzer{
	Name: "osExitMainAnalyzer",
	Doc:  "Check for direct os.Exit calls in main functions",
	Run:  run,
}

// osExitMainAnalyzer is a custom analyzer that checks for direct usages of
// os.Exit in the main function. This is generally discouraged in Go as it
// prevents deferred functions from running.
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
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

// hasDirectOsExitCall checks if an AST node (function body) contains
// a direct call to os.Exit.
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

// function main is the entrypoint of this program
// in order to launch the program from the command line, execute the following commands:
// $ go run cmd/staticlint/multichecker.go ./cmd/...
// $ go run cmd/staticlint/multichecker.go ./internal/...
func main() {
	// main initializes and runs the multichecker with a comprehensive
	// set of analyzers including standard, third-party, and custom analyzers.
	var mychecks = []*analysis.Analyzer{
		// all standard analyzers from "golang.org/x/tools/go/analysis/passes"
		// refer to their documentation for more details
		appends.Analyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		defers.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		slog.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,

		// the custom analyzer (see above)
		osExitMainAnalyzer,
	}

	// enabling all SA-analyzers from staticcheck
	for _, a := range staticcheck.Analyzers {
		if strings.HasPrefix(a.Analyzer.Name, "SA") {
			mychecks = append(mychecks, a.Analyzer)
		}
	}

	// enabling all S-analyzers from staticcheck
	for _, a := range simple.Analyzers {
		if strings.HasPrefix(a.Analyzer.Name, "S") {
			mychecks = append(mychecks, a.Analyzer)
		}
	}

	// enabling all ST-analyzers from staticcheck
	for _, a := range stylecheck.Analyzers {
		if strings.HasPrefix(a.Analyzer.Name, "ST") {
			mychecks = append(mychecks, a.Analyzer)
		}
	}

	// enabling the errcheck analyzer
	mychecks = append(mychecks, errcheck.Analyzer)

	// enabling the gocritic analyzer
	mychecks = append(mychecks, gocritic.Analyzer)

	// launch the analyzers
	multichecker.Main(
		mychecks...,
	)
}
