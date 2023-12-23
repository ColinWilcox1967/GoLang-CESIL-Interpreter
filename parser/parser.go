package parser

import (
	"fmt"
//	"os"
	"strings"
	"strconv"
	support "../support"
)

var (
	ProgramLabels map[string]int
	Accumulator int
	InstructionPointer int
	Variables map[string]int
	ProgramData []int

	ProgramDataStart int = -1 // undefined
	ProgramDataEnd int = -1 // undefined
	DataItemPointer int
	HaltProgram bool
	DataStartSet, DataEndSet bool
)


func doMarkStartOfData(lineNumber int) {
	ProgramDataStart = lineNumber
	DataItemPointer = 0
	DataStartSet = true
}

func doMarkEndOfData(lineNumber int) {
	ProgramDataEnd = lineNumber
	DataEndSet = true
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
			HaltProgram = true
		} else {
			Accumulator /= value
		}
	} else {
		value, err = support.StringToInteger(argument)
		if err == nil {
			if value == 0 {
				support.Message("Divide by zero")
				HaltProgram = true
			}
			Accumulator /= value
		}
	}

	return err
}

func doHalt() {
	HaltProgram = true
}

func doLine() {
	fmt.Println()
}

func doIn() {
	dataBlockSize := ProgramDataEnd-ProgramDataStart-1

	if DataItemPointer >= dataBlockSize {
		support.Message ("More data needed.")
		HaltProgram = true
	} else {
		Accumulator = ProgramData[DataItemPointer]
		DataItemPointer++
	}
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
	if pos != 0 {
		InstructionPointer = pos-1 // line 1 is actually offset zero.
	}
}

func doJumpIfNegative(argument string) {
	pos := ProgramLabels[argument] 
	if pos != 0 {
		if Accumulator < 0 {
			InstructionPointer = pos-1
		}		
	}
}

func doJumpIfZero(argument string) {
	pos := ProgramLabels[argument] 
	if pos != 0 {
		if Accumulator == 0 {
			InstructionPointer = pos-1
		}		
	}
}

func labelAlreadyDefined(label string) bool {
	return ProgramLabels[label] != 0
}

func preScanDataBlock(program []string) {

	// Check we have a matching pairt of start/end data block markers
	if ProgramDataStart >= 0 && ProgramDataEnd == -1 {
		str := fmt.Sprintf ("Missing end of data block marker.\n")
		support.Message(str)
		HaltProgram = true
	} else 
	if ProgramDataStart == -1 && ProgramDataEnd >= 0 {
		str := fmt.Sprintf ("Missing start of data block marker.\n")
		support.Message(str)
		HaltProgram = true
	}

	if !HaltProgram {
		// Check markers are correct way round
		if ProgramDataStart > ProgramDataEnd {
			str := fmt.Sprintf ("Data block start marker defined AFTER end marker.\n")
			support.Message(str)
			HaltProgram = true
		} else {
			// Now iterate through the intermediate program lines and fill out data block
			for lineIndex := ProgramDataStart+1; lineIndex < ProgramDataEnd; lineIndex++ {
				value, err := strconv.Atoi(program[lineIndex])
				if err == nil {
					ProgramData = append(ProgramData, value)
				} else {
					str := fmt.Sprintf ("Problem with data value '%s' on line %d.\n", program[lineIndex], lineIndex)
					support.Message(str)
					HaltProgram = true
				}
			}
		}
	}
}

func prescanLabels(program []string) {
	for lineNumber, line := range (program) {
		fields := strings.Fields (line)
		fieldCount := len(fields)

		switch (fieldCount) {
			case 0: // no fields = no labels!!
			
			case 1: if !validCommand(fields[0]) {
						if !labelAlreadyDefined(fields[0]) {
							if !DataStartSet { //skip over data block
								ProgramLabels[fields[0]] = lineNumber
							}
						} else {
							str := fmt.Sprintf("Label '%s' already defined at line %d.\n", strings.ToUpper(fields[0]), lineNumber+1)
							support.Message(str)
						}
					} else {
						switch(fields[0]) {
							case "%": if !DataStartSet {
											doMarkStartOfData(lineNumber)
									} else {
											str := fmt.Sprintf("Start of data block already specified at line %d.\n", ProgramDataStart)
											support.Message(str)
									}
							case "*": if !DataEndSet {
											doMarkEndOfData(lineNumber)
									} else {
											str := fmt.Sprintf("End of data block already specified at line %d.\n", ProgramDataEnd)
											support.Message(str)
									}
							default: 
						}
					}		 
			case 2,3:
			default: if !validCommand(fields[0]) {
						if !labelAlreadyDefined(fields[0]) {
							if !DataStartSet {
								ProgramLabels[fields[0]] = lineNumber
							}
						} 
					}	
			}
	}
}

