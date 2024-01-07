package compiler

import (
	"io"

	"github.com/devansh42/rugialang/util"
)

type scanner struct {
	// current scanner position position
	curPos  int
	lastErr error
	// peekEOF is error occured while peeking
	peekEOF   error
	src       []byte
	tokens    []Token
	line, col int
}

func newScanner(reader io.Reader) (*scanner, error) {
	var sc = scanner{
		curPos: -1,
		line:   1,
	}
	sc.src, sc.lastErr = io.ReadAll(reader)
	if sc.lastErr != nil {
		return nil, sc.lastErr
	}
	return &sc, nil
}

func (sc *scanner) cur() byte {
	if !sc.err() {
		return sc.src[sc.curPos]
	}
	return 0
}

func (sc *scanner) peek() byte {
	if sc.curPos+1 < len(sc.src) {
		return sc.src[sc.curPos+1]
	}
	sc.peekEOF = io.EOF
	return 0
}

func (sc *scanner) nothingAhead() bool {
	return sc.peekEOF == io.EOF
}

func (sc *scanner) advance() byte {
	if sc.curPos+1 < len(sc.src) {
		sc.curPos++
		return sc.src[sc.curPos]
	}
	sc.lastErr = io.EOF
	return 0
}

func (sc *scanner) err() bool {
	return sc.lastErr != nil
}

func (sc *scanner) moveNextLine() {
	sc.line++
	sc.col = 0
}

func (sc *scanner) addToken(tokenType TokenType, line, col int) {
	sc.tokens = append(sc.tokens, Token{
		tokenType: tokenType,
		row:       line,
		col:       col,
	})
}
func (sc *scanner) addLiterToken(tokenType TokenType, line, col int, bufStartInc, bufEndEx int) {
	tok := Token{
		tokenType: tokenType,
		row:       line,
		col:       col,
	}
	tok.literal = make([]byte, bufEndEx-bufStartInc)
	copy(tok.literal, sc.src[bufStartInc:bufEndEx])
	sc.tokens = append(sc.tokens, tok)
}
func (sc *scanner) addIdentifierToken(col int, id []byte) {
	sc.tokens = append(sc.tokens, Token{
		tokenType: Identifier,
		literal:   id,
		row:       sc.line,
		col:       col,
	})
}

func (sc *scanner) consumeComment() {
	next := sc.peek()
	if !sc.nothingAhead() {
		if next == '/' { // Line Comment
			sc.consumeSingleLineComment()
		} else if next == '*' { // Multiline comment
			sc.consumeMultilineComment()
		} else {
			sc.col++
			sc.addToken(Divide, sc.line, sc.col)
		}
	}
}

func (sc *scanner) consumeSingleLineComment() {
	sc.advance() // consume /
	for x := sc.advance(); !sc.err() && x != '\n'; x = sc.advance() {
	}
	if !sc.err() { // found new line
		sc.moveNextLine()
	}
}

func (sc *scanner) consumeMultilineComment() {
	sc.advance() // consume *
	readCol := 2 // read (/*)
	for {
		for x := sc.advance(); !sc.err() && x != '\n'; x = sc.advance() {
			readCol++
			if x == '*' && sc.peek() == '/' {
				readCol++
				// end of comment
				sc.col = readCol
				return
			}
		}
		if !sc.err() { // new line occured
			sc.moveNextLine()
		} else { // some error occured
			return
		}
	}
}

func (sc *scanner) consumeInterpretedStr() {
	startPos := sc.curPos + 1
	l := 0
	for c := sc.advance(); !sc.err() && c != '"'; c = sc.advance() {
		l++
	}
	if !sc.err() {
		sc.addLiterToken(StringLit, sc.line, sc.col+1, startPos, startPos+l)
		sc.col += (l + 2) // 2 extra for 2 " symbol
	}
}

func (sc *scanner) consumeRawStr() {
	startPos := sc.curPos + 1
	l := 0
	sc.col++ // col after consuming opening `
	col := sc.col
	line := sc.line

	for c := sc.advance(); !sc.err() && c != '`'; c = sc.advance() {
		l++
		sc.col++
		if c == '\n' {
			sc.moveNextLine()
		}
	}
	if !sc.err() {
		sc.col++ // col after consuming closing `
		sc.addLiterToken(StringLit, line, col, startPos, startPos+l)
		sc.col += (l + 2) // 2 extra for 2 " symbol
	}
}

func (sc *scanner) consumeInt() {
	startPos := sc.curPos
	l := 1
	for c := sc.peek(); !sc.nothingAhead(); c = sc.peek() {
		if !util.IsDigit(c) {
			break
		}
		l++
		sc.advance() // consume digit
	}
	sc.addLiterToken(IntLit, sc.line, sc.col+l, startPos, startPos+l)
	sc.col += l
}
func (sc *scanner) consumeOperator() bool {
	var currentLiteral = []byte{sc.cur()}
	perfectMatch, prefixMatch := util.Match(OperatorTokenTrie, currentLiteral)
	if perfectMatch && !prefixMatch { // Found exact Operator
		sc.addToken(OperatorTokenMap[string([]byte{sc.cur()})], sc.line, sc.col)
		sc.col++
		return true
	} else if prefixMatch { // operator can be combined operator e.g. << >> <= >=
		nextToken := sc.peek()
		if !sc.nothingAhead() {
			currentLiteral = append(currentLiteral, nextToken)
			perfectMatch, _ = util.Match(OperatorTokenTrie, currentLiteral)
			if perfectMatch {
				sc.addToken(OperatorTokenMap[string(currentLiteral[:])], sc.line, sc.col)
				sc.col += 2
				sc.advance() // consuming second byte
				return true
			}
		}
		sc.addToken(OperatorTokenMap[string([]byte{sc.cur()})], sc.line, sc.col)
		sc.col++
		return true
	}
	return false

}

