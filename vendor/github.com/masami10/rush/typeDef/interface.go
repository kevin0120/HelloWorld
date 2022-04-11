package typeDef

import (
	"strconv"
)

//为了创建工单所使用的Step模型
type StepTighteningPayload struct {
	TighteningTotal int                     `json:"tightening_total"`
	TighteningPoint []RoutingOperationPoint `json:"tightening_points"`
}

type Pset struct {
	Value int
}

func (pset *Pset) MarshalJSON() (data []byte, err error) {
	dd := strconv.Itoa(pset.Value)
	return []byte(dd), nil
}

func (pset *Pset) UnmarshalJSON(data []byte) (err error) {
	d := string(data)
	dd, err := strconv.Atoi(d)
	pset.Value = dd
	return err
}

type JobPoint struct {
	Seq                int     `json:"sequence"`
	PSet               Pset    `json:"pset"`
	X                  float64 `json:"x"`
	Y                  float64 `json:"y"`
	MaxRedoTimes       int     `json:"max_redo_times"`
	GroupSeq           int     `json:"group_sequence"`
	ConsuProductID     int64   `json:"consu_product_id"`
	ToleranceMin       float64 `json:"tolerance_min"`
	ToleranceMax       float64 `json:"tolerance_max"`
	ToleranceMinDegree float64 `json:"tolerance_min_degree"`
	ToleranceMaxDegree float64 `json:"tolerance_max_degree"`
}

// RoutingOperationPoint 是标准拧紧工步的点位信息
type RoutingOperationPoint struct {
	JobPoint
	BoltNumber          string   `json:"nut_no"`
	IsKey               bool     `json:"is_key"`
	KeyCount            int      `json:"key_num"`
	TighteningTool      []string `json:"tightening_unit"`
	TighteningPointName string   `json:"tightening_point_name"`
	// Deprecated
	ControllerSN []string `json:"controller_sn"`
}

type RoutingOperation struct {
	OperationID    int64                   `json:"id"`
	Job            int                     `json:"job"`
	MaxOpTime      int                     `json:"max_op_time"`
	Name           string                  `json:"name"`
	Img            string                  `json:"img"`
	ProductId      int64                   `json:"product_id"`
	ProductType    string                  `json:"product_type"`
	WorkCenterCode string                  `json:"workcenter_code"`
	VehicleTypeImg string                  `json:"vehicleTypeImg"`
	Points         []RoutingOperationPoint `json:"points"`
}
