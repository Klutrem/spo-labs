package optimizer

import (
	"lab1_2/triad"
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
	*triads = eliminateRedundantOperations(*triads)
	removeSameTriads(triads)
}
