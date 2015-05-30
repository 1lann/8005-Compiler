package main

func parseBlockClosure(lexerArray []string, set *compilerSet) error {
	cursor := 0
	for cursor < len(lexerArray) {
		token := lexerArray[cursor]
		if token != tokenCloseCurly {
			cursor++
			continue
		}

		if set.currentBlock < 0 {
			return errors.New("You cannot close a block you haven't " +
				" opened (Missing)")
		}

		// Clean up
		// Jobs: Add return statement
		// Unwind into embedded block if necessary (loops only I think)

		thisBlock := set.blocks[set.currentBlock]
		thisInstructionSet := set.blocks[set.currentBlock].instructions
		thisCursor := set.blocks[set.currentBlock].cursor
		switch thisBlock.blockType {
		// Return on functions
		case blockFunction:
			thisInstructionSet[thisCursor] = 8
			thisInstructionSet[thisCursor+1] = 0
			thisBlock.name
		// Ifs are virtually functions, return too
		case blockIf:

		// Loops are embedded, go back to the start
		case blockLoop:

		}

		set.blocks[set.currentBlock].cursor = thisCursor
	}

}
