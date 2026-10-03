package ui

import (
	"fmt"

	"github.com/ncruces/zenity"
)

func processZenityError(err error) error {
	if err == zenity.ErrCanceled {
		return fmt.Errorf("Operation cancelled by user")
	}
	return err
}
