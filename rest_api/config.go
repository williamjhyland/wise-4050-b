package rest_api

import "errors"

type CloudConfig struct {
	DeviceIP string `json:"device_ip"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (conf *CloudConfig) Validate(path string) ([]string, error) {
	if conf.DeviceIP == "" {
		return nil, errors.New("device IP is required")
	}
	if conf.Username == "" {
		return nil, errors.New("username is required")
	}
	if conf.Password == "" {
		return nil, errors.New("password is required")
	}
	return nil, nil
}
