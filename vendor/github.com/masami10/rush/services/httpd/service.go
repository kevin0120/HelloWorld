package httpd

import (
	stdContext "context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/arl/statsviz"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/websocket"
	"github.com/masami10/rush/services/diagnostic"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"go.uber.org/atomic"
)

const (
	// Root path for the API
	BasePath     = "/rush/v1"
	PathVersion2 = "/rush/v2"
)

type Diagnostic interface {
	NewHTTPServerErrorLogger() *log.Logger
	NewAccessLogger() *log.Logger
	StartingService()
	StoppedService()
	ShutdownTimeout()
	AuthenticationEnabled(enabled bool)

	ListeningOn(addr string, proto string)

	WriteBodyReceived(body string)

	HTTP(
		host string,
		username string,
		start time.Time,
		method string,
		uri string,
		proto string,
		status int,
		referer string,
		userAgent string,
		reqID string,
		duration time.Duration,
	)

	Error(msg string, err error)
	RecoveryError(
		msg string,
		err string,
		host string,
		username string,
		start time.Time,
		method string,
		uri string,
		proto string,
		status int,
		referer string,
		userAgent string,
		reqID string,
		duration time.Duration,
	)
}

type Service struct {
	addr string
	err  chan error

	methods Methods

	ApiDoc          string
	Handler         []*Handler
	shutdownTimeout time.Duration
	externalURL     string
	server          *iris.Application

	stop chan chan struct{}

	HandlerByNames map[string]int

	cors CorsConfig

	configValue atomic.Value

	diag   Diagnostic
	opened bool

	DiagService
	httpServerErrorLogger *log.Logger
}

func NewService(doc string, c Config, hostname string, d Diagnostic, disc *diagnostic.Service) (*Service, error) {

	port, _ := c.Port()
	u := url.URL{
		Host:   fmt.Sprintf("%s:%d", hostname, port),
		Scheme: "http",
	}
	s := &Service{
		ApiDoc:                doc,
		addr:                  c.BindAddress,
		externalURL:           u.String(),
		cors:                  c.Cors,
		server:                iris.New(),
		err:                   make(chan error, 1),
		HandlerByNames:        make(map[string]int),
		shutdownTimeout:       time.Duration(c.ShutdownTimeout),
		diag:                  d,
		DiagService:           disc,
		httpServerErrorLogger: d.NewHTTPServerErrorLogger(),
		opened:                false,
	}

	s.configValue.Store(c)

	crs := cors.New(cors.Options{
		AllowedOrigins:   s.cors.AllowedOrigins,
		AllowCredentials: s.cors.AllowCredentials,
		AllowedMethods:   s.cors.AllowedMethods,
	})

	s.server.WrapRouter(crs.ServeHTTP)

	s.methods = newHttpMethods(s)
	//TODO: iris升级后增加统一的validator进行参数的验证
	ac := s.diag.NewAccessLogger().Writer()
	//TODO: iris升级后使用accesslog中间件替代
	if c.AccessLog {
		c := logger.Config{
			Status:             true,
			IP:                 true,
			Method:             true,
			Path:               true,
			Query:              true,
			MessageContextKeys: []string{"logger_message"},
			MessageHeaderKeys:  []string{"User-Agent"},
			LogFunc: func(endTime time.Time, latency time.Duration, status, ip, method, path string, message interface{}, headerMessage interface{}) {
				output := logger.Columnize(endTime.Format("2006/01/02 - 15:04:05"), latency, status, ip, method, path, message, headerMessage)
				ac.Write([]byte(output)) //nolint
			},
		}
		excludeExtensions := []string{"doc"} //获取API文档接口不记录

		c.AddSkipper(func(ctx iris.Context) bool {
			path := ctx.Path()
			for _, ext := range excludeExtensions {
				if strings.HasSuffix(path, ext) {
					return true
				}
			}
			return false
		})

		l := logger.New(c)
		s.server.Logger().SetPrefix("[RUSH HTTPD ACCESSLOG]")
		s.server.Logger().SetLevel("debug")
		s.server.Use(l)
	}

	if err := s.addNewHandler(BasePath, c, d, disc); err != nil {
		return nil, err
	}

	if err := s.addNewHandler(PathVersion2, c, d, disc); err != nil {
		return nil, err
	}

	r := Route{
		RouteType:   ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/doc",
		HandlerFunc: s.methods.getDoc,
	}
	if err := s.AddNewHttpHandler(r); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) manage() {
	stopDone := <-s.stop
	timeout := s.shutdownTimeout
	ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
	defer cancel()
	_ = s.server.Shutdown(ctx)
	close(stopDone)
}

