package main

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Program struct {
	Pos lexer.Position

	Types []*TypeDeclaration `@@*`
	Vars  []*VarDeclaration  `@@*`
}

type TypeDeclaration struct {
	Pos lexer.Position

	TypeKeyword string   `"type"`
	Name        string   `@Ident`
	Fields      []*Field `"{" @@* "}"`
}

type VarDeclaration struct {
	Pos lexer.Position

	VarKeyword string `"var"`
	Name       string `@Ident`
	Type       string `"=" @Ident ";"`
}

type Field struct {
	Pos lexer.Position

	Name string `@Ident`
	Type string `":" @("byte" | "real") ";"`
}

var (
	lex = lexer.MustSimple([]lexer.SimpleRule{
		{"comment", `//.*|/\*.*?\*/`},
		{"whitespace", `\s+`},

		{"TypeKeyword", `\btype\b`},
		{"VarKeyword", `\bvar\b`},
		{"Ident", `[a-zA-Z_][a-zA-Z0-9_]*`},
		{"Punct", `[-,(){}=;:]`},
		{"Literal", `\b(byte|real)\b`},
	})
	parser = participle.MustBuild[Program](
		participle.Lexer(lex),
		participle.UseLookahead(2),
	)
)

const (
	ByteSize  = 1 // Размер типа byte в байтах
	RealSize  = 6 // Размер типа real в байтах
	Alignment = 2 // Кратность распределения
)

func calculateTypeSize(typ *TypeDeclaration) int {
	totalSize := 0
	for _, field := range typ.Fields {
		switch field.Type {
		case "byte":
			totalSize += ByteSize
		case "real":
			totalSize += RealSize
		}
	}
	// Выравнивание до кратности
	if totalSize%Alignment != 0 {
		totalSize += Alignment - (totalSize % Alignment)
	}
	return totalSize
}

func calculateVariableSize(varDecl *VarDeclaration, types map[string]int) int {
	size, exists := types[varDecl.Type]
	if !exists {
		panic(fmt.Sprintf("Unknown type: %s", varDecl.Type))
	}
	return size
}

func main() {
	const sample = `
type MyType {
    field1 : byte;
    field2 : real;
}

type AnotherType {
    field3 : real;
    field4 : byte;
}

var var1 = MyType;
var var2 = AnotherType;
var var3 = byte;
var var4 = real;
`

	// Парсинг входного текста
	ast, err := parser.ParseString("", sample)
	if err != nil {
		panic(err)
	}

	// Вывод исходного текста
	fmt.Println("Исходный текст:")
	fmt.Println(sample)

	// Расчет размеров типов
	types := map[string]int{
		"byte": ByteSize,
		"real": RealSize,
	}
	fmt.Println("\nРазмеры типов данных:")
	for _, typ := range ast.Types {
		size := calculateTypeSize(typ)
		types[typ.Name] = size
		fmt.Printf("Тип %s: %d байт\n", typ.Name, size)
	}

	// Расчет размеров переменных
	totalMemory := 0
	fmt.Println("\nРазмеры переменных:")
	for _, varDecl := range ast.Vars {
		size := calculateVariableSize(varDecl, types)
		fmt.Printf("Переменная %s: %d байт\n", varDecl.Name, size)
		totalMemory += size
	}

	// Вывод суммарного объема памяти
	fmt.Printf("\nСуммарный объем памяти: %d байт\n", totalMemory)
}
