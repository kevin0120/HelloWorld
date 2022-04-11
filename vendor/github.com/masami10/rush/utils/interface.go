package utils

// Service represents a service attached to the server.
type ICommonService interface {
	Open() error
	Close() error
}

type Diagnostic interface {
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}
