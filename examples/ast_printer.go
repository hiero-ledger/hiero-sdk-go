package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strconv"
	"strings"
)

const StatusIndex = 384
const ProtoCodeStructName = "ResponseCodeEnum_"

func main() {
	pbCodesFile := "../proto/services/response_code.pb.go" // Use full or relative path
	pbCodesFset := token.NewFileSet()

	// Parse the Go file
	pbCodesNode, err := parser.ParseFile(pbCodesFset, pbCodesFile, nil, parser.AllErrors)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	var newConstGenDecls []*ast.GenDecl
	var newSwitchCaseClause []*ast.CaseClause

	// Walk the AST to find const declarations related to ResponseCodeEnum
	ast.Inspect(pbCodesNode, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						if len(name.Name) >= len(ProtoCodeStructName) && name.Name[:len(ProtoCodeStructName)] == ProtoCodeStructName {
							if basicLit, ok := valueSpec.Values[len(valueSpec.Values)-1].(*ast.BasicLit); ok {
								if basicLit.Kind == token.INT {
									intValue, _ := strconv.Atoi(basicLit.Value)
									if intValue > StatusIndex {
										newValueSpec := &ast.ValueSpec{
											Names: []*ast.Ident{
												ast.NewIdent(strings.ReplaceAll(name.Name, ProtoCodeStructName, "")),
											},
											Type: ast.NewIdent("Status"),
											Values: []ast.Expr{
												&ast.BasicLit{
													Kind:  token.INT,
													Value: basicLit.Value,
												},
											},
										}
										newGenDecl := &ast.GenDecl{
											Tok:   token.CONST,
											Specs: []ast.Spec{newValueSpec},
										}
										newConstGenDecls = append(newConstGenDecls, newGenDecl)
									}
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	fmt.Println(newConstGenDecls)

	codesFile := "../sdk/status.go" // Use full or relative path
	codesFset := token.NewFileSet()

	// Parse the Go file
	codesNode, err := parser.ParseFile(codesFset, codesFile, nil, parser.AllErrors)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	// Walk the AST to find the const block
	ast.Inspect(codesNode, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			// Append the new constant declarations to the original const block
			for _, newGenDecl := range newConstGenDecls {
				genDecl.Specs = append(genDecl.Specs, newGenDecl.Specs...)
			}
		}
		return true
	})

	ast.Inspect(codesNode, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						fmt.Println(name.Name)
						fmt.Println(valueSpec.Values[len(valueSpec.Values)-1])
					}
				}
			}
		}
		return true
	})

	ast.Inspect(codesNode, func(n ast.Node) bool {
		if switchStmt, ok := n.(*ast.SwitchStmt); ok {
			for _, stmt := range switchStmt.Body.List {
				if caseClause, ok := stmt.(*ast.CaseClause); ok {
					for _, caseStmt := range caseClause.Body {
						if retStmt, ok := caseStmt.(*ast.ReturnStmt); ok {
							fmt.Println(retStmt.Return)
							fmt.Println(retStmt.Results[0])
						}
					}
				}
			}
		}
		return true
	})

	fmt.Println(newSwitchCaseClause)

	// Create a new file for output
	outFile, err := os.Create("output.go")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outFile.Close()

	// Print the Go source code from the AST to the file
	err = printer.Fprint(outFile, codesFset, codesNode)
	if err != nil {
		fmt.Println("Error printing file:", err)
		return
	}
}
