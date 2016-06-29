package DhtCrawler

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type DhtNode struct {
	node       *KNode
	table      *KTable
	network    *Network
	log        *log.Logger
	master     chan string
	krpc       *KRPC
	outChan    chan string
	expVarPort int
}

func NewDhtNode(id *Id, logger io.Writer, outHashIdChan chan string, master chan string, expVarPort int) *DhtNode {
	node := new(KNode)
	node.Id = *id

	dht := new(DhtNode)
	// 用于数据传递的chan
	dht.outChan = outHashIdChan
	// 生成一个日志器
	dht.log = log.New(logger, "", log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	// 自己的节点信息
	dht.node = node

	dht.table = new(KTable)
	// 初始化本地监听udp的端口信息
	dht.network = NewNetwork(dht)

	dht.krpc = NewKRPC(dht)
	dht.master = master

	dht.expVarPort = expVarPort

	return dht
}

func (dht *DhtNode) Run() {

	//当前DHT节点运转进程
	go func() {
		dht.network.Listening()
	}()

	//自动结交更多DHT节点进程进程
	go func() {
		dht.NodeFinder()
	}()


	// dht统计信息服务
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", dht.expVarPort), nil)
		if err != nil {
			dht.log.Println(err.Error())
			return
		}
		// 统计redis中hashinfo信息
		//for {
		//	go func() {
		//		rds := redis.Get()
		//		defer rds.Close()
		//
		//		inf, _ := rds.Do("info", "Keyspace")
		//		REDIS_HASH_INFO_COUNT.Set(fmt.Sprintf("%v", inf))
		//
		//	}()
		//	// 10s 统计一次
		//	time.Sleep(time.Second * 10)
		//}
	}()

	dht.log.Println("start expvar on", dht.expVarPort)
	dht.log.Println(fmt.Sprintf("DhtCrawler %s is runing...", dht.network.Conn.LocalAddr().String()))
	for {
		select {
		case msg := <-dht.master:
			dht.log.Println(msg)
		}
	}
}
