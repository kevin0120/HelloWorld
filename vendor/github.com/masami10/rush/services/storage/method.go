package storage

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"
	"xorm.io/xorm"
)

func (s *Service) getLastWorkOrderTimeTrack(orderCode string, canInsert bool) (*TimeTrack, error) {
	var tt TimeTrack
	session := s.eng.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return nil, errors.Errorf("Session Begin Error")
	}
	ss := session.Where("workorder_code = ?", orderCode).Limit(1).Desc("id")
	_, e := ss.Get(&tt)
	if e == nil {
		return &tt, nil
	}
	//不存在 创建
	if !canInsert {
		return nil, errors.Errorf("Can Not Found The Time Tracking Info for Order: %s", orderCode)
	}
	if id, err := session.Insert(tt); err != nil {
		return nil, errors.Wrap(err, "Insert TimeTrack Error")
	} else {
		tt.Id = id
	}
	return &tt, nil
}

func (s *Service) AppendTimeTrackEndTime(orderCode string, endTime time.Time) error {
	timeTrackingInfo, err := s.getLastWorkOrderTimeTrack(orderCode, false)
	if err != nil {
		return err
	}
	timeTrackingInfo.EndTime = endTime
	timeTrackingInfo.Duration = endTime.Sub(timeTrackingInfo.StartTime).Seconds() // 计算duration
	_, err = s.eng.ID(timeTrackingInfo.Id).Update(timeTrackingInfo)

	return err
}

func (s *Service) AppendTimeTrackStartTime(orderCode string, startTime time.Time) error {
	timeTrackingInfo, err := s.getLastWorkOrderTimeTrack(orderCode, true)
	if err != nil {
		return err
	}
	timeTrackingInfo.StartTime = startTime
	_, err = s.eng.ID(timeTrackingInfo.Id).Update(timeTrackingInfo)

	return err
}

func (s *Service) AppendProductivityLoss(code string, end time.Time, lossType, desc, lossCode string) error {
	session := s.eng.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return errors.Errorf("Session Begin Error")
	}
	productLoss := &ProductivityLoss{
		WorkOrderCode: code,
		EndTime:       end,
		LossCode:      lossCode,
		LossType:      lossType,
		LossDesc:      desc,
	}

	if _, err = session.Insert(productLoss); err != nil {
		return errors.Wrap(err, "Insert ProductivityLoss Error")
	}
	err = session.Commit()
	if err != nil {
		return errors.Wrap(err, "Insert ProductivityLoss Commit Error")
	}
	return nil
}

//參考文檔
//https://github.com/go-xorm/xorm/blob/master/README_CN.md
//http://gobook.io/read/github.com/go-xorm/manual-zh-CN/chapter-01/

func (s *Service) WorkorderSync(work *Workorders) (string, error) {

	err := s.validator.Struct(work)
	if err != nil {
		return "", errors.Wrapf(err, "loss workorder-steps information")
	}

	session := s.eng.NewSession().ForUpdate()
	if err = session.Begin(); err != nil {
		return "", errors.Wrapf(err, "Create Sesssion Error")
	}
	defer session.Close()
	if work.Status == WORKORDER_STATUS_CANCEL {
		//更新订单状态到取消
		if success, err := s.updateWorkOrderCancel(session, work.Code); !success {
			session.Rollback() //nolint: errcheck
			return "", errors.Wrapf(err, "updateWorkOrderCancel")
		}
		return work.Code, nil
	}
	success, err := s.DeleteWorkAndStep(session, work.Code, work.UniqueNum)
	if !success {
		session.Rollback() //nolint: errcheck
		return "", errors.Wrapf(err, "本地已有更新版本号对应的工单")
	}

	_, err = session.Insert(work)
	if err != nil {
		session.Rollback() //nolint: errcheck
		return "", errors.Wrapf(err, "store data fail")
	}

	for i := 0; i < len(work.Steps); i++ {
		work.Steps[i].WorkorderID = work.Id
		if _, err := session.Insert(work.Steps[i]); err != nil {
			session.Rollback() //nolint: errcheck
			return "", errors.Wrapf(err, "store data fail")
		}

	}
	// update WorkorderID
	if _, err = session.Where("id = ?", work.Id).Update(&Workorders{WorkorderID: work.Id}); err != nil {
		session.Rollback() //nolint: errcheck
		return "", errors.Wrapf(err, "store Update Work Order ID fail")
	}

	err = session.Commit()
	if err != nil {
		return "", errors.Wrapf(err, "commit fail")
	}

	return work.Code, nil
}

