package ErrCollect

import (
	"sync"
)

type EnumErrChanType int

const (
	ErrChanCount = 10

	EnumServerErr EnumErrChanType = 0
	EnumClientErr EnumErrChanType = 1
	//EnumServerWriteErr EnumErrChanType = 0
	//EnumServerReadErr  EnumErrChanType = 1
	//EnumClientWriteErr EnumErrChanType = 2
	//EnumClientReadErr  EnumErrChanType = 3
)

type ErrCollect struct {
	collectHandler map[EnumErrChanType][]error
	collectLock    sync.Mutex
}

//func init() {
//
//	if CollectErr == nil {
//		CollectErr = NewErrCollect()
//	}
//
//}

func NewErrCollect() *ErrCollect {

	return &ErrCollect{
		collectHandler: make(map[EnumErrChanType][]error),
		collectLock:    sync.Mutex{},
	}

}

func (this *ErrCollect) AddCollect(key EnumErrChanType) {
	this.collectLock.Lock()
	defer this.collectLock.Unlock()
	temp := make([]error, 0, 10)

	this.collectHandler[key] = temp
}

func (this *ErrCollect) DeleteCollect(key EnumErrChanType) {

	this.collectLock.Lock()
	defer this.collectLock.Unlock()
	_, ok := this.collectHandler[key]

	if ok {
		delete(this.collectHandler, key)
	}

}

func (this *ErrCollect) GetCollect(key EnumErrChanType) []error {

	this.collectLock.Lock()
	defer this.collectLock.Unlock()

	errChan, ok := this.collectHandler[key]

	if ok {
		//res := errChan[0]
		//errChan = errChan[1:]
		return errChan
	}
	return nil

}

func (this *ErrCollect) WriteError(key EnumErrChanType, err error) bool {

	this.collectLock.Lock()
	defer this.collectLock.Unlock()

	errChan, ok := this.collectHandler[key]

	if !ok {
		return false
	}

	if len(errChan) == ErrChanCount {
		return false
	}

	errChan = append(errChan, err)
	return true
}

func (this *ErrCollect) ReadError(key EnumErrChanType) (error, bool) {

	this.collectLock.Lock()
	defer this.collectLock.Unlock()

	errChan, okC := this.collectHandler[key]

	if !okC {
		return nil, false
	}

	if errChan == nil || len(errChan) == 0 {
		return nil, false
	}

	err := errChan[0]

	//data, okE := <-errChan

	errChan = errChan[1:]

	return err, true

}
