package main

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Program представляет собой корневую структуру для всего программы.
// Она содержит массивы типов данных и переменных, определенных в программе.
type Program struct {
	Pos lexer.Position // Позиция в исходном файле, где начинается описание программы.

	Types []*TypeDeclaration `@@*` // Массив типов данных, которые встречаются в программе.
	Vars  []*VarDeclaration  `@@*` // Массив переменных, которые встречаются в программе.
}

// TypeDeclaration описывает определение типа данных.
// Содержит ключевое слово "type", имя типа и список его полей.
type TypeDeclaration struct {
	Pos lexer.Position // Позиция типа в исходном файле.

	TypeKeyword string   `"type"` // Ключевое слово "type", обязательное для определения типа.
	Name        string   `@Ident` // Имя типа данных.
	Fields      []*Field `"{" @@* "}"`
	// Список полей типа данных, которые могут быть разных типов.
}

// VarDeclaration описывает объявление переменной.
// Содержит ключевое слово "var", имя переменной и тип с присваиванием.
type VarDeclaration struct {
	Pos lexer.Position // Позиция переменной в исходном файле.

	VarKeyword string `"var"`  // Ключевое слово "var", обязательное для объявления переменной.
	Name       string `@Ident` // Имя переменной.
	Type       string `"=" @Ident ";"`
	// Тип переменной, присваиваемый через знак "=".
}

// Field описывает поле внутри типа данных (например, внутри структуры).
// Каждое поле имеет имя и тип (byte, real).
type Field struct {
	Pos lexer.Position // Позиция поля в исходном файле.

	Name string `@Ident` // Имя поля.
	Type string `":" @("byte" | "real") ";"`
	// Тип поля (byte или real), обязательно завершается точкой с запятой.
}

// lex содержит правила для лексического анализа входного текста программы.
// Сюда входят правила для комментариев, пробельных символов и основных токенов.
var (
	lex = lexer.MustSimple([]lexer.SimpleRule{
		{"comment", `//.*|/\*.*?\*/`},       // Правило для комментариев в программе (однострочные и многострочные).
		{"whitespace", `\s+`},               // Правило для пробельных символов, которые мы игнорируем.
		{"TypeKeyword", `\btype\b`},         // Правило для ключевого слова "type".
		{"VarKeyword", `\bvar\b`},           // Правило для ключевого слова "var".
		{"Ident", `[a-zA-Z_][a-zA-Z0-9_]*`}, // Правило для идентификаторов (имена переменных и типов).
		{"Punct", `[-,(){}=;:]`},            // Правило для знаков препинания (например, скобки, запятые).
		{"Literal", `\b(byte|real)\b`},      // Правило для типов данных "byte" и "real".
	})
	// Parser создает парсер, который использует определенный лексер для обработки входных данных.
	Parser = participle.MustBuild[Program](
		participle.Lexer(lex),      // Используем лексер, определенный выше.
		participle.UseLookahead(2), // Используем двухсимвольный просмотр для корректного парсинга.
	)
)
