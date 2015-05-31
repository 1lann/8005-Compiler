package main

import (
	"errors"
)

func processFunctionDefinition(name string, set *compilerSet) error {
	object, found := set.functionMap[name]
	if found && object.block >= 0 {
		return errors.New("Function re-definition not allowed")
	}

	blockNum := len(set.blocks)
	object.block = blockNum

	newBlock := block{}
	newBlock.name = name

	set.blocks = append(set.blocks, newBlock)
	set.parentBlocks = append(set.parentBlocks, set.currentBlock)
	set.currentBlock = blockNum

	returnKey := substitutionFrame + name
	frameKey := substitutionPushNextPointer + substitutionGotoFunction + name
	gotoKey := substitutionFrameReturn + name

	set.appendInstruction(instruction{value: 10, pointerKey: gotoKey})
	set.appendInstruction(instruction{value: 0})
	set.appendInstruction(instruction{value: 9})
	set.appendInstruction(instruction{value: 10, pointerKey: returnKey})
	set.appendInstruction(instruction{value: 0, pointerKey: frameKey})

	return nil
}

func parseFunctionDefinition(lexerArray []string, set *compilerSet) error {
	cursor := 0
	for cursor < len(lexerArray) {
		token := lexerArray[cursor]
		if token != tokenFunctionDefine {
			cursor++
			continue
		}

		if len(lexerArray)-cursor <= 2 {
			return errors.New("Invalid function definition, note that " +
				"the \"{\" must be on the same line as the function " +
				"definition")
		}

		if isNumber(lexerArray[cursor+1]) {
			return errors.New("Function definition name cannot be a number")
		}

		if isToken(lexerArray[cursor+1]) {
			return errors.New("Expected function name, got \"" +
				lexerArray[cursor+1] + "\"")
		}

		if lexerArray[cursor+2] != tokenOpenCurly {
			return errors.New("Expected \"{\", got \"" +
				lexerArray[cursor+2] + "\"")
		}

		err := processFunctionDefinition(lexerArray[cursor+1], set)
		if err != nil {
			return err
		}

		consumeLexerArray(lexerArray, cursor, 3)
		cursor += 3
	}

	return nil
}
