package dispatcherbus

import (
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"sync"
)

// name: handlerName
type DispatcherMap map[string]*utils.DispatchHandlerStruct

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	diag Diagnostic

	dispatchers    map[string]*utils.Dispatcher
	dispatchersMtx sync.Mutex
}

func NewService(d Diagnostic) (*Service, error) {

	srv := &Service{
		diag:           d,
		dispatchers:    map[string]*utils.Dispatcher{},
		dispatchersMtx: sync.Mutex{},
	}

	return srv, nil
}

func (s *Service) Open() error {
	return nil
}

func (s *Service) Close() error {
	for key, v := range s.dispatchers {
		if v != nil {
			v.Release(key)
		}
	}

	return nil
}

func (s *Service) getDispatcher(name string) (*utils.Dispatcher, error) {
	s.dispatchersMtx.Lock()
	defer s.dispatchersMtx.Unlock()

	d, exist := s.dispatchers[name]
	if !exist {
		s.dispatchers[name] = utils.CreateDispatcher(utils.DefaultDispatcherBufLen)
		d = s.dispatchers[name]
		d.Start()
	}

	return d, nil
}

func (s *Service) Create(name string, length int) error {
	s.dispatchersMtx.Lock()
	defer s.dispatchersMtx.Unlock()

	_, exist := s.dispatchers[name]
	if exist {
		return errors.New("Dispatcher Already Exist")
	}

	s.dispatchers[name] = utils.CreateDispatcher(length)

	return nil
}

func (s *Service) Start(name string) error {
	dispatcher, err := s.getDispatcher(name)
	if err != nil {
		s.diag.Error("Start", err)
		return err
	}
	dispatcher.Start()
	return nil
}

func (s *Service) Release(name string, handlerID string) error {
	dispatcher, err := s.getDispatcher(name)
	if err != nil {
		return err
	}
	dispatcher.Release(handlerID)

	return nil
}

func (s *Service) Register(name string, handler *utils.DispatchHandlerStruct) {
	dispatcher, err := s.getDispatcher(name)
	if err != nil {
		// 如果dispatcher还没创建， 将handler加入注册列表等待创建后注册
		s.diag.Error("Register", errors.Errorf("Please Create Dispatcher: %s First", name))
		return
	}
	if handler.ID == "" {
		handler.ID = utils.GenerateID()
	}

	dispatcher.Register(handler.ID, handler.Handler)
}

func (s *Service) Dispatch(name string, data interface{}) error {
	dispatcher, err := s.getDispatcher(name)
	if err != nil {
		return err
	}

	return dispatcher.Dispatch(data)
}

// create, register and start
func (s *Service) LaunchDispatchersByHandlerMap(dispatcherMap DispatcherMap) {
	for name, handler := range dispatcherMap {
		err := s.Create(name, utils.DefaultDispatcherBufLen)
		if err != nil {
			s.diag.Debug(err.Error())
		}

		if handler != nil {
			s.Register(name, handler)
		}
		if err := s.Start(name); err != nil {
			s.diag.Error("Start Dispatcher Failed", err)
		}
	}
}

func (s *Service) ReleaseDispatchersByHandlerMap(dispatcherMap DispatcherMap) {
	for name, handler := range dispatcherMap {
		if handler == nil {
			continue
		}

		if err := s.Release(name, handler.ID); err != nil {
			s.diag.Error("Release Dispatcher Failed", err)
		}
	}
}
