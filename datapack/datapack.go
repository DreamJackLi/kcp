package datapack

import (
	"encoding/binary"
	"gitlab.dove.im/wx/cc_server_common/ErrCollect"
	"gitlab.dove.im/wx/cc_server_common/iface"
	"gitlab.dove.im/wx/cc_server_common/message"
	"net"
	"time"
)

const (
	ApiTypeLen    = 4
	PackLen       = 4
	ServerTimeOut = 2 * time.Minute
)

type DataPack struct {
	connection net.Conn
}

func NewDataPack(connection net.Conn) iface.IDataPack {

	return &DataPack{
		connection: connection,
	}
}

func (this *DataPack) GetHeadLen() uint32 {
	return PackLen
}

// PackData(apiType EnumKcpMessageApiType, data IMessage) ([]byte, error)
func (this *DataPack) PackData(data iface.IMessage) ([]byte, error) {

	packLenBufWrite := make([]byte, ApiTypeLen+PackLen, ApiTypeLen+PackLen+data.GetDataLen())

	binary.BigEndian.PutUint32(packLenBufWrite[:ApiTypeLen], data.GetApiType())
	binary.BigEndian.PutUint32(packLenBufWrite[ApiTypeLen:ApiTypeLen+PackLen], data.GetDataLen())

	packLenBufWrite = append(packLenBufWrite, data.GetDataBody()...)

	return packLenBufWrite, nil
}

func (this *DataPack) UnPackData(collectEnum ErrCollect.EnumErrChanType, coll *ErrCollect.ErrCollect) (iface.IMessage, error) {

	packLenBuf := make([]byte, ApiTypeLen+PackLen)

	_, err := this.ReadFromConn(int(ApiTypeLen), packLenBuf[0:ApiTypeLen])

	if err != nil {
		coll.WriteError(collectEnum, err)
		return nil, err
	}

	apiType := binary.BigEndian.Uint32(packLenBuf[0:ApiTypeLen])

	// 校验 APi Type

	if !iface.CheckApiType(apiType) {
		//fmt.Println("CheckApiType is False ")
		return nil, nil
	}

	_, err = this.ReadFromConn(int(PackLen), packLenBuf[ApiTypeLen:])

	if err != nil {
		coll.WriteError(collectEnum, err)
		return nil, err
	}

	packLen := binary.BigEndian.Uint32(packLenBuf[ApiTypeLen:])

	//fmt.Println("UnPackData packLen ", packLen)

	msgBuf := make([]byte, packLen)

	_, err = this.ReadFromConn(int(packLen), msgBuf)

	if err != nil {
		coll.WriteError(collectEnum, err)
		return nil, err
	}

	m := message.NewKcpMessage(iface.Uint32ToApiType(apiType), packLen, msgBuf)

	return m, nil
}

func (this *DataPack) ReadFromConn(targetNum int, tarByte []byte) ([]byte, error) {

	this.connection.SetReadDeadline(time.Now().Add(ServerTimeOut))
	packLenN, err := this.connection.Read(tarByte)
	//readCount , err := io.ReadFull(this.Conn , msgBuf)

	if err != nil {
		//ErrCollect.CollectErr.WriteError(ErrCollect.EnumServerReadErr, err)
		return nil, err
	}

	if packLenN < targetNum {

		_, err := this.LoopRead(targetNum, packLenN, this.connection, tarByte)

		if err != nil {
			//ErrCollect.CollectErr.WriteError(ErrCollect.EnumServerReadErr, err)
			return nil, err
		} else {
			return tarByte, nil
		}

	}

	return tarByte, nil

}

func (this *DataPack) LoopRead(targetNum, sourceNum int, conn net.Conn, dataBuf []byte) (int, error) {

	if sourceNum < targetNum {

		for {

			if sourceNum >= targetNum {
				break
			}

			this.connection.SetReadDeadline(time.Now().Add(ServerTimeOut))
			n, err := this.connection.Read(dataBuf[sourceNum:])

			if err != nil {
				return sourceNum, err
			}
			sourceNum += n
		}

		return sourceNum, nil

	} else {
		return sourceNum, nil
	}

}
