package parsec

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrEOF = errors.New("reached eof")

type Scanner interface {
	EOF() bool
	GetChar() rune
	Forward()
}

type Parser[OUT any] struct {
	Name string
	f    func(Scanner) (OUT, error)
}

func (p *Parser[OUT]) Parse(s Scanner) (OUT, error) {
	return p.f(s)
}

type StringScanner struct {
	s   string
	cur int
}

func NewStringScanner(s string) *StringScanner {
	return &StringScanner{s: s, cur: 0}
}

func (s *StringScanner) EOF() bool {
	return !(s.cur < len(s.s))
}

func (s *StringScanner) GetChar() rune {
	return rune(s.s[s.cur])
}

func (s *StringScanner) Forward() {
	s.cur++
}

func Char(c rune) Parser[rune] {
	return Parser[rune]{
		Name: fmt.Sprintf("char parser: '%c'", c),
		f: func(s Scanner) (rune, error) {
			if s.EOF() {
				return 0, ErrEOF
			}
			char := s.GetChar()
			if char != c {
				return 0, errors.New("char not matched")
			}
			s.Forward()
			return char, nil
		},
	}
}

func Sequence[IN, OUT any](mapper func([]IN) OUT, ps ...Parser[IN]) Parser[OUT] {
	return Parser[OUT]{
		Name: "sequence",
		f: func(s Scanner) (OUT, error) {
			o := new(OUT)

			var values []IN
			for _, p := range ps {
				value, err := p.Parse(s)
				if err != nil {
					return *o, err
				}
				values = append(values, value)
			}
			if mapper == nil {
				return *o, nil
			}
			return mapper(values), nil
		},
	}
}

func OneOf[T any](ps ...Parser[T]) Parser[T] {
	return Parser[T]{
		Name: "oneof",
		f: func(s Scanner) (T, error) {
			o := new(T)

			for _, p := range ps {
				value, err := p.Parse(s)
				if err != nil {
					continue
				}
				return value, nil
			}

			return *o, errors.New("no parser matched")
		},
	}
}

func OneOrMore[IN, OUT any](mapper func([]IN) OUT, parser Parser[IN]) Parser[OUT] {
	return Parser[OUT]{
		Name: "one or more",
		f: func(s Scanner) (OUT, error) {
			var values []IN
			var err error
			value, err := parser.Parse(s)
			if err != nil {
				return *new(OUT), errors.New("one or more should at least matched once")
			}
			values = append(values, value)
			for !s.EOF() {
				value, err := parser.Parse(s)
				if err != nil {
					break
				}
				values = append(values, value)
			}
			if mapper == nil {
				return *new(OUT), nil
			}
			return mapper(values), nil
		},
	}
}

func ZeroOrMore[IN, OUT any](mapper func([]IN) OUT, parser Parser[IN]) Parser[OUT] {
	return Parser[OUT]{
		Name: "zero or more",
		f: func(s Scanner) (OUT, error) {
			var values []IN
			for !s.EOF() {
				value, err := parser.Parse(s)
				if err != nil {
					break
				}
				values = append(values, value)
			}
			if mapper == nil {
				return *new(OUT), nil
			}
			return mapper(values), nil
		},
	}
}

func ZeroOrOne[IN, OUT any](mapper func(*IN) OUT, parser Parser[IN]) Parser[OUT] {
	return Parser[OUT]{
		Name: "zero or one",
		f: func(s Scanner) (OUT, error) {
			value, err := parser.Parse(s)
			if err != nil {
				if mapper == nil {
					return *new(OUT), nil
				}
				return mapper(nil), nil
			}
			if mapper == nil {
				return *new(OUT), nil
			}
			return mapper(&value), nil
		},
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

var IntParser = OneOrMore(func(cs []rune) int64 {
	var s string
	for _, c := range cs {
		s += string(c)
	}
	i, _ := strconv.Atoi(s)
	return int64(i)
}, DigitParser)
