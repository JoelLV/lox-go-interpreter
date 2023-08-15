package main

import (
	"errors"
	"fmt"
)

type Scanner struct {
	currIndex int
	srcCode   []string
	loxError  bool
}

func (scnr *Scanner) init(sourceCode []string) {
	/*Initializes a new Scanner
	type.
	*/
	scnr.currIndex = 0
	scnr.srcCode = sourceCode
	scnr.loxError = false
}

func (scnr *Scanner) runScanner() []Token {
	/*Runs scanner and returns an
	array of tokens to be used by
	the parser.
	*/
	defer func() {
		//Sets the field loxError to true when an
		//exception is catched.
		if r := recover(); r != nil {
			fmt.Println(r)
			scnr.loxError = true
		}
	}()
	var tokenArr []Token = []Token{}
	for ; scnr.currIndex < len(scnr.srcCode); scnr.currIndex++ {
		//Combines two token arrays.
		tokenArr = append(tokenArr, scnr.getTokensInLine(scnr.getCurrLine())...)
	}
	tokenArr = append(tokenArr, Token{line: scnr.getLineNum(), tokenType: EOF, lexeme: "EOF"})

	return tokenArr
}

func (scnr *Scanner) getTokensInLine(line string) []Token {
	/*Returns an array of tokens representing
	one line of source code.
	*/
	var SINGLE_LEXEMES = map[string]TokenType{
		"(": LEFT_PAREN,
		")": RIGHT_PAREN,
		"{": LEFT_BRACE,
		"}": RIGHT_BRACE,
		",": COMMA,
		"-": MINUS,
		"+": PLUS,
		";": SEMICOLON,
		"*": STAR,
		"%": MOD,
	}
	var OPERATOR_LEXEMES = map[string]TokenType{
		"!=": NOT_EQUAL,
		"!":  NOT,
		"==": EQUAL_EQUAL,
		"=":  EQUAL,
		"<":  LESS_THAN,
		"<=": LESS_EQUAL,
		">":  GREATER_THAN,
		">=": GREATER_EQUAL,
	}
	var tokenArr []Token

	for lineIndex := 0; lineIndex < len(line); lineIndex++ {
		var tokenType TokenType
		lexeme := ""
		currChar := string(line[lineIndex])

		if scnr.inMap(SINGLE_LEXEMES, currChar) {
			tokenType = SINGLE_LEXEMES[currChar]
			lexeme = currChar
		} else if scnr.inMap(OPERATOR_LEXEMES, currChar) {
			nextChar := scnr.nextChar(line, lineIndex)
			if currChar == "!" {
				if nextChar == "=" {
					tokenType = OPERATOR_LEXEMES[currChar+nextChar]
					lexeme = "!="
					lineIndex++
				} else {
					tokenType = OPERATOR_LEXEMES[currChar]
					lexeme = "="
				}
			} else if currChar == "=" {
				if nextChar == "=" {
					tokenType = OPERATOR_LEXEMES[currChar+nextChar]
					lexeme = "=="
					lineIndex++
				} else {
					tokenType = OPERATOR_LEXEMES[currChar]
					lexeme = "="
				}
			} else if currChar == "<" {
				if nextChar == "=" {
					tokenType = OPERATOR_LEXEMES[currChar+nextChar]
					lexeme = "<="
					lineIndex++
				} else {
					tokenType = OPERATOR_LEXEMES[currChar]
					lexeme = "<"
				}
			} else if currChar == ">" {
				if nextChar == "=" {
					tokenType = OPERATOR_LEXEMES[currChar+nextChar]
					lexeme = ">="
					lineIndex++
				} else {
					tokenType = OPERATOR_LEXEMES[currChar]
					lexeme = ">"
				}
			}
		} else if scnr.isIgnorable(currChar) {
			continue
		} else if currChar == "/" {
			nextChar := scnr.nextChar(line, lineIndex)
			if nextChar == "/" {
				lineIndex++
				break
			} else {
				tokenType = SLASH
				lexeme = currChar
			}
		} else if currChar == "\"" {
			lastQuote, err := scnr.getLastQuoteIndex(line, lineIndex+1)
			if err != nil {
				scnr.error(scnr.getLineNum(), "Unterminated string.")
			} else {
				tokenType = STRING
				lexeme = line[lineIndex : lastQuote+1]
				lineIndex = lastQuote
			}
		} else if scnr.isDigit(currChar) {
			lastNum := scnr.getLastNumIndex(line, lineIndex+1)
			lexeme = line[lineIndex : lastNum+1]
			tokenType = NUMBER
			lineIndex = lastNum
		} else if scnr.isAlpha(currChar) {
			tokenType, lineIndex, lexeme = scnr.getIdentifierType(line, lineIndex)
		} else {
			scnr.error(scnr.getLineNum(), "Unknown Character."+currChar)
		}
		tokenArr = append(tokenArr, scnr.getToken(tokenType, scnr.getLineNum(), lexeme))
	}

	return tokenArr
}

