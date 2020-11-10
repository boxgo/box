package errcode

var (
	// OK success, nil.
	OK error

	// 100001000-100001999 system reserve error
	ModSystem = Build(100, 1)

	// 100002000-100002999 config error
	ModConfig = Build(100, 2)

	// 100003000-100003999 rpc error
	ModRPC = Build(100, 3)

	// 100004000-100004999 schedule task error
	ModSchedule = Build(100, 4)

	// 100100000-100199999 db error
	ModDB        = Build(100, 100)
	ModMongoDB   = Build(100, 101)
	ModMysql     = Build(100, 102)
	ModPostresql = Build(100, 103)
	ModOracle    = Build(100, 104)
	ModRedis     = Build(100, 105)

	// 100200000-100299999 message queue error
	ModMQ       = Build(100, 200)
	ModKafka    = Build(100, 201)
	ModRocketMQ = Build(100, 202)
	ModRabbitMQ = Build(100, 203)

	// 200001000-200001999 user action error
	ErrUser = Build(200, 1)
)

var (
	ErrrModDBReadTimeout    = ModDB.Build(1, "db read timeout")
	ErrrModDBWriteTimeout   = ModDB.Build(2, "db write timeout")
	ErrrModDBDeleteTimeout  = ModDB.Build(3, "db delete timeout")
	ErrrModDBDeleteTimeout1 = ModDB.Build(3, "db delete timeout")
)
