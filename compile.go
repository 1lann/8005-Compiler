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
	blockFunction = "BLOCK_FUNCTION"
	blockLoop     = "BLOCK_LOOP"
	blockIf       = "BLOCK_IF"
)

const (
	substitutionFrame        = "SUBSTITUTION_FRAME" // Frames for functions
	substitutionReturnFrameA = "SUBSTITUTION_RETURN_FRAME_A"
	substitutionReturnFrameB = "SUBSTITUTION_RETURN_FRAME_B"

	substitiutionLoopExit = "SUBSTITUTION_LOOP_EXIT" // Also used for else statements

	substitutionGotoFunction = "SUBSTITUTION_GOTO_FUNCTION"
	substitutionGotoIf       = "SUBSTITUTION_GOTO_IF" // Ifs are one way and have hard-coded frames
)

type instruction struct {
	value int
	// substitutionType string
	line       int // Corresponding line number
	key        string
	pointerKey string // SubstitutionType:Key
}

type block struct {
	blockType string
	name      string
	// Return point locations are frameLocation + 1 and frameLocation + 3
	instructions []instruction
	cursor       int
}

// type programObject struct {
// 	frameSubstitutions       map[string]int // Points to start of return frame
// 	returnPointSubstitutions map[string]int // Points to points inside the return frame
// 	gotoSubstitutions        map[string]int // Points to start of object
// 	block                    int            // For blocks
// 	location                 int            // For variables
// }

type compilerSet struct {
	instructions [256]int
	blocks       []block
	cursor       int
	keyIota      int
	// functionMap  map[string]programObject // Maps names to address substitutions required
	// variableMap  map[string]programObject // Maps names to address substitutions required
	defineMap    map[string]int // Maps #defines
	currentLine  int
	filename     string
	currentBlock int
	parentBlocks []int
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

func warn(set *compilerSet, warning string) {
	fmt.Println(set.filename + ":" + strconv.Itoa(set.currentLine) + ": " +
		"warning: " + warning)
}

func compile(filename string, reader *bufio.Reader) ([256]int, error) {
	compiler := compilerSet{}
	compiler.filename = filename
	compiler.currentBlock = 0
	compiler.functionMap = make(map[string]programObject)
	compiler.variableMap = make(map[string]programObject)
	compiler.defineMap = make(map[string]int)
	compiler.currentBlock = -1

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
