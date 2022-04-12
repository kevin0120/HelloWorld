package tightening_device

import (
	"encoding/json"
	"fmt"
	"github.com/masami10/rush/utils"
	"time"

	"github.com/kataras/iris/v12/websocket"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/wsnotify"
)

var getPSetInfoWithDetail = utils.GetEnvBool("ENV_PSET_INFO_DETAIL", false)
var useExternalControl = utils.GetEnvBool("ENV_EXTERNAL_CONTROL", false)

func (s *Service) OnWS_TOOL_MODE_SELECT(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	req := ToolModeSelect{}
	_ = json.Unmarshal(byteData, &req)
	err := s.ToolModeSelect(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, 0, ""), s.diag)
}

func (s *Service) OnWS_TOOL_ENABLE(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	req := ToolControl{}
	_ = json.Unmarshal(byteData, &req)

	if useExternalControl {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, 0, ""), s.diag)
		return
	}

	err := s.ToolControl(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, 0, ""), s.diag)

	//FIXME: 此功能未来使用nodered控制 使能/切断使能后，控制工具对应LED灯。
	s.doDispatch(dispatcherbus.DispatcherToolEnable, req)

	if !req.Enable && s.config().SocketSelector.Enable {
		s.doSelectorPatch()
	}
}

func (s *Service) OnWS_TOOL_JOB(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req JobSet
	_ = json.Unmarshal(byteData, &req)
	err := s.ToolJobSet(&req, true)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, 0, ""), s.diag)

}

func (s *Service) OnWS_TOOL_PSET(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req PSetSet
	if err := json.Unmarshal(byteData, &req); err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
	}

	if useExternalControl {
		s.doTraceFromPSetReq(&req)
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, 0, ""), s.diag)
		return
	}

	if s.config().SocketSelector.Enable {
		//如果使能套筒选择器，不设定程序号，而是控制套筒选择器
		s.doTraceFromPSetReq(&req)
		s.doSelectorPatch()

		_msg := fmt.Sprintf("工步: %s, 拧紧点: %d, 请选择套筒选择器地址: %d", req.StepCode, req.Sequence, req.PSet)
		wsnotify.WSNotifyInfo(c, _msg, s.diag)
	} else {
		err := s.ToolPSetSet(&req)
		if err != nil {
			_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
			return
		}
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, 0, ""), s.diag)

}

func (s *Service) OnWS_TOOL_PSET_BATCH(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req PSetBatchSet
	_ = json.Unmarshal(byteData, &req)

	err := s.ToolPSetBatchSet(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, 0, ""), s.diag)
}

func (s *Service) OnWS_TOOL_PSET_LIST(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)
	var pSetList interface{}
	var req ToolInfo
	_ = json.Unmarshal(byteData, &req)

	pSetNums, err := s.GetToolPSetList(&req)

	if getPSetInfoWithDetail {
		var _pSetList []*PSetDetail
		for i := 0; i < len(pSetNums); i++ {
			pSetInfo := pSetNums[i]
			pSetDetail, err1 := s.GetToolPSetDetail(&ToolPSet{
				req,
				pSetInfo.ID,
			})
			if err1 != nil {
				s.diag.Error("error when GetToolPSetDetail", err1)
				continue
			}
			_pSetList = append(_pSetList, pSetDetail)
		}
		pSetList = _pSetList
	} else {
		var _pSetList []*PSetDetail
		for i := 0; i < len(pSetNums); i++ {
			_pSetList = append(_pSetList, &PSetDetail{
				PSetID:   pSetNums[i].ID,
				PSetName: pSetNums[i].Name,
			})
		}
		pSetList = _pSetList
	}

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SeqNumber, msg.Type, pSetList), s.diag)
}

func (s *Service) OnWS_TOOL_PSET_DETAIL(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req ToolPSet
	_ = json.Unmarshal(byteData, &req)

	psetDetail, err := s.GetToolPSetDetail(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SeqNumber, msg.Type, psetDetail), s.diag)
}

func (s *Service) OnWS_TOOL_JOB_LIST(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req ToolInfo
	_ = json.Unmarshal(byteData, &req)

	jobList, err := s.GetToolJobList(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SeqNumber, msg.Type, jobList), s.diag)
}

func (s *Service) OnWS_TOOL_JOB_DETAIL(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req ToolJob
	_ = json.Unmarshal(byteData, &req)

	jobDetail, err := s.GetToolJobDetail(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SeqNumber, msg.Type, jobDetail), s.diag)
}

//手动回补填写结果
func (s *Service) OnWS_TOOL_RESULT_MANUAL_SET(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var result TighteningResult
	_ = json.Unmarshal(byteData, &result)

	if err := result.ValidateSet(); err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	dbTool, err := s.storageService.GetTool(result.TighteningUnit)
	if err == nil {
		dbTool.Count = result.Count
		//fixme: xorm 无法设置0的情况
		if result.Count <= 0 {
			dbTool.Count = 1
		}
		if err1 := s.storageService.UpdateTool(&dbTool); err1 != nil {
			_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -3, err1.Error()), s.diag)
			return
		}
	}

	// 未设置userid会导致数据插入失败
	if result.UserID == 0 {
		result.UserID = 1
	}
	result.TighteningID = utils.GenerateID()
	//手动输入时间以rush收到为准
	result.UpdateTime = time.Now()
	// 处理数据
	result.StepResults = []StepData{
		StepData{
			PSetDefine{}, result,
		},
	}

	dbResult := result.ToDBResult()
	err = s.storageService.StorageInsertResult(dbResult)
	if err != nil {
		s.diag.Error("Handle Result With Curve Failed", err)
	}
	// 分发结果
	result.WorkorderID = dbResult.WorkorderID
	result.UserID = dbResult.UserID
	result.Batch = dbResult.Batch
	result.ID = dbResult.Id
	result.Count = dbResult.Count
	result.Seq = dbResult.Seq
	result.GroupSeq = dbResult.GroupSeq

	result.ScannerCode = dbResult.ScannerCode
	result.PointID = dbResult.PointID

	s.doDispatch(dispatcherbus.DispatcherResult, &result)
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, 0, ""), s.diag)
}
