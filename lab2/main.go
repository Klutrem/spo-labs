package main

import (
	"encoding/json"
	"fmt"
	lab1 "lab2/modules" // Импорт пакета лексера из первой лабораторной
	"os"
)

// Rule представляет правило продукции в грамматике
type Rule struct {
	Svertka string   // Символ продукции
	Posl    []string // Последовательность продукции
}

// Node представляет узел в дереве разбора
type Node struct {
	Type     string // Нетерминал или терминал
	Lexem    string // Значение лексемы
	Children []*Node
}

func main() {
	// Генерация токенов с использованием лексера из первой лабораторной
	lexems := lab1.RunLexer()

	// Проверка на ошибки лексера
	if hasErrors(lexems) {
		fmt.Println("Лексер обнаружил ошибки.")
		return
	}

	// Определение правил грамматики
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

	// Выполнение синтаксического анализа
	svertka, tree := parse(lexems, rules)

	// Вывод результата
	fmt.Println("Свертка:", svertka)
	if svertka == "S" {
		fmt.Println("Выражение корректно")
		saveTree(tree)
	} else {
		fmt.Println("Выражение некорректно")
	}
}

// hasErrors проверяет, были ли ошибки у лексера
func hasErrors(tokens []lab1.Token) bool {
	for _, token := range tokens {
		if token.Type == lab1.ErrorType {
			return true
		}
	}
	return false
}

// parse выполняет синтаксический анализ на основе заданных токенов и правил грамматики
func parse(tokens []lab1.Token, rules []Rule) (string, *Node) {
	var stack []*Node

	// Помещение всех токенов в стек
	for _, token := range tokens {
		stack = append(stack, &Node{Type: token.Type, Lexem: token.Value})
	}

	// Применение правил, пока стек не содержит только один элемент
	for len(stack) > 1 {
		applied := false

		for _, rule := range rules {
			if canApplyRule(stack, rule) {
				fmt.Printf("Применение правила: %s -> %v\n", rule.Svertka, rule.Posl)
				stack = applyRule(stack, rule)
				applied = true
				break
			} else {
				fmt.Printf("Невозможно применить правило: %s -> %v к стеку\n", rule.Svertka, rule.Posl)
			}
		}

		if !applied {
			fmt.Println("Ошибка: невозможно применить ни одно правило.")
			fmt.Println("Текущий стек:", stack)
			return "ошибка", nil
		}
	}

	return stack[0].Type, stack[0]
}

// canApplyRule проверяет, может ли правило быть применено к текущему стеку
func canApplyRule(stack []*Node, rule Rule) bool {
	if len(rule.Posl) > len(stack) {
		return false
	}

	for i, symbol := range rule.Posl {
		stackIndex := len(stack) - len(rule.Posl) + i
		if symbol != stack[stackIndex].Type {
			fmt.Printf("Несоответствие правила: символ правила %s, символ стека %s в индексе %d стека\n", symbol, stack[stackIndex].Type, stackIndex)
			return false
		}
	}

	return true
}

// applyRule применяет правило грамматики к текущему стеку
func applyRule(stack []*Node, rule Rule) []*Node {
	newNode := &Node{Type: rule.Svertka, Lexem: "", Children: []*Node{}}
	stackSize := len(stack)

	// Добавление дочерних элементов к новому узлу в правильном порядке
	for i := 0; i < len(rule.Posl); i++ {
		newNode.Children = append(newNode.Children, stack[stackSize-len(rule.Posl)+i])
	}

	// Удаление соответствующих элементов из стека
	stack = stack[:stackSize-len(rule.Posl)]

	// Добавление нового узла в стек
	stack = append(stack, newNode)

	fmt.Printf("Новый размер стека: %d\n", len(stack))
	printStack(stack)

	return stack
}

// printStack выводит текущий стек для отладки
func printStack(stack []*Node) {
	fmt.Print("Текущий стек: [ ")
	for _, node := range stack {
		fmt.Printf("{%s, %s} ", node.Type, node.Lexem)
	}
	fmt.Println("]")
}

// saveTree сохраняет дерево разбора в JSON файл
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
