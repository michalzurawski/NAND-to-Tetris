package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Encapsulates access to the input code. Reads an assembly language command,
// parses it, and provides convenient access to the command’s components
// (fields and symbols). In addition, removes all white space and comments.
type Parser struct {
	scanner *bufio.Scanner
	file    *os.File
}

type CommandType int

const (
	ADDRESS = iota
	COMMAND = iota
	LABEL   = iota
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
		text := parser.getAssemblyCode()
		if len(text) > 0 {
			return true
		}
	}
	return false
}

// Returns the type of current command
func (parser *Parser) GetCommandType() CommandType {
	text := parser.getAssemblyCode()
	switch first := text[0]; first {
	case '@':
		return ADDRESS
	case '(':
		return LABEL
	}
	return COMMAND
}

// Returns the symbol or decimal xxx of the current command @xxx or (xxx)
// Should be called only when CommandType is ADDRESS or LABEL
func (parser *Parser) GetSymbol() string {
	text := parser.getAssemblyCode()
	if text[len(text)-1] == ')' {
		text = text[:len(text)-1]
	}
	return text[1:]
}

// Returns the dest, comp and jump mnemonics in current COMMAND
// Should be called only when CommandType is COMMAND
func (parser *Parser) GetMnemonics() (string, string, string) {
	text := parser.getAssemblyCode()
	dest := ""
	jump := ""
	eqIndex := strings.Index(text, "=")
	semIndex := strings.Index(text, ";")
	if eqIndex != -1 {
		dest = text[:eqIndex]
	}
	if semIndex != -1 {
		jump = text[semIndex+1:]
	} else {
		semIndex = len(text)
	}
	comp := text[eqIndex+1 : semIndex]
	return dest, comp, jump
}

func (parser *Parser) getAssemblyCode() string {
	text := parser.scanner.Text()
	commentIndex := strings.Index(text, "//")
	if commentIndex != -1 {
		text = text[:commentIndex]
	}
	return strings.TrimSpace(text)
}
