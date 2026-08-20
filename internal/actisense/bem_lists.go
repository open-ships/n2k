package actisense

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	supportedPGNListVariant = uint32(0x00001100)
	rxPGNListVariant        = uint32(0x00001101)
	txPGNListVariant        = uint32(0x00001102)
	proprietaryListVariant  = uint32(0x00001103)
	maxPGNListEntries       = 255
	maxProprietaryBitmap    = 32
	ProprietaryDP0Base      = uint32(0x0000FF00)
	ProprietaryDP1Base      = uint32(0x0001FF00)
)

type PGNEnableFlag uint8

const (
	PGNDisabled    PGNEnableFlag = 0
	PGNEnabled     PGNEnableFlag = 1
	PGNRespondMode PGNEnableFlag = 2
)

const (
	RxPGNMaskAcceptAll = uint32(0xFFFFFFFF)
	TxPGNRateDefault   = uint32(0xFFFFFFFF)
	TxPGNRateEvent     = uint32(0)
)

type RxPGNState struct {
	PGN  uint32
	Flag PGNEnableFlag
	Mask uint32
}

func RxPGNEnableGet(pgn uint32) []byte {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, pgn)
	return data
}

func RxPGNEnableSet(pgn uint32, flag PGNEnableFlag, mask *uint32) []byte {
	length := 5
	if mask != nil {
		length = 9
	}
	data := make([]byte, length)
	binary.LittleEndian.PutUint32(data[:4], pgn)
	data[4] = byte(flag)
	if mask != nil {
		binary.LittleEndian.PutUint32(data[5:9], *mask)
	}
	return data
}

func DecodeRxPGNState(response BEMResponse) (RxPGNState, error) {
	if response.BEMID != BEMRxPGNEnable {
		return RxPGNState{}, fmt.Errorf("actisense: expected Rx PGN response, got BEM 0x%02X", response.BEMID)
	}
	if len(response.Data) < 9 {
		return RxPGNState{}, fmt.Errorf("actisense: Rx PGN response is %d bytes; expected 9", len(response.Data))
	}
	return RxPGNState{
		PGN: binary.LittleEndian.Uint32(response.Data[:4]), Flag: PGNEnableFlag(response.Data[4]),
		Mask: binary.LittleEndian.Uint32(response.Data[5:9]),
	}, nil
}

func TxPGNEnableSetFull(pgn uint32, flag PGNEnableFlag, rate *uint32) []byte {
	length := 5
	if rate != nil {
		length = 9
	}
	data := make([]byte, length)
	binary.LittleEndian.PutUint32(data[:4], pgn)
	data[4] = byte(flag)
	if rate != nil {
		binary.LittleEndian.PutUint32(data[5:9], *rate)
	}
	return data
}

// PGNListSelector chooses the Rx list, Tx list, or both lists for delete and
// default operations.
type PGNListSelector uint8

const (
	PGNListRx   PGNListSelector = 0
	PGNListTx   PGNListSelector = 1
	PGNListBoth PGNListSelector = 2
)

func (s PGNListSelector) validate() error {
	if s > PGNListBoth {
		return fmt.Errorf("actisense: invalid PGN list selector %d", s)
	}
	return nil
}

type PGNListParameters struct {
	RxMaximum uint16
	RxSession uint16
	RxActive  uint16
	TxMaximum uint16
	TxSession uint16
	TxActive  uint16
	RxSync    uint8
	TxSync    uint8
}

func (p PGNListParameters) RxSynchronized() bool { return p.RxSync == 0 }
func (p PGNListParameters) TxSynchronized() bool { return p.TxSync == 0 }
func (p PGNListParameters) Synchronized() bool {
	return p.RxSynchronized() && p.TxSynchronized()
}

func DecodePGNListParameters(response BEMResponse) (PGNListParameters, error) {
	if response.BEMID != BEMPGNListParameters {
		return PGNListParameters{}, fmt.Errorf("actisense: expected PGN List Parameters response, got BEM 0x%02X", response.BEMID)
	}
	if len(response.Data) < 14 {
		return PGNListParameters{}, fmt.Errorf("actisense: PGN List Parameters response is %d bytes; expected 14", len(response.Data))
	}
	return PGNListParameters{
		RxMaximum: binary.LittleEndian.Uint16(response.Data[0:2]),
		RxSession: binary.LittleEndian.Uint16(response.Data[2:4]),
		RxActive:  binary.LittleEndian.Uint16(response.Data[4:6]),
		TxMaximum: binary.LittleEndian.Uint16(response.Data[6:8]),
		TxSession: binary.LittleEndian.Uint16(response.Data[8:10]),
		TxActive:  binary.LittleEndian.Uint16(response.Data[10:12]),
		RxSync:    response.Data[12], TxSync: response.Data[13],
	}, nil
}

