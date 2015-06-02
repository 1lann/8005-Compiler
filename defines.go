package eightc

import (
	"errors"
	"strconv"
)

func mapDefines(lexerArray []string, cursor *int, set *compilerSet) error {
	token := lexerArray[*cursor]
	if token != tokenDefine {
		return nil
	}

	if len(lexerArray)-*cursor <= 2 {
		return errors.New("Expected #define key and value, got newline")
	}

	if isNumber(lexerArray[*cursor+1]) {
		return errors.New("#define key cannot be a number")
	}

	if isToken(lexerArray[*cursor+1]) {
		return errors.New("#define got unexpected token: " + lexerArray[*cursor+1])
	}

	if !isNumber(lexerArray[*cursor+2]) {
		return errors.New("#define value must be a number")
	}

	num, _ := strconv.Atoi(lexerArray[*cursor+2])

	if num > 255 {
		set.warn("> 255 number overflow")
	}

	set.defineMap[lexerArray[*cursor+1]] = num

	consumeLexerArray(lexerArray, cursor, 3)

	return nil
}
