package main

import (
	"fmt"
	"os"
	"strconv"
)

type CompilationEngine struct {
	vmWriter     *VMWriter
	tokenizer    *Tokenizer
	symbolTable  *SymbolTable
	className    string
	counterWhile int
	counterIf    int
}

// Creates new CompilationEngine
func NewCompilationEngine(fileName string) *CompilationEngine {
	tokenizer := NewTokenizer(fileName + ".jack")
	symbolTable := NewSymbolTable()
	vmWriter := NewVMWriter(fileName + ".vm")

	return &CompilationEngine{vmWriter: vmWriter, tokenizer: tokenizer, symbolTable: symbolTable}
}

// Closes the file
func (compilationEngine *CompilationEngine) Close() error {
	compilationEngine.tokenizer.Close()
	return compilationEngine.vmWriter.Close()
}

func (compilationEngine *CompilationEngine) CompileClass() {
	compilationEngine.tokenizer.Advance()
	compilationEngine.eatKeyword(CLASS)
	compilationEngine.className = compilationEngine.eatIdentifier() + "."
	compilationEngine.eatSymbol(LEFT_CURLY)
	for compilationEngine.isKeyword(STATIC, FIELD) {
		compilationEngine.compileClassVariableDeclaration()
	}
	for compilationEngine.isKeyword(CONSTRUCTOR, METHOD, FUNCTION) {
		compilationEngine.compileSubroutine()
	}
	compilationEngine.eatSymbol(RIGHT_CURLY)
}

func (compilationEngine *CompilationEngine) compileClassVariableDeclaration() {
	kind := compilationEngine.eatKeyword(STATIC, FIELD)
	_, typeVar := compilationEngine.compileType()
	compilationEngine.eatIdentifierDefinition(typeVar, kind)
	for compilationEngine.isSymbol(COMMA) {
		compilationEngine.eatSymbol(COMMA)
		compilationEngine.eatIdentifierDefinition(typeVar, kind)
	}
	compilationEngine.eatSymbol(SEMICOLON)
}

func (compilationEngine *CompilationEngine) compileSubroutine() {
	compilationEngine.symbolTable.StartSubroutine()
	compilationEngine.counterIf = 0
	compilationEngine.counterWhile = 0
	functionType := compilationEngine.tokenizer.GetKeyword()
	if functionType == METHOD {
		compilationEngine.symbolTable.Define("this", compilationEngine.className, ARG)
	}
	compilationEngine.eatKeyword(CONSTRUCTOR, METHOD, FUNCTION)
	if isType, _ := compilationEngine.compileType(); !isType {
		compilationEngine.eatKeyword(VOID)
	}
	name := compilationEngine.eatIdentifier()
	compilationEngine.eatSymbol(LEFT_PARANTHESIS)
	compilationEngine.compileParameterList()
	compilationEngine.eatSymbol(RIGHT_PARANTHESIS)
	compilationEngine.compileSubroutineBody(name, functionType)
}

func (compilationEngine *CompilationEngine) compileParameterList() {
	hasParameter, varType := compilationEngine.compileType()
	if hasParameter {
		compilationEngine.eatIdentifierDefinition(varType, ARG)
		for compilationEngine.isSymbol(COMMA) {
			compilationEngine.eatSymbol(COMMA)
			_, varType = compilationEngine.compileType()
			compilationEngine.eatIdentifierDefinition(varType, ARG)
		}
	}
}

func (compilationEngine *CompilationEngine) compileSubroutineBody(name string, functionType Keyword) {
	compilationEngine.eatSymbol(LEFT_CURLY)
	for compilationEngine.isKeyword(VAR) {
		compilationEngine.compileVariableDeclaration()
	}
	compilationEngine.vmWriter.WriteFunction(compilationEngine.className+name, compilationEngine.symbolTable.variableCounter)
	if functionType == CONSTRUCTOR {
		classSize := compilationEngine.symbolTable.fieldCounter
		compilationEngine.vmWriter.WritePush(CONST, classSize)
		compilationEngine.vmWriter.WriteCall("Memory.alloc", 1)
		compilationEngine.vmWriter.WritePop(POINTER, 0)
	} else if functionType == METHOD {
		compilationEngine.vmWriter.WritePush(ARG, 0)
		compilationEngine.vmWriter.WritePop(POINTER, 0)
	}
	compilationEngine.compileStatements()
	compilationEngine.eatSymbol(RIGHT_CURLY)
}

