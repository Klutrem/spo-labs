package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	lab1 "lab2/modules" // Importing the lexer package from the first lab
)

// Rule represents a production rule in the grammar
type Rule struct {
	Svertka string   // The production symbol
	Posl    []string // The production sequence
}

func main() {
	// Generate tokens using the lexer from the first lab
	lexems := lab1.RunLexer()

	// Convert tokens to the required format
	var ArrayLexems []lab1.Token
	for _, it := range lexems {
		if it.Type == "number" {
			ArrayLexems = append(ArrayLexems, lab1.Token{Type: "identifier", Value: it.Value})
		} else {
			ArrayLexems = append(ArrayLexems, it)
		}
	}

	// Check for lexer errors
	if hasErrors(ArrayLexems) {
		fmt.Println("Lexer encountered errors.")
		return
	}

	// Define grammar rules
	rules := []Rule{
		{Svertka: "E", Posl: []string{"identifier"}},
		{Svertka: "E", Posl: []string{"arithmetic", "l_bracket", "F", "r_bracket"}},
		{Svertka: "E", Posl: []string{"arithmetic", "l_bracket", "identifier", "r_bracket"}},
		{Svertka: "E", Posl: []string{"l_bracket", "identifier", "r_bracket"}},
		{Svertka: "E", Posl: []string{"l_bracket", "F", "r_bracket"}},
		{Svertka: "T", Posl: []string{"T", "arithmetic", "E"}},
		{Svertka: "T", Posl: []string{"E"}},
		{Svertka: "F", Posl: []string{"identifier", "arithmetic", "identifier"}},
		{Svertka: "F", Posl: []string{"F", "arithmetic", "T"}},
		{Svertka: "F", Posl: []string{"T"}},
		{Svertka: "S", Posl: []string{"F", "operator", "identifier", "delimiter"}},
		{Svertka: "S", Posl: []string{"F", "operator", "F", "delimiter"}},
	}

	// Perform parsing
	svertka, tree := parse(ArrayLexems, rules)

	// Print the result
	fmt.Println("Свертка:", svertka)
	if svertka == "S" {
		fmt.Println("Выражение корректно")
		saveTree(tree)
	} else {
		fmt.Println("Выражение некорректно")
	}
}

// hasErrors checks if lexer encountered errors
func hasErrors(tokens []lab1.Token) bool {
	for _, token := range tokens {
		if token.Type == "error" {
			return true
		}
	}
	return false
}

// parse function performs parsing based on the given tokens and grammar rules
func parse(tokens []lab1.Token, rules []Rule) (string, *Node) {
	return "", nil
}

// Node represents a node in the parse tree
type Node struct {
	Type     string // NonTerminal or Terminal
	Lexem    string // Lexeme value
	Children []*Node
}

// saveTree function saves the parse tree to a JSON file
func saveTree(tree *Node) {
	jsonData, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = ioutil.WriteFile("Tree.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Дерево сохранено в файл Tree.json")
}
