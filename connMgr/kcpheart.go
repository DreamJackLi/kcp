package connMgr

import (
	"gitlab.dove.im/wx/cc_server_common/iface"
	"time"
)

const (
	//  15s发一次  30s 检查一次
	CheckHeartDuration = 30 * time.Second
)

type KcpHeart struct {
	LastHeartTime int64
	//ExitChan      chan bool
	sendHeart iface.ISendHeart
	closed    bool
}

func NewKcpHeart(sendHeart iface.ISendHeart) *KcpHeart {
	return &KcpHeart{
		LastHeartTime: 0,
		//ExitChan:      make(chan bool),
		sendHeart: sendHeart,
	}
}

func (this *KcpHeart) SetLastHeartTime(lastTime int64) {
	this.LastHeartTime = lastTime
}

func (this *KcpHeart) StartHeart() {

	tick := time.Tick(CheckHeartDuration)

	for {

		if this.closed {
			return
		}

		select {

		case <-tick:
			// 检查 心跳
			curTime := time.Now().UnixNano()
			//temp := curTime - this.LastHeartTime
			//fmt.Println("Check Heart Time is ", temp)
			if curTime < this.LastHeartTime || curTime-this.LastHeartTime <= CheckHeartDuration.Nanoseconds() {
				// 心跳合法
				// 更新心跳
				this.LastHeartTime = curTime
			} else {
				// 心跳不合法，通知
				this.sendHeart.SetHeartStatue(true)
			}
		}

	}

}

func (this *KcpHeart) HeartStop() {
	if !this.closed {
		this.closed = true
		//this.ExitChan <- true
		//close(this.ExitChan)
	}

}
