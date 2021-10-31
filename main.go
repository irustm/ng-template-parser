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

type BoundAttribute struct {
	name  string
	value AstWithSource
}

type BoundEvent struct {
	name    string
	handler AstWithSource
}

type AstWithSource struct {
	source string
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
			if attr.Key[0] == '#' {
				element.References = append(element.References, Reference{value: attr.Val})
			} else {
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
