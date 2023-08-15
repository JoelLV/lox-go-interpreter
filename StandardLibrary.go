package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type ClockFunction struct {
	declaration Function
	closure     *Environment
}

func (clockFunc ClockFunction) arity() int {
	/*Returns the arity of a
	ClockFunction type.
	*/
	return 0
}

func (clockFunc ClockFunction) call(interpreter Interpreter, arguments []LoxValue) LoxValue {
	/*Returns the time represented as the
	number of seconds since epoch.
	*/
	now := time.Now()
	return now.Unix()
}

func (clockFunc ClockFunction) String() string {
	/*Returns a string representation
	of a clock function in lox.
	*/
	return "<native fn>"
}

type ToStringFunction struct {
	declaration Function
	closure     *Environment
}

func (toStringFunc ToStringFunction) arity() int {
	/*Returns the arity of a ToStringFunction
	type.
	*/
	return 1
}

func (toStringFunc ToStringFunction) call(interpreter Interpreter, arguments []LoxValue) LoxValue {
	/*Converts any lox type to its
	string representation.
	*/
	if integer, isInt := arguments[0].(int64); isInt {
		return fmt.Sprintf("%d", integer)
	} else if float, isFloat := arguments[0].(float64); isFloat {
		return fmt.Sprintf("%f", float)
	} else if boolean, isBool := arguments[0].(bool); isBool {
		if boolean {
			return "true"
		} else {
			return "false"
		}
	} else if userFunc, isUserFunc := arguments[0].(LoxFunction); isUserFunc {
		return userFunc.String()
	} else if clockFunc, isClockFunc := arguments[0].(ClockFunction); isClockFunc {
		return clockFunc.String()
	} else if stringFunc, isStringFunc := arguments[0].(ToStringFunction); isStringFunc {
		return stringFunc.String()
	} else if inputFunc, isInputFunc := arguments[0].(InputFunction); isInputFunc {
		return inputFunc.String()
	} else if parseFunc, isParseFunc := arguments[0].(ParseFunction); isParseFunc {
		return parseFunc.String()
	} else if instanceFunc, isInstanceFunc := arguments[0].(IsInstanceFunction); isInstanceFunc {
		return instanceFunc.String()
	} else {
		return arguments[0]
	}
}

func (toStringFunc ToStringFunction) String() string {
	/*Returns a human readable string
	representing the object toString.
	*/
	return "<native fn>"
}

type InputFunction struct {
	declaration Function
	closure     *Environment
}

func (inputFunc InputFunction) arity() int {
	/*Determines the arity of a
	input built-in function.
	*/
	return 0
}

func (inputFunc InputFunction) call(Interpreter Interpreter, arguments []LoxValue) LoxValue {
	/*Prompts user for input.
	 */
	reader := bufio.NewReader(os.Stdin)
	userInput, err := reader.ReadString('\n')
	userInput = strings.Replace(userInput, "\n", "", -1)
	userInput = strings.Replace(userInput, "\"", "", -1)
	if err == nil {
		return userInput
	} else {
		return ""
	}
}

func (inputFunc InputFunction) String() string {
	/*Returns a human readable string
	representing the object input.
	*/
	return "<native fn>"
}

type ParseFunction struct {
	declaration Function
	closure     *Environment
}

func (parseFunc ParseFunction) arity() int {
	/*Determines the arity of a
	parseString built-in function.
	*/
	return 2
}

func (parseFunc ParseFunction) call(interpreter Interpreter, arguments []LoxValue) LoxValue {
	/*Parses a string
	given the lox type.
	*/
	typeStr, typeIsString := arguments[0].(string)
	valueStr, valueIsString := arguments[1].(string)

	typeStr = strings.Replace(typeStr, "\"", "", -1)
	valueStr = strings.Replace(valueStr, "\"", "", -1)

	if typeIsString && valueIsString {
		switch typeStr {
		case "int":
			{
				valueInt, err := strconv.ParseInt(valueStr, 10, 64)
				if err == nil {
					return valueInt
				} else {
					panic(FunctionException{message: fmt.Sprintf("Cannot convert '%s' to int.", valueStr)})
				}
			}
		case "float":
			{
				valueFloat, err := strconv.ParseFloat(valueStr, 64)
				if err == nil {
					return valueFloat
				} else {
					panic(FunctionException{message: fmt.Sprintf("Cannot convert '%s' to float.", valueStr)})
				}
			}
		case "bool":
			{
				valueBool, err := strconv.ParseBool(valueStr)
				if err == nil {
					return valueBool
				} else {
					panic(FunctionException{message: fmt.Sprintf("Cannot convert '%s' to boolean.", valueStr)})
				}
			}
		case "string":
			{
				return valueStr
			}
		default:
			panic(FunctionException{message: fmt.Sprintf("Type '%s' is not supported.", typeStr)})
		}
	} else {
		panic(FunctionException{message: fmt.Sprintf("Arguments are not strings.")})
	}
}

func (parseFunc ParseFunction) String() string {
	/*Returns a human readable string
	representation of the built-in
	function parseString
	*/
	return "<native fn>"
}

type IsInstanceFunction struct {
	declaration Function
	closure     *Environment
}

func (isInstanceFunc IsInstanceFunction) arity() int {
	/*Determines the arity of
	built-in function isInstance.
	*/
	return 2
}

func (isInstanceFunc IsInstanceFunction) call(interpreter Interpreter, arguments []LoxValue) LoxValue {
	/*Determines whether a
	given type is the given value.
	*/
	typeStr, typeIsString := arguments[0].(string)

	typeStr = strings.Replace(typeStr, "\"", "", -1)
	if typeIsString {
		switch typeStr {
		case "int":
			{
				_, valueIsInt := arguments[1].(int64)
				return valueIsInt
			}
		case "float":
			{
				_, valueIsFloat := arguments[1].(float64)
				return valueIsFloat
			}
		case "boolean":
			{
				_, valueIsBool := arguments[1].(bool)
				return valueIsBool
			}
		case "string":
			{
				_, valueIsString := arguments[1].(string)
				return valueIsString
			}
		case "function":
			{
				_, valueIsFunc := arguments[1].(LoxCallable)
				return valueIsFunc
			}
		default:
			{
				panic(FunctionException{message: fmt.Sprintf("Type '%s' is not supported.", typeStr)})
			}
		}
	} else {
		panic(FunctionException{message: "Type argument must be string."})
	}
}

func (isInstanceFunc IsInstanceFunction) String() string {
	/*Returns the string representation
	of the built-in function isInstance
	*/
	return "<native fn>"
}
