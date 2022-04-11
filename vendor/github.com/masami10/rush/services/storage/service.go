package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"time"

	"github.com/masami10/rush/typeDef"

	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"xorm.io/xorm"
	xorm_log "xorm.io/xorm/log"
)

type Diagnostic interface {
	Error(msg string, err error)
	Info(msg string)
	Debug(msg string)
	OpenEngineSuccess(info string)
	Close()
	Closed()
}

type Service struct {
	diag        Diagnostic
	configValue atomic.Value
	eng         *xorm.Engine
	validator   *validator.Validate
	connected   *atomic.Bool
	glbIOWriter io.Writer
	glblogLevel string
}

func (s *Service) IsConnected() bool {
	return s.connected.Load()
}

func (s *Service) updateConnectStatus(status bool) {
	s.connected.Store(status)
}

func NewService(c Config, d Diagnostic, out io.Writer, lvl string) *Service {

	s := &Service{
		diag:        d,
		connected:   atomic.NewBool(false),
		glbIOWriter: out,
		glblogLevel: lvl,
	}

	s.configValue.Store(c)
	s.validator = validator.New()

	return s
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) doCreateOrUpdateTable(tb string, schema interface{}) error {
	if err := s.validateStoreService(); err != nil {
		return err
	}
	exist, err := s.eng.IsTableExist(tb)
	if !exist {
		s.diag.Info(fmt.Sprintf("Table: %s Is Not Exist, Will Be Created!", tb))
	}
	if err == nil && !exist {
		// create
		if err = s.eng.Sync2(schema); err != nil {
			return errors.Wrapf(err, "Create Table %s fail", tb)
		}
	}
	if exist && err == nil {
		// update
		//fixme: 现在xorm没有办法改变列的类型，只有手动修改
		if err = s.eng.Sync2(schema); err != nil {
			return errors.Wrapf(err, "Update Table %s fail", tb)
		}
	}
	return err
}

func (s *Service) ShowSql() {
	s.eng.ShowSQL(true)
}

type tb struct {
	name   string
	schema interface{}
}

var RushTables = []tb{
	{"workorders", Workorders{}}, {"time_track", TimeTrack{}},
	{"results", Results{}}, {"curves", Curves{}}, {"controllers", Controllers{}},
	{"tools", Tools{}}, {"routing_operations", RoutingOperations{}}, {"steps", Steps{}},
	{"productivity_loss", ProductivityLoss{}},
}

func (s *Service) manage() {
	init := true
	duration := 1 * time.Millisecond
	for {
		time.Sleep(duration)
		duration = 60 * time.Second
		if s.eng == nil {
			continue
		}
		if err := s.eng.Ping(); err != nil {
			s.updateConnectStatus(false)
			duration = 5 * time.Second
			continue
		}
		s.updateConnectStatus(true)
		if !init {
			continue
		}
		for _, tb := range RushTables {
			if err := s.doCreateOrUpdateTable(tb.name, tb.schema); err != nil {
				s.diag.Error("doCreateOrUpdateTable", err)
			}
		}
		init = false
	}
}

func logLevelFromName(lvl string) xorm_log.LogLevel {
	var level xorm_log.LogLevel = xorm_log.DEFAULT_LOG_LEVEL
	switch lvl {
	case "INFO", "info":
		level = xorm_log.LOG_INFO
	case "ERROR", "error":
		level = xorm_log.LOG_ERR
	case "DEBUG", "debug":
		level = xorm_log.LOG_DEBUG
	}

	return level
}

func (s *Service) Open() error {
	c := s.Config()
	if !c.Enable {
		return nil
	}

	info := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Url,
		c.DBName)
	engine, err := xorm.NewEngine("postgres", info)

	if err != nil {
		return errors.Wrapf(err, "Create postgres engine fail")
	}

	if engine != nil {
		logLevel := logLevelFromName(s.glblogLevel)

		engine.SetLogger(xorm_log.NewSimpleLogger3(s.glbIOWriter, "[RUSH Storage]", log.Ldate|log.Lmicroseconds, logLevel))

		//设置时区
		engine.DatabaseTZ = time.UTC
		engine.TZLocation = time.UTC

		//设置连接池信息
		engine.SetMaxOpenConns(c.MaxConnects) // always success
		engine.SetConnMaxLifetime(10 * time.Minute)
		engine.SetMaxIdleConns(c.MaxConnects / 5)

		// 缓存导致查询不到数据，移除

	}

	s.eng = engine

	s.diag.OpenEngineSuccess(info)

	go s.manage()

	go s.dataRetentionManage() //启动drop数据协程
	return nil
}

func (s *Service) Close() error {
	s.diag.Close()
	if s.eng != nil {
		s.eng.Close()
	}

	s.diag.Closed()

	return nil
}

func (s *Service) validateStoreService() error {

	if !s.Config().Enable || s.eng == nil || !s.IsConnected() {
		return errors.New("Store Engine Is Empty Or Is Not Connected!!!")
	}
	return nil
}

func (s *Service) UpdateRecord(bean interface{}, id int64, data map[string]interface{}) error {
	if err := s.validateStoreService(); err != nil {
		return err
	}

	RowsAffected, err := s.eng.Table(bean).Where("id=?", id).Update(data)
	if err != nil {
		return errors.Wrapf(err, "Update Data Fail")
	}

	s.diag.Debug(fmt.Sprintf("Update Data: %s, Row:%d Success!!!", reflect.TypeOf(bean).Name(), RowsAffected))

	return nil
}

func (s *Service) Store(data interface{}) error {
	if err := s.validateStoreService(); err != nil {
		return err
	}
	RowsAffected, err := s.eng.Insert(data)
	if err != nil || RowsAffected == 0 {
		return errors.Wrapf(err, "Insert New Data Fail")
	}

	s.diag.Debug(fmt.Sprintf("Insert Data: %s, Row:%d Success!!!", reflect.TypeOf(data).Name(), RowsAffected))

	return nil
}

