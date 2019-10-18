package main

import (
	"fmt"
	"gitlab.dove.im/wx/cc_server_common/iface"
	"gitlab.dove.im/wx/cc_server_common/tcpserver"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

//var totalTcpDelay map[int64]int64
//var lastTcpSecDelay map[int64]int64
//var lastTcpSnmp int64
//
//var muTcp sync.Mutex

func main() {

	//readChan := make(chan []byte , 100)

	//s := server.NewServer("127.0.0.1", "20002")
	//47.91.211.120
	s := tcpserver.NewTcpServer("0.0.0.0", "55555")
	//err := s.StartServer(func(conn iface.IConnection) {
	//	fmt.Println("New Conn  ConnID is ", conn.GetConnID())
	//
	//	go func() {
	//
	//		for{
	//			err := conn.WriteData([]byte("qweqwe"))
	//
	//			if err != nil {
	//				fmt.Println("Write Err is ", err)
	//				continue
	//			}
	//
	//			fmt.Printf("Write To Conn ID is %d\n"  , conn.GetConnID())
	//
	//			time.Sleep(1 * time.Second)
	//		}
	//
	//
	//	}()
	//
	//	//conn.WriteData([]byte("qweqwe"))
	//	fmt.Println("New Conn  Conn Remote is ", conn.GetRemoteAdd())
	//}, func(bytes []byte, connection iface.IConnection) {
	//
	//	fmt.Printf("New Msg : Conn ID is  %d Msg is %s \n", connection.GetConnID() , string(bytes))
	//	//connection.ReadData()
	//})
	logfd, err := os.OpenFile("./log/server.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {
		fmt.Println("Open Log File Error :", err)
		return
	}
	//

	tpcLog := log.New(logfd, "TcpServer :", log.Ldate|log.Ltime)

	//
	tick := time.Tick(15 * time.Second)
	var mu sync.Mutex
	totalSendPkg := 0
	totalSendByte := 0
	totalRecvByte := 0

	go func() {

		for {

			select {
			case <-tick:
				mu.Lock()

				tpcLog.Printf("-------> total totalSendPkg is %d \t totalSendByte Bytes is %d \t totalRecvByte is %d \n", totalSendPkg, totalSendByte, totalRecvByte)
				totalSendPkg = 0
				totalSendByte = 0
				totalRecvByte = 0
				mu.Unlock()
			}

		}

	}()

	err = s.StartServer(func(conn iface.IConnection) {

		tpcLog.Printf("-------- New Conn  ConnID is %d \t Ip is %s \t Port is %s \n  ", conn.GetConnID(), conn.GetRemoteAdd(), s.GetServerPort())
		go func() {

			for {
				readData, err := conn.ReadData()

				if err != nil {
					fmt.Println("------- Read Err is ", err)
					conn.Stop()
					return
					//continue
				}

				readDataLen := len(readData)
				buf := make([]byte, readDataLen/2+readDataLen)
				var sendBufNum int
				if readDataLen != 0 {
					sendBufNum = rand.Intn(readDataLen) + readDataLen/2
				}

				sendBuf := buf[:sendBufNum]

				rand.Read(sendBuf)
				// 立即回复
				err = conn.WriteData(sendBuf)
				if err != nil {
					fmt.Println("------- Write Err is ", err)
					conn.Stop()
					return
				}
				//tpcLog.Printf("------- ReadData To Conn ID is %d recvdataLen is %d \n", conn.GetConnID() , len(readData))
				mu.Lock()
				totalSendPkg++
				totalSendByte += sendBufNum
				totalRecvByte += readDataLen
				mu.Unlock()

			}

		}()

	}, nil)

	if err != nil {
		fmt.Println(err)
	}

}
