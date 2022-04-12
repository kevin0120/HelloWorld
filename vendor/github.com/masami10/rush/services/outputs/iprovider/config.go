package iprovider

import (
	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
	"os"
)

var factory_code string

func init() {
	factory_code = os.Getenv("factory_code")
}

type Plugin struct {
	Type     string            `yaml:"type"`
	Url      string            `yaml:"url"`
	Path     string            `yaml:"path"`
	Format   string            `yaml:"format"`
	User     string            `yaml:"user"`
	PassWord string            `yaml:"password"`
	Method   string            `yaml:"method"`
	Headers  map[string]string `yaml:"headers"`
}

func (p *Plugin) Hash() string {
	return fmt.Sprintf("%s:%s:%s", p.Type, p.Url, p.Path)
}

/*
新payload
*/

type PublishPackage struct {
	ReplaceMicroSeconds string         `json:"replace_microseconds"`
	Conf                PackageContent `json:"conf"`
}

type PackageContent struct {
	EntityID    string `json:"entity_id"`
	FactoryCode string `json:"factory_code"`
	// 工艺类型
	CraftType int `json:"craft_type"`

	Result PackageResult `json:"result"`
	Curve  PackageCurve  `json:"curve"`
	Task   Task          `json:"task"`
	Param  Param         `json:"curve_param"`
}

type Task struct {
	DagID      string `json:"dag_id"`
	TaskID     string `json:"task_id"`
	RealTaskID string `json:"real_task_id"`
	ExeDate    string `json:"exec_date"`
}

type Param struct {
	// 扭矩上限区间最小值
	TorqueUpLimitMin float32 `json:"torque_up_limit_min"`

	// 扭矩上限区间最大值
	TorqueUpLimitMax float32 `json:"torque_up_limit_max"`

	// 曲线比对阈值
	CompareThreshold float32 `json:"compare_threshold"`

	// 曲线斜率比对阈值
	SlopeThreshold float32 `json:"slope_threshold"`

	// 扭矩差值阈值
	TorqueDiffThreshold float32 `json:"torque_diff_threshold"`

	// 训练窗口
	Windows [][]float32 `json:"windows"`
}

type PackageResult struct {
	tightening_device.TighteningResult
}

type PackageCurve struct {
	CurM []float32 `json:"cur_m"`
	CurW []float32 `json:"cur_w"`
	CurT []float32 `json:"cur_t"`
	CurS []float32 `json:"cur_s"`
}

func (s *PublishPackage) FromDBResultAndCurve(dbResult *storage.Results, dbCurve *storage.Curves) {
	if dbResult == nil || dbCurve == nil {
		return
	}

	s.Conf.Result.NutNo = fmt.Sprintf("%s-%d-%d", dbResult.ControllerName, dbResult.Job, dbResult.Count)
	s.Conf.Result.TighteningID = dbResult.TighteningID
	s.Conf.Result.MeasureResult = dbResult.Result
	// 把多段数据保存在pset_info的list中
	var resultValues []tightening_device.TighteningResult
	_ = json.Unmarshal([]byte(dbResult.ResultValue), &resultValues)
	var psetDefines []tightening_device.PSetDefine
	_ = json.Unmarshal([]byte(dbResult.PSetDefine), &psetDefines)

	resultValueNum := len(resultValues)
	for idx := 0; idx < resultValueNum; idx++ {
		stepData := tightening_device.StepData{PSetDefine: psetDefines[idx], TighteningResult: resultValues[idx]}
		s.Conf.Result.StepResults = append(s.Conf.Result.StepResults, stepData)
	}
	// 原来的数据现在变成最后一段的数据
	resultValue := resultValues[len(resultValues)-1]
	s.Conf.Result.MeasureTorque = resultValue.MeasureTorque
	s.Conf.Result.MeasureAngle = resultValue.MeasureAngle
	psetDefine := psetDefines[len(resultValues)-1]
	s.Conf.Result.Strategy = resultValue.Strategy
	s.Conf.Result.TorqueMax = psetDefine.Mp
	s.Conf.Result.TorqueMin = psetDefine.Mm
	s.Conf.Result.AngleMax = psetDefine.Wp
	s.Conf.Result.AngleMin = psetDefine.Wm
	s.Conf.Result.AngleTarget = psetDefine.Wa
	s.Conf.Result.TorqueTarget = psetDefine.Ma

	s.Conf.Result.ToolSN = dbResult.ToolSN
	s.Conf.Result.Job = dbResult.Job
	s.Conf.Result.PSet = dbResult.PSet
	s.Conf.Result.UpdateTime = dbResult.UpdateTime
	s.Conf.Result.ControllerName = dbResult.ControllerName
	s.Conf.Result.ErrorCode = dbResult.ErrCode
	s.Conf.Result.BatchCount = dbResult.Count
	s.Conf.Result.Vin = dbResult.ExInfo

	s.Conf.Result.TotalToolCounter = dbResult.ToolCounter

	curveData := tightening_device.TighteningCurveContent{}
	_ = json.Unmarshal([]byte(dbCurve.CurveData), &curveData)
	s.Conf.Curve.CurM = curveData.CUR_M
	s.Conf.Curve.CurW = curveData.CUR_W
	s.Conf.Curve.CurT = curveData.CUR_T
	s.Conf.Curve.CurS = curveData.CUR_S

	pkgName := fmt.Sprintf("%s/%s/%d", s.Conf.Result.ToolSN, s.Conf.Result.TighteningID, s.Conf.Result.UpdateTime.Unix())
	s.Conf.EntityID = pkgName

	s.Conf.FactoryCode = factory_code
	s.Conf.Result.Batch = dbResult.Batch
}

type FileItem struct {
	Filename string
	Data     []byte
}
