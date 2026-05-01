package pgn

import (
	"fmt"
)

type RadioFrequencyModePower struct {
	Info MessageInfo `json:"info"`
	RxFrequency *float32 `json:"rxFrequency"`
	TxFrequency *float32 `json:"txFrequency"`
	RadioChannel string `json:"radioChannel"`
	TxPower *uint8 `json:"txPower"`
	Mode *uint16 `json:"mode"`
	ChannelBandwidth *uint16 `json:"channelBandwidth"`
}
func DecodeRadioFrequencyModePower(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val RadioFrequencyModePower
	val.Info = Info
	if v, err := stream.readUnsignedResolution(32, 10); err != nil {
		return nil, fmt.Errorf("parse failed for RadioFrequencyModePower-RxFrequency: %w", err)
	} else {
		val.RxFrequency = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 10); err != nil {
		return nil, fmt.Errorf("parse failed for RadioFrequencyModePower-TxFrequency: %w", err)
	} else {
		val.TxFrequency = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(48); err != nil {
		return nil, fmt.Errorf("parse failed for RadioFrequencyModePower-RadioChannel: %w", err)
	} else {
		val.RadioChannel = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for RadioFrequencyModePower-TxPower: %w", err)
	} else {
		val.TxPower = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for RadioFrequencyModePower-Mode: %w", err)
	} else {
		val.Mode = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for RadioFrequencyModePower-ChannelBandwidth: %w", err)
	} else {
		val.ChannelBandwidth = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeRadioFrequencyModePower(val *RadioFrequencyModePower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUnsignedResolution(val.RxFrequency, 32, 10)
	w.writeUnsignedResolution(val.TxFrequency, 32, 10)
	w.writeFixedString(val.RadioChannel, 48)
	w.writeUInt8(val.TxPower, 8)
	w.writeUInt16(val.Mode, 16)
	w.writeUInt16(val.ChannelBandwidth, 16)
	return w.Bytes(), w.Err()
}
func encodeRadioFrequencyModePowerAny(v any) ([]byte, error) {
	val, ok := v.(*RadioFrequencyModePower)
	if !ok {
		return nil, fmt.Errorf("expected *RadioFrequencyModePower, got %T", v)
	}
	return EncodeRadioFrequencyModePower(val)
}
