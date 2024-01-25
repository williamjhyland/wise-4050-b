package modbus_tcp

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/goburrow/modbus"
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var errUnimplemented = errors.New("unimplemented")
var Model = resource.NewModel("bill", "wise4050", "modbus")
var PrettyName = "WISE-4050 4DI/4DO 2.4G WiFi IoT Wireless I/O Module"
var Description = "WISE-4000 series is an Ethernet-based wired or wireless IoT device, which inte- grated with IoT data acquisition, processing, and publishing functions."

// Your sensor model type
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

// Get sensor reading
func (s *mySensor) Readings(ctx context.Context, _ map[string]interface{}) (map[string]interface{}, error) {
	s.RunModbusTCP()
	return map[string]interface{}{"setting": 0}, nil
}

func (s *mySensor) Close(ctx context.Context) error {
	s.logger.Infof("Shutting down %s", PrettyName)
	s.done <- true
	s.logger.Infof("Notifying monitor to shut down")
	s.wg.Wait()
	s.logger.Info("Monitor shut down")
	return nil
}

func (s *mySensor) RunModbusTCP() {
	// Connect to the Modbus device
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%s", s.sensorConfig.DeviceAddress, s.sensorConfig.Port)) // Replace with your device IP and port
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 1 // Set the Slave ID
	err := handler.Connect()
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	// Read DI Status from the correct address
	results, err := client.ReadCoils(0, 4) // Reading 4 register starting from address 301
	if err != nil {
		fmt.Printf("Failed to read: %v\n", err)
		return
	}
	ReadCoilValues(results)

	// Read DO Status from the correct address
	results, err = client.ReadCoils(16, 4) // Reading 4 register starting from address 301
	if err != nil {
		fmt.Printf("Failed to read: %v\n", err)
		return
	}
	ReadCoilValues(results)
}

func ReadCoilValues(byteVal []byte) {
	for i, byteVal := range byteVal {
		for bit := 0; bit < 8; bit++ {
			coilState := byteVal&(1<<bit) != 0
			fmt.Printf("Coil %d State: %t\n", i*8+bit, coilState)
			if i*8+bit >= 3 { // Assuming you only want to read the first 4 coils
				break
			}
		}
	}
}
