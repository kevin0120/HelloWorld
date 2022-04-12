package io

import "github.com/masami10/rush/services/device"

const (
	IoStatusOnline  = device.BaseDeviceStatusOnline
	IoStatusOffline = device.BaseDeviceStatusOffline

	IoTypeInput  = "input"
	IoTypeOutput = "output"

	IoModbustcp = "modbustcp"
)
