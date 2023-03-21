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
boolean := func(cs []rune) (bool, error) {
    var s string
    for _, c := range cs {
        s += string(c)
    }

    if s == "true" {
        return true, nil
    } else if s == "false" {
        return false, nil
    } else {
        return false, fmt.Errorf("expected boolean found %s", s)
    }
}
trueKeywordParser := Map(Sequence(Char('t'), Char('r'), Char('u'), Char('e')), boolean)
falseKeywordParser := Map(Sequence(Char('f'), Char('a'), Char('l'), Char('s'), Char('e')), boolean)
booleanParser := OneOf(trueKeywordParser, falseKeywordParser)
```