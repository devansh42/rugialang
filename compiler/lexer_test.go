package compiler

import (
	"strings"
	"testing"
)

var codeSnippet = `
package main

func main(){
	fmt.Println("Devansh Gupta")
	var x = 1+(2*3) // This a line comment

}
`
var c2 = `name=1>>>`

func TestLexer(t *testing.T) {
	sc, err := newScanner(strings.NewReader(codeSnippet))
	if err != nil {
		t.Errorf("err while reading input string: %v", err)
		return
	}
	sc.scan()
	t.Logf("Tokens:\n%v ", addEOFSemiColon(addSemiColons(sc.tokens)))

}

func TestTrie(t *testing.T) {
	t.Logf("%+v", OperatorTokenTrie)
}
