package server

import (
	"fmt"
	"testing"
)

func TestNewServer(t *testing.T) {

	//s := NewServer()
	//
	//s.ServerBegin()

	//tempS := make([]int , 0 , 3)
	//
	//tempS = append(tempS , 1)
	//tempS = append(tempS , 2)
	//tempS = append(tempS , 3)
	//PrintSlice(tempS)
	//fmt.Println("++++++++++++++++++++++++++")
	//tempS = tempS[0:0]
	//PrintSlice(tempS)

	testChan := make(chan int, 10)
	//testChan<- 2

	close(testChan)

	fmt.Println(len(testChan))

}

func PrintSlice(sli []int) {

	for _, v := range sli {
		fmt.Println(v)
	}

}
