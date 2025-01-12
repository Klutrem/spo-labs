package triad

import (
	"lab1_2/node"
	"lab1_2/types"
	"strings"
)

// ConvertNodeToTriads конвертирует узел (node) в триады и возвращает последний операнд.
// Аргументы:
// - nodeToConvert: узел для преобразования
// - triads: указатель на список, куда будут добавлены новые триады
func ConvertNodeToTriads(nodeToConvert node.Node, triads *[]Triad) Operand {
	var lastOperator Operand
	switch {
	case nodeToConvert.HasDelimeter():
		// Если узел содержит разделитель, конвертируем его дочерний узел
		ConvertNodeToTriads(nodeToConvert.Children[0], triads)
	case nodeToConvert.IsParenthesesExpression():
		// Если узел является выражением в скобках, обрабатываем внутреннее выражение
		innerNode := nodeToConvert.Children[1] // Внутреннее выражение в скобках
		lastOperator = ConvertNodeToTriads(innerNode, triads)
	case nodeToConvert.IsSimpleArithmetic():
		// Если узел является простым арифметическим выражением, конвертируем его в триаду
		ConvertArithmeticNodeToTriad(nodeToConvert, triads)
	case nodeToConvert.IsComplexExpression() || nodeToConvert.IsAssignment():
		// Если узел является сложным выражением или присваиванием, конвертируем его в триаду
		ConvertExpressionToTriad(nodeToConvert, triads)
	default:
		// Для всех остальных случаев рекурсивно обрабатываем дочерние узлы
		for _, child := range nodeToConvert.Children {
			lastOperator = ConvertNodeToTriads(child, triads)
		}
	}
	// Если последний оператор пустой, создаем ссылку на текущую длину списка триад
	if lastOperator.GetOperand() == "" {
		lastOperator = LinkOperand(len(*triads))
	}

	return lastOperator
}

// ConvertArithmeticNodeToTriad преобразует узел с арифметическим выражением в триаду
// Аргументы:
// - nodeToConvert: узел для преобразования
// - triads: указатель на список триад
// Возвращает:
// - Новую триаду
func ConvertArithmeticNodeToTriad(nodeToConvert node.Node, triads *[]Triad) Triad {
	operator := nodeToConvert.Children[1].Lexem                   // Оператор (например, "+", "-")
	operand := OperandFromSimpleNode(nodeToConvert.Children[0])   // Первый операнд
	operator2 := OperandFromSimpleNode(nodeToConvert.Children[2]) // Второй операнд

	// Создаем новую триаду
	newTriad := Triad{
		Operator: operator,
		Operand1: operand,
		Operand2: operator2,
	}

	// Добавляем триаду в список
	*triads = append(*triads, newTriad)
	return newTriad
}

// ConvertExpressionToTriad преобразует сложное выражение или присваивание в триаду
// Аргументы:
// - nodeToConvert: узел для преобразования
// - triads: указатель на список триад
// Возвращает:
// - Новую триаду
func ConvertExpressionToTriad(nodeToConvert node.Node, triads *[]Triad) Triad {
	var operand1, operand2 Operand
	var operator string

	// Конвертируем первый и второй операнды
	operand1 = ConvertToOperand(nodeToConvert.Children[0], triads)
	operand2 = ConvertToOperand(nodeToConvert.Children[2], triads)

	// Определяем оператор
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
	if n.IsSimple() && n.IsTerminal() {
		o = OperandFromSimpleNode(n)
	} else {
		o = ConvertNodeToTriads(n, triads)
	}
	return o
}
