package main

import (
	"strconv"
)

type Parser struct {
	tokens   []Token
	index    int
	loxError bool
}

func (parser *Parser) init(tokenList []Token) {
	/*Initializes a parser
	object.
	*/
	parser.tokens = tokenList
	parser.index = 0
	parser.loxError = false
}

func (parser *Parser) parseTokens() []Stmt {
	/*Parses all tokens
	in the slice of tokens of the parser
	object and returns a slice of statement
	ASTs.
	*/
	defer func() {
		recover()
	}()
	var statements []Stmt
	for !parser.atEnd() {
		statements = append(statements, parser.declaration())
	}

	return statements
}

func (parser *Parser) declaration() Stmt {
	/*Representation of declaration
	as a grammar rule.
	*/
	defer func() {
		if r := recover(); r != nil {
			parser.synchronize()
		}
	}()

	if parser.matchAndAdvance(VAR) {
		return parser.varDeclaration(false)
	} else if parser.matchAndAdvance(CONST) {
		return parser.constDeclaration()
	} else if parser.matchAndAdvance(FUN) {
		return parser.function()
	} else {
		return parser.statement()
	}
}

func (parser *Parser) constDeclaration() Stmt {
	/*Used to set a flag for
	variable declarations.
	*/
	parser.consume(VAR, "Expect 'var' keyword after 'const'")
	return parser.varDeclaration(true)
}

func (parser *Parser) varDeclaration(isConstant bool) Stmt {
	/*Representation of variable declaration
	as a grammar rule.
	*/
	var initializer Expr = nil
	name := parser.consume(IDENTIFIER, "Expect variable name")

	if parser.matchAndAdvance(EQUAL) {
		initializer = parser.expression()
	}

	parser.consume(SEMICOLON, "Expect ';' after variable declaration")
	return Var{name: name, initializer: initializer, isConst: isConstant}
}

func (parser *Parser) function() Stmt {
	/*Representation of a function declaration
	as a grammar rule.
	*/
	name := parser.consume(IDENTIFIER, "Expect function name")
	parser.consume(LEFT_PAREN, "Expect '(' after function name")
	var parameters []Token

	if !parser.matchTokenType(RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				parser.compilerError(parser.getCurrentToken(), "Can't have more than 255 parameters")
			}
			parameters = append(parameters, parser.consume(IDENTIFIER, "Expect parameter name"))

			if !parser.matchAndAdvance(COMMA) {
				break
			}
		}
	}
	parser.consume(RIGHT_PAREN, "Expect ')' after parameters")

	parser.consume(LEFT_BRACE, "Expect '{' before function body")
	body := parser.block()
	return Function{name: name, params: parameters, body: body}
}

func (parser *Parser) statement() Stmt {
	/*Representation of a statement as a
	grammar rule.
	*/
	if parser.matchAndAdvance(PRINT) {
		return parser.printStatement()
	} else if parser.matchAndAdvance(LEFT_BRACE) {
		return Block{statements: parser.block()}
	} else if parser.matchAndAdvance(IF) {
		return parser.ifStatement()
	} else if parser.matchAndAdvance(WHILE) {
		return parser.whileStatement()
	} else if parser.matchAndAdvance(FOR) {
		return parser.forStatement()
	} else if parser.matchAndAdvance(RETURN) {
		return parser.returnStatement()
	} else {
		return parser.expressionStatement()
	}
}

func (parser *Parser) printStatement() Stmt {
	/*Representation of print statement
	as a grammar rule.
	*/
	value := parser.expression()
	parser.consume(SEMICOLON, "Expect ';' after value")

	return Print{expression: value}
}

func (parser *Parser) block() []Stmt {
	/*Representation of block statement
	as a grammar rule.
	*/
	var statements []Stmt

	for !parser.atEnd() && !parser.matchTokenType(RIGHT_BRACE) {
		statements = append(statements, parser.declaration())
	}

	parser.consume(RIGHT_BRACE, "Expect '}' after block")
	return statements
}

func (parser *Parser) ifStatement() Stmt {
	/*Representation of an if statement
	as a grammar rule.
	*/
	parser.consume(LEFT_PAREN, "Expect '(' after 'if'")
	condition := parser.expression()
	parser.consume(RIGHT_PAREN, "Expect ')' after if condition")

	thenBranch := parser.statement()
	var elseBranch Stmt = nil

	if parser.matchAndAdvance(ELSE) {
		elseBranch = parser.statement()
	}

	return If{condition: condition, thenBranch: thenBranch, elseBranch: elseBranch}
}

func (parser *Parser) whileStatement() Stmt {
	/*Representation of a while statement
	as a grammar rule.
	*/
	parser.consume(LEFT_PAREN, "Expect '(' after 'while'")
	condition := parser.expression()
	parser.consume(RIGHT_PAREN, "Expect ')' after condition")
	body := parser.statement()

	return While{condition: condition, body: body}
}

