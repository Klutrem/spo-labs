package main

import (
	"encoding/json"
	"fmt"
	"lab1_2/code_generation"
	"lab1_2/node"
	"lab1_2/triad"
	"log"
	"os"
)

func main() {
	var nodes []node.Node
	// читаем файл, сгенерированный алгоритмом из прошлых ЛР
	fileData, err := os.ReadFile("../output.json")
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}

	if err := json.Unmarshal(fileData, &nodes); err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}
	println("Начальные выражения:")
	for _, node := range nodes {
		// Вывод начальной лексеммы
		fmt.Printf("%+v\n", node.Lexem)
	}

	var doubleTriads [][]triad.Triad

	println("Триады:")
	for _, node := range nodes {
		// Вывод начальной лексеммы
		var triads []triad.Triad
		triad.ConvertNodeToTriads(node, &triads)
		// Печать триад
		doubleTriads = append(doubleTriads, triads)
	}

	resultTriads := triad.MergeTriadList(doubleTriads...)
	for i, t := range resultTriads {
		fmt.Printf("%d: %s\n", i+1, t.ToString())
	}
	println("Код до оптимизации:")
	res := code_generation.GenerateAssemblyCode(resultTriads)
	println(res)
}
