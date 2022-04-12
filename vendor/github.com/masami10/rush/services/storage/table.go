package storage

import (
	"time"
)

// Deprecated: please use ProductivityLoss instead
type TimeTrack struct {
	Id            int64     `xorm:"pk autoincr notnull 'id'" json:"id"`
	WorkOrderCode string    `xorm:"varchar(128) index notnull 'workorder_code'" json:"workorder_code"`
	StartTime     time.Time `xorm:"'start_time'" json:"start_time"`
	EndTime       time.Time `xorm:"'end_time'" json:"end_time"`
	Duration      float64   `xorm:"'duration'" json:"duration"`
}

type ProductivityLoss struct {
	Id            int64     `xorm:"pk autoincr notnull 'id'" json:"id"`
	WorkOrderCode string    `xorm:"varchar(128) index notnull 'workorder_code'" json:"workorder_code"`
	StartTime     time.Time `xorm:"'start_time'" json:"start_time"`
	EndTime       time.Time `xorm:"'end_time'" json:"end_time"`
	Duration      float64   `xorm:"'duration'" json:"duration"`
	LossType      string    `xorm:"varchar(64) notnull 'loss_type'" json:"loss_type"`
	LossCode      string    `xorm:"varchar(128) 'loss_code'" json:"loss_code"`
	LossDesc      string    `xorm:"varchar(128) 'desc'" json:"desc"`
}

type Workorders struct {
	Id             int64  `xorm:"pk autoincr notnull 'id'" json:"id"`
	Origin         string `xorm:"varchar(256) 'origin'" json:"origin,omitempty"` // 当没有origin为空的时候，发送给hmi接口中得抑制这个字段
	WorkorderID    int64  `xorm:"bigint 'x_workorder_id'" json:"-"`
	HMISN          string `xorm:"varchar(64) 'hmi_sn'" json:"-"`
	WorkcenterCode string `xorm:"varchar(64) 'workcenter_code'" json:"workcenter"`
	Vin            string `xorm:"varchar(64) 'vin'" json:"vin"`
	Knr            string `xorm:"varchar(64) 'knr'" json:"-"`
	LongPin        string `xorm:"varchar(64) 'long_pin'" json:"-"`

	MaxOpTime    int   `xorm:"int 'max_op_time'" json:"-"`
	MaxSeq       int   `xorm:"int 'max_seq'" json:"-"`
	LastResultID int64 `xorm:"bigint 'last_result_id'" json:"-"`
	//WorkSheet      string    `xorm:"text 'work_sheet'"`
	ImageOPID      int64     `xorm:"bigint 'img_op_id'" json:"-"`
	VehicleTypeImg string    `xorm:"text 'vehicle_type_img'" json:"-"`
	UpdateTime     time.Time `xorm:"datetime 'update_time'" json:"-"`
	ProductID      int64     `xorm:"bigint 'product_id'" json:"-"`
	WorkcenterID   int64     `xorm:"bigint 'workcenter_id'" json:"-"`
	//deprecated
	UserID int64 `xorm:"bigint 'user_id'" json:"-"`

	UserCode string `xorm:"varchar(64) 'user_code'" json:"user_code"`
	JobID    int    `xorm:"bigint 'job_id'" json:"-"`
	Mode     string `xorm:"varchar(64) 'mode'" json:"-"`

	Consumes string `xorm:"text 'consumes'" json:"-"`

	// mo相关信息
	MO_EquipemntName  string `xorm:"varchar(64) 'equipment_name'" json:"-"` // 设备名
	MO_FactoryName    string `xorm:"varchar(64) 'factory_name'" json:"-"`   // 工厂代码
	MO_Year           int64  `xorm:"bigint 'year'" json:"-"`
	MO_Pin            int64  `xorm:"bigint 'pin'" json:"-"`
	MO_Pin_check_code int64  `xorm:"bigint 'pin_check_code'" json:"-"`
	MO_AssemblyLine   string `xorm:"varchar(64) 'assembly_line'" json:"-"`
	MO_Lnr            string `xorm:"varchar(64) 'lnr'" json:"lnr"` //车序
	MO_Model          string `xorm:"varchar(64) 'model'" json:"-"`

	Code                string    `xorm:"unique varchar(128) 'code'"  json:"code"`
	TrackCode           string    `xorm:"varchar(128) 'track_code'" json:"track_code"` //产成品追踪码
	ProductCode         string    `xorm:"varchar(128) 'product_code'" json:"product_code"  validate:"required"`
	FinishedProduct     string    `xorm:"varchar(128) 'finished_product'" json:"finished_product"  validate:"required"`
	DatePlannedStart    time.Time `xorm:"datetime 'date_planned_start'" json:"date_planned_start"`
	DatePlannedComplete time.Time `xorm:"datetime 'date_planned_complete'" json:"date_planned_complete"`
	UniqueNum           int64     `xorm:"bigint 'unique_num'" json:"unique_num"`
	Data                string    `xorm:"text 'data'" json:"data"`

	Name             string      `xorm:"varchar(64) 'name'" json:"-"`
	Desc             string      `xorm:"varchar(64) 'desc'" json:"-"`
	Payload          string      `xorm:"text" json:"-"`
	MarshalPayload   interface{} `xorm:"-" json:"-"`
	Workorder        string      `xorm:"text 'workorder'" json:"-"`
	Status           string      `xorm:"varchar(32) default 'todo' 'status' " json:"status"`
	ProductTypeImage string      `json:"product_type_image"`
	Created          time.Time   `xorm:"created" json:"-"`
	Updated          time.Time   `xorm:"updated" json:"-"`

	TimesTracking []TimeTrack `xorm:"-" json:"times_tracking"`

	Steps []Steps `json:"-" validate:"required"`
}

