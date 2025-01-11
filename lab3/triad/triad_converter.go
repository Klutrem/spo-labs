package triad

import (
	"lab1_2/node"
)

func ConvertNodeToTriads(nodeToConvert node.Node) []Triad {
	var triads []Triad

	if nodeToConvert.IsArithmetic() {
		triads = append(triads, ConvertArithmeticNodeToTriad(nodeToConvert))
	}
	return triads
}

func ConvertArithmeticNodeToTriad(nodeToConvert node.Node) Triad {
	return Triad{
		Operand:   nodeToConvert.Children[1].Lexem,
		Operator1: nodeToConvert.Children[0].Lexem,
		Operator2: nodeToConvert.Children[2].Lexem,
	}
}
