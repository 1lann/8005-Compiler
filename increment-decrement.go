package main

import (
	"errors"
	"strconv"
)

func parseIncrementDecrement(lexerArray []string, cursor *int, set *compilerSet) error {
	token := lexerArray[*cursor]
	if token != tokenR0 && token != tokenR1 {
		return nil
	}

	if len(lexerArray)-*cursor <= 1 {
		return nil
	}

	if lexerArray[*cursor+1] == tokenIncrement {
		if token == tokenR0 {
			set.appendInstruction(instruction{value: 1})
			consumeLexerArray(lexerArray, cursor, 2)
			return nil
		} else {
			set.appendInstruction(instruction{value: 3})
			consumeLexerArray(lexerArray, cursor, 2)
			return nil
		}
	}

	if lexerArray[*cursor+1] == tokenDecrement {
		if token == tokenR0 {
			set.appendInstruction(instruction{value: 2})
			consumeLexerArray(lexerArray, cursor, 2)
			return nil
		} else {
			set.appendInstruction(instruction{value: 4})
			consumeLexerArray(lexerArray, cursor, 2)
			return nil
		}
	}

	if lexerArray[*cursor+1] != tokenSelfAdd &&
		lexerArray[*cursor+1] != tokenSelfSub {
		return nil
	}

	if len(lexerArray)-*cursor <= 2 {
		return errors.New("Missing value to increment/decrement by")
	}

	if lexerArray[*cursor+2] == tokenR0 {
		return errors.New("You cannot increment/decrement by r0, try using " +
			":swap instead")
	}

	if lexerArray[*cursor+2] == tokenR1 {
		if token == tokenR1 {
			return errors.New("You cannot increment/decrement r1 by r1, " +
				"try using r0 = r1 instead")
		}

		if lexerArray[*cursor+1] == tokenSelfAdd {
			set.appendInstruction(instruction{value: 5})
			consumeLexerArray(lexerArray, cursor, 3)
			return nil
		} else {
			set.appendInstruction(instruction{value: 6})
			consumeLexerArray(lexerArray, cursor, 3)
			return nil
		}
	}

	if !isNumber(lexerArray[*cursor+2]) {
		return errors.New("Increment/decrement value must be a number or r1")
	}

	if lexerArray[*cursor+2] == "0" {
		return errors.New("Really? You want to increment/decrement by 0?")
	}

	incrementer := 0

	if token == tokenR0 {
		incrementer = 1
	} else {
		incrementer = 3
	}

	if lexerArray[*cursor+1] == tokenSelfSub {
		incrementer++
	}

	num, _ := strconv.Atoi(lexerArray[*cursor+2])

	for i := 0; i < num; i++ {
		set.appendInstruction(instruction{value: incrementer})
	}

	consumeLexerArray(lexerArray, cursor, 3)

	return nil
}