func (s *Service) WorkorderIn(in []byte) (string, error) {

	var work Workorders
	var workPayload WorkorderPayload
	err := json.Unmarshal(in, &work)
	if err != nil && !strings.Contains(err.Error(), "workcenter") {
		return "", err
	}
	err = json.Unmarshal(in, &workPayload)
	if err != nil {
		return "", err
	}

	workcenterBody, _ := json.Marshal(workPayload.Workcenter)
	var wc WorkCenterPayload
	if err = json.Unmarshal(workcenterBody, &wc); err != nil {
		return "", err
	}

	wp, err := json.Marshal(workPayload)
	if err != nil {
		return "", err
	}

	workorder1 := Workorders{
		//大众相关字段
		MO_Lnr:            work.MO_Lnr,
		Origin:            work.Origin,
		LongPin:           work.LongPin,
		Knr:               work.Knr,
		MO_Pin_check_code: work.MO_Pin_check_code,
		Vin:               work.Vin,
		//
		Code:                work.Code,
		TrackCode:           work.TrackCode,
		ProductCode:         work.ProductCode,
		DatePlannedStart:    work.DatePlannedStart,
		DatePlannedComplete: work.DatePlannedComplete,
		UniqueNum:           work.UniqueNum, //业务场景同一张工单发两次
		Workorder:           string(wp),
		WorkcenterCode:      wc.Code,
		FinishedProduct:     work.FinishedProduct,
		UpdateTime:          work.UpdateTime,
		Status:              WORKORDER_STATUS_TODO,
	}

	if work.Status == WORKORDER_STATUS_CANCEL {
		workorder1.Status = work.Status
	}

	var workorderMap map[string]interface{}
	var step []map[string]interface{}

	err = json.Unmarshal(in, &workorderMap)

	if err != nil {
		return "", errors.Wrap(err, "JSON Unmarshal WorkorderMap Error")
	}

	if _, exist := workorderMap["steps"]; exist {
		steps, _ := json.Marshal(workorderMap["steps"])
		if err = json.Unmarshal(steps, &step); err != nil {
			s.diag.Error("JSON Unmarshal WorkorderMap Steps Error", err)
		}
	}

	for i := 0; i < len(step); i++ {
		stepString, _ := json.Marshal(step[i])
		var msg Steps
		//var stepText StepTextPayload
		var stepTightening StepTighteningPayload
		var sp []byte
		if err = json.Unmarshal(stepString, &msg); err != nil {
			s.diag.Error("JSON Unmarshal stepString Error", err)
			continue
		}

		if msg.TestType == TestTypeTightening || msg.TestType == TestTypePromiscuousTightening {
			if err = json.Unmarshal(stepString, &stepTightening); err != nil {
				s.diag.Error("JSON Unmarshal stepString To Tightening Step Error", err)
			} else {
				sp, _ = json.Marshal(stepTightening)
			}
		} else {
			sp = stepString
		}

		payload, _ := json.Marshal(map[string]string{
			"title": msg.Title,
		})

		step := Steps{
			ImageRef:       msg.ImageRef,
			TestType:       msg.TestType,
			Code:           msg.Code,
			Sequence:       msg.Sequence,
			FailureMessage: msg.FailureMessage,
			Desc:           msg.Desc,
			Skippable:      msg.Skippable,
			Undoable:       msg.Undoable,
			Data:           msg.Data,

			Text:           msg.Text,
			ToleranceMin:   msg.ToleranceMin,
			ToleranceMax:   msg.ToleranceMax,
			Target:         msg.Target,
			ConsumeProduct: msg.ConsumeProduct,
			TimeLimit:      msg.TimeLimit,
			WorkorderID:    workorder1.Id,
			Step:           string(sp),

			Status:  "ready",
			Payload: string(payload),
		}

		workorder1.Steps = append(workorder1.Steps, step)
	}
	return s.WorkorderSync(&workorder1)
}

