package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Paths    []*FolderPair   `json:"paths"`
	BaiduYun *BaiduYunConfig `json:"baidu_yun"`
}

type BaiduYunConfig struct {
	AccessCode string `json:"access_code"`

	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type FolderPair struct {
	LocalPath   string `json:"local_path"`
	NetworkPath string `json:"network_path"`
}

func main() {
	// 读取配置文件
	configPath := ""
	for _, val := range os.Args {
		if strings.HasPrefix(val, "--config=") {
			configPath = strings.Split(val, "=")[1]
		}
	}

	if len(configPath) == 0 {
		fmt.Println("config path not found")
		return
	}

	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("read config file like:\n%v\n", string(buf))

	config := &Config{}
	if err = json.Unmarshal(buf, config); err != nil {
		panic(err)
	}

	// step 1 读取百度相关的文件路径

	// Step 2 批次读取 -> 操作本地读取的文件路径
	// TODO Step 2.1 预处理，文件过大就切割一下
	localPathList := make([]string, 0)
	for _, fp := range config.Paths {
		if err := filepath.Walk(fp.LocalPath, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			// 过滤隐藏文件
			if strings.HasPrefix(info.Name(), ".") {
				return nil
			}

			paths := strings.Split(path, fp.LocalPath)
			pathOne := paths[1]
			if strings.HasPrefix(pathOne, "/") {
				pathOne = pathOne[1:]
			}

			localPathList = append(localPathList, pathOne)
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