func (parser *Parser) forStatement() Stmt {
	/*Representation of a for loop statement
	as a grammar rule.
	*/
	parser.consume(LEFT_PAREN, "Expect '(' after 'for'")

	var initializer Stmt = nil
	if parser.matchAndAdvance(SEMICOLON) {
		initializer = nil
	} else if parser.matchAndAdvance(CONST) {
		initializer = parser.varDeclaration(true)
	} else if parser.matchAndAdvance(VAR) {
		initializer = parser.varDeclaration(false)
	} else {
		initializer = parser.expressionStatement()
	}

	var condition Expr = nil
	if !parser.atEnd() && !parser.matchTokenType(SEMICOLON) {
		condition = parser.expression()
	}
	parser.consume(SEMICOLON, "Expect ';' after loop condition")

	var increment Expr = nil
	if !parser.atEnd() && !parser.matchTokenType(RIGHT_PAREN) {
		increment = parser.expression()
	}
	parser.consume(RIGHT_PAREN, "Expect ')' after for clauses")

	body := parser.statement()

	if increment != nil {
		body = Block{statements: []Stmt{body, Expression{expression: increment}}}
	}
	if condition == nil {
		condition = Literal{value: true}
	}
	body = While{condition: condition, body: body}

	if initializer != nil {
		body = Block{statements: []Stmt{initializer, body}}
	}

	return body
}

func (parser *Parser) expressionStatement() Stmt {
	/*Representation of expression statement
	as a grammar rule.
	*/
	expr := parser.expression()
	parser.consume(SEMICOLON, "Expect ';' after expression")

	return Expression{expression: expr}
}

func (parser *Parser) returnStatement() Stmt {
	/*Representation of a return statement as a grammar
	rule.
	*/
	keyword := parser.previousToken()
	var value Expr = nil

	if !parser.matchTokenType(SEMICOLON) {
		value = parser.expression()
	}

	parser.consume(SEMICOLON, "Expect ';' after return value")
	return Return{keyword: keyword, value: value}
}

func (parser *Parser) expression() Expr {
	/*Representation of expression
	as a grammar rule.
	*/
	return parser.assignment()
}

func (parser *Parser) assignment() Expr {
	/*Representation of assignment
	as a grammar rule.
	*/
	expr := parser.orExpression()

	if parser.matchAndAdvance(EQUAL) {
		equals := parser.previousToken()
		value := parser.assignment()

		if _, isVar := expr.(Variable); isVar {
			name := expr.(Variable).name
			return Assign{name: name, value: value}
		} else {
			panic(parser.compilerError(equals, "Invalid assignment target"))
		}
	} else {
		return expr
	}
}

func (parser *Parser) orExpression() Expr {
	/*Representation of or expression
	as a grammar rule.
	*/
	expr := parser.andExpression()

	for parser.matchAndAdvance(OR) {
		operator := parser.previousToken()
		right := parser.andExpression()
		expr = Logical{left: expr, operator: operator, right: right}
	}

	return expr
}

func (parser *Parser) andExpression() Expr {
	/*Representation of an and expression
	as a grammar rule.
	*/
	expr := parser.equality()

	for parser.matchAndAdvance(AND) {
		operator := parser.previousToken()
		right := parser.equality()
		expr = Logical{left: expr, operator: operator, right: right}
	}

	return expr
}

func (parser *Parser) equality() Expr {
	/*Representation of equality
	as a grammar rule.
	*/
	expr := parser.comparison()

	for parser.matchAndAdvance(NOT_EQUAL) || parser.matchAndAdvance(EQUAL_EQUAL) {
		operator := parser.previousToken()
		rightSide := parser.comparison()
		expr = Binary{left: expr, operator: operator, right: rightSide}
	}

	return expr
}

func (parser *Parser) comparison() Expr {
	/*Representation of comparison as
	a grammar rule.
	*/
	expr := parser.term()

	for parser.matchAndAdvance(GREATER_EQUAL) || parser.matchAndAdvance(GREATER_THAN) ||
		parser.matchAndAdvance(LESS_EQUAL) || parser.matchAndAdvance(LESS_THAN) {
		operator := parser.previousToken()
		right := parser.term()
		expr = Binary{left: expr, operator: operator, right: right}
	}

	return expr
}

func (parser *Parser) term() Expr {
	/*Representation of term
	as a grammar rule.
	*/
	expr := parser.factor()

	for parser.matchAndAdvance(MINUS) || parser.matchAndAdvance(PLUS) {
		operator := parser.previousToken()
		right := parser.factor()
		expr = Binary{left: expr, operator: operator, right: right}
	}

	return expr
}

func (parser *Parser) factor() Expr {
	/*Representation of factor
	as a grammar rule/
	*/
	expr := parser.unary()

	for parser.matchAndAdvance(SLASH) || parser.matchAndAdvance(STAR) || parser.matchAndAdvance(MOD) {
		operator := parser.previousToken()
		right := parser.unary()
		expr = Binary{left: expr, operator: operator, right: right}
	}

	return expr
}

