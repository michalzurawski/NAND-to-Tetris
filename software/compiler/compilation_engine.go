package main

import (
	"fmt"
	"os"
	"strconv"
)

type CompilationEngine struct {
	file      *os.File
	tokenizer *Tokenizer
	open      int
}

// Creates new CompilationEngine
func NewCompilationEngine(fileName string) *CompilationEngine {
	tokenizer := NewTokenizer(fileName + ".jack")
	file, err := os.Create(fileName + ".vm")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not save file %s", fileName)
		os.Exit(1)
	}

	return &CompilationEngine{file: file, tokenizer: tokenizer}
}

// Closes the file
func (compilationEngine *CompilationEngine) Close() error {
	compilationEngine.tokenizer.Close()
	return compilationEngine.file.Close()
}

func (compilationEngine *CompilationEngine) CompileClass() {
	compilationEngine.writeOpen("<class>\n")
	compilationEngine.tokenizer.Advance()
	compilationEngine.eatKeyword(CLASS)
	compilationEngine.eatIdentifier()
	compilationEngine.eatSymbol(LEFT_CURLY)
	for compilationEngine.isKeyword(STATIC, FIELD) {
		compilationEngine.compileClassVariableDeclaration()
	}
	for compilationEngine.isKeyword(CONSTRUCTOR, METHOD, FUNCTION) {
		compilationEngine.compileSubroutine()
	}
	compilationEngine.eatSymbol(RIGHT_CURLY)
	compilationEngine.writeClose("</class>\n")
}

func (compilationEngine *CompilationEngine) compileClassVariableDeclaration() {
	compilationEngine.writeOpen("<classVarDec>\n")
	compilationEngine.eatKeyword(STATIC, FIELD)
	compilationEngine.compileType()
	compilationEngine.eatIdentifier()
	for compilationEngine.isSymbol(COMMA) {
		compilationEngine.eatSymbol(COMMA)
		compilationEngine.eatIdentifier()
	}
	compilationEngine.eatSymbol(SEMICOLON)
	compilationEngine.writeClose("</classVarDec>\n")
}

func (compilationEngine *CompilationEngine) compileSubroutine() {
	compilationEngine.writeOpen("<subroutineDec>\n")
	compilationEngine.eatKeyword(CONSTRUCTOR, METHOD, FUNCTION)
	if !compilationEngine.compileType() {
		compilationEngine.eatKeyword(VOID)
	}
	compilationEngine.eatIdentifier()
	compilationEngine.eatSymbol(LEFT_PARANTHESIS)
	compilationEngine.compileParameterList()
	compilationEngine.eatSymbol(RIGHT_PARANTHESIS)
	compilationEngine.compileSubroutineBody()
	compilationEngine.writeClose("</subroutineDec>\n")
}

func (compilationEngine *CompilationEngine) compileParameterList() {
	compilationEngine.writeOpen("<parameterList>\n")
	hasParameter := compilationEngine.compileType()
	if hasParameter {
		compilationEngine.eatIdentifier()
		for compilationEngine.isSymbol(COMMA) {
			compilationEngine.eatSymbol(COMMA)
			compilationEngine.compileType()
			compilationEngine.eatIdentifier()
		}
	}
	compilationEngine.writeClose("</parameterList>\n")
}

func (compilationEngine *CompilationEngine) compileSubroutineBody() {
	compilationEngine.writeOpen("<subroutineBody>\n")
	compilationEngine.eatSymbol(LEFT_CURLY)
	for compilationEngine.isKeyword(VAR) {
		compilationEngine.compileVariableDeclaration()
	}
	compilationEngine.compileStatements()
	compilationEngine.eatSymbol(RIGHT_CURLY)
	compilationEngine.writeClose("</subroutineBody>\n")
}

func (compilationEngine *CompilationEngine) compileVariableDeclaration() {
	compilationEngine.writeOpen("<varDec>\n")
	compilationEngine.eatKeyword(VAR)
	compilationEngine.compileType()
	compilationEngine.eatIdentifier()
	for compilationEngine.isSymbol(COMMA) {
		compilationEngine.eatSymbol(COMMA)
		compilationEngine.eatIdentifier()
	}
	compilationEngine.eatSymbol(SEMICOLON)
	compilationEngine.writeClose("</varDec>\n")
}

