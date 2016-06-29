package DhtCrawler

import (
	"bytes"
	"github.com/zeebo/bencode"
	"math"
	"net"
	"sync/atomic"
)

type action func(arg map[string]interface{}, raddr *net.UDPAddr)

type KRPC struct {
	Dht   *DhtNode
	Types map[string]action
	tid   uint32
}

func NewKRPC(dhtNode *DhtNode) *KRPC {
	krpc := new(KRPC)
	krpc.Dht = dhtNode

	return krpc
}

func (krpc *KRPC) GenTID() uint32 {
	return krpc.autoID() % math.MaxUint16
}

func (krpc *KRPC) autoID() uint32 {
	return atomic.AddUint32(&krpc.tid, 1)
}

// 将接收到的数据在这里处理
//
func (krpc *KRPC) Decode(data string, raddr *net.UDPAddr) error {
	val := make(map[string]interface{})

	if err := bencode.DecodeString(data, &val); err != nil {
		return err
	} else {
		var ok bool
		message := new(KRPCMessage)

		message.T, ok = val["t"].(string) //请求tid
		if !ok {
			return nil
		}

		message.Y, ok = val["y"].(string) //请求类型
		if !ok {
			return nil
		}

		message.Addr = raddr

		switch message.Y {
		case "q":
			query := new(Query)
			query.Y = val["q"].(string) //
			query.A = val["a"].(map[string]interface{})
			message.Addion = query
			break
		case "r":
			res := new(Response)
			res.R = val["r"].(map[string]interface{})
			message.Addion = res
			break
		case "e":
			// {"t":"aa", "y":"e", "e":[201, "A Generic Error Ocurred"]}
			krpc.Dht.log.Println("get error:", val["e"])
			break
		default:
			krpc.Dht.log.Println("invalid message")
			break
		}

		switch message.Y {
		case "q":
			krpc.Query(message)
			break
		case "r":
			krpc.Response(message)
			break
		}

	}
	return nil
}

func (krpc *KRPC) Response(msg *KRPCMessage) {
	if response, ok := msg.Addion.(*Response); ok {
		// 这是 find_node 发送后返回的节点信息
		if nodestr, ok := response.R["nodes"].(string); ok {
			nodes := ParseBytesStream([]byte(nodestr))
			for _, node := range nodes {
				krpc.Dht.table.Put(node)
			}
		}
	}
}

func (krpc *KRPC) Query(msg *KRPCMessage) {
	if query, ok := msg.Addion.(*Query); ok {
		if query.Y == "ping" {
			KRPC_PING_COUNT.Add(1)

		}
		// 查询端
		if query.Y == "get_peers" {
			KRPC_GET_PEERS_COUNT.Add(1)
			if infohash, ok := query.A["info_hash"].(string); ok {
				// 理论上这个 infohash十可靠的
				krpc.Dht.outChan <- Id(infohash).String()
				// 把自己现在保存的snode 发送给查询的地址
				nodes := ConvertByteStream(krpc.Dht.table.Snodes)
				data, _ := krpc.EncodingNodeResult(msg.T, "asdf13e", nodes)
				krpc.Dht.network.Send([]byte(data), msg.Addr)
			}
		}

		if query.Y == "announce_peer" {
			KRPC_ANNOUNCE_PEER_COUNT.Add(1)
			// 这里的infohash不是可靠的
			if infohash, ok := query.A["info_hash"].(string); ok {
				krpc.Dht.outChan <- Id(infohash).String()
			}
		}
	}
}

func ConvertByteStream(nodes []*KNode) []byte {
	buf := bytes.NewBuffer(nil)
	for _, v := range nodes {
		convertNodeInfo(buf, v)
	}
	return buf.Bytes()
}

func convertNodeInfo(buf *bytes.Buffer, v *KNode) {
	buf.Write(v.Id)
	convertIPPort(buf, v.Ip, v.Port)
}
func convertIPPort(buf *bytes.Buffer, ip net.IP, port int) {
	buf.Write(ip.To4())
	buf.WriteByte(byte((port & 0xFF00) >> 8))
	buf.WriteByte(byte(port & 0xFF))
}

func ParseBytesStream(data []byte) []*KNode {
	var nodes []*KNode = nil
	for j := 0; j < len(data); j = j + 26 {
		if j + 26 > len(data) {
			break
		}

		kn := data[j : j + 26]
		node := new(KNode)
		node.Id = Id(kn[0:20])
		node.Ip = kn[20:24]
		port := kn[24:26]
		node.Port = int(port[0]) << 8 + int(port[1])
		nodes = append(nodes, node)
	}
	return nodes
}

type KRPCMessage struct {
	T      string
	Y      string
	Addion interface{}
	Addr   *net.UDPAddr
}
/**
ping 
ping Query = {"t":"aa", "y":"q", "q":"ping", "a":{"id":"abcdefghij0123456789"}}
bencoded = d1:ad2:id20:abcdefghij0123456789e1:q4:ping1:t2:aa1:y1:qe
Response = {"t":"aa", "y":"r", "r": {"id":"mnopqrstuvwxyz123456"}}
bencoded = d1:rd2:id20:mnopqrstuvwxyz123456e1:t2:aa1:y1:re
 */
type Query struct {
	Y string // 请求的方法
	A map[string]interface{} // 附加的请求字段
}

type Response struct {
	R map[string]interface{}
}

func (krpc *KRPC) EncodingNodeResult(tid string, token string, nodes []byte) (string, error) {
	v := make(map[string]interface{})
	v["t"] = tid
	v["y"] = "r"
	args := make(map[string]string)
	args["id"] = string(krpc.Dht.node.Id)
	if token != "" {
		args["token"] = token
	}
	args["nodes"] = bytes.NewBuffer(nodes).String()
	v["r"] = args
	//krpc.Dht.log.Println(v)
	s, err := bencode.EncodeString(v)
	return s, err
}
