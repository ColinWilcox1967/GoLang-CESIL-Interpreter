package main

import (
	"fmt"
	"flag"
	"strings"
	"os"

	parser "./parser"
	loader "./loader"
	support "./support"
)

const (
	CESIL_VERSION = "1.0"
	DEFAULT_CESIL_FILE = "default.ces"
)

var (
	sourceFile string
)

func scanArguments () {

	flag.StringVar (&sourceFile, "file", DEFAULT_CESIL_FILE, "CESIL source file.")
	flag.Parse ()

	if len(sourceFile) == 0 {
		sourceFile = DEFAULT_CESIL_FILE
	}
}

func initParser () {
	parser.ProgramLabels = make(map[string]int)
	parser.Variables = make(map[string]int)
	parser.Accumulator = 0
	parser.InstructionPointer = 0
}

func showBanner () {
	fmt.Printf ("\nCESIL Language Interpreter, Version %s\n", CESIL_VERSION)
	fmt.Println ("(c) 2023 Colin Wilcox")
}

func main () {

	showBanner()
	scanArguments()

	initParser()

	fmt.Printf("\nSource File : '%s'\n\n", strings.ToUpper(sourceFile))

	lines, ok := loader.ReadFileAsLines(sourceFile)
	if !ok {
		support.Message("Problem reading source file.\n")
		os.Exit(-2)
	}

	ok = parser.Parse(lines)
	if (!ok) {
		support.Message("Problem parsing CESIL file.\n")
		os.Exit(-1)
	}

	fmt.Printf ("\n\nProgram ran successfully.\n")
}