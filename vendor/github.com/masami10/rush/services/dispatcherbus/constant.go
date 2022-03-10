package dispatcherbus

const (
	// ********************************Device***********************************
	// 设备状态(包括IO，条码枪，读卡器，拧紧控制器，拧紧工具等)发生变化时，会通过此分发器进行状态分发
	DispatcherDeviceStatus = "DISPATCHER_DEVICE_STATUS"

	// ********************************Scanner***********************************
	// 当收到条码数据时(来自条码枪，拧紧控制器条码等)，会通过此分发器进行条码分发
	DispatcherScannerData = "DISPATCHER_SCANNER_DATA"

	// ********************************Scanner***********************************
	// 当收到通用外设数据时，会通过此分发器进行条码分发
	DispatcherAnyDeviceInputData = "DISPATCHER_ANY_DEVICE_INPUT"

	// ********************************io***********************************
	// 当收到IO输入输出状态变化时(IO模块或拧紧控制器IO等)，会通过此分发器进行IO状态分发
	DispatcherIO = "DISPATCHER_IO"

	// ********************************Reader***********************************
	// 当收到读卡器数据时，会通过此分发器进行读卡器数据分发
	DispatcherReaderData = "DISPATCHER_READER_DATA"

	// ********************************Reader***********************************
	// 当收到使能请求或关闭是否分发
	DispatcherToolEnable = "DISPATCHER_TOOL_ENABLE"

	// ********************************Tightening***********************************
	// 当收到拧紧结果时，会通过此分发器进行拧紧结果分发
	DispatcherResult = "DISPATCHER_RESULT"

	// 当收到拧紧曲线时，会通过此分发器进行拧紧曲线分发
	DispatcherCurve = "DISPATCHER_CURVE"

	// 当收到控制器推送的JOB时，会通过此分发器进行分发
	DispatcherJob = "DISPATCHER_JOB"

	// 当收到控制器推送的维护通知时，会通过此分发器进行分发
	DispatcherToolMaintenance = "DISPATCHER_TOOL_MAINTENANCE"

	// 当收到控制器推送的工具错误信息时，会通过此分发器进行分发
	DispatcherToolError = "DISPATCHER_TOOL_ERROR"

	// ********************************Service***********************************
	// 当检测到服务状态变化时(aiis, odoo, 外部系统等)，会通过此分发器进行状态分发
	DispatcherServiceStatus = "DISPATCHER_SERVICE_STATUS"

	// ********************************Transport***********************************
	// 当Transport服务状态变化时， 会通过此分发器进行状态分发
	DispatcherTransportStatus = "DISPATCHER_TRANSPORT_STATUS"

	// ********************************WEBSOCKET***********************************
	// 当收到WebSocket请求时， 会通过此分发器进行请求分发
	DispatcherWsNotify = "DISPATCHER_WS_NOTIFY"

	// ********************************HMI***********************************
	//// 当收到HMI发来的开工请求时，会向此分发器发送数据。可以根据具体需求订阅并处理
	DispatcherOrderStart = "DISPATCHER_ORDER_START"
	//
	//// 当收到HMI发来的完工请求时，会向此分发器发送数据。可以根据具体需求订阅并处理
	DispatcherOrderFinish = "DISPATCHER_ORDER_FINISH"

	// 当收到HMI发来的完工阻塞
	DispatcherOrderPending = "DISPATCHER_ORDER_PENDING"

	// 当收到HMI发来的工单继续
	DispatcherOrderResume = "DISPATCHER_ORDER_RESUME"

	// 当收到HMI或者服务器发来的工单取消
	DispatcherOrderCancel = "DISPATCHER_ORDER_CANCEL"

	// 当收到下发的新工单时，会将新工单数据发到此分发器
	DispatcherOrderNew = "DISPATCHER_ORDER_NEW"

	// 当收到ODOO推送的工具保养通知时，会将保养信息发送到此分发器
	DispatcherMaintenanceInfo = "DISPATCHER_MAINTENANCE_INFO"

	// 套筒控制
	DispatcherSocketSelector = "DispatcherSocketSelector"

	//放行
	DispatcherWorkSiteLeaving = "DispatcherWorkSiteLeaving"
)
