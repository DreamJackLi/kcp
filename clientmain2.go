package main

import (
	"fmt"
	"gitlab.dove.im/wx/cc_server_common/tcpclient"
	"math/rand"
	"time"
)

const (
	LiveTimeCap = 50
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getLifeTime(pos int) (int, int) {
	maxLifeTime := 0
	minLifeTime := 0
	pos++
	if pos <= 10 {
		minLifeTime = 30
		maxLifeTime = 120
	} else if pos <= 30 && pos > 10 {
		minLifeTime = 180
		maxLifeTime = 1800
	} else if pos <= 50 && pos > 30 {
		minLifeTime = 1800
		maxLifeTime = 7200
	}

	return minLifeTime, maxLifeTime
}

func main() {

	for i := 0; i < LiveTimeCap; i++ {
		pos := i
		go func() {

			buf := make([]byte, 128*1024)
			r := true
			tick := time.Tick(100 * time.Millisecond)
			// 获取 当前连接 存活时间
			minLifeTime, maxLifeTime := getLifeTime(pos)

			for {

				randLifeTime := rand.Intn(maxLifeTime-minLifeTime) + minLifeTime
				tickLifeTime := time.NewTicker(time.Duration(randLifeTime) * time.Second)

				tcpCli := tcpclient.NewTcpClient()
				//err := tcpCli.StartTcpClient("127.0.0.1", "55555")
				err := tcpCli.StartTcpClient("47.91.211.120", "55555")
				if err != nil {
					fmt.Println("StartTcpClient Error : ", err)
					time.Sleep(5 * time.Second)
					continue
				}
				r = true
				fmt.Printf("Start TcpClient Pos is %d \t tickLifeTime is %d \n: ", pos, randLifeTime)
				for r {
					select {
					case <-tickLifeTime.C:
						fmt.Println("Client Time Over pos is ", pos)
						tcpCli.StopClient()
						tickLifeTime.Stop()
						r = false
						break
					case <-tick:
						dataLen := rand.Intn(2 * 1024)
						bb := buf[:dataLen]
						rand.Read(bb)
						err = tcpCli.WriteData(bb)
						if err != nil {
							fmt.Println("Write Data Err ,", err)
							tcpCli.StopClient()
							tickLifeTime.Stop()
							r = false
							break
						}
						//fmt.Println("Write Data ")

						_, err = tcpCli.ReadData()
						if err != nil {
							fmt.Println("client Read Err ", err)
							tcpCli.StopClient()
							tickLifeTime.Stop()
							r = false
							break
						}

					}
				}
			}

		}()

		time.Sleep(1 * time.Second)

	}

	select {}
}
