package eightc

import (
	"errors"
	"strconv"
)

func assembleProgram(set *compilerSet) error {
	// Assemble the program
	allInstructions := make([]instruction, 256)

	// "Flattens" instructions
	i := 0
	for _, block := range set.blocks {
		for k, instruct := range block.instructions {
			if i >= 255 {
				return errors.New(" The program does not fit in the 256 " +
					"byte instruction limit!")
			}
			allInstructions[i] = instruct
			i++

			if k == len(block.instructions)-1 &&
				len(instruct.pushPointerKey) > 0 {
				// Has a push pointer key at end of block, add a padding
				allInstructions[i] = instruction{value: 0, line: instruct.line}
				i++
			}
		}
	}

	// Build pointer key database
	pointerKeys := make(map[string]int)
	pushPointerKeys := []string{}
	for i, instruct := range allInstructions {
		for _, key := range instruct.pointerKey {
			if _, found := pointerKeys[key]; found {
				return errors.New(strconv.Itoa(instruct.line) + ": Duplicate " +
					"pointer key \"" + key + "\"")
			}

			pointerKeys[key] = i
		}

		for _, pushKey := range pushPointerKeys {
			if _, found := pointerKeys[pushKey]; found {
				return errors.New(strconv.Itoa(instruct.line) + ": Duplicate " +
					"pointer key \"" + pushKey + "\"")
			}
			pointerKeys[pushKey] = i
		}

		pushPointerKeys = []string{}
		for _, pushKey := range instruct.pushPointerKey {
			pushPointerKeys = append(pushPointerKeys, pushKey)
		}
	}

	// Resolve keys and form array
	for i, instruct := range allInstructions {
		set.instructions[i] = instruct.value

		if len(instruct.key) == 0 {
			continue
		}

		if address, found := pointerKeys[instruct.key]; !found {
			return errors.New(strconv.Itoa(instruct.line) + ": Could not " +
				"resolve key \"" + instruct.key + "\"")
		} else {
			set.instructions[i] = address
		}
	}

	return nil
}
