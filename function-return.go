package main

import "strconv"

func parseLoopBreak(lexerArray []string, set *compilerSet) error {
	cursor := 0
	for cursor < len(lexerArray) {
		token := lexerArray[cursor]
		if token != tokenReturn {
			cursor++
			continue
		}

		key := substitutionFrameReturn + thisBlock.name
		set.appendInstruction(instruction{value: 10})
		set.appendInstruction(instruction{value: 0})
		set.appendInstruction(instruction{value: 9})
		set.appendInstruction(instruction{value: 0, key: key})

		consumeLexerArray(lexerArray, cursor, 1)
		cursor += 2
	}
}
