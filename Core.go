package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func runLexer(srcCode []string) ([]Token, bool) {
	/*Runs the scanner using
	the given source code
	and returns the resulting
	array of tokens.
	*/
	var scnr Scanner
	scnr.init(srcCode)
	return scnr.runScanner(), scnr.loxError
}

func runParser(tokenArr []Token) ([]Stmt, bool) {
	/*Runs the parser using
	the given token array
	and returns the resulting
	AST. Returns a boolean
	representing whether an error
	occured or not.
	*/
	var parser Parser
	parser.init(tokenArr)
	return parser.parseTokens(), parser.loxError
}

func runInterpreter(interpreter Interpreter) Interpreter {
	/*Runs the given interpreter
	and returns the modified interpreter.
	*/
	interpreter.interpret()
	return interpreter

}

func runFile(filePath string) {
	/*Reads all text lines of given
	path and transfers all the lines
	to a single string array.
	*/
	var srcCode []string
	fileName, err := os.Open(filePath)

	if err != nil {
		log.Fatalf("Failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(fileName)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		srcCode = append(srcCode, scanner.Text())
	}
	fileName.Close()

	tokenArr, scnrError := runLexer(srcCode)
	if !scnrError {
		stmtArr, parserError := runParser(tokenArr)
		if !parserError {
			var interpreter Interpreter
			interpreter.init(stmtArr)
			runInterpreter(interpreter)
		}
	}
}

func runRepl() {
	/*Prompts the user to enter
	code. Everytime a new line is
	written, it gets evaluated
	by the interpreter.
	*/
	var userInput string
	var srcCode []string
	var err error
	var interpreter Interpreter

	reader := bufio.NewReader(os.Stdin)
	interpreter.init(nil)
	for {
		fmt.Print("> ")
		userInput, err = reader.ReadString('\n')
		srcCode = []string{userInput}

		if err == io.EOF {
			fmt.Println()
			break
		}
		tokenArr, scnrError := runLexer(srcCode)
		if !scnrError {
			stmtArr, parserError := runParser(tokenArr)
			if !parserError {
				interpreter.trees = stmtArr
				interpreter = runInterpreter(interpreter)
			}
		}
	}
}

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: plox [script]")
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runRepl()
	}
}
