package main

import (
	"bufio"
	// "errors"
	"fmt"
	"github.com/1lann/eightc"
	"os"
	"strconv"
)

func cliHandle() (*bufio.Reader, string) {
	if len(os.Args) < 2 {
		fmt.Println("No file specified! Use -h for more information")
		return nil, ""
	}

	fileName := os.Args[1]

	if fileName == "-h" || fileName == "--help" {
		fmt.Println("8005 Compiler by Jason Chu")
		fmt.Println("Usage: eightc [file name]")
		return nil, ""
	}

	file, fileErr := os.Open(fileName)

	if fileErr != nil {
		if os.IsNotExist(fileErr) {
			fmt.Println(fileName + ": file does not exist!")
			return nil, ""
		} else if os.IsPermission(fileErr) {
			fmt.Println(fileName + ": permission denied")
		} else if os.IsExist(fileErr) {
			fmt.Println(fileName + ": cannot compile folder")
		} else {
			fmt.Println(fileName + ": error while reading file")
		}
	}

	reader := bufio.NewReader(file)

	return reader, fileName
}

func main() {
	reader, filename := cliHandle()
	if reader == nil {
		return
	}

	out, err := eightc.Compile(filename, reader)
	if err != nil {
		fmt.Println("Failed to compile!")
		fmt.Println(err)
		return
	}

	textOut := ""

	for _, val := range out {
		textOut += strconv.Itoa(val) + " "
	}

	fmt.Println("Program compiled")
	fmt.Println(textOut)
}
