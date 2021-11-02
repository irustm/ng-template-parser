package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"strings"
)

type TextAttribute struct {
	name  string
	value string
}

type Reference struct {
	name  string
	value string
}

// https://github.com/angular/angular/blob/e112e320bf6c2b60e8ecea46f80bcaec593c65b7/packages/compiler/src/expression_parser/ast.ts
type BindingType int

const (
	Property BindingType = iota
	Attribute
	Class
	Style
	Animation
)

type BoundAttribute struct {
	name        string
	bindingType BindingType
	value       AstWithSource
}

type BoundEvent struct {
	name    string
	handler AstWithSource
}

type Interpolation struct {
	strings     []string
	expressions []PropertyRead
}
type BoundText struct {
	value AstWithSource
}

type AstParsed struct {
	name string
	args []interface{}
}

type AstWithSource struct {
	ast    AstParsed
	source string
}

// Ast parsed types
type PropertyRead struct {
	name string
}

type PropertyWrite struct {
	name  string
	value PropertyRead
}

type Text struct {
	value string
}

type Comment struct {
	value string
}

type Element struct {
	Name       string
	Attributes []TextAttribute
	Inputs     []BoundAttribute
	Outputs    []BoundEvent
	References []Reference
	Children   []interface{}
}

type Root struct {
	Nodes []interface{}
}

func main() {
	template := `
       <div #ida1 [class.asd]="containera" log="1" (click)="onClick($event)" [(aa)]="model">
		{{test}} 
       </div>`

	r := strings.NewReader(template)
	tokenizer := html.NewTokenizer(r)

	parse(tokenizer)

	//for _, element := range root.Nodes {
	//	fmt.Println(element)
	//}

}

func parse(tokenizer *html.Tokenizer) Root {
	root := Root{}
	tokenizer.Next()

	for {
		token := tokenizer.Token()

		if token.Type == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				//end of the file, break out of the loop
				break
			}
			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
		}

		root.Nodes = append(root.Nodes, walk(tokenizer, token))
	}

	return root
}

func walk(tokenizer *html.Tokenizer, token html.Token) interface{} {
	tokenType := token.Type

	println(token.String())

	if tokenType == html.TextToken {
		tokenizer.Next()
		return Text{value: token.Data}
	}

	if tokenType == html.CommentToken {
		tokenizer.Next()
		return Comment{value: token.Data}
	}

	if tokenType == html.StartTagToken {
		element := Element{Name: token.Data}

		// parse attributes
		for _, attr := range token.Attr {
			// Reference
			if attr.Key[0] == '#' {
				element.References = append(element.References, Reference{value: attr.Val})

				// Output
			} else if attr.Key[0] == '(' {
				name := attr.Key[1 : len(attr.Key)-1]

				element.Outputs = append(element.Outputs,
					BoundEvent{
						name:    name,
						handler: AstWithSource{source: attr.Val},
					})

				// Input
			} else if attr.Key[0] == '[' && attr.Key[1] != '(' {
				name := attr.Key[1 : len(attr.Key)-1]

				element.Inputs = append(element.Inputs,
					BoundAttribute{
						name:  name,
						value: AstWithSource{source: attr.Val},
					})

				// Input / Output
			} else if attr.Key[0] == '[' && attr.Key[1] == '(' {
				name := attr.Key[2 : len(attr.Key)-2]

				element.Inputs = append(element.Inputs,
					BoundAttribute{
						name:  name,
						value: AstWithSource{source: attr.Val},
					})
				element.Outputs = append(element.Outputs,
					BoundEvent{
						name:    name,
						handler: AstWithSource{source: attr.Val},
					})
			} else {
				// Error
				// TODO parse attributes
				fmt.Println(attr.Key, " = ", attr.Val)
			}
		}

		tokenType = tokenizer.Next()
		token = tokenizer.Token()

		for tokenType != html.EndTagToken {
			element.Children = append(element.Children, walk(tokenizer, token))
			tokenType = tokenizer.Token().Type
		}

		tokenType = tokenizer.Next()

		return element
	}

	return nil
}
