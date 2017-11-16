package update

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func DBUpdate(cur Version, successFun func(fileNames []string)) (fileNames []string, err error) {
	rel, err := DBCheckUpdate(cur, false)
	if err != nil {
		return
	}
	fileNames, err = DownloadFiles(rel.Assets, successFun)
	if err != nil {
		return
	}
	return
}

func DBCheckUpdate(cur Version, compare bool) (rel Release, err error) {
	resp, err := http.Get(apiDBUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var release Release
	json.Unmarshal(dat, &release)
	rel = release.ToRelease()
	if compare {
		if !rel.Version.After(cur) {
			err = errors.New("没有可用的更新")
		}
	}
	return
}