// Close closes the underlying listener.
func (s *Service) Close() error {
	defer s.diag.StoppedService()
	// If server is not set we were never started
	if s.server == nil || !s.opened {
		return nil
	}
	// Signal to manage loop we are stopping
	stopping := make(chan struct{})
	s.stop <- stopping

	<-stopping
	s.server = nil
	s.opened = false
	return nil
}

func (s *Service) serve() {
	if s.server == nil {
		panic("Http Server Must Be Init Before Invoke [server] Method!!!!")
	}
	err := s.server.Run(s.Addr(), iris.WithoutInterruptHandler, iris.WithoutServerError(iris.ErrServerClosed))
	// The listener was closed so exit
	// See https://github.com/golang/go/issues/4373
	if !strings.Contains(err.Error(), "closed") {
		s.err <- fmt.Errorf("listener failed: addr=%s, err=%w", s.URL(), err)
	} else {
		s.err <- nil
	}
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

// Open starts the service
func (s *Service) Open() error {
	s.diag.StartingService()

	s.stop = make(chan chan struct{})

	var envAddr = utils.GetEnv("ENV_RUNTIME_ADDR", ":6060")

	go s.manage()
	go s.serve()
	go func() {
		if utils.CheckRuntimeEnvIsDev() {
			mux := http.NewServeMux()

			err := statsviz.Register(mux,
				statsviz.Root("/debug"),
				statsviz.SendFrequency(500*time.Millisecond),
			)
			if err != nil {
				s.diag.Error("statsiviz start", err)
			}
			statsSrv := &http.Server{Addr: envAddr, Handler: mux}
			s.server.NewHost(statsSrv).ListenAndServe()
		}
	}()
	s.opened = true
	return nil
}

func (s *Service) Addr() iris.Runner {
	return iris.Addr(s.addr)
}

func (s *Service) Err() <-chan error {
	return s.err
}

func (s *Service) URL() string {

	return "http://" + s.server.ConfigurationReadOnly().GetVHost()
}

// URL that should resolve externally to the server HTTP endpoint.
// It is possible that the URL does not resolve correctly  if the hostname config setting is incorrect.
func (s *Service) ExternalURL() string {
	return s.externalURL
}

func (s *Service) GetHandlerByName(version string) (*Handler, error) {
	i, ok := s.HandlerByNames[version]
	if !ok {
		// Should be unreachable code
		return nil, fmt.Errorf("cannot get handler By %s", version)
	}

	return s.Handler[i], nil
}

func (s *Service) AddNewVersion(version string) error {
	c := s.config()
	return s.addNewHandler(version, c, s.diag, s.DiagService)
}

func (s *Service) AddNewRawHttpHandler(method, endpoint string, handler context.Handler) error {
	app := s.server
	if app == nil {
		return errors.New("AddNewRawHttpHandler: Please Initial Web Server First")
	}
	r := app.Handle(method, endpoint, handler)
	if r == nil {
		return errors.New("AddNewRawHttpHandler: Add Handler Error")
	}
	return nil
}

func (s *Service) addNewHandler(version string, c Config, d Diagnostic, disc DiagService) error {
	if _, ok := s.HandlerByNames[version]; ok {
		// Should be unreachable code
		return errors.New("Cannot Append Handler Twice")
	}
	crs := cors.New(cors.Options{
		AllowedOrigins:   s.cors.AllowedOrigins,
		AllowCredentials: s.cors.AllowCredentials,
		AllowedMethods:   s.cors.AllowedMethods,
	})
	s.server.WrapRouter(crs.ServeHTTP)
	p := s.server.Party(version).AllowMethods(iris.MethodOptions)
	if p == nil {
		return fmt.Errorf("fail to create the party %s", version)
	}
	h := NewHandler(
		c.LogEnabled,
		c.WriteTracing,
		d,
	)
	h.service = s.server
	h.DiagService = disc
	h.Version = version
	h.party = &p

	i := len(s.Handler)
	s.Handler = append(s.Handler, h)

	s.HandlerByNames[version] = i

	return nil
}

func (s *Service) AddNewWebsocketHandler(r WebsocketRoute) error {
	if s.server == nil || r.Server == nil {
		return errors.New("AddNewWebsocketHandler, Http Server Is Empty!")
	}
	s.server.Get(r.Pattern, websocket.Handler(r.Server))
	return nil
}

func (s *Service) AddNewHttpHandler(r Route) error {
	h, err := s.GetHandlerByName(BasePath)
	if err != nil {
		return err
	}
	err = h.addRoute(r)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) AddNewHttpHandlerWithVersion(r Route, version string) error {
	if version == "" {
		version = BasePath
	}

	h, err := s.GetHandlerByName(version)
	if err != nil {
		return err
	}
	err = h.addRoute(r)
	if err != nil {
		return err
	}
	return nil
}
