package enum

type EFileType int

const (
	UNKNOWN_FILE_TYPE = -1
	DIRECTORY         = 1
	FILE              = 2
	SYMBOLIC_LINK     = 3
)

func EFileTypeFromValue(value int) EFileType {
	switch value {
	case 1:
		return DIRECTORY
	case 2:
		return FILE
	case 3:
		return SYMBOLIC_LINK
	}
	return UNKNOWN_FILE_TYPE
}

func (s EFileType) IsValid() bool {
	switch s {
	case UNKNOWN_FILE_TYPE, DIRECTORY, FILE, SYMBOLIC_LINK:
		return true
	default:
		return false
	}
}

func (s EFileType) Int() int {
	if s.IsValid() {
		return int(s)
	}
	return UNKNOWN_FILE_TYPE
}

func (s EFileType) Equals(other EFileType) bool {
	return s.Int() == other.Int()
}

func (s EFileType) IsUnknown() bool {
	return s.Int() < 1 || s.Int() > 3 || !s.IsValid()
}
