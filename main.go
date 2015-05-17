package main

import (
	"bufio"
	// "errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func cliHandle() *bufio.Reader {
	if len(os.Args) < 2 {
		fmt.Println("No file specified! Use -h for more information")
		return nil
	}

	fileName := os.Args[1]

	if fileName == "-h" || fileName == "--help" {
		fmt.Println("8005 Compiler by Jason Chu")
		fmt.Println("Usage: 8005c [file name]")
		return nil
	}

	file, fileErr := os.Open(fileName)

	if fileErr != nil {
		if os.IsNotExist(fileErr) {
			fmt.Println(fileName + ": file does not exist!")
			return nil
		} else if os.IsPermission(fileErr) {
			fmt.Println(fileName + ": permission denied")
		} else if os.IsExist(fileErr) {
			fmt.Println(fileName + ": cannot compile folder")
		} else {
			fmt.Println(fileName + ": error while reading file")
		}
	}

	reader := bufio.NewReader(file)

	return reader
}

func main() {
	reader := cliHandle()
	if reader == nil {
		return
	}

	resp, err := compile(reader)
	if err != nil {
		return
	}
}
