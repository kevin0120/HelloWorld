package tightening_device

import (
	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/encryption"
	"sync"
	"sync/atomic"
	"time"

	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/httpd"
	"gopkg.in/resty.v1"

	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
)

const (
	ModelDesoutterCvi3        = "ModelDesoutterCvi3"
	ModelDesoutterCvi3Twin    = "ModelDesoutterCvi3Twin"
	ModelDesoutterCvi2        = "ModelDesoutterCvi2"
	ModelDesoutterCvi2R       = "ModelDesoutterCvi2R"
	ModelDesoutterCvi2L       = "ModelDesoutterCvi2L"
	ModelDesoutterDeltaWrench = "ModelDesoutterDeltaWrench"
	ModelDesoutterConnector   = "ModelDesoutterConnector"
	ModelCraneIQWrench        = "ModelCraneIQWrench"
	ModelLexenWrench          = "ModelLexenWrench"
)

var ENV_UNRESTRICT_TOOL_CTRL = utils.GetEnvBool("ENV_UNRESTRICT_TOOL_CTRL", false) //单工位单工具时可只用此环境变量

type Service struct {
	diag               Diagnostic
	configValue        atomic.Value
	runningControllers map[string]ITighteningController
	mtxDevices         sync.Mutex
	protocols          map[string]ITighteningProtocol

	storageService IStorageService
	dispatcherBus  Dispatcher
	deviceService  IDeviceService
	dispatcherMap  dispatcherbus.DispatcherMap

	ioService  IOService
	httpd      HTTPService
	httpClient *resty.Client

	// websocket请求处理器
	wsnotify.WSRequestHandlers
}

func (s *Service) ensureHttpClient() *resty.Client {
	if s.httpClient != nil {
		return s.httpClient
	}
	client := resty.New()
	client.SetRESTMode() // restful mode is default
	client.SetTimeout(8 * time.Second)
	client.SetContentLength(true)
	// Headers for all request
	client.
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second)

	s.httpClient = client
	return client
}

func (s *Service) loadTighteningController(c Config) {
	authPayload := encryption.GetAuthPayload()
	controllerNumber := authPayload.Controllers
	if controllerNumber == 0 {
		msg := "未获取到控制器授权，启动控制器失败!"
		s.diag.Error(msg, errors.New("privilege grant failed"))
		return
	}
	for k, deviceConfig := range c.Devices {
		if k >= controllerNumber && controllerNumber != -1 {
			msg := "可启动控制器数量达到上限!"
			s.diag.Error(msg, errors.New("privilege num empty"))
			break
		}
		p, err := s.getProtocol(deviceConfig.Protocol)
		if err != nil {
			s.diag.Error("loadTighteningController", err)
			continue
		}

		c, err := p.NewController(&c.Devices[k], s.dispatcherBus) //如果不传index或导致获取的配置信息有误
		if err != nil {
			s.diag.Error("Create Controller Failed", err)
			continue
		}

		sn := deviceConfig.SN
		if sn == "" {
			sn = fmt.Sprintf("%d", k+1)
		}

		c.SetSerialNumber(sn)
		s.ioService.AddModule(fmt.Sprintf(TIGHTENING_CONTROLLER_IO_SN_FORMAT, c.SerialNumber()), c.CreateIO())
		s.addController(sn, c)
	}
}

func NewService(c Config, d Diagnostic, protocols []ITighteningProtocol, dp Dispatcher, ds IDeviceService, db IStorageService, io IOService, httpd HTTPService) (*Service, error) {

	s := &Service{
		diag:               d,
		dispatcherBus:      dp,
		runningControllers: map[string]ITighteningController{},
		protocols:          map[string]ITighteningProtocol{},
		deviceService:      ds,
		storageService:     db,
		ioService:          io,
		httpd:              httpd,
	}

	s.setupGlobalDispatchers()
	s.setupWSRequestHandlers()
	s.setupHttpRoute()
	s.ensureHttpClient()

	s.configValue.Store(c)

	// 载入支持的协议
	for _, protocol := range protocols {
		s.protocols[protocol.Name()] = protocol
	}

	// 根据配置加载所有拧紧控制器
	s.loadTighteningController(c)
	return s, nil
}