type StepScannerInfo struct {
	Result string `json:"result"`
}

type Steps struct {
	Id             int64  `xorm:"pk autoincr notnull 'id'" json:"id"`
	Title          string `xorm:"-" json:"title"`
	Code           string `xorm:"varchar(128) 'code'" json:"code"`
	Sequence       int64  `xorm:"bigint 'sequence'" json:"sequence"`
	TestType       string `xorm:"varchar(128) 'test_type'" json:"test_type"`
	FailureMessage string `xorm:"varchar(128) 'failure_msg'" json:"failure_msg"`
	Desc           string `xorm:"text 'desc'" json:"desc"`
	Data           string `xorm:"text 'data'" json:"data"` // hmi更新的数据段
	// deprecated
	ImageRef       string  `xorm:"varchar(128) 'tightening_image_by_step_code'" json:"tightening_image_by_step_code" validate:"required"`
	Skippable      bool    `xorm:"varchar(64) 'skippable'" json:"skippable"`
	Undoable       bool    `xorm:"varchar(64) 'undoable'" json:"undoable"`
	Text           string  `xorm:"text 'text'" json:"text"`
	ToleranceMin   float64 `xorm:"Double 'tolerance_min'" json:"tolerance_min"`
	ToleranceMax   float64 `xorm:"Double 'tolerance_max'" json:"tolerance_max"`
	Target         float64 `xorm:"Double 'target'" json:"target"`
	ConsumeProduct string  `xorm:"varchar(128) 'consume_product'" json:"consume_product"`
	//deprecated
	Payload        string      `xorm:"text" json:"payload"`
	MarshalPayload interface{} `xorm:"-" json:"-"`

	WorkorderID int64  `xorm:"bigint 'x_workorder_id'" json:"-"`
	Step        string `xorm:"text 'step'" json:"-"`

	Image  string `xorm:"text" json:"img"`
	Status string `xorm:"varchar(32) default 'ready' 'status'" json:"status"`

	Created time.Time `xorm:"created" json:"-"`
	Updated time.Time `xorm:"updated" json:"-"`

	// 工步时间限制
	TimeLimit int64 `xorm:"bigint 'time_limit'" json:"time_limit"`
}

type Results struct {
	Id                 int64     `xorm:"pk autoincr notnull 'id'" gorm:"primaryKey"`
	HasUpload          bool      `xorm:"bool 'has_upload'"`
	Seq                int       `xorm:"int 'seq'"`
	GroupSeq           int       `xorm:"int 'group_sequence'"`
	ResultId           int64     `xorm:"bigint 'x_result_id'"`
	WorkorderID        int64     `xorm:"bigint 'x_workorder_id'"`
	StepID             int64     `xorm:"bigint 'step_id'"`
	UserID             int64     `xorm:"bigint 'user_id'"`
	ControllerSN       string    `xorm:"varchar(64) 'controller_sn'"`
	ControllerName     string    `xorm:"varchar(64) 'controller_name'"`
	ErrCode            string    `xorm:"varchar(64) 'err_code'"`
	ToolSN             string    `xorm:"varchar(64) 'tool_sn' unique(unique_result_everytime_tool_pk)"`
	TighteningUnit     string    `xorm:"varchar(64) 'tightening_unit'"`
	Result             string    `xorm:"varchar(32) 'result'"`
	Stage              string    `xorm:"varchar(32) 'stage'"`
	UpdateTime         time.Time `xorm:"datetime 'update_time' unique(unique_result_everytime_tool_pk)"`
	PSetDefine         string    `xorm:"text 'pset_define'"`
	ResultValue        string    `xorm:"text 'result_value'"`
	Count              int       `xorm:"int 'count'"`
	PSet               int       `xorm:"int 'pset'"`
	Job                int       `xorm:"int 'job'"`
	NutNo              string    `xorm:"varchar(64) 'nut_no'"`
	ConsuProductID     int64     `xorm:"bigint 'consu_product_id'"`
	ToleranceMinDegree float64   `xorm:"Double 'tolerance_min_degree'"`
	ToleranceMaxDegree float64   `xorm:"Double 'tolerance_max_degree'"`
	ToleranceMax       float64   `xorm:"Double 'tolerance_max'"`
	ToleranceMin       float64   `xorm:"Double 'tolerance_min'"`
	OffsetX            float64   `xorm:"Double 'offset_x'"`
	OffsetY            float64   `xorm:"Double 'offset_y'"`
	MaxRedoTimes       int       `xorm:"int 'max_redo_times'"`
	Batch              string    `xorm:"varchar(32) 'batch'"`
	ExInfo             string    `xorm:"text 'exinfo'"` // OP协议错误代码保存的地方
	Spent              float64   `xorm:"Double 'spent'"`
	TighteningID       string    `xorm:"varchar(128) 'tightening_id' unique(unique_result_everytime_tool_pk)"`
	CurveFile          string    `xorm:"varchar(256) 'curve_file'"`
	ToolCounter        int64     `xorm:"bigint 'tool_counter'"`
	// payload字段
	ScannerCode string `xorm:"varchar(256) 'scanner_code'"`
	PointID     string `xorm:"varchar(128) 'point_id'"`

	StepResult          string `xorm:"text 'stepResult'"` //分段拧紧结果
	TighteningPointName string `xorm:"varchar(128)" json:"tightening_point_name"`
}

