package main

import "github.com/irustm/ng-template-parser/ep"

func main() {
	lex := ep.Lexer{}
	tokens := lex.Tokenize("rr+jama")

	println(len(tokens))
}
