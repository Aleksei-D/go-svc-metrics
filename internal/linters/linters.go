// модуль linters предоставляет набор костомных чекеров
package linters

import (
	"go/ast"

	"github.com/cybozu-go/golang-custom-analyzer/pkg/eventuallycheck"
	"github.com/cybozu-go/golang-custom-analyzer/pkg/restrictpkg"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

var analyzers []*analysis.Analyzer = []*analysis.Analyzer{
	printf.Analyzer,
	shadow.Analyzer,
	structtag.Analyzer,
	unmarshal.Analyzer,
	unusedresult.Analyzer,
	eventuallycheck.Analyzer,
	restrictpkg.RestrictPackageAnalyzer,
	OsExitInMainAnalyzer,
}

// стилистические чекеры
var styleChecks map[string]bool = map[string]bool{
	"ST1000": true,
	"ST1001": true,
	"ST1002": true,
	"ST1003": true,
}

// проверки, которые использует gopls для автоматического рефакторинга
var quickFixChecks map[string]bool = map[string]bool{
	"QF1000": true,
	"QF1001": true,
}

type fileInspector struct {
	pass *analysis.Pass
}

type funcInsespector struct {
	pass     *analysis.Pass
	funcDecl *ast.FuncDecl
}

// NewCheckers возразает набор для мультичекера
func NewCheckers() []*analysis.Analyzer {
	var checks []*analysis.Analyzer

	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	for _, v := range stylecheck.Analyzers {
		if styleChecks[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}

	}

	for _, v := range quickfix.Analyzers {
		if quickFixChecks[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}

	}

	checks = append(checks, analyzers...)
	return checks
}

// OsExitInMainAnalyzer реализует интерфейс анализатора для проверки отсутствия вызова os.Exit в функции main  пакета main
var OsExitInMainAnalyzer = &analysis.Analyzer{
	Name: "osexitinmainanalyzer",
	Doc:  "checks for os.exit in packages 'main'",
	Run:  runOsExitInMain,
}

func runOsExitInMain(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}
	filesProcessing(pass)
	return nil, nil
}

func filesProcessing(pass *analysis.Pass) {
	fileInspector := fileInspector{pass: pass}
	for _, file := range pass.Files {
		ast.Inspect(file, fileInspector.run)
	}
}

func (f *fileInspector) run(node ast.Node) bool {
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		f.funcInspecting(funcDecl)
	}
	return true
}

func (f *fileInspector) funcInspecting(funcDecl *ast.FuncDecl) {
	if funcDecl.Name.Name == "main" {
		return
	}
	funcInsespector := funcInsespector{pass: f.pass, funcDecl: funcDecl}
	ast.Inspect(funcDecl, funcInsespector.run)
}

func (f *funcInsespector) run(node ast.Node) bool {
	if c, ok := node.(*ast.CallExpr); ok {
		f.callExprNodeInsepcting(c)
	}
	return true
}

func (f *funcInsespector) callExprNodeInsepcting(callExpr *ast.CallExpr) {
	s, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	if s.Sel.Name != "Exit" {
		return
	}

	i, ok := s.X.(*ast.Ident)
	if !ok {
		return
	}

	if i.Name == "os" {
		f.pass.Reportf(f.funcDecl.Pos(), "package %s should not contain a 'os.Exit' in main function", f.pass.Pkg.Name())
	}
}
