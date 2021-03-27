package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func findGoStmtInDecl(fset *token.FileSet, funcdecl *ast.FuncDecl) {
	fmt.Println("Deal main function")
	var i int
	for i = 0; i < len(funcdecl.Body.List); i++ {
		stmt := funcdecl.Body.List[i]
		switch stmt.(type) {
		case (*ast.GoStmt):
			gostmt, ok := stmt.(*ast.GoStmt)
			if !ok {
				fmt.Println("can't transfer to gostmt")
			}
			funexpr := gostmt.Call.Fun

			switch funexpr.(type) {
			case (*ast.Ident):
				ident, ok := funexpr.(*ast.Ident)
				if !ok {
					fmt.Println("can't transfer to ident")
				}
				fmt.Println("[Go Statement]", "<", gostmt.Go, ">", ident.Name)
				break
			default:
				fmt.Println(funexpr)
			}

		default:
			continue
		}
	}
}

func main() {
	fset := token.NewFileSet()
	filename := "examples/example1.go"
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		panic(err)
	}

	//Print the AST
	//ast.Print(fset, f)

	//find out all gostmt
	var decl ast.Decl
	var i int
	for i = 0; i < len(f.Decls); i++ {
		decl = f.Decls[i]
		switch decl.(type) {
		case *ast.FuncDecl:

			fmt.Println(i, "FuncDecl")
			funcdecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				fmt.Println("not a Funcdecl type")
				return
			}
			if funcdecl.Name.Name != "main" {
				fmt.Println(i, "[User Define Func]", funcdecl.Name.Name)
				continue
			}

			//handle func main
			fmt.Println(i, "[Main Func]")
			findGoStmtInDecl(fset, funcdecl)

		default:
			fmt.Println(i, "Other Decl")
		}

	}
}
