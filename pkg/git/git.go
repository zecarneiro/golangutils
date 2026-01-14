package git

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golangutils/pkg/logger"
	"golangutils/pkg/models"
	"golangutils/pkg/obj"
)

func GitGetLatestVersionRepo(owner string, repo string, isLatest bool) (models.GitRelease, error) {
	release := models.GitRelease{}
	urlsufix := "/latest"
	if !isLatest {
		urlsufix = ""
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases%s", owner, repo, urlsufix)
	logger.Info("Get Latest git version from url: " + url)
	myClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := myClient.Get(url)
	if err != nil {
		return release, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body) // response body is []byte
	if isLatest {
		err = obj.StringToObject(string(body), release)
	} else {
		var releaseArr []models.GitRelease
		err = obj.StringToObject(string(body), &releaseArr)
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