func (s *Service) WorkorderOut(orderCode string, workorderID int64) (interface{}, error) {

	var workorder Workorders
	var ss *xorm.Session
	if orderCode == "" {
		ss = s.eng.Alias("r").Where("r.id = ?", workorderID)
	} else {
		ss = s.eng.Alias("r").Where("r.code = ?", orderCode)
	}

	ok, err := ss.Get(&workorder)
	if err != nil {
		return nil, err
	}
	// 如果code没有查到订单就用说明发的是vin码，为了向下兼容
	if !ok && orderCode != "" {
		ss = s.eng.Alias("r").Where("r.vin = ?", orderCode)
		ok, e := ss.Get(&workorder)
		if e != nil || !ok {
			return nil, e
		}
	}

	var step []Steps
	err = s.eng.SQL("select * from steps where x_workorder_id = ?", workorder.Id).Find(&step)

	if err != nil {
		return nil, err
	}

	var steps []map[string]interface{}
	for i := 0; i < len(step); i++ {
		stepMap := stringToMap(step[i].Step)
		stepOut := strucToMap(step[i])
		stepOut["payload"] = stepMap

		if step[i].TestType != TestTypeTightening && step[i].TestType != TestTypePromiscuousTightening {
			delete(stepOut, "tightening_image_by_step_code")
			steps = append(steps, stepOut)
			continue
		}
		stepOut["image"] = step[i].Image
		delete(stepOut, "tightening_image_by_step_code")
		steps = append(steps, stepOut)
	}

	map2 := stringToMap(workorder.Workorder)
	workOrderOut := strucToMap(workorder)
	workOrderOut["steps"] = steps
	workOrderOut["payload"] = map2

	image2, _ := s.findOrderPicture(workorder.ProductCode)
	workOrderOut["product_type_image"] = image2
	//delete(workOrderOut,"product_code")
	//rr, _ := json.Marshal(workOrderOut)

	return workOrderOut, nil
}

func (s *Service) WorkorderStart(order string, workorderID int64) (StartRequest, error) {
	var startpayload StartRequest
	var workPayload WorkorderPayload
	var workorder Workorders
	var ss *xorm.Session
	if order == "" {
		ss = s.eng.Alias("r").Where("r.id = ?", workorderID)
	} else {
		ss = s.eng.Alias("r").Where("r.code = ?", order)
	}

	exit, e := ss.Get(&workorder)
	if e != nil || !exit {
		return startpayload, e
	}

	err := json.Unmarshal([]byte(workorder.Workorder), &workPayload)
	if err != nil {
		return startpayload, err
	}

	startpayload = StartRequest{
		WIPORDERNO:    workorder.Code,
		WIPORDERTYPE:  workPayload.WIPORDERTYPE,
		OPRSEQUENCENO: workPayload.OPRSEQUENCENO,
		UPDATEON:      time.Now(),
		UPDATEBY:      workPayload.STARTEMPLOYEE,
		LOCATION:      workPayload.LOCATION,
		RESOURCEGROUP: workPayload.RESOURCEGROUP,
		SKILL:         workPayload.SKILL,
		RESOURCENAME:  workPayload.RESOURCENAMES,
	}

	return startpayload, nil
}