func (compilationEngine *CompilationEngine) compileStatements() {
	compilationEngine.writeOpen("<statements>\n")
	for compilationEngine.isKeyword(LET, IF, WHILE, DO, RETURN) {
		switch compilationEngine.tokenizer.GetKeyword() {
		case LET:
			compilationEngine.compileLet()
		case IF:
			compilationEngine.compileIf()
		case WHILE:
			compilationEngine.compileWhile()
		case DO:
			compilationEngine.compileDo()
		case RETURN:
			compilationEngine.compileReturn()
		}
	}
	compilationEngine.writeClose("</statements>\n")
}

func (compilationEngine *CompilationEngine) compileLet() {
	compilationEngine.writeOpen("<letStatement>\n")
	compilationEngine.eatKeyword(LET)
	compilationEngine.eatIdentifier()
	if compilationEngine.isSymbol(LEFT_BRACKET) {
		compilationEngine.eatSymbol(LEFT_BRACKET)
		compilationEngine.compileExpression()
		compilationEngine.eatSymbol(RIGHT_BRACKET)
	}
	compilationEngine.eatSymbol(EQUAL)
	compilationEngine.compileExpression()
	compilationEngine.eatSymbol(SEMICOLON)
	compilationEngine.writeClose("</letStatement>\n")
}

func (compilationEngine *CompilationEngine) compileIf() {
	compilationEngine.writeOpen("<ifStatement>\n")
	compilationEngine.eatKeyword(IF)
	compilationEngine.eatSymbol(LEFT_PARANTHESIS)
	compilationEngine.compileExpression()
	compilationEngine.eatSymbol(RIGHT_PARANTHESIS)
	compilationEngine.eatSymbol(LEFT_CURLY)
	compilationEngine.compileStatements()
	compilationEngine.eatSymbol(RIGHT_CURLY)
	if compilationEngine.isKeyword(ELSE) {
		compilationEngine.eatKeyword(ELSE)
		compilationEngine.eatSymbol(LEFT_CURLY)
		compilationEngine.compileStatements()
		compilationEngine.eatSymbol(RIGHT_CURLY)
	}
	compilationEngine.writeClose("</ifStatement>\n")
}

func (compilationEngine *CompilationEngine) compileWhile() {
	compilationEngine.writeOpen("<whileStatement>\n")
	compilationEngine.eatKeyword(WHILE)
	compilationEngine.eatSymbol(LEFT_PARANTHESIS)
	compilationEngine.compileExpression()
	compilationEngine.eatSymbol(RIGHT_PARANTHESIS)
	compilationEngine.eatSymbol(LEFT_CURLY)
	compilationEngine.compileStatements()
	compilationEngine.eatSymbol(RIGHT_CURLY)
	compilationEngine.writeClose("</whileStatement>\n")
}

func (compilationEngine *CompilationEngine) compileDo() {
	compilationEngine.writeOpen("<doStatement>\n")
	compilationEngine.eatKeyword(DO)
	compilationEngine.eatIdentifier()
	if compilationEngine.isSymbol(DOT) {
		compilationEngine.eatSymbol(DOT)
		compilationEngine.eatIdentifier()
	}
	compilationEngine.eatSymbol(LEFT_PARANTHESIS)
	compilationEngine.compileExpressionList()
	compilationEngine.eatSymbol(RIGHT_PARANTHESIS)
	compilationEngine.eatSymbol(SEMICOLON)
	compilationEngine.writeClose("</doStatement>\n")
}

func (compilationEngine *CompilationEngine) compileReturn() {
	compilationEngine.writeOpen("<returnStatement>\n")
	compilationEngine.eatKeyword(RETURN)
	if !compilationEngine.isSymbol(SEMICOLON) {
		compilationEngine.compileExpression()
	}
	compilationEngine.eatSymbol(SEMICOLON)
	compilationEngine.writeClose("</returnStatement>\n")
}

func (compilationEngine *CompilationEngine) compileExpression() {
	compilationEngine.writeOpen("<expression>\n")
	compilationEngine.compileTerm()
	for compilationEngine.isSymbol(PLUS, MINUS, MULTIPLY, DIVIDE, AND, OR, LESS, GREATER, EQUAL) {
		compilationEngine.eatSymbol(PLUS, MINUS, MULTIPLY, DIVIDE, AND, OR, LESS, GREATER, EQUAL)
		compilationEngine.compileTerm()
	}
	compilationEngine.writeClose("</expression>\n")
}

