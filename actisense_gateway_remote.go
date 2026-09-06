package n2k

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"

	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/raw"
)

// ActisenseRemoteDevice addresses a remote device through the gateway's own
// identity using BST-94 / PGN 126720. Each request first verifies its return
// address with a random Echo challenge, then binds correlation to that address
// and connection. It never claims a virtual node or changes PGN enable lists.
func (s *ActisenseGatewaySession) ActisenseRemoteDevice(source uint8, options ...ActisenseRemoteOption) (*ActisenseRemoteDevice, error) {
	if s == nil {
		return nil, actisense.ErrSessionClosed
	}
	s.mu.Lock()
	closed, terminalErr := s.closed, s.terminalErr
	s.mu.Unlock()
	if terminalErr != nil {
		return nil, terminalErr
	}
	if closed {
		return nil, actisense.ErrSessionClosed
	}
	return s.remote.device(source, options)
}

type actisenseGatewayProbe struct {
	source    uint8
	epoch     uint64
	challenge [16]byte
	result    chan actisenseGatewayProbeResult
	cancel    context.CancelCauseFunc
}

type actisenseGatewayProbeResult struct {
	address uint8
	err     error
}

func (s *ActisenseGatewaySession) acquireRemoteProbe(ctx context.Context) error {
	select {
	case s.probeGate <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-s.ctx.Done():
		return actisense.ErrSessionClosed
	}
}

// probeGate remains held through pending-command registration. A live Echo
// command to this source excludes another probe because device errors contain
// no challenge and would otherwise have ambiguous ownership.
func (s *ActisenseGatewaySession) verifyRemoteAddress(ctx context.Context, source uint8) error {
	s.remote.mu.Lock()
	for key := range s.remote.pending {
		if key.source == source && key.bemID == actisense.BEMEcho {
			s.remote.mu.Unlock()
			return fmt.Errorf("%w: remote address probe conflicts with source %d Echo", actisense.ErrRequestInFlight, source)
		}
	}
	s.remote.mu.Unlock()
	probeCtx, cancel := context.WithCancelCause(ctx)
	defer cancel(context.Canceled)
	probe := &actisenseGatewayProbe{source: source, result: make(chan actisenseGatewayProbeResult, 1), cancel: cancel}
	_, _ = rand.Read(probe.challenge[:])
	s.mu.Lock()
	if s.closed || !s.connected {
		s.mu.Unlock()
		return ErrActisenseNotReady
	}
	probe.epoch = s.epoch
	s.remoteProbe = probe
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		if s.remoteProbe == probe {
			s.remoteProbe = nil
		}
		s.mu.Unlock()
	}()
	echo, err := actisense.EncodeEcho(probe.challenge[:])
	if err != nil {
		return err
	}
	inner, err := encodeActisenseRemoteCommand(actisense.BEMEcho, echo)
	if err != nil {
		return err
	}
	if err := s.sendRemote(probeCtx, probe.epoch, source, append([]byte{0x11, 0x99}, inner...)); err != nil {
		if cause := context.Cause(probeCtx); cause != nil {
			return cause
		}
		return err
	}
	select {
	case result := <-probe.result:
		if result.err != nil {
			return result.err
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		if !s.connected || s.closed || s.epoch != probe.epoch || s.remoteProbe != probe {
			return ErrActisenseRemoteEpochChanged
		}
		if !s.remoteAddressKnown || s.remoteAddress != result.address {
			s.remote.invalidate(ErrActisenseRemoteEpochChanged)
			s.identityEpoch++
		}
		s.remoteAddress, s.remoteAddressKnown = result.address, true
		return nil
	case <-probeCtx.Done():
		return fmt.Errorf("n2k: verifying Actisense gateway return address: %w", context.Cause(probeCtx))
	case <-s.ctx.Done():
		return actisense.ErrSessionClosed
	}
}

func (s *ActisenseGatewaySession) sendRemote(ctx context.Context, epoch uint64, destination uint8, payload []byte) error {
	writeCtx, cancel := s.writeContext(ctx)
	defer cancel()
	return s.transport.WriteMessageEpoch(writeCtx, epoch, actisenseRemotePGN, actisenseRemotePriority, destination, payload)
}

func (s *ActisenseGatewaySession) invalidateRemoteLocked(err error) {
	if err == nil {
		err = ErrActisenseRemoteEpochChanged
	}
	s.remoteAddressKnown = false
	s.identityEpoch++
	if probe := s.remoteProbe; probe != nil {
		probe.cancel(err)
		select {
		case probe.result <- actisenseGatewayProbeResult{err: err}:
		default:
		}
		s.remoteProbe = nil
	}
	if s.remote != nil {
		s.remote.invalidate(err)
	}
}

func (s *ActisenseGatewaySession) handleGatewayMessage(message actisense.Message) {
	if message.Direction != actisense.DirectionReceived || !message.HasSource || message.PGN != actisenseRemotePGN || len(message.Data) < 5 || len(message.Data) > 223 || message.Data[0] != 0x11 || message.Data[1] != 0x99 {
		return
	}
	if s.handleRemoteProbe(message) {
		return
	}
	destination := message.Destination
	s.remote.observe(raw.Observation{
		Kind: raw.KindMessage, PGN: message.PGN, Priority: message.Priority,
		Source: message.Source, Destination: &destination, Payload: message.Data,
		Direction: raw.DirectionReceived,
	})
}

func (s *ActisenseGatewaySession) handleRemoteProbe(message actisense.Message) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	probe := s.remoteProbe
	if probe == nil || probe.epoch != s.epoch || message.Source != probe.source || message.Destination > 251 {
		return false
	}
	datagram, err := actisense.DecodeRaw(message.Data[2:])
	if err != nil {
		return false
	}
	response, ok, err := actisense.DecodeBEMResponse(datagram)
	if !ok || err != nil || response.BSTID != actisense.BSTBEMResponse {
		return false
	}
	result := actisenseGatewayProbeResult{address: message.Destination}
	switch response.BEMID {
	case actisense.BEMEcho:
		if response.ErrorCode != 0 {
			result.err = &actisense.DeviceError{Command: response.BEMID, Code: response.ErrorCode}
		} else {
			echo, decodeErr := actisense.DecodeEcho(response)
			if decodeErr != nil || !bytes.Equal(echo, probe.challenge[:]) {
				return false
			}
		}
	case actisense.BEMNegativeAck:
		diagnostic, ok, decodeErr := actisense.DecodeDiagnostic(response)
		if !ok || decodeErr != nil || diagnostic.NegativeAck == nil || byte(diagnostic.NegativeAck.UniqueCommandID) != actisense.BEMEcho {
			return false
		}
		nack := diagnostic.NegativeAck
		result.err = &actisense.NegativeAckError{Command: actisense.BEMEcho, UniqueCommandID: nack.UniqueCommandID, DeviceCode: nack.ErrorCode}
	default:
		return false
	}
	select {
	case probe.result <- result:
	default:
	}
	return true
}
