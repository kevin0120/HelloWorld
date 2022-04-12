package utils

import (
	"crypto/tls"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func DropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func GetEnvInt32(key string, fallback int32) int32 {
	ret := fallback
	value, exists := os.LookupEnv(key)
	if !exists {
		return ret
	}
	if t, err := strconv.Atoi(value); err != nil { //nolint:gosec
		return ret
	} else {
		ret = int32(t)
	}
	return ret
}

func GetEnvInt(key string, fallback int) int {
	ret := fallback
	value, exists := os.LookupEnv(key)
	if !exists {
		return ret
	}
	if t, err := strconv.Atoi(value); err != nil { //nolint:gosec
		return ret
	} else {
		ret = t
	}
	return ret
}

func GetEnvBool(key string, fallback bool) bool {
	ret := fallback
	value, exists := os.LookupEnv(key)
	if !exists {
		return ret
	}
	if t, err := strconv.ParseBool(value); err != nil {
		return ret
	} else {
		ret = t
	}
	return ret
}

func CheckRuntimeEnvIsDev() bool {
	var env = GetEnv("ENV_RUNTIME", "dev")

	if env == "dev" || env == "development" {
		return true
	}
	return false
}

func FileIsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func CreateRetryClient(timeout time.Duration, maxRetry int, diag Diagnostic) *resty.Client {

	r := resty.New()
	if diag != nil {
		r.SetLogger(diag)
	}

	r.SetTimeout(timeout)
	r.SetContentLength(true)
	r.
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Brand":        "Oneshare",
		}).
		SetRetryCount(maxRetry).
		SetRetryWaitTime(time.Second * 2).
		SetRetryMaxWaitTime(20 * time.Second).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) //nolint:gosec

	return r
}

func AppendByteSlice(s []byte, t []byte) []byte {
	zlen := len(s) + len(t)
	z := make([]byte, zlen)
	copy(z, s)
	copy(z[len(s):], t)
	return z
}

func GenerateID() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}

func GetDateTime() (string, string) {
	stime := strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")
	return stime[0], stime[1]
}

func RushRound(x, unit float64) float64 {
	return float64(int64(x/unit+0.5)) * unit
}

func ReverseString(raw string) string {
	rt := ""
	for _, v := range raw {
		rt = string(v) + rt
	}
	return rt
}

func ArrayContains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}

	return false
}

func WaitGroupTimeout(wg *sync.WaitGroup, timeout time.Duration) error {
	if wg == nil {
		return errors.New("Wait Group Is Nil")
	}

	wg.Add(1)
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return nil // completed normally
	case <-time.After(timeout):
		return errors.New("Timeout") // timed out
	}
}
