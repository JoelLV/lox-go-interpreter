package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Interpreter struct {
	trees []Stmt
	env   *Environment
}

func (inter *Interpreter) init(stmtArr []Stmt) {
	/*Initializes a interpreter
	object.
	*/
	var env Environment
	env.init()
	inter.trees = stmtArr
	inter.env = &env

	inter.env.define("clock", ClockFunction{})
	inter.env.define("toString", ToStringFunction{})
	inter.env.define("input", InputFunction{})
	inter.env.define("parseString", ParseFunction{})
	inter.env.define("isInstance", IsInstanceFunction{})
}

func (inter *Interpreter) interpret() {
	/*Interprets the abstract syntax
	trees passed in the initialized of
	the interpreter.
	*/
	defer func() {
		if r := recover(); r != nil {
			exc := r.(LoxException)
			runtimeError(exc)
		}
	}()
	for i := 0; i < len(inter.trees); i++ {
		inter.execute(inter.trees[i])
	}
}

func (inter *Interpreter) execute(stmt Stmt) {
	/*Executes given statement.
	 */
	stmt.accept(*inter)
}

func (inter *Interpreter) evaluate(expr Expr) LoxValue {
	/*Evaluates given expression.
	 */
	return expr.accept(*inter)
}

func (inter *Interpreter) stringify(value LoxValue) string {
	/*Returns a human readable string representing
	the result of the interpretation.
	*/
	if value == nil {
		return "nil"
	} else if inter.isInt(value) {
		return fmt.Sprintf("%d", value.(int64))
	} else if inter.isFloat(value) {
		return fmt.Sprintf("%f", value.(float64))
	} else if boolean, isBool := value.(bool); isBool {
		return strconv.FormatBool(boolean)
	} else if _, isString := value.(string); isString {
		value := value.(string)
		return strings.Replace(value, "\"", "", -1)
	} else {
		if inter.isClockFunction(value) {
			function := value.(ClockFunction)
			return function.String()
		} else if inter.isToStringFunction(value) {
			function := value.(ToStringFunction)
			return function.String()
		} else if inter.isInputFunction(value) {
			function := value.(InputFunction)
			return function.String()
		} else if inter.isParseStringFunction(value) {
			function := value.(ParseFunction)
			return function.String()
		} else if inter.isInstanceFunction(value) {
			function := value.(IsInstanceFunction)
			return function.String()
		} else {
			function := value.(LoxFunction)
			return function.String()
		}
	}
}

func (inter *Interpreter) visitVarStmt(stmt Var) Stmt {
	/*Returns the evaluation of
	the variable declaration statement.
	*/
	var value LoxValue
	if inter.env.varExists(stmt.name) {
		panic(LoxException{token: stmt.name, message: fmt.Sprintf("Variable '%s' already exists.", stmt.name.lexeme)})
	}
	if stmt.initializer != nil {
		value = inter.evaluate(stmt.initializer)
	}
	inter.env.define(stmt.name.lexeme, value)
	if stmt.isConst {
		inter.env.constValues[stmt.name.lexeme] = true
	}

	return nil
}

func (inter *Interpreter) visitIfStmt(stmt If) Stmt {
	/*Returns the evaluation of the if statement.
	 */
	if inter.isTruthy(inter.evaluate(stmt.condition)) {
		inter.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		inter.execute(stmt.elseBranch)
	}

	return nil
}

func (inter *Interpreter) visitExpressionStmt(stmt Expression) Stmt {
	/*Returns the evaluation of
	the expression statement.
	*/
	inter.evaluate(stmt.expression)

	return nil
}

func (inter *Interpreter) visitPrintStmt(stmt Print) Stmt {
	/*Returns the evaluation of
	the print statement.
	*/
	value := inter.evaluate(stmt.expression)
	fmt.Println(inter.stringify(value))

	return nil
}

func (inter *Interpreter) visitWhileStmt(stmt While) Stmt {
	/*Returns the evaluation of
	a while statement.
	*/
	for inter.isTruthy(inter.evaluate(stmt.condition)) {
		inter.execute(stmt.body)
	}

	return nil
}

func (inter *Interpreter) visitBlockStmt(stmt Block) Stmt {
	/*Calls executeBlock method to
	execute all statements within
	the block.
	*/
	var blockEnv Environment
	blockEnv.init()
	blockEnv.enclosing = inter.env
	inter.executeBlock(stmt.statements, &blockEnv)

	return nil
}

func (inter *Interpreter) executeBlock(statements []Stmt, env *Environment) {
	/*Executes all statements
	within the block.
	*/
	previousEnv := inter.env

	defer func() {
		inter.env = previousEnv
	}()
	inter.env = env
	for i := 0; i < len(statements); i++ {
		inter.execute(statements[i])
	}
}

