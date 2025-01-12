package optimizer

import (
	"fmt"
	"lab1_2/triad"
	"lab1_2/types"
	"strconv"
)

// countValueIfPossible пытается свернуть триаду, если это возможно
func countValueIfPossible(t triad.Triad, index int, triads *[]triad.Triad) {
	// Шаг 1: Если операнд1 является ссылкой, заменяем её на константу
	if t.Operand1.IsLink() {
		changeLinkOperand(&t.Operand1, *triads)
	}

	// Шаг 2: Если операнд2 является ссылкой, заменяем её на константу
	if t.Operand2.IsLink() {
		changeLinkOperand(&t.Operand2, *triads)
	}

	// Шаг 3: Пытаемся заменить переменную на константу, если такая есть в таблице констант
	tryReplaceVariableWithConstant(&t.Operand1, &t.Operand2, t.Operator, constantsTable)
	tryReplaceVariableWithConstant(&t.Operand2, &t.Operand1, t.Operator, constantsTable)

	// Шаг 4: Если операнд1 — переменная, а операнд2 — число, и оператор присваивания, добавляем значение в таблицу констант
	if t.Operand1.IsVariable() && t.Operand2.IsNumber() && t.Operator == types.Identifier {
		num, err := strconv.Atoi(t.Operand2.GetOperand())
		if err != nil {
			fmt.Printf("Ошибка преобразования в число: %e", err)
		}
		constantsTable[t.Operand1.GetOperand()] = num
	}

	// Шаг 5: Если оба операнда — числа, то выполняем операцию и заменяем триаду на "константную" (C(K, 0))
	if t.Operand1.IsNumber() && t.Operand2.IsNumber() {
		result := performOperation(t) // Выполняем арифметическую операцию
		t = triad.Triad{
			Operator: "C",
			Operand1: triad.OperandFromString(fmt.Sprintf("%d", result)),
			Operand2: triad.OperandFromString("0"),
		}
	}
	(*triads)[index] = t // Обновляем триаду в исходном списке
}

// changeLinkOperand заменяет ссылочный операнд на константу, если это возможно
func changeLinkOperand(o *triad.Operand, triads []triad.Triad) {
	linkIndex := *o.GetLink()          // Получаем индекс ссылки
	linkedTriad := triads[linkIndex-1] // Находим связанную триаду
	if linkedTriad.Operator == "C" && linkedTriad.Operand2.GetOperand() == "0" {
		// Если триада содержит константу, заменяем операнд на её значение
		*o = triad.OperandFromString(linkedTriad.Operand1.GetOperand())
	}
}

// tryReplaceVariableWithConstant пытается заменить переменную на константу из таблицы констант
func tryReplaceVariableWithConstant(operand *triad.Operand, otherOperand *triad.Operand, operator string, constantsTable map[string]int) {
	// Если операнд — переменная, заменяем её на константу, если она есть в таблице
	if operand.IsVariable() && (!otherOperand.IsNumber() || operator != types.Identifier) {
		if value, exists := constantsTable[operand.GetOperand()]; exists {
			*operand = triad.NumberOperand(value) // Замена переменной на её значение
		}
	}
}

// performOperation выполняет арифметическую операцию над двумя числовыми операндами
func performOperation(triad triad.Triad) int {
	var operand1, operand2 int
	// Преобразуем строковые значения операндов в числа
	fmt.Sscanf(triad.Operand1.GetOperand(), "%d", &operand1)
	fmt.Sscanf(triad.Operand2.GetOperand(), "%d", &operand2)

	// Выбираем операцию в зависимости от оператора
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
