package enums

import "golangutils/pkg/common"

type FileType int

const (
	Directory    FileType = 1
	File         FileType = 2
	SymbolicLink FileType = 3
)

func GetFileTypeFromValue(value int) FileType {
	switch value {
	case 1:
		return Directory
	case 2:
		return File
	case 3:
		return SymbolicLink
	}
	return common.UnknownInt
}

func (f FileType) IsValid() bool {
	switch f {
	case Directory, File, SymbolicLink:
		return true
	default:
		return false
	}
}

func (f FileType) Int() int {
	if f.IsValid() {
		return int(f)
	}
	return common.UnknownInt
}

func (f FileType) Equals(other FileType) bool {
	return f.Int() == other.Int()
}
