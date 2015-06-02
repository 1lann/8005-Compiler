package eightc

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	substitutionFrame       = "SUBSTITUTION_FRAME:"        // Value inside frame
	substitutionFrameReturn = "SUBSTITUTION_FRAME_RETURN:" // For returns inside frames

	substitutionFunctionReturn = "SUBSTITUTION_FUNCTION_RETURN:" // Not even in the frame

	substitiutionConditionExit = "SUBSTITUTION_CONDITION_EXIT:"
	substitutionGotoFunction   = "SUBSTITUTION_GOTO_FUNCTION:"
	substitutionLoopStart      = "SUBSTITUTION_LOOP_START:"

	substitutionVariable = "SUBSTITUTION_VARIABLE:"
)

type instruction struct {
	value int
	// substitutionType string
	line           int      // Corresponding line number
	key            string   // SubstitutionType:Key
	pointerKey     []string // SubstitutionType:Key
	pushPointerKey []string // SubstitutionType:Key
}

type block struct {
	name string
	// Return point locations are frameLocation + 1 and frameLocation + 3
	instructions []instruction
}

func (set *compilerSet) warn(warning string) {
	fmt.Println(set.filename + ":" + strconv.Itoa(set.currentLine) + ": " +
		"warning: " + warning)
}

func (set *compilerSet) appendInstruction(instruct instruction) {
	instruct.line = set.currentLine
	set.blocks[set.currentBlock].instructions = append(
		set.blocks[set.currentBlock].instructions,
		instruct,
	)
}

func (set *compilerSet) addPointerKey(key string) {
	instructions := set.blocks[set.currentBlock].instructions
	instructions[len(instructions)-1].pointerKey = append(
		instructions[len(instructions)-1].pointerKey,
		key,
	)
}

func (set *compilerSet) pushPointerKey(key string) {
	instructions := set.blocks[set.currentBlock].instructions
	instructions[len(instructions)-1].pushPointerKey = append(
		instructions[len(instructions)-1].pushPointerKey,
		key,
	)
}

func (set *compilerSet) stepUpIota() int {
	newIota := set.getUniqueIota()
	set.parentIota = append(set.parentIota, set.currentIota)
	set.currentIota = newIota
	return newIota
}

func (set *compilerSet) stepDownIota() error {
	if len(set.parentIota) < 1 {
		return errors.New("Unexpected \"}\" not inside block")
	}
	set.currentIota = set.parentIota[len(set.parentIota)-1]
	set.parentIota = set.parentIota[:len(set.parentIota)-1]

	return nil
}

func (set *compilerSet) getUniqueIota() int {
	uniqueIota := set.newIota
	set.newIota++
	return uniqueIota
}
