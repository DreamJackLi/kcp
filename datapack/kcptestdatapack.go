package datapack

import (
	"encoding/binary"
	"gitlab.dove.im/wx/cc_server_common/message"
	"net"
)

type KcpTestDataPack struct {
	connection net.Conn
}

func NewKcpTestDataPack(connection net.Conn) *KcpTestDataPack {

	return &KcpTestDataPack{
		connection: connection,
	}
}

func (this *KcpTestDataPack) GetHeadLen() uint32 {

	return 0

}

func (this *KcpTestDataPack) UnPackData(oriData []byte) (*message.KcpTestMessage, error) {

	if len(oriData) < 12 {
		return nil, nil
	}

	dataLen := binary.BigEndian.Uint32(oriData[:4])
	curTime := binary.BigEndian.Uint64(oriData[4:12])
	packNum := binary.BigEndian.Uint32(oriData[12:16])

	m := &message.KcpTestMessage{
		PackNum:  packNum,
		DataLen:  dataLen,
		TestData: oriData[16:],
		CurTime:  curTime,
	}

	return m, nil

}

func (this *KcpTestDataPack) PackData(data *message.KcpTestMessage) ([]byte, error) {

	testMsg := make([]byte, 8+4+4, len(data.GetTestData())+8+4+4)

	binary.BigEndian.PutUint32(testMsg[:4], data.GetDataLen())
	binary.BigEndian.PutUint64(testMsg[4:12], data.GetCurTime())
	binary.BigEndian.PutUint32(testMsg[12:16], data.GetPackNum())
	testMsg = append(testMsg, data.GetTestData()...)

	return testMsg, nil
}
