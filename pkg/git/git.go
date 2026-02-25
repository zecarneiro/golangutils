package git

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/models"
	"golangutils/pkg/netc"
	"golangutils/pkg/obj"
)

func GithubGetLatestVersionRepo(owner string, repo string, isLatest bool) (models.GitRelease, error) {
	var release models.GitRelease
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
		release, err = obj.StringToObject[models.GitRelease](string(body))
	} else {
		var releaseArr []models.GitRelease
		releaseArr, err = obj.StringToObject[[]models.GitRelease](string(body))
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

func DownloadFromGithubRepo(owner string, repo string, version string, filename string) error {
	urlBase := "https://github.com/%s/%s/releases/download/%s/%s"
	url := fmt.Sprintf(urlBase, owner, repo, version, filename)
	dstOutput, err := file.GetCurrentDir()
	if err != nil {
		return err
	}
	dstOutput = file.ResolvePath(dstOutput, filename)
	err = netc.Download(url, dstOutput)
	if err != nil {
		logger.Error(err)
		version = fmt.Sprintf("v%s", version)
		logger.Info(fmt.Sprintf("Try again with change on the version: %s", version))
		url = fmt.Sprintf(urlBase, owner, repo, version, filename)
	}
	errV := netc.Download(url, dstOutput)
	if errV != nil {
		return errV
	}
	return errV
}
