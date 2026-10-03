package ui

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/models"

	"github.com/ncruces/zenity"
)

func notify(title string, message string, icon string) error {
	return processZenityError(zenity.Notify(message, zenity.Title(title), zenity.Icon(icon)))
}

func InfoNofity(message string, icon string) error {
	return notify("Information", message, icon)
}

func WarnNofity(message string, icon string) error {
	return notify("Warnning", message, icon)
}

func ErrorNofity(message string, icon string) error {
	return notify("Error", message, icon)
}

func OkNofity(message string, icon string) error {
	return notify("Success", message, icon)
}

func SelectFileWithFilters(title string, filters zenity.FileFilters) models.Response[string] {
	response := models.Response[string]{Data: "", Error: nil}
	selectedFilePath, err := zenity.SelectFile(zenity.Title(title), filters)
	err = processZenityError(err)
	if err != nil {
		response.Error = err
	} else {
		if file.IsFile(selectedFilePath) {
			response.Data = selectedFilePath
		} else {
			response.Error = fmt.Errorf(`Accept file only`)
		}
	}
	return response
}

func SelectFile(title string) models.Response[string] {
	return SelectFileWithFilters(title, zenity.FileFilters{{Name: "All Files", Patterns: []string{"*"}}})
}

func SelectFolder(title string) models.Response[string] {
	response := models.Response[string]{Data: "", Error: nil}
	selectedFolderPath, err := zenity.SelectFile(zenity.Title(title), zenity.Directory())
	err = processZenityError(err)
	if err != nil {
		response.Error = err
	} else {
		if file.IsDir(selectedFolderPath) {
			response.Data = selectedFolderPath
		} else {
			response.Error = fmt.Errorf(`Accept directory only`)
		}
	}
	return response
}

func InfoDialog(title string, message string) error {
	return processZenityError(zenity.Info(message, zenity.Title(title)))
}

func WarnDialog(title string, message string) error {
	return processZenityError(zenity.Warning(message, zenity.Title(title)))
}

func ErrorDialog(title string, message string) error {
	return processZenityError(zenity.Error(message, zenity.Title(title)))
}

func SelectList(title string, message string, entries []string) models.Response[string] {
	response := models.Response[string]{Data: "", Error: nil}
	selectedEntry, err := zenity.List(message, entries, zenity.Title(title),
		zenity.DisallowEmpty(), // User must select a entry before to exit or can cancel
	)
	err = processZenityError(err)
	if err != nil {
		response.Error = err
	} else {
		response.Data = selectedEntry
	}
	return response
}

func MultiSelectList(title string, message string, entries []string) models.Response[[]string] {
	response := models.Response[[]string]{Data: []string{}, Error: nil}
	selectedEntries, err := zenity.ListMultiple(message, entries, zenity.Title(title), zenity.CheckList())
	err = processZenityError(err)
	if err != nil {
		response.Error = err
	} else {
		response.Data = selectedEntries
	}
	return response
}

func Input(title string, message string, defaultValue string) models.Response[string] {
	response := models.Response[string]{Data: defaultValue, Error: nil}
	selectedInput, err := zenity.Entry(message, zenity.Title(title), zenity.EntryText(defaultValue))
	err = processZenityError(err)
	if err != nil {
		response.Error = err
	} else {
		response.Data = selectedInput
	}
	return response
}
