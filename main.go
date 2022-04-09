package main

import (
	"akt-wxss/parser"
	"fmt"
)

func main() {
	p := parser.NewParser()

	s := []byte(`body {
		background-color:#d0e4fe;
	}
	h1 {
		color:orange;
		text-align:center;
	}
	p {
		font-family:"Times New Roman";
		font-size:20rpx;
	}`)

	ast := p.Parse(s)

	ast.Traverse(func(node *parser.CSSDefinition) {
		fmt.Printf("%v", node)

		node.Selector.Selector = ".b"

		fmt.Printf("%v", node)
	})

	mini := ast.Minisize()

	fmt.Printf("%s", mini.String())

}
