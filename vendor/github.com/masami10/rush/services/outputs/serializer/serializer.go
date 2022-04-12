package serializer

import (
	"encoding/json"
	"github.com/gocarina/gocsv"
	"github.com/masami10/rush/services/outputs/iprovider"
	"strconv"
	"time"
)

const (
	JsonProviderFormat = "json"
	CsvProviderFormat  = "csv"
)

type Diagnostic = iprovider.Diagnostic
type PublishPackage = iprovider.PublishPackage
type Serializer = iprovider.Serializer

/*
所有的Serializer结构体放在这里
*/

type JsonSerializer struct {
	diag   Diagnostic
	suffix string
}

func (js *JsonSerializer) Serialize(data []PublishPackage) ([]iprovider.FileItem, error) {
	var ret []iprovider.FileItem
	for _, pkg := range data {
		fileName := ""
		if pkg.Conf.Result.TighteningID != "" {
			fileName = pkg.Conf.Result.TighteningID
		} else {
			fileName = strconv.Itoa(int(time.Now().Unix()))
		}
		fileName = fileName + js.suffix
		d, err := json.Marshal(pkg)
		if err != nil {
			js.diag.Error("error when Serialize pkg", err)
		}
		ret = append(ret, iprovider.FileItem{
			Filename: fileName,
			Data:     d,
		})
	}
	return ret, nil
}

func (js *JsonSerializer) ContentType() string {
	return "application/json"
}

type CsvSerializer struct {
	diag   Diagnostic
	suffix string
}

func (cs *CsvSerializer) Serialize(data []PublishPackage) ([]iprovider.FileItem, error) {
	var ret []iprovider.FileItem
	fileName := strconv.Itoa(int(time.Now().Unix())) + cs.suffix

	csvContent, err := gocsv.MarshalBytes(&data)
	if err != nil {
		return nil, err
	}
	ret = append(ret, iprovider.FileItem{
		Filename: fileName,
		Data:     csvContent,
	})
	return ret, nil
}

func (cs *CsvSerializer) ContentType() string {
	return "application/csv"
}

func NewSerializer(t string, diag Diagnostic) Serializer {
	switch t {
	case JsonProviderFormat:
		return &JsonSerializer{diag: diag, suffix: ".json"}
	case CsvProviderFormat:
		return &CsvSerializer{diag: diag, suffix: ".csv"}
	default:
		return &JsonSerializer{diag: diag, suffix: ".json"}
	}
}
