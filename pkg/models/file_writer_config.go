package models

type FileWriterConfig struct {
	File         string
	Data         string
	IsAppend     bool
	IsCreateDir  bool
	EncodingName string
	WithUtf8BOM  bool
}