func (inter *Interpreter) visitFunctionStmt(stmt Function) Stmt {
	/*Creates a LoxFunction object and
	defines an environment.
	*/
	function := LoxFunction{stmt, inter.env}
	inter.env.define(stmt.name.lexeme, function)

	return nil
}

func (inter *Interpreter) visitReturnStmt(stmt Return) Stmt {
	/*Executes the return statement by
	raising an exception with a
	Return object attached to it.
	*/
	var value LoxValue
	if stmt.value != nil {
		value = inter.evaluate(stmt.value)
	}
	panic(RuntimeReturn{value: value})
}

func (inter *Interpreter) visitAssignExpr(expr Assign) LoxValue {
	/*Returns the evaluation of an assignment
	expression.
	*/
	if inter.env.isConst(expr.name) {
		panic(LoxException{token: expr.name, message: fmt.Sprintf("Cannot reassign constant variable '%s'.", expr.name.lexeme)})
	}
	value := inter.evaluate(expr.value)
	inter.env.assign(expr.name, value)
	return value
}

func (inter *Interpreter) visitBinaryExpr(expr Binary) LoxValue {
	/*Returns the representation of a binary
	expression.
	*/
	left := inter.evaluate(expr.left)
	right := inter.evaluate(expr.right)

	switch expr.operator.tokenType {
	case MOD:
		{
			inter.checkNumberOperands(expr.operator, left, right)
			result := math.Mod(inter.convertNumToFloat(left), inter.convertNumToFloat(right))
			if inter.isInt(left) && inter.isInt(right) {
				return int64(result)
			} else {
				return result
			}
		}
	case MINUS:
		{
			inter.checkNumberOperands(expr.operator, left, right)
			if inter.isInt(left) && inter.isInt(right) {
				return left.(int64) - right.(int64)
			} else {
				leftFloat, rightFloat := inter.convertNumsToFloat(left, right)
				return leftFloat - rightFloat
			}
		}
	case SLASH:
		{
			inter.checkNumberOperands(expr.operator, left, right)
			if (math.Mod(inter.convertNumToFloat(left), inter.convertNumToFloat(right)) == 0) && inter.isInt(left) && inter.isInt(right) {
				return left.(int64) / right.(int64)
			} else {
				leftFloat, rightFloat := inter.convertNumsToFloat(left, right)
				return leftFloat / rightFloat
			}
		}
	case STAR:
		{
			inter.checkNumberOperands(expr.operator, left, right)
			if inter.isInt(left) && inter.isInt(right) {
				return left.(int64) * right.(int64)
			} else {
				leftFloat, rightFloat := inter.convertNumsToFloat(left, right)
				return leftFloat * rightFloat
			}
		}
	case PLUS:
		{
			if inter.isString(left) && inter.isString(right) {
				return fmt.Sprintf("%s%s", left.(string), right.(string))
			} else if inter.isInt(left) && inter.isInt(right) {
				return left.(int64) + right.(int64)
			} else if inter.isNumber(left) && inter.isNumber(right) {
				leftFloat, rightFloat := inter.convertNumsToFloat(left, right)
				return leftFloat + rightFloat
			} else {
				panic(LoxException{token: expr.operator, message: "Operands must be two numbers or two strings."})
			}
		}
	case GREATER_THAN:
		{
			inter.checkNumberOperands(expr.operator, left, right)
			if inter.isInt(left) && inter.isInt(right) {
				return left.(int64) > right.(int64)
			} else {
				leftFloat, rightFloat := inter.convertNumsToFloat(left, right)
				return leftFloat > rightFloat
			}
		}
	case GREATER_EQUAL:
		{
			inter.checkNumberOperands(expr.operator, left, right)
			if inter.isInt(left) && inter.isInt(right) {
				return left.(int64) >= right.(int64)
			} else {
				leftFloat, rightFloat := inter.convertNumsToFloat(left, right)
				return leftFloat >= rightFloat
			}
		}
	case LESS_THAN:
		{
			inter.checkNumberOperands(expr.operator, left, right)
			if inter.isInt(left) && inter.isInt(right) {
				return left.(int64) < right.(int64)
			} else {
				leftFloat, rightFloat := inter.convertNumsToFloat(left, right)
				return leftFloat < rightFloat
			}
		}
	case LESS_EQUAL:
		{
			inter.checkNumberOperands(expr.operator, left, right)
			if inter.isInt(left) && inter.isInt(right) {
				return left.(int64) <= right.(int64)
			} else {
				leftFloat, rightFloat := inter.convertNumsToFloat(left, right)
				return leftFloat <= rightFloat
			}
		}
	case NOT_EQUAL:
		{
			return !inter.isEqual(left, right)
		}
	case EQUAL_EQUAL:
		{
			return inter.isEqual(left, right)
		}
	default:
		return nil
	}
}

