package optimizer

import (
	"fmt"
	"lab1_2/triad"
	"lab1_2/types"
	"strconv"
)

// countValueIfPossible пытается свернуть триаду, если это возможно
func countValueIfPossible(t triad.Triad, index int, triads *[]triad.Triad) {
	// Шаг 1: Заменить операнд1 на константу, если он есть в triadResults
	if t.Operand1.IsLink() {
		changeLinkOperand(&t.Operand1, *triads)
	}

	// Шаг 2: Заменить операнд2 на константу, если он есть в triadResults
	if t.Operand2.IsLink() {
		changeLinkOperand(&t.Operand2, *triads)
	}

	tryReplaceVariableWithConstant(&t.Operand1, &t.Operand2, t.Operator, constantsTable)
	tryReplaceVariableWithConstant(&t.Operand2, &t.Operand1, t.Operator, constantsTable)

	if t.Operand1.IsVariable() && t.Operand2.IsNumber() && t.Operator == types.Identifier {
		num, err := strconv.Atoi(t.Operand2.GetOperand())
		if err != nil {
			fmt.Printf("Ошибка перевода в число: %e", err)
		}
		constantsTable[t.Operand1.GetOperand()] = num
	}

	// Шаг 3: Если все операнды являются константами, то выполняем операцию
	if t.Operand1.IsNumber() && t.Operand2.IsNumber() {
		// Выполняем операцию на операндах
		result := performOperation(t)
		// Обновляем текущую триаду на особую C(K, 0)
		t = triad.Triad{
			Operator: "C",
			Operand1: triad.OperandFromString(fmt.Sprintf("%d", result)),
			Operand2: triad.OperandFromString("0"),
		}
	}
	(*triads)[index] = t
}

func changeLinkOperand(o *triad.Operand, triads []triad.Triad) {
	linkIndex := *o.GetLink()
	linkedTriad := triads[linkIndex-1]
	if linkedTriad.Operator == "C" && linkedTriad.Operand2.GetOperand() == "0" {
		// Заменяем операнд2 на значение константы из триады
		*o = triad.OperandFromString(linkedTriad.Operand1.GetOperand())
	}
}

func tryReplaceVariableWithConstant(operand *triad.Operand, otherOperand *triad.Operand, operator string, constantsTable map[string]int) {
	if operand.IsVariable() && (!otherOperand.IsNumber() || operator != types.Identifier) {
		if value, exists := constantsTable[operand.GetOperand()]; exists {
			*operand = triad.NumberOperand(value) // Заменяем переменную на её значение
		}
	}
}

// performOperation выполняет операцию над двумя операндами
func performOperation(triad triad.Triad) int {
	var operand1, operand2 int
	fmt.Sscanf(triad.Operand1.GetOperand(), "%d", &operand1)
	fmt.Sscanf(triad.Operand2.GetOperand(), "%d", &operand2)

	// В зависимости от оператора выполняем соответствующую операцию
	switch triad.Operator {
	case "+":
		return operand1 + operand2
	case "*":
		return operand1 * operand2
	case "-":
		return operand1 - operand2
	case "/":
		return operand1 / operand2
	}
	return 0
}

func removeRedundantTriadsWithConstants(triads *[]triad.Triad) {
	var optimizedTriads []triad.Triad
	for _, t := range *triads {
		// Проверяем, является ли триада C(K, 0)
		if t.Operator == "C" && t.Operand2.GetOperand() == "0" {
			// Пропускаем эту триаду (не добавляем в результат)
			continue
		}
		// Добавляем триаду в оптимизированный список
		optimizedTriads = append(optimizedTriads, t)
	}
	// Обновляем исходный список триад
	*triads = optimizedTriads
}
