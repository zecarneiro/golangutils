package consolemenu

import (
	"fmt"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/str"
	"strings"
)

type ConsoleMenu struct {
	menus            []*ConsoleMenuEntry
	exitMenu         ConsoleMenuEntry
	stopMenu         ConsoleMenuEntry
	userInputMessage string
	entryCounter     int
	maxEntryLabelLen int

	Title               string
	PreBuildCaller      func()
	PostProcessorCaller func()
	InsertStopMenu      bool
	DoneCheckMsg        string
}

func New() *ConsoleMenu {
	entity := &ConsoleMenu{
		userInputMessage: "Insert an option",
		exitMenu:         ConsoleMenuEntry{Label: "Exit"},
		stopMenu:         ConsoleMenuEntry{Label: "Stop"},
	}
	entity.Reset()
	return entity
}

// -------------------- PRIVATE --------------------
func (m *ConsoleMenu) buildMenu() {
	counter := 0
	for _, menuEntry := range m.menus {
		if menuEntry.isSeparator {
			printSeparator()
		} else {
			m.printEntryMenu(*menuEntry)
			counter = menuEntry.EntryNumber
		}
	}
	printSeparator()
	if m.InsertStopMenu {
		counter++
		m.stopMenu.EntryNumber = counter
		m.printEntryMenu(m.stopMenu)
	}
	counter++
	m.exitMenu.EntryNumber = counter
	m.printEntryMenu(m.exitMenu)
}

func (m *ConsoleMenu) printEntryMenu(menuEntry ConsoleMenuEntry) {
	counterLabel := fmt.Sprintf(entryNumberAndLabelFormat, menuEntry.EntryNumber)
	currentCounterSize := len(counterLabel)
	maxCounterSize := len(fmt.Sprintf(entryNumberAndLabelFormat, m.entryCounter))
	// Align text for: 1. with 10. with 100. with etc...
	if maxCounterSize > currentCounterSize {
		counterLabel = fmt.Sprintf(`%s%s`, counterLabel, strings.Repeat(" ", maxCounterSize-currentCounterSize))
	}
	label := menuEntry.Label
	// Add Done message on label but aligned
	if len(m.DoneCheckMsg) > 0 && menuEntry.alreadyRan && m.maxEntryLabelLen > 0 {
		spacesLen := m.maxEntryLabelLen - len(label)
		spaces := strings.Repeat(" ", spacesLen)
		label = fmt.Sprintf(`%s%s %s`, label, spaces, m.DoneCheckMsg)
	}
	logger.Log(fmt.Sprintf(`%s %s`, counterLabel, label))
}

func (m *ConsoleMenu) processEntryMenu(entryNumber int) {
	for _, menuEntry := range m.menus {
		if menuEntry.EntryNumber == entryNumber && menuEntry.entryProcessor != nil {
			menuEntry.entryProcessor(*menuEntry)
			menuEntry.alreadyRan = true
			break
		}
	}
}

// -------------------- PUBLIC --------------------
func (m *ConsoleMenu) ExitLabel(label string) {
	if !str.IsEmpty(label) {
		m.exitMenu.Label = label
	}
}

func (m *ConsoleMenu) StopMenuLabel(label string) {
	if !str.IsEmpty(label) {
		m.stopMenu.Label = label
	}
}

func (m *ConsoleMenu) UserInputMsg(msg string) {
	if !str.IsEmpty(msg) {
		m.userInputMessage = msg
	}
}

func (m *ConsoleMenu) Reset() {
	m.entryCounter = 0
	m.maxEntryLabelLen = 0
	m.menus = []*ConsoleMenuEntry{}
}

func (m *ConsoleMenu) AddEntry(label string, processor func(entryMenu ConsoleMenuEntry)) {
	m.AddEntryWithData(label, "", processor)
}

func (m *ConsoleMenu) AddEntryWithData(label string, data any, processor func(entryMenu ConsoleMenuEntry)) {
	if str.IsEmpty(label) {
		logic.ProcessError(fmt.Errorf("Invalid given menu entry label"))
	}
	if processor == nil {
		logic.ProcessError(fmt.Errorf("Invalid given menu entry processor"))
	}
	labelLen := len(label)
	if labelLen > m.maxEntryLabelLen {
		m.maxEntryLabelLen = labelLen
	}
	m.entryCounter++
	m.menus = append(m.menus, &ConsoleMenuEntry{EntryNumber: m.entryCounter, Label: label, Data: data, entryProcessor: processor})
}

func (m *ConsoleMenu) AddSeparator() {
	m.menus = append(m.menus, &ConsoleMenuEntry{isSeparator: true})
}

func (m *ConsoleMenu) Start() {
	for {
		if m.PreBuildCaller != nil {
			m.PreBuildCaller()
		}
		if !str.IsEmpty(m.Title) {
			logger.Title(m.Title)
		}
		m.buildMenu()
		userResponse := readUserInput(m.userInputMessage, m.exitMenu.EntryNumber)
		if m.InsertStopMenu && userResponse == m.stopMenu.EntryNumber {
			break
		}
		if userResponse == m.exitMenu.EntryNumber {
			logic.Exit(0)
		}
		m.processEntryMenu(userResponse)
		if m.PostProcessorCaller != nil {
			m.PostProcessorCaller()
		}
	}
}
