package parsec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChar(t *testing.T) {
	s := NewStringScanner("c")
	parser := Char('c')
	c, err := parser.Parse(s)

	assert.NoError(t, err)

	assert.Equal(t, 'c', c)

	assert.Equal(t, 1, s.cur)
	assert.Equal(t, true, s.EOF())
}

func TestSeq(t *testing.T) {
	s := NewStringScanner("char")
	charKeywordParser := Sequence[rune, bool](nil, Char('c'), Char('h'), Char('a'), Char('r'))

	_, err := charKeywordParser.Parse(s)
	assert.NoError(t, err)
}

func TestOneOf(t *testing.T) {
	s := NewStringScanner("true")
	trueKeywordParser := Sequence(func(i []rune) bool { return true }, Char('t'), Char('r'), Char('u'), Char('e'))
	falseKeywordParser := Sequence(func(i []rune) bool { return false }, Char('f'), Char('a'), Char('l'), Char('s'), Char('e'))
	booleanParser := OneOf(trueKeywordParser, falseKeywordParser)

	b, err := booleanParser.Parse(s)
	assert.NoError(t, err)
	assert.True(t, b)
}

func TestOneOrMore(t *testing.T) {
	s := NewStringScanner("bbbb")
	bParser := Char('b')
	bsParser := OneOrMore(func(i []rune) (int, error) {
		return len(i), nil
	}, bParser)

	b, err := bsParser.Parse(s)
	assert.NoError(t, err)
	assert.Equal(t, 4, b)
}

func TestZeroOrMore(t *testing.T) {
	s := NewStringScanner("bbbb")
	bParser := Char('b')
	bsParser := ZeroOrMore(func(i []rune) int {
		return len(i)
	}, bParser)

	b, err := bsParser.Parse(s)
	assert.NoError(t, err)
	assert.Equal(t, 4, b)
}

func TestZeroOrOne(t *testing.T) {
	s := NewStringScanner("bbbb")
	bParser := Char('b')
	bsParser := ZeroOrOne(func(i *rune) bool {
		return i != nil
	}, bParser)

	b, err := bsParser.Parse(s)
	assert.NoError(t, err)
	assert.True(t, b)
}

func TestDigit(t *testing.T) {
	s := NewStringScanner("1")
	c, err := DigitParser.Parse(s)
	assert.NoError(t, err)
	assert.Equal(t, 1, c)

}

func TestInt(t *testing.T) {
	s := NewStringScanner("0123")
	c, err := IntParser.Parse(s)
	assert.NoError(t, err)
	assert.Equal(t, 123, c)
}

func TestFloat(t *testing.T) {
	s := NewStringScanner("0123.21")
	c, err := IntParser.Parse(s)
	assert.NoError(t, err)
	assert.Equal(t, 123, c)
}
