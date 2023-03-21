package parsec

import (
	"errors"
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
	charKeywordParser := Map(Sequence(Char('c'), Char('h'), Char('a'), Char('r')), func(cs []rune) (struct{}, error) {
		var s string
		for _, c := range cs {
			s += string(c)
		}

		if s != "char" {
			return struct{}{}, errors.New("not matched")
		}

		return struct{}{}, nil
	})

	_, err := charKeywordParser.Parse(s)
	assert.NoError(t, err)
}

func TestOneOf(t *testing.T) {
	s := NewStringScanner("true")
	trueKeywordParser := Map(Sequence(Char('t'), Char('r'), Char('u'), Char('e')), func(_ []rune) (bool, error) { return true, nil })
	falseKeywordParser := Map(Sequence(Char('f'), Char('a'), Char('l'), Char('s'), Char('e')), func(_ []rune) (bool, error) { return false, nil })
	booleanParser := OneOf(trueKeywordParser, falseKeywordParser)

	b, err := booleanParser.Parse(s)
	assert.NoError(t, err)
	assert.True(t, b)
}

func TestOneOrMore(t *testing.T) {
	s := NewStringScanner("bbbb")
	bParser := Char('b')
	bsParser := OneOrMore(bParser)

	bs, err := bsParser.Parse(s)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(bs))
}

func TestZeroOrMore(t *testing.T) {
	s := NewStringScanner("bbbb")
	bParser := Char('b')
	bsParser := ZeroOrMore(bParser)

	bs, err := bsParser.Parse(s)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(bs))
}

func TestZeroOrOne(t *testing.T) {
	s := NewStringScanner("s")
	bParser := Char('b')
	bsParser := ZeroOrOne(bParser)

	b, err := bsParser.Parse(s)
	assert.NoError(t, err)
	assert.Nil(t, b)
}

func TestDigit(t *testing.T) {
	s := NewStringScanner("1")
	c, err := DigitParser.Parse(s)
	assert.NoError(t, err)
	assert.Equal(t, '1', c)

}

func TestInt(t *testing.T) {
	s := NewStringScanner("0123")
	c, err := IntParser.Parse(s)
	assert.NoError(t, err)
	assert.Equal(t, 123, c)
}

func TestFloat(t *testing.T) {
	s := NewStringScanner("0123.21")
	c, err := FloatParser.Parse(s)
	assert.NoError(t, err)
	assert.Equal(t, 123.21, c)
}