func (s *Service) FindUnuploadResults(result_upload bool, result []string) ([]Results, error) {
	var results []Results

	if err := s.validateStoreService(); err != nil {
		return results, err
	}

	ss := s.eng.Alias("r").Where("r.has_upload = ?", result_upload).And("r.stage = ?", "final").In("r.result", result)

	e := ss.Find(&results)

	if e != nil {
		return results, e
	} else {
		return results, nil
	}
}

func (s *Service) ListUnUploadResults() ([]Results, error) {
	var results []Results

	if err := s.validateStoreService(); err != nil {
		return results, err
	}

	ss := s.eng.Alias("r").Where("r.has_upload = ?", false)

	e := ss.Find(&results)

	if e != nil {
		return results, e
	} else {
		return results, nil
	}
}

func (s *Service) GetTool(serial string) (Tools, error) {
	var tools Tools
	if err := s.validateStoreService(); err != nil {
		return tools, err
	}

	rt, err := s.eng.Alias("g").Where("g.serial = ?", serial).Get(&tools)

	if err != nil {
		return tools, err
	} else {
		if !rt {
			return tools, errors.New("found gun failed")
		} else {
			return tools, nil
		}
	}
}

func (s *Service) UpdateTool(gun *Tools) error {
	if err := s.validateStoreService(); err != nil {
		return err
	}

	g, err := s.GetTool(gun.Serial)
	if err == nil {
		// update
		_, err := s.eng.ID(g.Id).Update(gun)
		if err != nil {
			return err
		}
	} else {
		// insert
		_, err := s.eng.Insert(gun)
		if err != nil {
			return err
		}
	}

	return nil

}

func (s *Service) UpdateCurve(curve *Curves) (*Curves, error) {

	if err := s.validateStoreService(); err != nil {
		return curve, err
	}

	sql := "update `curves` set has_upload = ?, curve_file = ?, curve_data = ? where id = ?"
	_, err := s.eng.Exec(sql,
		curve.HasUpload, curve.CurveFile, curve.CurveData, curve.Id)

	if err != nil {
		return curve, err
	}

	return curve, nil
}

func (s *Service) ListCurvesByResult(result_id int64) ([]Curves, error) {
	var curves []Curves

	if err := s.validateStoreService(); err != nil {
		return curves, err
	}

	e := s.eng.Alias("c").Where("c.result_id = ?", result_id).Find(&curves)
	if e != nil {
		return curves, e
	}
	return curves, nil
}

func (s *Service) ListUnUploadCurves() ([]Curves, error) {
	var curves []Curves

	if err := s.validateStoreService(); err != nil {
		return curves, err
	}

	e := s.eng.Alias("c").Where("c.has_upload = ?", false).And("c.tightening_id != ?", "").Find(&curves)
	if e != nil {
		return curves, e
	}

	return curves, nil
}

func (s *Service) InsertWorkorder(workorder *Workorders, results *[]Results, checkWorkorder bool, checkResult bool, rawid bool) error {
	if err := s.validateStoreService(); err != nil {
		return err
	}

	session := s.eng.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		return err
	}
	// 执行事务

	// 保存新工单
	var raw_workorderid int64
	if workorder != nil {
		if checkWorkorder {
			has, _ := s.WorkorderExists(workorder.WorkorderID)
			if has {
				// 忽略存在的工单
				return nil
			}
		}

		_, err = session.Insert(workorder)
		if err != nil {
			session.Rollback() //nolint errcheck
			return errors.Wrapf(err, "store data fail")
		} else {
			s.diag.Debug(fmt.Sprintf("new workorder:%d", workorder.WorkorderID))
			raw_workorderid = workorder.Id
		}
	}

	// 预保存结果
	if results != nil {
		for _, v := range *results {

			if checkResult {
				has, _ := s.ResultExists(v.ResultId)
				if has {
					continue
				}
			}

			if rawid {
				v.WorkorderID = raw_workorderid
			}

			_, err = session.Insert(v)
			if err != nil {
				session.Rollback() //nolint errcheck
				return errors.Wrapf(err, "store data fail")
			} else {
				s.diag.Debug(fmt.Sprintf("new result:%d", v.ResultId))
			}
		}
	}

	err = session.Commit()
	if err != nil {
		return errors.Wrapf(err, "commit fail")
	}

	return nil
}

func (s *Service) DeleteResultsForJob(key string) error {

	if err := s.validateStoreService(); err != nil {
		return err
	}

	sql := fmt.Sprintf("delete from `results` where exinfo = '%s'", key)
	_, err := s.eng.Exec(sql)

	if err != nil {
		return err
	}
	return nil
}

