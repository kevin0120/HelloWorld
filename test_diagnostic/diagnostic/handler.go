package diagnostic

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"
)

type logLevel int

const (
	llInvalid logLevel = iota
	llDebug
	llError
	llInfo
)

type StaticLevelHandler struct {
	l     Logger
	level logLevel
}

func (h *StaticLevelHandler) Write(buf []byte) (int, error) {
	switch h.level {
	case llDebug:
		h.l.Debug(string(buf))
	case llError:
		h.l.Error(string(buf))
	case llInfo:
		h.l.Info(string(buf))
	default:
		return 0, errors.New("invalid log level")
	}

	return len(buf), nil
}

func Err(l Logger, msg string, err error, ctx []T) {
	if len(ctx) == 0 {
		l.Error(msg, Error(err))
		return
	}

	if len(ctx) == 1 {
		el := ctx[0]
		l.Error(msg, Error(err), String(el.Key, el.Value))
		return
	}

	if len(ctx) == 2 {
		x := ctx[0]
		y := ctx[1]
		l.Error(msg, Error(err), String(x.Key, x.Value), String(y.Key, y.Value))
		return
	}

	// Use the allocation version for any length
	fields := make([]Field, len(ctx)+1) // +1 for error
	fields[0] = Error(err)
	for i := 1; i < len(fields); i++ {
		kv := ctx[i-1]
		fields[i] = String(kv.Key, kv.Value)
	}

	l.Error(msg, fields...)
}

func Info(l Logger, msg string, ctx []T) {
	if len(ctx) == 0 {
		l.Info(msg)
		return
	}

	if len(ctx) == 1 {
		el := ctx[0]
		l.Info(msg, String(el.Key, el.Value))
		return
	}

	if len(ctx) == 2 {
		x := ctx[0]
		y := ctx[1]
		l.Info(msg, String(x.Key, x.Value), String(y.Key, y.Value))
		return
	}

	// Use the allocation version for any length
	fields := make([]Field, len(ctx))
	for i, kv := range ctx {
		fields[i] = String(kv.Key, kv.Value)
	}

	l.Info(msg, fields...)
}

func Debug(l Logger, msg string, ctx []T) {
	if len(ctx) == 0 {
		l.Debug(msg)
		return
	}

	if len(ctx) == 1 {
		el := ctx[0]
		l.Debug(msg, String(el.Key, el.Value))
		return
	}

	if len(ctx) == 2 {
		x := ctx[0]
		y := ctx[1]
		l.Debug(msg, String(x.Key, x.Value), String(y.Key, y.Value))
		return
	}

	// Use the allocation version for any length
	fields := make([]Field, len(ctx))
	for i, kv := range ctx {
		fields[i] = String(kv.Key, kv.Value)
	}

	l.Debug(msg, fields...)
}

// Cmd handler

type CmdHandler struct {
	l Logger
}

func (h *CmdHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *CmdHandler) RushStarting(version, branch, commit string) {
	h.l.Info("rush starting", String("version", version), String("branch", branch), String("commit", commit))
}

func (h *CmdHandler) GoVersion() {
	h.l.Info("go version", String("version", runtime.Version()))
}

func (h *CmdHandler) Info(msg string) {
	h.l.Info(msg)
}

type ServerHandler struct {
	l Logger
}

func (h *ServerHandler) Error(msg string, err error, ctx ...T) {
	Err(h.l, msg, err, ctx)
}

func (h *ServerHandler) Info(msg string, ctx ...T) {
	Info(h.l, msg, ctx)
}

func (h *ServerHandler) Debug(msg string, ctx ...T) {
	Debug(h.l, msg, ctx)
}

// HTTPD handler

type HTTPDHandler struct {
	l Logger
}

func (h *HTTPDHandler) NewHTTPServerErrorLogger() *log.Logger {
	s := &StaticLevelHandler{
		l:     h.l.With(String("service", "httpd_server_errors")),
		level: llError,
	}

	return log.New(s, "", log.LstdFlags)
}

func (h *HTTPDHandler) StartingService() {
	h.l.Info("starting HTTP service")
}

func (h *HTTPDHandler) StoppedService() {
	h.l.Info("closed HTTP service")
}

func (h *HTTPDHandler) ShutdownTimeout() {
	h.l.Error("shutdown timedout, forcefully closing all remaining connections")
}

func (h *HTTPDHandler) AuthenticationEnabled(enabled bool) {
	h.l.Info("authentication", Bool("enabled", enabled))
}

func (h *HTTPDHandler) ListeningOn(addr string, proto string) {
	h.l.Info("listening on", String("addr", addr), String("protocol", proto))
}

func (h *HTTPDHandler) WriteBodyReceived(body string) {
	h.l.Debug("write body received by handler: %s", String("body", body))
}

func (h *HTTPDHandler) HTTP(
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
) {
	h.l.Info("http request",
		String("host", host),
		String("username", username),
		Time("start", start),
		String("method", method),
		String("uri", uri),
		String("protocol", proto),
		Int("status", status),
		String("referer", referer),
		String("user-agent", userAgent),
		String("request-id", reqID),
		Duration("duration", duration),
	)
}

