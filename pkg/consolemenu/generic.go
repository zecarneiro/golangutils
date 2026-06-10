package consolemenu

import (
	"fmt"
	"golangutils/pkg/conv"
	"golangutils/pkg/logger"
)

const (
	entryNumberAndLabelFormat = `%d.`
)

func printSeparator() {
	fmt.Println("----")
}

func readUserInput(userInputMsg string, maxEntryNumer int) int {
	var response int
	printSeparator()
	for {
		var userResponse string
		fmt.Printf(`%s: `, userInputMsg)
		fmt.Scanln(&userResponse)
		userResponseInt, err := conv.StringToInt(userResponse)
		if err != nil || userResponseInt < 1 || userResponseInt > maxEntryNumer {
			logger.Warn("Invalid option inserted!")
		} else {
			response = userResponseInt
			break
		}
	}
	return response
}