func (inter *Interpreter) visitCallExpr(expr Call) LoxValue {
	/*Evaluates a call expression.
	 */
	callee := inter.evaluate(expr.callee)

	var arguments []LoxValue
	for i := 0; i < len(expr.arguments); i++ {
		arguments = append(arguments, inter.evaluate(expr.arguments[i]))
	}
	if _, isLoxCallable := callee.(LoxCallable); !isLoxCallable {
		panic(LoxException{token: expr.paren, message: "Can only call functions"})
	}

	defer func() {
		r := recover()
		if r != nil {
			if functionErr, isFuncError := r.(FunctionException); isFuncError {
				functionError(expr.paren.line, functionErr.message)
			} else {
				panic(r)
			}
		}
	}()

	if inter.isClockFunction(callee) {
		function := callee.(ClockFunction)
		if len(arguments) != function.arity() {
			panic(LoxException{token: expr.paren, message: fmt.Sprintf("Expected %d arguments but got %d", function.arity(), len(arguments))})
		}
		return function.call(*inter, arguments)
	} else if inter.isToStringFunction(callee) {
		function := callee.(ToStringFunction)
		if len(arguments) != function.arity() {
			panic(LoxException{token: expr.paren, message: fmt.Sprintf("Expected %d arguments but got %d", function.arity(), len(arguments))})
		}
		return function.call(*inter, arguments)
	} else if inter.isInputFunction(callee) {
		function := callee.(InputFunction)
		if len(arguments) != function.arity() {
			panic(LoxException{token: expr.paren, message: fmt.Sprintf("Expected %d arguments but got %d", function.arity(), len(arguments))})
		}
		return function.call(*inter, arguments)
	} else if inter.isParseStringFunction(callee) {
		function := callee.(ParseFunction)
		if len(arguments) != function.arity() {
			panic(LoxException{token: expr.paren, message: fmt.Sprintf("Expected %d arguments but got %d", function.arity(), len(arguments))})
		}
		return function.call(*inter, arguments)
	} else if inter.isInstanceFunction(callee) {
		function := callee.(IsInstanceFunction)
		if len(arguments) != function.arity() {
			panic(LoxException{token: expr.paren, message: fmt.Sprintf("Expected %d arguments but got %d", function.arity(), len(arguments))})
		}
		return function.call(*inter, arguments)
	} else {
		function := callee.(LoxFunction)
		if len(arguments) != function.arity() {
			panic(LoxException{token: expr.paren, message: fmt.Sprintf("Expected %d arguments but got %d", function.arity(), len(arguments))})
		}
		return function.call(*inter, arguments)
	}
}

func (inter *Interpreter) visitGroupingExpr(expr Grouping) LoxValue {
	/*Returns the evaluation of the expression
	enclosed in parenthesis.
	*/
	return inter.evaluate(expr.expression)
}

func (inter *Interpreter) visitLiteralExpr(expr Literal) LoxValue {
	/*Returns the value of a literal expression.
	 */
	if inter.isInt(expr.value) {
		return expr.value.(int64)
	} else if inter.isFloat(expr.value) {
		return expr.value.(float64)
	} else {
		return expr.value
	}
}

func (inter *Interpreter) visitUnaryExpr(expr Unary) LoxValue {
	/*Returns the evaluation of
	a unary expression.
	*/
	right := inter.evaluate(expr.right)

	if expr.operator.tokenType == MINUS {
		inter.checkNumberOperand(expr.operator, right)
		if inter.isInt(right) {
			return -right.(int64)
		} else {
			return -right.(float64)
		}
	} else if expr.operator.tokenType == NOT {
		return !inter.isTruthy(right)
	} else {
		return nil
	}
}

func (inter *Interpreter) visitLogicalExpr(expr Logical) LoxValue {
	/*Executes a logical expression.
	 */
	left := inter.evaluate(expr.left)

	if expr.operator.tokenType == OR {
		if inter.isTruthy(left) {
			return left
		}
	} else if expr.operator.tokenType == AND {
		if !inter.isTruthy(left) {
			return left
		}
	}

	return inter.evaluate(expr.right)
}

func (inter *Interpreter) visitVariableExpr(expr Variable) LoxValue {
	/*Returns the evaluation of
	the variable expression.
	*/
	return inter.env.get(expr.name)
}

func (inter *Interpreter) isTruthy(obj LoxValue) bool {
	/*Determines whether the object
	is truthy or falsey.
	*/
	if obj == nil {
		return false
	} else if _, isBool := obj.(bool); isBool {
		return obj.(bool)
	} else {
		return true
	}
}

