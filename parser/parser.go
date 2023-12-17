package parser

import (
	"fmt"
	"os"
	"strings"
	support "../support"
)


var (
	ProgramLabels map[string]int
	Accumulator int
	InstructionPointer int
	Variables map[string]int
	ProgramData []int

	ProgramDataStart, ProgramDataEnd int
	DataItemPointer int
)


func doMarkStartOfData(lineNumber int) {
	ProgramDataStart = lineNumber+1
	DataItemPointer = 0
}

func doMarkEndOfData(lineNumber int) {
	ProgramDataEnd = lineNumber - 1
}

func doAddDataItem(dataLine string, lineNumber int) {

	dataValue, err := support.StringToInteger(dataLine)
	if err == nil {
		ProgramData = append(ProgramData, dataValue)
	} else {
		str := fmt.Sprintf("Invalid data item '%s' at line %d.\n", dataLine, lineNumber)
		support.Message(str)
	}
}

func doAdd(argument string) error {
	var err error
	
	value, ok := Variables[argument]
	if ok {
		Accumulator += value
	} else {
		value, err = support.StringToInteger(argument)
		if err == nil {
			Accumulator += value
		}
	}

	return err
}

func doSubtract(argument string) error {
	var err error
	value, ok := Variables[argument]
	if ok {
		Accumulator -= value
	} else {
		value, err = support.StringToInteger(argument)
		if err == nil {
			Accumulator -= value
		}
	}

	return err
}

func doMultiply(argument string) error {
	var err error
	value, ok := Variables[argument]
	if ok {
		Accumulator *= value
	} else {	
		value, err = support.StringToInteger(argument)
		if err == nil {
			Accumulator *= value
		}
	}

	return err
}

func doDivide(argument string) error {
	var err error
	
	value, ok := Variables[argument]
	if ok {
		if value == 0 {
			support.Message("Divide by zero")
			os.Exit(-1)
		} else {
			Accumulator /= value
		}
	} else {
		value, err = support.StringToInteger(argument)
		if err == nil {
			if value == 0 {
				support.Message("Divide by zero")
				os.Exit(-1)
			}
			Accumulator /= value
		}
	}

	return err
}

func doHalt() {
	os.Exit(0)
}

func doLine() {
	fmt.Println()
}

func doIn() {
	Accumulator = ProgramData[DataItemPointer]
	DataItemPointer++
}

func doOut() {
	fmt.Printf("%d", Accumulator)
}

func doPrint(str string) {
	fmt.Printf(str)
}


func doLoad(argument string) error {
	value, ok := Variables[argument]
	if ok {
		Accumulator = value
	} else {
		Accumulator,_ = support.StringToInteger(argument)
	}
	return nil
}

func doStore(argument string) error {
	Variables[argument] = Accumulator
	return nil
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

func buildString(parts []string, offset int) string {
	var str string

	for partId := offset; partId < len(parts); partId++ {
		if parts[partId][0] == '"' {
			str += parts[partId][1:]
		} else
		if parts[partId][len(parts[partId])-1] == '"' {
			str += parts[partId][:len(parts[partId])-1]
		} else {
			str += parts[partId]
		}

		if partId < len(parts)-1 {
			str += " "
		}
	}

	return str
}

func Parse(program []string) bool {

	for lineNumber := 0; lineNumber < len(program); lineNumber++ {

		fields := strings.Fields (program[lineNumber])
		
		var label, instruction, argument string

		if len(fields) >= 3 {
			// LABEL COMMAND ARGUMENT
			label  = strings.ToUpper(fields[0])
			_, ok := ProgramLabels[label]
			
			if !ok {
				ProgramLabels[label] = lineNumber
			}

			instruction = fields[1]

			argument = buildString(fields, 2) // build parts into a printable string
						
		} else if len(fields) == 2 {
			// COMMAND ARGUMENT
			instruction = fields[0]

			if instruction == "PRINT" {
				argument = buildString(fields, 1)
			} else {
				argument = fields[1]
			}
		} else {
			instruction = program[lineNumber]
		}
		
		if len(instruction) > 0 {
			switch (instruction) {
				case "HALT": 		doHalt()
				case "LINE": 		doLine()
				case "LOAD": 		doLoad(argument)
				case "STORE":		doStore(argument)
				case "IN": 			doIn()
				case "OUT": 		doOut()
				case "PRINT":		doPrint(argument)
				case "ADD": 		doAdd(argument)
				case "SUBTRACT": 	doSubtract(argument)
				case "MULTIPLY": 	doMultiply(argument)
				case "DIVIDE": 		doDivide(argument)
				case "JUMP": 		doJump(argument)
				case "JINEG": 		doJumpIfNegative(argument)
				case "JIZERO": 		doJumpIfZero(argument)
				case "%":			doMarkStartOfData(lineNumber)
				case "*":			doMarkEndOfData(lineNumber)
				default:
					// Comment
					if program[lineNumber][0] == '(' {
						break
					}

					if lineNumber >= ProgramDataStart && lineNumber <= ProgramDataEnd {
						doAddDataItem(program[lineNumber], lineNumber)
					} else {
						fmt.Printf ("Unknown instruction in line %d : '%s'\n", lineNumber+1, instruction)
					}
			}
		}	
	}

	return true
}