func (s *Service) WorkorderFinished(order string, workorderID int64) (FinishedRequest, error) {
	var finishedpayload FinishedRequest
	var workPayload WorkorderPayload
	var workorder Workorders
	var ss *xorm.Session
	if order == "" {
		ss = s.eng.Alias("r").Where("r.id = ?", workorderID)
	} else {
		ss = s.eng.Alias("r").Where("r.code = ?", order)
	}

	exist, e := ss.Get(&workorder)
	if e != nil || !exist {
		return finishedpayload, e
	}

	err := json.Unmarshal([]byte(workorder.Workorder), &workPayload)
	if err != nil {
		return finishedpayload, err
	}

	finishedpayload = FinishedRequest{
		WIPORDERNO:    workorder.Code,
		WIPORDERTYPE:  workPayload.WIPORDERTYPE,
		OPRSEQUENCENO: workPayload.OPRSEQUENCENO,
		UPDATEON:      time.Now(),
		UPDATEBY:      workPayload.STARTEMPLOYEE,
		LOCATION:      workPayload.LOCATION,
		RESOURCEGROUP: workPayload.RESOURCEGROUP,
		SKILL:         workPayload.SKILL,
		RESOURCENAME:  workPayload.RESOURCENAMES,
	}

	return finishedpayload, nil
}

func strucToMap(in interface{}) (m map[string]interface{}) {
	j, _ := json.Marshal(in)
	err := json.Unmarshal(j, &m)
	if err != nil {
		return nil
	}
	return
}

func stringToMap(in string) (m map[string]interface{}) {
	err := json.Unmarshal([]byte(in), &m)
	if err != nil {
		return nil
	}
	return
}

func (s *Service) findStepPicture(ref string) (string, error) {

	var ro RoutingOperations
	ss := s.eng.Alias("r").Where("r.tightening_step_ref = ?", ref).Limit(1)
	_, e := ss.Get(&ro)
	if e != nil {
		return "", e
	}
	return ro.Img, nil
}

func (s *Service) findOrderPicture(ref string) (string, error) {
	ro, err := s.GetRoutingOperationViaProductTypeCode(ref)
	if err != nil {
		err := errors.Errorf("GetRoutingOperationViaProductTypeCode: Product Type: %s Fail", ref)
		s.diag.Error("findOrderPicture", err)
		return "", err
	}
	return ro.ProductTypeImage, nil
}

func (s *Service) GetRoutingOperationViaProductTypeCode(ProductType string) (*RoutingOperations, error) {
	var ro RoutingOperations
	ss := s.eng.Alias("r").Where("r.product_type = ?", ProductType).Limit(1)
	b, e := ss.Get(&ro)
	if e != nil || !b {
		return nil, errors.Errorf("Operation For Product Type Code: %s Is Not Existed", ProductType)
	}
	return &ro, nil
}

func (s *Service) GetRoutingOperationsViaProductTypeCode(ProductType string) ([]RoutingOperations, error) {
	var ro []RoutingOperations
	ss := s.eng.Alias("r").Where("r.product_type = ?", ProductType).OrderBy("id")
	e := ss.Find(&ro)
	if e != nil {
		return nil, errors.Errorf("Operation For Product Type Code: %s Is Not Existed", ProductType)
	}
	return ro, nil
}

func (s *Service) GetAllRoutingOperations() ([]RoutingOperations, error) {
	var ro []RoutingOperations
	ss := s.eng.Alias("r").OrderBy("id")
	e := ss.Find(&ro)
	if e != nil {
		return nil, errors.Wrap(e, "GetAllRoutingOperations Error")
	}
	return ro, nil
}

func (s *Service) GetResultByWorkOrderID(orderID int64) ([]Results, error) {

	var results []Results
	var ss *xorm.Session
	if orderID == 0 {
		return nil, errors.New("GetResultByWorkOrderID Order ID Is Required")
	} else {
		ss = s.eng.Alias("r").Where("r.x_workorder_id = ?", orderID)
	}

	e := ss.Find(&results)
	if e != nil {
		return nil, e
	} else {
		return results, nil
	}
}
