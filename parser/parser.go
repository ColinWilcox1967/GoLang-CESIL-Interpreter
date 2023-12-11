package parser

import (
	"fmt"
	"strconv"
)

var (
	ProgramLabels map[string]int
	Accumulator int
	Variables map[string]int
)

func doHalt () {

}

func doLine () {
	fmt.Println()
}

func doIn () {

}

func doOut () {
	fmt.Printf("%d",Accumulator)
}

func doPrint (str string) {
	fmt.Printf(str)
}


func doLoad (argument string) error {
	value, err := strconv.Atoi(argument)
	if err == nil {
		Accumulator = value
		return nil
	}

	return err
}

func doStore (argument string) error {
	
}

func Parse (program []string) bool {

	for lineNumber := 0; lineNumber < len(program); lineNumber++ {

		label := ""
		instruction := ""
		argument := ""

		switch (instruction) {
			case "HALT": 		doHalt()
			case "LINE": 		doLine()
			case "LOAD": 		doLoad(argument)
			case "STORE":		doStore(argument)
			case "IN": 			doIn()
			case "OUT": 		doOut()
			case "PRINT": 		doPrint(argument)
			case "ADD": 		doAdd(argument)
			case "SUBTRACT": 	doSubtract(argument)
			case "MULTIPLY": 	doMultiply(argument)
			case "DIVIDE": 		doDivide(argument)
			case "JUMP": 		doJump(argument)
			case "JINEG": 		doJumpIfNegative(argument)
			case "JIZERO": 		doJumpIfZero()

			default:
				 fmt.Printf ("Unknown instruction in line %d : '%s'\n", lineNumber+1, instruction)

		}
	}
}
