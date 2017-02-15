package main

import "strconv"

// Keeps a correspondence between symbolic labels and numeric addresses.
type SymbolTable struct {
	nextVariableAddress int
	table               map[string]int
}

func NewSymbolTable() *SymbolTable {
	table := make(map[string]int)
	virtualRegisters := 0x0010
	for i := 0; i < virtualRegisters; i++ {
		table["R"+strconv.Itoa(i)] = i
	}
	table["SP"] = 0x0000
	table["LCL"] = 0x0001
	table["ARG"] = 0x0002
	table["THIS"] = 0x0003
	table["THAT"] = 0x0004
	table["THAT"] = 0x0004
	table["SCREEN"] = 0x4000
	table["KBD"] = 0x6000
	return &SymbolTable{nextVariableAddress: virtualRegisters, table: table}
}

// Adds the pair (symbol, address) to the table.
func (symbolTable *SymbolTable) AddEntry(symbol string, address int) {
	symbolTable.table[symbol] = address
}

// True if contains given symbol
func (symbolTable *SymbolTable) HasSymbol(symbol string) bool {
	_, has := symbolTable.table[symbol]
	return has
}

// Returns the address associated with the symbol.
func (symbolTable *SymbolTable) GetAddress(symbol string) int {
	address, _ := symbolTable.table[symbol]
	return address
}
