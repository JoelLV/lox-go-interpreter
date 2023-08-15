package main

import (
	"fmt"
)

type Environment struct {
	values      map[string]LoxValue
	constValues map[string]bool
	enclosing   *Environment
}

func (env *Environment) init() {
	/*Initializes a new
	environment object.
	*/
	env.values = map[string]LoxValue{}
	env.constValues = map[string]bool{}
}

func (env *Environment) define(name string, value LoxValue) {
	/*Adds a new variable
	to the values dictionary.
	*/
	env.values[name] = value
}

func (env *Environment) get(name Token) LoxValue {
	/*Returns the value of values in the current
	environment given name.
	*/
	if env.varExists(name) {
		return env.values[name.lexeme]
	} else if env.enclosing != nil {
		return env.enclosing.get(name)
	} else {
		panic(LoxException{token: name, message: fmt.Sprintf("Undefined variable '%s'", name.lexeme)})
	}
}

func (env *Environment) varExists(name Token) bool {
	/*Determines if a variable exists in the environment.
	 */
	_, inMap := env.values[name.lexeme]
	return inMap
}

func (env *Environment) isConst(name Token) bool {
	/*Determines if a given variable is a constant
	variable with respect to this environment.
	*/
	_, inMap := env.constValues[name.lexeme]
	return inMap
}

func (env *Environment) assign(name Token, value LoxValue) {
	/*Reassigns existing variables in
	dictionary values with new value.
	*/
	if env.varExists(name) {
		env.values[name.lexeme] = value
	} else if env.enclosing != nil {
		env.enclosing.assign(name, value)
	} else {
		panic(LoxException{token: name, message: fmt.Sprintf("Undefined variable '%s'", name.lexeme)})
	}
}
