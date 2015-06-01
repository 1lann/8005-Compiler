package main

func parseFunctionReturn(lexerArray []string, cursor *int, set *compilerSet) error {
	token := lexerArray[*cursor]
	if token != tokenReturn {
		return nil
	}

	key := substitutionFrameReturn + set.blocks[set.currentBlock].name
	set.appendInstruction(instruction{value: 10})
	set.appendInstruction(instruction{value: 0})
	set.appendInstruction(instruction{value: 9})
	set.appendInstruction(instruction{key: key})

	consumeLexerArray(lexerArray, cursor, 1)

	return nil
}
