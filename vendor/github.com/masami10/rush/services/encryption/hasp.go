package encryption

/*
#ifdef linux
#cgo CFLAGS: -I../../libs/hasp/include/
#cgo LDFLAGS: ${SRCDIR}/../../libs/hasp/linux/libhasp_linux_x86_64_37458.a -lpthread
#endif

#ifdef _WIN32
#cgo CFLAGS: -I../../libs/hasp/include/
#cgo LDFLAGS: ${SRCDIR}/../../libs/hasp/win64/libhasp_windows_x64_37458.a
#endif

#include <stdio.h>
#include <string.h>

#include "hasp_api.h"
#include "hasp_vcode.h"

unsigned char str[1024];

hasp_status_t hasp_login_feature(int feature_id, hasp_handle_t* handle) {
	hasp_status_t   status;
	status = hasp_login(feature_id,
						(hasp_vendor_code_t)vendor_code,
						handle);
	return status;
}

hasp_status_t hasp_logout_feature(hasp_handle_t handle) {
	hasp_status_t   status;
	status = hasp_logout(handle);
	return status;
}

char* hasp_read_ro(hasp_handle_t handle, hasp_size_t offect, hasp_size_t fsize) {
	hasp_status_t   status;
    memset(str, 0, sizeof str);
	status = hasp_read(handle,
					HASP_FILEID_RO,
					offect,
					fsize,
					&str);
    if (status != HASP_STATUS_OK) {
        return "";
    }
	return str;
}
*/
import "C"
import (
	"encoding/json"
)

type AuthPayload struct {
	Controllers int `json:"controllers"` // 控制器数量
	// ...
}

func HaspLogin(featureID int, handle *C.hasp_handle_t) int {
	ret := C.hasp_login_feature(C.int(featureID), handle)
	return int(ret)
}

func HaspLogout(handle C.hasp_handle_t) int {
	ret := C.hasp_logout_feature(handle)
	return int(ret)
}

func GetAuthPayload() *AuthPayload {
	resp := &AuthPayload{
		Controllers: 0,
	}
	if IsDev {
		resp.Controllers = -1
		return resp
	}
	handle := C.hasp_handle_t(0)

	ret := HaspLogin(DefaultFeatureId, &handle)
	if ret == 0 {
		defer HaspLogout(handle)
		ret := C.hasp_read_ro(
			handle,
			C.hasp_size_t(controllersNumberOffset[0]),
			C.hasp_size_t(controllersNumberOffset[1]),
		)
		// 如果失败了ret为空 这里解析默认值
		_ = json.Unmarshal([]byte(C.GoString(ret)), resp)
		return resp
	}
	return resp
}

func GetPurviewFromFeature() bool {
	// dev mod always true
	if IsDev {
		return true
	}
	handle := C.hasp_handle_t(0)
	ret := HaspLogin(DefaultFeatureId, &handle)
	if ret == 0 {
		defer HaspLogout(handle)
		return true
	}
	return false
}
