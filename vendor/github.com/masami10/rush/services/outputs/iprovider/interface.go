package iprovider

type Provider interface {
	SetSerializer(serializer Serializer)
	Connect() error
	Close() error
	Write(pkg []PublishPackage) error
}

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
	Info(msg string)
}

type Serializer interface {
	Serialize(data []PublishPackage) ([]FileItem, error)
	ContentType() string
}
