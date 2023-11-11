package analyzer

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	gomockControllerType = "github.com/golang/mock/gomock.Controller"
	gomockPkg            = "github.com/golang/mock/gomock"
	finish               = "Finish"
	newControllerMethod  = "NewController"
	testingType          = "*testing.T"
)

// New returns new gomockcontrollerfinish analyzer.
func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "gomockcontrollerfinish",
		Doc:      "Checks whether an unnecessary call to .Finish() on gomock.Controller exists",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

	// track if gomock.NewController(t) is called in the current test file
	//
	// note that it doesn't have to be t as args
	var newControllerCalled bool

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}

		// check if it's a test file
		if !isTestFile(pass.Fset.Position(callExpr.Pos()).Filename) {
			return
		}

		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		selIdent, ok := selectorExpr.X.(*ast.Ident)
		if !ok {
			return
		}

		// get pkg name if expression is from Go package
		pkg, ok := pass.TypesInfo.ObjectOf(selIdent).(*types.PkgName)
		// check if it's gomock pkg and method is NewController
		if ok && strings.HasSuffix(pkg.Imported().Path(), gomockPkg) && selectorExpr.Sel.Name == newControllerMethod {
			if len(callExpr.Args) == 1 {
				if argType := pass.TypesInfo.TypeOf(callExpr.Args[0]); argType.String() == testingType {
					newControllerCalled = true
				}
			}
		}

		// check for unnecessary call to gomock.Controller.Finish()
		if newControllerCalled && strings.HasSuffix(pass.TypesInfo.TypeOf(selIdent).String(), gomockControllerType) && selectorExpr.Sel.Name == finish {
			pass.Reportf(selectorExpr.Sel.Pos(), "since go1.14, if you are passing a testing.T to NewController then calling Finish on gomock.Controller is no longer needed")
		}
	})

	return nil, nil
}

// isTestFile checks if the file is a test file based on its name.
func isTestFile(filename string) bool {
	return strings.HasSuffix(filename, "_test.go")
}