type ResultsWorkorders struct {
	Results    `xorm:"extends"`
	Workorders `xorm:"extends"`
}

type StepDataWorkorders struct {
	Steps      `xorm:"extends"`
	Workorders `xorm:"extends"`
}

type Curves struct {
	Id             int64     `xorm:"pk autoincr notnull 'id'"`
	ResultID       int64     `xorm:"bigint 'result_id'"`
	Count          int       `xorm:"int 'count'"`
	CurveFile      string    `xorm:"varchar(128) 'curve_file'"`
	CurveData      string    `xorm:"text 'curve_data'"`
	HasUpload      bool      `xorm:"bool 'has_upload'"`
	UpdateTime     time.Time `xorm:"datetime 'update_time'"`
	ToolSN         string    `xorm:"varchar(128) 'tool_sn'"`
	TighteningUnit string    `xorm:"varchar(128) 'tightening_unit'"`
	TighteningID   string    `xorm:"varchar(128) 'tightening_id'"`
}

type Controllers struct {
	Id           int64     `xorm:"pk autoincr notnull 'id'"`
	SN           string    `xorm:"varchar(128) 'controller_sn'"`
	Name         string    `xorm:"varchar(128) 'controller_name'"`
	LastID       string    `xorm:"varchar(128) 'last_id'"`
	TriggerStart time.Time `xorm:"datetime 'trigger_start'"`
	TriggerStop  time.Time `xorm:"datetime 'trigger_stop'"`
	Status       int       `xorm:"int 'status'"`
	Counter      int64     `xorm:"bigint 'counter'"`
}

type Tools struct {
	Id                 int64  `xorm:"pk autoincr notnull 'id'"`
	GunID              int64  `xorm:"bigint 'gun_id'"`
	Serial             string `xorm:"varchar(128) 'serial'"`
	CurrentWorkorderID int64  `xorm:"bigint 'workorder_id'"` //当前正在进行的工单
	Seq                int    `xorm:"bigint 'sequence'"`
	Count              int    `xorm:"int 'count'"` //FIXME: xorm 无法设置为0
	Total              int    `xorm:"int 'total'"`
	Mode               string `xorm:"varchar(128) 'mode'"`
	UserID             int64  `xorm:"bigint 'user_id'"`
	StepID             int64  `xorm:"bigint 'step_id'"`
	PointID            string `xorm:"varchar(128) 'point_id'"`
	ScannerCode        string `xorm:"varchar(128) 'scanner_code'"`
	Location           string `xorm:"text 'location'"`
	Batch              int    `xorm:"int 'batch'"`
}

type RoutingOperations struct {
	Id int64 `xorm:"pk autoincr notnull 'id'"`
	// Deprecated
	OperationID       int64  `xorm:"bigint 'operation_id'"`
	Job               int    `xorm:"bigint 'job'"`
	MaxOpTime         int    `xorm:"int 'max_op_time'"`
	Name              string `xorm:"varchar(256) 'name'"`
	Img               string `xorm:"text 'img'"`
	TighteningStepRef string `xorm:"varchar(256) 'tightening_step_ref'"`

	ProductId    int64 `xorm:"bigint 'product_id'"`
	WorkcenterID int64 `xorm:"bigint 'workcenter_id'"`

	ProductType      string `xorm:"varchar(256) 'product_type'"`
	ProductTypeImage string `xorm:"text 'product_type_image'"`

	WorkcenterCode string `xorm:"varchar(256) 'workcenter_code'"`
	VehicleTypeImg string `xorm:"text 'vehicle_type_img'"`
	Points         string `xorm:"text 'points'"`
	Steps          string `xorm:"text 'steps'"`
}
