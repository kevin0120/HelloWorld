package wsnotify

//deprecated
//type WSRegist struct {
//	HMISn string `json:"hmi_sn"`
//}

type WSMsg struct {
	SeqNumber uint64      `json:"sn"`
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
}

type WSReply struct {
	SeqNumber uint64 `json:"sn"`
	Result    int    `json:"result"`
	Msg       string `json:"msg"`
}
