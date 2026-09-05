package n2k

import (
	"context"
	"fmt"
	"testing"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/adapter"
	"github.com/open-ships/n2k/internal/decoder"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/require"
)

func TestReplayObservationsIndependentEpochs(t *testing.T) {
	for _, secondEpoch := range []uint64{1, 0} {
		t.Run(fmt.Sprint(secondEpoch), func(t *testing.T) {
			first := Observation{
				Kind: ObservationMessage, PGN: 127250, Source: 0x42,
				NetworkID: "network-a", AdapterID: "adapter-a", ConnectionEpoch: 10, ClaimEpoch: 10,
				Payload: []byte{0, 0x5c, 0x3d, 0xff, 0x7f, 0xff, 0x7f, 0xfc},
			}
			second := first.Clone()
			second.NetworkID, second.AdapterID = "network-b", "adapter-b"
			second.ConnectionEpoch, second.ClaimEpoch = secondEpoch, secondEpoch
			scanner := NewScanner(context.Background(), ReplayObservations([]Observation{first, second}))
			defer func() { _ = scanner.Close() }()
			var delivered []string
			for scanner.Next() {
				message := scanner.Message().(*pgn.VesselHeading)
				delivered = append(delivered, message.Info.NetworkID)
			}
			require.NoError(t, scanner.Err())
			t.Logf("replayed independent epochs (10,%d), delivered networks %v, scanner error %v", secondEpoch, delivered, scanner.Err())
			require.Equal(t, []string{"network-a", "network-b"}, delivered)
		})
	}
}

func TestCANAdapterFastPacketsCannotSpanEpochs(t *testing.T) {
	for _, changed := range []struct {
		name              string
		connection, claim uint64
	}{{"connection", 2, 1}, {"claim", 1, 2}} {
		t.Run(changed.name, func(t *testing.T) {
			collector := &epochPacketCollector{}
			assembler := adapter.NewCANAdapter()
			assembler.SetOutput(collector)
			id := framer.BuildCANID(126998, 6, 0x42, 255)
			info := pgn.MessageInfo{PGN: 126998, SourceId: 0x42, NetworkID: "network-a", AdapterID: "adapter-a", ConnectionEpoch: 1, ClaimEpoch: 1}
			first := can.Frame{ID: id, Length: 8, Data: [8]byte{0x20, 9, 3, 1, 'a', 3, 1, 'b'}}
			assembler.HandleMessageWithInfo(&first, info)
			info.ConnectionEpoch, info.ClaimEpoch = changed.connection, changed.claim
			last := can.Frame{ID: id, Length: 4, Data: [8]byte{0x21, 3, 1, 'z'}}
			assembler.HandleMessageWithInfo(&last, info)
			// A continuation from a new epoch must not consume the old prefix.
			require.Empty(t, collector.packets)
		})
	}
}

func TestReplayObservationsInterleavedFastPacketEpochs(t *testing.T) {
	for _, sameNetwork := range []bool{false, true} {
		t.Run(fmt.Sprint(sameNetwork), func(t *testing.T) {
			id := framer.BuildCANID(126998, 6, 0x42, 255)
			first := framer.FrameFastPacket(id, []byte{3, 1, 'a', 3, 1, 'b', 3, 1, 'c'}, 0)
			second := framer.FrameFastPacket(id, []byte{3, 1, 'x', 3, 1, 'y', 3, 1, 'z'}, 0)
			require.Len(t, first, 2)
			require.Len(t, second, 2)
			secondNetwork := "network-b"
			if sameNetwork {
				secondNetwork = "network-a"
			}
			observation := func(frame can.Frame, network string, epoch uint64) Observation {
				return Observation{Kind: ObservationFrame, Frame: &frame, NetworkID: network, AdapterID: "replay", ConnectionEpoch: epoch, ClaimEpoch: epoch}
			}
			observations := []Observation{
				observation(first[0], "network-a", 10),
				observation(second[0], secondNetwork, 1),
				observation(first[1], "network-a", 10),
				observation(second[1], secondNetwork, 1),
			}
			scanner := NewScanner(context.Background(), ReplayObservations(observations))
			defer func() { _ = scanner.Close() }()
			var decoded []string
			for scanner.Next() {
				message := scanner.Message().(*pgn.ConfigurationInformation)
				decoded = append(decoded, message.InstallationDescription1+message.InstallationDescription2+message.ManufacturerInformation)
			}
			require.NoError(t, scanner.Err())
			require.Equal(t, []string{"abc", "xyz"}, decoded)
		})
	}
}

type epochPacketCollector struct{ packets []decoder.Packet }

func (c *epochPacketCollector) Decode(packet decoder.Packet) {
	c.packets = append(c.packets, packet)
}
