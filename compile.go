package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
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

type compilerSet struct {
	instructions [256]int
	blocks       []block
	cursor       int
	newIota      int
	currentIota  int
	defineMap    map[string]int // Maps #defines
	currentLine  int
	filename     string
	currentBlock int
	parentBlocks []int
	parentIota   []int
}

const (
	stopInst = iota
	incrementR0Inst
	incrementR1Inst
	decrementR0Inst
	decrementR1Inst
	equalsZeroGotoInst
	notEqualsZeroGotoInst
)

// Disallow expressions because "they are too hard" ;)
// Setting/multing/adding r0
// That's it!
// We should also add #defines.

func isNumber(str string) bool {
	result, err := strconv.Atoi(str)
	if err != nil {
		return false
	}

	if strconv.Itoa(result) == str {
		return true
	}

	return false
}

func isToken(str string) bool {
	if len(str) > 6 && str[:6] == "TOKEN_" {
		return true
	}
	return false
}

func compileLine(line string, set *compilerSet) error {
	line = strings.Replace(line, "\t", " ", -1)

	parts := strings.Split(strings.ToLower(line), " ")

	lexerArray := make([]string, 0)

	for _, part := range parts {
		tokens := getTokens(part)
		lexerArray = append(lexerArray, tokens...)
	}

	if set.currentBlock < 0 {
		if err := mapDefines(lexerArray, set); err != nil {
			return err
		}

		if err := parseFunctionDefinition(lexerArray, set); err != nil {
			return err
		}
	} else {
		if err := parseStaticSet(lexerArray, set); err != nil {
			return err
		}
	}

	for i := 0; i < len(lexerArray); i++ {
		if lexerArray[i] != tokenConsumed {
			return errors.New("Invalid syntax at: " + lexerArray[i])
		}
	}

	return nil
}

func consumeLexerArray(lexerArray []string, cursor int, length int) {
	for i := cursor; i < length; i++ {
		lexerArray[i] = tokenConsumed
	}
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

func (set *compilerSet) getCurrentIota() int {

}

func (set *compilerSet) stepUpIota() int {

}

func (set *compilerSet) setupDownIota() {
	set.currentIota = set.parentIota[len(set.parentIota)-1]
	set.parentIota = set.parentIota[:len(set.parentIota)-1]
}

func (set *compilerSet) getUniqueIota() int {
	uniqueIota := set.newIota
	set.newIota++
	return uniqueIota
}

func compile(filename string, reader *bufio.Reader) ([256]int, error) {
	compiler := compilerSet{}
	compiler.filename = filename
	compiler.currentBlock = 0
	compiler.defineMap = make(map[string]int)
	compiler.blocks = append(compiler.blocks, block{name: "main"})

	for {
		lineData, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				if compiler.currentBlock >= 0 {
					return [256]int{},
						errors.New("Reached end of file without closing " +
							"block (You're missing a \"}\" somewhere)")
				}

				assembleProgram(&compiler)

				break
			}

			fmt.Println("Error reading file:", err)
			return [256]int{}, err
		}

		line := string(lineData)
		compiler.currentLine++

		if err := compileLine(line, &compiler); err != nil {
			err = errors.New(compiler.filename + ":" +
				strconv.Itoa(compiler.currentLine) + ": " + err.Error())
			return [256]int{}, err
		}
	}

	return compiler.instructions, nil
}
