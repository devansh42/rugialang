package compiler

import (
	"fmt"

	"github.com/devansh42/rugialang/util"
)

type TokenType int

type Token struct {
	tokenType TokenType
	// literal represnts optional literal values for literal types
	literal []byte
	row     int
	col     int
}

const (
	// Arithmatic
	Plus TokenType = iota + 1
	Minus
	Mul
	Divide
	Mod
	// Assignment
	Assign
	ShortAssign
	// Comparison
	Eq
	Not
	NotEq
	LTH
	GTH
	LTEQ
	GTEQ
	// Terminator
	SemiColon
	WhiteSpace
	EOL
	// Brackets
	OpenParen
	CloseParen
	OpenBrace
	CloseBrace
	OpenBracket
	CloseBracket
	// Puntuation
	Comma
	Colon
	Period
	// Bitwise Ops
	BitwiseAND
	BitwiseOR
	BitwiseXOR
	LeftShift
	RightShift
	// BooleanOperator
	AND
	OR

	// Keywords

	Break
	Case
	Const
	Continue
	Default
	Else
	For
	Func
	Goto
	If
	Import
	Map
	Package
	Range
	Return
	Struct
	Switch
	Type
	Var

	//	Literal

	RuneLit
	StringLit
	IntLit
	FlotLit
	// Identifier
	Identifier
)

var OperatorTokenTrie = buildTrie(OperatorTokenMap)
var KeywordTokenTrie = buildTrie(KeywordTokenMap)

var OperatorTokenMap = map[string]TokenType{
	".":  Period,
	",":  Comma,
	":=": ShortAssign,
	":":  Colon,
	"+":  Plus,
	"-":  Minus,
	"*":  Mul,
	"%":  Mod,
	";":  SemiColon,
	"/":  Divide,
	"(":  OpenParen,
	")":  CloseParen,
	"[":  OpenBracket,
	"]":  CloseBracket,
	"{":  OpenBrace,
	"}":  CloseBrace,
	"<<": LeftShift,
	"<=": LTEQ,
	"<":  LTH,
	">>": RightShift,
	">=": GTEQ,
	">":  GTH,
	"==": Eq,
	"=":  Assign,
	"!=": NotEq,
	"!":  Not,
	"&&": AND,
	"||": OR,
	"&":  BitwiseAND,
	"|":  BitwiseOR,
	"^":  BitwiseXOR}
var KeywordTokenMap = map[string]TokenType{
	"break":    Break,
	"case":     Case,
	"const":    Const,
	"continue": Continue,
	"default":  Default,
	"else":     Else,
	"for":      For,
	"goto":     Goto,
	"if":       If,
	"import":   Import,
	"map":      Map,
	"package":  Package,
	"range":    Range,
	"return":   Return,
	"struct":   Struct,
	"switch":   Switch,
	"type":     Type,
	"var":      Var,
	"_":        Identifier,
}

func buildTrie(m map[string]TokenType) *util.Trie {
	trie := util.NewTrie()
	for k := range m {
		trie.Add([]byte(k))
	}
	return trie
}

func (tk Token) String() string {
	return fmt.Sprintf("(%d,%v, L%d, C%d)", tk.tokenType, tk.literal, tk.row, tk.col)
}
