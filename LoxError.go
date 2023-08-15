package main

import (
	"fmt"
	"os"
)

type LoxException struct {
	message string
	token   Token
}

type FunctionException struct {
	message string
}

func loxError(token Token, errMessage string, atEnd bool) LoxException {
	/*Displays error for user to handle.
	 */
	if atEnd {
		fmt.Printf("[line %d] Error at end: %s.\n", token.line-1, errMessage)
	} else {
		fmt.Printf("[line %d] Error at '%s': %s.\n", token.line, token.lexeme, errMessage)
	}
	return LoxException{message: errMessage, token: token}
}

func runtimeError(errMessage LoxException) {
	/*Handles runtime errors
	and displays them.
	*/
	fmt.Printf("%s\n[line %d] ", errMessage.message, errMessage.token.line)
	os.Exit(70)
}

func functionError(line int, message string) {
	/*Handles built in
	function errors and displays
	them.
	*/
	fmt.Printf("%s\n[line %d] ", message, line)
	os.Exit(70)
}
