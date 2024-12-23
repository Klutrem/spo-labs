package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	terminal    = "terminal"
	nonTerminal = "NonTerminal"
)

type Node struct {
	Type     string `json:"type"`     // Тип узла (Terminal или NonTerminal)
	Lexem    string `json:"lexem"`    // Лексема (строка, представляющая символ)
	Children []Node `json:"children"` // Дочерние элементы (если есть)
}

func main() {
	var nodes []Node
	// читаем файл, сгенерированный алгоритмом из прошлых ЛР
	fileData, err := os.ReadFile("../output.json")
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}

	if err := json.Unmarshal(fileData, &nodes); err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}
	var result string
	println("Начальные выражения:")
	for _, node := range nodes {
		// Вывод начальной лексеммы
		fmt.Printf("%+v\n", node.Lexem)
		res, err := parseNode(node)
		if err != nil {
			fmt.Printf("Ошибка парсинга ноды: %s\n", err.Error())
		}
		result += "\n" + res
	}
	fmt.Printf("Результат парсинга: %s\n", result)
}

// parseNode: парсит узел в зависимости от его типа
func parseNode(node Node) (string, error) {
	// Если это простая нода, обрабатываем её
	if isSimpleNode(node) {
		return convertSimpleNodeToAsm(node)
	}

	// иначе ищем простую ноду
	return parseNonTerminalNode(node)
}

// parseNonTerminalNode: ищет терминал среди дочерних элементов и обрабатывает его
func parseNonTerminalNode(node Node) (string, error) {
	// Если у ноды есть дети, проверим их
	for _, child := range node.Children {
		// Если дочерний элемент - терминал, обрабатываем его
		if isSimpleNode(child) {
			return convertSimpleNodeToAsm(child)
		}
		// Если дочерний элемент - нетерминал, рекурсивно обрабатываем его
		if child.Type == nonTerminal {
			return parseNonTerminalNode(child)
		}
	}
	// Если мы не нашли терминал среди детей, ошибка
	return "", fmt.Errorf("Не удалось найти терминал среди дочерних элементов: %+v", node.Lexem)
}

func isSimpleNode(node Node) bool {
	var isSimple bool = false
	for _, child := range node.Children {
		if len(child.Children) == 0 && child.Type != terminal {
			isSimple = true
		} else {
			isSimple = false
		}
	}
	return isSimple
}

// convertSimpleNodeToAsm: преобразует лексему в ассемблерный код
func convertSimpleNodeToAsm(node Node) (string, error) {
	// Присваивание (например, a := 3)
	if strings.ContainsAny(node.Lexem, identifier) && len(node.Children) >= 3 {
		left := node.Children[0].Lexem
		right := node.Children[2].Lexem

		// Генерируем MOV команду для присваивания
		return generateMov(left, right), nil
	}

	// Дополнительные проверки для операторов
	if strings.ContainsAny(node.Lexem, OperatorChars) && len(node.Children) > 0 {
		// Например, если есть оператор, то мы можем добавить логику для работы с операциями
		return processOperator(node), nil
	}

	// Проверка для других лексем (например, чисел или идентификаторов)
	if strings.ContainsAny(node.Lexem, Alphanumeric) {
		// Здесь можно добавить логику для других типов выражений
		return "", fmt.Errorf("Не тот тип")
	}

	// Если лексема не распознана, вернуть ошибку
	return "", fmt.Errorf("Неизвестная лексема для преобразования в ассемблер: %v", node.Lexem)
}

// generateMov: генерирует команду MOV для присваивания
func generateMov(left, right string) string {
	return fmt.Sprintf("MOV %s, %s", left, right)
}

// processOperator: обработка операторов (например, +, -, *, /)
func processOperator(node Node) string {
	// Допустим, просто обработаем простые операторы
	// Например: a + b
	left := node.Children[0].Lexem
	operator := node.Children[1].Lexem
	right := node.Children[1].Lexem

	// Генерируем команду, например: ADD a, b
	return fmt.Sprintf("%s %s, %s", operator, left, right)
}
