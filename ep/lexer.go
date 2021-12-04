package ep

import (
	"errors"
	"github.com/irustm/ng-template-parser/chars"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

// https://github.com/angular/angular/blob/master/packages/compiler/src/expression_parser/lexer.ts

type TokenType int

const (
	Character TokenType = iota
	Identifier
	PrivateIdentifier
	Keyword
	String
	Operator
	Number
	Error
)

var KEYWORDS = []string{"var", "let", "as", "null", "undefined", "true", "false", "if", "else", "this"}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type Lexer struct {
}

func (l Lexer) Tokenize(text string) []Token {
	var scanner = newScanner(text)
	var tokens []Token

	var token, err = scanner.scanToken()

	for err == nil {
		tokens = append(tokens, token)
		token, err = scanner.scanToken()

	}
	return tokens
}

type Token struct {
	Index     int
	End       int
	TypeToken TokenType
	NumValue  int
	StrValue  string
}

func (t Token) isCharacter(code int) bool {
	return t.TypeToken == Character && t.NumValue == code
}

func (t Token) isNumber() bool {
	return t.TypeToken == Number
}

func (t Token) isString() bool {
	return t.TypeToken == String
}

func (t Token) isOperator(operator string) bool {
	return t.TypeToken == Operator && t.StrValue == operator
}

func (t Token) isIdentifier() bool {
	return t.TypeToken == Identifier
}

func (t Token) isPrivateIdentifier() bool {
	return t.TypeToken == PrivateIdentifier
}

func (t Token) isKeyword() bool {
	return t.TypeToken == Keyword
}

func (t Token) isKeywordLet() bool {
	return t.TypeToken == Keyword && t.StrValue == "let"
}

func (t Token) isKeywordAs() bool {
	return t.TypeToken == Keyword && t.StrValue == "as"
}

func (t Token) isKeywordNull() bool {
	return t.TypeToken == Keyword && t.StrValue == "null"
}

func (t Token) isKeywordUndefined() bool {
	return t.TypeToken == Keyword && t.StrValue == "undefined"
}

func (t Token) isKeywordTrue() bool {
	return t.TypeToken == Keyword && t.StrValue == "true"
}

func (t Token) isKeywordFalse() bool {
	return t.TypeToken == Keyword && t.StrValue == "false"
}

func (t Token) isKeywordThis() bool {
	return t.TypeToken == Keyword && t.StrValue == "this"
}

func (t Token) isError() bool {
	return t.TypeToken == Error
}

func (t Token) toNumber() int {
	if t.TypeToken == Number {
		return t.NumValue
	}

	return -1
}

func (t Token) toString() string {
	switch t.TypeToken {
	case Character:
	case Identifier:
	case Keyword:
	case Operator:
	case PrivateIdentifier:
	case String:
	case Error:
		{
			return t.StrValue
		}
	case Number:
		{
			return string(t.NumValue)
		}
	default:
		{
			return ""
		}
	}
	return ""
}

func newCharacterToken(index int, end int, code int) Token {
	return Token{
		index,
		end,
		Character,
		code,
		string(code),
	}
}

func newIdentifierToken(index int, end int, text string) Token {
	return Token{index, end, Identifier, 0, text}
}

func newPrivateIdentifierToken(index int, end int, text string) Token {
	return Token{index, end, PrivateIdentifier, 0, text}
}

func newKeywordToken(index int, end int, text string) Token {
	return Token{index, end, Keyword, 0, text}
}

func newOperatorToken(index int, end int, text string) Token {
	return Token{index, end, Operator, 0, text}
}

func newStringToken(index int, end int, text string) Token {
	return Token{index, end, String, 0, text}
}

func newNumberToken(index int, end int, n int) Token {
	return Token{index, end, Number, n, ""}
}

func newErrorToken(index int, end int, message string) Token {
	return Token{index, end, Error, 0, message}
}

var EOF Token = Token{-1, -1, Character, 0, ""}

type scanner struct {
	length int
	peek   int
	index  int
	input  string
}

func newScanner(input string) scanner {
	res := scanner{length: len(input), peek: 0, index: -1, input: input}

	return res
}

func (s *scanner) advance() {
	s.index++

	if s.index >= s.length {
		s.peek = chars.VEOF
	} else {
		s.peek = int(s.input[s.index])
	}
}

func (s *scanner) scanToken() (Token, error) {
	var input = s.input
	var length = s.length
	var peek = s.peek
	var index = s.index

	for peek <= chars.VSPACE {

		index++

		if index >= length {
			peek = chars.VEOF
			break
		} else {
			peek = int(input[index])
		}
	}

	s.peek = peek
	s.index = index

	if index >= length {
		return Token{}, errors.New("null")
	}

	// Handle identifiers and numbers.
	if isIdentifierStart(peek) {
		return s.scanIdentifier()
	}
	if chars.IsDigit(peek) {
		return s.scanNumber(index)
	}

	var start = index

	switch peek {
	case chars.VPERIOD:
		{
			s.advance()
			if chars.IsDigit(s.peek) {
				return s.scanNumber(start)
			} else {
				return newCharacterToken(start, s.index, chars.VPERIOD), nil
			}
		}
	case chars.VLPAREN, chars.VRPAREN, chars.VLBRACE, chars.VRBRACE, chars.VLBRACKET, chars.VRBRACKET, chars.VCOMMA, chars.VCOLON, chars.VSEMICOLON:
		{
			return s.scanCharacter(start, peek)
		}
	case chars.VSQ, chars.VDQ:
		{
			return s.scanString()
		}
	case chars.VHASH:
		{
			return s.scanPrivateIdentifier()
		}

	case chars.VPLUS, chars.VMINUS, chars.VSTAR, chars.VSLASH, chars.VPERCENT, chars.VCARET:
		{
			return s.scanOperator(start, string(peek))
		}

	case chars.VQUESTION:
		{
			return s.scanQuestion(start)
		}
	case chars.VLT, chars.VGT:
		{
			return s.scanComplexOperator(start, string(peek), chars.VEQ, "=")

		}
	case chars.VBANG, chars.VEQ:
		{
			return s.scanComplexOperatorThree(start, string(peek), chars.VEQ, "=", chars.VEQ, "=")
		}
	case chars.VAMPERSAND:
		{
			return s.scanComplexOperator(start, "&", chars.VAMPERSAND, "&")
		}
	case chars.VBAR:
		{
			return s.scanComplexOperator(start, "|", chars.VBAR, "|")

		}
	case chars.VNBSP:
		{

			for chars.IsWhitespace(s.peek) {
				s.advance()
			}

			return s.scanToken()
		}
	}

	s.advance()

	return Token{}, errors.New("Unexpected character [" + string(peek) + "]")
}

func (s *scanner) scanCharacter(start int, code int) (Token, error) {
	s.advance()
	return newCharacterToken(start, s.index, code), nil
}

func (s *scanner) scanOperator(start int, str string) (Token, error) {
	s.advance()
	return newOperatorToken(start, s.index, str), nil
}

func (s *scanner) scanComplexOperator(start int, one string, twoCode int, two string) (Token, error) {
	s.advance()

	var str string = one

	if s.peek == twoCode {
		s.advance()
		str += two
	}

	return newOperatorToken(start, s.index, str), nil
}

func (s *scanner) scanComplexOperatorThree(start int, one string, twoCode int, two string, threeCode int, three string) (Token, error) {
	s.advance()

	var str string = one

	if s.peek == twoCode {
		s.advance()
		str += two
	}

	if s.peek == threeCode {
		s.advance()
		str += three
	}

	return newOperatorToken(start, s.index, str), nil
}

func (s *scanner) scanIdentifier() (Token, error) {
	var start = s.index
	s.advance()
	for isIdentifierPart(s.peek) {
		s.advance()
	}
	var str = s.input[start:s.index]

	if contains(KEYWORDS, str) {
		return newKeywordToken(start, s.index, str), nil
	} else {
		return newIdentifierToken(start, s.index, str), nil
	}
}

/** Scans an ECMAScript private identifier. */
func (s *scanner) scanPrivateIdentifier() (Token, error) {
	var start = s.index
	s.advance()
	for isIdentifierStart(s.peek) {
		return Token{}, errors.New("Invalid character [#]")
	}

	for isIdentifierPart(s.peek) {
		s.advance()
	}

	var identifierName = s.input[start:s.index]
	return newPrivateIdentifierToken(start, s.index, identifierName), nil
}

func (s *scanner) scanNumber(start int) (Token, error) {
	var simple = s.index == start
	var hasSeparators = false
	s.advance()

	for {
		if chars.IsDigit(s.peek) {
			// Do nothing.
		} else if s.peek == chars.V_ {
			// Separators are only valid when they're surrounded by digits. E.g. `1_0_1` is
			// valid while `_101` and `101_` are not. The separator can't be next to the decimal
			// point or another separator either. Note that it's unlikely that we'll hit a case where
			// the underscore is at the start, because that's a valid identifier and it will be picked
			// up earlier in the parsing. We validate for it anyway just in case.
			if !chars.IsDigit(int(s.input[s.index-1])) || !chars.IsDigit(int(s.input[s.index+1])) {
				return Token{}, errors.New("Invalid numeric separator")
			}
			hasSeparators = true
		} else if s.peek == chars.VPERIOD {
			simple = false
		} else if isExponentStart(s.peek) {
			s.advance()

			if isExponentSign(s.peek) {
				s.advance()
			}
			if !chars.IsDigit(s.peek) {
				return Token{}, errors.New("Invalid exponent")
			}
			simple = false
		} else {
			break
		}

		s.advance()
	}

	var str = s.input[start:s.index]
	if hasSeparators {
		str = strings.ReplaceAll(str, "_", "")
	}

	var value int
	if simple {
		value = parseIntAutoRadix(str)
	} else {
		value = parseFloat(str)
	}

	return newNumberToken(start, s.index, value), nil
}

func (s *scanner) scanString() (Token, error) {
	var start = s.index
	var quote = s.peek
	s.advance()

	var buffer string = ""
	var marker int = s.index
	var input string = s.input

	for s.peek != quote {
		if s.peek == chars.VBACKSLASH {
			buffer += input[marker:s.index]
			s.advance()
			var unescapedCode int
			// Workaround for TS2.1-introduced type strictness
			s.peek = s.peek // ??
			if s.peek == chars.Vu {
				var hex = input[s.index+1 : s.index+5]
				matcher, _ := regexp.Compile("^[0-9a-f]+$")

				if matcher.MatchString(hex) {
					n := new(big.Int)
					n.SetString(hex, 16)
					unescapedCode = int(n.Uint64())
				} else {
					return Token{}, errors.New("Invalid unicode escape " + hex)
				}

				for i := 0; i < 5; i++ {
					s.advance()
				}
			} else {
				unescapedCode = unescape(s.peek)
				s.advance()
			}

			buffer += string(unescapedCode)
			marker = s.index
		} else if s.peek == chars.VEOF {
			return Token{}, errors.New("Unterminated quote")
		} else {
			s.advance()
		}
	}

	var last = input[marker:s.index]
	s.advance()

	return newStringToken(start, s.index, buffer+last), nil
}

func (s *scanner) scanQuestion(start int) (Token, error) {
	s.advance()
	var str = "?"
	// Either `a ?? b` or 'a?.b'.
	if s.peek == chars.VQUESTION || s.peek == chars.VPERIOD {
		if s.peek == chars.VPERIOD {
			str += "."
		} else {
			str += "?"
		}
		s.advance()
	}
	return newOperatorToken(start, s.index, str), nil
}

func (s *scanner) error(message string, offset int) (Token, error) {
	position := s.index + offset

	return newErrorToken(position, s.index, `Lexer Error: `+message+` at column `+string(position)+` in expression [`+s.input+`]`), nil
}

func isIdentifierStart(code int) bool {
	return (chars.Va <= code && code <= chars.Vz) || (chars.VA <= code && code <= chars.VZ) ||
		(code == chars.V_) || (code == chars.VDOLLAR)
}

func IsIdentifier(input string) bool {
	if len(input) == 0 {
		return false
	}

	var scanner = newScanner(input)

	if isIdentifierStart(scanner.peek) {
		return false
	}
	scanner.advance()

	for scanner.peek != chars.VEOF {
		if isIdentifierPart(scanner.peek) {
			return false
		}
		scanner.advance()
	}

	return true
}

func isIdentifierPart(code int) bool {
	return chars.IsAsciiLetter(code) || chars.IsDigit(code) || (code == chars.V_) || (code == chars.VDOLLAR)
}

func isExponentStart(code int) bool {
	return (code == chars.Ve) || (code == chars.VE)
}

func isExponentSign(code int) bool {
	return (code == chars.VMINUS) || (code == chars.VPLUS)
}

func unescape(code int) int {
	switch code {
	case chars.Vn:
		return chars.VLF
	case chars.Vf:
		return chars.VFF
	case chars.Vr:
		return chars.VCR
	case chars.Vt:
		return chars.VTAB
	case chars.Vv:
		return chars.VVTAB
	default:
		return code
	}
}

func parseIntAutoRadix(text string) int {
	result, err := strconv.Atoi(text)

	if err != nil {
		// error
	}

	return result
}

func parseFloat(text string) int {
	result, err := strconv.ParseFloat(text, 10)

	if err != nil {
		// error
	}

	// TODO,
	return int(result)
}
