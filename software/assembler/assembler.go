package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	usage := "Usage: " + os.Args[0] + " name of the file"
	if len(os.Args) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	fileName := os.Args[1]
	parser := NewParser(fileName)
	defer parser.Close()
	symbolTable := NewSymbolTable()
	currentInstruction := 0

	for parser.Advance() {
		switch parser.GetCommandType() {
		case ADDRESS:
			currentInstruction++
		case COMMAND:
			currentInstruction++
		case LABEL:
			symbol := parser.GetSymbol()
			symbolTable.AddEntry(symbol, currentInstruction)
		}
	}

	parser = NewParser(fileName)
	defer parser.Close()
	fileSave, err := os.Create(fileName[:len(fileName)-3] + "hack")
	if err != nil {
		fmt.Println("Could not save file", fileName)
		os.Exit(1)
	}
	defer fileSave.Close()

	for parser.Advance() {
		switch parser.GetCommandType() {
		case ADDRESS:
			symbol := parser.GetSymbol()
			address := getAddress(symbol, symbolTable)
			fileSave.WriteString(GetACommand(address) + "\n")
		case COMMAND:
			dest, comp, jump := parser.GetMnemonics()
			fileSave.WriteString(GetCCommand(dest, comp, jump) + "\n")
		}
	}
}

func getAddress(symbol string, symbolTable *SymbolTable) int {
	address, err := strconv.Atoi(symbol)
	if err != nil {
		if symbolTable.HasSymbol(symbol) {
			address = symbolTable.GetAddress(symbol)
		} else {
			symbolTable.AddEntry(symbol, symbolTable.nextVariableAddress)
			address = symbolTable.nextVariableAddress
			symbolTable.nextVariableAddress++
		}
	}
	return address
}
