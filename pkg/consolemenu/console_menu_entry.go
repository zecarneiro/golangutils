package consolemenu

type ConsoleMenuEntry struct {
	EntryNumber int
	Label       string
	Data        any

	entryProcessor func(ConsoleMenuEntry)
	isSeparator    bool
	alreadyRan     bool
}

func GetData[T any](e *ConsoleMenuEntry) T {
	return e.Data.(T)
}
