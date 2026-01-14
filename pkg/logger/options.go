package logger

func WithKeepLine(keep bool) {
	keepLine = keep
}

func WithHeaderLength(length int) {
	headerLength = length
}

func WithSeparatorLength(length int) {
	separatorLength = length
}

func WithLogFile(filepath string) {
	logFile = filepath
}
