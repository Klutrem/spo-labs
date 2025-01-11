package triad

import (
	"lab1_2/node"
	"lab1_2/types"
	"strings"
)

// Основная функция для конвертации узлов в триады
func ConvertNodeToTriads(nodeToConvert node.Node, triads *[]Triad) Operand {
	var lastOperator Operand
	switch {
	case nodeToConvert.HasDelimeter():
		ConvertNodeToTriads(nodeToConvert.Children[0], triads)
	case nodeToConvert.IsParenthesesExpression():
		// Обработка выражения в скобках
		innerNode := nodeToConvert.Children[1] // Внутреннее выражение в скобках
		lastOperator = ConvertNodeToTriads(innerNode, triads)
	case nodeToConvert.IsSimpleArithmetic():
		ConvertArithmeticNodeToTriad(nodeToConvert, triads)
	case nodeToConvert.IsComplexExpression() || nodeToConvert.IsAssignment():
		ConvertExpressionToTriad(nodeToConvert, triads)
	default:
		for _, child := range nodeToConvert.Children {
			lastOperator = ConvertNodeToTriads(child, triads)
		}
	}
	if lastOperator.GetOperand() == "" {
		lastOperator = LinkOperand(len(*triads))
	}

	return lastOperator
}

// Преобразование арифметического выражения в триаду
func ConvertArithmeticNodeToTriad(nodeToConvert node.Node, triads *[]Triad) Triad {
	operator := nodeToConvert.Children[1].Lexem
	operand := OperandFromSimpleNode(nodeToConvert.Children[0])
	operator2 := OperandFromSimpleNode(nodeToConvert.Children[2])

	// Создание триады и добавление в список
	newTriad := Triad{
		Operator: operator,
		Operand1: operand,
		Operand2: operator2,
	}

	*triads = append(*triads, newTriad)
	return newTriad
}

// Преобразование выражения в триаду
func ConvertExpressionToTriad(nodeToConvert node.Node, triads *[]Triad) Triad {
	var operand1, operand2 Operand
	var operator string

	operand1 = ConvertToOperand(nodeToConvert.Children[0], triads)
	operand2 = ConvertToOperand(nodeToConvert.Children[2], triads)

	if strings.ContainsAny(nodeToConvert.Children[1].Lexem, types.OperatorChars) {
		operator = nodeToConvert.Children[1].Lexem
	} else if nodeToConvert.Children[1].Lexem == types.Identifier {
		operator = types.Identifier
	}

	// Создаем и добавляем новую триаду
	newTriad := Triad{
		Operator: operator,
		Operand1: operand1,
		Operand2: operand2,
	}

	*triads = append(*triads, newTriad)
	return newTriad
}

func ConvertToOperand(n node.Node, triads *[]Triad) Operand {
	var o Operand
	if n.IsSimple() {
		o = OperandFromSimpleNode(n)
	} else {
		o = ConvertNodeToTriads(n, triads)
	}
	return o
}
