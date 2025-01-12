package main

import (
	"fmt"
	"io/ioutil"
	"os"
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
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		return
	}

	filename := os.Args[1]

	// Чтение файла
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Failed to read file: %s", err))
	}

	// Парсинг входного текста
	ast, err := Parser.ParseString("", string(data))
	if err != nil {
		panic(err)
	}

	// Вывод исходного текста
	fmt.Println("Исходный текст:")
	fmt.Println(string(data))

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