type SupportedPGN struct {
	Index uint8
	PGN   uint32
}

type SupportedPGNList struct {
	TransferID      uint8
	DatabaseVersion uint16
	Entries         []SupportedPGN
}

type supportedPGNPart struct {
	transferID uint8
	database   uint16
	total      uint8
	first      uint8
	entries    []SupportedPGN
}

func decodeSupportedPGNPart(response BEMResponse) (supportedPGNPart, error) {
	if response.BEMID != BEMSupportedPGNList {
		return supportedPGNPart{}, fmt.Errorf("actisense: expected Supported PGN List response, got BEM 0x%02X", response.BEMID)
	}
	data := response.Data
	if len(data) < 10 {
		return supportedPGNPart{}, fmt.Errorf("actisense: Supported PGN List response is %d bytes; expected at least 10", len(data))
	}
	if variant := binary.LittleEndian.Uint32(data[1:5]); variant != supportedPGNListVariant {
		return supportedPGNPart{}, fmt.Errorf("actisense: Supported PGN List variant is 0x%08X; expected 0x%08X", variant, supportedPGNListVariant)
	}
	part := supportedPGNPart{transferID: data[0], database: binary.LittleEndian.Uint16(data[5:7]), total: data[7], first: data[8]}
	count := int(data[9])
	needed := 10 + count*4
	if len(data) < needed {
		return supportedPGNPart{}, fmt.Errorf("actisense: Supported PGN List response needs %d bytes for %d entries; got %d", needed, count, len(data))
	}
	part.entries = make([]SupportedPGN, 0, count)
	for index := 0; index < count; index++ {
		offset := 10 + index*4
		part.entries = append(part.entries, SupportedPGN{
			Index: data[offset],
			PGN:   uint32(data[offset+1]) | uint32(data[offset+2])<<8 | uint32(data[offset+3])<<16,
		})
	}
	return part, nil
}

type supportedPGNAccumulator struct {
	result      SupportedPGNList
	seen        []bool
	initialized bool
	received    int
}

func (a *supportedPGNAccumulator) feed(part supportedPGNPart) (bool, error) {
	if !a.initialized {
		a.result.TransferID = part.transferID
		a.result.DatabaseVersion = part.database
		a.result.Entries = make([]SupportedPGN, int(part.total))
		a.seen = make([]bool, int(part.total))
		a.initialized = true
	} else if part.transferID != a.result.TransferID {
		return false, fmt.Errorf("actisense: Supported PGN List transfer ID changed from %d to %d", a.result.TransferID, part.transferID)
	} else if int(part.total) != len(a.result.Entries) {
		return false, fmt.Errorf("actisense: Supported PGN List total changed from %d to %d", len(a.result.Entries), part.total)
	} else if part.database != a.result.DatabaseVersion {
		return false, fmt.Errorf("actisense: Supported PGN List database version changed from %d to %d", a.result.DatabaseVersion, part.database)
	}
	end := int(part.first) + len(part.entries)
	if end > len(a.result.Entries) {
		return false, fmt.Errorf("actisense: Supported PGN List range %d:%d exceeds total %d", part.first, end, len(a.result.Entries))
	}
	for index, entry := range part.entries {
		slot := int(part.first) + index
		a.result.Entries[slot] = entry
		if !a.seen[slot] {
			a.seen[slot] = true
			a.received++
		}
	}
	return a.received == len(a.result.Entries), nil
}

type RxPGNListEntry struct {
	Index uint8
	Mask  uint8
}

type TxPGNListEntry struct {
	Index    uint8
	Priority uint8
	RateMS   uint16
}

type ProprietaryPGNList struct {
	DP0Bitmap   [maxProprietaryBitmap]byte
	DP1Bitmap   [maxProprietaryBitmap]byte
	EnabledPGNs []uint32
}

type RxPGNEnableList struct {
	TransferID         uint8
	Entries            []RxPGNListEntry
	Proprietary        ProprietaryPGNList
	ProprietaryPresent bool
}

type TxPGNEnableList struct {
	TransferID         uint8
	Entries            []TxPGNListEntry
	Proprietary        ProprietaryPGNList
	ProprietaryPresent bool
}

type f2Variant uint8

const (
	f2Standard f2Variant = iota
	f2Proprietary
)

type rxF2Part struct {
	transferID uint8
	variant    f2Variant
	total      uint8
	first      uint8
	entries    []RxPGNListEntry
	dp0        []byte
	dp1        []byte
}

