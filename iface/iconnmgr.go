package iface

type IConnMgr interface {
	AddConn(connID int64, conn IConnection)
	RemoveConn(connID int64)
	GetConn(connID int64) IConnection
	StopAllConn()
	StopConn(connID int64)
}
