package message

type KcpTestMessage struct {
	PackNum  uint32
	DataLen  uint32
	TestData []byte
	CurTime  uint64
}

func NewKcpTestMessage(packNum uint32, testData []byte, curTime uint64) *KcpTestMessage {

	return &KcpTestMessage{
		PackNum:  packNum,
		DataLen:  uint32(len(testData)),
		TestData: testData,
		CurTime:  curTime,
	}

}

func (this *KcpTestMessage) GetTestData() []byte {
	return this.TestData
}

func (this *KcpTestMessage) GetCurTime() uint64 {
	return this.CurTime
}
func (this *KcpTestMessage) GetDataLen() uint32 {
	return this.DataLen
}
func (this *KcpTestMessage) GetPackNum() uint32 {
	return this.PackNum
}
