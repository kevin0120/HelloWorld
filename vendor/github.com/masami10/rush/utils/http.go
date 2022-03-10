package utils

import (
	"encoding/json"

	"github.com/kataras/iris/v12"
)


type RushErrResp struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func ErrResponse(ctx iris.Context, respCode int, errCode string, errMsg string) {
	ctx.StatusCode(respCode)
	errResp := RushErrResp{
		ErrorCode: errCode,
		ErrorMsg:  errMsg,
	}

	body, _ := json.Marshal(errResp)
	_, _ = ctx.Write(body)
}