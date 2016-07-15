package search

import (
	"github.com/huichen/wukong/engine"
	"github.com/huichen/wukong/types"
	"zhongzi/config"
)

type MySearcher struct {
	*engine.Engine
}

var searcher *engine.Engine

func _init() {
	searcher.Init(types.EngineInitOptions{
		UsePersistentStorage: true,
		PersistentStorageFolder: config.SEARCH_STORAGE_FOLDER,
		PersistentStorageShards:config.SEARCH_STORAGE_SHARDS,
		SegmenterDictionaries: config.SEARCH_DIC_FILE,
		//StopTokenFile:         "../../data/stop_tokens.txt",
		IndexerInitOptions: &types.IndexerInitOptions{
			IndexType: types.LocationsIndex,
		},
	})
}

func NewMySearcher() *MySearcher {
	if (searcher == nil) {
		searcher = new(engine.Engine)
		_init()
	}

	return &MySearcher{searcher}
}


// 添加索引
func (this *MySearcher)AddIndexDocumetn(oid uint64, content string, field interface{}) {
	this.IndexDocument(oid, types.DocumentIndexData{
		Content: content,
		Fields:field,
	})
}


// 搜索信息
func (this *MySearcher)SearchToOid() {

}





