# go-parsec

Parser combinator library for golang.


## Example

### Integer Parser

```go
var IntParser = Map(OneOrMore(DigitParser), func(cs []rune) (int, error) {
	i, err := strconv.Atoi(string(cs))
	if err != nil {
		return 0, err
	}
	return i, nil
})
```

### Boolean Parser

```go
s := NewStringScanner("true")
trueKeywordParser := Map(Sequence(Char('t'), Char('r'), Char('u'), Char('e')), func(_ []rune) (bool, error) { return true, nil })
falseKeywordParser := Map(Sequence(Char('f'), Char('a'), Char('l'), Char('s'), Char('e')), func(_ []rune) (bool, error) { return false, nil })
booleanParser := OneOf(trueKeywordParser, falseKeywordParser)
```
