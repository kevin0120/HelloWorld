package tightening_device

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/masami10/rush/services/httpd"

	"github.com/masami10/rush/services/storage"
)

type HMICommonResponse = httpd.HMICommonResponse

const (
	TIGHTENING_OPENPROTOCOL = "OpenProtocol"
	TIGHTENING_AUDIVW       = "Audi/VW"
	TIGHTENING_CVIMONITOR   = "CVIMonitor"
	TIGHTENING_CVINETWEB    = "CVINetWeb"
)

const (
	TIGHTENING_ERR_NOT_SUPPORTED = "Not Supported "
	TIGHTENING_ERR_UNKNOWN       = "Error Unknown "
	TIGHTENING_ERR_TIMEOUT       = "Timeout "

	TIGHTENING_CONTROLLER_IO_SN_FORMAT = "%s_io"
)

const (
	RESULT_OK  = "OK"
	RESULT_NOK = "NOK"
	RESULT_LSN = "LSN"
	RESULT_AK2 = "AK2"
)

const (
	STRATEGY_AD  = "AD"
	STRATEGY_AW  = "AW"
	STRATEGY_ADW = "ADW"
	STRATEGY_LN  = "LN"
	STRATEGY_AK2 = "AK2"
)

// type
const (
	TIGHTENING_DEVICE_TYPE_CONTROLLER = "controller"
	TIGHTENING_DEVICE_TYPE_TOOL       = "tool"
	MODE_PSET                         = "pset"
	//MODE_JOB                          = "job"
	RESULT_PASS      = "pass"
	RESULT_FAIL      = "fail"
	RESULT_EXCEPTION = "exception"
)

const (
	SocketSelectorTrigger = "SocketSelectorTrigger"
	SocketSelectorClear   = "SocketSelectorClear"
)

type PSetDefine struct {
	//PsetStrategy string  `json:"strategy"`
	Mp float64 `json:"M+"` // 最大扭矩
	Mm float64 `json:"M-"` // 最小扭矩
	Ms float64 `json:"MS"` // 扭矩阈值
	Ma float64 `json:"MA"` // 目标扭矩
	Wp float64 `json:"W+"` // 最大角度
	Wm float64 `json:"W-"` //最小角度
	Wa float64 `json:"WA"` //目标角度
}

type ResultValue struct {
	Mi float64 `json:"measure_torque"` //扭矩
	Wi float64 `json:"measure_angle"`  //角度
	Ti float64 `json:"measure_time"`   //時間
}

type ControllerOutput struct {
	OutputNo int    `json:"no"`
	Status   string `json:"status"`
}

type BaseResult struct {
	// 工具序列号
	ToolSN         string `json:"tool_sn"`
	TighteningUnit string `json:"tightening_unit"`
	// 实际结果
	MeasureResult string `json:"measure_result"`

	// 实际扭矩
	MeasureTorque float64 `json:"measure_torque"`

	// 实际角度
	MeasureAngle float64 `json:"measure_angle"`

	// 实际耗时
	MeasureTime float64 `json:"measure_time"`

	// 批次信息
	Batch string `json:"batch"`

	// 当前点位次序
	Seq int `json:"seq"`

	// 当前点位次序
	GroupSeq int `json:"group_seq"`

	// 当前拧紧次数
	Count int `json:"count"`

	ScannerCode string `json:"scanner_code"`
}

type JobInfo struct {
	// job号
	Job int `json:"job"`
}

type StepData struct {
	PSetDefine
	TighteningResult
}

type TighteningResult struct {
	BaseResult
	JobInfo

	RemoteAddr string `json:"remote_addr"`

	// 控制器序列号
	ControllerSN string `json:"controller_sn"`

	// 控制器名
	ControllerName string `json:"controller_name"`

	// 当前批次
	BatchCount int `json:"batch_count"`

	// 错误代码
	ErrorCode string `json:"error_code"`

	// 工具通道号
	ChannelID int `json:"channel_id"`

	// 收到时间
	UpdateTime time.Time `json:"update_time"`

	// pset号
	PSet int `json:"pset"`

	// 拧紧ID
	TighteningID string `json:"tightening_id"`

	// 拧紧策略
	Strategy string `json:"strategy"`

	// 最大扭矩
	TorqueMax float64 `json:"torque_max"`

	// 最小扭矩
	TorqueMin float64 `json:"torque_min"`

	// 扭矩阈值
	TorqueThreshold float64 `json:"torque_threshold"`

	// 目标扭矩
	TorqueTarget float64 `json:"torque_target"`

	// 最大角度
	AngleMax float64 `json:"angle_max"`

	// 最小角度
	AngleMin float64 `json:"angle_min"`

	// 目标角度
	AngleTarget float64 `json:"angle_target"`

	// 工单ID
	WorkorderID int64 `json:"workorder_id"`

	// 用户ID
	UserID int64 `json:"user_id"`

	// 螺栓编号
	NutNo string `json:"nut_no"`

	// 条码
	Vin string `json:"vin"`

	LastCalibrationDate string `json:"last_calibration_date"`
	TotalToolCounter    int64  `json:"total_tool_counter"`

	// 结果id
	ID int64 `json:"id"`

	// 点位ID
	PointID string `json:"point_id"`

	// 工作模式
	Mode string `json:"mode"`

	// WorkStation
	WorkStation string `json:"work_station"`
	// 分步拧紧的结果 批次号:结果
	StepResults []StepData `json:"step_results"`
}

