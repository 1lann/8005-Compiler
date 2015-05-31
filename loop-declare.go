package main

import "strconv"

func parseLoopDefinition(lexerArray []string, set *compilerSet) error {
	cursor := 0
	for cursor < len(lexerArray) {
		token := lexerArray[cursor]
		if token != tokenLoop {
			token++
			continue
		}

		if cursor == len(lexerArray)-1 {
			return errors.New("Expected \"{\", got newline")
		}

		if lexerArray[cursor+1] != tokenOpenCurly {
			return errors.New("Expected \"{\", got \"" +
				lexerArray[cursor+1] + "\"")
		}

		pointerKey := substitutionPushNextPointer +
			substitutionLoopStart + strconv.Itoa(set.keyIota)

		instructions := set.blocks[set.currentBlock].instructions
		instructions[len(instructions)-1].pointerKey = pointerKey

		consumeLexerArray(lexerArray, cursor, 2)

		cursor += 2
	}
}
