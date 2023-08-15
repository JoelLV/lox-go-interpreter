package main

import "fmt"

type LoxFunction struct {
	declaration Function
	closure     *Environment
}

func (loxFunc LoxFunction) call(interpreter Interpreter, arguments []LoxValue) (returnVal LoxValue) {
	/*Implements the method call
	from the interface LoxCallable.
	Creates a new environment relative
	to the function and executes the function.
	*/
	var env Environment
	env.init()
	env.enclosing = loxFunc.closure
	returnVal = nil

	for i := 0; i < len(loxFunc.declaration.params); i++ {
		env.define(loxFunc.declaration.params[i].lexeme, arguments[i])
	}
	defer func() {
		r := recover()
		if r != nil {
			if _, isFuncReturn := r.(RuntimeReturn); isFuncReturn {
				returnVal = r.(RuntimeReturn).value
			} else {
				panic(r)
			}
		}
	}()
	interpreter.executeBlock(loxFunc.declaration.body, &env)
	return returnVal
}

func (loxFunc LoxFunction) arity() int {
	/*Determines the arity of the function
	 */
	return len(loxFunc.declaration.params)
}

func (loxFunc LoxFunction) String() string {
	/*Returns a human readable string
	representing the object LoxFunction.
	*/
	return fmt.Sprintf("<fn %s>", loxFunc.declaration.name.lexeme)
}