func (r *TighteningResult) ValidateSet() error {
	if r.Count <= 0 {
		return errors.New("Count Should Be Greater Than 0 ")
	}

	if r.MeasureResult != RESULT_OK && r.MeasureResult != RESULT_NOK {
		return errors.New("MeasureResult Error ")
	}

	if r.MeasureTorque <= 0 {
		return errors.New("MeasureTorque Should Be Greater Than 0 ")
	}

	return nil
}

func (r *TighteningResult) ToDBResult() *storage.Results {
	ll := len(r.StepResults)
	psetDefines := make([]PSetDefine, ll)
	for i := 0; i < ll; i++ {
		pSetInfo := r.StepResults[i]
		psetDefines[i] = pSetInfo.PSetDefine
	}
	strPsetDefine, _ := json.Marshal(psetDefines)
	resultValues := make([]TighteningResult, ll)

	for i := 0; i < ll; i++ {
		pSetInfo := r.StepResults[i]
		resultValues[i] = pSetInfo.TighteningResult

	}
	strResultValue, _ := json.Marshal(resultValues)
	var stepResult string
	if ss, err := json.Marshal(r.StepResults); err == nil {
		stepResult = string(ss)
	}

	return &storage.Results{
		Job:            r.Job,
		ControllerName: r.ControllerName,
		ErrCode:        r.ErrorCode,
		ExInfo:         r.Vin,
		Stage:          r.Strategy,
		ToolCounter:    r.TotalToolCounter,
		HasUpload:      false,
		Seq:            r.Seq,
		WorkorderID:    r.WorkorderID,
		Result:         r.MeasureResult,
		ToolSN:         r.ToolSN,
		TighteningUnit: r.TighteningUnit,
		ControllerSN:   r.ControllerSN,
		TighteningID:   r.TighteningID,
		UpdateTime:     r.UpdateTime,
		PSetDefine:     string(strPsetDefine),
		ResultValue:    string(strResultValue),
		Count:          r.BatchCount,
		PSet:           r.PSet,
		Batch:          r.Batch,
		Spent:          r.MeasureTime,
		UserID:         r.UserID,
		NutNo:          r.NutNo,
		ScannerCode:    r.ScannerCode,
		StepResult:     stepResult,
	}
}

type TighteningCurve struct {
	// 工具序列号
	ToolSN         string `json:"tool_sn"`
	TighteningUnit string `json:"tightening_unit"`
	// 拧紧ID
	TighteningID string `json:"tightening_id"`

	// 收到时间
	UpdateTime time.Time `json:"update_time"`

	TighteningCurveContent
}

func (s *TighteningCurve) GenerateTimeCurveByCoef(coef float32) {
	for k := range s.CUR_M {
		s.CUR_T = append(s.CUR_T, float32(k+1)*coef)
	}
}

func NewTighteningCurve() *TighteningCurve {
	return &TighteningCurve{}
}

func (c *TighteningCurve) ToDBCurve() *storage.Curves {
	curveContent, _ := json.Marshal(c.TighteningCurveContent)

	return &storage.Curves{
		HasUpload:      false,
		UpdateTime:     c.UpdateTime,
		CurveData:      string(curveContent),
		ToolSN:         c.ToolSN,
		TighteningUnit: c.TighteningUnit,
		TighteningID:   c.TighteningID,
		//CurveFile:  fmt.Sprintf("%s_%s.json", c.ToolSN, c.TighteningID),
	}
}

type TighteningCurveContent struct {
	// 实际拧紧结果(ok/nok)
	Result string `json:"result"`

	CUR_M []float32 `json:"cur_m"` // 扭矩
	CUR_W []float32 `json:"cur_w"` // 角度
	CUR_T []float32 `json:"cur_t"` // 时间
	CUR_S []float32 `json:"cur_s"` // 转速
}

type PSetDetail struct {
	PSetID            int     `json:"pset"`
	PSetName          string  `json:"pset_name"`
	RotationDirection string  `json:"rotation_direction"`
	BatchSize         int     `json:"batch_size"`
	TorqueMin         float64 `json:"torque_min"`
	TorqueMax         float64 `json:"torque_max"`
	TorqueTarget      float64 `json:"torque_target"`
	AngleMin          float64 `json:"angle_min"`
	AngleMax          float64 `json:"angle_max"`
	AngleTarget       float64 `json:"angle_target"`
}

type JobDetail struct {
	JobID             int       `json:"job"`
	JobName           string    `json:"job_name"`
	OrderStrategy     string    `json:"order_strategy"`
	CountType         string    `json:"count_type"`
	LockAtJobDone     bool      `json:"lock_at_job_done"`
	UseLineControl    bool      `json:"use_line_control"`
	RepeatJob         bool      `json:"repeat_job"`
	LooseningStrategy string    `json:"loosening_strategy"`
	Steps             []JobStep `json:"steps"`
}

type JobStep struct {
	StepName  string `json:"step_name"`
	ChannelID int    `json:"channel_id"`
	PSetID    int    `json:"pset_id"`
	BatchSize int    `json:"batch_size"`
	Socket    int    `json:"socket"`
}

type ToolMaintenanceInfo struct {
	ToolSN               string `json:"serial_no"`
	ControllerSN         string `json:"controller_sn"`
	TotalTighteningCount int    `json:"times"`
	CountSinLastService  int    `json:"sin_last_service"`
}

type SocketSelectorReq struct {
	PSetSet
	Type string `json:"type"`
}

type Equipment struct {
	EquipmentSN string   `json:"equipment_sn"`
	Location    Location `json:"location"`
	Type        string   `json:"type"`
}

type Location struct {
	EquipmentSN string `json:"equipment_sn"`
	Input       int    `json:"io_input"`
	Output      int    `json:"io_output"`
}

type PSetInfo struct {
	ID   int
	Name string
}
