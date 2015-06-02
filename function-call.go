package eightc

import (
	"errors"
	"strconv"
)

func parseFunctionCall(lexerArray []string, cursor *int, set *compilerSet) error {
	token := lexerArray[*cursor]
	if token != tokenFunction {
		return nil
	}

	if *cursor == len(lexerArray)-1 {
		return errors.New("Expected function name, got newline")
	}

	if isNumber(lexerArray[*cursor+1]) {
		return errors.New("Expected function name, got number")
	}

	if isToken(lexerArray[*cursor+1]) {
		return errors.New("Expected function name, got symbol")
	}

	switch lexerArray[*cursor+1] {
	case "swap":
		key := substitutionVariable + strconv.Itoa(set.tempMemory)

		set.appendInstruction(instruction{value: 12})
		set.appendInstruction(instruction{key: key})
		set.appendInstruction(instruction{value: 15})
		set.appendInstruction(instruction{key: key})
		set.appendInstruction(instruction{value: 14})
		set.appendInstruction(instruction{key: key})
	case "printChar":
		set.appendInstruction(instruction{value: 17})
	case "printInt":
		set.appendInstruction(instruction{value: 7})
	case "ring":
		set.appendInstruction(instruction{value: 16})
	default:
		key := substitutionFunctionReturn + strconv.Itoa(set.getUniqueIota())
		frameKey := substitutionFrame + lexerArray[*cursor+1]
		startKey := substitutionGotoFunction + lexerArray[*cursor+1]

		set.appendInstruction(instruction{value: 10})
		set.appendInstruction(instruction{key: key})
		set.appendInstruction(instruction{value: 12})
		set.appendInstruction(instruction{key: frameKey})
		set.appendInstruction(instruction{value: 10})
		set.appendInstruction(instruction{value: 0})
		set.appendInstruction(instruction{value: 9})
		set.appendInstruction(instruction{key: startKey})
		set.pushPointerKey(key)
	}

	consumeLexerArray(lexerArray, cursor, 2)

	return nil
}