func (compilationEngine *CompilationEngine) compileVariableDeclaration() {
	compilationEngine.eatKeyword(VAR)
	_, typVar := compilationEngine.compileType()
	compilationEngine.eatIdentifierDefinition(typVar, VAR)
	for compilationEngine.isSymbol(COMMA) {
		compilationEngine.eatSymbol(COMMA)
		compilationEngine.eatIdentifierDefinition(typVar, VAR)
	}
	compilationEngine.eatSymbol(SEMICOLON)
}

func (compilationEngine *CompilationEngine) compileStatements() {
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
}

func (compilationEngine *CompilationEngine) compileLet() {
	compilationEngine.eatKeyword(LET)
	name := compilationEngine.eatIdentifier()
	isArray := false
	if compilationEngine.isSymbol(LEFT_BRACKET) {
		compilationEngine.eatSymbol(LEFT_BRACKET)
		compilationEngine.compileExpression()
		compilationEngine.eatSymbol(RIGHT_BRACKET)
		compilationEngine.vmWriter.WritePush(compilationEngine.symbolTable.GetVariableInfo(name))
		compilationEngine.vmWriter.WriteArithmetic(PLUS)
		isArray = true
	}
	compilationEngine.eatSymbol(EQUAL)
	compilationEngine.compileExpression()
	compilationEngine.eatSymbol(SEMICOLON)
	if isArray {
		compilationEngine.vmWriter.WritePop(TEMP, 0)
		compilationEngine.vmWriter.WritePop(POINTER, 1)
		compilationEngine.vmWriter.WritePush(TEMP, 0)
		compilationEngine.vmWriter.WritePop(THAT, 0)
	} else {
		compilationEngine.vmWriter.WritePop(compilationEngine.symbolTable.GetVariableInfo(name))
	}
}

func (compilationEngine *CompilationEngine) compileIf() {
	compilationEngine.counterIf++
	compilationEngine.eatKeyword(IF)
	compilationEngine.eatSymbol(LEFT_PARANTHESIS)
	compilationEngine.compileExpression()
	compilationEngine.eatSymbol(RIGHT_PARANTHESIS)
	label := compilationEngine.getLabelIfTrue()
	labelFalse := compilationEngine.getLabelIfFalse()
	compilationEngine.vmWriter.WriteIf(label)
	compilationEngine.vmWriter.WriteGoto(labelFalse)
	compilationEngine.vmWriter.WriteLabel(label)
	labelEnd := compilationEngine.getLabelIfEnd()
	compilationEngine.eatSymbol(LEFT_CURLY)
	compilationEngine.compileStatements()
	compilationEngine.eatSymbol(RIGHT_CURLY)
	if compilationEngine.isKeyword(ELSE) {
		compilationEngine.vmWriter.WriteGoto(labelEnd)
		compilationEngine.vmWriter.WriteLabel(labelFalse)
		compilationEngine.eatKeyword(ELSE)
		compilationEngine.eatSymbol(LEFT_CURLY)
		compilationEngine.compileStatements()
		compilationEngine.eatSymbol(RIGHT_CURLY)
		compilationEngine.vmWriter.WriteLabel(labelEnd)
	} else {
		compilationEngine.vmWriter.WriteLabel(labelFalse)
	}
}

func (compilationEngine *CompilationEngine) compileWhile() {
	compilationEngine.counterWhile++
	label := compilationEngine.getLabelWhile()
	compilationEngine.vmWriter.WriteLabel(label)
	compilationEngine.eatKeyword(WHILE)
	compilationEngine.eatSymbol(LEFT_PARANTHESIS)
	compilationEngine.compileExpression()
	compilationEngine.eatSymbol(RIGHT_PARANTHESIS)
	compilationEngine.vmWriter.WriteArithmetic(NOT)
	labelEnd := compilationEngine.getLabelWhileEnd()
	compilationEngine.vmWriter.WriteIf(labelEnd)
	compilationEngine.eatSymbol(LEFT_CURLY)
	compilationEngine.compileStatements()
	compilationEngine.eatSymbol(RIGHT_CURLY)
	compilationEngine.vmWriter.WriteGoto(label)
	compilationEngine.vmWriter.WriteLabel(labelEnd)
}

func (compilationEngine *CompilationEngine) compileDo() {
	compilationEngine.eatKeyword(DO)
	name := compilationEngine.eatIdentifier()
	count := 0
	if compilationEngine.isSymbol(DOT) {
		compilationEngine.eatSymbol(DOT)
		if compilationEngine.symbolTable.HasVariable(name) {
			compilationEngine.vmWriter.WritePush(compilationEngine.symbolTable.GetVariableInfo(name))
			name = compilationEngine.symbolTable.GetVariableType(name)
			count = 1
		}
		name += "." + compilationEngine.eatIdentifier()
	} else {
		compilationEngine.vmWriter.WritePush(POINTER, 0)
		name = compilationEngine.className + name
		count = 1

	}
	compilationEngine.eatSymbol(LEFT_PARANTHESIS)
	count += compilationEngine.compileExpressionList()
	compilationEngine.eatSymbol(RIGHT_PARANTHESIS)
	compilationEngine.eatSymbol(SEMICOLON)
	compilationEngine.vmWriter.WriteCall(name, count)
	compilationEngine.vmWriter.WritePop(TEMP, 0)
}

