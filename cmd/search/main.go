package main

import (
	"net/http"
	"log"

	"zhongzi/config"
)

func main() {

	http.HandleFunc("/add", Add)

	log.Println("start listen and serve on ", config.SEARCH_HTTP_LISTEN)
	log.Println(http.ListenAndServe(config.SEARCH_HTTP_LISTEN, nil))
}

func Add(rep http.ResponseWriter, req *http.Request) {




}
