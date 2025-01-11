package triad

import "fmt"

type Triad struct {
	Operand   string
	Operator1 string
	Operator2 string
}

func (t Triad) ToString() string {
	return fmt.Sprintf("%s(%s, %s)", t.Operand, t.Operator1, t.Operator2)
}
