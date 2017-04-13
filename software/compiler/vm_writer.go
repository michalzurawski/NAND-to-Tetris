package main

import (
	"fmt"
	"os"
	"strconv"
)

type VMWriter struct {
	file *os.File
}

func NewVMWriter(fileName string) *VMWriter {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Could not save file", fileName)
		os.Exit(1)
	}
	return &VMWriter{file: file}
}

func (vmWriter *VMWriter) Close() error {
	return vmWriter.file.Close()
}

func (vmWriter *VMWriter) WritePush(keyword Keyword, index int) {
	vmWriter.writeln("push " + segment[keyword] + " " + strconv.Itoa(index))
}

func (vmWriter *VMWriter) WritePop(keyword Keyword, index int) {
	vmWriter.writeln("pop " + segment[keyword] + " " + strconv.Itoa(index))
}

func (vmWriter *VMWriter) WriteArithmetic(symbol Symbol) {
	if symbol == DIVIDE {
		vmWriter.WriteCall("Math.divide", 2)
	} else if symbol == MULTIPLY {
		vmWriter.WriteCall("Math.multiply", 2)
	} else {
		vmWriter.writeln(sym[symbol])
	}
}

func (vmWriter *VMWriter) WriteLabel(label string) {
	vmWriter.writeln("label " + label)
}

func (vmWriter *VMWriter) WriteGoto(label string) {
	vmWriter.writeln("goto " + label)
}

func (vmWriter *VMWriter) WriteIf(label string) {
	vmWriter.writeln("if-goto " + label)
}

func (vmWriter *VMWriter) WriteCall(name string, nArgs int) {
	vmWriter.writeln("call " + name + " " + strconv.Itoa(nArgs))
}

func (vmWriter *VMWriter) WriteFunction(name string, nLocal int) {
	vmWriter.writeln("function " + name + " " + strconv.Itoa(nLocal))
}

func (vmWriter *VMWriter) WriteReturn() {
	vmWriter.writeln("return")
}

func (vmWriter *VMWriter) writeln(value string) {
	vmWriter.file.WriteString(value + "\n")
}

var segment = map[Keyword]string{
	CONST:   "constant",
	ARG:     "argument",
	VAR:     "local",
	LOCAL:   "local",
	STATIC:  "static",
	FIELD:   "this",
	THIS:    "this",
	THAT:    "that",
	POINTER: "pointer",
	TEMP:    "temp",
}

var sym = map[Symbol]string{
	AND:     "and",
	PLUS:    "add",
	OR:      "or",
	MINUS:   "sub",
	NEG:     "neg",
	EQUAL:   "eq",
	GREATER: "gt",
	LESS:    "lt",
	NOT:     "not",
}
