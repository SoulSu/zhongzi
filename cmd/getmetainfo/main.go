package main

import (
	"net/http"
	"zhongzi/meta"
	"log"
	"fmt"
	"encoding/json"
)

func main() {

	http.HandleFunc("/get", get_metainfo)

	http.ListenAndServe(":9091", nil)

}

func get_metainfo(res http.ResponseWriter, req *http.Request) {

	metainfo := req.URL.Query().Get("hash")

	bt := server.New(metainfo)

	mtinfo, err := bt.GetTorrentMetaInfo()
	if err != nil {
		log.Println("get torrent metainfo err:", err.Error())
	}

	b,err:=json.Marshal(mtinfo)
	if err != nil {
		log.Println("get torrent metainfo err:", err.Error())
	}

	fmt.Fprint(res, string(b))

}

