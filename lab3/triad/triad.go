package triad

import "fmt"

// Triad представляет триаду, которая состоит из:
// - Оператора (например, "+", "-", "*" и т.д.)
// - Двух операндов (Operand1 и Operand2), которые являются входными значениями для операции
type Triad struct {
	Operator string  // Оператор, выполняющий операцию
	Operand1 Operand // Первый операнд (входное значение 1)
	Operand2 Operand // Второй операнд (входное значение 2)
}

// ToString возвращает строковое представление триады в формате:
// "Оператор(Операнд1, Операнд2)"
func (t Triad) ToString() string {
	return fmt.Sprintf("%s(%s, %s)", t.Operator, t.Operand1.GetOperand(), t.Operand2.GetOperand())
}

// Equals проверяет равенство двух триад. Две триады считаются равными, если:
// - Их операторы совпадают
// - Оба операнда совпадают по значениям
// Не обойтись простым сравнением структур, поскольку в операндах используются ссылки
// и при сравнении они могут ссылаться на разные адреса в памяти
func (t Triad) Equals(tr Triad) bool {
	return t.Operand1.GetOperand() == tr.Operand1.GetOperand() &&
		t.Operand2.GetOperand() == tr.Operand2.GetOperand() &&
		t.Operator == tr.Operator
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
