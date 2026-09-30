//go:build darwin && !cgo

package n2k

import "go.bug.st/serial"

func serialDevices() ([]SerialDevice, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, err
	}
	devices := make([]SerialDevice, 0, len(ports))
	for _, port := range ports {
		devices = append(devices, SerialDevice{Name: port})
	}
	return devices, nil
}
