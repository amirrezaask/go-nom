package nom

type Scanner interface {
	EOF() bool
	GetChar() rune
	Forward()
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