func (compilationEngine *CompilationEngine) compileTerm() {
	compilationEngine.writeOpen("<term>\n")
	if compilationEngine.isKeyword() {
		compilationEngine.eatKeyword()
	} else if compilationEngine.isIdentifier() {
		compilationEngine.eatIdentifier()
		if compilationEngine.isSymbol(LEFT_BRACKET) {
			compilationEngine.eatSymbol(LEFT_BRACKET)
			compilationEngine.compileExpression()
			compilationEngine.eatSymbol(RIGHT_BRACKET)
		} else if compilationEngine.isSymbol(LEFT_PARANTHESIS) {
			compilationEngine.eatSymbol(LEFT_PARANTHESIS)
			compilationEngine.compileExpressionList()
			compilationEngine.eatSymbol(RIGHT_PARANTHESIS)
		} else if compilationEngine.isSymbol(DOT) {
			compilationEngine.eatSymbol(DOT)
			compilationEngine.eatIdentifier()
			compilationEngine.eatSymbol(LEFT_PARANTHESIS)
			compilationEngine.compileExpressionList()
			compilationEngine.eatSymbol(RIGHT_PARANTHESIS)
		}
	} else if compilationEngine.tokenizer.GetTokenType() == STRING_CONST {
		compilationEngine.eatString()
	} else if compilationEngine.tokenizer.GetTokenType() == INT_CONST {
		compilationEngine.eatInteger()
	} else if compilationEngine.isSymbol(LEFT_PARANTHESIS) {
		compilationEngine.eatSymbol(LEFT_PARANTHESIS)
		compilationEngine.compileExpression()
		compilationEngine.eatSymbol(RIGHT_PARANTHESIS)
	} else if compilationEngine.isSymbol(MINUS, TILDE) {
		compilationEngine.eatSymbol(MINUS, TILDE)
		compilationEngine.compileTerm()
	}
	compilationEngine.writeClose("</term>\n")
}

func (compilationEngine *CompilationEngine) compileExpressionList() {
	compilationEngine.writeOpen("<expressionList>\n")
	if !compilationEngine.isSymbol() || compilationEngine.isSymbol(MINUS, TILDE, LEFT_PARANTHESIS) {
		compilationEngine.compileExpression()
		for compilationEngine.isSymbol(COMMA) {
			compilationEngine.eatSymbol(COMMA)
			compilationEngine.compileExpression()
		}
	}
	compilationEngine.writeClose("</expressionList>\n")
}

func (compilationEngine *CompilationEngine) compileType() bool {
	if compilationEngine.isKeyword(INT, CHAR, BOOLEAN) {
		compilationEngine.eatKeyword(INT, CHAR, BOOLEAN)
	} else if compilationEngine.isIdentifier() {
		compilationEngine.eatIdentifier()
	} else {
		return false
	}
	return true
}

func (compilationEngine *CompilationEngine) eatKeyword(keywords ...Keyword) {
	if !compilationEngine.isKeyword(keywords...) {
		compilationEngine.writeError("Expected keyword") // TODO: List expected keywords
	}
	keyword := compilationEngine.tokenizer.GetKeyword()
	for i := 0; i < compilationEngine.open; i++ {
		compilationEngine.file.WriteString(" ")
	}
	compilationEngine.file.WriteString("<keyword> " + keywordsToString[keyword] + " </keyword>\n")
	compilationEngine.tokenizer.Advance()
}

func (compilationEngine *CompilationEngine) isKeyword(keywords ...Keyword) bool {
	if compilationEngine.tokenizer.GetTokenType() != KEYWORD {
		return false
	} else if len(keywords) == 0 {
		return true
	}
	for _, keyword := range keywords {
		if compilationEngine.tokenizer.GetKeyword() == keyword {
			return true
		}
	}
	return false
}

func (compilationEngine *CompilationEngine) eatSymbol(symbols ...Symbol) {
	if !compilationEngine.isSymbol(symbols...) {
		fmt.Fprintf(os.Stderr, "At line %d\n", compilationEngine.tokenizer.GetTokenType())
		fmt.Fprintf(os.Stderr, "At line %s\n", compilationEngine.tokenizer.text)
		fmt.Fprintf(os.Stderr, "Expected symbol ")
		for _, symbol := range symbols {
			fmt.Fprintf(os.Stderr, "%s, ", symbolsToString[symbol])
		}
		fmt.Fprintf(os.Stderr, "\n")
		panic("")
	}
	symbol := compilationEngine.tokenizer.GetSymbol()
	for i := 0; i < compilationEngine.open; i++ {
		compilationEngine.file.WriteString(" ")
	}
	compilationEngine.file.WriteString("<symbol> " + symbolsToString[symbol] + " </symbol>\n")
	compilationEngine.tokenizer.Advance()
}

