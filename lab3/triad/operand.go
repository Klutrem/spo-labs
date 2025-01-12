package triad

import (
	"fmt"
	"lab1_2/node"
	"lab1_2/types"
	"strconv"
	"strings"
)

type Operand struct {
	element string
	linkTo  *int
}

func (o Operand) GetOperand() string {
	if o.IsLink() {
		return fmt.Sprintf("^%v", *o.linkTo)
	}
	return o.element
}

func (o Operand) IsLink() bool {
	return o.linkTo != nil
}

func (o Operand) IsNumber() bool {
	return strings.ContainsAny(o.element, types.Numbers)
}

func (o Operand) IsVariable() bool {
	return strings.ContainsAny(o.element, types.Alphabet)
}

func (o Operand) GetLink() *int {
	return o.linkTo
}

func (o Operand) SetLink(link int) {
	o.linkTo = &link
}

func OperandFromString(s string) Operand {
	return Operand{
		element: s,
	}
}

func NumberOperand(n int) Operand {
	return Operand{
		element: strconv.Itoa(n),
	}
}

func OperandFromSimpleNode(n node.Node) Operand {
	return Operand{
		element: n.Lexem,
	}
}

func LinkOperand(index int) Operand {
	return Operand{
		linkTo: &index,
	}
}