func (s *Service) getProtocol(protocolName string) (ITighteningProtocol, error) {
	if p, ok := s.protocols[protocolName]; !ok {
		return nil, errors.New("Protocol Is Not Support")
	} else {
		return p, nil
	}
}

func (s *Service) setupHttpRoute() {

	// @TODO 重复的接口
	//r = httpd.Route{
	//	RouteType:   httpd.ROUTE_TYPE_HTTP,
	//	Method:      "PUT",
	//	Pattern:     "/tool-enable",
	//	HandlerFunc: s.putToolEnable,
	//}
	//s.httpd.AddNewHttpHandler(r)

	r := httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/tool-pset",
		HandlerFunc: s.putToolPSet,
	}
	if err := s.httpd.AddNewHttpHandler(r); err != nil {
		s.diag.Error("AddNewHttpHandler tool-pset Error", err)
	}
}

func (s *Service) setupGlobalDispatchers() {
	s.dispatcherMap = dispatcherbus.DispatcherMap{
		dispatcherbus.DispatcherCurve: utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DispatcherJob:   utils.CreateDispatchHandlerStruct(nil),
	}
}

func (s *Service) setupWSRequestHandlers() {
	s.WSRequestHandlers = wsnotify.WSRequestHandlers{
		Diag: s.diag,
	}

	s.SetupHandlers(wsnotify.WSRequestHandlerMap{
		wsnotify.WS_TOOL_MODE_SELECT:       s.OnWS_TOOL_MODE_SELECT,
		wsnotify.WS_TOOL_ENABLE:            s.OnWS_TOOL_ENABLE,
		wsnotify.WS_TOOL_JOB:               s.OnWS_TOOL_JOB,
		wsnotify.WS_TOOL_PSET:              s.OnWS_TOOL_PSET,
		wsnotify.WS_TOOL_PSET_BATCH:        s.OnWS_TOOL_PSET_BATCH,
		wsnotify.WS_TOOL_PSET_LIST:         s.OnWS_TOOL_PSET_LIST,
		wsnotify.WS_TOOL_PSET_DETAIL:       s.OnWS_TOOL_PSET_DETAIL,
		wsnotify.WS_TOOL_JOB_LIST:          s.OnWS_TOOL_JOB_LIST,
		wsnotify.WS_TOOL_JOB_DETAIL:        s.OnWS_TOOL_JOB_DETAIL,
		wsnotify.WS_TOOL_RESULT_MANUAL_SET: s.OnWS_TOOL_RESULT_MANUAL_SET,
	})
}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	s.dispatcherBus.LaunchDispatchersByHandlerMap(s.dispatcherMap)
	s.initDispatcherRegisters()

	// 启动所有拧紧控制器
	s.startupControllers()

	return nil
}

func (s *Service) Close() error {

	s.dispatcherBus.ReleaseDispatchersByHandlerMap(s.dispatcherMap)

	// 关闭所有控制器
	s.shutdownControllers()

	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) initDispatcherRegisters() {
	// 注册websocket请求
	s.dispatcherBus.Register(dispatcherbus.DispatcherWsNotify, utils.CreateDispatchHandlerStruct(s.HandleWSRequest))

	// 套筒控制
	s.dispatcherBus.Register(dispatcherbus.DispatcherSocketSelector, utils.CreateDispatchHandlerStruct(s.handlerSocketSelector))
}

func (s *Service) doDispatch(name string, data interface{}) {
	if err := s.dispatcherBus.Dispatch(name, data); err != nil {
		s.diag.Error(fmt.Sprintf("doDispatch: %s", name), err)
	}
}

