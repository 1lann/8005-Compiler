package main

import (
	"errors"
	"strconv"
)

func writeStaticSetterInstruction(rtoken string, value int, set *compilerSet) {
	blockCursor := set.blocks[set.currentBlock].cursor
	blockInstructions := set.blocks[set.currentBlock].instructions
	if rtoken == tokenR0 {
		blockInstructions[blockCursor] = 10
		blockInstructions[blockCursor+1] = value
		set.blocks[set.currentBlock].cursor += 2
	} else if rtoken == tokenR1 {
		blockInstructions[blockCursor] = 11
		blockInstructions[blockCursor] = value
		set.blocks[set.currentBlock].cursor += 2
	}
}

func parseStaticSet(lexerArray []string, set *compilerSet) error {
	cursor := 0
	for cursor < len(lexerArray) {
		token := lexerArray[cursor]
		if token == tokenR0 || token == tokenR1 {
			cursor++
			continue
		}

		if len(lexerArray)-cursor <= 1 {
			return errors.New("Unknown r0/r1 use")
		}

		if lexerArray[cursor+1] != tokenSet {
			cursor++
			continue
		}

		if len(lexerArray)-cursor <= 2 {
			return errors.New("Missing value to set register to")
		}

		if isToken(lexerArray[cursor+2]) {
			return errors.New("Setter got unexpected token: " + lexerArray[cursor+2])
		} else if isNumber(lexerArray[cursor+2]) {
			num, _ := strconv.Atoi(lexerArray[cursor+2])
			writeStaticSetterInstruction(token, num, set)
		} else {
			// Is either a #define or variable
			deref := lexerArray[cursor+2]
			if value, found := set.defineMap[deref]; found {
				writeStaticSetterInstruction(token, value, set)
			} else {
				return errors.New("Unknown dereferencable object \"" +
					deref + "\". Did you mean to " +
					"use a key-store variable through \"<-\" instead?")
			}
		}

		consumeLexerArray(lexerArray, cursor, 3)
		cursor += 3
	}

	return nil
}
