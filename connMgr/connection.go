package connMgr

import (
	"gitlab.dove.im/wx/cc_server_common/ErrCollect"
	"gitlab.dove.im/wx/cc_server_common/datapack"
	"gitlab.dove.im/wx/cc_server_common/iface"
	"gitlab.dove.im/wx/cc_server_common/message"
	"time"

	//"cc_server_common/server"
	"errors"
	"fmt"
	"net"
)

const (
	ConnectionChanCount = 1000
	HeartChanCount      = 100
	//PackLen = 4
)

type ReadMessage struct {
	Data []byte
	Err  error
}

type Connection struct {
	KcpServer      iface.IServer
	Conn           net.Conn
	ConnID         int64
	isClosed       bool
	WriteChan      chan []byte
	WriteHeartChan chan []byte
	ReadChan       chan iface.IMessage
	//CloseChan      chan bool
	NewMsgHandler    iface.NewMessageHandle
	DataPack         iface.IDataPack
	heartCheck       *KcpHeart
	HeartTimeOut     bool // true 为 超时  false是未超时
	ErrCollectServer *ErrCollect.ErrCollect
}

func NewConnection(kcpServer iface.IServer, conn net.Conn, connID int64, msgHandler iface.NewMessageHandle) *Connection {

	c := &Connection{
		KcpServer: kcpServer,
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		WriteChan: make(chan []byte, ConnectionChanCount),
		ReadChan:  make(chan iface.IMessage, ConnectionChanCount),
		//CloseChan:      make(chan bool, 2),
		NewMsgHandler:    msgHandler,
		DataPack:         datapack.NewDataPack(conn),
		HeartTimeOut:     false,
		WriteHeartChan:   make(chan []byte, HeartChanCount),
		ErrCollectServer: ErrCollect.NewErrCollect(),
	}
	c.ErrCollectServer.AddCollect(ErrCollect.EnumServerErr)
	heart := NewKcpHeart(c)

	c.heartCheck = heart

	return c
}

func (this *Connection) Start() {
	//CollectErr

	go this.beginRead()
	go this.beginWrite()
	go this.heartCheck.StartHeart()
}

func (this *Connection) beginRead() {

	defer this.Stop()

	for {

		if this.isClosed {
			return
		}

		select {
		//case <-this.CloseChan:
		//	return
		default:

			msg, err := this.DataPack.UnPackData(ErrCollect.EnumServerErr, this.ErrCollectServer)
			if err != nil {
				this.ErrCollectServer.WriteError(ErrCollect.EnumServerErr, err)
				continue
			}

			if msg == nil {
				continue
			}

			this.heartCheck.SetLastHeartTime(time.Now().UnixNano())
			//fmt.Println("Read msg is  " , string(msgBuf))
			if this.NewMsgHandler != nil {
				//this.NewMsgHandler(msg, this)
			} else {
				if msg.GetApiType() != iface.ApiTypeToUint32(iface.EnumApiHeart) {
					this.SetReadData(msg)
					//this.ReadChan <- msg
				} else {
					//fmt.Println("Heart Msg.... ")
					this.SetHeart([]byte("2"))
				}
			}

		}

	}

}

func (this *Connection) SetHeartStatue(heartStatue bool) {
	this.HeartTimeOut = heartStatue
	this.Stop()
	this.ErrCollectServer.WriteError(ErrCollect.EnumServerErr, errors.New("Heart TimeOut"))
}

func (this *Connection) beginWrite() {

	defer this.Stop()
	defer this.KcpServer.GetConnMgr().RemoveConn(this.ConnID)

	for {
		if this.isClosed {
			return
		}
		select {
		//case <-this.CloseChan:
		//	return
		case data := <-this.WriteChan:

			m := message.NewKcpMessage(iface.EnumApiNormal, uint32(len(data)), data)

			this.SendData(m)

		case data := <-this.WriteHeartChan:

			m := message.NewKcpMessage(iface.EnumApiHeart, uint32(len(data)), data)

			this.SendData(m)
		}

	}

}

func (this *Connection) SendData(m iface.IMessage) {

	msg, err := this.DataPack.PackData(m)
	if err != nil {
		this.ErrCollectServer.WriteError(ErrCollect.EnumServerErr, err)
		return
	}

	_, err = this.Conn.Write(msg)

	if err != nil {
		this.ErrCollectServer.WriteError(ErrCollect.EnumServerErr, err)
		fmt.Println("beginWrite  Write Err :", err)
		return
	}
}

func (this *Connection) Stop() {

	if this.isClosed {
		return
	}
	this.isClosed = true
	(this.Conn).Close()

	//this.CloseChan <- true

	this.KcpServer.GetConnMgr().RemoveConn(this.ConnID)
	this.heartCheck.HeartStop()
	close(this.WriteChan)
	close(this.ReadChan)
	//close(this.CloseChan)
	close(this.WriteHeartChan)

	this.ErrCollectServer.DeleteCollect(ErrCollect.EnumServerErr)

}

func (this *Connection) ReadData() ([]byte, error) {

	if this.isClosed {
		this.ErrCollectServer.WriteError(ErrCollect.EnumServerErr, errors.New("Connection is Close"))
		return nil, errors.New("Connection is Close")
	}
	//if len(this.CloseChan) == 0 {
	//
	//	return nil, errors.New("Chan is Empty")
	//}

	err, _ := this.ErrCollectServer.ReadError(ErrCollect.EnumServerErr)

	if err != nil {
		return nil, err
	}

	data, ok := <-this.ReadChan

	if ok {
		return data.GetDataBody(), nil
	} else {
		return nil, errors.New("Chan is closed")
	}

}

func (this *Connection) SetReadData(data iface.IMessage) {

	if this.isClosed {
		this.ErrCollectServer.WriteError(ErrCollect.EnumServerErr, errors.New("Read Chan is Close"))
		return
	}

	err, _ := this.ErrCollectServer.ReadError(ErrCollect.EnumServerErr)

	if err != nil {
		return
	}

	this.ReadChan <- data

}

func (this *Connection) SetHeart(data []byte) {

	if this.isClosed {
		this.ErrCollectServer.WriteError(ErrCollect.EnumServerErr, errors.New("Read Chan is Close"))
		return
	}
	this.WriteHeartChan <- data
}

func (this *Connection) WriteData(data []byte) error {

	err, _ := this.ErrCollectServer.ReadError(ErrCollect.EnumServerErr)

	if err != nil {
		return err
	}

	if this.isClosed {
		this.ErrCollectServer.WriteError(ErrCollect.EnumServerErr, errors.New("Connection is Close"))
		return errors.New("Connection is Close")
	}

	if len(this.WriteChan) >= ConnectionChanCount {

		return errors.New("Server WriteChan is Full")
	}

	this.WriteChan <- data
	return nil
}

func (this *Connection) GetConnID() int64 {
	return this.ConnID
}

func (this *Connection) GetConn() net.Conn {
	return this.Conn
}

func (this *Connection) CloseState() bool {
	return this.isClosed
}

func (this *Connection) GetRemoteAdd() string {

	if this.isClosed {
		return ""
	}
	return this.Conn.RemoteAddr().String()
}