func (compilationEngine *CompilationEngine) compileReturn() {
	compilationEngine.eatKeyword(RETURN)
	if !compilationEngine.isSymbol(SEMICOLON) {
		compilationEngine.compileExpression()
	} else {
		compilationEngine.vmWriter.WritePush(CONST, 0)
	}
	compilationEngine.eatSymbol(SEMICOLON)
	compilationEngine.vmWriter.WriteReturn()
}

func (compilationEngine *CompilationEngine) compileExpression() {
	compilationEngine.compileTerm()
	for compilationEngine.isSymbol(PLUS, MINUS, MULTIPLY, DIVIDE, AND, OR, LESS, GREATER, EQUAL) {
		symbol := compilationEngine.eatSymbol(PLUS, MINUS, MULTIPLY, DIVIDE, AND, OR, LESS, GREATER, EQUAL)
		compilationEngine.compileTerm()
		compilationEngine.vmWriter.WriteArithmetic(symbol)
	}
}

func (compilationEngine *CompilationEngine) compileTerm() {
	if compilationEngine.isKeyword() {
		keyword := compilationEngine.eatKeyword()
		if keyword == TRUE || keyword == FALSE || keyword == NULL {
			compilationEngine.vmWriter.WritePush(CONST, 0)
			if keyword == TRUE {
				compilationEngine.vmWriter.WriteArithmetic(NOT)
			}
		} else if keyword == THIS {
			compilationEngine.vmWriter.WritePush(POINTER, 0)
		}
	} else if compilationEngine.isIdentifier() {
		name := compilationEngine.eatIdentifier()
		if compilationEngine.isSymbol(LEFT_BRACKET) {
			compilationEngine.eatSymbol(LEFT_BRACKET)
			compilationEngine.compileExpression()
			compilationEngine.eatSymbol(RIGHT_BRACKET)
			compilationEngine.vmWriter.WritePush(compilationEngine.symbolTable.GetVariableInfo(name))
			compilationEngine.vmWriter.WriteArithmetic(PLUS)
			compilationEngine.vmWriter.WritePop(POINTER, 1)
			compilationEngine.vmWriter.WritePush(THAT, 0)
		} else if compilationEngine.isSymbol(LEFT_PARANTHESIS) {
			compilationEngine.eatSymbol(LEFT_PARANTHESIS)
			count := compilationEngine.compileExpressionList()
			compilationEngine.eatSymbol(RIGHT_PARANTHESIS)
			compilationEngine.vmWriter.WriteCall(name, count)
		} else if compilationEngine.isSymbol(DOT) {
			compilationEngine.eatSymbol(DOT)
			count := 0
			if compilationEngine.symbolTable.HasVariable(name) {
				compilationEngine.vmWriter.WritePush(compilationEngine.symbolTable.GetVariableInfo(name))
				name = compilationEngine.symbolTable.GetVariableType(name)
				count = 1
			}
			name += "." + compilationEngine.eatIdentifier()
			compilationEngine.eatSymbol(LEFT_PARANTHESIS)
			count += compilationEngine.compileExpressionList()
			compilationEngine.eatSymbol(RIGHT_PARANTHESIS)
			compilationEngine.vmWriter.WriteCall(name, count)
		} else {
			compilationEngine.vmWriter.WritePush(compilationEngine.symbolTable.GetVariableInfo(name))
		}
	} else if compilationEngine.tokenizer.GetTokenType() == STRING_CONST {
		stringConst := compilationEngine.eatString()
		compilationEngine.vmWriter.WritePush(CONST, len(stringConst))
		compilationEngine.vmWriter.WriteCall("String.new", 1)
		for _, c := range stringConst {
			compilationEngine.vmWriter.WritePush(CONST, int(c))
			compilationEngine.vmWriter.WriteCall("String.appendChar", 2)
		}
	} else if compilationEngine.tokenizer.GetTokenType() == INT_CONST {
		value := compilationEngine.eatInteger()
		compilationEngine.vmWriter.WritePush(CONST, value)
	} else if compilationEngine.isSymbol(LEFT_PARANTHESIS) {
		compilationEngine.eatSymbol(LEFT_PARANTHESIS)
		compilationEngine.compileExpression()
		compilationEngine.eatSymbol(RIGHT_PARANTHESIS)
	} else if compilationEngine.isSymbol(MINUS, NOT) {
		symbol := compilationEngine.eatSymbol(MINUS, NOT)
		if symbol == MINUS {
			symbol = NEG
		}
		compilationEngine.compileTerm()
		compilationEngine.vmWriter.WriteArithmetic(symbol)
	}
}

