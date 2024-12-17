package golangutils

import "fmt"

const ()

type Logger struct {
	keepLine                                                                                                bool
	colorReset, colorRed, colorGreen, colorYellow, colorBlue, colorPurple, colorCyan, colorWhite, colorGray string
}

func NewLogger() Logger {
	return Logger{
		keepLine:    false,
		colorReset:  "\033[0m",
		colorRed:    "\033[31m",
		colorGreen:  "\033[32m",
		colorYellow: "\033[33m",
		colorBlue:   "\033[34m",
		colorPurple: "\033[35m",
		colorCyan:   "\033[36m",
		colorWhite:  "\033[37m",
		colorGray:   "\033[90m",
	}
}

func (l *Logger) EnableKeepLine() {
	l.keepLine = true
}

func (l *Logger) disableKeepLine() {
	l.keepLine = false
}

func (l *Logger) Log(data string) {
	if l.keepLine {
		fmt.Print(data)
	} else {
		fmt.Println(data)
	}
	l.disableKeepLine()
}

func (l *Logger) Debug(data string) {
	l.EnableKeepLine()
	l.Log("[DEBUG] ")
	l.Log(data)
}

func (l *Logger) Warn(data string) {
	l.EnableKeepLine()
	l.Log(fmt.Sprintf("[%sWARN%s] ", string(l.colorYellow), string(l.colorReset)))
	l.Log(data)
}

func (l *Logger) Error(data string) {
	l.EnableKeepLine()
	l.Log(fmt.Sprintf("[%sERROR%s] ", string(l.colorRed), string(l.colorReset)))
	l.Log(data)
}

func (l *Logger) Info(data string) {
	l.EnableKeepLine()
	l.Log(fmt.Sprintf("[%sINFO%s] ", string(l.colorBlue), string(l.colorReset)))
	l.Log(data)
}

func (l *Logger) Ok(data string) {
	l.EnableKeepLine()
	l.Log(fmt.Sprintf("[%sOK%s] ", string(l.colorGreen), string(l.colorReset)))
	l.Log(data)
}

func (l *Logger) Prompt(data string) {
	l.EnableKeepLine()
	l.Log(fmt.Sprintf("%s>>>%s ", string(l.colorGray), string(l.colorReset)))
	l.Log(data)
}

func (l *Logger) Title(data string) {
	messageLen := len(data)
	separator := ""
	for i := 1; i < messageLen+7; i++ {
		separator = "#" + separator
	}
	l.Log(separator)
	l.Log("## " + data + " ##")
	l.Log(separator)
}

func (l *Logger) Header(data string) {
	l.Log("## " + data + " ##")
}

func (l *Logger) Separator(length int) {
	data := "# "
	if length < 6 {
		length = 6
	}
	for i := 1; i < length-4; i++ {
		data += "-"
	}
	data += " #"
	l.Log(data)
}

func (l *Logger) Help(appName string, description string, usages []string, args []string, others []string) {
	l.Log(appName + " - " + description)
	if len(usages) > 0 {
		l.Log("USAGE: ")
		for _, usage := range usages {
			l.Log("\t" + appName + " " + usage)
		}
	}
	if len(args) > 0 {
		l.Log("ARGS: ")
		args = append(args, "--help|-h: Show Help")
		for _, arg := range args {
			l.Log("\t" + arg)
		}
	}
	if len(others) > 0 {
		for _, other := range others {
			l.Log("\t" + other)
		}
	}
}

func (l *Logger) HasAndLogError(err error) bool {
	if err != nil {
		l.Error(err.Error())
		return true
	}
	return false
}
