package main

import (
	"encoding/json"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type TextAttribute struct {
	Name  string
	Value string
}

type Reference struct {
	Name  string
	Value string
}

// https://github.com/angular/angular/blob/e112e320bf6c2b60e8ecea46f80bcaec593c65b7/packages/compiler/src/expression_parser/ast.ts
type BindingType int

const (
	BindingTypeProperty BindingType = iota
	BindingTypeAttribute
	BindingTypeClass
	BindingTypeStyle
	BindingTypeAnimation
)

type BoundAttribute struct {
	Name        string
	BindingType BindingType
	Value       AstWithSourcePropertyRead
}

type BoundEvent struct {
	Name        string
	BindingType BindingType
	Handler     interface{}
}

type Interpolation struct {
	Strings     []string
	Expressions []PropertyRead
}
type BoundText struct {
	Value AstWithSourceInterpolation
}

type AstWithSourcePropertyRead struct {
	Ast    PropertyRead
	Source string
}

type AstWithSourceMethodCall struct {
	Ast    MethodCall
	Source string
}

type AstWithSourceInterpolation struct {
	Ast    Interpolation
	Source string
}

type AstWithSourcePropertyWrite struct {
	Ast    PropertyWrite
	Source string
}

// Ast parsed types

type PropertyRead struct {
	Name string
}

type PropertyWrite struct {
	Name  string
	Value PropertyRead
}

type MethodCall struct {
	Name string
	Args []PropertyRead
}

type Text struct {
	Value string
}

type Comment struct {
	Value string
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
	//r := strings.NewReader(template)

	f, _ := os.Open("template.html")

	t := time.Now()

	tokenizer := html.NewTokenizer(f)

	start := time.Now()

	data := parse(tokenizer)

	elapsed := t.Sub(start)
	println(elapsed.Milliseconds())

	file, _ := json.Marshal(data)
	//
	_ = ioutil.WriteFile("out.json", file, 0644)
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

	if tokenType == html.TextToken {
		data := token.Data
		tokenizer.Next()

		if strings.Contains(data, "{{") {
			var strings []string
			var expressions []PropertyRead
			var buff []uint8

			for i := 0; i < len(data); i++ {
				if data[i] == '{' {
					// add latest
					strings = append(strings, string(buff))
					buff = nil
					i += 1
				} else if data[i] == '}' {
					// add latest
					expressions = append(expressions, PropertyRead{Name: string(buff)})
					buff = nil
					i += 1
				} else {
					buff = append(buff, data[i])
				}
			}

			return BoundText{Value: AstWithSourceInterpolation{
				Ast:    Interpolation{Strings: strings, Expressions: expressions},
				Source: data,
			}}
		}

		return Text{Value: data}
	}

	if tokenType == html.CommentToken {
		tokenizer.Next()
		return Comment{Value: token.Data}
	}

	if tokenType == html.StartTagToken {
		element := Element{Name: token.Data}

		// parse attributes
		for _, attr := range token.Attr {
			// Reference
			if attr.Key[0] == '#' {
				element.References = append(element.References, Reference{Value: attr.Val})

				// Output
			} else if attr.Key[0] == '(' {
				name := attr.Key[1 : len(attr.Key)-1]
				source := attr.Val
				sourceSplit := strings.Split(source, "(")

				handlerName := sourceSplit[0]
				propName := sourceSplit[1][:len(sourceSplit[1])-1]

				var outputArgs []PropertyRead
				outputArgs = append(outputArgs, PropertyRead{Name: propName})

				element.Outputs = append(element.Outputs,
					BoundEvent{
						Name:    name,
						Handler: AstWithSourceMethodCall{Source: source, Ast: MethodCall{Name: handlerName, Args: outputArgs}},
					})

				// Input
			} else if attr.Key[0] == '[' && attr.Key[1] != '(' {
				name := attr.Key[1 : len(attr.Key)-1]
				// [class.asd]
				nameStrings := strings.Split(name, ".")
				inputTypeString := nameStrings[0]
				bindingType := BindingTypeProperty

				if inputTypeString == "class" {
					bindingType = BindingTypeClass
				}

				if inputTypeString == "style" {
					bindingType = BindingTypeStyle
				}

				if inputTypeString == "attr" {
					bindingType = BindingTypeAttribute
				}

				if bindingType != BindingTypeProperty {
					name = nameStrings[1]
				}

				element.Inputs = append(element.Inputs,
					BoundAttribute{
						Name:        name,
						BindingType: bindingType,
						Value:       AstWithSourcePropertyRead{Source: attr.Val, Ast: PropertyRead{Name: attr.Val}},
					})

				// Input / Output
			} else if attr.Key[0] == '[' && attr.Key[1] == '(' {
				name := attr.Key[2 : len(attr.Key)-2]

				element.Inputs = append(element.Inputs,
					BoundAttribute{
						Name:        name,
						BindingType: BindingTypeProperty,
						Value: AstWithSourcePropertyRead{
							Source: attr.Val,
							Ast:    PropertyRead{Name: attr.Val},
						},
					})

				handlerName := name + "Change"
				ast := PropertyWrite{Name: attr.Val, Value: PropertyRead{Name: "$event"}}

				element.Outputs = append(element.Outputs,
					BoundEvent{
						Name:        handlerName,
						BindingType: BindingTypeProperty,
						Handler:     AstWithSourcePropertyWrite{Source: attr.Val + "=$event", Ast: ast},
					})
			} else {
				element.Attributes = append(element.Attributes,
					TextAttribute{
						Name:  attr.Val,
						Value: attr.Key,
					})

				// TODO Error
				//fmt.Println(attr.Key, " = ", attr.Val)
			}
		}

		tokenType = tokenizer.Next()

		for tokenType != html.EndTagToken {
			token = tokenizer.Token()
			element.Children = append(element.Children, walk(tokenizer, token))
			tokenType = tokenizer.Token().Type
		}

		tokenType = tokenizer.Next()

		return element
	}

	return nil
}
