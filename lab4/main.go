package main

import (
	"fmt"
	"log"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Структуры для грамматики

type Program struct {
	Types []*Type `@@*`
}

type Type struct {
	Name   string   `"type" @Ident`
	Fields []*Field `@@*`
}

type Field struct {
	VarName string `@Ident "="`
	Expr    Expr   `@@`
}

type Expr interface {
	Size(alignment int) int
}

type ScalarType struct {
	Type string `@("byte" | "real")`
}

func (s ScalarType) Size(alignment int) int {
	size := 0
	switch s.Type {
	case "byte":
		size = 1
	case "real":
		size = 6
	}

	// Выравнивание
	if size%alignment != 0 {
		padding := alignment - (size % alignment)
		size += padding
	}
	return size
}

type Union struct {
	Fields []*UnionField `"union" @@* "end"`
}

type UnionField struct {
	VarName string `@Ident ":"`
	Type    string `@("byte" | "real")`
}

func (u Union) Size(alignment int) int {
	totalSize := 0
	for _, field := range u.Fields {
		var fieldSize int
		switch field.Type {
		case "byte":
			fieldSize = 1
		case "real":
			fieldSize = 6
		}

		// Выравнивание
		if fieldSize%alignment != 0 {
			padding := alignment - (fieldSize % alignment)
			fieldSize += padding
		}

		totalSize += fieldSize
	}
	return totalSize
}

func main() {
	// Лексер для грамматики
	lexer := lexer.MustSimple([]lexer.SimpleRule{
		{`Ident`, `[a-zA-Z][a-zA-Z_\d]*`},
		{`Type`, `byte|real`},
		{`Union`, `union`},
		{`End`, `end`},
		{`Punct`, `[,]`},
		{`Whitespace`, `\s+`},
		{`Semicolon`, `;`}, // Точка с запятой
	})

	// Парсер с использованием библиотеки participle
	parser := participle.MustBuild[Program](
		participle.Lexer(lexer),
		participle.Union[Expr](ScalarType{}, Union{}),
	)

	// Пример входного текста
	input := `
type MyType
  var1 = byte;
  var2 = real;
  var3 = union
    var4: byte;
    var5: real;
  end

type AnotherType
  var6 = byte;
`

	// Вывод исходного текста
	fmt.Println("Исходный текст:")
	fmt.Println(input)

	// Парсинг строки
	program, err := parser.ParseString("", input)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	// Задаем кратность выравнивания (2 байта)
	alignment := 2

	// Вывод размеров типов данных
	fmt.Println("\nРазмеры типов данных:")
	for _, t := range program.Types {
		fmt.Printf("Тип '%s':\n", t.Name)
		for _, field := range t.Fields {
			size := field.Expr.Size(alignment)
			fmt.Printf("  Переменная '%s': %d байт (с учетом выравнивания)\n", field.VarName, size)
		}
	}

	// Вывод размеров объединений (если есть)
	for _, t := range program.Types {
		for _, field := range t.Fields {
			if union, ok := field.Expr.(Union); ok {
				fmt.Printf("  Объединение '%s': %d байт (с учетом выравнивания)\n", field.VarName, union.Size(alignment))
			}
		}
	}

	// Вычисление итогового объема памяти
	var totalSize int
	for _, t := range program.Types {
		for _, field := range t.Fields {
			totalSize += field.Expr.Size(alignment)
		}
	}

	// Итоговый объем памяти
	fmt.Printf("\nИтоговый объем памяти: %d байт (с учетом выравнивания)\n", totalSize)
}
