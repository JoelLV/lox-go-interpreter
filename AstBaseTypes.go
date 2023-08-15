package main

type Expr interface {
	accept(visitor Interpreter) LoxValue
}

type Stmt interface {
	accept(visitor Interpreter) LoxValue
}

type LoxCallable interface {
	call(interpreter Interpreter, arguments []LoxValue) LoxValue
}

type LoxValue interface {
}
