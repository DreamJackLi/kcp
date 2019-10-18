package iface

import "gitlab.dove.im/wx/cc_server_common/ErrCollect"

type IDataPack interface {
	GetHeadLen() uint32
	UnPackData(collectEnum ErrCollect.EnumErrChanType, coll *ErrCollect.ErrCollect) (IMessage, error)
	PackData(data IMessage) ([]byte, error)
}
