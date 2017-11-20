package update

import (
	"fmt"
	"os"
	"io"
	"archive/zip"
	"bytes"
	"time"
	"io/ioutil"
	"regexp"
	"net/http"
	"strings"
)

type ResponseData struct {
	Release    `json:"data"`
	Msg string `json:"msg"`
	Ret int    `json:"ret"`
}

type Version struct {
	Major int `json:"major,omitempty"`
	Minor int `json:"minor,omitempty"`
	Patch int `json:"patch,omitempty"`
}

type Release struct {
	Version
	Body    string  `json:"body"`
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Url string `json:"browser_download_url"`
}

const (
	tagFmt      = "v%d.%d.%d"
	getApi      = "https://gitee.com/beewit/app"
	APPSpread   = "spread"
	APPSpreadDB = "spread-db"
	//apiUrl   = "https://gitee.com/api/v5/repos/beewit/spread/releases/latest?access_token=kdw2HGxYpTzVrdKpbQbV"
	//apiDBUrl = "https://gitee.com/api/v5/repos/beewit/spread-db/releases/latest?access_token=kdw2HGxYpTzVrdKpbQbV"
)

var ApiUrl string

func Logs(errStr string) {
	errStr = time.Now().Format("2006-01-02 15:04:05") + "   " + errStr
	file, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND, 0x644)
	defer file.Close()
	if err != nil {
		println(errStr)
	} else {
		file.Write([]byte(errStr))
	}
}

func init() {
	var err error
	ApiUrl, err = GetApi()
	if err != nil {
		Logs(fmt.Sprintf("update/update.go  获取GetApi失败，ERROR:%s", err.Error()))
	}
}

func GetApiByType(app string) (url string) {
	if ApiUrl == "" {
		return
	}
	if strings.Contains(ApiUrl, "?") {
		url = fmt.Sprintf("%s&app=%s", ApiUrl, app)
		return
	}
	url = fmt.Sprintf("%s?app=%s", ApiUrl, app)
	return
}

func GetApi() (url string, err error) {
	var resp *http.Response
	var dat []byte
	resp, err = http.Get(getApi)
	if err != nil {
		return
	}
	dat, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if regName := regexp.MustCompile(`<meta content='(.*)' name='Description'>`).FindAllStringSubmatch(string(dat), -1); len(regName) == 1 {
		url = regName[0][1]
		if url == "" {
			return
		}
		url = strings.Trim(url, "")
		reStr := "http(s)?://([\\w-]+\\.)+[\\w-]+(/[\\w- ./?%&=]*)?"
		flog, _ := regexp.Match(reStr, []byte(url))
		if !flog {
			url = ""
		}
	}
	return
}

func CopyFile(byte []byte, dst string) (w int64, err error) {
	Rename(dst+".old", dst+time.Now().Format("2006-01-02 15:04:05")+".old")
	Rename(dst, dst+".old")
	dstFile, err := os.Create(dst)
	if err != nil {
		return
	}
	defer dstFile.Close()
	return io.Copy(dstFile, bytes.NewReader(byte))
}

func Rename(dst, oldDst string) bool {
	flog, err := PathExists(dst)
	if err == nil && flog {
		err = os.Rename(dst, oldDst)
		if err == nil {
			return true
		}
	}
	return false
}

func Unzip(fileName string) {
	File, err := zip.OpenReader(fileName)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range File.File {
		info := v.FileInfo()
		if info.IsDir() {
			err := os.MkdirAll(v.Name, 0644)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}
		srcFile, err := v.Open()
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer srcFile.Close()
		newFile, err := os.Create(v.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		io.Copy(newFile, srcFile)
		newFile.Close()
	}
	defer File.Close()
}

func CreateFile(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (v Version) After(o Version) bool {
	if v.Major != o.Major {
		return v.Major > o.Major
	} else if v.Minor != o.Minor {
		return v.Minor > o.Minor
	}
	return v.Patch > o.Patch
}

func (gr Release) ToRelease() (rel Release) {
	var major, minor, patch int
	fmt.Sscanf(gr.TagName, tagFmt, &major, &minor, &patch)
	rel.Version = Version{major, minor, patch}
	rel.Body = gr.Body
	rel.TagName = gr.TagName
	for _, ga := range gr.Assets {
		rel.Assets = append(rel.Assets, Asset{ga.Url})
	}
	return
}