func (compilationEngine *CompilationEngine) compileExpressionList() int {
	count := 0
	if !compilationEngine.isSymbol() || compilationEngine.isSymbol(MINUS, NOT, LEFT_PARANTHESIS) {
		count++
		compilationEngine.compileExpression()
		for compilationEngine.isSymbol(COMMA) {
			count++
			compilationEngine.eatSymbol(COMMA)
			compilationEngine.compileExpression()
		}
	}
	return count
}

func (compilationEngine *CompilationEngine) compileType() (bool, string) {
	var typeVar string
	if compilationEngine.isKeyword(INT, CHAR, BOOLEAN) {
		typeVar = compilationEngine.tokenizer.GetStringValue()
		compilationEngine.eatKeyword(INT, CHAR, BOOLEAN)
	} else if compilationEngine.isIdentifier() {
		typeVar = compilationEngine.tokenizer.GetIdentifier()
		compilationEngine.eatIdentifier()
	} else {
		return false, ""
	}
	return true, typeVar
}

func (compilationEngine *CompilationEngine) eatKeyword(keywords ...Keyword) Keyword {
	if !compilationEngine.isKeyword(keywords...) {
		compilationEngine.writeError("Expected keyword")
	}
	keyword := compilationEngine.tokenizer.GetKeyword()
	compilationEngine.tokenizer.Advance()
	return keyword
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

func (compilationEngine *CompilationEngine) eatSymbol(symbols ...Symbol) Symbol {
	if !compilationEngine.isSymbol(symbols...) {
		panic("")
	}
	symbol := compilationEngine.tokenizer.GetSymbol()
	compilationEngine.tokenizer.Advance()
	return symbol
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

func (compilationEngine *CompilationEngine) eatString() string {
	if compilationEngine.tokenizer.GetTokenType() != STRING_CONST {
		compilationEngine.writeError("Expected string constant")
	}
	stringConst := compilationEngine.tokenizer.GetStringValue()
	compilationEngine.tokenizer.Advance()
	return stringConst
}

func (compilationEngine *CompilationEngine) eatInteger() int {
	if compilationEngine.tokenizer.GetTokenType() != INT_CONST {
		compilationEngine.writeError("Expected integer constant")
	}
	value := compilationEngine.tokenizer.GetIntegerValue()
	compilationEngine.tokenizer.Advance()
	return value
}

func (compilationEngine *CompilationEngine) eatIdentifier() string {
	if !compilationEngine.isIdentifier() {
		compilationEngine.writeError("Expected identifier")
	}
	identifier := compilationEngine.tokenizer.GetIdentifier()
	compilationEngine.tokenizer.Advance()
	return identifier
}

func (compilationEngine *CompilationEngine) eatIdentifierDefinition(variableType string, kind Keyword) {
	if !compilationEngine.isIdentifier() {
		compilationEngine.writeError("Expected identifier")
	}
	identifier := compilationEngine.tokenizer.GetIdentifier()
	compilationEngine.tokenizer.Advance()
	compilationEngine.symbolTable.Define(identifier, variableType, kind)
}

func (compilationEngine *CompilationEngine) isIdentifier() bool {
	return compilationEngine.tokenizer.GetTokenType() == IDENTIFIER
}

func (compilationEngine *CompilationEngine) getLabelIfTrue() string {
	return compilationEngine.getLabelIf("TRUE")
}

func (compilationEngine *CompilationEngine) getLabelIfFalse() string {
	return compilationEngine.getLabelIf("FALSE")
}

func (compilationEngine *CompilationEngine) getLabelIfEnd() string {
	return compilationEngine.getLabelIf("END")
}

func (compilationEngine *CompilationEngine) getLabelIf(suffix string) string {
	return "IF_" + suffix + strconv.Itoa(compilationEngine.counterIf-1)
}

func (compilationEngine *CompilationEngine) getLabelWhile() string {
	return "WHILE_EXP" + strconv.Itoa(compilationEngine.counterWhile-1)
}

func (compilationEngine *CompilationEngine) getLabelWhileEnd() string {
	return "WHILE_END" + strconv.Itoa(compilationEngine.counterWhile-1)
}

func (compilationEngine *CompilationEngine) writeError(format string) {
	fmt.Fprint(os.Stderr, format)
	os.Exit(1)
}
