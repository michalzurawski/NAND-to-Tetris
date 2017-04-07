package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	usage := "Usage: " + os.Args[0] + " name of the directory containg .jack files"
	if len(os.Args) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	directoryName := os.Args[1]
	files, _ := ioutil.ReadDir(directoryName)

	for _, file := range files {
		idx := strings.LastIndex(file.Name(), ".jack")
		if idx == -1 {
			continue
		}
		compilationEngine := NewCompilationEngine(directoryName + file.Name()[0:idx])
		defer compilationEngine.Close()

		compilationEngine.CompileClass()
	}
}
