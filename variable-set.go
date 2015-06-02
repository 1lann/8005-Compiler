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

	found := false
	for _, block := range set.blocks {
		if block.name == lexerArray[*cursor+2] {
			if len(block.instructions) != 1 ||
				len(block.instructions[0].pointerKey) != 1 ||
				block.instructions[0].pointerKey[0] != substitutionVariable+
					lexerArray[*cursor+2] {
				return errors.New("Cannot have variable and function with " +
					"the same name")
			}
			found = true
			break
		}
	}

	if !found {
		variableBlock := block{name: lexerArray[*cursor+2]}
		variableBlock.instructions = append(variableBlock.instructions,
			instruction{pointerKey: []string{substitutionVariable +
				lexerArray[*cursor+2]}, line: set.currentLine})
		set.blocks = append(set.blocks, variableBlock)
	}

	set.appendInstruction(instruction{value: 12 + instructDiff})
	set.appendInstruction(instruction{key: substitutionVariable +
		lexerArray[*cursor+2]})
	consumeLexerArray(lexerArray, cursor, 3)
	return nil
}
