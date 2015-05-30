package main

import (
	"errors"
	"strconv"
)

func mapDefines(lexerArray []string, set *compilerSet) error {
	cursor := 0
	for cursor < len(lexerArray) {
		token := lexerArray[cursor]
		if token != tokenDefine {
			cursor++
			continue
		}

		if len(lexerArray)-cursor <= 2 {
			return errors.New("Expected #define key and value, got new line")
		}

		if isNumber(lexerArray[cursor+1]) {
			return errors.New("#define key cannot be a number")
		}

		if isToken(lexerArray[cursor+1]) {
			return errors.New("#define got unexpected token: " + lexerArray[cursor+1])
		}

		if !isNumber(lexerArray[cursor+2]) {
			return errors.New("#define value must be a number")
		}

		num, _ := strconv.Atoi(lexerArray[cursor+2])

		if num > 255 {
			warn(set, "> 255 number overflow")
		}

		set.defineMap[lexerArray[cursor+1]] = num

		consumeLexerArray(lexerArray, cursor, 3)
		cursor += 3
	}

	return nil
}
