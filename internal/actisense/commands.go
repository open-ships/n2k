package actisense

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	defaultCommandTimeout  = 5 * time.Second
	defaultMultiInactivity = 500 * time.Millisecond
	maxRemoteEchoPayload   = 206
)

// Requester is the envelope-independent BEM request seam. Session implements
// it for local BST-A1/A0 traffic; the public package supplies a PGN-126720
// implementation for remote devices.
type Requester interface {
	Request(context.Context, byte, []byte) (BEMResponse, error)
	RequestMulti(context.Context, byte, []byte, time.Duration, func([]BEMResponse) (bool, error)) ([]BEMResponse, error)
}

type CommandSetConfig struct {
	Timeout         time.Duration
	MultiInactivity time.Duration
	Remote          bool
}

// CommandSet is the typed, envelope-independent Actisense BEM Interface.
// Persistent and rebooting verbs are methods but are never called implicitly.
type CommandSet struct {
	requester  Requester
	timeout    time.Duration
	inactivity time.Duration
	remote     bool

	mu         sync.RWMutex
	lastModel  ModelID
	lastSerial uint32
}

func NewCommandSet(requester Requester, config CommandSetConfig) *CommandSet {
	if config.Timeout <= 0 {
		config.Timeout = defaultCommandTimeout
	}
	if config.MultiInactivity <= 0 {
		config.MultiInactivity = defaultMultiInactivity
	}
	return &CommandSet{requester: requester, timeout: config.Timeout, inactivity: config.MultiInactivity, remote: config.Remote}
}

func (c *CommandSet) requestContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, c.timeout)
}

func (c *CommandSet) remember(response BEMResponse) {
	c.mu.Lock()
	c.lastModel = ModelID(response.ModelID)
	c.lastSerial = response.SerialNumber
	c.mu.Unlock()
}

// DeviceCapabilities returns the last model identity observed on this handle.
// It is Unknown until at least one response has arrived.
func (c *CommandSet) DeviceCapabilities() DeviceCapabilities {
	c.mu.RLock()
	model := c.lastModel
	c.mu.RUnlock()
	return CapabilitiesForModel(model)
}

func (c *CommandSet) DeviceSerialNumber() uint32 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastSerial
}

func (c *CommandSet) request(ctx context.Context, command byte, data []byte) (BEMResponse, error) {
	if c == nil || c.requester == nil {
		return BEMResponse{}, errors.New("actisense: command handle is closed")
	}
	requestCtx, cancel := c.requestContext(ctx)
	defer cancel()
	response, err := c.requester.Request(requestCtx, command, data)
	if response.BEMID != 0 {
		c.remember(response)
	}
	return response, err
}

func (c *CommandSet) requestMulti(ctx context.Context, command byte, data []byte, complete func(BEMResponse) (bool, error)) ([]BEMResponse, error) {
	if c == nil || c.requester == nil {
		return nil, errors.New("actisense: command handle is closed")
	}
	requestCtx, cancel := c.requestContext(ctx)
	defer cancel()
	responses, err := c.requester.RequestMulti(requestCtx, command, data, c.inactivity, func(responses []BEMResponse) (bool, error) {
		response := responses[len(responses)-1]
		c.remember(response)
		return complete(response)
	})
	return responses, err
}

// RawRequest sends one BEM command through this handle's local or remote
// envelope. Prefer typed methods when a command is known.
func (c *CommandSet) RawRequest(ctx context.Context, command byte, data []byte) (BEMResponse, error) {
	return c.request(ctx, command, data)
}

// RawRequestMulti collects a bounded response train for a caller-defined BEM
// command. complete receives each response in arrival order; returning true
// ends the train. On failure, all received responses are returned with the error.
func (c *CommandSet) RawRequestMulti(ctx context.Context, command byte, data []byte, complete func(BEMResponse) (bool, error)) ([]BEMResponse, error) {
	if complete == nil {
		return nil, errors.New("actisense: a multi-response request requires a completion function")
	}
	return c.requestMulti(ctx, command, data, complete)
}

func (c *CommandSet) GetOperatingMode(ctx context.Context) (OperatingMode, error) {
	response, err := c.request(ctx, BEMOperatingMode, nil)
	if err != nil {
		return 0, err
	}
	return DecodeOperatingMode(response)
}

func (c *CommandSet) SetOperatingMode(ctx context.Context, mode OperatingMode) error {
	response, err := c.request(ctx, BEMOperatingMode, OperatingModeSet(mode))
	if err != nil {
		return err
	}
	accepted, err := DecodeOperatingMode(response)
	if err != nil {
		return err
	}
	if accepted != mode {
		return fmt.Errorf("actisense: device acknowledged operating mode %d; requested %d", accepted, mode)
	}
	return nil
}

