package tightening_device

import (
	"context"
	"fmt"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"strings"
)

var ENV_PSET_BATCH_ENABLE = utils.GetEnvBool("ENV_PSET_BATCH_ENABLE", true)
var ENV_PSET_WITH_ENABLE = utils.GetEnvBool("ENV_PSET_WITH_ENABLE", true)

type JobSet struct {
	TighteningUnit string `json:"tightening_unit"`
	ControllerSN   string `json:"controller_sn"`
	ToolSN         string `json:"tool_sn"`
	WorkorderID    int64  `json:"workorder_id"`
	Total          int    `json:"total"`
	StepID         int64  `json:"step_id"`
	Job            int    `json:"job"`
	UserID         int64  `json:"user_id"`
}

func (s *JobSet) Validate() error {
	if s.ControllerSN == "" {
		return errors.New("Controller SerialNumber is required")
	}

	if s.Job <= 0 {
		return errors.New("Job Should Be Greater Than Zero")
	}

	return nil
}

type PSetBatchSet struct {
	ControllerSN   string `json:"controller_sn" validate:"required"`
	TighteningUnit string `json:"tightening_unit"`
	ToolSN         string `json:"tool_sn" validate:"required"`
	PSet           int    `json:"pset" validate:"required"`
	Batch          int    `json:"batch" validate:"required"`
}

//todo: 接口中哪些字段是的验证条件还不确认
type PSetSet struct {
	ToolControl
	StepCode      string `json:"workstep_code" validate:"-"`
	WorkorderID   int64  `json:"workorder_id" validate:"-"`
	WorkorderCode string `json:"workorder_code" validate:"-"`
	UserID        int64  `json:"user_id" validate:"gt=0"`
	PSet          int    `json:"pset" validate:"gt=0"`
	Sequence      uint   `json:"sequence" validate:"gt=0"`
	Count         int    `json:"count" validate:"-"` //拧紧结果计数，发送请求时应该不传递
	Batch         int    `json:"batch" validate:"-"`
	Total         int    `json:"total" validate:"-"`
	IP            string `json:"ip" validate:"-"`
	PointID       string `json:"point_id" validate:"-"`
	ScannerCode   string `json:"scanner_code" validate:"-"`
}

type ToolControl struct {
	TighteningUnit string `json:"tightening_unit"`
	ControllerSN   string `json:"controller_sn" validate:"required"`
	ToolSN         string `json:"tool_sn" validate:"required"`
	Enable         bool   `json:"enable" validate:"required"`
}

type ToolModeSelect struct {
	TighteningUnit string `json:"tightening_unit"`
	ControllerSN   string `json:"controller_sn"`
	ToolSN         string `json:"tool_sn"`
	Mode           string `json:"mode"`
}

func (s *ToolModeSelect) Validate() error {
	if s.ControllerSN == "" || (s.ToolSN == "" && s.TighteningUnit == "") {
		return errors.New("Controller SerialNumber or Tool SerialNumber is required")
	}

	return nil
}

func (s *Service) ToolControl(req *ToolControl) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	if req.TighteningUnit == "" && req.ToolSN != "" {
		req.TighteningUnit = req.ToolSN
	}

	tool, err := s.getTool(req.ControllerSN, req.TighteningUnit)
	if err != nil {
		return err
	}

	return tool.ToolControl(req.Enable)
}

func (s *Service) ToolJobSet(req *JobSet, update bool) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	err := req.Validate()
	if err != nil {
		return err
	}

	if req.TighteningUnit == "" {
		req.TighteningUnit = req.ToolSN
	}

	tool, err := s.getTool(req.ControllerSN, req.TighteningUnit)
	if err != nil {
		return err
	}

	if req.UserID == 0 {
		req.UserID = 1
	}

	if update {
		_ = s.storageService.UpdateTool(&storage.Tools{
			Serial:             req.TighteningUnit,
			CurrentWorkorderID: req.WorkorderID,
			Total:              req.Total,
			UserID:             req.UserID,
		})
	}
	return tool.SetJob(req.Job)
}

