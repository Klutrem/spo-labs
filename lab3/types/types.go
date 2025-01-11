package types

const (
	Delimiter     = ";"                                                              // Символ разделителя
	Alphabet      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"           // Алфавит
	Alphanumeric  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // Алфавит и цифры
	Numbers       = "0123456789"                                                     // Цифры
	OperatorChars = "+-*/"                                                           // Символы операторов
	Parentheses   = "()"                                                             // Символы скобок
	Identifier    = ":="
)

const (
	Terminal    = "terminal"
	NonTerminal = "NonTerminal"
)
