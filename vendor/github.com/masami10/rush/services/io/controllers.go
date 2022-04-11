package io

import (
	"github.com/kataras/iris/v12"
	"github.com/masami10/rush/services/httpd"
)

func (s *Service) putIOSet(ctx iris.Context) {
	req := IoSet{}
	resp := HMICommonResponse{StatusCode: iris.StatusBadRequest}

	if err := ctx.ReadJSON(&req); err != nil {
		resp.Message = "putIOSet Read JSON Error"
		resp.Extra = err.Error()
		_ = httpd.NewCommonResponseBody(&resp, ctx)
		return
	}

	if err := s.Write(req.SN, req.Index, req.Status); err != nil {
		resp.Message = "putIOSet Write IO Data Error"
		resp.Extra = err.Error()
		_ = httpd.NewCommonResponseBody(&resp, ctx)
		return
	}

	ctx.StatusCode(iris.StatusOK)
}
