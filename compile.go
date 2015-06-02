package eightc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	blockFunction = iota
	blockIf
	blockLoop
)

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
	currentType  int
	parentBlocks []int
	parentIota   []int
	parentTypes  []int
	tempMemory   int
}

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

func convertChars(str string) string {
	currentRune := ' '

	for currentRune <= '~' {
		str = strings.Replace(str, "'"+string(currentRune)+"'",
			strconv.Itoa(int(currentRune)), -1)

		currentRune++
	}

	return str
}

func convertNegative(str string) string {
	for i := -99; i < 0; i++ {
		str = strings.Replace(str, strconv.Itoa(i), strconv.Itoa(256+i), -1)
	}

	return str
}

func compileLine(line string, set *compilerSet) error {
	line = strings.Replace(line, "\t", " ", -1)
	line = convertChars(line)
	line = convertNegative(line)

	parts := strings.Split(line, " ")

	lexerArray := make([]string, 0)

	for _, part := range parts {
		part = strings.Replace(part, "\n", "", 1)
		tokens, comment := getTokens(part)
		lexerArray = append(lexerArray, tokens...)
		if comment {
			break
		}
	}

	lastCursor := -1
	lexerCursor := 0

	for lexerCursor != lastCursor && lexerCursor < len(lexerArray) {
		lastCursor = lexerCursor

		if set.currentBlock == 0 {
			if err := mapDefines(lexerArray, &lexerCursor, set); err != nil {
				return err
			}
			if lastCursor != lexerCursor {
				continue
			}

			if err := parseFunctionDefinition(lexerArray, &lexerCursor, set); err != nil {
				return err
			}
			if lastCursor != lexerCursor {
				continue
			}

		}

		if err := parseStaticSet(lexerArray, &lexerCursor, set); err != nil {
			return err
		}
		if lastCursor != lexerCursor {
			continue
		}

		if err := parseIncrementDecrement(lexerArray, &lexerCursor, set); err != nil {
			return err
		}
		if lastCursor != lexerCursor {
			continue
		}

		if err := parseFunctionCall(lexerArray, &lexerCursor, set); err != nil {
			return err
		}
		if lastCursor != lexerCursor {
			continue
		}

		if err := parseGetVariable(lexerArray, &lexerCursor, set); err != nil {
			return err
		}
		if lastCursor != lexerCursor {
			continue
		}

		if err := parseSetVariable(lexerArray, &lexerCursor, set); err != nil {
			return err
		}
		if lastCursor != lexerCursor {
			continue
		}

		if err := parseIfDefinition(lexerArray, &lexerCursor, set); err != nil {
			return err
		}
		if lastCursor != lexerCursor {
			continue
		}

		if err := parseLoopDefinition(lexerArray, &lexerCursor, set); err != nil {
			return err
		}
		if lastCursor != lexerCursor {
			continue
		}

		if err := parseBlockClosure(lexerArray, &lexerCursor, set); err != nil {
			return err
		}
		if lastCursor != lexerCursor {
			continue
		}

		if err := parseLoopBreak(lexerArray, &lexerCursor, set); err != nil {
			return err
		}
		if lastCursor != lexerCursor {
			continue
		}

		if err := parseFunctionReturn(lexerArray, &lexerCursor, set); err != nil {
			return err
		}
		if lastCursor != lexerCursor {
			continue
		}
	}

	for i := 0; i < len(lexerArray); i++ {
		if lexerArray[i] != tokenConsumed {
			return errors.New("Invalid syntax at: " + lexerArray[i])
		}
	}

	return nil
}

func consumeLexerArray(lexerArray []string, cursor *int, length int) {
	for i := *cursor; i < *cursor+length; i++ {
		lexerArray[i] = tokenConsumed
	}

	*cursor += length
}

func Compile(filename string, reader *bufio.Reader) ([256]int, error) {
	compiler := compilerSet{}
	compiler.filename = filename
	compiler.currentBlock = 0
	compiler.defineMap = make(map[string]int)
	compiler.blocks = append(compiler.blocks, block{name: "main"})

	// TODO Can be replaced by an allocate memory function
	compiler.tempMemory = compiler.getUniqueIota()
	tempMemoryBlock := block{}
	tempMemoryBlock.instructions = append(tempMemoryBlock.instructions,
		instruction{value: 0})
	compiler.blocks = append(compiler.blocks, tempMemoryBlock)
	compiler.currentBlock = 1
	compiler.addPointerKey(substitutionVariable +
		strconv.Itoa(compiler.tempMemory))
	compiler.currentBlock = 0

	allErrors := []error{}

	for {
		lineData, err := reader.ReadBytes('\n')

		line := string(lineData)
		compiler.currentLine++

		if err := compileLine(line, &compiler); err != nil {
			err = errors.New(compiler.filename + ":" +
				strconv.Itoa(compiler.currentLine) + ": " + err.Error())
			allErrors = append(allErrors, err)
		}

		if err != nil {
			if err == io.EOF {
				if compiler.currentBlock > 0 {
					allErrors = append(allErrors,
						errors.New(compiler.filename+": Reached end of "+
							"file without closing block (You're missing a "+
							"\"}\" somewhere)"),
					)
				}

				if len(allErrors) > 10 {
					concat := ""

					for i := 0; i < 10; i++ {
						concat += allErrors[i].Error() + "\n"
					}

					concat += "and " + strconv.Itoa(len(allErrors)-10) +
						" additional errors..."

					return [256]int{},
						errors.New(concat)
				} else if len(allErrors) > 0 {
					concat := ""

					for i := 0; i < len(allErrors); i++ {
						concat += allErrors[i].Error() + "\n"
					}

					concat = concat[:len(concat)-1]

					return [256]int{},
						errors.New(concat)
				}

				compiler.appendInstruction(instruction{value: 0})
				compiler.addPointerKey(substitutionFrameReturn + "main")

				if err := assembleProgram(&compiler); err != nil {
					return [256]int{},
						errors.New(compiler.filename + ":" + err.Error())
				}

				break
			}

			fmt.Println("Error reading file:", err)
			return [256]int{}, err
		}
	}

	return compiler.instructions, nil
}
