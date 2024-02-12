package rest_api

import (
    "context"
    "encoding/base64"
    "errors"
    "fmt"
    "net/http"
    "sync"
)

var errUnimplemented = errors.New("unimplemented")
var Model = resource.NewModel("bill", "wise4050", "restapi")
var PrettyName = "WISE-4050 4DI/4DO 2.4G WiFi IoT Wireless I/O Module"
var Description = "WISE-4000 series is an Ethernet-based wired or wireless IoT device, which inte- grated with IoT data acquisition, processing, and publishing functions."

type mySensor struct {
    Config      *CloudConfig
    httpClient  *http.Client
    mu          sync.RWMutex
    cancelFunc  context.CancelFunc
}

func init() {
	resource.RegisterComponent(
		sensor.API,
		Model,
		resource.Registration[sensor.Sensor, *CloudConfig]{Constructor: NewSensor})
}

func NewSensor(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (sensor.Sensor, error) {
	logger.Infof("Starting %s %s", PrettyName)
	s := &mySensor{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}
	if err := s.Reconfigure(ctx, deps, conf); err != nil {
		return nil, err
	}
	return s, nil
}


func (s *mySensor) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	var err error
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sensorConfig, err = resource.NativeConfig[*CloudConfig](conf)
	if err != nil {
		return err
	}
	s.logger.Debugf("Reconfiguring %s", PrettyName)

	return err
}

func (s *RESTSensor) Readings(ctx context.Context) ([]byte, error) {
    url := fmt.Sprintf("http://%s/di_value/slot_0", s.Config.DeviceIP)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    auth := base64.StdEncoding.EncodeToString([]byte(s.Config.Username + ":" + s.Config.Password))
    req.Header.Add("Authorization", "Basic "+auth)

    resp, err := s.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("error: status code %d", resp.StatusCode)
    }

    // Process response body to return meaningful sensor readings
    // ...

    return nil, errors.New("unimplemented")
}

func (s *RESTSensor) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
    // Implement command functionality here
    // ...

    return nil, errors.New("unimplemented")
}

func (s *RESTSensor) Close(ctx context.Context) error {
    if s.cancelFunc != nil {
        s.cancelFunc()
    }
    return nil
}