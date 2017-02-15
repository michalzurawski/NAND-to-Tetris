package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Encapsulates access to the input code. Reads an Virtual Machine command,
// parses it, and provides convenient access to the command’s components
// (fields and symbols). In addition, removes all white space and comments.
type CodeWriter struct {
	file               *os.File
	labelPrefix        string
	labelCount         int
	commandTranslation map[ArithmeticCommand]string
	segmentTranslation map[Segment]string
}

// Opens the input file
func NewCodeWriter(fileName string) *CodeWriter {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Could not save file", fileName)
		os.Exit(1)
	}

	commandTranslation := make(map[ArithmeticCommand]string)
	commandTranslation[NEG] = "M=-M"
	commandTranslation[NOT] = "M=!M"
	commandTranslation[ADD] = "M=D+M"
	commandTranslation[SUB] = "M=M-D"
	commandTranslation[AND] = "M=D&M"
	commandTranslation[OR] = "M=D|M"
	commandTranslation[EQ] = "D;JEQ"
	commandTranslation[GT] = "D;JGT"
	commandTranslation[LT] = "D;JLT"

	segmentTranslation := make(map[Segment]string)
	segmentTranslation[ARGUMENT] = "@ARG"
	segmentTranslation[LOCAL] = "@LCL"
	segmentTranslation[THIS] = "@THIS"
	segmentTranslation[THAT] = "@THAT"
	segmentTranslation[POINTER] = "@THIS"
	segmentTranslation[TEMP] = "@R5"

	return &CodeWriter{commandTranslation: commandTranslation, segmentTranslation: segmentTranslation, file: file}
}

// Closes the file
func (codeWriter *CodeWriter) Close() error {
	return codeWriter.file.Close()
}

// Sets the file name of current parsing file
func (codeWriter *CodeWriter) SetFileName(fileName string) {
	labelPrefix := strings.Replace(fileName, "/", "_", -1)
	codeWriter.labelPrefix = labelPrefix
}

// Writes bootsrap
func (codeWriter *CodeWriter) WriteInit() {
	codeWriter.WriteComment("SP = 256")
	codeWriter.write("@256")
	codeWriter.write("D=A")
	codeWriter.write("@SP")
	codeWriter.write("M=D")
	codeWriter.WriteComment("call Sys.init")
	codeWriter.WriteCall("Sys.init", 0)
}

// Writes the assembly code that is the translation of the given ARITHMETIC command
func (codeWriter *CodeWriter) WriteArithmetic(command ArithmeticCommand) {
	// a = pop()
	codeWriter.write("@SP")
	if command == NEG || command == NOT {
		// push(command(a))
		codeWriter.write("A=M-1")
		codeWriter.writeCommandTranslation(command)
	} else {
		// b = pop()
		codeWriter.write("AM=M-1")
		codeWriter.write("D=M")
		codeWriter.write("A=A-1")
		if command == ADD || command == SUB || command == AND || command == OR {
			// push(command(a,b))
			codeWriter.writeCommandTranslation(command)
		} else {
			firstLabel := codeWriter.nextLabel()
			secondLabeL := codeWriter.nextLabel()
			// if (command(b, a)) a = -1 else a = 0
			codeWriter.write("D=M-D")
			codeWriter.write("@" + firstLabel)
			codeWriter.writeCommandTranslation(command)
			codeWriter.write("D=0")
			codeWriter.write("@" + secondLabeL)
			codeWriter.write("0;JMP")
			codeWriter.write("(" + firstLabel + ")")
			codeWriter.write("D=-1")
			codeWriter.write("(" + secondLabeL + ")")
			// push(a)
			codeWriter.write("@SP")
			codeWriter.write("A=M-1")
			codeWriter.write("M=D")
		}
	}
}

