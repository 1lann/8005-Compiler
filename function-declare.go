package eightc

import (
	"errors"
)

func processFunctionDefinition(name string, set *compilerSet) error {
	blockNum := len(set.blocks)

	newBlock := block{}
	newBlock.name = name

	set.blocks = append(set.blocks, newBlock)
	set.parentBlocks = append(set.parentBlocks, set.currentBlock)
	set.currentBlock = blockNum

	set.parentTypes = append(set.parentTypes, blockFunction)

	returnKey := substitutionFrame + name
	pushPointerKey := substitutionGotoFunction + name
	gotoKey := substitutionFrameReturn + name

	set.appendInstruction(instruction{value: 10})
	set.addPointerKey(gotoKey)
	set.appendInstruction(instruction{value: 0})
	set.appendInstruction(instruction{value: 9})
	set.appendInstruction(instruction{value: 0})
	set.addPointerKey(returnKey)
	set.pushPointerKey(pushPointerKey)

	return nil
}

func parseFunctionDefinition(lexerArray []string, cursor *int, set *compilerSet) error {
	token := lexerArray[*cursor]
	if token != tokenFunctionDefine {
		return nil
	}

	if len(lexerArray)-*cursor <= 2 {
		return errors.New("Invalid function definition, note that " +
			"the \"{\" must be on the same line as the function " +
			"definition")
	}

	if isNumber(lexerArray[*cursor+1]) {
		return errors.New("Function definition name cannot be a number")
	}

	if isToken(lexerArray[*cursor+1]) {
		return errors.New("Expected function name, got \"" +
			lexerArray[*cursor+1] + "\"")
	}

	if lexerArray[*cursor+2] != tokenOpenCurly {
		return errors.New("Expected \"{\", got \"" +
			lexerArray[*cursor+2] + "\"")
	}

	for _, block := range set.blocks {
		if block.name == lexerArray[*cursor+2] {
			return errors.New("Attempt to redefine function \"" +
				lexerArray[*cursor+2] + "\"")
		}
	}

	err := processFunctionDefinition(lexerArray[*cursor+1], set)
	if err != nil {
		return err
	}

	consumeLexerArray(lexerArray, cursor, 3)

	return nil
}