func decodeRxF2Part(response BEMResponse) (rxF2Part, error) {
	if response.BEMID != BEMRxPGNEnableListF2 {
		return rxF2Part{}, fmt.Errorf("actisense: expected Rx PGN List F2 response, got BEM 0x%02X", response.BEMID)
	}
	data := response.Data
	if len(data) < 5 {
		return rxF2Part{}, fmt.Errorf("actisense: Rx PGN List F2 response is %d bytes; expected at least 5", len(data))
	}
	part := rxF2Part{transferID: data[0]}
	switch variant := binary.LittleEndian.Uint32(data[1:5]); variant {
	case rxPGNListVariant:
		part.variant = f2Standard
		if len(data) < 8 {
			return rxF2Part{}, errors.New("actisense: Rx PGN List F2 standard response is shorter than 8 bytes")
		}
		part.total, part.first = data[5], data[6]
		count := int(data[7])
		if len(data) < 8+count*2 {
			return rxF2Part{}, fmt.Errorf("actisense: Rx PGN List F2 standard response is truncated for %d entries", count)
		}
		part.entries = make([]RxPGNListEntry, count)
		for index := range count {
			part.entries[index] = RxPGNListEntry{Index: data[8+index*2], Mask: data[9+index*2]}
		}
	case proprietaryListVariant:
		part.variant = f2Proprietary
		var err error
		part.dp0, part.dp1, err = decodeProprietaryBitmaps(data[5:])
		if err != nil {
			return rxF2Part{}, fmt.Errorf("actisense: Rx PGN List F2: %w", err)
		}
	default:
		return rxF2Part{}, fmt.Errorf("actisense: Rx PGN List F2 variant 0x%08X is unsupported", variant)
	}
	return part, nil
}

type txF2Part struct {
	transferID uint8
	variant    f2Variant
	total      uint8
	first      uint8
	entries    []TxPGNListEntry
	dp0        []byte
	dp1        []byte
}

func decodeTxF2Part(response BEMResponse) (txF2Part, error) {
	if response.BEMID != BEMTxPGNEnableListF2 {
		return txF2Part{}, fmt.Errorf("actisense: expected Tx PGN List F2 response, got BEM 0x%02X", response.BEMID)
	}
	data := response.Data
	if len(data) < 5 {
		return txF2Part{}, fmt.Errorf("actisense: Tx PGN List F2 response is %d bytes; expected at least 5", len(data))
	}
	part := txF2Part{transferID: data[0]}
	switch variant := binary.LittleEndian.Uint32(data[1:5]); variant {
	case txPGNListVariant:
		part.variant = f2Standard
		if len(data) < 8 {
			return txF2Part{}, errors.New("actisense: Tx PGN List F2 standard response is shorter than 8 bytes")
		}
		part.total, part.first = data[5], data[6]
		count := int(data[7])
		if len(data) < 8+count*4 {
			return txF2Part{}, fmt.Errorf("actisense: Tx PGN List F2 standard response is truncated for %d entries", count)
		}
		part.entries = make([]TxPGNListEntry, count)
		for index := range count {
			offset := 8 + index*4
			part.entries[index] = TxPGNListEntry{Index: data[offset], Priority: data[offset+1], RateMS: binary.LittleEndian.Uint16(data[offset+2 : offset+4])}
		}
	case proprietaryListVariant:
		part.variant = f2Proprietary
		var err error
		part.dp0, part.dp1, err = decodeProprietaryBitmaps(data[5:])
		if err != nil {
			return txF2Part{}, fmt.Errorf("actisense: Tx PGN List F2: %w", err)
		}
	default:
		return txF2Part{}, fmt.Errorf("actisense: Tx PGN List F2 variant 0x%08X is unsupported", variant)
	}
	return part, nil
}

func decodeProprietaryBitmaps(data []byte) ([]byte, []byte, error) {
	if len(data) < 1 {
		return nil, nil, errors.New("proprietary response has no DP0 length")
	}
	dp0Length := int(data[0])
	if dp0Length > maxProprietaryBitmap || len(data) < 1+dp0Length+1 {
		return nil, nil, fmt.Errorf("proprietary DP0 length %d is invalid for %d bytes", dp0Length, len(data))
	}
	dp0 := append([]byte(nil), data[1:1+dp0Length]...)
	offset := 1 + dp0Length
	dp1Length := int(data[offset])
	offset++
	if dp1Length > maxProprietaryBitmap || len(data) < offset+dp1Length {
		return nil, nil, fmt.Errorf("proprietary DP1 length %d is invalid for %d remaining bytes", dp1Length, len(data)-offset)
	}
	return dp0, append([]byte(nil), data[offset:offset+dp1Length]...), nil
}

