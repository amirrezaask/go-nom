package nom

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChar(t *testing.T) {
	parser := Char('c')
	tail, c, err := parser("c")

	assert.NoError(t, err)
	assert.Equal(t, 'c', c)
	assert.Empty(t, tail)
}
func TestTag(t *testing.T) {
	charParser := Tag("char")
	tail, out, err := charParser("char")
	assert.NoError(t, err)
	assert.Equal(t, "char", out)
	assert.Empty(t, tail)
}

func TestSeq(t *testing.T) {
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

	tail, _, err := charKeywordParser("char")
	assert.NoError(t, err)
	assert.Empty(t, tail)
}

func TestOneOf(t *testing.T) {
	trueKeywordParser := Map(Sequence(Char('t'), Char('r'), Char('u'), Char('e')), func(_ []rune) (bool, error) { return true, nil })
	falseKeywordParser := Map(Sequence(Char('f'), Char('a'), Char('l'), Char('s'), Char('e')), func(_ []rune) (bool, error) { return false, nil })
	booleanParser := OneOf(trueKeywordParser, falseKeywordParser)

	tail, b, err := booleanParser("true")
	assert.NoError(t, err)
	assert.True(t, b)
	assert.Empty(t, tail)
}

func TestOneOf2(t *testing.T) {
	trueKeyword := Tag("true")
	truerKeyword := Tag("truer")
	booleanParser := OneOf(truerKeyword, trueKeyword)

	tail, b, err := booleanParser("truer")
	assert.NoError(t, err)
	assert.Equal(t, "truer", b)
	assert.Empty(t, tail)
}

func TestOneOrMore(t *testing.T) {
	bParser := Char('b')
	bsParser := OneOrMore(bParser)

	tail, bs, err := bsParser("bbbb")
	assert.NoError(t, err)
	assert.Equal(t, 4, len(bs))
	assert.Empty(t, tail)
}

func TestZeroOrMore(t *testing.T) {
	bParser := Char('b')
	bsParser := ZeroOrMore(bParser)

	tail, bs, err := bsParser("bbbb")
	assert.NoError(t, err)
	assert.Equal(t, 4, len(bs))
	assert.Empty(t, tail)
}

func TestZeroOrOne(t *testing.T) {
	bParser := Char('b')
	bsParser := ZeroOrOne(bParser)

	tail, b, err := bsParser("s")
	assert.NoError(t, err)
	assert.Nil(t, b)
	assert.Equal(t, "s", tail)
}

func TestDigit(t *testing.T) {
	tail, c, err := DigitParser("1")
	assert.NoError(t, err)
	assert.Equal(t, '1', c)
	assert.Empty(t, tail)

}

func TestInt(t *testing.T) {
	tail, c, err := IntParser("0123")
	assert.NoError(t, err)
	assert.Equal(t, 123, c)
	assert.Empty(t, tail)
}

func TestFloat(t *testing.T) {
	tail, c, err := FloatParser("0123.21")
	assert.NoError(t, err)
	assert.Equal(t, 123.21, c)
	assert.Empty(t, tail)
}
