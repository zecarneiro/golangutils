package generic

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/logger"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type GitRelease struct {
	TagName string `json:"tag_name,omitempty"`
	Version string
}

func ProcessError(err error) {
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func GitGetLatestVersionRepo(owner string, repo string, isLatest bool) (GitRelease, error) {
	release := GitRelease{}
	urlsufix := "/latest"
	if !isLatest {
		urlsufix = ""
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases%s", owner, repo, urlsufix)
	logger.Info("Get Latest git version from url: " + url)
	var myClient = &http.Client{Timeout: 10 * time.Second}
	resp, err := myClient.Get(url)
	if err != nil {
		return release, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body) // response body is []byte
	if isLatest {
		release, err = common.StringToObject[GitRelease](string(body))
	} else {
		var releaseArr []GitRelease
		releaseArr, err = common.StringToObject[[]GitRelease](string(body))
		if len(releaseArr) > 0 {
			release = releaseArr[0]
		}
	}
	if err != nil {
		return release, err
	}
	release.Version = strings.TrimPrefix(release.TagName, "v")
	return release, nil
}

func Confirm(message string, isNoDefault bool) bool {
	yesNoMsg := "[y/N]"
	if !isNoDefault {
		yesNoMsg = "[Y/n]"
	}
	fmt.Printf("%s %s: ", message, yesNoMsg)
	var response string
	fmt.Scanln(&response)
	response = strings.Trim(response, " ")
	if response == "Y" || response == "y" {
		return true
	} else if len(response) == 0 {
		return common.Ternary(isNoDefault, false, true)
	}
	return false
}
