package mix

type Sign bool

func (s *Sign) String() string {
	if *s {
		return "+"
	} else {
		return "-"
	}
}
