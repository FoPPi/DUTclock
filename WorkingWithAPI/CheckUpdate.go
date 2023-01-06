package WorkingWithAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func CheckUpdateOnGitHub() (newVersion bool, release string) {
	LatestRelease, err := TakeJSON()
	if err != nil {
		fmt.Println(err)
	}

	if LatestRelease.TagName != Version {
		return true, LatestRelease.HtmlUrl
	} else {
		return false, ""
	}
}

func TakeJSON() (*LatestReleaseJSON, error) {

	// 1/1/1576/NEXT
	url := "https://api.github.com/repos/FoPPi/DUTclock/releases/latest"

	// Get request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("No response from request")
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error", err)
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	var result *LatestReleaseJSON
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
		return nil, err
	}

	return result, err
}

type LatestReleaseJSON struct {
	HtmlUrl string `json:"html_url"`
	TagName string `json:"tag_name"`
}
