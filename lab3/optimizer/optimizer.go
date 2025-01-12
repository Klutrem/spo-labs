package optimizer

import (
	"fmt"
	"lab1_2/triad"
	"lab1_2/types"
	"strconv"
)

// Мапа для хранения значений переменных, которые стали известны
var triadResults map[int]int = make(map[int]int)
var constantsTable map[string]int = make(map[string]int)

// OptimizeTriads выполняет свертку триад
func OptimizeTriads(triads *[]triad.Triad) {
	// Проход по всем триадам
	for index, triad := range *triads {
		// Пробуем свертку для текущей триады
		countValueIfPossible(triad, index, triads)
	}
}

// countValueIfPossible пытается свернуть триаду, если это возможно
func countValueIfPossible(t triad.Triad, index int, triads *[]triad.Triad) {
	// Шаг 1: Заменить операнд1 на константу, если он есть в triadResults
	if t.Operand1.IsLink() {
		triadValue, found := triadResults[*t.Operand1.GetLink()]
		if found {
			// Заменяем операнд1 на константу
			t.Operand1 = triad.OperandFromString(fmt.Sprintf("%d", triadValue))
		}
	}

	// Шаг 2: Заменить операнд2 на константу, если он есть в triadResults
	if t.Operand2.IsLink() {
		triadValue, found := triadResults[*t.Operand2.GetLink()]
		if found {
			// Заменяем операнд2 на константу
			t.Operand2 = triad.OperandFromString(fmt.Sprintf("%d", triadValue))
		}
	}

	if t.Operand1.IsVariable() {
		if value, exists := constantsTable[t.Operand1.GetOperand()]; exists {
			t.Operand1 = triad.NumberOperand(value) // Заменяем переменную на её значение
		}
	}

	// Если операнд2 - переменная и она есть в таблице, заменяем её на значение из таблицы
	if t.Operand2.IsVariable() {
		if value, exists := constantsTable[t.Operand2.GetOperand()]; exists {
			t.Operand2 = triad.NumberOperand(value) // Заменяем переменную на её значение
		}
	}

	if t.Operand1.IsVariable() && t.Operand2.IsNumber() && t.Operator == types.Identifier {
		num, err := strconv.Atoi(t.Operand2.GetOperand())
		if err != nil {
			fmt.Printf("Ошибка перевода в число: %e", err)
		}
		triadResults[index+1] = num
		constantsTable[t.Operand1.GetOperand()] = num
	}

	// Шаг 3: Если все операнды являются константами, то выполняем операцию
	if t.Operand1.IsNumber() && t.Operand2.IsNumber() {
		// Выполняем операцию на операндах
		result := performOperation(t)
		// Сохраняем результат в triadResults
		triadResults[index+1] = result
		// Обновляем текущую триаду на особую C(K, 0)
		t = triad.Triad{
			Operator: "C",
			Operand1: triad.OperandFromString(fmt.Sprintf("%d", result)),
			Operand2: triad.OperandFromString("0"),
		}
	}
	(*triads)[index] = t
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