// Writes the assembly code that is the translation of the given POP command
func (codeWriter *CodeWriter) WritePop(segment Segment, index string) {
	if segment == STATIC {
		// *(static+index) = pop()
		codeWriter.write("@SP")
		codeWriter.write("AM=M-1")
		codeWriter.write("D=M")
		codeWriter.write("@" + codeWriter.labelPrefix + "." + index)
		codeWriter.write("M=D")
	} else {
		codeWriter.write("@" + index)
		codeWriter.write("D=A")
		codeWriter.writeSegmentTranslation(segment)
		if segment == POINTER || segment == TEMP {
			// segment + i = pop()
			codeWriter.write("D=A+D")
		} else {
			// *(segment + i) = pop()
			codeWriter.write("D=M+D")
		}
		codeWriter.write("@R13")
		codeWriter.write("M=D")
		codeWriter.write("@SP")
		codeWriter.write("AM=M-1")
		codeWriter.write("D=M")
		codeWriter.write("@R13")
		codeWriter.write("A=M")
		codeWriter.write("M=D")
	}
}

// Writes the assembly code that is the translation of the given PUSH
func (codeWriter *CodeWriter) WritePush(segment Segment, index string) {
	if segment == STATIC {
		// push(*(static+index))
		codeWriter.write("@" + codeWriter.labelPrefix + "." + index)
		codeWriter.write("D=M")
	} else {
		codeWriter.write("@" + index)
		codeWriter.write("D=A")
		if segment != CONSTANT {
			codeWriter.writeSegmentTranslation(segment)
			if segment == POINTER || segment == TEMP {
				// push(segment+index)
				codeWriter.write("A=A+D")
			} else {
				// push(*(semgent+index))
				codeWriter.write("A=M+D")
			}
			codeWriter.write("D=M")
		}
	}
	codeWriter.write("@SP")
	codeWriter.write("M=M+1")
	codeWriter.write("A=M-1")
	codeWriter.write("M=D")
}

// Writes the assembly code that is translation of call command
func (codeWriter *CodeWriter) WriteCall(functionName string, argNumber int) {
	label := codeWriter.nextLabel()
	codeWriter.WritePush(CONSTANT, label)
	// codeWriter.WritePush(CONSTANT, "LCL")
	codeWriter.write("@LCL")
	codeWriter.write("D=M")
	codeWriter.write("@SP")
	codeWriter.write("M=M+1")
	codeWriter.write("A=M-1")
	codeWriter.write("M=D")
	// codeWriter.WritePush(CONSTANT, "ARG")
	codeWriter.write("@ARG")
	codeWriter.write("D=M")
	codeWriter.write("@SP")
	codeWriter.write("M=M+1")
	codeWriter.write("A=M-1")
	codeWriter.write("M=D")
	// codeWriter.WritePush(CONSTANT, "THIS")
	codeWriter.write("@THIS")
	codeWriter.write("D=M")
	codeWriter.write("@SP")
	codeWriter.write("M=M+1")
	codeWriter.write("A=M-1")
	codeWriter.write("M=D")
	// codeWriter.WritePush(CONSTANT, "THAT")
	codeWriter.write("@THAT")
	codeWriter.write("D=M")
	codeWriter.write("@SP")
	codeWriter.write("M=M+1")
	codeWriter.write("A=M-1")
	codeWriter.write("M=D")
	// ARG = SP - 5 - argNumber
	codeWriter.write("@SP")
	codeWriter.write("D=M")
	codeWriter.write("@" + strconv.Itoa(5+argNumber))
	codeWriter.write("D=D-A")
	codeWriter.write("@ARG")
	codeWriter.write("M=D")
	// LCL = SP
	codeWriter.write("@SP")
	codeWriter.write("D=M")
	codeWriter.write("@LCL")
	codeWriter.write("M=D")
	// goto functionName
	codeWriter.WriteGoto(functionName)
	// (Return_address)
	codeWriter.WriteLabel(label)
}

