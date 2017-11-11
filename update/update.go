package update

import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"errors"
	"os"
	"io"
	"archive/zip"
	"bytes"
	"net/url"
	"strings"
	"github.com/henrylee2cn/pholcus/common/mahonia"
	"fmt"
	"github.com/astaxie/beego/logs"
)

type Version struct {
	Major int `json:"major,omitempty"`
	Minor int `json:"minor,omitempty"`
	Patch int `json:"patch,omitempty"`
}

type Release struct {
	Version
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Url string `json:"browser_download_url"`
}

const (
	tagFmt = "v%d.%d.%d"
	apiUrl = "https://gitee.com/api/v5/repos/beewit/spread/releases/latest?access_token=kdw2HGxYpTzVrdKpbQbV"
)

func Upload(cur Version) (error) {
	rel, err := CheckUpload(cur)
	if err != nil {
		return err
	}
	err = DownloadFile(rel.Assets)
	if err != nil {
		return err
	}
	return nil
}

var Log = logs.GetBeeLogger()

func init() {
	conf := fmt.Sprintf(
		`{
			"filename": "%s",
			"maxdays": %s,
			"daily": %s,
			"rotate": %s,
			"level": %s,
			"separate": "[%s]"
		}`,
		"spread-update.log",
		"10",
		"true",
		"true",
		"7",
		"error",
	)

	logs.SetLogger(logs.AdapterMultiFile, conf)
	logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)
}

func CheckUpload(cur Version) (rel Release, err error) {
	Log.Info("正在检测版本..")
	resp, err := http.Get(apiUrl)
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
	rel = release.toRelease()

	if !rel.Version.after(cur) {
		err = errors.New("没有可用的更新")
	}
	return
}
func DownloadFile(assets []Asset) error {
	if len(assets) <= 0 {
		return errors.New("未获取到下载文件")
	}
	for i, asset := range assets {
		if i < len(assets)-1 {
			u, err := url.Parse(asset.Url)
			if err != nil {
				return err
			}
			v, err := url.ParseQuery(u.RawQuery)
			if err != nil {
				return err
			}
			newUrl := v.Get("u")
			fmt.Println("Downloading", newUrl)
			resp, err := http.Get(newUrl)

			if err != nil {
				return err
			}
			println(resp.Header.Get("content-type"))
			dis := resp.Header.Get("content-disposition")
			if !strings.Contains(dis, "attachment;filename=") {
				return errors.New("不是有效的下载文件")
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			fileName := strings.Replace(dis, "attachment;filename=", "", 1)
			CopyFile(b, fileName)
			if strings.Contains(fileName, ".zip") {
				Unzip(fileName)
			}
		}
	}
	return nil
}

func CopyFile(byte []byte, dst string) (w int64, err error) {
	//判断文件是否存在，存在则进行更改文件名作为历史备份
	flog, err := PathExists(dst)
	if err == nil && flog {
		os.Rename(dst, dst+".old")
	}
	dstFile, err := os.Create(dst)
	if err != nil {
		return
	}
	defer dstFile.Close()
	return io.Copy(dstFile, bytes.NewReader(byte))
}

func Unzip(fileName string) {
	File, err := zip.OpenReader(fileName)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range File.File {
		info := v.FileInfo()
		if info.IsDir() {
			err := os.MkdirAll(mahonia.NewDecoder("gb18030").ConvertString(v.Name), 0644)
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
		newFile, err := os.Create(mahonia.NewDecoder("gb18030").ConvertString(v.Name))
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

func (v Version) after(o Version) bool {
	if v.Major != o.Major {
		return v.Major > o.Major
	} else if v.Minor != o.Minor {
		return v.Minor > o.Minor
	}
	return v.Patch > o.Patch
}

func (gr Release) toRelease() (rel Release) {
	var major, minor, patch int
	fmt.Sscanf(gr.TagName, tagFmt, &major, &minor, &patch)
	rel.Version = Version{major, minor, patch}
	for _, ga := range gr.Assets {
		rel.Assets = append(rel.Assets, Asset{ga.Url})
	}
	return
}
