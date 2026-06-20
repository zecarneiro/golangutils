package netc

import (
	"golangutils/pkg/logic"
	"net/http"
	"time"
)

func HasInternetWithoutErr() bool {
	status, err := HasInternet()
	return logic.Ternary(err != nil, false, status)
}

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
