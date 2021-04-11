package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	// ansin "atomc/src/analizator/sintactic"
)

type TokenType int

const (
	CtReal TokenType = iota
	CtInt
	CtChar
	CtString
	Id
	End
	Div
	Add
	Sub
	Mul
	Dot
	And
	Or
	Not
	NotEq
	Equal
	Assign
	Less
	LessEq
	Greater
	GreaterEq
	Comma
	Semicolon
	Lpar
	Rpar
	Lbracket
	Rbracket
	Lacc
	Racc
	Error
	Break
	Char
	Double
	Else
	For
	If
	Int
	Return
	Struct
	Void
	While
)

var constLookup = map[TokenType]string{
	CtReal:    "CtReal",
	CtInt:     "CtInt",
	CtChar:    "CtChar",
	CtString:  "CtString",
	Id:        "Id",
	End:       "End",
	Div:       "Div",
	Add:       "Add",
	Sub:       "Sub",
	Mul:       "Mul",
	Dot:       "Dot",
	And:       "And",
	Or:        "Or",
	Not:       "Not",
	NotEq:     "NotEq",
	Equal:     "Equal",
	Assign:    "Assign",
	Less:      "Less",
	LessEq:    "LessEq",
	Greater:   "Greater",
	GreaterEq: "GreaterEq",
	Comma:     "Comma",
	Semicolon: "Semicolon",
	Lpar:      "Lpar",
	Rpar:      "Rpar",
	Lbracket:  "Lbracket",
	Rbracket:  "Rbracket",
	Lacc:      "Lacc",
	Racc:      "Racc",
	Error:     "Error",
	Break:     "Break",
	Char:      "Char",
	Double:    "Double",
	Else:      "Else",
	For:       "For",
	If:        "If",
	Int:       "Int",
	Return:    "Return",
	Struct:    "Struct",
	Void:      "Void",
	While:     "While",
}

type Token struct {
	tokenType TokenType
	value     interface{}
	line      uint
}

