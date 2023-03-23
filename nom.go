package nom

import (
	"errors"
	"fmt"
)

type Empty struct{}

var ErrEOF = errors.New("reached eof")

func eof(s string) bool {
	return s == ""
}

type Parser[OUT any] func(string) (string, OUT, error)

func Char(c rune) Parser[rune] {
	return func(s string) (string, rune, error) {
		// fmt.Println(s)
		// fmt.Println(c)
		if s == "" {
			return "", 0, ErrEOF
		}
		char := s[0]
		if rune(char) != c {
			return s, 0, fmt.Errorf("char not matched exp:'%c' got:'%c'", c, char)
		}
		return s[1:], c, nil
	}
}

func Sequence[T any](ps ...Parser[T]) Parser[[]T] {
	return func(s string) (string, []T, error) {
		tail := s
		var values []T
		for _, p := range ps {
			var value T
			var err error
			tail, value, err = p(tail)
			if err != nil {
				return s, nil, err
			}
			values = append(values, value)
		}
		return tail, values, nil
	}
}

func OneOf[T any](ps ...Parser[T]) Parser[T] {
	return func(s string) (string, T, error) {
		o := new(T)
		for _, p := range ps {
			tail, value, err := p(s)
			if err != nil {
				continue
			}
			return tail, value, nil
		}

		return s, *o, errors.New("no parser matched")
	}
}

func OneOrMore[T any](parser Parser[T]) Parser[[]T] {
	return func(s string) (string, []T, error) {
		var values []T
		var err error
		tail, value, err := parser(s)
		if err != nil {
			return s, nil, errors.New("one or more should at least matched once")
		}
		values = append(values, value)
		for !eof(tail) {
			var value T
			var err error
			tail, value, err = parser(tail)
			if err != nil {
				break
			}
			values = append(values, value)
		}
		return tail, values, nil
	}
}

func ZeroOrMore[T any](parser Parser[T]) Parser[[]T] {
	return func(s string) (string, []T, error) {
		var values []T
		tail := s
		for !eof(tail) {
			var value T
			var err error
			tail, value, err = parser(tail)
			if err != nil {
				break
			}
			values = append(values, value)
		}
		return tail, values, nil
	}
}

func ZeroOrOne[T any](parser Parser[T]) Parser[*T] {
	return func(s string) (string, *T, error) {

		tail, value, err := parser(s)
		if err != nil {
			return s, nil, nil
		}
		return tail, &value, nil
	}
}

func Map[IN, OUT any](p Parser[IN], f func(i IN) (OUT, error)) Parser[OUT] {
	return func(s string) (string, OUT, error) {

		tail, res, err := p(s)
		if err != nil {
			return s, *new(OUT), err
		}
		out, err := f(res)
		if err != nil {
			return s, *new(OUT), err
		}
		return tail, out, nil
	}
}

func Tag(tag string) Parser[string] {
	var chars []Parser[rune]
	for _, c := range tag {
		chars = append(chars, Char(c))
	}
	seq := Sequence(chars...)
	return Map(seq, func(rs []rune) (string, error) {
		return string(rs), nil
	})
}
