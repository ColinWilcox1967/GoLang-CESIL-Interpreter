package parser

import (
	"fmt"
	"strconv"
	"os"
)

var (
	ProgramLabels map[string]int
	Accumulator int
	InstructionPointer int
	Variables map[string]int
)



func doAdd(argument string) error {
	value, err := stringToInteger(argument)
	if err == nil {
		Accumulator += value
	}

	return err
}

func doSubtract(argument string) error {
	value, err := stringToInteger(argument)
	if err == nil {
		Accumulator -= value
	}

	return err
}

func doMultiply(argument string) error {
	value, err := stringToInteger(argument)
	if err == nil {
		Accumulator *= value
	}

	return err
}

func doDivide(argument string) error {
	value, err := stringToInteger(argument)
	if err == nil {
		if value == 0 {
			message("Divide by zero")
			os.Exit(-1)
		}
		Accumulator /= value
	}

	return err
}

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

func doStore(argument string) error {
	
}

func doJump(argument string) {
	pos := ProgramLabels[argument]
	if pos != -1 {
		InstructionPointer = pos
	}

}

func doJumpIfNegative(argument string) {
	pos := ProgramLabels[argument] 
	if pos != -1 {
		if Accumulator < 0 {
			InstructionPointer = pos
		}		
	}
}

func doJumpIfZero(argument string) {
	pos := ProgramLabels[argument] 
	if pos != -1 {
		if Accumulator == 0 {
			InstructionPointer = pos
		}		
	}
}

func Parse (program []string) bool {

	for lineNumber := 0; lineNumber < len(program); lineNumber++ {

	//	label := ""
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
			case "JIZERO": 		doJumpIfZero(argument)

			default:
				 fmt.Printf ("Unknown instruction in line %d : '%s'\n", lineNumber+1, instruction)

		}
	}
}
