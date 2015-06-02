package eightc

import "strconv"

func parseLoopBreak(lexerArray []string, cursor *int, set *compilerSet) error {
	token := lexerArray[*cursor]
	if token != tokenBreak {
		return nil
	}

	loopKey := set.currentIota

	depth := 0

	if set.currentType == blockIf {
		depth++
	}

	if set.currentType != blockLoop {
		for i := len(set.parentTypes) - 1; i >= 0; i-- {
			if set.parentTypes[i] == blockLoop {
				break
			} else if set.parentTypes[i] == blockIf {
				depth++
			}
		}

		loopKey = set.parentIota[len(set.parentIota)-depth]
	}

	key := substitiutionConditionExit +
		strconv.Itoa(loopKey)

	set.appendInstruction(instruction{value: 8})
	set.appendInstruction(instruction{key: key})
	set.appendInstruction(instruction{value: 9})
	set.appendInstruction(instruction{key: key})

	consumeLexerArray(lexerArray, cursor, 1)

	return nil
}
