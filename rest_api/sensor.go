package rest_api

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var errUnimplemented = errors.New("unimplemented")
var Model = resource.NewModel("bill", "advantech-wise-4050", "restapi")
var PrettyName = "WISE-4050 4DI/4DO 2.4G WiFi IoT Wireless I/O Module"
var Description = "WISE-4000 series is an Ethernet-based wired or wireless IoT device, which inte- grated with IoT data acquisition, processing, and publishing functions."

type mySensor struct {
	resource.Named
	logger       logging.Logger
	readings     string // the value to be returned by the sensor.readings() method
	mu           sync.RWMutex
	cancelCtx    context.Context
	cancelFunc   func()
	monitor      func()
	done         chan bool
	wg           sync.WaitGroup
	sensorConfig *CloudConfig
	client       *http.Client
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

func (s *mySensor) Readings(ctx context.Context, _ map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("http://%s/di_value/slot_0", s.sensorConfig.DeviceIP)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(s.sensorConfig.Username + ":" + s.sensorConfig.Password))
	req.Header.Add("Authorization", "Basic "+auth)

	resp, err := s.client.Do(req)
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

func (s *mySensor) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	// Implement command functionality here
	// ...

	return nil, errors.New("unimplemented")
}

func (s *mySensor) Close(ctx context.Context) error {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	return nil
}
