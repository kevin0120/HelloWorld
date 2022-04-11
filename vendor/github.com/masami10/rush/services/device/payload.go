package device

type Status struct {
	SN             string      `json:"sn"`
	TighteningUnit string      `json:"tightening_unit"` // use by tool control
	Type           string      `json:"type"`
	Status         string      `json:"status"`
	Children       interface{} `json:"children"`
	Config         interface{} `json:"config"`
	Data           interface{} `json:"data"`
}

type AnyDeviceData struct {
	SN     string      `json:"sn"`
	Source string      `json:"src"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
