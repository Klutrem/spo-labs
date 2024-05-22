package lab1

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Константы для разделителей и наборов символов
const (
	Delimiter     = ";"                                                              // Символ разделителя
	Alphabet      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"           // Алфавит
	Alphanumeric  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // Алфавит и цифры
	Numbers       = "0123456789"                                                     // Цифры
	OperatorChars = "+-*/"                                                           // Символы операторов
	Parentheses   = "()"                                                             // Символы скобок
)

// Типы лексем
const (
	DelimiterType   = "delimiter"   // Разделитель
	IdentifierType  = "identifier"  // Идентификатор
	ConstType       = "const"       // Константа
	AssignmentType  = "assignment"  // Присваивание
	OperatorType    = "operator"    // Оператор
	ParenthesesType = "parentheses" // Скобки
	NumberType      = "number"      // Число
	ErrorType       = "error"       // Ошибка
	CommentType     = "comment"     // Комментарий
)

// Структура TokenSaver для хранения и печати токенов
type TokenSaver struct {
	tokens []Token
}

// Структура Token представляет лексему с ее типом и значением
type Token struct {
	Type  string // Тип
	Value string // Значение
}

// Add добавляет токен в TokenSaver
func (ts *TokenSaver) Add(tokenType, value string) {
	ts.tokens = append(ts.tokens, Token{Type: tokenType, Value: value})
}

// Print выводит все токены из TokenSaver
func (ts *TokenSaver) Print() {
	fmt.Println("Токены:")
	for _, token := range ts.tokens {
		fmt.Printf("{ Тип: %s, Значение: %s }\n", token.Type, token.Value)
	}
}

func main() {
	RunLexer()
}

// RunLexer обрабатывает ввод и генерирует токены
func RunLexer() []Token {
	state := DelimiterType         // Начальное состояние - разделитель
	token := ""                    // Текущий токен
	tokenSaver := &TokenSaver{}    // Сохранитель токенов
	characters := readCharacters() // Получаем ввод

	for i := 0; i < len(characters); i++ {
		char := string(characters[i]) // Текущий символ
		switch state {
		case DelimiterType:
			handleDelimiter(char, tokenSaver) // Обработка разделителя
			if char != " " && char != "\n" && char != "\r" {
				token = char                // Начинаем новый токен
				state = GetLexemeType(char) // Определяем тип лексемы
			}
		case IdentifierType, NumberType, OperatorType, ParenthesesType, ErrorType:
			if IsValid(char, state) {
				token += char // Продолжаем собирать текущий токен
			} else {
				tokenSaver.Add(state, token) // Завершаем текущий токен и сохраняем
				token = ""                   // Начинаем новый токен
				state = DelimiterType        // Возвращаемся к обработке разделителей
				i--                          // Повторно обрабатываем текущий символ
			}
		case AssignmentType:
			handleAssignment(char, tokenSaver, &state, &token, &i, characters) // Обработка оператора присваивания
		case ConstType:
			handleConst(char, tokenSaver, &state, &token, &i, characters) // Обработка константы
		case CommentType:
			if char == "\n" {
				state = DelimiterType // Обработка комментария до конца строки
			}
		}
	}
	if token != "" && state != CommentType && token != "#" {
		tokenSaver.Add(state, token) // Добавляем последний токен, если он есть
	}
	tokenSaver.Print() // Печать всех токенов
	return tokenSaver.tokens
}

// handleDelimiter обрабатывает разделитель
func handleDelimiter(char string, tokenSaver *TokenSaver) {
	if char == Delimiter {
		tokenSaver.Add(DelimiterType, Delimiter) // Добавляем разделитель в токены
	}
}

// GetLexemeType возвращает тип лексемы для заданного символа
func GetLexemeType(char string) string {
	if strings.Contains(Alphabet, char) {
		return IdentifierType // Идентификатор
	} else if strings.Contains(Numbers, char) {
		return NumberType // Число
	} else if char == "'" {
		return ConstType // Константа
	} else if char == "#" {
		return CommentType // Комментарий
	} else if strings.Contains(OperatorChars, char) {
		return OperatorType // Оператор
	} else if strings.Contains(Parentheses, char) {
		return ParenthesesType // Скобки
	} else if char == ";" {
		return DelimiterType // Разделитель
	} else if char == ":" {
		return AssignmentType // Оператор присваивания
	}
	return ErrorType // Ошибка
}

// IsValid проверяет, является ли символ допустимым в текущем состоянии
func IsValid(char string, state string) bool {
	switch state {
	case IdentifierType:
		return strings.Contains(Alphanumeric, char) // Допустимы буквы и цифры
	case NumberType:
		return strings.Contains(Numbers, char) // Допустимы только цифры
	case ConstType:
		return char != "'" // Допустимы любые символы, кроме '
	case OperatorType:
		return strings.Contains(OperatorChars, char) // Допустимы только символы операторов
	case ParenthesesType:
		return strings.Contains(Parentheses, char) // Допустимы только символы скобок
	}
	return false // В остальных случаях недопустимо
}

// handleAssignment обрабатывает оператор присваивания
func handleAssignment(char string, tokenSaver *TokenSaver, state *string, token *string, i *int, characters []byte) {
	if char == "=" && (*i)+1 < len(characters) && characters[(*i)-1] == ':' {
		*state = AssignmentType                // Устанавливаем тип лексемы - оператор присваивания
		*token = ":="                          // Задаем значение токена как оператор присваивания
		(*i)++                                 // Пропускаем символ '=', так как он уже обработан
		tokenSaver.Add(AssignmentType, *token) // Добавляем оператор присваивания в токены
		*state = DelimiterType                 // Возвращаемся к обработке разделителей
	} else {
		tokenSaver.Add(OperatorType, *token) // Добавляем текущий токен как оператор
		*token = ""                          // Начинаем новый токен
		*state = ErrorType                   // Устанавливаем тип лексемы - ошибка
	}
}

// handleConst обрабатывает константу
func handleConst(char string, tokenSaver *TokenSaver, state *string, token *string, i *int, characters []byte) {
	if strings.Contains(Alphabet, char) && checkConstBrackets(*i, characters) {
		*state = ConstType                  // Устанавливаем тип лексемы - константа
		*token = fmt.Sprint("'", char, "'") // Задаем значение токена как константу
		(*i)++                              // Пропускаем символ константы
		(*i)++                              // Пропускаем символ "'", так как он уже обработан
		tokenSaver.Add(ConstType, *token)   // Добавляем константу в токены
		*state = DelimiterType              // Возвращаемся к обработке разделителей
	} else {
		tokenSaver.Add(ErrorType, *token) // Добавляем текущий токен как ошибку
		*token = ""                       // Начинаем новый токен
		*state = ErrorType                // Устанавливаем тип лексемы - ошибка
	}
}

// checkConstBrackets проверяет соответствие скобок в константе
func checkConstBrackets(i int, characters []byte) bool {
	return (string(characters[i-1]) == "'" && string(characters[i+1]) == "'") // Проверяем наличие скобок для константы
}

// readCharacters считывает символы из стандартного ввода
func readCharacters() []byte {
	reader := bufio.NewReader(os.Stdin) // Создаем новый Reader для стандартного ввода
	var characters []byte
	for {
		input, err := reader.ReadString('\n') // Считываем строку из ввода
		if err != nil {
			break // Завершаем цикл, если произошла ошибка или достигнут конец ввода
		}
		characters = append(characters, []byte(input)...) // Добавляем символы из строки в массив символов
	}
	return characters // Возвращаем массив символов
}
