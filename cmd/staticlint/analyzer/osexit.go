package analyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"strings"
)

var AnalyzerOsExit = &analysis.Analyzer{
	Name: "checkosexit",
	Doc:  "check for os.Exit in main function",
	Run:  runExitCheck,
}

func runExitCheck(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		////  find main pkg
		if strings.ToLower(f.Name.Name) == "main" {
			//  find main func in child nodes
			ast.Inspect(f, func(node ast.Node) bool {
				//  check if node is func declaration and func name == main
				if funcMain, ok := node.(*ast.FuncDecl); ok && funcMain.Name.Name == "main" {
					checkForExit(node, pass)
				}
				return true
			})
		}
	}
	return nil, nil
}

//  checkForExit checks not nil node childs for os.Exit
func checkForExit(node ast.Node, pass *analysis.Pass) {
	ast.Inspect(node, func(node ast.Node) bool {
		// some nodes may be nil
		if node != nil {
			// find (stand-alone) expression
			if call, ok := node.(*ast.CallExpr); ok {
				// find os.Exit
				// todo: check for 'os' aliases
				if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
					if pkg, ok := selector.X.(*ast.Ident); ok && pkg.String() == "os" && selector.Sel.Name == "Exit" {
						pass.Reportf(node.Pos(), "os.Exit call in main.go")
						return true
					}
				}
			}
		}
		return true
	})
}
