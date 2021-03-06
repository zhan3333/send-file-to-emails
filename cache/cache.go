package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"send-fiule-to-emails/util"
)

var filePath string
var data = map[string]interface{}{}

func Init(cacheFilePath string) {
	filePath = cacheFilePath
	if !util.IsFileExists(cacheFilePath) {
		f, err := os.Create(cacheFilePath)
		defer func() { _ = f.Close() }()
		if err != nil {
			panic(fmt.Sprintf("创建缓存文件失败: %+v\n", err))
		}
	}
	if err := load(); err != nil {
		panic(fmt.Sprintf("读取缓存失败: %s", err.Error()))
	}
}

func Set(key string, value interface{}) {
	data[key] = value
	save()
}

func Get(key string) interface{} {
	return data[key]
}

func save() {
	b, _ := json.Marshal(data)
	_ = ioutil.WriteFile(filePath, b, os.ModePerm)
}

func load() error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	if len(content) != 0 {
		err = json.Unmarshal(content, &data)
		if err != nil {
			return err
		}
	}
	return nil
}
