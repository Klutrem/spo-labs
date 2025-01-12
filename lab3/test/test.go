package main

import (
	"lab1_2/node"
	"lab1_2/triad"
)

func main() {
	// Пример выражения a := b - c
	nodeToConvert := node.Node{
		Type:  "NonTerminal",
		Lexem: "a := b - c ",
		Children: []node.Node{
			{
				Type:  "Terminal",
				Lexem: "a",
			},
			{
				Type:  "Terminal",
				Lexem: ":=",
			},
			{
				Type:  "NonTerminal",
				Lexem: "b - c + d",
				Children: []node.Node{
					{
						Type:  "NonTerminal",
						Lexem: "c + d",
						Children: []node.Node{
							{
								Type:  "Terminal",
								Lexem: "c",
							},
							{
								Type:  "Terminal",
								Lexem: "+",
							},
							{
								Type:  "Terminal",
								Lexem: "d",
							},
						},
					},

					{
						Type:  "Terminal",
						Lexem: "-",
					},

					{
						Type:  "Terminal",
						Lexem: "b",
					},
				},
			},
		},
	}
	var triads []triad.Triad
	triad.ConvertNodeToTriads(nodeToConvert, &triads)

}
