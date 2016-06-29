package DhtCrawler

import "expvar"

var (
	// 接收到ping的个数
	KRPC_PING_COUNT = expvar.NewInt("ping_count")
	// 接收到getpeer的个数
	KRPC_GET_PEERS_COUNT = expvar.NewInt("get_peer_count")
	// 接收到announce_peer的个数
	KRPC_ANNOUNCE_PEER_COUNT = expvar.NewInt("announce_peer_count")

	// hashinfo信息
	REDIS_HASH_INFO_COUNT =expvar.NewString("redis_hashinfo_count")

)

