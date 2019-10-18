package client

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestNewKcpClient(t *testing.T) {

	//client := NewKcpClient()
	//
	//client.BeginClient()
	//

	//temp := make([]int, 0, 4)
	//
	//temp = append(temp, 1)
	//temp = append(temp, 2)
	//temp = append(temp, 3)
	//temp = append(temp, 4)
	//
	////temp = append(temp[:2], temp[3:]...)
	////
	//for _, v := range temp {
	//
	//	fmt.Println(v)
	//
	//}
	//
	//temp1 := make([]int, 0, 1)
	//temp1 = append(temp1, 1)
	//
	//temp1 = temp1[1:]
	////
	////temp = temp[0:0]
	////
	//for _, v := range temp1 {
	//
	//	fmt.Println(v)
	//
	//}

	//timeNano := time.Now().UnixNano()
	//
	//testByte := make([]byte, 8)
	//
	//binary.BigEndian.PutUint64(testByte, uint64(timeNano))
	//
	//var temp uint64
	//buf := bytes.NewBuffer(testByte)
	//
	//err := binary.Read(buf, binary.BigEndian, &temp)
	//
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//
	//	fmt.Println("Success Read ", temp)
	//
	//}

	//curTime := time.Now().UnixNano()
	//fmt.Println(curTime)
	//
	//fmt.Println(uint64(curTime))

	//tick := time.Tick(1 * time.Second)
	//
	//temp := make(chan int, 10)
	//
	//go func() {
	//
	//	for {
	//		temp <- 1
	//		time.Sleep(1 * time.Second)
	//	}
	//
	//}()
	//
	//for {
	//
	//	select {
	//	case <-tick:
	//
	//		fmt.Println("tick..........Start")
	//		time.Sleep(2 * time.Second)
	//		fmt.Println("tick..........End")
	//	case data := <-temp:
	//		fmt.Println("Print data ", data)
	//	}
	//
	//}
	//rand.Seed(time.Now().UnixNano())
	//buf := make([]byte , 10)
	//sendBuf := buf[:rand.Intn(5)+2]
	//
	//rand.Read(sendBuf)

	//temp := make([]int , 2)
	//temp = append(temp , 1)
	//temp = append(temp , 2)
	//
	//temp1 := temp[:2]
	//
	//fmt.Println(temp1)

	//for i := 0; i < 3; i++ {
	//
	//	//pos := i
	//	go func() {
	//		fmt.Println(i)
	//		time.Sleep(3 * time.Second)
	//	}()
	//}
	//
	//select {}

	//tempS := make([]int32, 0, 4)
	//tempS = append(tempS, 1)
	//tempS = append(tempS, 2)
	//tempS = append(tempS, 3)
	//tempS = append(tempS, 12)
	//
	//var temp int32
	//temp = 12
	//
	//res := CheckInSlice(temp, tempS)
	//
	//fmt.Println(res)

	curTime := time.Now().Unix()
	fmt.Println(curTime)
	time.Sleep(1 * time.Second)

	newTime := time.Now().Unix()
	fmt.Println(newTime)

	temp := newTime - curTime

	//

	fmt.Println(temp)

	temp1 := 30 * time.Second
	fmt.Println(temp1.Seconds())

}

func CheckInSlice(ele interface{}, dataSlice interface{}) int {

	dataSliceValue := reflect.ValueOf(dataSlice)

	for i := 0; i < dataSliceValue.Len(); i++ {
		checkRes := CheckDeepEqual(dataSliceValue.Index(i), ele)

		if checkRes {
			return i
		}
	}

	return -1

}

func CheckDeepEqual(value reflect.Value, ele interface{}) bool {

	switch value.Type().Name() {

	case "int32":
		v, ok := value.Interface().(int32)
		if !ok {
			return false
		} else {
			return reflect.DeepEqual(v, ele)
		}
	}

	return false
}
