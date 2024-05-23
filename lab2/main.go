package main

import (
	"encoding/json"
	"fmt"
	lab1 "lab2/modules" // Importing the lexer package from the first lab
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
		{Svertka: "T", Posl: []string{"E"}},
		{Svertka: "E", Posl: []string{"(", "F", ")"}},
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

	// Push all tokens to the stack
	for _, token := range tokens {
		stack = append(stack, &Node{Type: token.Type, Lexem: token.Value})
	}

	// Apply rules until the stack contains only one element
	for len(stack) > 1 {
		applied := false

		for _, rule := range rules {
			if canApplyRule(stack, rule) {
				fmt.Printf("Applying rule: %s -> %v\n", rule.Svertka, rule.Posl)
				stack = applyRule(stack, rule)
				applied = true
				break
			}
		}

		if !applied {
			fmt.Println("Error: could not apply any rule.")
			return "error", nil
		}
	}

	return stack[0].Type, stack[0]
}

// canApplyRule checks if a rule can be applied to the current stack
func canApplyRule(stack []*Node, rule Rule) bool {
	if len(rule.Posl) > len(stack) {
		return false
	}

	for i, symbol := range rule.Posl {
		if symbol != stack[len(stack)-len(rule.Posl)+i].Type {
			return false
		}
	}

	return true
}

// applyRule applies a grammar rule to the current stack
func applyRule(stack []*Node, rule Rule) []*Node {
	newNode := &Node{Type: rule.Svertka, Lexem: "", Children: []*Node{}}
	stackSize := len(stack)

	// Add the children to the new node in reverse order (to maintain the correct sequence)
	for i := len(rule.Posl) - 1; i >= 0; i-- {
		newNode.Children = append([]*Node{stack[stackSize-len(rule.Posl)+i]}, newNode.Children...)
	}

	// Remove the matched elements from the stack
	stack = stack[:stackSize-len(rule.Posl)]

	// Add the new node to the stack
	stack = append(stack, newNode)

	fmt.Printf("New stack size: %d\n", len(stack))
	printStack(stack)

	return stack
}

// printStack prints the current stack for debugging purposes
func printStack(stack []*Node) {
	fmt.Print("Current stack: [ ")
	for _, node := range stack {
		fmt.Printf("{%s, %s} ", node.Type, node.Lexem)
	}
	fmt.Println("]")
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