// ---------------------- ANLEX --------------------------------------
func getNextToken(text *string, curPos *uint, currLine *uint) Token {

	var tokenStr string = ""
	var state uint = 0
	var tokenChar byte

	for {
		if int(*curPos) == len(*text) {
			return Token{
				tokenType: End,
				line:      *currLine,
			}
		}
		var c byte = (*text)[*curPos]
		*curPos += 1

		switch state {
		case 0:
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_') {
				state = 30
				tokenStr += string(c)
			} else if c == '\x00' {
				*curPos -= 1
				return Token{
					tokenType: End,
					line:      *currLine,
				}
			} else if c == '+' {
				return Token{
					tokenType: Add,
					line:      *currLine,
				}
			} else if c == '-' {
				return Token{
					tokenType: Sub,
					line:      *currLine,
				}
			} else if c == '*' {
				return Token{
					tokenType: Mul,
					line:      *currLine,
				}
			} else if c == '.' {
				return Token{
					tokenType: Dot,
					line:      *currLine,
				}
			} else if c == ',' {
				return Token{
					tokenType: Comma,
					line:      *currLine,
				}
			} else if c == ';' {
				return Token{
					tokenType: Semicolon,
					line:      *currLine,
				}
			} else if c == '(' {
				return Token{
					tokenType: Lpar,
					line:      *currLine,
				}
			} else if c == ')' {
				return Token{
					tokenType: Rpar,
					line:      *currLine,
				}
			} else if c == '[' {
				return Token{
					tokenType: Lbracket,
					line:      *currLine,
				}
			} else if c == ']' {
				return Token{
					tokenType: Rbracket,
					line:      *currLine,
				}
			} else if c == '{' {
				return Token{
					tokenType: Lacc,
					line:      *currLine,
				}
			} else if c == '}' {
				return Token{
					tokenType: Racc,
					line:      *currLine,
				}
			} else if c == '&' {
				state = 15
			} else if c == '|' {
				state = 16
			} else if c == '!' {
				state = 17
			} else if c == '=' {
				state = 18
			} else if c == '<' {
				state = 19
			} else if c == '>' {
				state = 20
			} else if c == ' ' || c == '\r' || c == '\n' || c == '\t' {
				state = 0
				if c == '\n' {
					*currLine += 1
				}
			} else if c == '/' {
				state = 12
			} else if c == '0' {
				state = 2
				tokenStr += string(c)
			} else if c >= '1' && c <= '9' {
				state = 1
				tokenStr += string(c)
			} else if c == '\'' {
				state = 21
			} else if c == '"' {
				state = 25
			} else {
				*curPos -= 1
				return Token{
					tokenType: Error,
					line:      *currLine,
				}
			}
			break
		case 1:
			if c >= '0' && c <= '9' {
				tokenStr += string(c)
			} else if c == '.' {
				state = 7
				tokenStr += string(c)
			} else if c == 'e' || c == 'E' {
				state = 9
				tokenStr += string(c)
			} else {
				*curPos -= 1
				int_nr, err := strconv.Atoi(tokenStr)
				if err != nil {
					int_nr = 0
				}
				return Token{
					tokenType: CtInt,
					value:     int_nr,
					line:      *currLine,
				}
			}
			break

		case 2:
			// hex
			if c == 'x' {
				state = 4
				tokenStr += string(c)
			} else if c >= '0' && c <= '7' {
				state = 3
				tokenStr += string(c)
			} else if c == '8' || c == '9' {
				state = 6
				tokenStr += string(c)
			} else if c == 'e' || c == 'E' {
				state = 9
				tokenStr += string(c)
			} else if c == '.' {
				state = 7
				tokenStr += string(c)
			} else {
				return Token{
					tokenType: CtInt,
					value:     0,
					line:      *currLine,
				}
			}

			break
		case 3:
			// Octal
			if c >= '0' && c <= '7' {
				tokenStr += string(c)
			} else if c == '8' || c == '9' {
				state = 6
				tokenStr += string(c)
			} else if c == 'e' || c == 'E' {
				state = 9
				tokenStr += string(c)
			} else if c == '.' {
				state = 7
				tokenStr += string(c)
			} else {
				*curPos -= 1
				int_nr, err := strconv.ParseInt(tokenStr, 0, 64)
				if err != nil {
					int_nr = 0
				}
				return Token{
					tokenType: CtInt,
					value:     int_nr,
					line:      *currLine,
				}
			}
			break
		case 4:
			if (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') || (c >= '0' && c <= '9') {
				state = 5
				tokenStr += string(c)
			} else {
				return Token{
					tokenType: Error,
					line:      *currLine,
				}
			}
			break
		case 5:
			// hex
			if (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') || (c >= '0' && c <= '9') {
				state = 5
				tokenStr += string(c)
			} else {
				*curPos -= 1
				int_nr, err := strconv.ParseInt(tokenStr, 0, 64)
				if err != nil {
					int_nr = 0
				}
				return Token{
					tokenType: CtInt,
					value:     int_nr,
					line:      *currLine,
				}
			}
			break
		case 6:
			if c >= '0' && c <= '9' {
				tokenStr += string(c)
			} else if c == '.' {
				state = 7
				tokenStr += string(c)
			} else if c == 'e' || c == 'E' {
				state = 9
				tokenStr += string(c)
			} else {
				return Token{
					tokenType: Error,
					line:      *currLine,
				}
			}
			break

		case 7:
			if c >= '0' && c <= '9' {
				state = 8
				tokenStr += string(c)
			} else {
				return Token{
					tokenType: Error,
					line:      *currLine,
				}
			}
			break
		case 8:
			if c >= '0' && c <= '9' {
				state = 8
				tokenStr += string(c)
			} else if c == 'e' || c == 'E' {
				state = 9
				tokenStr += string(c)
			} else {
				*curPos -= 1
				float_nr, err := strconv.ParseFloat(tokenStr, 64)
				if err != nil {
					float_nr = 0.0
				}
				return Token{
					tokenType: CtReal,
					value:     float_nr,
					line:      *currLine,
				}
			}
			break
		case 9:
			if c == '+' || c == '-' {
				state = 10
				tokenStr += string(c)
			} else if c >= '0' && c <= '9' {
				state = 11
				tokenStr += string(c)
			} else {
				return Token{
					tokenType: Error,
					line:      *currLine,
				}
			}

			break
		case 10:
			if c >= '0' && c <= '9' {
				state = 11
				tokenStr += string(c)
			} else {
				return Token{
					tokenType: Error,
					line:      *currLine,
				}
			}
			break
		case 11:
			if c >= '0' && c <= '9' {
				state = 11
				tokenStr += string(c)
			} else {
				*curPos -= 1
				float_nr, err := strconv.ParseFloat(tokenStr, 64)
				if err != nil {
					float_nr = 0.0
				}
				return Token{
					tokenType: CtReal,
					value:     float_nr,
					line:      *currLine,
				}
			}
			break
		case 12:
			if c == '*' {
				state = 13
			} else if c == '/' {
				state = 29
			} else {
				*curPos -= 1
				return Token{
					tokenType: Div,
					line:      *currLine,
				}
			}
			break
		case 13:
			if c == '*' {
				state = 14
			} else if c == '\n' {
				*currLine += 1
			} else {
			}
			break
		case 14:
			if c == '*' {
				state = 14
			} else if c == '/' {
				state = 0
			} else if c == '\n' {
				state = 13
				*currLine += 1
			} else {
				state = 13
			}
			break
		case 15:
			if c == '&' {
				return Token{
					tokenType: And,
					line:      *currLine,
				}
			} else {
				*curPos -= 1
				return Token{
					tokenType: Error,
					line:      *currLine,
				}
			}
			break
		case 16:
			if c == '|' {
				return Token{
					tokenType: Or,
					line:      *currLine,
				}
			} else {
				*curPos -= 1
				return Token{
					tokenType: Error,
					line:      *currLine,
				}
			}
			break
		case 17:
			if c == '=' {
				return Token{
					tokenType: NotEq,
					line:      *currLine,
				}
			} else {
				*curPos -= 1
				return Token{
					tokenType: Not,
					line:      *currLine,
				}
			}
			break
		case 18:
			if c == '=' {
				return Token{
					tokenType: Equal,
					line:      *currLine,
				}
			} else {
				*curPos -= 1
				return Token{
					tokenType: Assign,
					line:      *currLine,
				}
			}
			break
		case 19:
			if c == '=' {
				return Token{
					tokenType: LessEq,
					line:      *currLine,
				}
			} else {
				*curPos -= 1
				return Token{
					tokenType: Less,
					line:      *currLine,
				}
			}
			break
		case 20:
			if c == '=' {
				return Token{
					tokenType: GreaterEq,
					line:      *currLine,
				}
			} else {
				*curPos -= 1
				return Token{
					tokenType: Greater,
					line:      *currLine,
				}
			}
			break
		case 21:
			if c == '\\' {
				state = 22
			} else {
				state = 24
				tokenChar = c
			}
			break
		case 22:
			if c == 'a' {
				state = 23
				tokenChar = '\x07'
			} else if c == 'b' {
				state = 23
				tokenChar = '\x08'
			} else if c == 't' {
				state = 23
				tokenChar = '\x09'
			} else if c == 'n' {
				state = 23
				tokenChar = '\x0A'
			} else if c == 'v' {
				state = 23
				tokenChar = '\x0B'
			} else if c == 'f' {
				state = 23
				tokenChar = '\x0C'
			} else if c == 'r' {
				state = 23
				tokenChar = '\x0D'
			} else if c == '0' {
				state = 23
				tokenChar = '\x00'
			} else if c == '?' || c == '"' || c == '\'' || c == '\\' {
				state = 23
				tokenChar = c
			} else {
				*curPos -= 1
				return Token{
					tokenType: Error,
					line:      *currLine,
				}
			}
			break
		case 23:
			if c == '\'' {
				return Token{
					tokenType: CtChar,
					value:     tokenChar,
					line:      *currLine,
				}
			} else {
				*curPos -= 1
				return Token{
					tokenType: Error,
					line:      *currLine,
				}
			}
			break
		case 24:
			if c == '\'' {
				return Token{
					tokenType: CtChar,
					value:     tokenChar,
					line:      *currLine,
				}
			} else {
				*curPos -= 1
				return Token{
					tokenType: Error,
					line:      *currLine,
				}
			}
			break
		// CtString
		case 25:
			if c == '\\' {
				state = 26
			} else if c == '"' {
				return Token{
					tokenType: CtString,
					value:     "",
					line:      *currLine,
				}
			} else if c == '\n' {
				*currLine += 1
				return Token{
					tokenType: Error,
					line:      *currLine,
				}
			} else {
				state = 28
				tokenStr += string(c)
			}
			break
		case 26:
			if c == 'a' {
				state = 27
				tokenStr += string('\x07')
			} else if c == 'b' {
				state = 27
				tokenStr += string('\x08')
			} else if c == 't' {
				state = 27
				tokenStr += string('\x09')
			} else if c == 'n' {
				state = 27
				tokenStr += string('\x0A')
			} else if c == 'v' {
				state = 27
				tokenStr += string('\x0B')
			} else if c == 'f' {
				state = 27
				tokenStr += string('\x0C')
			} else if c == 'r' {
				state = 27
				tokenStr += string('\x0D')
			} else if c == '0' {
				state = 27
				tokenStr += string('\x00')
			} else if c == '?' || c == '"' || c == '\'' || c == '\\' {
				state = 27
				tokenStr += string(c)
			} else {
				*curPos -= 1
				return Token{
					tokenType: Error,
					line:      *currLine,
				}
			}
			break
		case 27:
			if c == '\\' {
				state = 26
			} else if c == '"' {
				return Token{
					tokenType: CtString,
					value:     tokenStr,
					line:      *currLine,
				}
			} else if c == '\n' {
				*currLine += 1
				return Token{
					tokenType: Error,
					line:      *currLine,
				}
			} else {
				tokenStr += string(c)
				state = 28
			}
			break
		case 28:
			if c == '\\' {
				state = 26
			} else if c == '"' {
				return Token{
					tokenType: CtString,
					value:     tokenStr,
					line:      *currLine,
				}
			} else if c == '\n' {
				*currLine += 1
				return Token{
					tokenType: Error,
					line:      *currLine,
				}
			} else {
				tokenStr += string(c)
			}

			break

		case 29:
			if c == '\n' || c == '\r' || c == '\x00' {
				state = 0
			} else {
			}
			break
		case 30:
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
				state = 30
				tokenStr += string(c)
			} else {
				*curPos -= 1
				switch tokenStr {
				case "break":
					return Token{
						tokenType: Break,
						line:      *currLine,
					}
				case "char":
					return Token{
						tokenType: Char,
						line:      *currLine,
					}

				case "double":
					return Token{
						tokenType: Double,
						line:      *currLine,
					}

				case "else":
					return Token{
						tokenType: Else,
						line:      *currLine,
					}

				case "for":
					return Token{
						tokenType: For,
						line:      *currLine,
					}

				case "if":
					return Token{
						tokenType: If,
						line:      *currLine,
					}

				case "int":
					return Token{
						tokenType: Int,
						line:      *currLine,
					}

				case "return":
					return Token{
						tokenType: Return,
						line:      *currLine,
					}

				case "struct":
					return Token{
						tokenType: Struct,
						line:      *currLine,
					}

				case "void":
					return Token{
						tokenType: Void,
						line:      *currLine,
					}

				case "while":
					return Token{
						tokenType: While,
						line:      *currLine,
					}

				default:
					return Token{
						tokenType: Id,
						value:     tokenStr,
						line:      *currLine,
					}
				}
			}
		// CtChar
		default:
			// invalid state
			return Token{
				tokenType: Error,
				line:      *currLine,
			}
		}
	}
}

