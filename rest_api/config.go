package rest_api

type RESTConfig struct {
    DeviceIP   string `json:"device_ip"`
    Username   string `json:"username"`
    Password   string `json:"password"`
}

func (conf *RESTConfig) Validate() error {
    if conf.DeviceIP == "" {
        return errors.New("device IP is required")
    }
    if conf.Username == "" {
        return errors.New("username is required")
    }
    if conf.Password == "" {
        return errors.New("password is required")
    }
    return nil
}
