package update

import (
	"fmt"
	"os"
	"github.com/henrylee2cn/pholcus/common/mahonia"
	"io"
	"archive/zip"
	"bytes"
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
	apiDBUrl = "https://gitee.com/api/v5/repos/beewit/spread-db/releases/latest?access_token=kdw2HGxYpTzVrdKpbQbV"
)

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
	for _, ga := range gr.Assets {
		rel.Assets = append(rel.Assets, Asset{ga.Url})
	}
	return
}
