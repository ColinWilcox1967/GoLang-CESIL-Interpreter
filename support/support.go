package support

import (
	"fmt"
	"strconv"
)

func Message(message string) {
	str := fmt.Sprintf("*** Error : %s", message)
	fmt.Println(str)
}

func StringToInteger(str string) (int, error) {
	return strconv.Atoi(str)
}