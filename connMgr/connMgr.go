package connMgr

import (
	"fmt"
	"gitlab.dove.im/wx/cc_server_common/iface"
	"sync"
)

type ConnMgr struct {
	connMap  map[int64]iface.IConnection
	connLock sync.Mutex
}

func NewConnMgr() *ConnMgr {
	return &ConnMgr{
		connMap:  make(map[int64]iface.IConnection),
		connLock: sync.Mutex{},
	}
}

func (this *ConnMgr) AddConn(connID int64, conn iface.IConnection) {

	this.connLock.Lock()
	defer this.connLock.Unlock()
	_, ok := this.connMap[connID]
	if !ok {
		this.connMap[connID] = conn
	}

}

func (this *ConnMgr) RemoveConn(connID int64) {

	this.connLock.Lock()
	defer this.connLock.Unlock()

	_, ok := this.connMap[connID]

	if ok {
		delete(this.connMap, connID)
		fmt.Println("RemoveConn ConnID is ", connID)
	}

}

func (this *ConnMgr) GetConn(connID int64) iface.IConnection {

	this.connLock.Lock()
	defer this.connLock.Unlock()
	conn, ok := this.connMap[connID]

	if ok {
		return conn
	} else {
		return nil
	}

}

func (this *ConnMgr) StopAllConn() {
	this.connLock.Lock()
	defer this.connLock.Unlock()
	for k, v := range this.connMap {
		v.Stop()
		delete(this.connMap, k)
	}
}

func (this *ConnMgr) StopConn(connID int64) {
	this.connLock.Lock()
	defer this.connLock.Unlock()
	_, ok := this.connMap[connID]

	if !ok {
		return
	}

	delete(this.connMap, connID)
}
