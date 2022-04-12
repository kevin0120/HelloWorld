package io

const (
	ModelMoxaE1212 = "MOXA_E1212" // 8进8出
	ModelMoxaE1242 = "MOXA_E1242" // 4进4出
	ModelMoxaE1210 = "MOXA_E1210" // 16进0出
	ModelMoxaE1211 = "MOXA_E1211" // 0进16出
)

const (
	ReadTypeDiscretes = "READ_TYPE_DISCRETES"
	ReadTypeCoils     = "READ_TYPE_COILS"
	//ReadTeypHoldingRegisters = "READ_TEYP_HOLDING_REGISTERS"
	//ReadTeypInputRegisters   = "READ_TEYP_INPUT_REGISTERS"

	WriteTypeSingleCoil = "WRITE_TYPE_SINGLE_COIL"
	//WriteTypeMultiCoil      = "WRITE_TYPE_MULTI_COIL"
	//WriteTypeSingleRegister = "WRITE_TYPE_SINGLE_REGISTER"
	//WriteTypeMultiRegister  = "WRITE_TYPE_MULTI_REGISTER"
)

type VendorCfg struct {
	InputNum      uint16
	InputAddress  uint16
	InputReadType string

	OutputNum      uint16
	OutputAddress  uint16
	OutputReadType string

	WriteType string
}

type Vendor interface {
	Type() string
	Cfg() VendorCfg
}

var VendorModels = map[string]Vendor{
	ModelMoxaE1212: &MOXA{model: ModelMoxaE1212},
	ModelMoxaE1242: &MOXA{model: ModelMoxaE1242},
	ModelMoxaE1210: &MOXA{model: ModelMoxaE1210},
	ModelMoxaE1211: &MOXA{model: ModelMoxaE1211},
}

type MOXA struct {
	model string
}

func (s *MOXA) Type() string {
	switch s.model {
	case ModelMoxaE1212:
		return IoModbustcp

	case ModelMoxaE1242:
		return IoModbustcp

	case ModelMoxaE1210:
		return IoModbustcp

	case ModelMoxaE1211:
		return IoModbustcp
	}

	return IoModbustcp
}

// inputNum, outputNum
func (s *MOXA) Cfg() VendorCfg {
	switch s.model {
	case ModelMoxaE1212:
		return VendorCfg{
			InputNum:       8,
			InputAddress:   0,
			InputReadType:  ReadTypeDiscretes,
			OutputNum:      8,
			OutputAddress:  0,
			OutputReadType: ReadTypeCoils,
			WriteType:      WriteTypeSingleCoil,
		}

	case ModelMoxaE1242:
		return VendorCfg{
			InputNum:       4,
			InputAddress:   0,
			InputReadType:  ReadTypeDiscretes,
			OutputNum:      4,
			OutputAddress:  0,
			OutputReadType: ReadTypeCoils,
			WriteType:      WriteTypeSingleCoil,
		}

	case ModelMoxaE1210:
		return VendorCfg{
			InputNum:       16,
			InputAddress:   0,
			InputReadType:  ReadTypeDiscretes,
			OutputNum:      0,
			OutputAddress:  0,
			OutputReadType: ReadTypeCoils,
			WriteType:      WriteTypeSingleCoil,
		}

	case ModelMoxaE1211:
		return VendorCfg{
			InputNum:       0,
			InputAddress:   0,
			InputReadType:  ReadTypeDiscretes,
			OutputNum:      16,
			OutputAddress:  0,
			OutputReadType: ReadTypeCoils,
			WriteType:      WriteTypeSingleCoil,
		}
	}

	return VendorCfg{}
}
