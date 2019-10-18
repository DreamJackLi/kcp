package iface

import "net"

type IConnection interface {
	GetConnID() int64
	GetConn() net.Conn
	ReadData() ([]byte, error)
	WriteData(data []byte) error
	Stop()
	CloseState() bool
	GetRemoteAdd() string
}
