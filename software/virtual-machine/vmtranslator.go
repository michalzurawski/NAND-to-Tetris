package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	usage := "Usage: " + os.Args[0] + " name of the directory containg .vm files"
	if len(os.Args) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	directoryName := os.Args[1]
	outputFile := getOutputFileName(directoryName)

	codeWriter := NewCodeWriter(outputFile)
	defer codeWriter.Close()
	codeWriter.WriteInit()

	files, _ := ioutil.ReadDir(directoryName)

	for _, file := range files {
		if lenFile := len(file.Name()); lenFile < 4 || file.Name()[lenFile-3:] != ".vm" {
			continue
		}
		parser := NewParser(directoryName + file.Name())
		defer parser.Close()
		codeWriter.SetFileName(directoryName + file.Name())

		for parser.Advance() {
			codeWriter.WriteComment(parser.GetVMCommand())
			switch parser.GetCommandType() {
			case ARITHMETIC:
				codeWriter.WriteArithmetic(parser.GetArithmeticCommand())
			case POP:
				codeWriter.WritePop(parser.GetSegment(), parser.GetSecondArgument())
			case PUSH:
				codeWriter.WritePush(parser.GetSegment(), parser.GetSecondArgument())
			case IF:
				codeWriter.WriteIf(parser.GetFirstArgument())
			case GOTO:
				codeWriter.WriteGoto(parser.GetFirstArgument())
			case LABEL:
				codeWriter.WriteLabel(parser.GetFirstArgument())
			case CALL:
				codeWriter.WriteCall(parser.GetFirstArgument(), parser.GetSecondArgumentAsInt())
			case FUNCTION:
				codeWriter.WriteFunction(parser.GetFirstArgument(), parser.GetSecondArgumentAsInt())
			case RETURN:
				codeWriter.WriteReturn()
			}
		}
	}
}

func getOutputFileName(directoryName string) string {
	if directoryName[len(directoryName)-1] == '/' {
		directoryName = directoryName[:len(directoryName)-1]
	}
	outputFile := directoryName
	if idx := strings.LastIndex(outputFile, "/"); idx != -1 {
		outputFile = outputFile[idx+1:]
	}
	return directoryName + "/" + outputFile + ".asm"
}
