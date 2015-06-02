package eightc

import (
	"errors"
	"strconv"
)

func writeStaticSetterInstruction(rtoken string, value int, set *compilerSet) {
	if rtoken == tokenR0 {
		set.appendInstruction(instruction{value: 10})
		set.appendInstruction(instruction{value: value})
	} else if rtoken == tokenR1 {
		set.appendInstruction(instruction{value: 11})
		set.appendInstruction(instruction{value: value})
	}
}

func parseStaticSet(lexerArray []string, cursor *int, set *compilerSet) error {
	token := lexerArray[*cursor]
	if token != tokenR0 && token != tokenR1 {
		return nil
	}

	if *cursor == len(lexerArray)-1 {
		return errors.New("Invalid r0/r1 use")
	}

	if lexerArray[*cursor+1] != tokenSet {
		return nil
	}

	if len(lexerArray)-*cursor <= 2 {
		return errors.New("Missing value to set register to")
	}

	if lexerArray[*cursor+2] == tokenR0 || lexerArray[*cursor+2] == tokenR1 {
		if lexerArray[*cursor+2] == tokenR0 {
			set.appendInstruction(instruction{value: 12})
		} else {
			set.appendInstruction(instruction{value: 13})
		}
		set.appendInstruction(instruction{key: substitutionVariable + strconv.Itoa(set.tempMemory)})

		if token == tokenR0 {
			set.appendInstruction(instruction{value: 14})
		} else {
			set.appendInstruction(instruction{value: 15})
		}
		set.appendInstruction(instruction{key: substitutionVariable + strconv.Itoa(set.tempMemory)})
	} else if isToken(lexerArray[*cursor+2]) {
		return errors.New("Setter got unexpected token: " +
			lexerArray[*cursor+2])
	} else if isNumber(lexerArray[*cursor+2]) {
		num, _ := strconv.Atoi(lexerArray[*cursor+2])
		writeStaticSetterInstruction(token, num, set)
	} else {
		// Is either a #define or variable
		deref := lexerArray[*cursor+2]
		if value, found := set.defineMap[deref]; found {
			writeStaticSetterInstruction(token, value, set)
		} else {
			return errors.New("Unknown dereferencable object \"" +
				deref + "\". Did you mean to " +
				"use a key-store variable through \"<-\" instead?")
		}
	}

	consumeLexerArray(lexerArray, cursor, 3)

	return nil
}
