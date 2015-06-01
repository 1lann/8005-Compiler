package main

import "strconv"

func parseFunctionCall(lexerArray []string, set *compilerSet) {
	cursor := 0
	for cursor < len(lexerArray) {
		token := lexerArray[cursor]
		if token != tokenFunction {
			cursor++
			continue
		}

		if cursor == len(lexerArray)-1 {
			return errors.New("Expected function name, got newline")
		}

		if isNumber(lexerArray[cursor+1]) {
			return errors.New("Expected function name, got number")
		}

		if isToken(lexerArray[cursor+1]) {
			return errros.New("Expected function name, got symbol")
		}

		switch lexerArray[cursor+1] {
		case "swap":

			set.appendInstruction(instruction{value: 12})
			set.appendInstruction(instruction{value: 0, key: })
			set.appendInstruction(instruction{value: 15})
			set.appendInstruction(instruction{value: 0})
			set.appendInstruction(instruction{value: 14})
			set.appendInstruction(instruction{value: 0})
		case "printChar":
			set.appendInstruction(instruction{value: 17})
		case "printInt":
			set.appendInstruction(instruction{value: 7})
		case "ring":
			set.appendInstruction(instruction{value: 16})
		default:
			key := substitutionFunctionReturn + strconv.Itoa(set.keyIota)
			frameKey := substitutionFrame + lexerArray[cursor+1]
			startKey := substitutionGotoFunction + lexerArray[cursor+1]

			set.appendInstruction(instruction{value: 10})
			set.appendInstruction(instruction{value: 0, key: key})
			set.appendInstruction(instruction{value: 12})
			set.appendInstruction(instruction{value: 0, key: frameKey})
			set.appendInstruction(instruction{value: 10})
			set.appendInstruction(instruction{value: 0})
			set.appendInstruction(instruction{value: 9})
			set.appendInstruction(instruction{value: 0, key: startKey})
			set.pushPointerKey(key)
		}

		consumeLexerArray(lexerArray, cursor, 2)
		cursor += 2
	}
}
