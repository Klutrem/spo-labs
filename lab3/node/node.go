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

func (node Node) ContainsNonTerminal() bool {
	var contains bool = false
	for _, child := range node.Children {
		if child.Type == types.NonTerminal {
			contains = true
		}
	}
	return contains
}

func (node Node) IsSimple() bool {
	var isSimple bool = false
	for _, child := range node.Children {
		if len(child.Children) == 0 && child.Type != types.Terminal && child.Lexem != ";" {
			isSimple = true
		} else {
			isSimple = false
		}
	}
	return isSimple
}

func (node Node) IsArithmetic() bool {
	if !node.ContainsNonTerminal() && len(node.Children) == 3 {
		return strings.ContainsAny(node.Children[1].Lexem, types.OperatorChars)
	}
	return false
}

func (node Node) ContainsLink() bool {
	if node.ContainsNonTerminal() {

	}
	return false
}
