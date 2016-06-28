package DhtCrawler

import (
	"fmt"
	"github.com/zeebo/bencode"
	"net"
	"time"
)

var BOOTSTRAP []string = []string{
	"67.215.246.10:6881",
	"212.129.33.50:6881",
	"82.221.103.244:6881"}

func (dhtNode *DhtNode) FindNode(node *KNode) {
	var id Id
	if node.Id != nil {
		id = node.Id.Neighbor()
	} else {
		id = dhtNode.node.Id.Neighbor()
	}
	tid := dhtNode.krpc.GenTID()
	v := make(map[string]interface{})
	v["t"] = fmt.Sprintf("%d", tid)
	v["y"] = "q"
	v["q"] = "find_node"
	args := make(map[string]string)
	args["id"] = string(id)
	args["target"] = string(GenerateID())
	v["a"] = args
	data, err := bencode.EncodeString(v)
	if err != nil {
		dhtNode.log.Fatalln(err)
	}

	raddr := new(net.UDPAddr)
	raddr.IP = node.Ip
	raddr.Port = node.Port
	// 向节点发送信息
	err = dhtNode.network.Send([]byte(data), raddr)
	if err != nil {
		dhtNode.log.Println(err)
	}
}

func (dhtNode *DhtNode) NodeFinder() {

	for {
		//	dhtNode.log.Println(len(dhtNode.table.Nodes), "port: ==== ", dhtNode.node.Port)
		// 查找节点信息 向bootstrap里面的配置节点发送 find_node 信息

		// 如果桶里面节点为空
		if len(dhtNode.table.Nodes) == 0 {
			for _, host := range BOOTSTRAP {
				raddr, err := net.ResolveUDPAddr("udp", host)
				if err != nil {
					dhtNode.log.Fatalf("Resolve DNS error, %s\n", err)
					return
				}
				node := new(KNode)
				node.Port = raddr.Port
				node.Ip = raddr.IP
				node.Id = nil

				dhtNode.FindNode(node)
			}
		} else {
			for _, node := range dhtNode.table.Nodes {
				dhtNode.FindNode(node)
			}
			// 清空节点 再次随机找节点
			// 为了方便更多节点匹配出来
			dhtNode.table.Nodes = nil
			time.Sleep(1 * time.Second)
		}

	}
}
