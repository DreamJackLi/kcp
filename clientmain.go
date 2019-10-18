package main

import (
	"fmt"
	"gitlab.dove.im/wx/cc_server_common/tcpclient"
	"math/rand"
	"time"
)

func main() {

	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt, os.Kill)
	//

	//for i := 0 ; i < 2000; i++{
	//	cli := client.NewKcpClient()
	//
	//	err := cli.StartClient("127.0.0.1", "20000")
	//
	//	if err != nil {
	//		fmt.Println("Client err", err)
	//		return
	//	}
	//
	//	go func() {
	//		for {
	//			err := cli.WriteData([]byte("client->" + time.Now().String()))
	//
	//			if err != nil {
	//
	//				continue
	//			}
	//			fmt.Println("Write Data ")
	//			time.Sleep(1 * time.Second)
	//		}
	//
	//	}()
	//
	//	go func() {
	//		for {
	//
	//			data, err := cli.ReadData()
	//			if err != nil {
	//				fmt.Println("client Read Err " , err)
	//				return
	//			}
	//			fmt.Println("client recv data is ", string(data))
	//
	//		}
	//	}()
	//}

	//fd, err := os.Open("111.jpg")
	//
	//if err != nil {
	//	fmt.Println("Open File Error :", err)
	//	return
	//}
	//
	//fileLen, _ := fd.Seek(0, io.SeekEnd)
	//
	//fileBuf := make([]byte, fileLen)
	//
	//fd.Seek(0, io.SeekStart)
	//
	//fd.Read(fileBuf)
	//
	//fmt.Println("File Size is ", fileLen)

	//cli := client.NewKcpClient()
	//
	//err = cli.StartClient("127.0.0.1", "20002")
	//
	//if err != nil {
	//	fmt.Println("Client err", err)
	//	return
	//}

	//go func() {
	//	for {
	//
	//		err = cli.WriteData(fileBuf)
	//		fmt.Println("", err)
	//		if err != nil {
	//			fmt.Println("Write Data Err ,", err)
	//			cli.StopClient()
	//			return
	//		}
	//		time.Sleep(1 * time.Second)
	//		fmt.Println("Write Data ")
	//
	//	}
	//
	//}()
	//
	//go func() {
	//	for {
	//
	//		data, err := cli.ReadData()
	//		if err != nil {
	//			fmt.Println("client Read Err ", err)
	//			cli.StopClient()
	//			return
	//		}
	//
	//		// 记录
	//		fmt.Println("client recv data is ", string(data))
	//
	//	}
	//}()

	//if err != nil {
	//	fmt.Println("Start Client Error " , err)
	//}

	buf := make([]byte, 128*1024)

	recvN := 0
	tick := time.Tick(50 * time.Millisecond)
	//tcpCli := tcpclient.NewTcpClient()
	//err := tcpCli.StartTcpClient("47.91.211.120" , "55555")
	//err := tcpCli.StartTcpClient("127.0.0.1" , "55555")
	//if err != nil {
	//	fmt.Println("StartTcpClient Error : " , err)
	//	return
	//}

	//err := tcpCli.StartTcpClient("127.0.0.1" , "55555")
	r := true
	go func() {
		for {
			tcpCli := tcpclient.NewTcpClient()
			//err := tcpCli.StartTcpClient("47.91.211.120" , "55555")
			err := tcpCli.StartTcpClient("127.0.0.1", "55555")
			r = true
			if err != nil {
				fmt.Println("StartTcpClient Error : ", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for r {
				select {
				case <-tick:
					dataLen := rand.Intn(5*1024) + 1024

					bb := buf[:dataLen]
					rand.Read(bb)
					recvN++
					err = tcpCli.WriteData(bb)
					if err != nil {
						fmt.Println("Write Data Err ,", err)
						r = false
						break
					}
					fmt.Println("Write Data ")

					_, err = tcpCli.ReadData()
					if err != nil {
						fmt.Println("client Read Err ", err)
						r = false
						break
					}
				}
			}

		}
	}()

	//go func() {
	//	for {
	//
	//		err = tcpCli.WriteData(fileBuf)
	//		fmt.Println("", err)
	//		if err != nil {
	//			fmt.Println("Write Data Err ,", err)
	//			tcpCli.StopClient()
	//			return
	//		}
	//		time.Sleep(1 * time.Second)
	//		fmt.Println("Client Write Data ")
	//
	//	}
	//
	//}()
	//
	//go func() {
	//	for {
	//
	//		_, err := tcpCli.ReadData()
	//		recvN++
	//		if err != nil {
	//			fmt.Println("client Read Err ", err)
	//			tcpCli.StopClient()
	//			return
	//		}
	//
	//		// 记录
	//		fmt.Println("client recv data is ", recvN)
	//
	//	}
	//}()

	select {}

}
