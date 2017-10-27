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
	for _, t := range os.Args[1:] {
		iface := GenInterface(t)
		printer.Fprint(os.Stdout, token.NewFileSet(), iface)
		fmt.Println()
	}
}

func GenInterface(structs string) *ast.GenDecl {
	fset := token.NewFileSet()
	f, err := parser.ParseDir(fset, ".", NoTest, 0)

	if err != nil {
		panic(err)
	}

	iface := &ast.InterfaceType{
		Methods: &ast.FieldList{
			List: make([]*ast.Field, 0),
		},
	}
	for _, pkg := range f {
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
				if typeName == structs && x.Name.IsExported() {

					iface.Methods.List = append(iface.Methods.List, &ast.Field{
						Names: []*ast.Ident{
							x.Name,
						},
						Type: x.Type,
					})
				}
			}

			return true
		})
	}
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: structs,
				},
				Type: iface,
			},
		},
	}
}

func NoTest(f os.FileInfo) bool {
	isTest, err := regexp.MatchString("_test.go", f.Name())
	if err != nil {
		panic(err)
	}
	return !isTest
}
