package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Encapsulates access to the input code. Reads an Virtual Machine command,
// parses it, and provides convenient access to the command’s components
// (fields and symbols). In addition, removes all white space and comments.
type Parser struct {
	scanner *bufio.Scanner
	file    *os.File
}

type CommandType int

const (
	ARITHMETIC CommandType = iota
	PUSH       CommandType = iota
	POP        CommandType = iota
	LABEL      CommandType = iota
	GOTO       CommandType = iota
	IF         CommandType = iota
	FUNCTION   CommandType = iota
	RETURN     CommandType = iota
	CALL       CommandType = iota
)

type ArithmeticCommand int

const (
	ADD ArithmeticCommand = iota
	SUB ArithmeticCommand = iota
	NEG ArithmeticCommand = iota
	NOT ArithmeticCommand = iota
	AND ArithmeticCommand = iota
	OR  ArithmeticCommand = iota
	EQ  ArithmeticCommand = iota
	LT  ArithmeticCommand = iota
	GT  ArithmeticCommand = iota
)

type Segment int

const (
	ARGUMENT Segment = iota
	LOCAL    Segment = iota
	THIS     Segment = iota
	THAT     Segment = iota
	POINTER  Segment = iota
	TEMP     Segment = iota
	STATIC   Segment = iota
	CONSTANT Segment = iota
)

// Opens the input file
func NewParser(fileName string) *Parser {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open file %s", fileName)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(file)
	return &Parser{scanner: scanner, file: file}
}

// Closes the file
func (parser *Parser) Close() error {
	return parser.file.Close()
}

// Reads the next command from the input and makes it current
// Returns true if there are more commands in the input
func (parser *Parser) Advance() bool {
	for parser.scanner.Scan() {
		text := parser.GetVMCommand()
		if len(text) > 0 {
			return true
		}
	}
	return false
}

// Returns the type of current command
func (parser *Parser) GetCommandType() CommandType {
	text := parser.GetVMCommand()
	switch firstWord := strings.Split(text, " ")[0]; firstWord {
	case "label":
		return LABEL
	case "goto":
		return GOTO
	case "if-goto":
		return IF
	case "function":
		return FUNCTION
	case "call":
		return CALL
	case "return":
		return RETURN
	case "pop":
		return POP
	case "push":
		return PUSH
	}
	return ARITHMETIC
}

// Returns arithmetic command
func (parser *Parser) GetArithmeticCommand() ArithmeticCommand {
	text := parser.GetVMCommand()
	switch text {
	case "add":
		return ADD
	case "sub":
		return SUB
	case "and":
		return AND
	case "or":
		return OR
	case "not":
		return NOT
	case "neg":
		return NEG
	case "eq":
		return EQ
	case "gt":
		return GT
	case "lt":
		return LT
	}
	panic("Unkown arithmetic command: " + text)
}

func (parser *Parser) GetSegment() Segment {
	text := parser.GetVMCommand()
	segment := strings.Split(text, " ")[1]
	switch segment {
	case "argument":
		return ARGUMENT
	case "local":
		return LOCAL
	case "this":
		return THIS
	case "that":
		return THAT
	case "pointer":
		return POINTER
	case "temp":
		return TEMP
	case "static":
		return STATIC
	case "constant":
		return CONSTANT
	}
	panic("Unkown semgent type " + segment)
}

// Returns the first argument of current command.
// Should not be called if the current command is RETURN
func (parser *Parser) GetFirstArgument() string {
	text := parser.GetVMCommand()
	return strings.Split(text, " ")[1]
}

// Returns the second argument of current command.
func (parser *Parser) GetSecondArgumentAsInt() int {
	text := parser.GetSecondArgument()
	number, _ := strconv.Atoi(text)
	return number
}

// Returns the second argument of current command.
// Should be called only if the current command is PUSH, POP, FUNCTION or CALL.
func (parser *Parser) GetSecondArgument() string {
	text := parser.GetVMCommand()
	return strings.Split(text, " ")[2]
}

func (parser *Parser) GetVMCommand() string {
	text := parser.scanner.Text()
	commentIndex := strings.Index(text, "//")
	if commentIndex != -1 {
		text = text[:commentIndex]
	}
	return strings.TrimSpace(text)
}
