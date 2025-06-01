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

var STATICCHECK_SA = []string{"sa1000", "sa1001", "sa1002", "sa1003", "sa1004", "sa1005", "sa1006", "sa1007", "sa1008", "sa1010", "sa1011", "sa1012", "sa1013", "sa1014", "sa1015", "sa1016", "sa1017", "sa1018", "sa1019", "sa1020", "sa1021", "sa1023", "sa1024", "sa1025", "sa1026", "sa1027", "sa1028", "sa1029", "sa1030", "sa1031", "sa1032", "sa2000", "sa2001", "sa2002", "sa2003", "sa3000", "sa3001", "sa4000", "sa4001", "sa4003", "sa4004", "sa4005", "sa4006", "sa4008", "sa4009", "sa4010", "sa4011", "sa4012", "sa4013", "sa4014", "sa4015", "sa4016", "sa4017", "sa4018", "sa4019", "sa4020", "sa4021", "sa4022", "sa4023", "sa4024", "sa4025", "sa4026", "sa4027", "sa4028", "sa4029", "sa4030", "sa4031", "sa4032", "sa5000", "sa5001", "sa5002", "sa5003", "sa5004", "sa5005", "sa5007", "sa5008", "sa5009", "sa5010", "sa5011", "sa5012", "sa6000", "sa6001", "sa6002", "sa6003", "sa6005", "sa6006", "sa9001", "sa9002", "sa9003", "sa9004", "sa9005", "sa9006", "sa9007", "sa9008", "sa9009"}

// ErrOsExitAnalizer анализатор фиксирующий использование функции os.Exit в main функциях
var ErrOsExitAnalizer = &analysis.Analyzer{
	Name: "err_os_exit_analizer",
	Doc:  "check call of os.Exit in the main.main functions",
	Run:  checkOsExit,
}

func checkOsExit(pass *analysis.Pass) (interface{}, error) {
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
	mychecks := make([]*analysis.Analyzer, 0, len(staticcheck.Analyzers))
	for _, analys := range staticcheck.Analyzers {
		if strings.HasPrefix(strings.ToUpper(analys.Analyzer.Name), "SA") {
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
