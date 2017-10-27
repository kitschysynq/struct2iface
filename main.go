// Package main provides a parser which returns all public methods of a struct as an interface
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"regexp"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseDir(fset, ".", nil, 0)

	if err != nil {
		panic(err)
	}

	iface := &ast.InterfaceType{
		Methods: &ast.FieldList{
			List: make([]*ast.Field, 0),
		},
	}
	for p, pkg := range f {
		fmt.Printf("Package: %q\n", p)
		ast.Inspect(pkg, func(n ast.Node) bool {
			switch x := n.(type) {

			// Function declaration
			case *ast.FuncDecl:
				if x.Recv == nil {
					break
				}
				var typeName string
				if recv, ok := x.Recv.List[0].Type.(*ast.StarExpr); ok {
					ident := recv.X.(*ast.Ident)
					typeName = ident.Name
				}
				if recv, ok := x.Recv.List[0].Type.(*ast.Ident); ok {
					typeName = recv.Name
				}
				if typeName == "Client" && x.Name.IsExported() {
					iface.Methods.List = append(iface.Methods.List, &ast.Field{
						Names: []*ast.Ident{
							x.Name,
						},
						Type: x.Type,
					})
					fmt.Printf("function %q with receiver type %q\n", x.Name.String(), typeName)
				}
			}

			return true
		})
	}
	printer.Fprint(os.Stdout, token.NewFileSet(), iface)
}

func NoTest(f os.FileInfo) bool {
	isTest, err := regexp.MatchString("_test.go$", f.Name())
	if err != nil {
		return false
	}
	return isTest
}
