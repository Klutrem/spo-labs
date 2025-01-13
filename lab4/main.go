package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	ByteSize  = 1 // Размер типа byte в байтах
	RealSize  = 6 // Размер типа real в байтах
	Alignment = 2 // Кратность распределения (выравнивание)
)

// Функция для вычисления размера типа данных
func calculateTypeSize(typ *TypeDeclaration) int {
	totalSize := 0
	// Проходим по всем полям типа
	for _, field := range typ.Fields {
		switch field.Type {
		case "byte":
			totalSize += ByteSize // Если тип поля byte, добавляем 1 байт
		case "real":
			totalSize += RealSize // Если тип поля real, добавляем 6 байт
		}
	}
	// Выравнивание до кратности
	if totalSize%Alignment != 0 {
		totalSize += Alignment - (totalSize % Alignment) // Если размер не кратен Alignment, выравниваем
	}
	return totalSize
}

// Функция для вычисления размера переменной
func calculateVariableSize(varDecl *VarDeclaration, types map[string]int) int {
	size, exists := types[varDecl.Type]
	if !exists {
		// Если тип переменной неизвестен, выводим ошибку
		panic(fmt.Sprintf("Unknown type: %s", varDecl.Type))
	}
	if size%Alignment != 0 {
		size += Alignment - (size % Alignment) // Если размер не кратен Alignment, выравниваем
	}
	return size
}

func main() {
	// Проверяем, что файл передан в качестве аргумента
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		return
	}

	filename := os.Args[1]

	// Чтение содержимого файла
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Failed to read file: %s", err))
	}

	// Парсинг исходного текста в AST
	ast, err := Parser.ParseString("", string(data))
	if err != nil {
		panic(err) // Если произошла ошибка при парсинге, выводим её
	}

	// Выводим исходный текст программы
	fmt.Println("Исходный текст:")
	fmt.Println(string(data))

	// Расчет размеров типов данных
	types := map[string]int{
		"byte": ByteSize, // Определяем размер типа byte
		"real": RealSize, // Определяем размер типа real
	}
	fmt.Println("\nРазмеры типов данных:")
	// Для каждого типа данных рассчитываем его размер
	for _, typ := range ast.Types {
		size := calculateTypeSize(typ)
		types[typ.Name] = size
		fmt.Printf("Тип %s: %d байт\n", typ.Name, size)
	}

	// Расчет размеров переменных
	totalMemory := 0
	fmt.Println("\nРазмеры переменных:")
	// Для каждой переменной вычисляем её размер
	for _, varDecl := range ast.Vars {
		size := calculateVariableSize(varDecl, types)
		fmt.Printf("Переменная %s: %d байт\n", varDecl.Name, size)
		totalMemory += size
	}

	// Выводим суммарный объем памяти для всех переменных
	fmt.Printf("\nСуммарный объем памяти: %d байт\n", totalMemory)
}