func (s *Service) GetControllers() map[string]ITighteningController {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	return s.runningControllers
}

func (s *Service) addController(controllerSN string, controller ITighteningController) {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	_, exist := s.runningControllers[controllerSN]
	if exist {
		return
	}

	s.runningControllers[controllerSN] = controller
}

func (s *Service) getController(controllerSN string) (ITighteningController, error) {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	td, exist := s.runningControllers[controllerSN]
	if !exist {
		return nil, errors.New(fmt.Sprintf("Controller %s Not Found", controllerSN))
	}

	return td, nil
}

func (s *Service) getTool(controllerSN, tighteningUnit string) (ITighteningTool, error) {
	controller, err := s.getController(controllerSN)
	if err != nil {
		return nil, err
	}

	tool, err := controller.GetToolViaSerialNumber(tighteningUnit)
	if err == nil {
		return tool, nil
	}
	if !ENV_UNRESTRICT_TOOL_CTRL {
		return nil, err
	}
	tool, err = s.getFirstTool(controller)

	return tool, err
}

func (s *Service) getFirstTool(controller ITighteningController) (ITighteningTool, error) {
	if controller == nil {
		return nil, errors.New("getFirstTool: Controller Is Nil")
	}

	for _, v := range controller.Children() {
		return v.(ITighteningTool), nil
	}

	return nil, errors.New("getFirstTool: Controller's Tool Not Found")
}

func (s *Service) GetControllerByToolSN(toolSN string) (ITighteningController, error) {
	for k, c := range s.runningControllers {
		for _, t := range c.Children() {
			if t.SerialNumber() == toolSN {
				return s.runningControllers[k].(ITighteningController), nil
			}
		}
	}

	return nil, errors.New("Controller Not Found")
}

func (s *Service) GetToolByToolSN(toolSN string) (ITighteningTool, error) {
	for _, c := range s.runningControllers {
		for _, t := range c.Children() {
			if t.SerialNumber() == toolSN {
				return t.(ITighteningTool), nil
			}
		}
	}

	return nil, errors.Errorf("Tool: %s Is  Not Found", toolSN)
}

func (s *Service) startupControllers() {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	for sn, c := range s.runningControllers {
		if _, err := s.storageService.CreateController(sn, c.Name()); err != nil {
			s.diag.Error(fmt.Sprintf("CreateController Serial Number: %s Error", sn), err)
			continue
		}
		err := c.Start()
		if err != nil {
			s.diag.Error(fmt.Sprintf("Startup Controller Serial Number: %s Failed", sn), err)
		}

		s.deviceService.AddDevice(sn, c)
	}
}

func (s *Service) shutdownControllers() {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	for _, c := range s.runningControllers {
		err := c.Stop()
		if err != nil {
			s.diag.Error("Shutdown Controller Failed", err)
		}
	}
}

func (s *Service) ToolLedControl(toolSN string, enable bool) error {
	locationBody, err := s.storageService.GetToolLocation(toolSN)
	if err != nil {
		s.diag.Error(fmt.Sprintf("Can Not Found ToolLocation VIA Tool SerialNumber: %s", toolSN), err)
		return err
	}

	var location Location
	_ = json.Unmarshal([]byte(locationBody), &location)
	var status uint16 = 0
	if enable {
		status = 1
	}

	err = s.ioService.Write(location.EquipmentSN, uint16(location.Output), status)
	if err != nil {
		s.diag.Error(fmt.Sprintf("IO Write Equipment: %s Error", location.EquipmentSN), err)
	}

	return err
}

func (s *Service) doSelectorPatch() {
	// io套筒什么都不做直接返回
	switch s.config().SocketSelector.Type {
	case IOSelectorType:
		return
	default:
		s.doDispatch(dispatcherbus.DispatcherSocketSelector, SocketSelectorReq{
			Type: SocketSelectorClear,
		})
	}
}
