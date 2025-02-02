package code_generation

import (
	"fmt"
	"lab1_2/triad"
	"strings"
)

var updatedLinkOperands map[int]string = make(map[int]string)

// GenerateAssemblyCode генерирует ассемблерный код на основе списка триад
func GenerateAssemblyCode(triads []triad.Triad) string {
	var assemblyCode strings.Builder

	// Переменная для отслеживания временных регистров (например, AX, BX, CX, DX)
	regIndex := 0

	// Обрабатываем каждую триаду
	for i, triad := range triads {
		operand1 := operandToString(triad.Operand1, i, triads)
		operand2 := operandToString(triad.Operand2, i, triads)
		switch triad.Operator {
		case "*":
			// Умножение: *(B,C)
			assemblyCode.WriteString(fmt.Sprintf("mul %s, %s\n", operand1, operand2))
			regIndex++ // Увеличиваем индекс для следующей операции

		case "+":
			// Сложение: +(операнд1, операнд2)
			assemblyCode.WriteString(fmt.Sprintf("add %s,%s\n", operand1, operand2))
			regIndex++ // Увеличиваем индекс для следующей операции

		case "-":
			// Вычитание: -(операнд1, операнд2)
			assemblyCode.WriteString(fmt.Sprintf("sub %s,%s\n", operand1, operand2))
			regIndex++ // Увеличиваем индекс для следующей операции

		case "/":
			// деление: /(операнд1, операнд2)
			assemblyCode.WriteString(fmt.Sprintf("div %s,%s\n", operand1, operand2))
			regIndex++ // Увеличиваем индекс для следующей операции

		case ":=":
			// Присваивание: :=(A, операнд)
			assemblyCode.WriteString(fmt.Sprintf("mov %s, %s\n", operand1, operand2))

			// Сохраняем результат присваивания
			updatedLinkOperands[i] = operand1
		default:
			assemblyCode.WriteString(fmt.Sprintf("Unknown operator: %s\n", triad.Operator))
		}
	}

	return assemblyCode.String()
}

func operandToString(o triad.Operand, index int, triads []triad.Triad) string {
	if o.IsLink() {
		// Проверяем, есть ли результат в `updatedLinkOperands`
		if linkedOperand, ok := updatedLinkOperands[*o.GetLink()]; ok {
			return linkedOperand
		}

		// Если результат не найден, возвращаем операнд из ссылки
		return triads[*o.GetLink()-1].Operand1.GetOperand()
	}

	return o.GetOperand()
}