func buildString(parts []string, offset int) string {
	var str string

	// append remaining fragments together separated by a space
	for partId := offset; partId < len(parts); partId++ {
		str += parts[partId]
		if partId < len(parts)-1 {
			str += " "
		}
	}

	// trim off the leading ... 
	if str[0] == '"' {
		str = str[1:]
	}

	// and trailing double quotes
	if str[len(str)-1] == '"' {
		str = str[:len(str)-1]
	}

	return str
}

func validCommand(command string) bool {
	supportCommands := [...]string{"%","(","*","HALT","LINE","IN", "OUT","LOAD","STORE","PRINT","ADD","DIVIDE","MULTIPLY","SUBTRACT","JUMP","JINEG","JIZERO"}

	for i := 0; i < len(supportCommands); i++ {
		if command == supportCommands[i] {
			return true
		}
	}
	return false
}

func Parse(program []string) bool {

	prescanLabels(program)
	preScanDataBlock(program)

	InstructionPointer = 0
	
	for !HaltProgram {
		fields := strings.Fields(program[InstructionPointer])
			
		fieldCount := len(fields)
	
		// There line formatting options are:
		// (1) LABEL COMMAND
		// (2) LABEL COMMAND ARGUMENT
		// (3) COMMAND ARGUMENT
		// (4) COMMAND
		// (5) <empty line> (skip)
		// (6) Comment line (skip)
		// (7) LABEL

		// case #5
		if fieldCount == 0 {
			InstructionPointer++
		} 
		
		// Option #4 or #7
		if fieldCount == 1 {
			if validCommand(fields[0]) {
				switch (fields[0]) {
					case "IN": doIn()
					case "OUT": doOut()
					case "LINE": doLine()
					case "HALT": doHalt()
					default:
				}
				
			} else {

			}
			InstructionPointer++
		}
		
		// Options #1 or #3
		if fieldCount == 2 {
			command := fields[0]
			if !validCommand(command) {
				command = fields[1]
			}
			if validCommand(command) {
				switch (command) {
					case "IN": 			doIn()
					case "OUT": 		doOut()
					case "LINE": 		doLine()
					case "HALT": 		doHalt()
					case "LOAD": 		doLoad(fields[1])
					case "STORE":		doStore(fields[1])
					case "PRINT":		argument := buildString(fields, 1)
										doPrint(argument)
					case "ADD": 		doAdd(fields[1])
					case "SUBTRACT": 	doSubtract(fields[1])
					case "MULTIPLY": 	doMultiply(fields[1])
					case "DIVIDE": 		doDivide(fields[1])
					case "JUMP": 		doJump(fields[1])
					case "JINEG": 		doJumpIfNegative(fields[1])
					case "JIZERO": 		doJumpIfZero(fields[1])
					default:
				}
			}
			InstructionPointer++
		}

		// Option #2
		if fieldCount >= 3 {
			var label string
				if fields[0] == "PRINT" {
				argument := buildString(fields, 1)
				doPrint(argument)
			} else {
					if fields[1] == "PRINT" {
					label = fields[0]
					_, exists := ProgramLabels[label]
					if !exists {
						if !DataStartSet {
							ProgramLabels[label] = InstructionPointer+1
						}
					}
						argument := buildString(fields,2)
					doPrint(argument)
				} else {
					label := fields[0]
					command := fields[1]
					argument := fields[2]
	
					if validCommand(command) {
						_, exists := ProgramLabels[label]
						if !exists {
							if !DataStartSet {
								ProgramLabels[label] = InstructionPointer+1
							}
						} 
						switch (command) {
							case "IN": 			doIn()
							case "OUT": 		doOut()
							case "LINE": 		doLine()
							case "HALT": 		doHalt()
							case "LOAD": 		doLoad(argument)
							case "STORE":		doStore(argument)
							case "ADD": 		doAdd(argument)
							case "SUBTRACT": 	doSubtract(argument)
							case "MULTIPLY": 	doMultiply(argument)
							case "DIVIDE": 		doDivide(argument)
							case "JUMP": 		doJump(argument)
							case "JINEG": 		doJumpIfNegative(argument)
							case "JIZERO": 		doJumpIfZero(argument)
							case "%":			doMarkStartOfData(InstructionPointer)
							case "*":			doMarkEndOfData(InstructionPointer)
							default:
						}
					}
				}
			} 
			InstructionPointer++
		}

		if InstructionPointer >= len(program){
			HaltProgram = true
		}
	}

	return true
}
