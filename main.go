package main

import
(
	"fmt"
	"os"
	"runtime"

	"zhongzi/DhtCrawler"
	"zhongzi/save/redis"
)

var (
	//主进程
	master = make(chan string)
	//爬虫输出抓取到的hashIds通道
	outHashIdChan = make(chan string)

	max_process = runtime.NumCPU()
)

func init() {

	zhongzi_log_file, err := os.OpenFile("zhongzi.log", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666);
	if err != nil {
		panic(err.Error())
	}

	//defer zhongzi_log_file.Close()

	//开启的dht节点
	for i := 0; i < max_process; i++ {
		go func() {
			id := DhtCrawler.GenerateID()
			dhtNode := DhtCrawler.NewDhtNode(&id, zhongzi_log_file, outHashIdChan, master, 8888)
			dhtNode.Run()
		}()
	}



}

func redis_set(key, value string) {
	rds := redis.Get()

	defer rds.Close()

	rds.Do("SET", key, value)
}

func main() {

	runtime.GOMAXPROCS(max_process)

	for {
		select {
		//输出爬虫抓取的HashId结果
		case hashId := <-outHashIdChan:
			fmt.Println(hashId)
			redis_set(hashId, hashId)
		case msg := <-master:
			fmt.Println(msg)
		}
	}
}
