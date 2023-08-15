package main

type Block struct {
	statements []Stmt
}

func (blockObj Block) accept(visitor Interpreter) LoxValue {
	return visitor.visitBlockStmt(blockObj)
}

type Expression struct {
	expression Expr
}

func (expressionObj Expression) accept(visitor Interpreter) LoxValue {
	return visitor.visitExpressionStmt(expressionObj)
}

type If struct {
	condition Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (ifObj If) accept(visitor Interpreter) LoxValue {
	return visitor.visitIfStmt(ifObj)
}

type Print struct {
	expression Expr
}

func (printObj Print) accept(visitor Interpreter) LoxValue {
	return visitor.visitPrintStmt(printObj)
}

type While struct {
	condition Expr
	body Stmt
}

func (whileObj While) accept(visitor Interpreter) LoxValue {
	return visitor.visitWhileStmt(whileObj)
}

type Var struct {
	name Token
	initializer Expr
	isConst bool
}

func (varObj Var) accept(visitor Interpreter) LoxValue {
	return visitor.visitVarStmt(varObj)
}

type Function struct {
	name Token
	params []Token
	body []Stmt
}

func (functionObj Function) accept(visitor Interpreter) LoxValue {
	return visitor.visitFunctionStmt(functionObj)
}

type Return struct {
	keyword Token
	value Expr
}

func (returnObj Return) accept(visitor Interpreter) LoxValue {
	return visitor.visitReturnStmt(returnObj)
}

