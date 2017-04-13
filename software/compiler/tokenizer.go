package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Encapsulates access to the input code. Reads an Virtual Machine command,
// parses it, and provides convenient access to the command’s components
// (fields and symbols). In addition, removes all white space and comments.
type Tokenizer struct {
	scanner *bufio.Scanner
	file    *os.File
	text    string
}

type TokenType int

const (
	KEYWORD      TokenType = iota
	SYMBOL       TokenType = iota
	IDENTIFIER   TokenType = iota
	INT_CONST    TokenType = iota
	STRING_CONST TokenType = iota
)

type Keyword int

const (
	CLASS       Keyword = iota
	METHOD      Keyword = iota
	FUNCTION    Keyword = iota
	CONSTRUCTOR Keyword = iota
	INT         Keyword = iota
	BOOLEAN     Keyword = iota
	CHAR        Keyword = iota
	VOID        Keyword = iota
	VAR         Keyword = iota
	STATIC      Keyword = iota
	FIELD       Keyword = iota
	LET         Keyword = iota
	DO          Keyword = iota
	IF          Keyword = iota
	ELSE        Keyword = iota
	WHILE       Keyword = iota
	RETURN      Keyword = iota
	TRUE        Keyword = iota
	FALSE       Keyword = iota
	NULL        Keyword = iota
	THIS        Keyword = iota
	ARG         Keyword = iota
	CONST       Keyword = iota
	THAT        Keyword = iota
	LOCAL       Keyword = iota
	POINTER     Keyword = iota
	TEMP        Keyword = iota
)

var keywords = map[string]Keyword{
	"class":       CLASS,
	"constructor": CONSTRUCTOR,
	"function":    FUNCTION,
	"method":      METHOD,
	"field":       FIELD,
	"static":      STATIC,
	"var":         VAR,
	"int":         INT,
	"char":        CHAR,
	"boolean":     BOOLEAN,
	"void":        VOID,
	"true":        TRUE,
	"false":       FALSE,
	"null":        NULL,
	"this":        THIS,
	"let":         LET,
	"do":          DO,
	"if":          IF,
	"else":        ELSE,
	"while":       WHILE,
	"return":      RETURN,
}

type Symbol int

const (
	LEFT_PARANTHESIS  Symbol = iota
	RIGHT_PARANTHESIS Symbol = iota
	LEFT_BRACKET      Symbol = iota
	RIGHT_BRACKET     Symbol = iota
	LEFT_CURLY        Symbol = iota
	RIGHT_CURLY       Symbol = iota
	DOT               Symbol = iota
	COMMA             Symbol = iota
	SEMICOLON         Symbol = iota
	PLUS              Symbol = iota
	MINUS             Symbol = iota
	MULTIPLY          Symbol = iota
	DIVIDE            Symbol = iota
	AND               Symbol = iota
	OR                Symbol = iota
	LESS              Symbol = iota
	GREATER           Symbol = iota
	EQUAL             Symbol = iota
	NOT               Symbol = iota
	NEG               Symbol = iota
)

var symbols = map[byte]Symbol{
	'(': LEFT_PARANTHESIS,
	')': RIGHT_PARANTHESIS,
	'[': LEFT_BRACKET,
	']': RIGHT_BRACKET,
	'{': LEFT_CURLY,
	'}': RIGHT_CURLY,
	'.': DOT,
	',': COMMA,
	';': SEMICOLON,
	'+': PLUS,
	'-': MINUS,
	'*': MULTIPLY,
	'/': DIVIDE,
	'&': AND,
	'|': OR,
	'<': LESS,
	'>': GREATER,
	'=': EQUAL,
	'~': NOT,
}

// Opens the input file
func NewTokenizer(fileName string) *Tokenizer {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open file %s", fileName)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(split)
	return &Tokenizer{scanner: scanner, file: file}
}

// Closes the file
func (tokenizer *Tokenizer) Close() error {
	return tokenizer.file.Close()
}

// Reads the next command from the input and makes it current
// Returns true if there are more commands in the input
func (tokenizer *Tokenizer) Advance() bool {
	for tokenizer.scanner.Scan() {
		tokenizer.text = strings.TrimSpace(tokenizer.scanner.Text())
		if len(tokenizer.text) > 0 {
			return true
		}
	}
	return false
}

// Returns the type of current token
func (tokenizer *Tokenizer) GetTokenString() string {
	if isKeyword(tokenizer.text) {
		return "keyword"
	} else if isSymbol(tokenizer.text) {
		return "symbol"
	} else if isInteger(tokenizer.text) {
		return "integerConstant"
	} else if isString(tokenizer.text) {
		return "stringConstant"
	}
	return "identifier"
}

// Returns the type of current token
func (tokenizer *Tokenizer) GetTokenType() TokenType {
	if isKeyword(tokenizer.text) {
		return KEYWORD
	} else if isSymbol(tokenizer.text) {
		return SYMBOL
	} else if isInteger(tokenizer.text) {
		return INT_CONST
	} else if isString(tokenizer.text) {
		return STRING_CONST
	}
	return IDENTIFIER
}

func (tokenizer *Tokenizer) GetKeyword() Keyword {
	return keywords[tokenizer.text]
}

func (tokenizer *Tokenizer) GetSymbol() Symbol {
	return symbols[tokenizer.text[0]]
}

func (tokenizer *Tokenizer) GetIdentifier() string {
	return tokenizer.text
}

func (tokenizer *Tokenizer) GetIntegerValue() int {
	num, _ := strconv.Atoi(tokenizer.text)
	return num
}

func (tokenizer *Tokenizer) GetStringValue() string {
	return tokenizer.text[1 : len(tokenizer.text)-1]
}

func isKeyword(text string) bool {
	_, ok := keywords[text]
	return ok
}

func isSymbol(text string) bool {
	_, ok := symbols[text[0]]
	return ok
}

func isInteger(text string) bool {
	_, err := strconv.Atoi(text)
	return err == nil
}

func isString(text string) bool {
	return text[0] == '"' && text[len(text)-1] == '"'
}

func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, " \n"); i >= 0 {
		if idx := bytes.Index(data[0:i], []byte("//")); idx >= 0 {
			return bytes.IndexByte(data, '\n') + 1, []byte(""), nil
		}
		if idx := bytes.Index(data[0:i], []byte("/*")); idx >= 0 {
			return bytes.Index(data, []byte("*/")) + 2, []byte(""), nil
		}
		if idx := bytes.IndexAny(data[0:i], "*/[]{}()+-,.;&|~<>="); idx >= 0 {
			if idx == 0 {
				return idx + 1, dropCR(data[0 : idx+1]), nil
			}
			return idx, dropCR(data[0:idx]), nil
		}
		if idx := bytes.IndexByte(data[0:i], '"'); idx >= 0 {
			end := bytes.IndexByte(data[idx+1:], '"') + idx + 2
			return end, dropCR(data[0:end]), nil
		}

		return i + 1, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
