package loader

import (
    "bufio"
    "fmt"
    "os"
)

func ReadFileAsLines (filePath string) []string {
 
	var lines []string

    readFile, err := os.Open(filePath)
  
    if err != nil {
        fmt.Println(err)
    }
	defer readFile.Close ()

    fileScanner := bufio.NewScanner(readFile)
 
    fileScanner.Split(bufio.ScanLines)
  
    for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
    }
  
    return lines
}