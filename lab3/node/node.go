package node

import (
	"lab1_2/types"
	"strings"
)

type Node struct {
	Type     string `json:"type"`     // Тип узла (Terminal или NonTerminal)
	Lexem    string `json:"lexem"`    // Лексема (строка, представляющая символ)
	Children []Node `json:"children"` // Дочерние элементы (если есть)
}

// ContainsNonTerminal проверяет, содержит ли нода хотя бы один дочерний нода типа NonTerminal
func (node Node) ContainsNonTerminal() bool {
	for _, child := range node.Children {
		if !child.IsTerminal() {
			return true
		}
	}
	return false
}

func (node Node) IsTerminal() bool {
	return node.Type == types.Terminal
}

// IsSimple проверяет, является ли нода "простым" (терминалом или узлом с простыми лексемами)
func (node Node) IsSimple() bool {
	if len(node.Children) == 0 {
		return true
	}
	for _, child := range node.Children {
		if child.ContainsNonTerminal() || !child.IsTerminal() {
			return false
		}
	}
	return true
}

// IsSimpleArithmetic проверяет, является ли нода арифметическим выражением (например, "a + b")
func (node Node) IsSimpleArithmetic() bool {
	if len(node.Children) == 3 && !node.ContainsNonTerminal() {
		// Для арифметического выражения должно быть 3 элемента: переменная, оператор и переменная
		return strings.ContainsAny(node.Children[1].Lexem, types.OperatorChars)
	}
	return false
}

func (node Node) HasDelimeter() bool {
	for _, child := range node.Children {
		if child.HasDelimeter() || child.Lexem == types.Delimiter {
			return true
		}
	}
	return false
}

// IsAssignment проверяет, является ли нода присваиванием (например, "a := 3")
func (node Node) IsAssignment() bool {
	if len(node.Children) == 3 {
		return node.Children[1].Lexem == types.Identifier
	}
	return false
}

// IsParenthesesExpression проверяет, является ли выражение в скобках (например, "(a + b)")
func (node Node) IsParenthesesExpression() bool {
	if len(node.Children) >= 3 {
		return node.Children[0].Lexem == "(" || node.Children[2].Lexem == ")"
	}
	return false
}

// IsComplexExpression проверяет, является ли нода сложным выражением с несколькими операторами
func (node Node) IsComplexExpression() bool {
	if node.Type == types.NonTerminal {
		return true
	}
	// Проверка на наличие арифметических выражений в дочерних узлах
	return node.ContainsNonTerminal()
}

// ContainsLink проверяет, ссылается ли нода на другие узлы
func (node Node) ContainsLink() bool {
	for _, child := range node.Children {
		if child.Type == types.NonTerminal {
			return true
		}
	}
	return false
}

// PrintLexems выводит лексемы всех узлов в дереве
func (node Node) PrintLexems() {
	if node.Type == types.Terminal {
		print(node.Lexem + " ")
	}
	for _, child := range node.Children {
		child.PrintLexems()
	}
}