func getTokens(text *string) []Token {
	var tokens []Token

	var currLine uint = 1
	var curPos uint = 0

outer:
	for {
		t := getNextToken(text, &curPos, &currLine)
		tokens = append(tokens, t)
		switch t.tokenType {
		case End:
			break outer
		default:
			continue
		}
	}

	return tokens
}

func printTokens(tokens []Token) {
	fmt.Printf("%-10s %-10s \t %-10s\n", "line", "token", "value")
	fmt.Printf("%s\n", strings.Repeat("-", 30))
	for _, token := range tokens {
		if token.value == nil {
			fmt.Printf("%-10d %-10s\n", token.line, constLookup[token.tokenType])
			continue
		}
		if token.tokenType == CtChar {
			if char, ok := token.value.(uint8); ok {
				fmt.Printf("%-10d %-10s\t %-10s\n", token.line, constLookup[token.tokenType], string(rune(int(char))))
			}
			continue
		}

		fmt.Printf("%-10d %-10s\t %-10v\n", token.line, constLookup[token.tokenType], token.value)
		// https://stackoverflow.com/questions/13094690/how-many-spaces-for-tab-character-t

	}
}

// ---------------------- ANSIN --------------------------------------

var tokens []Token

var currTokenId int = 0

func tokenErr(msg string) {
	fmt.Printf("error in line %d: %s\n", tokens[currTokenId].line, msg)
	os.Exit(1)
}

