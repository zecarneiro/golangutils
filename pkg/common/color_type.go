package common

type ColorType string

const (
	Reset  ColorType = "\033[0m"
	Red    ColorType = "\033[31m"
	Green  ColorType = "\033[32m"
	Yellow ColorType = "\033[33m"
	Blue   ColorType = "\033[34m"
	Purple ColorType = "\033[35m"
	Cyan   ColorType = "\033[36m"
	White  ColorType = "\033[37m"
	Gray   ColorType = "\033[90m"
)

func (c ColorType) IsValid() bool {
	switch c {
	case Reset, Red, Green, Yellow, Blue, Purple, Cyan, White, Gray:
		return true
	default:
		return false
	}
}

func (c ColorType) String() string {
	if c.IsValid() {
		return string(c)
	}
	return ""
}

func (c ColorType) Equals(other ColorType) bool {
	return c.String() == other.String()
}
