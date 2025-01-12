package main

import (
	"log"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
)

// Определим структуры для грамматики
type Program struct {
	Type string `parser:"'type'"`
	L    []Type `parser:"@@*"`
	Var  string `parser:"'var'"`
	R    []Var  `parser:"@@*"`
}

type Type struct {
	T string `parser:"'t' '=' 'c' | 't' '=' 'D'"`
}

type Var struct {
	K     string `parser:"'K'"`
	Colon string `parser:":'c' | 't' | 'D'"`
}

type D struct {
	Union string `parser:"'union'"`
	F     []F    `parser:"@@*"`
	End   string `parser:"'end'"`
}

type F struct {
	E []E `parser:"@@*"`
}

type E struct {
	K     string `parser:"'K'"`
	Colon string `parser:":'c' | 't'"`
}

var (
	// Определим лексер с состояниями
	def = lexer.MustStateful(lexer.Rules{
		"Root": {
			{Name: `Type`, Pattern: `type`, Action: nil},
			{Name: `Var`, Pattern: `var`, Action: nil},
		},
		"Type": {
			{Name: "T", Pattern: `t`, Action: nil},
			{Name: "Equal", Pattern: `=`, Action: nil},
			{Name: "C", Pattern: `c`, Action: nil},
			{Name: "D", Pattern: `D`, Action: nil},
		},
		"Var": {
			{Name: "K", Pattern: `K`, Action: nil},
			{Name: "Colon", Pattern: `:`, Action: nil},
			{Name: "C", Pattern: `c`, Action: nil},
			{Name: "T", Pattern: `t`, Action: nil},
			{Name: "D", Pattern: `D`, Action: nil},
		},
		"Union": {
			{Name: "Union", Pattern: `union`, Action: nil},
			{Name: "End", Pattern: `end`, Action: nil},
		},
	})

	// Построим парсер
	parser = participle.MustBuild[Program](participle.Lexer(def),
		participle.Elide("Whitespace"))
)

func main() {
	// Тестовая строка
	input := `type t=c var K: t K: c`

	// Парсим строку
	pr, err := parser.ParseString("err", input)
	if err != nil {
		log.Fatal(err)
	}

	// Печатаем результат
	repr.Println(pr)
}