func (sc *scanner) consumeKeyword() []byte {
	var currentLiteral []byte

	for ; !sc.err(); sc.advance() {

		currentLiteral = append(currentLiteral, sc.cur())
		perfectMatch, prefixMatch := util.Match(KeywordTokenTrie, currentLiteral)

		// Note: We are considering that perfectMatch and prefixMatch are mutually exclusive
		// as we are doing in case of operators because so far no substring of keywords are also a keyword
		// which is not the case with operators
		// e.g. < is a substring of << and <=
		if perfectMatch {
			nextToken := sc.peek()
			if !sc.nothingAhead() && (util.IsAlpha(nextToken) || util.IsDigit(nextToken) || util.IsUnderscore(nextToken)) {
				return currentLiteral // it's an identifer
			}
			sc.addToken(KeywordTokenMap[string(currentLiteral)], sc.line, sc.col+len(currentLiteral))
			sc.col += len(currentLiteral)
			return nil
		} else if prefixMatch {
			continue
		} else {
			return currentLiteral
		}
	}
	return nil
}

func (sc *scanner) consumeIdentifier(curLiteral []byte) {
	for c := sc.peek(); !sc.nothingAhead() && (util.IsAlpha(c) || util.IsDigit(c) || util.IsUnderscore(c)); c = sc.peek() {
		curLiteral = append(curLiteral, c)
		sc.advance() // consume
	}
	sc.addIdentifierToken(sc.col, curLiteral)
	sc.col += len(curLiteral)
}

func (sc *scanner) consumeIdentifierOperatorOrKeyword() {
	foundOperator := sc.consumeOperator()
	if foundOperator {
		return
	}
	if !sc.err() && !sc.nothingAhead() {
		idLiteral := sc.consumeKeyword()
		if len(idLiteral) > 0 {
			sc.consumeIdentifier(idLiteral)
		}
	}

}

func (sc *scanner) scan() {
	for c := sc.advance(); !sc.err(); c = sc.advance() {
		switch c {

		case '\n':

			sc.moveNextLine()
			sc.addToken(EOL, sc.line, sc.col)

		case '\t', ' ', '\r', '\v', '\f':
			sc.col++ // TODO: Special arrangement for tabs
		case '/':

			sc.consumeComment()
			break
		case '"':
			sc.consumeInterpretedStr()
			break

		case '`':

			sc.consumeRawStr()
			break
		case ';':
			sc.addToken(SemiColon, sc.line, sc.col)
			sc.col++
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':

			sc.consumeInt()

		default:
			sc.consumeIdentifierOperatorOrKeyword()
		}
	}
}

// Algo:
// addSemicolons [] = []
// addSemicolons [TokenPos WhiteSpace _ _] = []
// addSemicolons [TokenPos Comment _ _] = []
// addSemicolons [x] = [x]
// addSemicolons (TokenPos WhiteSpace _ _:tail_) = addSemicolons tail_
// addSemicolons (TokenPos Comment _ _:tail_) = addSemicolons tail_
// addSemicolons (TokenPos EOL _ _:x:tail_) = addSemicolons (x:tail_)
// addSemicolons (TokenPos SemiColon _ _ :TokenPos SemiColon x y:tail_) = addSemicolons $ TokenPos SemiColon x y: tail_ -- merging semicolons incase user has already put an semicolon
// addSemicolons (x:TokenPos EOL sp ep  :tail_) = case getTk x of
//
//		Identifier _  -> fn
//		IntLit _ -> fn
//		RuneLit _ -> fn
//		FloatLit _ -> fn
//		StringLit _ -> fn
//		Break  -> fn
//		Continue -> fn
//		Return -> fn
//		RightBrace -> fn
//		RightSupScript -> fn
//		RightParen -> fn
//		_ -> x : addSemicolons tail_
//	where
//	fn = x : TokenPos SemiColon sp ep : addSemicolons tail_
//
// addSemicolons (x:tail_) = x : addSemicolons tail_
func addSemiColons(tokens []Token) []Token {
	if len(tokens) == 0 {
		return nil
	}
	if tokens[0].tokenType == EOL {
		return addSemiColons(tokens[1:])
	}

	if len(tokens) > 1 && tokens[0].tokenType == SemiColon && tokens[1].tokenType == SemiColon {
		return addSemiColons(tokens[1:])

	}
	if len(tokens) > 1 && tokens[1].tokenType == EOL {
		switch tokens[0].tokenType {
		case Identifier, IntLit, RuneLit, FlotLit, StringLit, Break, Continue, Return, CloseBrace, CloseBracket, CloseParen:
			tokens[1].tokenType = SemiColon
			return addSemiColons(tokens)
		default:
			val := append([]Token{}, tokens[0])
			return append(val, addSemiColons(tokens[2:])...)
		}
	}
	val := append([]Token{}, tokens[0])
	return append(val, addSemiColons(tokens[1:])...)
}

// Algo:
// addLastSemiColon [] = []
// addLastSemiColon [TokenPos SemiColon sp ep] = [TokenPos SemiColon sp ep]
// addLastSemiColon [TokenPos tk sp ep] = TokenPos tk sp ep : [ TokenPos SemiColon ep $ incSourceColumn ep 1 ]
// addLastSemiColon (x:tail_) = x: addLastSemiColon tail_

func addEOFSemiColon(tokens []Token) []Token {
	if len(tokens) == 0 {
		return nil
	}
	lastToken := tokens[len(tokens)-1]
	if lastToken.tokenType != SemiColon {
		tokens = append(tokens, Token{
			tokenType: SemiColon,
			row:       lastToken.row,
			col:       lastToken.col + 1,
		})
	}
	return tokens
}