func (parser *Parser) unary() Expr {
	/*Representation of unary as
	a grammar rule.
	*/
	if parser.matchAndAdvance(NOT) || parser.matchAndAdvance(MINUS) {
		operator := parser.previousToken()
		right := parser.unary()
		return Unary{operator: operator, right: right}
	}

	return parser.call()
}

func (parser *Parser) call() Expr {
	/*Representation of a function call
	as a grammar rule.
	*/
	expr := parser.primary()

	for {
		if parser.matchAndAdvance(LEFT_PAREN) {
			expr = parser.finishCall(expr)
		} else {
			break
		}
	}
	return expr
}

func (parser *Parser) finishCall(callee Expr) Expr {
	/*Parses the arguments inside
	the function call
	*/
	var arguments []Expr

	if !parser.matchTokenType(RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				parser.compilerError(parser.getCurrentToken(), "Can't have more than 255 arguments")
			}
			arguments = append(arguments, parser.expression())
			if !parser.matchAndAdvance(COMMA) {
				break
			}
		}
	}
	closingParen := parser.consume(RIGHT_PAREN, "Expect ')' after arguments")
	return Call{callee: callee, paren: closingParen, arguments: arguments}
}

func (parser *Parser) primary() Expr {
	/*Representation of primary as
	a grammar rule.
	*/
	if parser.matchAndAdvance(FALSE) {
		return Literal{value: false}
	} else if parser.matchAndAdvance(TRUE) {
		return Literal{value: true}
	} else if parser.matchAndAdvance(NIL) {
		return Literal{value: nil}
	} else if parser.matchAndAdvance(STRING) {
		return Literal{value: parser.previousToken().lexeme}
	} else if parser.matchAndAdvance(NUMBER) {
		literalInt, err := strconv.ParseInt(parser.previousToken().lexeme, 10, 64)
		if err == nil {
			return Literal{value: literalInt}
		} else {
			literalFloat, _ := strconv.ParseFloat(parser.previousToken().lexeme, 64)
			return Literal{value: literalFloat}
		}
	} else if parser.matchAndAdvance(LEFT_PAREN) {
		expr := parser.expression()
		parser.consume(RIGHT_PAREN, "Expect ')' after expression")
		return Grouping{expression: expr}
	} else if parser.matchAndAdvance(IDENTIFIER) {
		return Variable{name: parser.previousToken()}
	} else {
		panic(parser.compilerError(parser.getCurrentToken(), "Expect expression"))
	}
}

func (parser *Parser) getCurrentToken() Token {
	/*Returns current token
	according to the index field
	of the parser.
	*/
	return parser.tokens[parser.index]
}

func (parser *Parser) getCurrentTokenType() TokenType {
	/*Returns the token type
	of the current token according
	to the index field of the parser.
	*/
	return parser.getCurrentToken().tokenType
}

func (parser *Parser) matchAndAdvance(tokenType TokenType) bool {
	/*Returns true and advances the parser index
	if the current token matches the token parameter.
	Otherwise, it returns false.
	*/
	if parser.matchTokenType(tokenType) {
		parser.advanceIndex()
		return true
	} else {
		return false
	}
}

func (parser *Parser) matchTokenType(tokenType TokenType) bool {
	/*Returns true if the token type passed matches
	the current token in parser.
	*/
	return parser.getCurrentTokenType() == tokenType
}

func (parser *Parser) advanceIndex() {
	/*Advances the parser index.
	 */
	parser.index++
}

func (parser *Parser) consume(tokenType TokenType, errMessage string) Token {
	/*Consumes the expected token according
	to the passed token type. Otherwise it
	throws the passed error message.
	*/
	if !parser.atEnd() && parser.matchTokenType(tokenType) {
		parser.advanceIndex()
		return parser.previousToken()
	} else {
		panic(parser.compilerError(parser.getCurrentToken(), errMessage))
	}
}

func (parser *Parser) compilerError(token Token, message string) LoxException {
	/*Sets parser field loxError to true and
	prints the error message. Returns a lox
	error type.
	*/
	parser.loxError = true
	return loxError(token, message, parser.atEnd())
}

func (parser *Parser) previousToken() Token {
	/*Returns the previous token of the parser.
	 */
	return parser.tokens[parser.index-1]
}

func (parser *Parser) atEnd() bool {
	/*Returns true if the current
	token is EOF. Otherwise it
	returns false.
	*/
	return parser.getCurrentTokenType() == EOF
}

func (parser *Parser) synchronize() bool {
	/*Synchronizes the parser
	after an error has been detected.
	*/
	if !parser.atEnd() {
		parser.index++
	}

	for !parser.atEnd() {
		if parser.previousToken().tokenType == SEMICOLON {
			break
		} else {
			tokenType := parser.getCurrentTokenType()

			switch tokenType {
			case FUN:
				return true
			case VAR:
				return true
			case FOR:
				return true
			case IF:
				return true
			case WHILE:
				return true
			case PRINT:
				return true
			case RETURN:
				return true
			}

			parser.index++
		}
	}
	return true
}
