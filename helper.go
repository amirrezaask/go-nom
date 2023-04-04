package nom

import "strconv"

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

var IntParser = Transform(OneOrMore(DigitParser), func(cs []rune) (int, error) {
	i, err := strconv.Atoi(string(cs))
	if err != nil {
		return 0, err
	}
	return i, nil
})

var FloatParser = Transform(Transform(Sequence(
	OneOrMore(DigitParser),
	Transform(Sequence(Transform(Char('.'), func(r rune) ([]rune, error) { return []rune{r}, nil }), OneOrMore(DigitParser)),
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