func (c *CommandSet) GetProductInfo(ctx context.Context) (ProductInfo, error) {
	accumulator := &productInfoAccumulator{}
	_, err := c.requestMulti(ctx, BEMProductInfo, nil, accumulator.feed)
	return accumulator.result, err
}

func (c *CommandSet) GetPortInventory(ctx context.Context) (PortInventory, error) {
	accumulator := &portInventoryAccumulator{}
	_, err := c.requestMulti(ctx, BEMPortInventory, nil, accumulator.feed)
	return accumulator.result, err
}

func (c *CommandSet) GetPortBaudrate(ctx context.Context, port uint8) (PortBaudrate, error) {
	response, err := c.request(ctx, BEMPortBaudrate, PortBaudrateGet(port))
	if err != nil {
		return PortBaudrate{}, err
	}
	return DecodePortBaudrate(response)
}

// SetPortBaudrate changes the session and stored rates explicitly. Passing a
// sentinel such as BaudRateNoChange preserves the corresponding value. The
// stored value may be written to non-volatile memory by device firmware.
func (c *CommandSet) SetPortBaudrate(ctx context.Context, port uint8, sessionBaud, storeBaud uint32) (PortBaudrate, error) {
	response, err := c.request(ctx, BEMPortBaudrate, PortBaudrateSet(port, sessionBaud, storeBaud))
	if err != nil {
		return PortBaudrate{}, err
	}
	state, err := DecodePortBaudrate(response)
	if err != nil {
		return PortBaudrate{}, err
	}
	if state.PortNumber != port {
		return state, fmt.Errorf("actisense: device acknowledged Port Baudrate port %d; requested %d", state.PortNumber, port)
	}
	return state, nil
}

func (c *CommandSet) GetPortPCodes(ctx context.Context) ([]PortPCode, error) {
	response, err := c.request(ctx, BEMPortPCode, nil)
	if err != nil {
		return nil, err
	}
	return DecodePortPCodes(response)
}

func (c *CommandSet) SetPortPCodes(ctx context.Context, codes []PortPCode) ([]PortPCode, error) {
	data, err := EncodePortPCodes(codes)
	if err != nil {
		return nil, err
	}
	response, err := c.request(ctx, BEMPortPCode, data)
	if err != nil {
		return nil, err
	}
	return DecodePortPCodes(response)
}

func (c *CommandSet) GetCANConfig(ctx context.Context) (CANConfig, error) {
	response, err := c.request(ctx, BEMCANConfig, nil)
	if err != nil {
		return CANConfig{}, err
	}
	return DecodeCANConfig(response)
}

// SetCANConfig changes the device's NMEA 2000 NAME and preferred address.
// Some firmware persists this operation itself; the library never calls it
// automatically.
func (c *CommandSet) SetCANConfig(ctx context.Context, config CANConfig) (CANConfig, error) {
	response, err := c.request(ctx, BEMCANConfig, CANConfigSet(config))
	if err != nil {
		return CANConfig{}, err
	}
	accepted, err := DecodeCANConfig(response)
	if err != nil {
		return CANConfig{}, err
	}
	if accepted != config {
		return accepted, fmt.Errorf("actisense: device acknowledged CAN Config %+v; requested %+v", accepted, config)
	}
	return accepted, nil
}

func (c *CommandSet) GetCANInfoField(ctx context.Context, field CANInfoField) (string, error) {
	command, err := field.command()
	if err != nil {
		return "", err
	}
	response, err := c.request(ctx, command, nil)
	if err != nil {
		return "", err
	}
	return DecodeCANInfoField(response, field)
}

// SetCANInfoField writes installation description 1 or 2. Manufacturer
// information is read-only. Firmware may persist these fields immediately.
func (c *CommandSet) SetCANInfoField(ctx context.Context, field CANInfoField, text string) (string, error) {
	if field == CANInfoManufacturerInformation {
		return "", errors.New("actisense: manufacturer CAN information is read-only")
	}
	command, err := field.command()
	if err != nil {
		return "", err
	}
	data, err := EncodeCANInfoField(text)
	if err != nil {
		return "", err
	}
	response, err := c.request(ctx, command, data)
	if err != nil {
		return "", err
	}
	accepted, err := DecodeCANInfoField(response, field)
	if err != nil {
		return "", err
	}
	if accepted != text {
		return accepted, fmt.Errorf("actisense: device acknowledged CAN Info text %q; requested %q", accepted, text)
	}
	return accepted, nil
}