func fillProprietaryList(result *ProprietaryPGNList, dp0, dp1 []byte) {
	copy(result.DP0Bitmap[:], dp0)
	copy(result.DP1Bitmap[:], dp1)
	result.EnabledPGNs = expandBitmap(dp0, ProprietaryDP0Base)
	result.EnabledPGNs = append(result.EnabledPGNs, expandBitmap(dp1, ProprietaryDP1Base)...)
}

func expandBitmap(bitmap []byte, base uint32) []uint32 {
	result := make([]uint32, 0)
	for index, value := range bitmap {
		for bit := uint8(0); bit < 8; bit++ {
			if value&(1<<bit) != 0 {
				result = append(result, base+uint32(index*8)+uint32(bit))
			}
		}
	}
	return result
}

type rxF2Accumulator struct {
	result            RxPGNEnableList
	seen              []bool
	standardReceived  int
	standardSeen      bool
	initialized       bool
	expectProprietary bool
}

func (a *rxF2Accumulator) feed(response BEMResponse) (bool, error) {
	part, err := decodeRxF2Part(response)
	if err != nil {
		return false, err
	}
	if !a.initialized {
		a.result.TransferID = part.transferID
		a.initialized = true
	} else if part.transferID != a.result.TransferID {
		return false, fmt.Errorf("actisense: Rx PGN List F2 transfer ID changed from %d to %d", a.result.TransferID, part.transferID)
	}
	if part.variant == f2Proprietary {
		fillProprietaryList(&a.result.Proprietary, part.dp0, part.dp1)
		a.result.ProprietaryPresent = true
		return a.complete(), nil
	}
	if !a.standardSeen {
		a.result.Entries = make([]RxPGNListEntry, int(part.total))
		a.seen = make([]bool, int(part.total))
		a.standardSeen = true
	} else if int(part.total) != len(a.result.Entries) {
		return false, fmt.Errorf("actisense: Rx PGN List F2 total changed from %d to %d", len(a.result.Entries), part.total)
	}
	end := int(part.first) + len(part.entries)
	if end > len(a.result.Entries) {
		return false, fmt.Errorf("actisense: Rx PGN List F2 range %d:%d exceeds total %d", part.first, end, len(a.result.Entries))
	}
	for index, entry := range part.entries {
		slot := int(part.first) + index
		a.result.Entries[slot] = entry
		if !a.seen[slot] {
			a.seen[slot] = true
			a.standardReceived++
		}
	}
	return a.complete(), nil
}

func (a *rxF2Accumulator) complete() bool {
	standardComplete := a.standardSeen && a.standardReceived == len(a.result.Entries)
	return standardComplete && (!a.expectProprietary || a.result.ProprietaryPresent)
}

type txF2Accumulator struct {
	result            TxPGNEnableList
	seen              []bool
	standardReceived  int
	standardSeen      bool
	initialized       bool
	expectProprietary bool
}

func (a *txF2Accumulator) feed(response BEMResponse) (bool, error) {
	part, err := decodeTxF2Part(response)
	if err != nil {
		return false, err
	}
	if !a.initialized {
		a.result.TransferID = part.transferID
		a.initialized = true
	} else if part.transferID != a.result.TransferID {
		return false, fmt.Errorf("actisense: Tx PGN List F2 transfer ID changed from %d to %d", a.result.TransferID, part.transferID)
	}
	if part.variant == f2Proprietary {
		fillProprietaryList(&a.result.Proprietary, part.dp0, part.dp1)
		a.result.ProprietaryPresent = true
		return a.complete(), nil
	}
	if !a.standardSeen {
		a.result.Entries = make([]TxPGNListEntry, int(part.total))
		a.seen = make([]bool, int(part.total))
		a.standardSeen = true
	} else if int(part.total) != len(a.result.Entries) {
		return false, fmt.Errorf("actisense: Tx PGN List F2 total changed from %d to %d", len(a.result.Entries), part.total)
	}
	end := int(part.first) + len(part.entries)
	if end > len(a.result.Entries) {
		return false, fmt.Errorf("actisense: Tx PGN List F2 range %d:%d exceeds total %d", part.first, end, len(a.result.Entries))
	}
	for index, entry := range part.entries {
		slot := int(part.first) + index
		a.result.Entries[slot] = entry
		if !a.seen[slot] {
			a.seen[slot] = true
			a.standardReceived++
		}
	}
	return a.complete(), nil
}

func (a *txF2Accumulator) complete() bool {
	standardComplete := a.standardSeen && a.standardReceived == len(a.result.Entries)
	return standardComplete && (!a.expectProprietary || a.result.ProprietaryPresent)
}
