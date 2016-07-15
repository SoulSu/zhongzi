package search

import (
	//"github.com/huichen/sego"
	//"fmt"
	"github.com/huichen/wukong/engine"
	"github.com/huichen/wukong/types"
	//"os"
	"log"
	"runtime/debug"
)

func main() {
	//E:\linux_share\mygo\src\github.com\huichen\wukong\engine\engine.go
	log.SetFlags(log.Ldate|log.Llongfile)
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	//var segmenter sego.Segmenter
	//fmt.Println(os.Getwd())
	//segmenter.LoadDictionary("../../github.com/huichen/sego/data/dictionary.txt")
	//
	text := []byte("中华人民共和国中央人民政府支持普通模式和搜索模式两种分词")
	//segments := segmenter.Segment(text)
	//
	//// 处理分词结果
	//// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	//fmt.Println(sego.SegmentsToString(segments, false))

	var searcher = engine.Engine{}
	searcher.Init(types.EngineInitOptions{
		UsePersistentStorage: true,
		PersistentStorageFolder: "src/zz",
		PersistentStorageShards:255,
		SegmenterDictionaries: "../../github.com/huichen/sego/data/dictionary.txt",
		//StopTokenFile:         "../../data/stop_tokens.txt",
		IndexerInitOptions: &types.IndexerInitOptions{
			IndexType: types.LocationsIndex,
		},
	})


	searcher.IndexDocument(3,types.DocumentIndexData{
		Content: string(text),
	})

	searcher.IndexDocument(0, types.DocumentIndexData{Content: "此次百度收购将成中国互联网最大并购"})
	searcher.IndexDocument(1, types.DocumentIndexData{Content: "百度宣布拟全资收购91无线业务"})
	searcher.IndexDocument(2, types.DocumentIndexData{Content: "百度是中国最大的搜索引擎"})

	srep := searcher.Search(types.SearchRequest{Text:"中国"})

	log.Println(srep)

}
