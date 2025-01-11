package triad

import (
	"fmt"
	"lab1_2/node"
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

func (o Operand) GetLink() *int {
	return o.linkTo
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