// Writes the assembly code that is translation of function command
func (codeWriter *CodeWriter) WriteFunction(functionName string, localNumber int) {
	codeWriter.WriteLabel(functionName)
	for i := 0; i < localNumber; i++ {
		codeWriter.WritePush(CONSTANT, "0")
	}
}

// Writes the assembly code that is translation of return command
func (codeWriter *CodeWriter) WriteReturn() {
	// frame = LCL
	codeWriter.write("@LCL")
	codeWriter.write("D=M")
	codeWriter.write("@R13")
	codeWriter.write("M=D")
	// returnAddress = *(frame - 5)
	codeWriter.write("@5")
	codeWriter.write("A=D-A")
	codeWriter.write("D=M")
	codeWriter.write("@R14")
	codeWriter.write("M=D")
	// *ARG = pop
	codeWriter.write("@SP")
	codeWriter.write("AM=M-1")
	codeWriter.write("D=M")
	codeWriter.write("@ARG")
	codeWriter.write("A=M")
	codeWriter.write("M=D")
	// SP = ARG+1
	codeWriter.write("D=A+1")
	codeWriter.write("@SP")
	codeWriter.write("M=D")
	// THAT = *(frame - 1)
	codeWriter.write("@R13")
	codeWriter.write("D=M")
	codeWriter.write("@1")
	codeWriter.write("A=D-A")
	codeWriter.write("D=M")
	codeWriter.write("@THAT")
	codeWriter.write("M=D")
	// THIS = *(frame - 2)
	codeWriter.write("@R13")
	codeWriter.write("D=M")
	codeWriter.write("@2")
	codeWriter.write("A=D-A")
	codeWriter.write("D=M")
	codeWriter.write("@THIS")
	codeWriter.write("M=D")
	// ARG = *(frame - 3)
	codeWriter.write("@R13")
	codeWriter.write("D=M")
	codeWriter.write("@3")
	codeWriter.write("A=D-A")
	codeWriter.write("D=M")
	codeWriter.write("@ARG")
	codeWriter.write("M=D")
	// LCL = *(frame - 4)
	codeWriter.write("@R13")
	codeWriter.write("D=M")
	codeWriter.write("@4")
	codeWriter.write("A=D-A")
	codeWriter.write("D=M")
	codeWriter.write("@LCL")
	codeWriter.write("M=D")
	// goto returnAddress
	codeWriter.write("@R14")
	codeWriter.write("A=M")
	codeWriter.write("0;JMP")
}

// Writes the assembly code that is translation of if-goto command
func (codeWriter *CodeWriter) WriteIf(label string) {
	// if pop() != 0 goto label
	codeWriter.write("@SP")
	codeWriter.write("AM=M-1")
	codeWriter.write("D=M")
	codeWriter.write("@" + label)
	codeWriter.write("D;JNE")
}

// Writes the assembly code that is translation of goto command
func (codeWriter *CodeWriter) WriteGoto(label string) {
	// goto label
	codeWriter.write("@" + label)
	codeWriter.write("0;JMP")
}

// Writes the assembly code that is translation of label command
func (codeWriter *CodeWriter) WriteLabel(label string) {
	codeWriter.write("(" + label + ")")
}

// Writes the comment
func (codeWriter *CodeWriter) WriteComment(comment string) {
	codeWriter.write("// " + comment)
}

func (codeWriter *CodeWriter) nextLabel() string {
	codeWriter.labelCount++
	return "__internal__" + codeWriter.labelPrefix + strconv.Itoa(codeWriter.labelCount)
}

func (codeWriter *CodeWriter) writeCommandTranslation(command ArithmeticCommand) {
	codeWriter.write(codeWriter.commandTranslation[command])
}

func (codeWriter *CodeWriter) writeSegmentTranslation(segment Segment) {
	codeWriter.write(codeWriter.segmentTranslation[segment])
}

func (codeWriter *CodeWriter) write(command string) {
	codeWriter.file.WriteString(command + "\n")
}
