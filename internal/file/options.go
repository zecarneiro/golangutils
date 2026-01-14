package file

type Option func(*File)

func WithLogger(log *LoggerUtils) Option {
	return func(f *File) {
		f.
	}
}