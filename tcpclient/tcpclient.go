package tcpclient

import (
	"errors"
	"fmt"
	"gitlab.dove.im/wx/cc_server_common/ErrCollect"
	"gitlab.dove.im/wx/cc_server_common/client"
	"gitlab.dove.im/wx/cc_server_common/connMgr"
	"gitlab.dove.im/wx/cc_server_common/datapack"
	"gitlab.dove.im/wx/cc_server_common/iface"
	"gitlab.dove.im/wx/cc_server_common/message"
	"net"
	"time"
)

const (
	ConnectionClientChanCount = 1000
	tcpHeartDuration          = 15
)

type TcpClient struct {
	ServerIP         string
	ServerPort       string
	ClientConn       net.Conn
	ConnID           int64
	isClose          bool
	WriteChan        chan []byte
	ReadChan         chan iface.IMessage
	DataPack         iface.IDataPack
	HeartCheck       *connMgr.KcpHeart
	ErrCollectClient *ErrCollect.ErrCollect
}

func NewTcpClient() *TcpClient {

	return &TcpClient{
		isClose:          false,
		WriteChan:        make(chan []byte, client.ConnectionClientChanCount),
		ReadChan:         make(chan iface.IMessage, client.ConnectionClientChanCount),
		DataPack:         nil,
		HeartCheck:       nil,
		ErrCollectClient: ErrCollect.NewErrCollect(),
	}
}

func (this *TcpClient) StartTcpClient(ip, port string) error {
	this.ServerPort = port
	this.ServerIP = ip
	serverAddr := fmt.Sprintf("%s:%s", this.ServerIP, this.ServerPort)
	fmt.Println("BeginClient serverAddr is:", serverAddr)

	udpSess, err := net.Dial("tcp4", serverAddr)
	//clientConn, err := kcp.Dial(serverAddr)

	if err != nil {
		this.StopClientWhenConnErr()
		fmt.Println("BeginClient Dial Err:", err)
		return err
	}

	this.ClientConn = udpSess
	this.HeartCheck = connMgr.NewKcpHeart(this)
	this.DataPack = datapack.NewDataPack(this.ClientConn)
	//CollectErr
	this.ErrCollectClient.AddCollect(ErrCollect.EnumClientErr)

	go this.clientTcpReadMessage()
	go this.clientTcpWriteMessage()
	go this.HeartCheck.StartHeart()
	return nil
}

func (this *TcpClient) clientTcpWriteMessage() {

	tick := time.Tick(tcpHeartDuration * time.Second)

	for {
		if this.isClose {
			return
		}

		select {

		case <-tick:
			mH := []byte("1")
			m := &message.KcpMessage{
				ApiType:  iface.EnumApiHeart,
				DataLen:  uint32(len(mH)),
				DataBody: mH,
			}
			this.WriteMessage(m)

		case data, ok := <-this.WriteChan:

			if !ok {
				return
			}
			m := &message.KcpMessage{
				ApiType:  iface.EnumApiNormal,
				DataLen:  uint32(len(data)),
				DataBody: data,
			}
			//fmt.Printf("Client Send Message APiType is %d, DataLen is %d\n ", m.ApiType, m.DataLen)
			this.WriteMessage(m)

		}
	}

}

func (this *TcpClient) clientTcpReadMessage() {

	for {
		if this.isClose {
			return
		}
		data, err := this.DataPack.UnPackData(ErrCollect.EnumClientErr, this.ErrCollectClient)

		if err != nil {
			this.ErrCollectClient.WriteError(ErrCollect.EnumClientErr, err)
			continue
		}
		this.HeartCheck.SetLastHeartTime(time.Now().UnixNano())

		if data.GetApiType() != iface.ApiTypeToUint32(iface.EnumApiHeart) {
			this.SetReadChanData(data)
		}

	}

}

func (this *TcpClient) SetReadChanData(iMessage iface.IMessage) {

	err, _ := this.ErrCollectClient.ReadError(ErrCollect.EnumClientErr)

	if err != nil {
		return
	}

	if this.isClose {
		this.ErrCollectClient.WriteError(ErrCollect.EnumClientErr, errors.New("Connection is Close"))
		return
	}

	this.ReadChan <- iMessage
}

func (this *TcpClient) WriteData(data []byte) error {

	err, _ := this.ErrCollectClient.ReadError(ErrCollect.EnumClientErr)

	if err != nil {
		return err
	}

	if this.isClose {
		return errors.New("Connection is Close")
	}

	if len(this.WriteChan) >= client.ConnectionClientChanCount {
		return errors.New("Client WriteChan is Full")
	}

	this.WriteChan <- data

	return nil
}

func (this *TcpClient) ReadData() ([]byte, error) {

	err, _ := this.ErrCollectClient.ReadError(ErrCollect.EnumClientErr)

	if err != nil {
		return nil, err
	}

	data, ok := <-this.ReadChan

	if ok {
		return data.GetDataBody(), nil
	} else {
		return nil, errors.New("ReadData Chan is Close")
	}
}

func (this *TcpClient) WriteMessage(m iface.IMessage) {

	msg, err := this.DataPack.PackData(m)

	if err != nil {
		fmt.Println("BeginClient ClientConn.Write err :", err)
		this.ErrCollectClient.WriteError(ErrCollect.EnumClientErr, err)
		return
	}

	_, err = this.ClientConn.Write(msg)

	if err != nil {
		fmt.Println("BeginClient ClientConn.Write err :", err)
		this.ErrCollectClient.WriteError(ErrCollect.EnumClientErr, err)
		return
	}
}

func (this *TcpClient) StopClient() (err error) {
	if this.isClose {
		return nil
	}
	this.isClose = true
	if this.ClientConn != nil {
		err = this.ClientConn.Close()
	}
	if this.HeartCheck != nil {
		this.HeartCheck.HeartStop()
	}

	close(this.ReadChan)
	close(this.WriteChan)
	this.ErrCollectClient.DeleteCollect(ErrCollect.EnumClientErr)
	return nil
}

func (this *TcpClient) SetHeartStatue(heartStatue bool) {
	this.StopClient()
	this.ErrCollectClient.WriteError(ErrCollect.EnumClientErr, errors.New("Heart TimeOut"))
}

func (this *TcpClient) StopClientWhenConnErr() {
	if this.isClose {
		return
	}
	this.isClose = true
	close(this.ReadChan)
	close(this.WriteChan)
	this.ErrCollectClient.DeleteCollect(ErrCollect.EnumClientErr)
}
