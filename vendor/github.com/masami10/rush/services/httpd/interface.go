package httpd

type DiagService interface {
	SetLogLevelFromName(lvl string) error
}