func (s *Service) ToolPSetBatchSet(req *PSetBatchSet) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	if req.PSet == 0 {
		return errors.New("ToolPSetBatchSet Pset Must Be Greater Than Zero")
	}

	if req.TighteningUnit == "" && req.ToolSN != "" {
		req.TighteningUnit = req.ToolSN
	}

	tool, err := s.getTool(req.ControllerSN, req.TighteningUnit)
	if err != nil {
		return err
	}

	return tool.SetPSetBatch(req.PSet, req.Batch)
}

func (s *Service) doTraceFromPSetReq(req *PSetSet) {
	var orderid int64
	var stepid int64

	if req.UserID == 0 {
		req.UserID = 1
	}

	if req.WorkorderCode != "" {
		order, err := s.storageService.GetWorkorderByCode(req.WorkorderCode)
		if err != nil {
			s.diag.Error("doTraceFromPSetReq.GetWorkorderByCode Failed", err)
		} else {
			orderid = order.Id
		}
	}

	if req.StepCode != "" {
		step, err := s.storageService.GetStepByCodeAndWorkorderID(req.StepCode, orderid)
		if err != nil {
			s.diag.Error("doTraceFromPSetReq.GetStepByCodeAndWorkorderID Failed", err)
		} else {
			stepid = step.Id
		}
	}

	_ = s.storageService.UpdateTool(&storage.Tools{
		Serial:             req.TighteningUnit,
		CurrentWorkorderID: orderid,
		Seq:                int(req.Sequence),
		Count:              req.Count,
		UserID:             req.UserID,
		Total:              req.Total,
		StepID:             stepid,
		ScannerCode:        req.ScannerCode,
		Batch:              req.Batch,
	})
}

func (s *Service) ToolPSetByIP(req *PSetSet) error {
	if req == nil {
		return errors.New("ToolPSetByIP Req Is Nil")
	}

	tool, err := s.findToolbyIP(req.IP)
	if err != nil {
		s.diag.Error("doTraceFromPSetReq.GetStepByCodeAndWorkorderID Failed", err)
		return err
	}

	if req.UserID == 0 {
		req.UserID = 1
	}

	controller := tool.GetParentService().(ITighteningController)
	if err = s.ToolPSetBatchSet(&PSetBatchSet{
		ControllerSN: controller.SerialNumber(),
		ToolSN:       tool.SerialNumber(),
		PSet:         req.PSet,
		Batch:        1,
	}); err != nil {
		s.diag.Error("PSet Batch Set Failed", err)
	}

	err = s.storageService.UpdateTool(&storage.Tools{
		Serial:  tool.SerialNumber(),
		Seq:     int(req.Sequence),
		Count:   req.Count,
		UserID:  req.UserID,
		Total:   req.Total,
		PointID: req.PointID,
	})
	if err != nil {
		s.diag.Error("UpdateTool Error ", err)
		return err
	}
	return nil
}

func (s *Service) ToolPSetSet(req *PSetSet) error {

	dopset := func(tool ITighteningTool) error {
		ctx := context.WithValue(context.Background(), "psetReq", req)
		return tool.SetPSet(ctx, req.PSet)
	}

	if req == nil {
		return errors.New("Req Is Nil")
	}

	if req.TighteningUnit == "" && req.ToolSN != "" {
		req.TighteningUnit = req.ToolSN
	}

	s.doTraceFromPSetReq(req)

	controller, err := s.getController(req.ControllerSN)
	if err != nil {
		return err
	}

	tool, err := s.getTool(req.ControllerSN, req.TighteningUnit)
	if err != nil {
		return err
	}

	if req.PSet <= 0 {
		return errors.New("ToolPSetSet.ToolPSetBatchSet Pset Must Be Greater Than 0!!!")
	}

	if ENV_PSET_BATCH_ENABLE {
		err := s.ToolPSetBatchSet(&PSetBatchSet{
			ControllerSN:   req.ControllerSN,
			ToolSN:         req.ToolSN,
			TighteningUnit: req.TighteningUnit,
			PSet:           req.PSet,
			Batch:          req.Batch,
		})
		if err != nil {
			if !strings.Contains(err.Error(), TIGHTENING_ERR_NOT_SUPPORTED) {
				s.diag.Error("ToolPSetSet.ToolPSetBatchSet", err)
				return err
			} else {
				if err := dopset(tool); err != nil {
					return err
				}
			}
			return nil
		}
		if controller.Model() != ModelLexenWrench {
			if err := dopset(tool); err != nil {
				return err
			}
		}
	} else {
		if err := dopset(tool); err != nil {
			return err
		}
	}

	if req.UserID <= 0 {
		req.UserID = 1
	}

	//FIXME:手动模式下sequence可能为0
	// if req.Sequence <= 0 {
	// 	err := errors.New("Sequence Is Less Than Zero")
	// 	s.diag.Error("ToolPSetSet", err)
	// 	return err
	// }

	if req.Enable || ENV_PSET_WITH_ENABLE {
		_ = tool.ToolControl(true)
	}

	s.diag.Info(fmt.Sprintf("Pset Request Pset Number: %d Success!!!", req.PSet))

	return nil
}

