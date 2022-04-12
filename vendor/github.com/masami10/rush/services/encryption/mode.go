// +build release

package encryption

// 只有在release模式下启动鉴权
// go build -tags release
func init() {
	IsDev = false
}
