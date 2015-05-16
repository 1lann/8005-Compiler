package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fileName := os.Args[1]

	if fileName == "-h" || fileName == "--help" {
		fmt.Println("8005 Compiler by Jason Chu")
		fmt.Println("Usage: 8005c [file name]")
		return
	}

	if len(fileName) == 0 {
		fmt.Println("No file specified!")
		return
	}

	text, fileErr := ioutil.ReadFile("/tmp/dat")

	if fileErr != nil {
		if os.IsNotExist(fileErr) {
			fmt.Println(fileName + ": file does not exist!")
			return
		} else if os.IsPermission(fileErr) {
			fmt.Println(fileName + ": permission denied")
		} else if os.IsExist(fileErr) {
			fmt.Println(fileName + ": cannot compile folder")
		}
	}

}
