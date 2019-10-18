package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/xtaci/kcp-go"

	"gitlab.dove.im/wx/cc_server_common/client"
	"gitlab.dove.im/wx/cc_server_common/datapack"
	"gitlab.dove.im/wx/cc_server_common/message"
)

var totalDelay map[int64]int64
var lastSecDelay map[int64]int64
var lastSnmp *kcp.Snmp

var mu sync.Mutex

func main() {

	totalDelay = make(map[int64]int64)
	lastSecDelay = make(map[int64]int64)
	lastSnmp = kcp.DefaultSnmp.Copy()

	rand.Seed(time.Now().UnixNano())

	fd, err := os.Open("111.jpg")
	if err != nil {
		fmt.Println("Open File Error :", err)
		return
	}

	fileLen, _ := fd.Seek(0, io.SeekEnd)

	fileBuf := make([]byte, fileLen)

	fd.Seek(0, io.SeekStart)

	fd.Read(fileBuf)

	fmt.Println("File Size is ", fileLen)

	//countData := countdata.NewCountData()

	cli := client.NewKcpClient()

	err = cli.StartClient("192.168.30.115", "20002")
	//
	if err != nil {
		fmt.Println("Client err", err)
		return
	}

	go func() {
		tick := time.Tick(1 * time.Second)

		n := 0
		for {
			select {
			case <-tick:
				n++

				mu.Lock()

				curSnmp := kcp.DefaultSnmp.Copy()

				ss := make([]string, 0, len(lastSecDelay)+len(totalDelay)+3)

				ss = append(ss, fmt.Sprintf("\n===================: %d", n))
				ss = append(ss, fmt.Sprintf("last recvPkgs: %d, recvBytes: %d, byte2: %d", curSnmp.OutPkts-lastSnmp.OutPkts, curSnmp.OutBytes-lastSnmp.OutBytes, curSnmp.BytesSent-lastSnmp.BytesSent))
				// ids := make([]int, 0, len(lastSecDelay))
				// for k := range lastSecDelay {
				// 	ids = append(ids, int(k))
				// }
				// sort.Ints(ids)
				// for _, k := range ids {
				// 	v := lastSecDelay[int64(k)]
				// 	ss = append(ss, fmt.Sprintf("%04d => %d", k, v))
				// }

				if n%10 == 0 {
					ss = append(ss, fmt.Sprintf("total recvPkgs: %d, recvBytes: %d, byte2: %d", curSnmp.OutPkts, curSnmp.OutBytes, curSnmp.BytesSent))

					headers := curSnmp.Header()
					s := curSnmp.ToSlice()

					for k, v := range headers {
						ss = append(ss, fmt.Sprintf("%s => %s", v, s[k]))
					}

					// ids := make([]int, 0, len(totalDelay))
					// for k := range totalDelay {
					// 	ids = append(ids, int(k))
					// }
					// sort.Ints(ids)
					// for _, k := range ids {
					// 	//						v := totalDelay[int64(k)]
					// 	//						ss = append(ss, fmt.Sprintf("%04d => %d", k, v))
					// }
				}

				lastSecDelay = make(map[int64]int64)

				lastSnmp = curSnmp

				mu.Unlock()

				fmt.Println(strings.Join(ss, "\n"))
			}
		}

	}()

	//pack := NewKcpTestDataPack()
	//	go countData.StartOneSecondCount()
	//	go countData.StartTenSecondCount()

	kcpTestPack := datapack.NewKcpTestDataPack(cli.ClientConn)
	n := 0

	for cnt := 0; cnt < 10000; {

		//time.Sleep(1 * time.Second)
		// 获取一个 1~32 之间的数
		//fmt.Println("Client Send Time1 ", nowTime)

		dataRatio := rand.Intn(16) + 16
		dataLen := dataRatio * 64

		nowTime1 := time.Now().UnixNano()
		//		fmt.Println("Client Send Time1 ", nowTime1)
		n++
		kcpTest := message.NewKcpTestMessage(uint32(n), fileBuf[:dataLen], uint64(nowTime1))

		// 打包
		kcpTestData, err := kcpTestPack.PackData(kcpTest)

		//fmt.Println("---------------Client WriteData 100 is ", kcpTestData[:100])
		//fmt.Println("---------------Client WriteData -100 is ", kcpTestData[len(kcpTestData)-101:])

		err = cli.WriteData(kcpTestData)

		if err != nil {
			fmt.Println("Write Data Err ,", err)
			cli.StopClient()
			return
		}

		//fmt.Println("Write Data ")

		_, err = cli.ReadData()
		if err != nil {
			fmt.Println("client Read Err ", err)
			cli.StopClient()
			return
		}

		//fmt.Println("---------------Client ReadData 100 is ", data[:100])
		//fmt.Println("---------------Client ReadData -100 is ", data[len(data)-101:])
		// 记录 当前时间
		//  EnumServerErr
		//msg, err := kcpTestPack.UnPackData(data)
		//
		//if err != nil {
		//	fmt.Println("client Read Err ", err)
		//	cli.StopClient()
		//	return
		//}
		//
		//if msg == nil {
		//	continue
		//}

		nowTime := time.Now().UnixNano()
		handlerTime := nowTime - nowTime1 //int64(msg.GetCurTime())
		handlerTime /= 1000
		handlerTime /= 1000
		//fmt.Println("client Handler data is ", handlerTime)
		//fmt.Println("client Cur Time is ", nowTime)
		//fmt.Println("client recv data is CurTime ", msg.GetCurTime())

		mu.Lock()

		if v, ok := totalDelay[handlerTime]; ok {
			totalDelay[handlerTime] = v + 1
		} else {
			totalDelay[handlerTime] = 1
		}

		if v, ok := lastSecDelay[handlerTime]; ok {
			lastSecDelay[handlerTime] = v + 1
		} else {
			lastSecDelay[handlerTime] = 1
		}

		mu.Unlock()
	}

	select {}

}
