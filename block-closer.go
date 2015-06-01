package main

import (
	"errors"
	"strconv"
)

func parseBlockClosure(lexerArray []string, cursor *int, set *compilerSet) error {
	token := lexerArray[*cursor]
	if token != tokenCloseCurly {
		return nil
	}

	if len(set.parentTypes) == 0 {
		return errors.New("You cannot close a block you haven't " +
			"yet opened (Unexpected \"}\")")
	}

	// Clean up
	// Jobs: Add return statement
	// Unwind into embedded block if necessary (loops only I think)

	thisBlock := set.blocks[set.currentBlock]
	parentType := set.parentTypes[len(set.parentTypes)-1]
	switch parentType {
	// Return on functions
	case blockFunction:
		key := substitutionFrameReturn + thisBlock.name
		set.appendInstruction(instruction{value: 10})
		set.appendInstruction(instruction{value: 0})
		set.appendInstruction(instruction{value: 9})
		set.appendInstruction(instruction{key: key})

		set.currentBlock = set.parentBlocks[len(set.parentBlocks)-1]
		set.parentBlocks = set.parentBlocks[:len(set.parentBlocks)-1]

	case blockIf:
		pushPointerKey := substitiutionConditionExit +
			strconv.Itoa(set.currentIota)

		if err := set.stepDownIota(); err != nil {
			return err
		}

		set.pushPointerKey(pushPointerKey)
	// Loops are embedded, go back to the start
	case blockLoop:
		key := substitutionLoopStart + strconv.Itoa(set.currentIota)
		pushPointerKey := substitiutionConditionExit +
			strconv.Itoa(set.currentIota)

		if err := set.stepDownIota(); err != nil {
			return err
		}

		set.appendInstruction(instruction{value: 8})
		set.appendInstruction(instruction{key: key})
		set.appendInstruction(instruction{value: 9})
		set.appendInstruction(instruction{key: key})
		set.pushPointerKey(pushPointerKey)
	default:
		return errors.New("Unknown closing block type")
	}

	set.parentTypes = set.parentTypes[:len(set.parentTypes)-1]

	consumeLexerArray(lexerArray, cursor, 1)

	return nil
}
