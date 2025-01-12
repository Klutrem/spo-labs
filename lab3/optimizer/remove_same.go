package optimizer

import (
	"lab1_2/triad"
)

// Основной алгоритм исключения лишних операций
func eliminateRedundantOperations(triads []triad.Triad) []triad.Triad {
	dep := make(map[string]int)           // хранит зависимости переменных
	depTriads := make([]int, len(triads)) // хранит зависимости для каждой триады
	result := []triad.Triad{}             // итоговый результат

	for i, t := range triads {
		// Шаг 1: Замена операндов, ссылающихся на SAME
		if t.Operand1.IsLink() {
			linkedTriad := triads[*t.Operand1.GetLink()-1]
			if linkedTriad.Operator == "SAME" {
				// Проверяем, есть ли уже ссылка в операнде
				if t.Operand1.IsLink() {
					// Если ссылка уже есть, то мы просто заменяем ее на нужную
					t.Operand1.SetLink(*linkedTriad.Operand1.GetLink())
				} else {
					// Если ссылки нет, то мы создаем новую ссылку
					t.Operand1 = triad.LinkOperand(*linkedTriad.Operand1.GetLink())
				}
			}
		}
		if t.Operand1.IsLink() {
			linkedTriad := triads[*t.Operand2.GetLink()-1]
			if linkedTriad.Operator == "SAME" {
				// Проверяем, есть ли уже ссылка в операнде
				if t.Operand2.IsLink() {
					// Если ссылка уже есть, то мы просто заменяем ее на нужную
					t.Operand2.SetLink(*linkedTriad.Operand2.GetLink())
				} else {
					// Если ссылки нет, то мы создаем новую ссылку
					t.Operand2 = triad.LinkOperand(*linkedTriad.Operand2.GetLink())
				}
			}
		}

		// Шаг 2: Вычисление числа зависимости текущей триады
		depTriads[i] = 1 + max(calcDependency(t.Operand1, dep), calcDependency(t.Operand2, dep))

		// Шаг 3: Проверка на идентичность с более ранней триадой
		redundant := false
		for j := 0; j < i; j++ {
			if depTriads[i] == depTriads[j] && triads[i] == triads[j] {
				redundant = true
				result = append(result, triad.Triad{
					Operator: "SAME",
					Operand1: triad.LinkOperand(j + 1), // Ссылка на триаду с номером i
					Operand2: triad.NumberOperand(0),
				})
				break
			}
		}

		if !redundant {
			// Заменяем текущую триаду на SAME(j, 0)
			result = append(result, t)

		}

		// Шаг 4: Присваивание числа зависимости переменным
		if t.Operator == ":=" {
			dep[t.Operand1.GetOperand()] = i
		}
	}

	return result
}

// Функция для расчета числа зависимости для операндов
func calcDependency(operand triad.Operand, dep map[string]int) int {
	if operand.IsVariable() {
		// возвращаем зависимость переменной
		return dep[operand.GetOperand()]
	}
	// Для констант или значений без зависимости возвращаем 0
	return 0
}

// Вспомогательная функция для поиска максимума
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Функция для удаления всех SAME триад и обновления ссылок
func removeSameTriads(triads *[]triad.Triad) {
	var result []triad.Triad
	// Массив для отслеживания замененных ссылок
	// Ключ - это индекс операнда, который ссылается на SAME триаду, а значение - индекс, на который нужно обновить ссылку
	linkUpdates := make(map[int]int)

	// Проходим по всем триадам
	for i, t := range *triads {
		if t.Operator == "SAME" {
			// Если это SAME триада, то мы обновляем ссылку для всех последующих триад, которые на нее ссылаются
			if t.Operand1.IsLink() {
				// Получаем индекс ссылки на триаду SAME
				linkedIndex := *t.Operand1.GetLink()
				// Обновляем ссылку на настоящую триаду
				linkUpdates[i+1] = linkedIndex
			}
			// Пропускаем SAME триаду, так как она не должна быть в результате
			continue
		}

		// Проверяем, если операнды ссылаются на удаленную SAME триаду, заменяем ссылки
		if t.Operand1.IsLink() {
			linkedIndex := *t.Operand1.GetLink()
			if newLink, exists := linkUpdates[linkedIndex]; exists {
				// Обновляем ссылку на новую триаду
				t.Operand1 = triad.LinkOperand(newLink)
			}
		}

		if t.Operand2.IsLink() {
			linkedIndex := *t.Operand2.GetLink()
			if newLink, exists := linkUpdates[linkedIndex]; exists {
				// Обновляем ссылку на новую триаду
				t.Operand2 = triad.LinkOperand(newLink)
			}
		}

		// Добавляем триаду в итоговый результат
		result = append(result, t)
	}

	*triads = result
}
