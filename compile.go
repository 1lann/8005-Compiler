func compileLine(line string, set *compilerSet) {
	line = strings.Replace(line, "\t", " ", -1)

	parts := strings.Split(strings.ToLower(line), " ")

	lexerArray := make([]string, 0)

	ip := 0

	for _, part := range parts {
		tokens = getTokens(part)
		lexerArray = append(lexerArray, tokens...)
	}
}

func compile(reader *bufio.Reader) ([255]int, error) {
	compilerSet

	for {
		lineData, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println("Error reading file:", err)
			return instructions, err
		}

		line := string(lineData)
		compileLine(line, &set)
	}
}
