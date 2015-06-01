package main

import "strconv"

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
		blockInstructions := set.blocks[set.currentBlock].instructions
		switch thisBlock.blockType {
		// Return on functions
		case blockFunction:
			key := substitutionFrameReturn + thisBlock.name
			set.appendInstruction(instruction{value: 10})
			set.appendInstruction(instruction{value: 0})
			set.appendInstruction(instruction{value: 9})
			set.appendInstruction(instruction{value: 0, key: key})

			set.currentBlock = set.parentBlocks[len(set.parentBlocks)-1]
			set.parentBlocks = set.parentBlocks[:len(set.parentBlocks)-1]

		case blockIf:
			pushPointerKey := substitiutionConditionExit +
				strconv.Itoa(set.keyIota)

			set.keyIota++
			set.pushPointerKey(pushPointerKey)
		// Loops are embedded, go back to the start
		case blockLoop:
			key := substitutionLoopStart + strconv.Itoa(set.keyIota)
			pushPointerKey := substitutionLoopStart +
				strconv.Itoa(set.keyIota)
			set.keyIota++

			set.appendInstruction(instruction{value: 8})
			set.appendInstruction(instruction{value: 0, key: key})
			set.appendInstruction(instruction{value: 9})
			set.appendInstruction(instruction{value: 0, key: key})
			set.pushPointerKey(pushPointerKey)
		}

		consumeLexerArray(lexerArray, cursor, 1)
		cursor++
	}

}