func (c *CommandSet) GetRxPGN(ctx context.Context, pgn uint32) (RxPGNState, error) {
	if err := validatePGN(pgn); err != nil {
		return RxPGNState{}, err
	}
	response, err := c.request(ctx, BEMRxPGNEnable, RxPGNEnableGet(pgn))
	if err != nil {
		return RxPGNState{}, err
	}
	return DecodeRxPGNState(response)
}

// SetRxPGN accepts one of the four documented PGN masks or a default/no-change
// sentinel. A nil mask selects the device default. Source filtering is not supported.
func (c *CommandSet) SetRxPGN(ctx context.Context, pgn uint32, flag PGNEnableFlag, mask *uint32) (RxPGNState, error) {
	if err := validatePGNAndFlag(pgn, flag); err != nil {
		return RxPGNState{}, err
	}
	if mask != nil && !validRxPGNMask(*mask) && *mask != RxPGNMaskDefault && *mask != RxPGNMaskNoChange {
		return RxPGNState{}, fmt.Errorf("actisense: unsupported Rx PGN mask 0x%08X", *mask)
	}
	response, err := c.request(ctx, BEMRxPGNEnable, RxPGNEnableSet(pgn, flag, mask))
	if err != nil {
		return RxPGNState{}, err
	}
	state, err := DecodeRxPGNState(response)
	if err != nil {
		return RxPGNState{}, err
	}
	if state.PGN != pgn || state.Flag != flag || !validRxPGNMask(state.Mask) || (mask != nil && validRxPGNMask(*mask) && state.Mask != *mask) {
		return state, fmt.Errorf("actisense: device acknowledged Rx PGN state %+v; request was PGN %d flag %d", state, pgn, flag)
	}
	return state, nil
}

func (c *CommandSet) GetTxPGN(ctx context.Context, pgn uint32) (TxPGNState, error) {
	if err := validatePGN(pgn); err != nil {
		return TxPGNState{}, err
	}
	response, err := c.request(ctx, BEMTxPGNEnable, TxPGNEnableGet(pgn))
	if err != nil {
		return TxPGNState{}, err
	}
	return DecodeTxPGNState(response)
}

// SetTxPGN uses milliseconds for rates 1-65534 and zero for event-driven data.
// A nil rate or any value >= 65535 leaves the current rate unchanged. This
// command has no restore-default rate sentinel.
func (c *CommandSet) SetTxPGN(ctx context.Context, pgn uint32, flag PGNEnableFlag, rate *uint32) (TxPGNState, error) {
	if err := validatePGNAndFlag(pgn, flag); err != nil {
		return TxPGNState{}, err
	}
	response, err := c.request(ctx, BEMTxPGNEnable, TxPGNEnableSetFull(pgn, flag, rate))
	if err != nil {
		return TxPGNState{}, err
	}
	state, err := DecodeTxPGNState(response)
	if err != nil {
		return TxPGNState{}, err
	}
	if state.PGN != pgn || state.Enabled != uint8(flag) || (rate != nil && *rate < TxPGNRateNoChange && state.Rate != *rate) {
		return state, fmt.Errorf("actisense: device acknowledged Tx PGN state %+v; request was PGN %d flag %d", state, pgn, flag)
	}
	return state, nil
}

func validatePGN(pgn uint32) error {
	if pgn > 0x3FFFF || pgn&0x20000 != 0 {
		return fmt.Errorf("actisense: PGN %d is outside the NMEA 2000 range", pgn)
	}
	return nil
}

func validatePGNAndFlag(pgn uint32, flag PGNEnableFlag) error {
	if err := validatePGN(pgn); err != nil {
		return err
	}
	if flag > PGNRespondMode {
		return fmt.Errorf("actisense: PGN enable flag %d is invalid", flag)
	}
	return nil
}

// GetSupportedPGNs performs the command's caller-driven chunk walk. Each
// sub-list is a distinct request and all transfer invariants are checked.
func (c *CommandSet) GetSupportedPGNs(ctx context.Context) (SupportedPGNList, error) {
	accumulator := &supportedPGNAccumulator{}
	var index, transferID uint8
	for requests := 0; requests <= maxPGNListEntries; requests++ {
		response, err := c.request(ctx, BEMSupportedPGNList, []byte{index, transferID})
		if err != nil {
			return accumulator.result, err
		}
		part, err := decodeSupportedPGNPart(response)
		if err != nil {
			return accumulator.result, err
		}
		done, err := accumulator.feed(part)
		if err != nil {
			return accumulator.result, err
		}
		if done {
			return accumulator.result, nil
		}
		if len(part.entries) == 0 {
			return accumulator.result, errors.New("actisense: Supported PGN List made no progress")
		}
		index = part.first + uint8(len(part.entries))
		transferID = part.transferID
	}
	return accumulator.result, errors.New("actisense: Supported PGN List walk exceeded 256 requests")
}