func (scnr *Scanner) getCurrLine() string {
	/*Returns current string line
	according to the current index of the scanner.
	*/
	return scnr.srcCode[scnr.currIndex]
}

func (scnr *Scanner) getLineNum() int {
	/*Returns the index of scanner
	plus 1 representing the line
	of source code the scanner is scanning.
	*/
	return scnr.currIndex + 1
}

func (scnr *Scanner) getToken(tknType TokenType, line int, lexeme string) Token {
	/*Constructs a new
	token object with given values
	and returns it.
	*/
	return Token{tokenType: tknType, line: line, lexeme: lexeme}
}

func (scnr *Scanner) inMap(mapToSearch map[string]TokenType, currChar string) bool {
	/*Determines whether
	a character is in the given
	map.
	*/
	_, isInMap := mapToSearch[currChar]
	return isInMap
}

func (scnr *Scanner) nextChar(str string, index int) string {
	/*Returns next character in the string.
	Returns empty string if index is equal to the length - 1
	*/
	if index >= len(str)-1 {
		return ""
	} else {
		return string(str[index+1])
	}
}

func (scnr *Scanner) isIgnorable(char string) bool {
	/*Determines whether a given
	character is an ignorable character.
	*/
	var IGNORE_CHARS = []string{
		" ", "\r", "\t", "\n", "",
	}
	for i := 0; i < len(IGNORE_CHARS); i++ {
		if IGNORE_CHARS[i] == char {
			return true
		}
	}
	return false
}

func (scnr *Scanner) getLastQuoteIndex(line string, currIndex int) (int, error) {
	/*Determines the index of the second
	double quote in the string given the
	start index. Returns an error if it
	reaches the end of the string.
	*/
	for i := currIndex; i < len(line); i++ {
		if string(line[i]) == "\"" {
			return i, nil
		}
	}
	return 0, errors.New("Unterminated string")
}

func (scnr *Scanner) error(lineNum int, message string) {
	/*Creates an error message
	and then throws an exception when
	a lox error is found.
	*/
	scnr.loxError = true
	errMessage := fmt.Sprintf("[line %d] %s", lineNum, message)
	panic(errMessage)
}

func (scnr *Scanner) isDigit(char string) bool {
	/*Determines whether a given
	character is a digit.
	*/
	return char >= "0" && char <= "9"
}

func (scnr *Scanner) getLastNumIndex(line string, index int) int {
	/*Returns the index where
	the last digit was found.
	*/
	decimalFound := false
	for i := index; i < len(line); i++ {
		if scnr.isDigit(string(line[i])) {
			continue
		} else {
			if decimalFound && string(line[i]) == "." {
				scnr.error(scnr.currIndex, "Number cannot have two decimals.")
			}

			if i < len(line)-1 && string(line[i]) == "." && scnr.isDigit(string(line[i+1])) {
				decimalFound = true
				continue
			} else {
				return i - 1
			}
		}
	}
	return len(line) - 1
}

func (scnr *Scanner) isAlpha(char string) bool {
	/*Determines whether a given
	character is alphabetic.
	*/
	return (char >= "a" && char <= "z") || (char >= "A" && char <= "Z")
}

func (scnr *Scanner) getIdentifierType(line string, startIndex int) (TokenType, int, string) {
	/*Determines whether a word
	is a reserved word or identifier.
	Returns the appropriate token type, the lexeme,
	and the index where the last identifier
	character was found.
	*/
	var RESERVED_WORDS = map[string]TokenType{
		"and":    AND,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
		"const":  CONST,
	}
	i := startIndex
	for ; i < len(line); i++ {
		char := string(line[i])

		if !scnr.isAlpha(char) && !scnr.isDigit(char) {
			break
		}
	}
	word := line[startIndex:i]
	var tokenType TokenType

	if scnr.inMap(RESERVED_WORDS, word) {
		tokenType = RESERVED_WORDS[word]
	} else {
		tokenType = IDENTIFIER
	}

	return tokenType, i - 1, word
}
