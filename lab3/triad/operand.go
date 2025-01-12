package triad

import (
	"fmt"
	"lab1_2/node"
	"lab1_2/types"
	"strconv"
	"strings"
)

// Operand представляет операнд триады.
// Это может быть:
// - Простое значение (например, число или переменная), хранящееся в поле `element`
// - Ссылка на другую триаду, хранящаяся в поле `linkTo`
type Operand struct {
	element string // Значение операнда (например, число, переменная или строка)
	linkTo  *int   // Указатель на индекс триады, на которую ссылается данный операнд
}

// GetOperand возвращает строковое представление операнда.
// Если это ссылка, то возвращается "^<индекс>".
// Если это значение, то возвращается само значение.
func (o Operand) GetOperand() string {
	if o.IsLink() {
		return fmt.Sprintf("^%v", *o.linkTo) // Формат ссылки
	}
	return o.element // Простое значение
}

// IsLink проверяет, является ли операнд ссылкой на другую триаду.
func (o Operand) IsLink() bool {
	return o.linkTo != nil
}

// IsNumber проверяет, является ли операнд числом.
// Использует набор символов, определённый в `types.Numbers`.
func (o Operand) IsNumber() bool {
	return strings.ContainsAny(o.element, types.Numbers)
}

// IsVariable проверяет, является ли операнд переменной.
// Использует набор символов, определённый в `types.Alphabet`.
func (o Operand) IsVariable() bool {
	return strings.ContainsAny(o.element, types.Alphabet)
}

// GetLink возвращает указатель на индекс ссылки операнда.
// Если операнд не является ссылкой, возвращает `nil`.
func (o Operand) GetLink() *int {
	return o.linkTo
}

// SetLink устанавливает ссылку на указанную триаду.
func (o Operand) SetLink(link int) {
	o.linkTo = &link
}

// OperandFromString создаёт операнд из строки.
// Пример: "x" -> Operand{element: "x"}
func OperandFromString(s string) Operand {
	return Operand{
		element: s,
	}
}

// NumberOperand создаёт операнд из числа.
// Число преобразуется в строку для хранения в `element`.
// Пример: 42 -> Operand{element: "42"}
func NumberOperand(n int) Operand {
	return Operand{
		element: strconv.Itoa(n),
	}
}

// OperandFromSimpleNode создаёт операнд из узла (Node).
// Использует поле `Lexem` узла как значение операнда.
// Пример: Node{Lexem: "x"} -> Operand{element: "x"}
func OperandFromSimpleNode(n node.Node) Operand {
	return Operand{
		element: n.Lexem,
	}
}

// LinkOperand создаёт операнд-ссылку на указанную триаду.
// Пример: 3 -> Operand{linkTo: &3}
func LinkOperand(index int) Operand {
	return Operand{
		linkTo: &index,
	}
}
