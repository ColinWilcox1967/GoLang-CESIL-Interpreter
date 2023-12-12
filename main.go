package main

import (
	"fmt"
	parser "./parser"
	loader "./loader"
)

const (
	CESIL_VERSION = "1.0"
)

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

	initParser()

	lines := loader.ReadFileAsLines("cesil.txt")
}