func (s *Service) ToolModeSelect(req *ToolModeSelect) error {
	if req == nil {
		return errors.New("Req Is Nil")
	}

	err := req.Validate()
	if err != nil {
		return err
	}
	if req.TighteningUnit == "" {
		req.TighteningUnit = req.ToolSN
	}

	tool, err := s.getTool(req.ControllerSN, req.TighteningUnit)
	if err != nil {
		return err
	}

	return tool.ModeSelect(req.Mode)
}

type ToolInfo struct {
	TighteningUnit string `json:"tightening_unit"`
	ControllerSN   string `json:"controller_sn"`
	ToolSN         string `json:"tool_sn"`
}

type ToolPSet struct {
	ToolInfo
	PSet int `json:"pset"`
}

type ToolJob struct {
	ToolInfo
	Job int `json:"job"`
}

func (s *Service) GetToolPSetList(req *ToolInfo) ([]PSetInfo, error) {
	if req == nil {
		return nil, errors.New("Req Is Nil")
	}
	if req.TighteningUnit == "" && req.ToolSN != "" {
		req.TighteningUnit = req.ToolSN
	}

	tool, err := s.getTool(req.ControllerSN, req.TighteningUnit)
	if err != nil {
		return nil, err
	}

	return tool.GetPSetList()
}

func (s *Service) GetToolPSetDetail(req *ToolPSet) (*PSetDetail, error) {
	if req == nil {
		return nil, errors.New("Req Is Nil")
	}

	if req.TighteningUnit == "" && req.ToolSN != "" {
		req.TighteningUnit = req.ToolSN
	}

	tool, err := s.getTool(req.ControllerSN, req.TighteningUnit)
	if err != nil {
		return nil, err
	}

	return tool.GetPSetDetail(req.PSet)
}

func (s *Service) GetToolJobList(req *ToolInfo) ([]int, error) {
	if req == nil {
		return nil, errors.New("Req Is Nil")
	}

	if req.TighteningUnit == "" && req.ToolSN != "" {
		req.TighteningUnit = req.ToolSN
	}

	tool, err := s.getTool(req.ControllerSN, req.TighteningUnit)
	if err != nil {
		return nil, err
	}

	return tool.GetJobList()
}

func (s *Service) GetToolJobDetail(req *ToolJob) (*JobDetail, error) {
	if req == nil {
		return nil, errors.New("Req Is Nil")
	}

	if req.TighteningUnit == "" && req.ToolSN != "" {
		req.TighteningUnit = req.ToolSN
	}

	tool, err := s.getTool(req.ControllerSN, req.TighteningUnit)
	if err != nil {
		return nil, err
	}

	return tool.GetJobDetail(req.Job)
}

func (s *Service) findToolbyIP(ip string) (ITighteningTool, error) {
	for _, controller := range s.runningControllers {
		tool, err := controller.GetToolViaIP(ip)
		if err == nil {
			return tool, nil
		}
	}

	return nil, errors.New("findToolbyIP: Not Found")
}
