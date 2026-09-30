//go:build !darwin || cgo

package n2k

import "go.bug.st/serial/enumerator"

func serialDevices() ([]SerialDevice, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, err
	}
	devices := make([]SerialDevice, 0, len(ports))
	for _, port := range ports {
		if port == nil {
			continue
		}
		devices = append(devices, SerialDevice{
			Name:         port.Name,
			IsUSB:        port.IsUSB,
			VID:          port.VID,
			PID:          port.PID,
			SerialNumber: port.SerialNumber,
			Product:      port.Product,
		})
	}
	return devices, nil
}