func (h *HTTPDHandler) RecoveryError(
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
) {
	h.l.Error(
		msg,
		String("err", err),
		String("host", host),
		String("username", username),
		Time("start", start),
		String("method", method),
		String("uri", uri),
		String("protocol", proto),
		Int("status", status),
		String("referer", referer),
		String("user-agent", userAgent),
		String("request-id", reqID),
		Duration("duration", duration),
	)
}

func (h *HTTPDHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

// Aiis Handler

type AiisHandler struct {
	l Logger
}

func (h *AiisHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *AiisHandler) Debug(msg string) {
	h.l.Debug(msg)
}

func (h *AiisHandler) Info(msg string) {
	h.l.Info(msg)
}

func (h *AiisHandler) PutResultDone() {
	h.l.Debug("Put Result to AIIS successful")
}

// Odoo Handler

type OdooHandler struct {
	l Logger
}

func (h *OdooHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *OdooHandler) Info(msg string) {
	h.l.Info(msg)
}

func (h *OdooHandler) Debug(msg string) {
	h.l.Debug(msg)
}

func (h *OdooHandler) CreateWOSuccess(id int64) {
	h.l.Debug(fmt.Sprintf("Create WO successful %d", id))
}

// Minio Handler

type MinioHandler struct {
	l Logger
}

func (h *MinioHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *MinioHandler) Debug(msg string) {
	h.l.Debug(msg)
}

// hmi handler
type HmiHandler struct {
	l Logger
}

func (h *HmiHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *HmiHandler) Debug(msg string) {
	h.l.Debug(msg)
}

func (h *HmiHandler) Close() {
	h.l.Info("hmi server closing")
}

func (h *HmiHandler) Closed() {
	h.l.Info("hmi server closed")
}

func (h *HmiHandler) Disconnect(id string) {
	h.l.Info("hmi Connection disconnected", String("ID", id))
}

// Websocket Handler

type WsHandler struct {
	l Logger
}

func (h *WsHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *WsHandler) Disconnect(id string) {
	h.l.Info("ws Connection disconnected", String("ID", id))
}

func (h *WsHandler) Close() {
	h.l.Info("ws server closing")
}

func (h *WsHandler) OnMessage(msg string) {
	h.l.Debug(msg)
}

func (h *WsHandler) Closed() {
	h.l.Info("ws server closed")
}

type StorageHandler struct {
	l Logger
}

func (h *StorageHandler) Info(msg string) {
	h.l.Info(msg)
}

func (h *StorageHandler) Debug(msg string) {
	h.l.Debug(msg)
}

func (h *StorageHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *StorageHandler) OpenEngineSuccess(msg string) {
	h.l.Info("Open Engine success", String("T", msg))
}

func (h *StorageHandler) Close() {
	h.l.Info("Storage Service closing")
}

func (h *StorageHandler) Closed() {
	h.l.Info("Storage Service closed")
}

//Audi / VW
type AudiVWHandler struct {
	l Logger
}

func (h *AudiVWHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *AudiVWHandler) Debug(msg string) {
	h.l.Debug(msg)
}

func (h *AudiVWHandler) Info(msg string) {
	h.l.Info(msg)
}

func (h *AudiVWHandler) StartManager() {
	h.l.Info("start Manage for write")
}

//Openprotocol
type OpenProtocolHandler struct {
	l Logger
}

func (h *OpenProtocolHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *OpenProtocolHandler) Debug(msg string) {
	h.l.Debug(msg)
}

func (h *OpenProtocolHandler) Info(msg string) {
	h.l.Info(msg)
}

// Controller
type ControllerHandler struct {
	l Logger
}

func (h *ControllerHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *ControllerHandler) Debug(msg string) {
	h.l.Debug(msg)
}

type ScannerHandler struct {
	l Logger
}

func (h *ScannerHandler) Info(msg string) {
	h.l.Info(msg)
}

func (h *ScannerHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *ScannerHandler) Debug(msg string) {
	h.l.Debug(msg)
}

type BrokerHandler struct {
	l Logger
}

func (h *BrokerHandler) Info(msg string) {
	h.l.Info(msg)
}

func (h *BrokerHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *BrokerHandler) Debug(msg string) {
	h.l.Debug(msg)
}

type IOHandler struct {
	l Logger
}

func (h *IOHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *IOHandler) Debug(msg string) {
	h.l.Debug(msg)
}

type ReaderHandler struct {
	l Logger
}

func (h *ReaderHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *ReaderHandler) Debug(msg string) {
	h.l.Debug(msg)
}

type TighteningDeviceHandler struct {
	l Logger
}

func (h *TighteningDeviceHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *TighteningDeviceHandler) Debug(msg string) {
	h.l.Debug(msg)
}

type DeviceHandler struct {
	l Logger
}

func (h *DeviceHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *DeviceHandler) Debug(msg string) {
	h.l.Debug(msg)
}

// 客制化项目模块使用的句柄
type CustomizeHandler struct {
	l Logger
}

func (h *CustomizeHandler) Info(msg string) {
	h.l.Info(msg)
}

func (h *CustomizeHandler) Error(msg string, err error) {
	h.l.Error(msg, Error(err))
}

func (h *CustomizeHandler) Debug(msg string) {
	h.l.Debug(msg)
}
