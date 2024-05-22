package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Constants for delimiters and character sets
const (
	Delimiter     = ";"
	Alphabet      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Alphanumeric  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Numbers       = "0123456789"
	OperatorChars = "+-*/"
	Parentheses   = "()"
)

// Lexeme types
const (
	DelimiterType   = "delimiter"
	IdentifierType  = "identifier"
	ConstType       = "const"
	AssignmentType  = "assignment"
	OperatorType    = "operator"
	ParenthesesType = "parentheses"
	NumberType      = "number"
	ErrorType       = "error"
	CommentType     = "comment"
)

// TokenSaver struct to store and print tokens
type TokenSaver struct {
	tokens []Token
}

// Token represents a lexical token with its type and value
type Token struct {
	Type  string
	Value string
}

// Add adds a token to the TokenSaver
func (ts *TokenSaver) Add(tokenType, value string) {
	ts.tokens = append(ts.tokens, Token{Type: tokenType, Value: value})
}

// Print prints all tokens in the TokenSaver
func (ts *TokenSaver) Print() {
	fmt.Println("Tokens:")
	for _, token := range ts.tokens {
		fmt.Printf("{ Type: %s, Value: %s }\n", token.Type, token.Value)
	}
}

func main() {
	runLexer()
}

// runLexer processes input and generates tokens
func runLexer() {
	state := DelimiterType
	token := ""
	tokenSaver := &TokenSaver{}
	characters := readCharacters()

	for i := 0; i < len(characters); i++ {
		char := string(characters[i])
		switch state {
		case DelimiterType:
			handleDelimiter(char, tokenSaver)
			if char != " " && char != "\n" && char != "\r" {
				token = char
				state = getLexemeType(char)
			}
		case IdentifierType, NumberType, OperatorType, ParenthesesType:
			if isValid(char, state) {
				token += char
			} else {
				tokenSaver.Add(state, token)
				token = ""
				state = DelimiterType
				i-- // Re-process current character
			}
		case AssignmentType:
			handleAssignment(char, tokenSaver, &state, &token, &i, characters)
		case ConstType:
			handleConst(char, tokenSaver, &state, &token, &i, characters)
		case CommentType:
			if char == "\n" {
				state = DelimiterType
			}
		}
	}
	if token != "" && state != CommentType {
		tokenSaver.Add(state, token)
	}
	tokenSaver.Print()
}

func handleDelimiter(char string, tokenSaver *TokenSaver) {
	if char == Delimiter {
		tokenSaver.Add(DelimiterType, Delimiter)
	}
}

func getLexemeType(char string) string {
	if strings.Contains(Alphabet, char) {
		return IdentifierType
	} else if strings.Contains(Numbers, char) {
		return NumberType
	} else if char == "'" {
		return ConstType
	} else if char == "#" {
		return CommentType
	} else if strings.Contains(OperatorChars, char) {
		return OperatorType
	} else if strings.Contains(Parentheses, char) {
		return ParenthesesType
	} else if char == ";" {
		return DelimiterType
	} else if char == ":" {
		return AssignmentType
	}
	return ErrorType
}

func isValid(char string, state string) bool {
	switch state {
	case IdentifierType:
		return strings.Contains(Alphanumeric, char)
	case NumberType:
		return strings.Contains(Numbers, char)
	case ConstType:
		return char != "'"
	case OperatorType:
		return strings.Contains(OperatorChars, char)
	case ParenthesesType:
		return strings.Contains(Parentheses, char)
	}
	return false
}

func handleAssignment(char string, tokenSaver *TokenSaver, state *string, token *string, i *int, characters []byte) {
	if char == "=" && (*i)+1 < len(characters) && characters[(*i)-1] == ':' {
		*state = AssignmentType
		*token = ":="
		(*i)++ // Skip '=' character
		tokenSaver.Add(AssignmentType, *token)
		*state = DelimiterType
	} else {
		tokenSaver.Add(OperatorType, *token)
		*token = ""
		*state = ErrorType
	}
}

func handleConst(char string, tokenSaver *TokenSaver, state *string, token *string, i *int, characters []byte) {
	if strings.Contains(Alphabet, char) && checkConstBrackets(*i, characters) {
		*state = ConstType
		*token = fmt.Sprint("'", char, "'")
		(*i)++ // Skip 'const character
		(*i)++ // Skip ' character
		tokenSaver.Add(ConstType, *token)
		*state = DelimiterType
	} else {
		tokenSaver.Add(ErrorType, *token)
		*token = ""
		*state = ErrorType
	}
}

func checkConstBrackets(i int, characters []byte) bool {
	return (string(characters[i-1]) == "'" && string(characters[i+1]) == "'")
}

// readCharacters reads characters from stdin
func readCharacters() []byte {
	reader := bufio.NewReader(os.Stdin)
	var characters []byte
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		characters = append(characters, []byte(input)...)
	}
	return characters
}
