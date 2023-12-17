package loader

import (
    "bufio"
    "fmt"
    "os"
)

func ReadFileAsLines (filePath string) ([]string, bool) {
 
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

//    for i:=0; i < len(lines); i++ {
  //      fmt.Printf ("%s\n",lines[i])
    //}

    
  
    ok := err == nil
    return lines, ok
}