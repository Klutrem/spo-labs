package optimizer

import (
	"fmt"
	"lab1_2/triad"
	"lab1_2/types"
	"strconv"
)

// Мапа для хранения значений переменных, которые стали известны
var constantsTable map[string]int = make(map[string]int)

// OptimizeTriads выполняет свертку триад
func OptimizeTriads(triads *[]triad.Triad) {
	// Проход по всем триадам
	for index, triad := range *triads {
		// Пробуем свертку для текущей триады
		countValueIfPossible(triad, index, triads)
	}
	removeRedundantTriadsWithConstants(triads)
}

// countValueIfPossible пытается свернуть триаду, если это возможно
func countValueIfPossible(t triad.Triad, index int, triads *[]triad.Triad) {
	// Шаг 1: Заменить операнд1 на константу, если он есть в triadResults
	if t.Operand1.IsLink() {
		linkIndex := *t.Operand1.GetLink()
		linkedTriad := (*triads)[linkIndex-1]
		if linkedTriad.Operator == "C" && linkedTriad.Operand2.GetOperand() == "0" {
			// Заменяем операнд1 на значение константы из триады
			t.Operand1 = triad.OperandFromString(linkedTriad.Operand1.GetOperand())
		}
	}

	// Шаг 2: Заменить операнд2 на константу, если он есть в triadResults
	if t.Operand2.IsLink() {
		linkIndex := *t.Operand2.GetLink()
		linkedTriad := (*triads)[linkIndex-1]
		if linkedTriad.Operator == "C" && linkedTriad.Operand2.GetOperand() == "0" {
			// Заменяем операнд2 на значение константы из триады
			t.Operand2 = triad.OperandFromString(linkedTriad.Operand1.GetOperand())
		}
	}

	if t.Operand1.IsVariable() && (!t.Operand2.IsNumber() || t.Operator != types.Identifier) {
		if value, exists := constantsTable[t.Operand1.GetOperand()]; exists {
			t.Operand1 = triad.NumberOperand(value) // Заменяем переменную на её значение
		}
	}

	// Если операнд2 - переменная и она есть в таблице, заменяем её на значение из таблицы
	if t.Operand2.IsVariable() && (!t.Operand1.IsNumber() || t.Operator != types.Identifier) {
		if value, exists := constantsTable[t.Operand2.GetOperand()]; exists {
			t.Operand2 = triad.NumberOperand(value) // Заменяем переменную на её значение
		}
	}

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
