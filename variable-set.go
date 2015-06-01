package main

import (
	"errors"
	"strconv"
)

func parseSetVariable(lexerArray []string, cursor *int, set *compilerSet) error {
	token := lexerArray[*cursor]
	if token != tokenR0 && token != tokenR1 {
		return nil
	}

	if len(lexerArray)-*cursor <= 1 {
		return nil
	}

	if lexerArray[*cursor+1] != tokenStore {
		return nil
	}

	if len(lexerArray)-*cursor <= 2 {
		return errors.New("Missing variable name to store to")
	}

	if isToken(lexerArray[*cursor+2]) {
		return errors.New("Expected number or variable name, got token \"" +
			lexerArray[*cursor+2] + "\"")
	}

	instructDiff := 0

	if token == tokenR1 {
		instructDiff = 1
	}

	if isNumber(lexerArray[*cursor+2]) {
		num, _ := strconv.Atoi(lexerArray[*cursor+2])

		set.appendInstruction(instruction{value: 12 + instructDiff})
		set.appendInstruction(instruction{value: num})
		consumeLexerArray(lexerArray, cursor, 3)
		return nil
	}

	variableBlock := block{}
	variableBlock.instructions = append(variableBlock.instructions,
		instruction{pointerKey: []string{substitutionVariable +
			lexerArray[*cursor+2]}})
	set.blocks = append(set.blocks, variableBlock)

	set.appendInstruction(instruction{value: 12 + instructDiff})
	set.appendInstruction(instruction{key: substitutionVariable +
		lexerArray[*cursor+2]})
	consumeLexerArray(lexerArray, cursor, 3)
	return nil
}
