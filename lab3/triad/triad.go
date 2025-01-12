package triad

import "fmt"

type Triad struct {
	Operator string  // Операнд (результат операции)
	Operand1 Operand // Первый операнд (операнды для операции)
	Operand2 Operand // Второй операнд (операнды для операции)
}

func (t Triad) ToString() string {
	return fmt.Sprintf("%s(%s, %s)", t.Operator, t.Operand1.GetOperand(), t.Operand2.GetOperand())
}

func (t Triad) Equals(tr Triad) bool {
	return t.Operand1.GetOperand() == tr.Operand1.GetOperand() && t.Operand2.GetOperand() == tr.Operand2.GetOperand() && t.Operator == tr.Operator
}

// MergeTriadList объединяет несколько списков триад и корректирует ссылки
func MergeTriadList(triadsList ...[]Triad) []Triad {
	var outputTriads []Triad
	offset := 0 // Смещение для обновления ссылок

	for _, triads := range triadsList {
		for _, triad := range triads {
			// Копируем триаду, чтобы избежать изменения исходного списка
			newTriad := triad

			// Обновляем ссылки, если они есть
			if newTriad.Operand1.IsLink() {
				newIndex := *newTriad.Operand1.linkTo + offset
				newTriad.Operand1 = LinkOperand(newIndex)
			}
			if newTriad.Operand2.IsLink() {
				newIndex := *newTriad.Operand2.linkTo + offset
				newTriad.Operand2 = LinkOperand(newIndex)
			}

			// Добавляем обновленную триаду в результирующий список
			outputTriads = append(outputTriads, newTriad)
		}
		// Обновляем смещение на длину текущего списка триад
		offset += len(triads)
	}

	return outputTriads
}
