package update

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Update(cur Version) (fileNames []string, err error) {
	rel, err := CheckUpdate(cur, false)
	if err != nil {
		return
	}
	fileNames, err = DownloadFiles("", rel, nil)
	if err != nil {
		return
	}
	return
}

func CheckUpdate(cur Version, compare bool) (rel Release, err error) {
	defer func() {
		Logs("结束更新程序...")
		if err := recover(); err != nil {
			errStr := fmt.Sprintf("《程序出现严重错误，终止运行！》，ERROR：%v", err)
			Logs(errStr)
		}
	}()
	apiUrl := GetApiByType(APPSpread)
	Logs(fmt.Sprintf("CheckUpdate Get %s...", apiUrl))
	resp, err := http.Get(apiUrl)
	Logs(fmt.Sprintf("CheckUpdate HTTP Success %s...", apiUrl))
	if err != nil {
		return
	}
	Logs(fmt.Sprintf("CheckUpdate Get Success %s...", apiUrl))
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	Logs(fmt.Sprintf("CheckUpdate resp.Body %s...", string(dat)))
	Logs("CheckUpdate ToRelease...")
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

func DownloadFiles(dir string, rel Release, successFun func(fileNames []string, rel Release)) (fileNames []string, err error) {
	Logs("DownloadFiles...")
	var fileName string
	defer func() {
		if err2 := recover(); err2 != nil {
			//回滚已修改的文件
			RollbackFile(fileNames)
		}
		if err != nil {
			//回滚已修改的文件
			RollbackFile(fileNames)
		} else {
			if successFun == nil {
				//删除历史文件
				RemoveFile(fileNames)
			} else {
				successFun(fileNames, rel)
			}
		}
	}()
	assets := rel.Assets
	if len(assets) <= 0 {
		err = errors.New("未获取到下载文件")
		return
	}
	for _, asset := range assets {
		fileName, err = DownloadFile(dir, asset)
		if err != nil {
			return
		}
		fileNames = append(fileNames, fileName)
	}
	return
}

func DownloadFile(dir string, asset Asset) (fileName string, err error) {
	resp, err := http.Get(asset.Url)

	if err != nil {
		return
	}
	dis := resp.Header.Get("content-disposition")
	if !strings.Contains(dis, "attachment;filename=") {
		err = errors.New("不是有效的下载文件")
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fileName = strings.Replace(dis, "attachment;filename=", "", 1)
	if dir != "" {
		fileName = dir + "/" + fileName
	}
	Logs(fmt.Sprintf("DownloadFile %s ...", fileName))
	_, err = CopyFile(b, fileName)
	if err != nil {
		return
	}
	if strings.Contains(fileName, ".zip") {
		Unzip(fileName)
	}
	return
}

func RollbackFile(fileNames []string) {
	if len(fileNames) > 0 {
		for _, name := range fileNames {
			err := os.Rename(name, name+".new")
			if err != nil {
				err = os.Rename(name+".old", name)
				if err != nil {
					os.Remove(name + ".new")
				}
			}
		}
	}
}

func RemoveFile(fileNames []string) {
	if len(fileNames) > 0 {
		for _, name := range fileNames {
			os.Remove(name + ".old")
		}
	}
}
