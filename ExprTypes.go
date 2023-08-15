package main

type Assign struct {
	name Token
	value Expr
}

func (assignObj Assign) accept(visitor Interpreter) LoxValue {
	return visitor.visitAssignExpr(assignObj)
}

type Binary struct {
	left Expr
	operator Token
	right Expr
}

func (binaryObj Binary) accept(visitor Interpreter) LoxValue {
	return visitor.visitBinaryExpr(binaryObj)
}

type Call struct {
	callee Expr
	paren Token
	arguments []Expr
}

func (callObj Call) accept(visitor Interpreter) LoxValue {
	return visitor.visitCallExpr(callObj)
}

type Grouping struct {
	expression Expr
}

func (groupingObj Grouping) accept(visitor Interpreter) LoxValue {
	return visitor.visitGroupingExpr(groupingObj)
}

type Literal struct {
	value LoxValue
}

func (literalObj Literal) accept(visitor Interpreter) LoxValue {
	return visitor.visitLiteralExpr(literalObj)
}

type Logical struct {
	left Expr
	operator Token
	right Expr
}

func (logicalObj Logical) accept(visitor Interpreter) LoxValue {
	return visitor.visitLogicalExpr(logicalObj)
}

type Unary struct {
	operator Token
	right Expr
}

func (unaryObj Unary) accept(visitor Interpreter) LoxValue {
	return visitor.visitUnaryExpr(unaryObj)
}

type Variable struct {
	name Token
}

func (variableObj Variable) accept(visitor Interpreter) LoxValue {
	return visitor.visitVariableExpr(variableObj)
}

