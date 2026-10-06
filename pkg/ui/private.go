package ui

import (
	"errors"

	"github.com/ncruces/zenity"
)

func processZenityError(err error) error {
	if err == zenity.ErrCanceled {
		return errors.New(userCancellationMessage)
	}
	return err
}