func (compilationEngine *CompilationEngine) isSymbol(symbols ...Symbol) bool {
	if compilationEngine.tokenizer.GetTokenType() != SYMBOL {
		return false
	}
	for _, symbol := range symbols {
		if compilationEngine.tokenizer.GetSymbol() == symbol {
			return true
		}
	}
	return len(symbols) == 0
}

func (compilationEngine *CompilationEngine) eatString() {
	if compilationEngine.tokenizer.GetTokenType() != STRING_CONST {
		compilationEngine.writeError("Expected string constant")
	}
	for i := 0; i < compilationEngine.open; i++ {
		compilationEngine.file.WriteString(" ")
	}
	compilationEngine.file.WriteString("<stringConstant> " + compilationEngine.tokenizer.GetStringValue() + " </stringConstant>\n")
	compilationEngine.tokenizer.Advance()
}

func (compilationEngine *CompilationEngine) eatInteger() {
	if compilationEngine.tokenizer.GetTokenType() != INT_CONST {
		compilationEngine.writeError("Expected integer constant")
	}
	for i := 0; i < compilationEngine.open; i++ {
		compilationEngine.file.WriteString(" ")
	}
	compilationEngine.file.WriteString("<integerConstant> " + strconv.Itoa(compilationEngine.tokenizer.GetIntegerValue()) + " </integerConstant>\n")
	compilationEngine.tokenizer.Advance()
}

func (compilationEngine *CompilationEngine) eatIdentifier() {
	if !compilationEngine.isIdentifier() {
		fmt.Fprintf(os.Stderr, "At line %d\n", compilationEngine.tokenizer.GetTokenType())
		fmt.Fprintf(os.Stderr, "At line %s\n", compilationEngine.tokenizer.text)
		compilationEngine.writeError("Expected identifier")
	}
	for i := 0; i < compilationEngine.open; i++ {
		compilationEngine.file.WriteString(" ")
	}
	compilationEngine.file.WriteString("<identifier> " + compilationEngine.tokenizer.GetIdentifier() + " </identifier>\n")
	compilationEngine.tokenizer.Advance()
}

func (compilationEngine *CompilationEngine) isIdentifier() bool {
	return compilationEngine.tokenizer.GetTokenType() == IDENTIFIER
}

func (compilationEngine *CompilationEngine) writeError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args)
	os.Exit(1)
}

func (compilationEngine *CompilationEngine) writeOpen(format string) {
	for i := 0; i < compilationEngine.open; i++ {
		compilationEngine.file.WriteString(" ")
	}
	compilationEngine.file.WriteString(format)
	compilationEngine.open += 2
}

func (compilationEngine *CompilationEngine) writeClose(format string) {
	compilationEngine.open -= 2
	for i := 0; i < compilationEngine.open; i++ {
		compilationEngine.file.WriteString(" ")
	}
	compilationEngine.file.WriteString(format)
}

var symbolsToString = map[Symbol]string{
	LEFT_PARANTHESIS:  "(",
	RIGHT_PARANTHESIS: ")",
	LEFT_BRACKET:      "[",
	RIGHT_BRACKET:     "]",
	LEFT_CURLY:        "{",
	RIGHT_CURLY:       "}",
	DOT:               ".",
	COMMA:             ",",
	SEMICOLON:         ";",
	PLUS:              "+",
	MINUS:             "-",
	MULTIPLY:          "*",
	DIVIDE:            "/",
	AND:               "&amp;",
	OR:                "|",
	LESS:              "&lt;",
	GREATER:           "&gt;",
	EQUAL:             "=",
	TILDE:             "~",
}

var keywordsToString = map[Keyword]string{
	CLASS:       "class",
	CONSTRUCTOR: "constructor",
	FUNCTION:    "function",
	METHOD:      "method",
	FIELD:       "field",
	STATIC:      "static",
	VAR:         "var",
	INT:         "int",
	CHAR:        "char",
	BOOLEAN:     "boolean",
	VOID:        "void",
	TRUE:        "true",
	FALSE:       "false",
	NULL:        "null",
	THIS:        "this",
	LET:         "let",
	DO:          "do",
	IF:          "if",
	ELSE:        "else",
	WHILE:       "while",
	RETURN:      "return",
}
