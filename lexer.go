package eightc

const (
	tokenR0             = "TOKEN_R0"
	tokenR1             = "TOKEN_R1"
	tokenStore          = "TOKEN_STORE"
	tokenLoad           = "TOKEN_LOAD"
	tokenSelfMult       = "TOKEN_SELF_MULT"
	tokenSelfAdd        = "TOKEN_SELF_ADD"
	tokenSelfSub        = "TOKEN_SELF_SUB"
	tokenMult           = "TOKEN_MULT"
	tokenAdd            = "TOKEN_ADD"
	tokenSub            = "TOKEN_SUB"
	tokenDiv            = "TOKEN_DIV"
	tokenIf             = "TOKEN_IF"
	tokenLoop           = "TOKEN_LOOP"
	tokenBreak          = "TOKEN_BREAK"
	tokenReturn         = "TOKEN_RETURN"
	tokenEquals         = "TOKEN_EQUALS"
	tokenNotEquals      = "TOKEN_NOT_EQUALS"
	tokenIncrement      = "TOKEN_INCREMENT"
	tokenDecrement      = "TOKEN_DECREMENT"
	tokenFunction       = "TOKEN_FUNCTION_CALL"
	tokenSet            = "TOKEN_SET"
	tokenOpenCurly      = "TOKEN_OPEN_CURLY"
	tokenCloseCurly     = "TOKEN_CLOSE_CURLY"
	tokenOpenBracket    = "TOKEN_OPEN_BRACKET"
	tokenCloseBracket   = "TOKEN_CLOSE_BRACKET"
	tokenDefine         = "TOKEN_DEFINE"
	tokenQuote          = "TOKEN_QUOTE"
	tokenFunctionDefine = "TOKEN_FUNCTION_DEFINE"
	tokenConsumed       = "TOKEN_CONSUMED" // Seperator used by the compiler
)

var tokenMapping map[string]string
var identifierMapping map[string]string
var singleTokenMapping map[string]string

func init() {
	tokenMapping = make(map[string]string)
	tokenMapping["->"] = tokenStore
	tokenMapping["<-"] = tokenLoad
	tokenMapping["*="] = tokenSelfMult
	tokenMapping["-="] = tokenSelfSub
	tokenMapping["+="] = tokenSelfAdd
	tokenMapping["++"] = tokenIncrement
	tokenMapping["--"] = tokenDecrement
	tokenMapping["=="] = tokenEquals
	tokenMapping["!="] = tokenNotEquals

	identifierMapping = make(map[string]string)
	identifierMapping["if"] = tokenIf
	identifierMapping["loop"] = tokenLoop
	identifierMapping["for"] = tokenLoop
	identifierMapping["break"] = tokenBreak
	identifierMapping["return"] = tokenReturn
	identifierMapping["#define"] = tokenDefine
	identifierMapping["func"] = tokenFunctionDefine

	singleTokenMapping = make(map[string]string)
	singleTokenMapping["-"] = tokenSub
	singleTokenMapping["+"] = tokenAdd
	singleTokenMapping["*"] = tokenMult
	singleTokenMapping["="] = tokenSet
	singleTokenMapping[":"] = tokenFunction
	singleTokenMapping["{"] = tokenOpenCurly
	singleTokenMapping["}"] = tokenCloseCurly
	singleTokenMapping["("] = tokenOpenBracket
	singleTokenMapping[")"] = tokenCloseBracket
	singleTokenMapping["\""] = tokenQuote
}

func getTokens(str string) ([]string, bool) {
	tokenList := []string{}
	comment := false

	for len(str) >= 2 {
		matched := false

		for keyword, token := range identifierMapping {
			if str == keyword {
				tokenList = append(tokenList, token)
				str = str[len(keyword):]
				matched = true
				break
			}
		}

		if matched {
			continue
		}

		substr2 := str[:2]

		for keyword, token := range tokenMapping {
			if substr2 == keyword {
				tokenList = append(tokenList, token)
				str = str[len(keyword):]
				matched = true
				break
			}
		}

		if matched {
			continue
		}

		substr1 := str[:1]

		for keyword, token := range singleTokenMapping {
			if substr1 == keyword {
				tokenList = append(tokenList, token)
				str = str[len(keyword):]
				matched = true
				break
			}
		}

		if matched {
			continue
		}

		subAlphaNum := ""

		for i := 0; i < len(str); i++ {
			if (str[i] >= 'a' && str[i] <= 'z') ||
				(str[i] >= 'A' && str[i] <= 'Z') ||
				(str[i] >= '0' && str[i] <= '9') {
				subAlphaNum += string(str[i])
			} else {
				break
			}
		}

		if subAlphaNum == "r0" {
			tokenList = append(tokenList, tokenR0)
		} else if subAlphaNum == "r1" {
			tokenList = append(tokenList, tokenR1)
		} else if len(subAlphaNum) > 0 {
			tokenList = append(tokenList, subAlphaNum)
		}

		if str == "//" {
			comment = true
			break
		}

		str = str[len(subAlphaNum):]

		if len(subAlphaNum) <= 0 {
			if str[0] == ';' {
				str = str[1:]
			} else {
				// fmt.Println("Lexer warning: Skipping character \"" + str[:1] + "\"")
				tokenList = append(tokenList, str[:1])
				str = str[1:]
			}
		}
	}

	if !comment {
		for keyword, token := range singleTokenMapping {
			if str == keyword {
				tokenList = append(tokenList, token)
				str = str[len(keyword):]
				break
			}
		}

		subAlphaNum := ""

		for i := 0; i < len(str); i++ {
			if (str[i] >= 'a' && str[i] <= 'z') ||
				(str[i] >= 'A' && str[i] <= 'Z') ||
				(str[i] >= '0' && str[i] <= '9') {
				subAlphaNum += string(str[i])
			} else {
				break
			}
		}

		if len(subAlphaNum) > 0 {
			tokenList = append(tokenList, subAlphaNum)
		}
	}

	return tokenList, comment
}
