package iface

type NewConnHandle func(conn IConnection)
type NewMessageHandle func([]byte, IConnection)

type IServer interface {
	StartServer(newConnHandle NewConnHandle, newMsgHandle NewMessageHandle) error
	//StartServerByUrl(ip , port string) error
	StopServer()

	GetConnMgr() IConnMgr
	GetServerPort() string
	//IOServer
}

type IOServer interface {
	//ReadData(connID int64)([]byte , error)
	WriteData(connID int64, data []byte) error
}