func consume(code TokenType) bool {
	if tokens[currTokenId].tokenType == code {
		currTokenId += 1
		return true
	}
	return false
}

func unit() bool {
	for {
		if declStruct() || declFunc() || declVar() {

		} else {
			break
		}
	}
	if consume(End) {
		return true
	}
	return false
}
func declStruct() bool {
	if consume(Struct) {
		if consume(Id) {
			if consume(Lacc) {
				for {
					if declVar() {

					} else {
						break
					}
				}
				if consume(Racc) {
					if consume(Semicolon) {
						return true
					} else {
						tokenErr("expected `;`")
					}
				} else {
					tokenErr("expected `}`")
				}
			} else {
				tokenErr("expected `{`")
			}
		} else {
			tokenErr("expected identifier")
		}
	}
	return false
}
func declVar() bool {
	if typeBase() {
		if consume(Id) {
			arrayDecl()
			for {
				if consume(Comma) {
					if consume(Id) {
						arrayDecl()
					} else {
						tokenErr("expected identifier")
					}
				} else {
					break
				}
			}
			if consume(Semicolon) {
				return true
			} else {
				tokenErr("expected `;`")
			}
		} else {
			tokenErr("expected identifier")
		}
	}

	return false
}
func typeBase() bool {

	startId := currTokenId

	if consume(Int) || consume(Double) || consume(Char) {
		return true
	}
	if consume(Struct) {
		if consume(Id) {
			return true
		} else {
			tokenErr("expected identifier")
		}
	} else {
		currTokenId = startId
	}

	return false
}
func arrayDecl() bool {
	if consume(Lbracket) {
		expr()
		if consume(Rbracket) {
			return true
		} else {
			tokenErr("expected `]`")
		}
	}
	return false
}
func typeName() bool {
	if typeBase() {
		arrayDecl()
		return true
	}
	return false
}
func declFunc() bool {
	return true
}
func funcArg() bool {
	return true
}
func stm() bool {
	return true
}
func stmCompound() bool {
	return true
}

func expr() bool {
	return true
}
func exprAssign() bool {
	return true
}
func exprOr() bool {
	return true
}
func exprAnd() bool {
	return true
}
func exprEq() bool {
	return true
}
func exprRel() bool {
	return true
}
func exprAdd() bool {
	return true
}
func exprMul() bool {
	return true
}
func exprCast() bool {
	return true
}
func exprUnary() bool {
	return true
}
func exprPostfix() bool {
	return true
}
func exprPrimary() bool {
	return true
}

func ansin() {
	if unit() {
	} else {
		tokenErr("top level error")
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("usage: %s [options] file\n", os.Args[0])
		os.Exit(1)
	}

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)

	tokens = getTokens(&text)
	printTokens(tokens)

}
