package message

import "gitlab.dove.im/wx/cc_server_common/iface"

//type KcpMessageApiType uint32

//const (
//	EnumApiHeart  = 0
//	EnumApiNormal = 1
//)

type KcpMessage struct {
	ApiType  iface.EnumKcpMessageApiType
	DataLen  uint32
	DataBody []byte

	//DataPack iface.IDataPack
}

func NewKcpMessage(apiType iface.EnumKcpMessageApiType, dataLen uint32, dataBody []byte) iface.IMessage {

	return &KcpMessage{
		ApiType:  apiType,
		DataLen:  dataLen,
		DataBody: dataBody,
	}

}

/*

GetApiType() uint32
	GetDataBody() []byte
	GetDataLen() uint32

*/

func (this *KcpMessage) GetApiType() uint32 {
	return uint32(this.ApiType)
}

func (this *KcpMessage) GetDataBody() []byte {
	return this.DataBody
}

func (this *KcpMessage) GetDataLen() uint32 {
	return this.DataLen
}

func (this *KcpMessage) SetApiType(apiType iface.EnumKcpMessageApiType) {
	this.ApiType = apiType
}

func (this *KcpMessage) SetDataBody(data []byte) {
	this.DataBody = data
}

func (this *KcpMessage) SetDataLen(dataLen uint32) {
	this.DataLen = dataLen
}
