package nom

import (
	"errors"
	"strconv"
)

type Empty struct{}

var ErrEOF = errors.New("reached eof")

type Parser[OUT any] func(Scanner) (OUT, error)

func Char(c rune) Parser[rune] {
	return func(s Scanner) (rune, error) {
		if s.EOF() {
			return 0, ErrEOF
		}
		char := s.GetChar()
		if char != c {
			return 0, errors.New("char not matched")
		}
		s.Forward()
		return char, nil
	}
}

func Sequence[T any](ps ...Parser[T]) Parser[[]T] {
	return func(s Scanner) ([]T, error) {
		var values []T
		for _, p := range ps {
			value, err := p(s)
			if err != nil {
				return nil, err
			}
			values = append(values, value)
		}
		return values, nil
	}
}

func OneOf[T any](ps ...Parser[T]) Parser[T] {
	return func(s Scanner) (T, error) {
		o := new(T)

		for _, p := range ps {
			value, err := p(s)
			if err != nil {
				continue
			}
			return value, nil
		}

		return *o, errors.New("no parser matched")
	}
}

func OneOrMore[T any](parser Parser[T]) Parser[[]T] {
	return func(s Scanner) ([]T, error) {
		var values []T
		var err error
		value, err := parser(s)
		if err != nil {
			return nil, errors.New("one or more should at least matched once")
		}
		values = append(values, value)
		for !s.EOF() {
			value, err := parser(s)
			if err != nil {
				break
			}
			values = append(values, value)
		}
		return values, nil
	}
}

func ZeroOrMore[T any](parser Parser[T]) Parser[[]T] {
	return func(s Scanner) ([]T, error) {
		var values []T
		for !s.EOF() {
			value, err := parser(s)
			if err != nil {
				break
			}
			values = append(values, value)
		}
		return values, nil
	}
}

func ZeroOrOne[T any](parser Parser[T]) Parser[*T] {
	return func(s Scanner) (*T, error) {
		value, err := parser(s)
		if err != nil {
			return nil, nil
		}
		return &value, nil
	}
}

func Map[IN, OUT any](p Parser[IN], f func(i IN) (OUT, error)) Parser[OUT] {
	return func(s Scanner) (OUT, error) {
		res, err := p(s)
		if err != nil {
			return *new(OUT), err
		}

		return f(res)
	}
}

func Flatten[T any](ps ...Parser[[]T]) Parser[[]T] {
	return func(s Scanner) ([]T, error) {
		var output []T
		for _, p := range ps {
			values, err := p(s)
			if err != nil {
				return nil, err
			}
			output = append(output, values...)
		}

		return output, nil
	}
}

var DigitParser = OneOf(
	Char('0'),
	Char('1'),
	Char('2'),
	Char('3'),
	Char('4'),
	Char('5'),
	Char('6'),
	Char('7'),
	Char('8'),
	Char('9'),
)

var IntParser = Map(OneOrMore(DigitParser), func(cs []rune) (int, error) {
	i, err := strconv.Atoi(string(cs))
	if err != nil {
		return 0, err
	}
	return i, nil
})

var FloatParser = Map(Map(Sequence(
	OneOrMore(DigitParser),
	Map(Sequence(Map(Char('.'), func(r rune) ([]rune, error) { return []rune{r}, nil }), OneOrMore(DigitParser)),
		func(rss [][]rune) ([]rune, error) {
			var out []rune
			for _, rs := range rss {
				out = append(out, rs...)
			}
			return out, nil
		})), func(rss [][]rune) ([]rune, error) {
	var out []rune
	for _, rs := range rss {
		out = append(out, rs...)
	}
	return out, nil
},
), func(rs []rune) (float64, error) {
	i, err := strconv.ParseFloat(string(rs), 64)
	if err != nil {
		return 0, err
	}
	return i, nil
})
