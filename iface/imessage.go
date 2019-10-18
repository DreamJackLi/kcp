package iface

type EnumKcpMessageApiType uint32

const (
	EnumEmpty     EnumKcpMessageApiType = 0
	EnumApiHeart  EnumKcpMessageApiType = 1
	EnumApiNormal EnumKcpMessageApiType = 2
)

func ApiTypeToUint32(apiType EnumKcpMessageApiType) uint32 {
	return uint32(apiType)
}

func Uint32ToApiType(apiInt uint32) EnumKcpMessageApiType {
	switch apiInt {
	case 1:
		return EnumApiHeart
	case 2:
		return EnumApiNormal
	}

	return EnumEmpty
}

func CheckApiType(apiInt uint32) bool {

	switch apiInt {
	case 1, 2:
		return true
	}

	return false
}

type IMessage interface {
	GetApiType() uint32
	GetDataBody() []byte
	GetDataLen() uint32
	SetApiType(apiType EnumKcpMessageApiType)
	SetDataBody(data []byte)
	SetDataLen(dataLen uint32)
}
