package parser

import (
	"fmt"
	"strconv"
)

func message(message string) {
	str := fmt.Sprintf("*** Error : %s\n", message)
	fmt.Println(str)
}

func stringToInteger(str string) (int, error) {
	return strconv.Atoi(str)
}