package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"strconv"
	"strings"
)

const StatusIndex = 98
const ProtoCodeStructName = "HederaFunctionality_"

func main() {
	// Paths to the input files
	pbCodesFile := "../../../proto/services/basic_types.pb.go"
	codesFile := "../../../sdk/request_type.go"

	// Parse the protobuf-generated Go file and extract new constants and case clauses
	pbCodesNode := parseFile(pbCodesFile)
	newConstGenDecls, newSwitchCaseClauses := extractNewConstants(pbCodesNode)

	// Parse the Go file where modifications will be applied
	codesNode := parseFile(codesFile)

	// Modify the AST to include the new constants and switch cases
	updateConstBlock(codesNode, newConstGenDecls)
	updateSwitchCases(codesNode, newSwitchCaseClauses)

	// Write the modified AST back to a new Go file
	writeToFile("output.go", codesNode)
}

// Parses a Go source file into an AST
func parseFile(filename string) *ast.File {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("Error parsing file %s: %v", filename, err)
	}
	return node
}

// Extracts new constants and corresponding switch case clauses from the parsed AST
func extractNewConstants(node *ast.File) ([]*ast.GenDecl, []*ast.CaseClause) {
	var newConstGenDecls []*ast.GenDecl
	var newSwitchCaseClauses []*ast.CaseClause

	ast.Inspect(node, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for _, name := range valueSpec.Names {
				if !strings.HasPrefix(name.Name, ProtoCodeStructName) {
					continue
				}

				basicLit, ok := valueSpec.Values[len(valueSpec.Values)-1].(*ast.BasicLit)
				if !ok || basicLit.Kind != token.INT {
					continue
				}

				intValue, _ := strconv.Atoi(basicLit.Value)
				if intValue <= StatusIndex {
					continue
				}

				// Extract constant name without prefix
				cleanName := strings.TrimPrefix(name.Name, ProtoCodeStructName)

				// Create a new constant declaration
				newValueSpec := &ast.ValueSpec{
					Names:  []*ast.Ident{ast.NewIdent(cleanName)},
					Values: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: basicLit.Value}},
				}
				newConstGenDecls = append(newConstGenDecls, &ast.GenDecl{
					Tok:   token.CONST,
					Specs: []ast.Spec{newValueSpec},
				})

				// Create a new switch case clause
				newCase := &ast.CaseClause{
					List: []ast.Expr{ast.NewIdent(cleanName)},
					Body: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: strconv.Quote(cleanName),
								},
							},
						},
					},
				}
				newSwitchCaseClauses = append(newSwitchCaseClauses, newCase)
			}
		}
		return true
	})

	return newConstGenDecls, newSwitchCaseClauses
}

// Adds new constant declarations to the existing const block in the AST
func updateConstBlock(node *ast.File, newConstGenDecls []*ast.GenDecl) {
	ast.Inspect(node, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			return true
		}

		// Append new constants to the existing const block
		for _, newGenDecl := range newConstGenDecls {
			genDecl.Specs = append(genDecl.Specs, newGenDecl.Specs...)
		}
		return false // Stop further traversal after modifying the first const block
	})
}

func updateSwitchCases(node *ast.File, newSwitchCaseClauses []*ast.CaseClause) {
	ast.Inspect(node, func(n ast.Node) bool {
		switchStmt, ok := n.(*ast.SwitchStmt)
		if !ok {
			return true
		}

		// Convert []*ast.CaseClause to []ast.Stmt
		var newStmts []ast.Stmt
		for _, caseClause := range newSwitchCaseClauses {
			newStmts = append(newStmts, caseClause)
		}

		// Append new cases to the switch statement
		switchStmt.Body.List = append(switchStmt.Body.List, newStmts...)

		return false // Stop further traversal after modifying the first switch statement
	})
}

// Writes the modified AST back to a new Go source file
func writeToFile(filename string, node *ast.File) {
	outFile, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error creating file %s: %v", filename, err)
	}
	defer outFile.Close()

	// Print the modified AST to the file
	fset := token.NewFileSet()
	if err := printer.Fprint(outFile, fset, node); err != nil {
		log.Fatalf("Error writing to file %s: %v", filename, err)
	}
}
