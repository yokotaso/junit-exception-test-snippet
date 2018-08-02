package main

import (
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
	"log"
	"bufio"
	"strings"
)

var tokens = []string{
	"TEST_ANNOTATION", "L_PARENTHESES", "R_PARENTHESES", "EXPECTED_ATTRIBUTE",
}

var tokenMap map[string]int

func init() {
	tokenMap = make(map[string]int)
	for id, name := range tokens {
		tokenMap[name] = id
	}
}

func getToken(tokenType int) lexmachine.Action {
	return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
		return s.Token(tokenType, string(m.Bytes), m), nil
	}
}

func skip(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
	return nil, nil
}

func newLexer() *lexmachine.Lexer {
	lexer := lexmachine.NewLexer()
	lexer.Add([]byte(`@Test`), getToken(tokenMap["TEST_ANNOTATION"]))
	lexer.Add([]byte(`\(`), getToken(tokenMap["L_PARENTHESES"]))
	lexer.Add([]byte(`\)`), getToken(tokenMap["R_PARENTHESES"]))
	lexer.Add([]byte(`expected\s+=\s+[A-Z].*\.class`), getToken(tokenMap["EXPECTED_ATTRIBUTE"]))
	lexer.Add([]byte(`\s+`), skip)
	err := lexer.Compile()
	if err != nil {
		log.Fatal(err)
	}
	return lexer
}

type parsedToken struct {
	TokenValue  string
	StartColumn int
	EndColumn   int
}

type Snippet struct {
	indentEnd      int
	originalText   []byte
	exceptionClass string
}

func (s *Snippet) Write(writer *bufio.Writer) {
	indent := string(s.originalText[0:s.indentEnd])
	if _, err := writer.Write([]byte(indent + "// TODO(REMOVE ME) Original Annotation: " + string(s.originalText[s.indentEnd:])));
		err != nil {
		log.Fatal(err)
	}

	assertJSnippet := []byte(indent + "// TODO(REWRITE WITH ME) " + "assertThatThrownBy(() -> /** TODO */).isInstanceOf(" + s.exceptionClass + ");\n")
	if _, err := writer.Write(assertJSnippet);
		err != nil {
		log.Fatal(err)
	}

	if _, err := writer.Write([]byte(indent + "@Test\n"));
		err != nil {
		log.Fatal(err)
	}

}

func ParseTestAnnotation(text []byte) (bool, *Snippet) {
	lexer := newLexer()
	scanner, err := lexer.Scanner(text)
	if err != nil {
		log.Fatal(err)
	}
	parsedTokens := make(map[int]parsedToken)
	for token, err, eof := scanner.Next(); !eof; token, err, eof = scanner.Next() {
		if err != nil {
			log.Fatal(err)
		}

		token := token.(*lexmachine.Token)
		parsedTokens[token.Type] = parsedToken{
			StartColumn: token.StartColumn,
			EndColumn:   token.EndColumn,
			TokenValue:  token.Value.(string),
		}
	}

	testAnnotationToken, ok := parsedTokens[tokenMap["TEST_ANNOTATION"]];
	if !ok {
		return false, nil
	}

	expectedAttribute, ok := parsedTokens[tokenMap["EXPECTED_ATTRIBUTE"]]
	if ! ok {
		return false, nil
	}

	exceptionClass := ""
	for _, attribute := range strings.Fields(expectedAttribute.TokenValue) {
		if strings.HasSuffix(attribute, ".class") {
			exceptionClass = attribute
		}
	}
	return true, &Snippet{
		indentEnd:      testAnnotationToken.StartColumn - 1,
		originalText:   text,
		exceptionClass: exceptionClass,
	}
}
