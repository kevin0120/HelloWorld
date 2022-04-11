package tightening_device

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/masami10/rush/services/httpd"
)

var validate *validator.Validate = nil

func init() {
	validate = validator.New()
}

func (s *Service) putToolEnable(ctx iris.Context) {
	var req ToolControl
	resp := HMICommonResponse{StatusCode: iris.StatusBadRequest}

	if err := ctx.ReadJSON(&req); err != nil {
		resp.Message = "putToolEnable Read JSON Error"
		resp.Extra = err.Error()
		_ = httpd.NewCommonResponseBody(&resp, ctx)
		return
	}

	if err := validate.Struct(req); err != nil {
		resp.Message = "putToolEnable Validate Error"
		resp.Extra = err.Error()
		_ = httpd.NewCommonResponseBody(&resp, ctx)
	}

	if err := s.ToolControl(&req); err != nil {
		resp.Message = "putToolEnable ToolControl Error"
		resp.Extra = err.Error()
		_ = httpd.NewCommonResponseBody(&resp, ctx)
		return
	}

	ctx.StatusCode(iris.StatusOK)
}

func (s *Service) putToolPSet(ctx iris.Context) {
	var req PSetSet
	resp := HMICommonResponse{StatusCode: iris.StatusBadRequest}

	if err := ctx.ReadJSON(&req); err != nil {
		resp.Message = "putToolPSet Read JSON Error"
		resp.Extra = err.Error()
		_ = httpd.NewCommonResponseBody(&resp, ctx)
		return
	}

	if err := validate.Struct(req); err != nil {
		resp.Message = "putToolPSet Validate Error"
		resp.Extra = err.Error()
		_ = httpd.NewCommonResponseBody(&resp, ctx)
	}

	if err := s.ToolPSetSet(&req); err != nil {
		resp.Message = "putToolPSet ToolPSetSet Error"
		resp.Extra = err.Error()
		_ = httpd.NewCommonResponseBody(&resp, ctx)
		return
	}

	ctx.StatusCode(iris.StatusOK)
}
