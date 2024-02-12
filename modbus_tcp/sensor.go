package modbus_tcp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/goburrow/modbus"
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var errUnimplemented = errors.New("unimplemented")
var Model = resource.NewModel("bill", "advantech-wise-4050", "modbus")
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
	readings, err := s.RunModbusTCP()
	if err != nil {
		return nil, err
	}
	return readings, nil
}

// DoCommand can be implemented to extend sensor functionality but returns unimplemented in this example.
func (s *mySensor) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	fmt.Printf("DOING!")
	// Connect to the Modbus device
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%s", s.sensorConfig.DeviceAddress, s.sensorConfig.Port))
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 1
	err := handler.Connect()
	defer handler.Close()
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return nil, err
	}
	fmt.Printf("DOING!")

	client := modbus.NewClient(handler)

	// Iterate through the command map and write to each specified coil
	for i := 1; i <= s.sensorConfig.DO.Length; i++ {
		coilKey := fmt.Sprintf("coil%d", i)
		if coilValue, ok := cmd[coilKey]; ok {
			coilAddr := uint16(s.sensorConfig.DO.BaseAddress + i - 1)
			var value uint16
			if coilValue.(bool) {
				value = 0xFF00 // ON value for Modbus coil
			} else {
				value = 0x0000 // OFF value for Modbus coil
			}
			_, err := client.WriteSingleCoil(coilAddr, value)
			if err != nil {
				return nil, err
			}
		}
	}

	return map[string]interface{}{"status": "success"}, nil
}

// The close method is executed when the component is shut down
func (s *mySensor) Close(ctx context.Context) error {
	s.logger.Infof("Shutting down %s", PrettyName)
	return nil
}

func (s *mySensor) RunModbusTCP() (map[string]interface{}, error) {
	// Connect to the Modbus device
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%s", s.sensorConfig.DeviceAddress, s.sensorConfig.Port)) // Replace with your device IP and port
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 1 // Set the Slave ID
	err := handler.Connect()
	defer handler.Close()
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return nil, err
	}

	client := modbus.NewClient(handler)
	result := make(map[string]interface{})

	// Read DI Status from the correct address
	diCoilStates, err := readCoilStates(client, s.sensorConfig.DI.BaseAddress, s.sensorConfig.DI.Length) // Adjust starting address and length
	if err != nil {
		return nil, err
	}

	fmt.Printf(strings.Join(diCoilStates, ", "))
	result["inputCoils"] = strings.Join(diCoilStates, ", ")

	// Read DO Status from the correct address
	doCoilStates, err := readCoilStates(client, s.sensorConfig.DO.BaseAddress, s.sensorConfig.DO.Length) // Adjust starting address and length
	if err != nil {
		return nil, err
	}
	fmt.Printf(strings.Join(doCoilStates, ", "))
	result["outputCoils"] = strings.Join(doCoilStates, ", ")

	return result, nil
}

func readCoilStates(client modbus.Client, startAddr, numCoils int) ([]string, error) {
	uStartAddr := uint16(startAddr)
	uNumCoils := uint16(numCoils)
	results, err := client.ReadCoils(uStartAddr, uNumCoils)
	if err != nil {
		fmt.Printf("Failed to read: %v\n", err)
		return nil, err
	}

	return decodeCoilValues(results, numCoils), nil
}

func decodeCoilValues(byteVal []byte, numCoils int) []string {
	coilStates := make([]string, 0)
	coilsProcessed := 0
	for i, val := range byteVal {
		for bit := 0; bit < 8; bit++ {
			if coilsProcessed >= numCoils {
				break // Exit if we've read all the coils we need
			}
			coilState := val&(1<<bit) != 0
			stateStr := "False"
			if coilState {
				stateStr = "True"
			}
			coilStates = append(coilStates, fmt.Sprintf("Coil %d: %s", i*8+bit+1, stateStr))
			coilsProcessed++
		}
	}
	return coilStates
}
