package main

import (
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
	Parser = participle.MustBuild[Program](
		participle.Lexer(lex),
		participle.UseLookahead(2),
	)
)
