package io

const (
	OutputStatusOff   = 0
	OutputStatusOn    = 1
	OutputStatusFlash = 2
)

type IoSet struct {
	SN     string `json:"sn"`
	Index  uint16 `json:"index"`
	Status uint16 `json:"status"`
}

type IoData struct {
	Inputs  string `json:"inputs"`
	Outputs string `json:"outputs"`
}

type IoConfig struct {
	InputNum  uint16 `json:"input_num"`
	OutputNum uint16 `json:"output_num"`
}

type IoContact struct {
	Src     string `json:"src"`
	SN      string `json:"sn"`
	Inputs  string `json:"inputs"`
	Outputs string `json:"outputs"`
}
