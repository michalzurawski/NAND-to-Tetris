package main

import (
  "fmt"
  "os"
  "strconv"
)

func main(){
  usage := "Usage: " + os.Args[0] + " name of the file"
  if len(os.Args) != 2 {
    fmt.Println(usage)
    os.Exit(1)
  }

  fileName := os.Args[1]
  file, err := os.Open(fileName)
  if err != nil {
    fmt.Println("Could not open file", fileName);
    os.Exit(1)
  }
  fileSave, err := os.Create(fileName[:len(fileName)-3] + "hack")
  if err != nil {
    fmt.Println("Could not save file", fileName);
    os.Exit(1)
  }
  defer fileSave.Close()

  parser := NewParser(file)
  symbolTable := NewSymbolTable()
  currentInstruction := 0

  for parser.Advance() {
    switch parser.GetCommandType() {
    case A_COMMAND:
      currentInstruction++
    case C_COMMAND:
      currentInstruction++
    case L_COMMAND:
      symbol := parser.GetSymbol()
      symbolTable.AddEntry(symbol, currentInstruction)
    default:
      panic("Unexpceted parser error")
    }
  }

  file.Close()
  file, _ = os.Open(fileName)
  defer file.Close()
  parser = NewParser(file)

  for parser.Advance() {
    switch parser.GetCommandType() {
    case A_COMMAND:
      symbol := parser.GetSymbol()
      address := getAddress(symbol, symbolTable)
      fileSave.WriteString(GetACommand(address) + "\n")
    case C_COMMAND:
      dest, comp, jump := parser.GetMnemonics()
      fileSave.WriteString(GetCCommand(dest, comp, jump) + "\n")
    case L_COMMAND:
    default:
      panic("Unexpceted parser error")
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
