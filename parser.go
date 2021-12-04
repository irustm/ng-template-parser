package main

import (
	"encoding/json"
	"github.com/SimplePEG/Go/rd"
	"github.com/SimplePEG/Go/speg"
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

var expressionParser rd.ParserFunc
var expressionRule speg.GrammarRules

// https://github.com/angular/angular/blob/master/packages/compiler/src/expression_parser/ast.ts
type BindingType int

const (
	BindingTypeProperty BindingType = iota
	BindingTypeAttribute
	BindingTypeClass
	BindingTypeStyle
	BindingTypeAnimation
)

type LiteralPrimitive struct {
	value interface{}
}

type BinaryCondition struct {
	Operation string
	Left      interface{}
	Right     interface{}
}

type Conditional struct {
	Condition interface{}
	trueExp   interface{}
	falseExp  interface{}
}

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
	Expressions []rd.Ast
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

func ParsePegExpression(text string) rd.Ast {
	result, _ := expressionParser(&rd.State{
		Text:     text,
		Position: 0,
		Rules:    expressionRule.Rules,
	})

	return result
}

func templateParse() {
	//r := strings.NewReader(template)
	var spegParser = speg.NewSPEGParser()
	var gAst, gErr = spegParser.ParseGrammar(GetNgParserGrammar())
	if !gErr {
		gparser, grule := speg.GetParser(gAst)

		expressionParser = gparser
		expressionRule = grule
	}

	f, _ := os.Open("template.html")

	tokenizer := html.NewTokenizer(f)

	println("started")

	start := time.Now()

	data := parse(tokenizer)

	endTime := time.Now()

	elapsed := endTime.Sub(start)

	println(elapsed.String())

	//file, _ := json.MarshalIndent(data, "", " ")
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
			var expressions []rd.Ast
			var buff []uint8

			for i := 0; i < len(data); i++ {
				if data[i] == '{' {
					// add latest
					strings = append(strings, string(buff))
					buff = nil
					i += 1
				} else if data[i] == '}' {
					expressions = append(expressions, ParsePegExpression(string(buff)))
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

func GetNgParserGrammar() string {
	return `GRAMMAR angular_expression

root         -> expression EOF; 
expression   -> ws (conditional / condition / binary / literal_map / property_array / method / value) ws pipe? ws;

conditional  -> ws (condition / binary / method / value) ws "?" ws trueExp ws ":" ws falseExp ws;

condition ->  (binary / method / value) ws condition_operation ws (binary / method / value);
condition_operation -> "===" / "==" / "!==" / "!=" / ">" / "<" ;

trueExp -> conditional / binary / method / value;
falseExp -> conditional / binary / method / value;

method -> (method_call / method_call_with_args) ;
method_call -> property_read "()" property_keyread_reader?;
method_call_with_args -> property_read "(" method_args ")";
method_args -> method_arg (value_separator method_arg)*;
method_arg -> conditional / binary / method / value;

value ->
   primitive_values / property_read / array / literal_map / literal_map_empty;

binary -> left ws (operation ws right)+;
left -> method / value / ("(" binary ")") / binary;
right -> method / value / ("(" binary ")") / binary;

pipe -> pipe_literal+;
pipe_literal -> "|" ws pipe_name ws (":" pipe_args)? ws;
pipe_name -> [A-Za-z0-9]+;
pipe_args -> value;

literal_map -> "{" ws string ws ":" ws value ws "}";
literal_map_empty -> "{" ws "}";

begin_array -> ws "[" ws;
end_array -> ws "]" ws;
value_separator -> ws "," ws;

array ->
    begin_array (value (value_separator value)*)? end_array;
    
property_array -> 
    begin_array ((conditional / value ) (value_separator (conditional / value))*)? end_array;

primitive_values ->
    true /
    false /
    null /
    undefined /
    number /
    string;

false -> "false";
null  -> "null";
undefined  -> "undefined";
true  -> "true";

number -> minus? int frac? exp?;
decimal_point -> ".";
digit1_9 -> [1-9];
e -> [eE];
exp -> e (minus / plus)? DIGIT+;
frac -> decimal_point DIGIT+;
int -> zero / (digit1_9 DIGIT*);
minus -> "-";
plus -> "+";
divider -> "/";
multi -> "*";
zero -> "0";

operation -> minus / plus / divider / multi;


property_read -> property_keyread / property_literal;
string -> quotation_mark char* quotation_mark;
property_literal -> [A-Za-z0-9_.]+;
property_keyread -> property_literal property_keyread_reader;
property_keyread_reader -> ("[" property_keyread_key "]")+;
property_keyread_key -> string / DIGIT+;

char ->
    unescaped /
    (escape ("\"" / "\\" / "/" / "b" / "f" / "n" / "r" / "t" / ("u" HEXDIG HEXDIG HEXDIG HEXDIG)));

escape -> "\\";
quotation_mark -> "\"";
unescaped -> [^\0-\x1F\x22\x5C];

DIGIT  -> [0-9];
HEXDIG -> [0-9a-fA-F];

ws -> [ \n\r]*;`
}
