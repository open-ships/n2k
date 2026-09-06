package actisense

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	BEMPortDuplicateDelete = 0x14
	BEMRxPGNEnableListF1   = 0x48
	BEMTxPGNEnableListF1   = 0x49
)

type PortDuplicateDelete uint8

const (
	PortDuplicateDeleteOff      PortDuplicateDelete = 0
	PortDuplicateDeleteOn       PortDuplicateDelete = 1
	PortDuplicateDeleteNoChange PortDuplicateDelete = 255
)

func decodePortDuplicateDelete(response BEMResponse) ([]PortDuplicateDelete, error) {
	if response.BEMID != BEMPortDuplicateDelete || len(response.Data) < 1 || int(response.Data[0])+1 != len(response.Data) || response.Data[0] > 223 {
		return nil, errors.New("actisense: invalid Port Duplicate Delete response length or command")
	}
	settings := make([]PortDuplicateDelete, int(response.Data[0]))
	for i, value := range response.Data[1:] {
		if value > 1 {
			return nil, fmt.Errorf("actisense: invalid duplicate delete state %d for port %d", value, i)
		}
		settings[i] = PortDuplicateDelete(value)
	}
	return settings, nil
}

// GetPortDuplicateDelete reads the documented BEM-14 port filter settings.
// Support depends on device firmware; unsupported devices return a BEM error.
func (c *CommandSet) GetPortDuplicateDelete(ctx context.Context) ([]PortDuplicateDelete, error) {
	response, err := c.request(ctx, BEMPortDuplicateDelete, nil)
	if err != nil {
		return nil, err
	}
	return decodePortDuplicateDelete(response)
}

// SetPortDuplicateDelete supplies every port, using NoChange for unchanged
// entries. Firmware applies and persists this operation immediately.
func (c *CommandSet) SetPortDuplicateDelete(ctx context.Context, settings []PortDuplicateDelete) ([]PortDuplicateDelete, error) {
	if len(settings) == 0 || len(settings) > 223 {
		return nil, errors.New("actisense: duplicate delete requires 1-223 port settings")
	}
	data := make([]byte, len(settings))
	for i, value := range settings {
		if value != PortDuplicateDeleteOff && value != PortDuplicateDeleteOn && value != PortDuplicateDeleteNoChange {
			return nil, fmt.Errorf("actisense: invalid duplicate delete setting %d for port %d", value, i)
		}
		data[i] = byte(value)
	}
	response, err := c.request(ctx, BEMPortDuplicateDelete, data)
	if err != nil {
		return nil, err
	}
	accepted, err := decodePortDuplicateDelete(response)
	if err != nil {
		return nil, err
	}
	if len(accepted) != len(settings) {
		return accepted, errors.New("actisense: duplicate delete acknowledgement has a different port count")
	}
	for i, value := range data {
		if value != byte(PortDuplicateDeleteNoChange) && byte(accepted[i]) != value {
			return accepted, fmt.Errorf("actisense: duplicate delete setting for port %d was not accepted", i)
		}
	}
	return accepted, nil
}

type RxPGNListF1Entry struct {
	PGN  uint32
	Mask uint32
}

type TxPGNListF1Entry struct {
	PGN       uint32
	RateMS    uint32
	TimeoutMS uint32
	Priority  uint8
}

// PartsReceived identifies which fields of a partial F1 result are populated.
type RxPGNEnableListF1 struct {
	PartsReceived int
	Entries       []RxPGNListF1Entry
}

type TxPGNEnableListF1 struct {
	PartsReceived int
	Entries       []TxPGNListF1Entry
}

type f1Accumulator struct {
	command byte
	parts   int
	count   int
	model   uint16
	serial  uint32
	origin  BEMOrigin
}

func (a *f1Accumulator) feed(response BEMResponse, width int) ([]uint32, error) {
	if response.BEMID != a.command || len(response.Data) < 1 {
		return nil, errors.New("actisense: invalid F1 response command or missing count")
	}
	count := int(response.Data[0])
	if count > 50 || len(response.Data) != 1+width*count {
		return nil, errors.New("actisense: F1 count exceeds 50 or response length is inconsistent")
	}
	if a.parts > 0 && (count != a.count || response.ModelID != a.model || response.SerialNumber != a.serial || response.Origin != a.origin) {
		return nil, errors.New("actisense: F1 response count or device identity changed during transfer")
	}
	values := make([]uint32, count)
	for i := range count {
		if width == 1 {
			values[i] = uint32(response.Data[1+i])
		} else {
			values[i] = binary.LittleEndian.Uint32(response.Data[1+4*i:])
		}
	}
	a.count, a.model, a.serial, a.origin = count, response.ModelID, response.SerialNumber, response.Origin
	a.parts++
	return values, nil
}

// GetRxPGNEnableListF1 reads the two ordered legacy response parts. F1 was
// discontinued at firmware 2.500 and is not implemented by NGX/W2K products.
// Prefer GetRxPGNEnableList. No automatic fallback or list mutation occurs.
func (c *CommandSet) GetRxPGNEnableListF1(ctx context.Context) (RxPGNEnableListF1, error) {
	a := f1Accumulator{command: BEMRxPGNEnableListF1}
	result := RxPGNEnableListF1{}
	_, err := c.requestMulti(ctx, a.command, nil, func(response BEMResponse) (bool, error) {
		values, err := a.feed(response, 4)
		if err != nil {
			return false, err
		}
		for _, value := range values {
			if a.parts == 1 {
				if err := validatePGN(value); err != nil {
					return false, err
				}
			} else if !validRxPGNMask(value) {
				return false, fmt.Errorf("actisense: invalid F1 Rx PGN mask 0x%08X", value)
			}
		}
		if a.parts == 1 {
			result.Entries = make([]RxPGNListF1Entry, a.count)
		}
		for i, value := range values {
			if a.parts == 1 {
				result.Entries[i].PGN = value
			} else {
				result.Entries[i].Mask = value
			}
		}
		result.PartsReceived = a.parts
		return a.parts == 2, nil
	})
	return result, err
}

// GetTxPGNEnableListF1 reads PGNs, rates, timeouts, and priorities in order.
// Prefer GetTxPGNEnableList on current firmware. Partial results retain the
// completed parts when a response train fails or times out.
func (c *CommandSet) GetTxPGNEnableListF1(ctx context.Context) (TxPGNEnableListF1, error) {
	a := f1Accumulator{command: BEMTxPGNEnableListF1}
	result := TxPGNEnableListF1{}
	_, err := c.requestMulti(ctx, a.command, nil, func(response BEMResponse) (bool, error) {
		width := 4
		if a.parts == 3 {
			width = 1
		}
		values, err := a.feed(response, width)
		if err != nil {
			return false, err
		}
		for _, value := range values {
			if a.parts == 1 {
				if err := validatePGN(value); err != nil {
					return false, err
				}
			} else if a.parts == 4 && value > 7 {
				return false, fmt.Errorf("actisense: invalid F1 Tx priority %d", value)
			}
		}
		if a.parts == 1 {
			result.Entries = make([]TxPGNListF1Entry, a.count)
		}
		for i, value := range values {
			switch a.parts {
			case 1:
				result.Entries[i].PGN = value
			case 2:
				result.Entries[i].RateMS = value
			case 3:
				result.Entries[i].TimeoutMS = value
			case 4:
				result.Entries[i].Priority = uint8(value)
			}
		}
		result.PartsReceived = a.parts
		return a.parts == 4, nil
	})
	return result, err
}
