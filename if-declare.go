package main

import "strconv"

func parseIfDefinition(lexerArray []string, set *compilerSet) error {
	cursor := 0
	for cursor < len(lexerArray) {
		token := lexerArray[cursor]
		if token != tokenIf {
			cursor++
			continue
		}

		if len(lexerArray)-cursor <= 4 {
			return errors.New("Invalid if statement syntax")
		}

		matchCloseBracket := false
		startCursor := cursor

		if lexerArray[cursor+1] == tokenOpenBracket {
			cursor++
			matchCloseBracket = true
			if len(lexerArray)-cursor <= 5 {
				return errors.New("Invalid if statement syntax")
			}
		}

		if lexerArray[cursor+1] != tokenR0 {
			return errors.New("If statements can only compare r0")
		}

		if lexerArray[cursor+3] != "0" {
			return errors.New("If statements can be compared with 0")
		}

		isEquals := (lexerArray[cursor+2] == tokenEquals)

		if lexerArray[cursor+2] != tokenEquals &&
			lexerArray[cursor+2] != tokenNotEquals {
			return errors.New("Compare operators can only be == and !=")
		}

		if matchCloseBracket && lexerArray[cursor+4] != tokenCloseBracket {
			return errors.New("Missing matching closing bracket")
		} else if matchCloseBracket {
			cursor += 5
		} else {
			cursor += 4
		}

		if lexerArray[cursor] != tokenOpenCurly {
			return errors.New("Expected \"{\", got \"" +
				lexerArray[cursor] + "\"")
		}

		key := substitiutionConditionExit + strconv.Itoa(set.keyIota)

		if isEquals {
			set.appendInstruction(instruction{value: 8})
			set.appendInstruction(instruction{value: 0, key: key})
		} else {
			set.appendInstruction(instruction{value: 9})
			set.appendInstruction(instruction{value: 0, key: key})
		}

		if matchCloseBracket {
			consumeLexerArray(lexerArray, startCursor, 7)
		} else {
			consumeLexerArray(lexerArray, startCursor, 5)
		}

		cursor++
	}
}
