package encryption

import "C"
import "time"

const (
	CheckAuthInterval = 30 * time.Minute
	DefaultFeatureId  = 2
)

// IsDev 用来改版编译参数
var IsDev = true

var (
	controllersNumberOffset = []uint{0, 1024}
)

func CheckAuthorityTicker() {
	for {
		<-time.After(CheckAuthInterval)
		if !GetPurviewFromFeature() {
			panic("程序未授权!")
		}
	}
}
