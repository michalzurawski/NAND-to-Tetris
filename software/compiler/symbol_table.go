package main

type VariableInformation struct {
	kind         Keyword
	variableType string
	index        int
}

type SymbolTable struct {
	class           map[string]VariableInformation
	subroutine      map[string]VariableInformation
	staticCounter   int
	fieldCounter    int
	argumentCounter int
	variableCounter int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{class: make(map[string]VariableInformation)}
}

func (symbolTable *SymbolTable) StartSubroutine() {
	symbolTable.subroutine = make(map[string]VariableInformation)
	symbolTable.argumentCounter = 0
	symbolTable.variableCounter = 0
}

func (symbolTable *SymbolTable) Define(name, variableType string, kind Keyword) {
	switch kind {
	case STATIC:
		symbolTable.class[name] = VariableInformation{kind: kind, index: symbolTable.staticCounter, variableType: variableType}
		symbolTable.staticCounter++
	case FIELD:
		symbolTable.class[name] = VariableInformation{kind: kind, index: symbolTable.fieldCounter, variableType: variableType}
		symbolTable.fieldCounter++
	case ARG:
		symbolTable.subroutine[name] = VariableInformation{kind: kind, index: symbolTable.argumentCounter, variableType: variableType}
		symbolTable.argumentCounter++
	case VAR:
		symbolTable.subroutine[name] = VariableInformation{kind: kind, index: symbolTable.variableCounter, variableType: variableType}
		symbolTable.variableCounter++
	}
}

func (symbolTable *SymbolTable) GetVariableInfo(name string) (Keyword, int) {
	info, has := symbolTable.subroutine[name]
	if has {
		return info.kind, info.index
	}
	return symbolTable.class[name].kind, symbolTable.class[name].index
}

func (symbolTable *SymbolTable) GetVariableType(name string) string {
	info, has := symbolTable.subroutine[name]
	if has {
		return info.variableType
	}
	return symbolTable.class[name].variableType
}

func (symbolTable *SymbolTable) HasVariable(name string) bool {
	_, has := symbolTable.subroutine[name]
	if !has {
		_, has = symbolTable.class[name]
	}
	return has
}
