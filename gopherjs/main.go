package main

import (
	"bufio"
	"bytes"
	"github.com/1lann/eightc"
	"github.com/gopherjs/gopherjs/js"
	"strconv"
)

func main() {
	js.Global.Set("eightc", map[string]interface{}{
		"Compile": compile,
	})
}

func compile(str string) (string, bool) {
	byteReader := bytes.NewBufferString(str)
	reader := bufio.NewReader(byteReader)

	out, err := eightc.Compile("input", reader)
	if err != nil {
		return err.Error(), false
	}

	textOut := ""

	for _, val := range out {
		textOut += strconv.Itoa(val) + " "
	}

	return textOut, true
}
