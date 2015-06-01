package main

import (
	"errors"
	"strconv"
)

func parseLoopDefinition(lexerArray []string, cursor *int, set *compilerSet) error {
	token := lexerArray[*cursor]
	if token != tokenLoop {
		return nil
	}

	if *cursor == len(lexerArray)-1 {
		return errors.New("Expected \"{\", got newline")
	}

	if lexerArray[*cursor+1] != tokenOpenCurly {
		return errors.New("Expected \"{\", got \"" +
			lexerArray[*cursor+1] + "\"")
	}

	pointerKey := substitutionLoopStart + strconv.Itoa(set.stepUpIota())
	set.parentTypes = append(set.parentTypes, blockLoop)

	set.pushPointerKey(pointerKey)

	consumeLexerArray(lexerArray, cursor, 2)

	return nil
}
