package main

import (
	"bytes"
	"go/ast"
	"go/printer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"strings"
)

func isTestFile(pkgPath string) bool {
	return strings.Contains(pkgPath, "_test.go")
}

func isLocalPackage(pkgPath string) bool {
	const projectPath = "github.com/xChygyNx/metrical"
	projectPathPrefix := projectPath + "/"
	return pkgPath == projectPath || strings.HasPrefix(pkgPath, projectPathPrefix)
}

// ErrOsExitAnalizer анализатор фиксирующий использование функции os.Exit в main функциях
var ErrOsExitAnalizer = &analysis.Analyzer{
	Name: "err_os_exit_analizer",
	Doc:  "check call of os.Exit in the main.main functions",
	Run:  checkOsExit,
}

func checkOsExit(pass *analysis.Pass) (interface{}, error) {
	pkgPath := pass.Pkg.Path()

	// Проверяем, что пакет находится внутри проекта
	if !isLocalPackage(pkgPath) {
		return nil, nil
	}

	// Проверяем, что это не тест файлы
	if isTestFile(pkgPath) {
		return nil, nil
	}
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.CallExpr:
				buf := bytes.NewBuffer(make([]byte, 0))
				printer.Fprint(buf, pass.Fset, x)
				funcName := strings.Trim(buf.String(), " ")
				if strings.HasPrefix(funcName, "os.Exit") {
					return true
				}
			}
			return false
		})
	}
	return nil, nil
}

func main() {
	needStaticChecks := map[string]bool{
		"S1003":  true, // Replace call to strings.Index with strings.Contains
		"ST1008": true, // A function’s error value should be its last return value
		"QF1004": true, // Use strings.ReplaceAll instead of strings.Replace with n == -1
	}
	mychecks := make([]*analysis.Analyzer, 0, len(staticcheck.Analyzers))
	for _, analys := range staticcheck.Analyzers {
		_, ok := needStaticChecks[strings.ToUpper(analys.Analyzer.Name)]
		if strings.HasPrefix(strings.ToUpper(analys.Analyzer.Name), "SA") || ok {
			mychecks = append(mychecks, analys.Analyzer)
		}
	}

	publicAnalyzers := []*analysis.Analyzer{
		printf.Analyzer,    // анализатор, фиксирующий несоответствие масок вывода типам выводимых агументов
		shadow.Analyzer,    // анализатор, фиксирующий "затенение" внешних переменных локальными переменными
		structtag.Analyzer, // анализатор проверяющий соответствие тэгов полям в структурах
	}
	mychecks = append(mychecks, publicAnalyzers...)
	mychecks = append(mychecks, ErrOsExitAnalizer)
	multichecker.Main(
		mychecks...,
	)
}