func (s *Service) WorkorderExists(id int64) (bool, error) {
	if err := s.validateStoreService(); err != nil {
		return false, err
	}

	exist, err := s.eng.Exist(&Workorders{WorkorderID: id})
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (s *Service) ResultExists(id int64) (bool, error) {
	if err := s.validateStoreService(); err != nil {
		return false, err
	}

	exist, err := s.eng.Exist(&Results{ResultId: id})
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (s *Service) GetResultByID(id int64) (*Results, error) {
	if err := s.validateStoreService(); err != nil {
		return nil, err
	}
	result := Results{}
	get, err := s.eng.Alias("r").Where("r.id = ?", id).Limit(1).Get(&result)

	if err != nil {
		return nil, err
	}

	if !get {
		return nil, errors.New("result does not exist")
	} else {
		return &result, nil
	}
}

func (s *Service) GetResultRangeTime(start time.Time, end time.Time) ([]*Results, error) {
	if err := s.validateStoreService(); err != nil {
		return nil, err
	}
	var result []*Results
	err := s.eng.Alias("r").Where("r.update_time between ? and ?", start, end).Find(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) ListWorkorders(hmi_sn string, workcenterCode string, status string) ([]Workorders, error) {
	var workorders []Workorders

	if err := s.validateStoreService(); err != nil {
		return workorders, err
	}

	sql := "select * from workorders where workorders.x_workorder_id > 0 "
	if status != "" {
		sql = sql + fmt.Sprintf(" and workorders.status = '%s'", status)
	}

	sql = sql + " order by update_time, lnr asc"

	err := s.eng.SQL(sql).Find(&workorders)

	return workorders, err
}

func (s *Service) GetWorkOrder(id int64, raw bool) (*Workorders, error) {

	var workOrders Workorders

	if err := s.validateStoreService(); err != nil {
		return nil, err
	}

	key := "x_workorder_id"
	if raw {
		key = "id"
	}

	get, err := s.eng.Alias("w").Where(fmt.Sprintf("w.%s = ?", key), id).Get(&workOrders)

	if err != nil {
		return &workOrders, errors.Wrap(err, "GetWorkOrder")
	}
	if !get {
		return nil, errors.Errorf("GetWorkOrderworkOrders: %d does not exist", id)
	} else {
		return &workOrders, nil
	}
}

func (s *Service) GetStep(id int64) (Steps, error) {

	var step Steps

	if err := s.validateStoreService(); err != nil {
		return step, err
	}

	get, err := s.eng.Alias("w").Where("w.id = ?", id).Get(&step)

	if err != nil {
		return step, err
	}
	if !get {
		return step, errors.New("Step does not exist")
	} else {
		return step, nil
	}
}

func (s *Service) FindWorkOrder(hmi_sn string, workcenter_code string, code string) (*Workorders, error) {

	var workOrders []Workorders

	if err := s.validateStoreService(); err != nil {
		return nil, err
	}

	get, err := s.eng.Alias("w").Where("w.long_pin = ? or w.vin = ? or w.knr = ?", code, code, code).And("w.x_workorder_id > ?", 0).Asc("w.update_time").Asc("w.lnr").Get(&workOrders)

	if err != nil || !get {
		return nil, err
	}
	ll := len(workOrders)
	for i := 0; i < ll; i++ {
		wo := workOrders[i]
		if wo.Status == WORKORDER_STATUS_TODO {
			return &wo, nil
		}
	}
	return &workOrders[0], nil
}

func (s *Service) GetOperationByID(opid int64) (RoutingOperations, error) {

	var routingOp RoutingOperations

	if err := s.validateStoreService(); err != nil {
		return routingOp, err
	}

	get, err := s.eng.Alias("r").Where("r.id = ?", opid).Get(&routingOp)

	if err != nil {
		return routingOp, err
	}
	if !get {
		return routingOp, errors.New("workorder does not exist")
	} else {
		return routingOp, nil
	}
}

func (s *Service) FindNextWorkorder(hmi_sn, workcenter_code string) (Workorders, error) {

	var workorder Workorders

	if err := s.validateStoreService(); err != nil {
		return workorder, err
	}

	get, err := s.eng.Alias("w").Where("w.status = ?", "todo").And("w.x_workorder_id > ?", 0).Asc("w.update_time").Asc("w.lnr").Get(&workorder)

	if err != nil {
		return workorder, err
	}
	if !get {
		return workorder, errors.New("workorder not found")
	} else {
		return workorder, nil
	}
}

//deprecated
func (s *Service) UpdateWorkorderUserID(code string, userID string) error {
	sql := "update `workorders` set user_id = ? where code = ?"
	_, err := s.eng.Exec(sql,
		userID,
		code)

	return err
}

func (s *Service) UpdateWorkorderUserCode(code string, userCode string) error {
	if err := s.validateStoreService(); err != nil {
		return err
	}

	sql := "update `workorders` set user_code = ? where code = ?"
	_, err := s.eng.Exec(sql,
		userCode,
		code)

	return err
}

func (s *Service) UpdateResult(result *Results) (int64, error) {

	if err := s.validateStoreService(); err != nil {
		return 0, err
	}

	sql := "update `results` set controller_sn = ?, result = ?, has_upload = ?, stage = ?, update_time = ?, pset_define = ?, result_value = ?, count = ?, batch = ?, tool_sn = ?, spent = ?, tightening_id = ? where id = ?"
	r, err := s.eng.Exec(sql,
		result.ControllerSN,
		result.Result,
		result.HasUpload,
		result.Stage,
		result.UpdateTime,
		result.PSetDefine,
		result.ResultValue,
		result.Count,
		result.Batch,
		result.ToolSN,
		result.Spent,
		result.TighteningID,
		result.Id)

	if err != nil {

		return 0, errors.Wrapf(err, "Update result fail for id : %d", result.Id)
	} else {
		id, _ := r.RowsAffected()
		return id, nil
	}
}

func (s *Service) UpdateWorkOrderStatus(workorder *Workorders) (*Workorders, error) {
	session := s.eng.NewSession().ForUpdate()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return nil, err
	}

	if workorder.Status == WORKORDER_STATUS_CANCEL {
		if _, err := s.updateWorkOrderCancel(session, workorder.Code); err != nil {
			return workorder, err
		}
	} else {
		if _, err := updateWorkOrderStatusViaID(session, workorder.Id, workorder.Status); err != nil {
			return workorder, err
		}
	}
	err := session.Commit()
	if err != nil {
		return nil, err
	}
	return workorder, nil
}

func (s *Service) UpdateResultByCount(id int64, count int, flag bool) error {

	var err error = nil
	var r sql.Result
	if count > 0 {
		sql := "update `results` set has_upload = ? where id = ? and count = ?"
		r, err = s.eng.Exec(sql, flag, id, count)
	} else {
		sql := "update `results` set has_upload = ? where id = ?"
		r, err = s.eng.Exec(sql, flag, id)
	}

	if err != nil {
		return err
	}

	id, err = r.RowsAffected()

	if id == 0 {
		err = errors.Errorf("Update Result has Upload: %v No Affect!", flag)
		return err
	}

	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteRoutingOperations(rds []RoutingOperationDelete) error {
	var e error = nil
	for _, v := range rds {
		sql := fmt.Sprintf("delete from `routing_operations` where operation_id = %d and product_type = '%s'", v.OperationID, v.ProductType)
		if _, err := s.eng.Exec(sql); err != nil {
			e = errors.Wrap(err, sql)
		}
	}

	return e
}

func (s *Service) DeleteAllRoutingOperations() error {
	_, err := s.eng.Exec("delete from `routing_operations`")

	return err
}

func (s *Service) DeleteInvalidResults(keep time.Time) {
	sql := fmt.Sprintf("delete from `results` where update_time < '%s'", keep.Format("2006-01-02 15:04:05"))
	_, err := s.eng.Exec(sql)

	if err != nil {
		s.diag.Error("DeleteInvalidResults", err)
	}
}

func (s *Service) DeleteInvalidCurves(keep time.Time) {

	sql := fmt.Sprintf("delete from `curves` where update_time < '%s'", keep.Format("2006-01-02 15:04:05"))
	_, err := s.eng.Exec(sql)

	if err != nil {
		s.diag.Error("DeleteInvalidCurves", err)
	}
}

// DeleteInvalidWorkorders delete work orders and steps when status = 'done'
func (s *Service) DeleteInvalidWorkorders(keep time.Time) {
	keepTime := keep.Format("2006-01-02 15:04:05")

	session := s.eng.NewSession()
	if err := session.Begin(); err != nil {
		s.diag.Error("error when delete work orders tx begin.", err)
		return
	}
	exec1 := fmt.Sprintf("DELETE FROM steps WHERE x_workorder_id in (SELECT x_workorder_id FROM workorders where status = 'done' and updated < '%s')", keepTime)
	_, err := s.eng.Exec(exec1)
	if err != nil {
		s.diag.Error("error when delete work orders clear steps", err)
		_ = session.Rollback()
		return
	}
	exec2 := fmt.Sprintf("delete from `workorders` where status = 'done' and updated < '%s'", keepTime)
	_, err = s.eng.Exec(exec2)
	if err != nil {
		s.diag.Error("error when delete work orders", err)
		_ = session.Rollback()
		return
	}

	err = session.Commit()
	if err != nil {
		s.diag.Error("error when delete work orders commit", err)
		return
	}
}

func (s *Service) InitWorkorderForJob(workorder_id int64) error {
	sql := "update `results` set result = ?, stage = ?, count = ? where x_workorder_id = ?"
	_, err := s.eng.Exec(sql, RESULT_NONE, RESULT_STAGE_INIT, 0, workorder_id)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Service) FindTargetResultForJob(workorder_id int64) (Results, error) {
	var results []Results

	ss := s.eng.Alias("r").Where("r.x_workorder_id = ?", workorder_id).And("r.stage = ?", RESULT_STAGE_INIT).OrderBy("r.seq")

	e := ss.Find(&results)

	if e != nil {
		return Results{}, e
	} else {
		if len(results) > 0 {
			return results[0], nil
		} else {
			return Results{}, errors.New("result not found")
		}
	}
}

func (s *Service) FindTargetResultForJobManual(raw_workorder_id int64) (Results, error) {
	var results []Results

	//ss := s.eng.Alias("r").Where("r.x_workorder_id = ?", raw_workorder_id).OrderBy("r.update_time").OrderBy("r.seq").OrderBy("r.count").Desc("r.update_time").Desc("r.seq").Desc("r.count")

	sql := fmt.Sprintf("select * from results where results.x_workorder_id = %d order by results.update_time desc, results.tightening_id desc, results.seq desc, results.count desc", raw_workorder_id)

	e := s.eng.SQL(sql).Find(&results)

	if e != nil {
		return Results{}, e
	} else {
		if len(results) > 0 {
			return results[0], nil
		} else {
			return Results{}, errors.New("result not found")
		}
	}
}

func (s *Service) CreateController(controller_sn, controllerName string) (Controllers, error) {
	var controller Controllers

	exist, err := s.eng.Alias("c").Where("c.controller_sn = ?", controller_sn).Get(&controller)

	if err != nil {
		return controller, err
	}
	if !exist {
		// 创建
		controller.SN = controller_sn
		controller.LastID = "0"
		controller.Name = controllerName
		err = s.Store(controller)
		if err != nil {
			s.diag.Error("CreateController", err)
		}
	} else {
		if controller.Name != controllerName {
			err = s.UpdateRecord(Controllers{}, controller.Id, map[string]interface{}{
				"controller_name": controllerName,
			})
			if err != nil {
				s.diag.Error("UpdateController", err)
			}
		}
	}
	return controller, err
}

func (s *Service) UpdateControllerStatus(controllerSN string, status int) error {
	sql := "update `controllers` set status = ? where controller_sn = ?"
	_, err := s.eng.Exec(sql,
		status,
		controllerSN)

	if err != nil {
		s.diag.Error("UpdateControllerStatus", err)
	}
	return err
}

func (s *Service) UpdateControllerCounter(controllerSN string, counter int64) error {
	sql := "update `controllers` set counter = ? where controller_sn = ?"
	_, err := s.eng.Exec(sql,
		counter,
		controllerSN)

	if err != nil {
		s.diag.Error("UpdateControllerCounter", err)
	}
	return err
}

func (s *Service) UpdateTightening(id int64, last_id string) error {
	sql := "update `controllers` set last_id = ? where id = ?"
	_, err := s.eng.Exec(sql,
		last_id,
		id)

	if err != nil {
		s.diag.Error("UpdateTightening", err)
	}
	return err
}

func (s *Service) ResetTightening(controller_sn string) error {
	sql := "update `controllers` set last_id = ? where controller_sn = ?"
	_, err := s.eng.Exec(sql,
		"0",
		controller_sn)

	if err != nil {
		s.diag.Error("ResetTightening", err)
	}
	return err
}

func (s *Service) UpdateRoutingOperations(ro *RoutingOperations) error {
	affected, err := s.eng.ID(ro.Id).Update(ro)

	if err != nil {
		s.diag.Error("UpdateRoutingOperations", err)
	}
	s.diag.Debug(fmt.Sprintf("Success UpdateRoutingOperations: %d", affected))
	return err
}

func (s *Service) GetRoutingOperations(name, model, step string) (RoutingOperations, error) {

	var ro RoutingOperations

	exist, err := s.eng.Alias("r").Where("r.name = ?", name).And("r.product_type = ?", model).And("r.tightening_step_ref = ?", step).Get(&ro)

	if err != nil {
		s.diag.Error("GetRoutingOperations", err)
		return ro, err
	}
	if !exist {
		return ro, errors.Errorf("Can Not Found RoutingOperations Model: %s, Name: %s", model, name)
	}
	return ro, nil
}

func (s *Service) FindRoutingOperations(workcenter_code, cartype string, job int) (RoutingOperations, error) {

	var ros []RoutingOperations

	ss := s.eng.Alias("r").Where("r.workcenter_code = ?", workcenter_code).And("r.product_type = ? or r.job = ?", cartype, job)

	err := ss.Find(&ros)

	if err != nil {
		s.diag.Error("FindRoutingOperations", err)
		return RoutingOperations{}, err
	}
	if len(ros) > 0 {
		return ros[0], nil
	} else {
		return RoutingOperations{}, errors.Errorf("FindRoutingOperations Model: %s", cartype)
	}
}

func (s *Service) FindLocalStepData(workOrderCode string, limit int, stepTestType string) ([]StepDataWorkorders, error) {
	var results []StepDataWorkorders

	sql := "select * from steps, workorders where steps.x_workorder_id = workorders.id "
	if stepTestType != "" {
		sql = sql + fmt.Sprintf(" and steps.test_type = '%s'", stepTestType)
	}

	if workOrderCode != "" {
		sql = sql + fmt.Sprintf(" and workorders.code = '%s'", workOrderCode)
	}

	sql = sql + " order by steps.updated desc"

	err := s.eng.SQL(sql).Find(&results)

	if limit > 0 && limit < len(results) {
		return results[0:limit], err
	}

	return results, err
}

func (s *Service) FindLocalResults(HMI string, limit int) ([]ResultsWorkorders, error) {
	var results []ResultsWorkorders

	sql := "select * from results, workorders where results.x_workorder_id = workorders.id "
	//if HMI != "" {
	//	sql = sql + fmt.Sprintf(" and workorders.hmi_sn = '%s'", HMI)
	//}

	sql = sql + " order by results.update_time desc"

	err := s.eng.SQL(sql).Find(&results)

	if limit > 0 && limit < len(results) {
		return results[0:limit], err
	}

	return results, err
}

func (s *Service) IsMultiResult(workorderID int64, batch string) bool {
	var results []Results

	ss := s.eng.Alias("r").Where("r.x_workorder_id = ?", workorderID).And("r.batch = ?", batch)

	e := ss.Find(&results)

	if e != nil {
		return false
	} else {
		if len(results) > 1 {
			return true
		} else {
			return false
		}
	}
}

func (s *Service) GetController(sn string) (interface{}, error) {

	var rt_controller Controllers

	rt, err := s.eng.Alias("c").Where("c.controller_sn = ?", sn).Get(&rt_controller)

	if err != nil {
		return nil, err
	} else {
		if !rt {
			return nil, nil
		} else {
			return rt_controller, nil
		}
	}
}

func (s *Service) ResetResult(id int64) error {
	sql := "update `results` set result = ?, has_upload = ?, stage = ?, count = ? where id = ?"

	_, err := s.eng.Exec(sql,
		RESULT_NONE,
		false,
		RESULT_STAGE_INIT,
		1,
		id)

	return err
}

func (s *Service) DeleteCurvesByResult(id int64) error {

	sql := fmt.Sprintf("delete from `curves` where result_id = %d", id)
	_, err := s.eng.Exec(sql)

	if err != nil {
		return err
	}
	return nil
}

func (s *Service) dataRetentionManage() {
	c := s.Config()
	if !c.Enable {
		return
	}
	for {
		start := time.Now()

		keep := time.Now().Add(time.Duration(c.DataKeep) * -1)

		if err := s.validateStoreService(); err == nil {
			// 清理过期结果
			s.DeleteInvalidResults(keep)

			// 清理过期波形
			s.DeleteInvalidCurves(keep)

			// 清理过期工单
			s.DeleteInvalidWorkorders(keep)
		}

		diff := time.Since(start) // 执行的间隔时间

		time.Sleep(time.Duration(c.VacuumPeriod) - diff)
	}
}

func (s *Service) Workorders(par []byte) ([]Workorders, error) {
	orderPar := WorkorderListPar{}
	err := json.Unmarshal(par, &orderPar)
	if err != nil {
		return nil, err
	}

	q := s.eng.Alias("w")

	if len(orderPar.StatusFilter) > 0 {
		q = q.In("status", orderPar.StatusFilter)
	} else if orderPar.Status != "" {
		q = q.Where("w.status = ?", orderPar.Status)
	} else {
		q = q.Where("w.status != ?", orderPar.Status)
	}

	if orderPar.Time_from != "" {
		q = q.And("w.date_planned_start >= ?", orderPar.Time_from)
	}
	if orderPar.Time_to != "" {
		q = q.And("w.date_planned_complete <= ?", orderPar.Time_to)
	}
	orderBy := orderPar.OrderBy
	if orderBy == "" {
		orderBy = "id"
	}
	if !orderPar.IsAsc {
		q.Desc(orderBy)
	} else {
		q.Asc(orderBy)
	}
	if orderPar.Page_size == 0 {
		orderPar.Page_size = 20
	}
	q.Limit(orderPar.Page_size, orderPar.Page_no*orderPar.Page_size)

	var rt []Workorders

	err = q.Find(&rt)

	for i := 0; i < len(rt); i++ {
		rt[i].ProductTypeImage, err = s.findOrderPicture(rt[i].ProductCode)
	}

	if err != nil {
		return nil, err
	} else {
		return rt, nil
	}
}

func (s *Service) Steps(workorderID int64) ([]Steps, error) {
	var rt []Steps

	q := s.eng.Alias("s").Where("s.workorder_id = ?", workorderID).Asc("s.id")

	err := q.Find(&rt)

	if err != nil {
		return nil, err
	} else {
		return rt, nil
	}
}

func (s *Service) WorkorderStep(workorderID int64) (*WorkorderStep, error) {
	workorder, err := s.GetWorkOrder(workorderID, true)
	if err != nil {
		return nil, err
	}

	var p interface{}
	err = json.Unmarshal([]byte(workorder.Payload), &p)
	if err == nil {
		workorder.MarshalPayload = p
	}

	steps, err := s.Steps(workorderID)
	if err != nil {
		return nil, err
	}

	for k, v := range steps {
		err := json.Unmarshal([]byte(v.Step), &steps[k].MarshalPayload)
		if err != nil {
			s.diag.Error(err.Error(), err)
		}
	}

	return &WorkorderStep{
		Workorders: *workorder,
		Steps:      steps,
	}, nil
}

func updateWorkOrderStatusViaID(ss *xorm.Session, id int64, status string) (bool, error) {
	sql := fmt.Sprintf(`update %s set status='%s' where id=%d`, "workorders", status, id)
	if _, err := ss.Exec(sql); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func updateWorkOrderStatus(ss *xorm.Session, code, status string) (bool, error) {
	sql := fmt.Sprintf("update %s set status='%s' where code='%s'", "workorders", status, code)
	if _, err := ss.Exec(sql); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (s *Service) updateWorkOrderCancel(ss *xorm.Session, code string) (bool, error) {
	var workorder Workorders

	exist, e := ss.Alias("r").Where("r.code = ?", code).Get(&workorder)
	if e != nil {
		return true, e
	}
	if !exist {
		return true, nil
	}

	if _, err := updateWorkOrderStatus(ss, workorder.Code, WORKORDER_STATUS_CANCEL); err != nil {
		s.diag.Error("Update Work Order Status Cancel Error", err)
		return false, err
	}

	if _, err := ss.Table(new(Steps)).Where("x_workorder_id = ?", workorder.Id).Update(map[string]interface{}{"status": STEP_STATUS_CANCEL}); err != nil {
		s.diag.Error("Update Steps Status Cancel Error", err)
		return false, err
	}

	return true, nil
}

func (s *Service) DeleteWorkAndStep(ss *xorm.Session, code string, uniqueNum int64) (bool, error) {
	var workorder Workorders

	exist, e := ss.Alias("r").Where("r.code = ?", code).Get(&workorder)
	if e != nil {
		return true, e
	}
	if !exist {
		return true, nil
	}
	if workorder.UniqueNum > uniqueNum {
		return false, nil
	}
	sql1 := "delete from `workorders` where id = ?"
	sql2 := "delete from `steps` where x_workorder_id = ?"

	if _, err := ss.Exec(sql1, workorder.Id); err != nil {
		s.diag.Error(fmt.Sprintf("Delete Work Order: %s Error", workorder.Code), err)
		return false, err
	}

	if _, err := ss.Exec(sql2, workorder.Id); err != nil {
		s.diag.Error(fmt.Sprintf("Delete Work Order: %s Related Work Step Error", workorder.Code), err)
		return false, err
	}
	return true, nil
}

func (s *Service) UpdateStepStatus(step *Steps) (*Steps, error) {

	sql := "update `steps` set status = ? where id = ?"
	_, err := s.eng.Exec(sql,
		step.Status,
		step.Id)

	if err != nil {
		return step, err
	} else {
		return step, nil
	}
}

func (s *Service) UpdateStepData(step *Steps) (*Steps, error) {

	sql := "update `steps` set data = ? where id = ?"
	_, err := s.eng.Exec(sql,
		step.Data,
		step.Id)

	if err != nil {
		return step, err
	} else {
		return step, nil
	}
}

func (s *Service) UpdateOrderData(order *Workorders) (*Workorders, error) {

	sql := "update `workorders` set data = ? where id = ?"
	_, err := s.eng.Exec(sql,
		order.Data,
		order.Id)

	if err != nil {
		return order, err
	} else {
		return order, nil
	}
}

func (s *Service) GetLastIncompleteCurve(toolSN string) (*Curves, error) {
	curve := Curves{}
	e := s.eng.Alias("c").Where("c.tightening_id = ?", "").And("c.tool_sn = ?", toolSN).Desc("c.update_time").Find(&curve)

	if e != nil {
		return &curve, e
	} else {
		return &curve, nil
	}
}

func (s *Service) getLasIncompleteCurveForTool(session *xorm.Session, TighteningUnit string) (*Curves, error) {

	curve := Curves{}
	exist, err := session.Alias("c").Where("c.tightening_id = ?", "").And("c.tightening_unit = ?", TighteningUnit).Desc("c.update_time").Get(&curve)

	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("Curve Not Found")
	}

	return &curve, nil
}

func (s *Service) updateIncompleteCurve(session *xorm.Session, curve *Curves) error {

	sql := "update `curves` set tightening_id = ?, curve_file = ?,tool_sn = ? where id = ?"
	_, err := session.Exec(sql,
		curve.TighteningID,
		curve.CurveFile,
		curve.ToolSN,
		curve.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) StorageInsertResult(result *Results) error {

	if !s.Config().Enable {
		return nil
	}

	session := s.eng.NewSession()
	defer session.Close()

	// 执行事务
	err := session.Begin()
	if err != nil {
		return err
	}

	// 保存结果
	_, err = session.Insert(result)
	if err != nil {
		return err
	}

	err = session.Commit()
	if err != nil {
		return errors.Wrapf(err, "Commit Failed")
	}

	return nil
}

func (s *Service) UpdateIncompleteCurveAndSaveResult(result *Results) error {

	if !s.Config().Enable {
		return nil
	}

	if result.Result == RESULT_LSN {
		return nil
	}

	session := s.eng.NewSession()
	defer session.Close()

	// 执行事务
	err := session.Begin()
	if err != nil {
		return err
	}

	// 获取最近不完整的曲线
	curve, err := s.getLasIncompleteCurveForTool(session, result.TighteningUnit)

	if err == nil {
		// 更新曲线
		curve.TighteningID = result.TighteningID
		curve.CurveFile = fmt.Sprintf("%s_%s.json", result.ToolSN, result.TighteningID)
		curve.ToolSN = result.ToolSN
		err = s.updateIncompleteCurve(session, curve)
		if err != nil {
			s.diag.Error("UpdateIncompleteCurveAndSaveResult Update Curve Error", err)
		} else {
			result.CurveFile = curve.CurveFile
		}
	}

	// 保存结果
	_, err = session.Insert(result)
	if err != nil {
		return err
	}

	err = session.Commit()
	if err != nil {
		return errors.Wrapf(err, "Commit Failed")
	}

	return nil
}

func (s *Service) getLasIncompleteResultForTool(session *xorm.Session, TighteningUnit string) (*Results, error) {
	result := Results{}
	if !s.Config().Enable {
		return nil, errors.New("Not Enabled")
	}

	exist, err := session.Alias("r").Where("r.curve_file = ?", "").And("r.tightening_unit = ?", TighteningUnit).
		And("r.result != ?", RESULT_LSN).And("r.result != ?", RESULT_AK2).Desc("r.update_time").Get(&result)

	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("Result Not Found")
	}

	return &result, nil
}

func (s *Service) updateIncompleteResult(session *xorm.Session, result *Results) error {

	sql := "update `results` set curve_file = ? where id = ?"
	_, err := session.Exec(sql,
		result.CurveFile,
		result.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateIncompleteResultAndSaveCurve(curve *Curves, result *Results) error {
	if !s.Config().Enable {
		return nil
	}

	session := s.eng.NewSession()
	defer session.Close()

	// 执行事务
	err := session.Begin()
	if err != nil {
		return err
	}

	result.CurveFile = fmt.Sprintf("%s_%s.json", result.ToolSN, result.TighteningID)
	err = s.updateIncompleteResult(session, result)
	if err != nil {
		s.diag.Error("UpdateIncompleteResultAndSaveCurve updateIncompleteResult Error", err)
	} else {
		curve.CurveFile = result.CurveFile
		curve.TighteningID = result.TighteningID
		curve.ToolSN = result.ToolSN
	}

	// 保存曲线
	_, err = session.Insert(curve)
	if err != nil {
		return err
	}

	err = session.Commit()
	if err != nil {
		return errors.Wrapf(err, "UpdateIncompleteResultAndSaveCurve Commit Fail")
	}

	return nil
}

func (s *Service) getIncompleteResults(session *xorm.Session, toolSN string) ([]Results, error) {
	results := []Results{}

	err := session.Alias("r").Where("r.curve_file = ?", "").And("r.tool_sn = ?", toolSN).Desc("r.update_time").Find(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *Service) getIncompleteCurves(session *xorm.Session, toolSN string) ([]Curves, error) {
	curves := []Curves{}

	err := session.Alias("c").Where("c.tightening_id = ?", "").And("c.tool_sn = ?", toolSN).Desc("c.update_time").Find(&curves)
	if err != nil {
		return nil, err
	}

	return curves, nil
}

func (s *Service) ClearToolResultAndCurve(toolSN string) error {
	if !s.Config().Enable {
		return nil
	}

	session := s.eng.NewSession()
	defer session.Close()

	// 执行事务
	err := session.Begin()
	if err != nil {
		return err
	}

	// 获取所有不完整的结果
	results, err := s.getIncompleteResults(session, toolSN)
	if err != nil {
		return err
	}

	// 处理不完整的结果
	for i, r := range results {
		r.CurveFile = fmt.Sprintf("%s_%s.json", r.ToolSN, r.TighteningID)
		err = s.updateIncompleteResult(session, &results[i])
		if err != nil {
			return err
		}
	}

	// 获取所有不完整的曲线
	curves, err := s.getIncompleteCurves(session, toolSN)
	if err != nil {
		return err
	}

	// 处理不完整的曲线
	for i, c := range curves {
		c.TighteningID = utils.GenerateID()
		c.CurveFile = fmt.Sprintf("invalid_%s_%s.json", c.ToolSN, c.TighteningID)
		err = s.updateIncompleteCurve(session, &curves[i])
		if err != nil {
			return err
		}
	}

	err = session.Commit()
	if err != nil {
		return errors.Wrapf(err, "Commit Fail")
	}

	return nil
}

func (s *Service) GetWorkorderByCode(code string) (*Workorders, error) {
	order := Workorders{}
	if !s.Config().Enable {
		return &order, nil
	}

	rt, err := s.eng.Alias("w").Where("w.code = ?", code).Limit(1).Get(&order)

	if err != nil {
		return nil, err
	} else {
		if !rt {
			return nil, errors.New("Workorder Does Not Exist")
		} else {
			return &order, nil
		}
	}
}

func (s *Service) GetWorkorderByTraceCode(code string) (*Workorders, error) {
	order := Workorders{}
	if !s.Config().Enable {
		return &order, nil
	}

	rt, err := s.eng.Alias("w").Where("w.track_code = ?", code).Limit(1).Get(&order)

	if err != nil {
		return nil, err
	} else {
		if !rt {
			return nil, errors.New("Workorder Does Not Exist")
		} else {
			return &order, nil
		}
	}
}

func (s *Service) GetStepsByWorkorderID(orderID int64) ([]Steps, error) {
	var steps []Steps
	if err := s.eng.Alias("s").Where("s.x_workorder_id = ?", orderID).Find(&steps); err != nil {
		return nil, err
	}

	return steps, nil
}

func (s *Service) GetResultsByStepID(stepID int64) ([]Results, error) {
	var results []Results
	if err := s.eng.Alias("r").Where("r.step_id = ?", stepID).Find(&results); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *Service) GetStepByCodeAndWorkorderID(code string, workorderID int64) (*Steps, error) {
	step := Steps{}
	var err error
	rt := false
	if !s.Config().Enable {
		return &step, nil
	}
	if workorderID != 0 {
		rt, err = s.eng.Alias("s").Where("s.code = ?", code).And("s.x_workorder_id = ?", workorderID).Limit(1).Get(&step)
	} else {
		rt, err = s.eng.Alias("s").Where("s.code = ?", code).Limit(1).Get(&step)

	}

	if err != nil {
		return nil, err
	}
	if !rt {
		return nil, errors.New("Step Does Not Exist")
	}
	return &step, nil
}

// seq, count
func (s *Service) CalBatch(workorderID int64) (int, int) {
	result, err := s.FindTargetResultForJobManual(workorderID)
	if err != nil {
		return 1, 1
	}

	if result.Result == RESULT_OK {
		return result.GroupSeq + 1, 1
	} else {
		return result.GroupSeq, result.Count + 1
	}
}

func (s *Service) PatchResultFromDB(result *Results, mode string) error {
	dbTool, err := s.GetTool(result.TighteningUnit)
	if err != nil {
		return err
	}

	result.ScannerCode = dbTool.ScannerCode
	result.PointID = dbTool.PointID
	if dbTool.CurrentWorkorderID != 0 {

		if mode == typeDef.MODE_JOB {
			result.Seq, result.Count = s.CalBatch(dbTool.CurrentWorkorderID)
		} else {
			result.Seq = dbTool.Seq
			result.Count = dbTool.Count
		}

		result.WorkorderID = dbTool.CurrentWorkorderID
		result.UserID = dbTool.UserID
		result.Batch = fmt.Sprintf("%d/%d", result.Seq, dbTool.Total)
		result.StepID = dbTool.StepID

		dbStep, err := s.GetStep(dbTool.StepID)
		if err != nil {
			s.diag.Error("Get Step Failed", err)
			return err
		}

		consume, err := s.GetConsumeBySeqInStep(&dbStep, result.Seq)
		if err != nil {
			s.diag.Error("Get Consume Failed", err)
			return err
		}

		result.NutNo = consume.NutNo
		result.TighteningPointName = consume.TighteningPointName
	}

	return nil
}

func (s *Service) GetConsumeBySeqInStep(step *Steps, seq int) (*StepComsume, error) {
	if step == nil {
		return nil, errors.New("Step Is Nil")
	}

	ts := TighteningStep{}
	if err := json.Unmarshal([]byte(step.Step), &ts); err != nil {
		return nil, err
	}

	if len(ts.TighteningPoints) == 0 {
		return nil, errors.New("Consumes Is Empty")
	}

	for k, v := range ts.TighteningPoints {
		if v.Seq == seq {
			return &ts.TighteningPoints[k], nil
		}
	}

	return nil, errors.New("Consume Not Found")
}

func (s *Service) GetResultByTighteningID(toolSN string, tid string) (*Results, error) {
	var result Results
	rt, err := s.eng.Where("tightening_id = ?", tid).And("tool_sn = ?", toolSN).Desc("update_time").Limit(1).Get(&result)

	if err != nil {
		return nil, err
	}
	if !rt {
		return nil, errors.New("Result Not Found")
	}
	return &result, nil
}

func (s *Service) GetLastResultByStepIDAndSeq(stepID int64, seq int) (*Results, error) {
	result := Results{}
	rt, err := s.eng.Alias("r").Where("r.step_id = ?", stepID).And("r.seq = ?", seq).Desc("r.id").Limit(1).Get(&result)

	if err != nil {
		return nil, err
	}
	if !rt {
		return nil, errors.New("Result Not Found")
	}
	return &result, nil
}

func (s *Service) UpdateToolLocation(toolSN string, location string) error {
	tool, err := s.GetTool(toolSN)
	if err != nil {
		return err
	}

	tool.Location = location
	_, err = s.eng.ID(tool.Id).Update(tool)
	return err
}

func (s *Service) GetToolLocation(toolSN string) (string, error) {
	tool, err := s.GetTool(toolSN)
	if err != nil {
		return "", err
	}

	return tool.Location, nil
}

func (s *Service) GetCurveByfile(file string) (*Curves, error) {
	curve := Curves{}
	rt, err := s.eng.Alias("c").Where("c.curve_file = ?", file).Limit(1).Get(&curve)

	if err != nil {
		return nil, err
	}
	if !rt {
		return nil, errors.New("Curve Not Found")
	}
	return &curve, nil
}

func (s *Service) GetResultByTighteningIDAndToolSN(tid string, toolSN string) (*Results, error) {
	return s.GetResultByTighteningID(toolSN, tid)
}

func (s *Service) CreateRecord(dbRecord interface{}) error {
	_, err := s.eng.Insert(dbRecord)
	return err
}

func (s *Service) CreateRecords(dbRecords ...interface{}) error {
	_, err := s.eng.Insert(dbRecords)
	return err
}

func (s *Service) ExistedScannerCodeFromStep(code string) bool {
	r, err := s.eng.Exec("SELECT * FROM steps WHERE test_type='register_byproducts' AND position(? in data) > 0", code)

	if err != nil {
		return false
	}
	if rows, err := r.RowsAffected(); err == nil {
		if rows > 0 {
			return true
		}
	}

	return false
}