func (c *CommandSet) DeletePGNLists(ctx context.Context, selector PGNListSelector) error {
	if err := selector.validate(); err != nil {
		return err
	}
	_, err := c.request(ctx, BEMDeletePGNLists, []byte{byte(selector)})
	return err
}

func (c *CommandSet) ActivatePGNLists(ctx context.Context) error {
	_, err := c.request(ctx, BEMActivatePGNLists, nil)
	return err
}

func (c *CommandSet) DefaultPGNLists(ctx context.Context, selector PGNListSelector) error {
	if err := selector.validate(); err != nil {
		return err
	}
	_, err := c.request(ctx, BEMDefaultPGNLists, []byte{byte(selector)})
	return err
}

func (c *CommandSet) GetPGNListParameters(ctx context.Context) (PGNListParameters, error) {
	response, err := c.request(ctx, BEMPGNListParameters, nil)
	if err != nil {
		return PGNListParameters{}, err
	}
	return DecodePGNListParameters(response)
}

func (c *CommandSet) GetRxPGNEnableList(ctx context.Context) (RxPGNEnableList, error) {
	accumulator := &rxF2Accumulator{}
	_, err := c.requestMulti(ctx, BEMRxPGNEnableListF2, nil, func(response BEMResponse) (bool, error) {
		accumulator.expectProprietary = CapabilitiesForModel(ModelID(response.ModelID)).ProprietaryPGNEnableListF2
		return accumulator.feed(response)
	})
	return accumulator.result, err
}

func (c *CommandSet) GetTxPGNEnableList(ctx context.Context) (TxPGNEnableList, error) {
	accumulator := &txF2Accumulator{}
	_, err := c.requestMulti(ctx, BEMTxPGNEnableListF2, nil, func(response BEMResponse) (bool, error) {
		accumulator.expectProprietary = CapabilitiesForModel(ModelID(response.ModelID)).ProprietaryPGNEnableListF2
		return accumulator.feed(response)
	})
	return accumulator.result, err
}

func (c *CommandSet) Echo(ctx context.Context, payload []byte) ([]byte, error) {
	if c.remote && len(payload) > maxRemoteEchoPayload {
		return nil, fmt.Errorf("actisense: remote Echo payload is %d bytes; PGN 126720 permits at most %d", len(payload), maxRemoteEchoPayload)
	}
	data, err := EncodeEcho(payload)
	if err != nil {
		return nil, err
	}
	response, err := c.request(ctx, BEMEcho, data)
	if err != nil {
		return nil, err
	}
	echoed, err := DecodeEcho(response)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(echoed, payload) {
		return echoed, errors.New("actisense: Echo response does not match the request")
	}
	return echoed, nil
}

func (c *CommandSet) GetTotalTime(ctx context.Context) (uint32, error) {
	response, err := c.request(ctx, BEMTotalTime, nil)
	if err != nil {
		return 0, err
	}
	return DecodeTotalTime(response)
}

func (c *CommandSet) SetTotalTime(ctx context.Context, seconds, passkey uint32) (uint32, error) {
	response, err := c.request(ctx, BEMTotalTime, TotalTimeSet(seconds, passkey))
	if err != nil {
		return 0, err
	}
	accepted, err := DecodeTotalTime(response)
	if err != nil {
		return 0, err
	}
	if accepted != seconds {
		return accepted, fmt.Errorf("actisense: device acknowledged total time %d; requested %d", accepted, seconds)
	}
	return accepted, nil
}

// Reinitialize explicitly asks the device to reboot. It is never invoked by
// connection setup, restore, or Close.
func (c *CommandSet) Reinitialize(ctx context.Context) error {
	_, err := c.request(ctx, BEMReinitialize, nil)
	return err
}

// CommitEEPROM explicitly persists staged session settings to EEPROM. It is
// never invoked automatically.
func (c *CommandSet) CommitEEPROM(ctx context.Context) error {
	_, err := c.request(ctx, BEMCommitEEPROM, nil)
	return err
}

// CommitFlash explicitly persists staged session settings to flash. It is
// never invoked automatically.
func (c *CommandSet) CommitFlash(ctx context.Context) error {
	_, err := c.request(ctx, BEMCommitFlash, nil)
	return err
}
