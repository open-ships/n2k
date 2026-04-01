# n2k

[![Tests](https://github.com/open-ships/n2k/actions/workflows/test.yaml/badge.svg)](https://github.com/open-ships/n2k/actions/workflows/test.yaml)
[![Go version](https://img.shields.io/github/go-mod/go-version/open-ships/n2k)](go.mod)

A Go library for decoding [NMEA 2000](https://www.nmea.org/content/STANDARDS/NMEA_2000) marine network messages into strongly-typed data structures. Raw CAN bus frames go in, Go structs come out.

## Installation

```bash
go get github.com/open-ships/n2k
```

## How It Works

The library implements a three-stage decoding pipeline:

```
CAN frames  -->  Endpoint  -->  Adapter  -->  Decoder  -->  Go structs
(hardware)      (transport)    (assembly)    (typing)      (your code)
```

**Endpoint** manages the connection to the NMEA 2000 gateway. Two implementations are included:
- `socketcanendpoint` -- Linux SocketCAN (kernel CAN drivers)
- `usbcanendpoint` -- USB-CAN Analyzer dongles (serial)

**Adapter** converts raw CAN frames into intermediate packets, handling CAN ID extraction, fast-packet assembly for multi-frame messages, and PGN identification.

**Decoder** translates packets into typed Go structs (e.g. `VesselHeading`, `WindData`, `EngineParametersRapidUpdate`) using the PGN definitions generated from the canboat database.

## Usage

```go
// Create the pipeline stages
endpoint := socketcanendpoint.NewSocketCANEndpoint(logger, "can0")
adapter := canadapter.NewCANAdapter()
decoder := pkt.NewPacketStruct()

// Wire them together
adapter.SetOutput(decoder)
decoder.SetOutput(yourHandler)
endpoint.SetOutput(adapter)

// Start receiving
err := endpoint.Run(ctx)
```

## Sniffer

A built-in CLI tool that reads from a SocketCAN interface and prints decoded messages as JSON to stdout, one per line.

```bash
go run ./cmd/sniffer.go -i can0
```

Pipe through `jq` for pretty-printing or filtering:

```bash
go run ./cmd/sniffer.go -i can0 | jq 'select(.Info.PGN == 127250)'
```

## Packages

| Package | Purpose |
|---------|---------|
| `pkg/endpoint` | Hardware abstraction for CAN data sources |
| `pkg/canbus` | Low-level CAN bus I/O (SocketCAN, USB-CAN) |
| `pkg/adapter` | CAN frame to NMEA 2000 packet conversion |
| `pkg/pkt` | Packet to typed Go struct decoding |
| `pkg/pgn` | PGN registry, decoders, and generated type definitions |
| `pkg/units` | Type-safe physical unit conversions (distance, velocity, temperature, etc.) |

## Development

Requires [just](https://github.com/casey/just) for running build commands.

```bash
just test          # run tests
just test-race     # run tests with race detector
just test-cover    # generate coverage report
just lint          # run linter
just fmt           # format code
```

## License

Apache 2.0 -- see [LICENSE](./LICENSE).

## Acknowledgments

This project is a fork of [boatkit-io/n2k](https://github.com/boatkit-io/n2k), which built the original Go implementation of this NMEA 2000 decoding pipeline.

The PGN definitions and decoders at the core of this library are generated from the [canboat](https://github.com/canboat/canboat) project's open-source NMEA 2000 database. canboat reverse-engineered the NMEA 2000 protocol through network observation and public sources, producing the comprehensive PGN catalog that makes libraries like this one possible. For deeper understanding of NMEA 2000 message semantics, field definitions, and manufacturer-specific PGNs, refer to the canboat documentation.
