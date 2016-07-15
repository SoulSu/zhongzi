package config

import (
	"os"
	"log"
)

// 搜索相关
const (
	SEARCH_HTTP_LISTEN = ":8081"
	SEARCH_DIC_FILE = "dictonary.txt"
	SEARCH_DATA_SAVE_PATH = ""

	SEARCH_STORAGE_FOLDER = "data/search"
	SEARCH_STORAGE_SHARDS = 256

)

func ProjectPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Panic(err.Error())
	}

	return pwd
}
