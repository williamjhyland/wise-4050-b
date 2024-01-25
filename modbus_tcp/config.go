package modbus_tcp

import "errors"

type CloudConfig struct {
	DeviceAddress string `json:"device_address"`
	Port          string `json:"port"`
	DI            *DI    `json:"digital_inputs"`
	DO            *DO    `json:"digital_outputs"`
}

// This is for Digital Input Parameters
type DI struct {
	BaseAddress int `json:"base_address"`
	Length      int `json:"length"`
}

// This is for Digital Input Parameters
type DO struct {
	BaseAddress int `json:"base_address"`
	Length      int `json:"length"`
}

func (conf *CloudConfig) Validate(path string) ([]string, error) {
	if conf.DI == nil {
		return nil, errors.New("digital_inputs are required")
	}

	if conf.DO == nil {
		return nil, errors.New("digital_outputs are required")
	}

	if conf.Port == "" {
		return nil, errors.New("Port is required")
	}

	if conf.DeviceAddress == "" {
		return nil, errors.New("DeviceAddress is required")
	}
	return nil, nil
}
