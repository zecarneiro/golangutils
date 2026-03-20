package netc

import (
	"net/http"
	"time"
)

func HasInternet() (bool, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	_, err := client.Get("https://www.google.com")
	if err != nil {
		return false, err
	}
	return true, nil
}
