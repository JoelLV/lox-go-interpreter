package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var stmtTypes = [][]string{
	{"Block", "statements []Stmt"},
	{"Expression", "expression Expr"},
	{"If", "condition Expr", "thenBranch Stmt", "elseBranch Stmt"},
	{"Print", "expression Expr"},
	{"While", "condition Expr", "body Stmt"},
	{"Var", "name Token", "initializer Expr", "isConst bool"},
	{"Function", "name Token", "params []Token", "body []Stmt"},
	{"Return", "keyword Token", "value Expr"},
}

var exprTypes = [][]string{
	{"Assign", "name Token", "value Expr"},
	{"Binary", "left Expr", "operator Token", "right Expr"},
	{"Call", "callee Expr", "paren Token", "arguments []Expr"},
	{"Grouping", "expression Expr"},
	{"Literal", "value LoxValue"},
	{"Logical", "left Expr", "operator Token", "right Expr"},
	{"Unary", "operator Token", "right Expr"},
	{"Variable", "name Token"},
}

func createAstTypes(varType string) {
	filename := fmt.Sprintf("../%sTypes.go", varType)
	file, err := os.Create(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	_, writeErr := file.WriteString("package main\n\n")

	if writeErr != nil {
		log.Fatal(writeErr)
	}

	var varTypeToWrite [][]string
	if varType == "Stmt" {
		varTypeToWrite = stmtTypes
	} else {
		varTypeToWrite = exprTypes
	}
	for i := 0; i < len(varTypeToWrite); i++ {
		typeStr := fmt.Sprintf("type %s struct {\n", varTypeToWrite[i][0])
		for ii := 1; ii < len(varTypeToWrite[i]); ii++ {
			typeStr += fmt.Sprintf("	%s\n", varTypeToWrite[i][ii])
		}
		typeStr += "}\n\n"
		typeStr += fmt.Sprintf("func (%sObj %s) accept(visitor Interpreter) LoxValue {\n", strings.ToLower(varTypeToWrite[i][0]), varTypeToWrite[i][0])
		typeStr += fmt.Sprintf("	return visitor.visit%s%s(%sObj)\n", varTypeToWrite[i][0], varType, strings.ToLower(varTypeToWrite[i][0]))
		typeStr += fmt.Sprintf("}\n\n")
		file.WriteString(typeStr)
	}
}

func main() {
	createAstTypes("Stmt")
	createAstTypes("Expr")
}
