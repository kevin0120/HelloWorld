package diagnostic

import (
	"errors"
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type nopCloser struct {
	f io.Writer
}

func (c *nopCloser) Write(b []byte) (int, error) { return c.f.Write(b) }
func (c *nopCloser) Close() error                { return nil }

type Service struct {
	c Config

	Logger Logger

	f      io.WriteCloser
	stdout io.Writer
	stderr io.Writer

	//SessionService *SessionService

	levelMu sync.RWMutex
	level   string
}

func NewService(c Config, stdout, stderr io.Writer) *Service {
	return &Service{
		c:      c,
		stdout: stdout,
		stderr: stderr,
	}
}

func (s *Service) NewServerHandler() *ServerHandler {
	return &ServerHandler{
		l: s.Logger.With(String("source", "srv")),
	}
}

func (s *Service) NewCmdHandler() *CmdHandler {
	return &CmdHandler{
		l: s.Logger.With(String("service", "run")),
	}
}

func (s *Service) NewHTTPDHandler() *HTTPDHandler {
	return &HTTPDHandler{
		l: s.Logger.With(String("service", "http")),
	}
}

func (s *Service) NewAudiVWHandler() *AudiVWHandler {
	return &AudiVWHandler{
		l: s.Logger.With(String("service", "AudiVW")),
	}
}

func (s *Service) NewCustomizeHandler(projectCode string) *CustomizeHandler {
	return &CustomizeHandler{
		l: s.Logger.With(String("service", projectCode)),
	}
}

func (s *Service) NewOpenProtocolHandler() *OpenProtocolHandler {
	return &OpenProtocolHandler{
		l: s.Logger.With(String("service", "OpenProtocol")),
	}
}

func (s *Service) NewControllerHandler() *ControllerHandler {
	return &ControllerHandler{
		l: s.Logger.With(String("service", "Controller")),
	}
}

func (s *Service) NewAiisHandler() *AiisHandler {
	return &AiisHandler{
		l: s.Logger.With(String("service", "aiis")),
	}
}

func (s *Service) NewOdooHandler() *OdooHandler {
	return &OdooHandler{
		l: s.Logger.With(String("service", "odoo")),
	}
}

func (s *Service) NewMinioHandler() *MinioHandler {
	return &MinioHandler{
		l: s.Logger.With(String("service", "minio")),
	}
}

func (s *Service) NewWebsocketHandler() *WsHandler {
	return &WsHandler{
		l: s.Logger.With(String("service", "websocket")),
	}
}

func (s *Service) NewHMIHandler() *HmiHandler {
	return &HmiHandler{
		l: s.Logger.With(String("service", "hmi")),
	}
}

func (s *Service) NewStorageHandler() *StorageHandler {
	return &StorageHandler{
		l: s.Logger.With(String("service", "storage")),
	}
}

func (s *Service) NewScannerHandler() *ScannerHandler {
	return &ScannerHandler{
		l: s.Logger.With(String("service", "scanner")),
	}
}

func (s *Service) NewBrokerHandler() *BrokerHandler {
	return &BrokerHandler{
		l: s.Logger.With(String("service", "broker")),
	}
}

func (s *Service) NewIOHandler() *IOHandler {
	return &IOHandler{
		l: s.Logger.With(String("service", "io")),
	}
}

func (s *Service) NewReaderHandler() *ReaderHandler {
	return &ReaderHandler{
		l: s.Logger.With(String("service", "reader")),
	}
}

func (s *Service) NewTighteningDeviceHandler() *TighteningDeviceHandler {
	return &TighteningDeviceHandler{
		l: s.Logger.With(String("service", "tightening_device")),
	}
}

func (s *Service) NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{
		l: s.Logger.With(String("service", "device")),
	}
}

func BootstrapMainHandler() *CmdHandler {
	s := NewService(NewConfig(), nil, os.Stderr)
	// Should never error
	_ = s.Open()

	return s.NewCmdHandler()
}

func logLevelFromName(lvl string) Level {
	var level Level
	switch lvl {
	case "INFO", "info":
		level = InfoLevel
	case "ERROR", "error":
		level = ErrorLevel
	case "DEBUG", "debug":
		level = DebugLevel
	}

	return level
}

func (s *Service) Open() error {
	s.levelMu.Lock()
	s.level = s.c.Level
	s.levelMu.Unlock()

	levelF := func(lvl Level) bool {
		s.levelMu.RLock()
		defer s.levelMu.RUnlock()
		return lvl >= logLevelFromName(s.level)
	}

	switch s.c.File {
	case "STDERR":
		s.f = &nopCloser{f: s.stderr}
	case "STDOUT":
		s.f = &nopCloser{f: s.stdout}
	default:
		dir := path.Dir(s.c.File)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
		}
		//
		//f, err := os.OpenFile(s.c.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
		//if err != nil {
		//	return err
		//}
		rl, err := rotatelogs.New(
			fmt.Sprintf("%s", s.c.File),
			rotatelogs.WithMaxAge(time.Duration(s.c.MaxAge)),
			rotatelogs.WithRotationTime(time.Duration(s.c.Rotate)))
		if err != nil {
			return err
		}

		s.f = rl
	}

	l := NewServerLogger(s.f)
	l.SetLevelF(levelF)

	s.Logger = NewMultiLogger(
		l,
	)

	return nil
}

func (s *Service) Close() error {
	if s.f != nil {
		return s.f.Close()
	}
	return nil
}

func (s *Service) SetLogLevelFromName(lvl string) error {
	s.levelMu.Lock()
	defer s.levelMu.Unlock()
	level := strings.ToUpper(lvl)
	switch level {
	case "INFO", "ERROR", "DEBUG":
		s.level = level
	default:
		return errors.New("invalid log level")
	}

	return nil
}
