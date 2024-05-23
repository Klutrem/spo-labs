package main

import (
	"encoding/json"
	"fmt"
	lab1 "lab2/modules"
	"os"
)

// Rule represents a production rule in the grammar
type Rule struct {
	Svertka string   // The production symbol
	Posl    []string // The production sequence
}

// Node represents a node in the parse tree
type Node struct {
	Type     string // NonTerminal or Terminal
	Lexem    string // Lexeme value
	Children []*Node
}

func main() {
	// Generate tokens using the lexer from the first lab
	lexems := lab1.RunLexer()

	// Check for lexer errors
	if hasErrors(lexems) {
		fmt.Println("Lexer encountered errors.")
		return
	}

	// Define grammar rules
	rules := []Rule{
		{Svertka: "S", Posl: []string{"identifier", ":=", "F"}},
		{Svertka: "F", Posl: []string{"F", "+", "T"}},
		{Svertka: "F", Posl: []string{"T"}},
		{Svertka: "T", Posl: []string{"T", "*", "E"}},
		{Svertka: "T", Posl: []string{"T", "/", "E"}},
		{Svertka: "T", Posl: []string{"E"}},
		{Svertka: "E", Posl: []string{"(", "F", ")"}},
		{Svertka: "E", Posl: []string{"-", "(", "F", ")"}},
		{Svertka: "E", Posl: []string{"identifier"}},
	}

	// Perform parsing
	svertka, tree := parse(lexems, rules)

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
		if token.Type == lab1.ErrorType {
			return true
		}
	}
	return false
}

// parse function performs parsing based on the given tokens and grammar rules
func parse(tokens []lab1.Token, rules []Rule) (string, *Node) {
	var stack []*Node
	index := 0

	for index < len(tokens) {
		token := tokens[index]

		// Try to apply rules
		applied := false
		for _, rule := range rules {
			if canApplyRule(token, rule) {
				stack = applyRule(rule, stack, tokens, &index)
				applied = true
				break
			}
		}

		if !applied {
			fmt.Println("Error: could not apply any rule.")
			return "error", nil
		}
	}

	if len(stack) != 1 {
		fmt.Println("Error: invalid parse tree.")
		return "error", nil
	}

	return stack[0].Type, stack[0]
}

// canApplyRule checks if a rule can be applied to the current token and tokens sequence
func canApplyRule(token lab1.Token, rule Rule) bool {
	return rule.Posl[0] == token.Type
}

// applyRule applies a grammar rule to the current token and tokens sequence
func applyRule(rule Rule, stack []*Node, tokens []lab1.Token, index *int) []*Node {
	newNode := &Node{Type: rule.Svertka, Lexem: "", Children: []*Node{}}

	for _, symbol := range rule.Posl {
		if symbol == tokens[*index].Type {
			newNode.Children = append(newNode.Children, &Node{Type: tokens[*index].Type, Lexem: tokens[*index].Value})
			*index++
		} else {
			newNode.Children = append(newNode.Children, stack[len(stack)-1])
			stack = stack[:len(stack)-1]
		}
	}

	stack = append(stack, newNode)
	return stack
}

// saveTree function saves the parse tree to a JSON file
func saveTree(tree *Node) {
	jsonData, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = os.WriteFile("Tree.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Дерево сохранено в файл Tree.json")
}
