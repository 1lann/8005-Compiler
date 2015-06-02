package eightc

import (
	"errors"
	"strconv"
)

func parseIfDefinition(lexerArray []string, cursor *int, set *compilerSet) error {
	token := lexerArray[*cursor]
	if token != tokenIf {
		return nil
	}

	if len(lexerArray)-*cursor <= 4 {
		return errors.New("Invalid if statement syntax")
	}

	matchCloseBracket := false

	if lexerArray[*cursor+1] == tokenOpenBracket {
		consumeLexerArray(lexerArray, cursor, 1)
		matchCloseBracket = true
		if len(lexerArray)-*cursor <= 5 {
			return errors.New("Invalid if statement syntax")
		}
	}

	if lexerArray[*cursor+1] != tokenR0 {
		return errors.New("If statements can only compare r0")
	}

	if lexerArray[*cursor+3] != "0" {
		return errors.New("If statements can only be compared with 0")
	}

	isEquals := (lexerArray[*cursor+2] == tokenEquals)

	if lexerArray[*cursor+2] != tokenEquals &&
		lexerArray[*cursor+2] != tokenNotEquals {
		return errors.New("Compare operators can only be == and !=")
	}

	if matchCloseBracket && lexerArray[*cursor+4] != tokenCloseBracket {
		return errors.New("Missing matching closing bracket")
	} else if matchCloseBracket {
		consumeLexerArray(lexerArray, cursor, 5)
	} else {
		consumeLexerArray(lexerArray, cursor, 4)
	}

	if lexerArray[*cursor] != tokenOpenCurly {
		return errors.New("Expected \"{\", got \"" +
			lexerArray[*cursor] + "\"")
	}

	key := substitiutionConditionExit + strconv.Itoa(set.stepUpIota())
	set.parentTypes = append(set.parentTypes, blockIf)

	if isEquals {
		set.appendInstruction(instruction{value: 8})
		set.appendInstruction(instruction{key: key})
	} else {
		set.appendInstruction(instruction{value: 9})
		set.appendInstruction(instruction{key: key})
	}

	consumeLexerArray(lexerArray, cursor, 1)

	return nil
}
