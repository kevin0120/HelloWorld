package utils

import (
	"log"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/atomic"
)

const (
	DefaultDispatcherBufLen = 128
	DispatcherWorkers       = 16
)

type dispatchHandler func(data interface{})

type DispatchHandlerStruct struct {
	ID      string
	Handler dispatchHandler
}

func CreateDispatchHandlerStruct(h dispatchHandler) *DispatchHandlerStruct {
	if h == nil {
		return nil
	}
	return &DispatchHandlerStruct{
		ID:      GenerateID(),
		Handler: h,
	}
}

// bufLen: 缓冲长度
func CreateDispatcher(bufLen int) *Dispatcher {
	return &Dispatcher{
		mtx:         &sync.Mutex{},
		open:        atomic.NewBool(false),
		buf:         make(chan *dispatcherStruct, bufLen),
		closing:     make(chan bool, DispatcherWorkers), //协程数的关闭通道
		dispatchers: map[string]dispatchHandler{},
	}
}

type Dispatcher struct {
	mtx     *sync.Mutex
	open    *atomic.Bool
	buf     chan *dispatcherStruct
	closing chan bool

	dispatchers map[string]dispatchHandler
}

func (s *Dispatcher) removeHandler(handler string) {
	if _, ok := s.dispatchers[handler]; ok {
		s.mtx.Lock()
		defer s.mtx.Unlock()
		delete(s.dispatchers, handler)
	}
}

func (s *Dispatcher) getOpen() bool {
	return s.open.Load()
}

func (s *Dispatcher) setOpen(open bool) {
	s.open.Store(open)
}

func (s *Dispatcher) Start() {
	if !s.getOpen() {
		for i := 0; i < DispatcherWorkers; i++ {
			go s.manage()
		}
		s.setOpen(true)
	}
}

func (s *Dispatcher) Release(handler string) {
	if !s.getOpen() {
		return
	}
	s.removeHandler(handler)

	//如果都为空了 说明要将所有协程关闭
	if len(s.dispatchers) == 0 {
		for i := 0; i < DispatcherWorkers; i++ {
			s.closing <- true
		}
		s.setOpen(false)
	}
}

func (s *Dispatcher) Register(key string, dispatcher dispatchHandler) string {
	if key == "" {
		key = GenerateID()
	}
	s.dispatchers[key] = dispatcher
	return key
}

func (s *Dispatcher) Dispatch(data interface{}) error {
	if !s.getOpen() {
		msg := "Dispatcher Is Not Opened!!!"
		log.Println(msg)
		return errors.New(msg)
	}
	s.doDispatch(data)
	return nil
}

type dispatcherStruct struct {
	h    dispatchHandler
	data interface{}
}

//todo: 限制协程的数量
func (s *Dispatcher) doDispatch(data interface{}) {
	for _, v := range s.dispatchers {
		d := &dispatcherStruct{
			h:    v,
			data: data,
		}
		s.buf <- d
	}
}

func (s *Dispatcher) manage() {
	for {
		select {
		case hh := <-s.buf:
			if hh.h != nil && hh.data != nil {
				hh.h(hh.data)
			}
		case <-s.closing:
			log.Println("Dispatcher Is Closed!!!")
			return
		}
	}
}
