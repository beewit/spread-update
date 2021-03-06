package update

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func DBUpdate(dir string, cur Version, successFun func(fileNames []string, rel Release)) (fileNames []string, err error) {
	rel, err := DBCheckUpdate(cur, true)
	if err != nil {
		return
	}
	fileNames, err = DownloadFiles(dir, rel, successFun)
	if err != nil {
		return
	}
	return
}

func DBCheckUpdate(cur Version, compare bool) (rel Release, err error) {
	apiDBUrl := GetApiByType(APPSpreadDB)
	resp, err := http.Get(apiDBUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var rd ResponseData
	json.Unmarshal(dat, &rd)
	rel = rd.Release.ToRelease()
	if compare {
		if !rel.Version.After(cur) {
			err = errors.New("没有可用的更新")
		}
	}
	return
}