func (inter *Interpreter) checkNumberOperand(operator Token, operand LoxValue) {
	/*Checks if operand is a number. Otherwise
	it throws an exception.
	*/
	if !inter.isNumber(operand) {
		panic(LoxException{token: operator, message: "Operand must be a number"})
	}
}

func (inter *Interpreter) checkNumberOperands(operator Token, left LoxValue, right LoxValue) {
	/*Checks if both left and right are
	numbers, otherwise it throws an
	exception.
	*/
	if !inter.isNumber(left) || !inter.isNumber(right) {
		panic(LoxException{token: operator, message: "Operands must be numbers"})
	}
}

func (inter *Interpreter) isNumber(value LoxValue) bool {
	/*Checks if the value is either a float
	or an int.
	*/
	return inter.isInt(value) || inter.isFloat(value)
}

func (inter *Interpreter) convertNumToInt(value LoxValue) int64 {
	/*Converts a LoxValue that is a number
	to a int64.
	*/
	if inter.isFloat(value) {
		return int64(value.(float64))
	} else {
		return value.(int64)
	}
}

func (inter *Interpreter) convertNumsToFloat(value1 LoxValue, value2 LoxValue) (float64, float64) {
	/*Converts two LoxValues that are numbers
	to float64 numbers.
	*/
	return inter.convertNumToFloat(value1), inter.convertNumToFloat(value2)
}

func (inter *Interpreter) convertNumToFloat(value LoxValue) float64 {
	/*Converts a LoxValue that is a number
	to a float64.
	*/
	if inter.isInt(value) {
		return float64(value.(int64))
	} else {
		return value.(float64)
	}
}

func (inter *Interpreter) isInt(value LoxValue) bool {
	/*Determines if the passed number
	is an integer or not.
	*/
	_, isInt := value.(int64)
	return isInt
}

func (inter *Interpreter) isFloat(value LoxValue) bool {
	/*Determines if a value is
	a float number.
	*/
	_, isFloat := value.(float64)
	return isFloat
}

func (inter *Interpreter) isString(value LoxValue) bool {
	/*Determines if a value is
	a string.
	*/
	_, isString := value.(string)
	return isString
}

func (inter *Interpreter) isUserFunc(value LoxValue) bool {
	/*Determines if a value is
	a function.
	*/
	_, isFunc := value.(LoxFunction)
	return isFunc
}

func (inter *Interpreter) isClockFunction(value LoxValue) bool {
	/*Determines if a value is
	a clock built-in function.
	*/
	_, isClockFunc := value.(ClockFunction)
	return isClockFunc
}

func (inter *Interpreter) isToStringFunction(value LoxValue) bool {
	/*Determines if a value is
	a toString built-in function.
	*/
	_, isToString := value.(ToStringFunction)
	return isToString
}

func (inter *Interpreter) isInputFunction(value LoxValue) bool {
	/*Determines if a value is a
	input built-in function.
	*/
	_, isInput := value.(InputFunction)
	return isInput
}

func (inter *Interpreter) isParseStringFunction(value LoxValue) bool {
	/*Determines if a value is a
	parseString built-in function.
	*/
	_, isParseFunc := value.(ParseFunction)
	return isParseFunc
}

func (inter *Interpreter) isInstanceFunction(value LoxValue) bool {
	/*Determines if a value is a
	isInstance built-in function.
	*/
	_, isInstanceFunc := value.(IsInstanceFunction)
	return isInstanceFunc
}

func (inter *Interpreter) isEqual(left LoxValue, right LoxValue) bool {
	/*Checks if left and right are
	equal according to the rules
	of the lox language.
	*/
	if left == nil && right == nil {
		return true
	} else if left == nil {
		return false
	} else if inter.isInt(left) && inter.isInt(right) {
		return left.(int64) == right.(int64)
	} else if inter.isString(left) && inter.isString(right) {
		return left.(string) == right.(string)
	} else if inter.isNumber(left) && inter.isNumber(right) {
		return inter.convertNumToFloat(left) == inter.convertNumToFloat(right)
	} else if inter.isUserFunc(left) && inter.isUserFunc(right) {
		return left.(LoxFunction).declaration.name.lexeme == right.(LoxFunction).declaration.name.lexeme
	} else if inter.isClockFunction(left) && inter.isClockFunction(right) {
		return true
	} else if inter.isToStringFunction(left) && inter.isToStringFunction(right) {
		return true
	} else if inter.isParseStringFunction(left) && inter.isParseStringFunction(right) {
		return true
	} else {
		return false
	}
}
