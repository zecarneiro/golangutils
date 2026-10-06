package ui

import "golangutils/pkg/str"

const USER_CANCELLATION_MESSAGE = "Operation cancelled by user"

var userCancellationMessage = USER_CANCELLATION_MESSAGE

func WithUserCancelMessage(message string) {
	if str.IsEmpty(message) && userCancellationMessage != USER_CANCELLATION_MESSAGE {
		userCancellationMessage = USER_CANCELLATION_MESSAGE
	} else {
		userCancellationMessage = message
	}
}
