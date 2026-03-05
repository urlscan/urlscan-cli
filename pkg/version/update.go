package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
)

const (
	repoOwner = "urlscan"
	repoName  = "urlscan-cli"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func CheckLatest() (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)

	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", nil
	}
	defer func() {
		closeErr := resp.Body.Close()
		if err == nil {
			err = closeErr
		}
	}()

	var release githubRelease
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(release.TagName, "v"), err
}

func IsNewer(current, latest string) bool {
	// for dev version (= an empty version), always show update notice for testing purpose
	if current == "" {
		return true
	}

	a, err := semver.NewVersion(current)
	if err != nil {
		return false
	}
	b, err := semver.NewVersion(latest)
	if err != nil {
		return false
	}
	return a.Compare(b) == -1
}
