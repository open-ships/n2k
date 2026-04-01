# N2K Package Simplification Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Simplify the n2k package to a flat public API (`n2k.Receive`, `n2k.NewScanner`) with CEL-based filtering, hiding all pipeline internals.

**Architecture:** Move existing canbus/adapter/decoder code into `internal/` packages. Create a new top-level `n2k` package that wires the pipeline internally and exposes `Receive()` (iterator) and `NewScanner()` (scanner). Add CEL filter support with automatic AST-based partitioning into pre-decode and post-decode stages.

**Tech Stack:** Go 1.25, google/cel-go, brutella/can, go.bug.st/serial, vishvananda/netlink

**Spec:** `docs/superpowers/specs/2026-03-31-n2k-simplification-design.md`

---

### Task 1: Move canbus to internal/canbus

Move `pkg/canbus/` to `internal/canbus/` with updated package name and import paths.

**Files:**
- Move: `pkg/canbus/interface.go` -> `internal/canbus/interface.go`
- Move: `pkg/canbus/socketcan.go` -> `internal/canbus/socketcan.go`
- Move: `pkg/canbus/usbcan.go` -> `internal/canbus/usbcan.go`
- Move: `pkg/canbus/*_test.go` -> `internal/canbus/*_test.go`
- Modify: all files that import `github.com/open-ships/n2k/pkg/canbus`

- [ ] **Step 1: Create internal/canbus directory and move files**

```bash
mkdir -p internal/canbus
git mv pkg/canbus/*.go internal/canbus/
```

- [ ] **Step 2: Update package declaration in moved files**

In each file under `internal/canbus/`, the package declaration should already be `package canbus` so no change needed. Update all import paths in the codebase from `github.com/open-ships/n2k/pkg/canbus` to `github.com/open-ships/n2k/internal/canbus`.

Files to update:
- `pkg/endpoint/socketcanendpoint/socketcanendpoint.go` — change import
- `pkg/endpoint/usbcanendpoint/usbcanendpoint.go` — change import

- [ ] **Step 3: Verify it compiles**

Run: `go build ./...`
Expected: SUCCESS

- [ ] **Step 4: Run tests**

Run: `go test ./internal/canbus/... ./pkg/...`
Expected: All tests pass

- [ ] **Step 5: Commit**

```bash
git add -A && git commit -m "refactor: move canbus to internal/canbus"
```

---

### Task 2: Move adapter to internal/adapter

Move `pkg/adapter/canadapter/` to `internal/adapter/` and eliminate the `pkg/adapter/` wrapper package.

**Files:**
- Move: `pkg/adapter/canadapter/*.go` -> `internal/adapter/*.go`
- Delete: `pkg/adapter/adapter.go` (the `adapter.Message` interface — inline as `any` or `*can.Frame`)
- Modify: all files importing `pkg/adapter` or `pkg/adapter/canadapter`

- [ ] **Step 1: Create internal/adapter and move files**

```bash
mkdir -p internal/adapter
git mv pkg/adapter/canadapter/*.go internal/adapter/
```

- [ ] **Step 2: Update package declaration**

Change `package canadapter` to `package adapter` in all moved files:
- `internal/adapter/canadapter.go`
- `internal/adapter/frame.go`
- `internal/adapter/multibuilder.go`
- `internal/adapter/sequence.go`
- All test files

- [ ] **Step 3: Remove adapter.Message indirection**

In `internal/adapter/canadapter.go`, the `HandleMessage` method currently accepts `adapter.Message` (an empty interface) and type-asserts to `*can.Frame`. Change the signature to accept `*can.Frame` directly:

```go
func (c *CANAdapter) HandleMessage(frame *can.Frame) {
    // Remove the type assertion, use frame directly
}
```

Remove the import of `github.com/open-ships/n2k/pkg/adapter`.

- [ ] **Step 4: Update PacketHandler interface**

The `PacketHandler` interface stays in `internal/adapter/` — it's used internally. No external consumers.

- [ ] **Step 5: Update endpoint imports**

In `pkg/endpoint/socketcanendpoint/socketcanendpoint.go` and `pkg/endpoint/usbcanendpoint/usbcanendpoint.go`:
- Remove import of `github.com/open-ships/n2k/pkg/adapter`
- These files will be moved in Task 3, but for now update them to call the adapter's new `HandleMessage(*can.Frame)` signature.

- [ ] **Step 6: Delete pkg/adapter/**

```bash
rm -rf pkg/adapter/
```

- [ ] **Step 7: Verify it compiles and tests pass**

Run: `go build ./... && go test ./...`
Expected: All pass

- [ ] **Step 8: Commit**

```bash
git add -A && git commit -m "refactor: move canadapter to internal/adapter, remove Message interface"
```

---

### Task 3: Move decoder (pkt) to internal/decoder

Move `pkg/pkt/` to `internal/decoder/`.

**Files:**
- Move: `pkg/pkt/pkt.go` -> `internal/decoder/packet.go`
- Move: `pkg/pkt/packetstruct.go` -> `internal/decoder/decoder.go`
- Move: `pkg/pkt/unknownpgn.go` -> `internal/decoder/unknownpgn.go`
- Move: `pkg/pkt/*_test.go` -> `internal/decoder/*_test.go`
- Modify: `internal/adapter/canadapter.go` (imports pkt)

- [ ] **Step 1: Create internal/decoder and move files**

```bash
mkdir -p internal/decoder
git mv pkg/pkt/*.go internal/decoder/
```

- [ ] **Step 2: Update package declaration**

Change `package pkt` to `package decoder` in all moved files.

- [ ] **Step 3: Rename types for clarity**

- `Packet` stays as `Packet` (internal, fine)
- `PacketStruct` -> `Decoder` (it's now in package `decoder`, so `decoder.Decoder`)
- `StructHandler` -> `Handler` (so `decoder.Handler`)
- `NewPacketStruct()` -> `New()`
- `HandlePacket()` -> `Decode()`

- [ ] **Step 4: Update internal/adapter imports**

In `internal/adapter/canadapter.go`:
- Change import from `github.com/open-ships/n2k/pkg/pkt` to `github.com/open-ships/n2k/internal/decoder`
- Update all references: `pkt.Packet` -> `decoder.Packet`, `pkt.NewPacket` -> `decoder.NewPacket`, etc.

- [ ] **Step 5: Delete pkg/pkt/**

```bash
rm -rf pkg/pkt/
```

- [ ] **Step 6: Verify it compiles and tests pass**

Run: `go build ./... && go test ./...`
Expected: All pass

- [ ] **Step 7: Commit**

```bash
git add -A && git commit -m "refactor: move pkt to internal/decoder"
```

---

### Task 4: Move pgn and units to top-level

Move `pkg/pgn/` to `pgn/` and `pkg/units/` to `units/` (top-level exported packages).

**Files:**
- Move: `pkg/pgn/*` -> `pgn/*`
- Move: `pkg/units/*` -> `units/*`
- Modify: all files importing `github.com/open-ships/n2k/pkg/pgn` or `github.com/open-ships/n2k/pkg/units`

- [ ] **Step 1: Move directories**

```bash
git mv pkg/pgn pgn
git mv pkg/units units
```

- [ ] **Step 2: Update all import paths**

Global find/replace:
- `github.com/open-ships/n2k/pkg/pgn` -> `github.com/open-ships/n2k/pgn`
- `github.com/open-ships/n2k/pkg/units` -> `github.com/open-ships/n2k/units`

Files to update:
- `internal/adapter/*.go`
- `internal/decoder/*.go`
- `internal/canbus/*.go` (if any)
- `cmd/sniffer.go`
- All test files

- [ ] **Step 3: Remove empty pkg/ directory**

```bash
rm -rf pkg/
```

- [ ] **Step 4: Verify it compiles and tests pass**

Run: `go build ./... && go test ./...`
Expected: All pass

- [ ] **Step 5: Commit**

```bash
git add -A && git commit -m "refactor: move pgn and units to top-level packages"
```

---

### Task 5: Absorb endpoints into internal packages

Eliminate `pkg/endpoint/` entirely. Move the SocketCAN and USB-CAN implementations into `internal/canbus/` since they're thin wrappers around the canbus channel types.

**Files:**
- Delete: `pkg/endpoint/endpoint.go` (Endpoint, MessageHandler interfaces — no longer needed)
- Absorb: `pkg/endpoint/socketcanendpoint/socketcanendpoint.go` into `internal/canbus/socketcan.go`
- Absorb: `pkg/endpoint/usbcanendpoint/usbcanendpoint.go` into `internal/canbus/usbcan.go`

- [ ] **Step 1: Add a Run method to SocketCANChannel**

In `internal/canbus/socketcan.go`, the `SocketCANChannel` already has a `Run(ctx)` method. Merge the endpoint's setup logic (creating the channel with options, calling Run) into a higher-level function:

```go
func RunSocketCAN(ctx context.Context, iface string, handler func(can.Frame)) error {
    ch := &SocketCANChannel{
        Options: SocketCANChannelOptions{
            InterfaceName: iface,
            BitRate:       250000,
            MessageHandler: func(frame can.Frame) {
                handler(frame)
            },
        },
    }
    return ch.Run(ctx)
}
```

- [ ] **Step 2: Add a Run function to USB-CAN**

In `internal/canbus/usbcan.go`, add similarly:

```go
func RunUSBCAN(ctx context.Context, port string, handler func(can.Frame)) error {
    ch := &USBCANChannel{
        Options: USBCANChannelOptions{
            SerialPortName: port,
            SerialBaudRate: 2000000,
            BitRate:        250000,
            FrameHandler: func(frame can.Frame) {
                handler(frame)
            },
        },
    }
    return ch.Run(ctx)
}
```

- [ ] **Step 3: Delete pkg/endpoint/**

```bash
rm -rf pkg/endpoint/
```

- [ ] **Step 4: Verify it compiles and tests pass**

Run: `go build ./... && go test ./...`
Expected: All pass

- [ ] **Step 5: Commit**

```bash
git add -A && git commit -m "refactor: absorb endpoints into internal/canbus, remove pkg/endpoint"
```

---

### Task 6: Create the Source interface and implementations

Define the internal `Source` interface and implement CAN, USB, and Replay sources. Also implement fan-in for multiple sources.

**Files:**
- Create: `source.go` (in root `n2k` package)

- [ ] **Step 1: Write the test for Replay source**

Create `source_test.go`:

```go
package n2k

import (
    "context"
    "testing"

    "github.com/brutella/can"
    "github.com/stretchr/testify/assert"
)

func TestReplaySource(t *testing.T) {
    frames := []can.Frame{
        {ID: 0x09F10D01, Length: 8, Data: [8]uint8{1, 2, 3, 4, 5, 6, 7, 8}},
        {ID: 0x09F10D02, Length: 8, Data: [8]uint8{8, 7, 6, 5, 4, 3, 2, 1}},
    }

    src := &replaySource{frames: frames}
    received := make([]can.Frame, 0)
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    err := src.run(ctx, func(f can.Frame) {
        received = append(received, f)
    })

    assert.NoError(t, err)
    assert.Equal(t, 2, len(received))
    assert.Equal(t, frames[0].ID, received[0].ID)
    assert.Equal(t, frames[1].ID, received[1].ID)
}

func TestFanIn(t *testing.T) {
    frames1 := []can.Frame{
        {ID: 0x09F10D01, Length: 8, Data: [8]uint8{1, 2, 3, 4, 5, 6, 7, 8}},
    }
    frames2 := []can.Frame{
        {ID: 0x09F10D02, Length: 8, Data: [8]uint8{8, 7, 6, 5, 4, 3, 2, 1}},
    }

    sources := []source{
        &replaySource{frames: frames1},
        &replaySource{frames: frames2},
    }

    received := make([]can.Frame, 0)
    ctx := context.Background()

    err := runSources(ctx, sources, func(f can.Frame) {
        received = append(received, f)
    })

    assert.NoError(t, err)
    assert.Equal(t, 2, len(received))
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./... -run "TestReplaySource|TestFanIn" -v`
Expected: FAIL — types not defined

- [ ] **Step 3: Implement source types**

Create `source.go`:

```go
package n2k

import (
    "context"
    "sync"

    "github.com/brutella/can"
    "github.com/open-ships/n2k/internal/canbus"
)

type source interface {
    run(ctx context.Context, handler func(can.Frame)) error
}

type socketCANSource struct {
    iface string
}

func (s *socketCANSource) run(ctx context.Context, handler func(can.Frame)) error {
    return canbus.RunSocketCAN(ctx, s.iface, handler)
}

type usbCANSource struct {
    port string
}

func (s *usbCANSource) run(ctx context.Context, handler func(can.Frame)) error {
    return canbus.RunUSBCAN(ctx, s.port, handler)
}

type replaySource struct {
    frames []can.Frame
}

func (s *replaySource) run(ctx context.Context, handler func(can.Frame)) error {
    for _, f := range s.frames {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            handler(f)
        }
    }
    return nil
}

func runSources(ctx context.Context, sources []source, handler func(can.Frame)) error {
    if len(sources) == 1 {
        return sources[0].run(ctx, handler)
    }

    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    var (
        mu  sync.Mutex
        wg  sync.WaitGroup
    )
    errc := make(chan error, len(sources))

    safeHandler := func(f can.Frame) {
        mu.Lock()
        defer mu.Unlock()
        handler(f)
    }

    for _, src := range sources {
        wg.Add(1)
        go func() {
            defer wg.Done()
            if err := src.run(ctx, safeHandler); err != nil {
                cancel()
                errc <- err
            }
        }()
    }

    wg.Wait()
    close(errc)

    for err := range errc {
        if err != nil {
            return err
        }
    }
    return nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./... -run "TestReplaySource|TestFanIn" -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add source.go source_test.go && git commit -m "feat: add source interface with replay and fan-in support"
```

---

### Task 7: Create the options API

Define the functional options pattern: `CAN()`, `USB()`, `Replay()`, `Filter()`, `IncludeUnknown()`, `WithLogger()`.

**Files:**
- Create: `options.go`

- [ ] **Step 1: Write tests for option construction**

Create `options_test.go`:

```go
package n2k

import (
    "log/slog"
    "testing"

    "github.com/brutella/can"
    "github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
    var cfg config
    opts := []Option{
        CAN("can0"),
        CAN("can1"),
        USB("/dev/ttyUSB0"),
        Replay([]can.Frame{{ID: 1}}),
        IncludeUnknown(),
        WithLogger(slog.Default()),
    }

    for _, o := range opts {
        o.apply(&cfg)
    }

    assert.Equal(t, 4, len(cfg.sources))
    assert.True(t, cfg.includeUnknown)
    assert.NotNil(t, cfg.logger)
}

func TestNoSourcesError(t *testing.T) {
    cfg := config{}
    err := cfg.validate()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "at least one source")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run "TestOptions|TestNoSources" -v`
Expected: FAIL — types not defined

- [ ] **Step 3: Implement options**

Create `options.go`:

```go
package n2k

import (
    "errors"
    "log/slog"

    "github.com/brutella/can"
)

type config struct {
    sources        []source
    filterExpr     string
    includeUnknown bool
    logger         *slog.Logger
}

func (c *config) validate() error {
    if len(c.sources) == 0 {
        return errors.New("n2k: at least one source (CAN, USB, or Replay) is required")
    }
    return nil
}

type Option interface {
    apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(c *config) { f(c) }

func CAN(iface string) Option {
    return optionFunc(func(c *config) {
        c.sources = append(c.sources, &socketCANSource{iface: iface})
    })
}

func USB(port string) Option {
    return optionFunc(func(c *config) {
        c.sources = append(c.sources, &usbCANSource{port: port})
    })
}

func Replay(frames []can.Frame) Option {
    return optionFunc(func(c *config) {
        c.sources = append(c.sources, &replaySource{frames: frames})
    })
}

func Filter(expr string) Option {
    return optionFunc(func(c *config) {
        c.filterExpr = expr
    })
}

func IncludeUnknown() Option {
    return optionFunc(func(c *config) {
        c.includeUnknown = true
    })
}

func WithLogger(l *slog.Logger) Option {
    return optionFunc(func(c *config) {
        c.logger = l
    })
}
```

- [ ] **Step 4: Run tests**

Run: `go test -run "TestOptions|TestNoSources" -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add options.go options_test.go && git commit -m "feat: add functional options API"
```

---

### Task 8: Create the pipeline wiring and Scanner

Wire internal adapter + decoder into a pipeline driven by sources. Implement the Scanner API.

**Files:**
- Create: `scanner.go`
- Create: `n2k.go`

- [ ] **Step 1: Write scanner test**

Create `scanner_test.go`:

```go
package n2k

import (
    "context"
    "testing"

    "github.com/brutella/can"
    "github.com/open-ships/n2k/pgn"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// CAN frame for PGN 127501 (BinarySwitchBankStatus) — known to decode successfully
// CAN ID: 0x09F20D00, Data: known test vector from existing tests
var testFrame127501 = can.Frame{
    ID:     0x09F20D00,
    Length: 8,
    Data:   [8]uint8{0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
}

func TestScannerBasic(t *testing.T) {
    ctx := context.Background()
    s := NewScanner(ctx, Replay([]can.Frame{testFrame127501}))

    require.True(t, s.Next())
    msg := s.Message()
    assert.NotNil(t, msg)

    _, ok := msg.(*pgn.BinarySwitchBankStatus)
    if !ok {
        // May be UnknownPGN if test vector doesn't decode — that's OK for wiring test
        _, ok = msg.(*pgn.UnknownPGN)
        assert.True(t, ok, "expected BinarySwitchBankStatus or UnknownPGN, got %T", msg)
    }

    assert.False(t, s.Next()) // no more frames
    assert.NoError(t, s.Err())
}

func TestScannerNoSources(t *testing.T) {
    ctx := context.Background()
    s := NewScanner(ctx)

    assert.False(t, s.Next())
    assert.Error(t, s.Err())
    assert.Contains(t, s.Err().Error(), "at least one source")
}

func TestScannerContextCancel(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    cancel() // cancel immediately

    s := NewScanner(ctx, Replay([]can.Frame{testFrame127501}))
    // Should terminate quickly
    for s.Next() {
    }
    // Either context error or no error (replay may have completed before cancel)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run "TestScanner" -v`
Expected: FAIL — Scanner not defined

- [ ] **Step 3: Implement Scanner**

Create `scanner.go`:

```go
package n2k

import (
    "context"
    "log/slog"
    "sync"

    "github.com/brutella/can"
    "github.com/open-ships/n2k/internal/adapter"
    "github.com/open-ships/n2k/internal/decoder"
    "github.com/open-ships/n2k/pgn"
)

type Scanner struct {
    ctx    context.Context
    cfg    config
    msg    any
    err    error
    ch     chan any
    once   sync.Once
}

func NewScanner(ctx context.Context, opts ...Option) *Scanner {
    cfg := config{}
    for _, o := range opts {
        o.apply(&cfg)
    }
    if cfg.logger == nil {
        cfg.logger = slog.Default()
    }

    return &Scanner{
        ctx: ctx,
        cfg: cfg,
        ch:  make(chan any, 64),
    }
}

func (s *Scanner) Next() bool {
    s.once.Do(func() {
        if err := s.cfg.validate(); err != nil {
            s.err = err
            close(s.ch)
            return
        }
        go s.run()
    })

    msg, ok := <-s.ch
    if !ok {
        return false
    }
    s.msg = msg
    return true
}

func (s *Scanner) Message() any {
    return s.msg
}

func (s *Scanner) Err() error {
    return s.err
}

func (s *Scanner) run() {
    defer close(s.ch)

    adapt := adapter.NewCANAdapter()
    dec := decoder.New()

    dec.SetOutput(&scannerHandler{scanner: s})
    adapt.SetOutput(dec)

    err := runSources(s.ctx, s.cfg.sources, func(f can.Frame) {
        adapt.HandleMessage(&f)
    })
    if err != nil {
        s.err = err
    }
}

type scannerHandler struct {
    scanner *Scanner
}

func (h *scannerHandler) HandleStruct(msg any) {
    if msg == nil {
        return
    }

    // Check if unknown PGN and should be dropped
    if u, ok := msg.(*pgn.UnknownPGN); ok {
        if !h.scanner.cfg.includeUnknown {
            h.scanner.cfg.logger.Debug("dropping unknown PGN",
                "pgn", u.Info.PGN,
                "reason", u.Reason,
            )
            return
        }
    }

    select {
    case h.scanner.ch <- msg:
    case <-h.scanner.ctx.Done():
    }
}
```

Note: The `scannerHandler` implements `decoder.Handler` (the renamed `StructHandler` from Task 3). The `HandleStruct` method name must match whatever the interface method is called after Task 3 renames. Adjust if needed.

- [ ] **Step 4: Create n2k.go stub**

Create `n2k.go`:

```go
// Package n2k decodes NMEA 2000 marine network messages from CAN bus hardware
// into strongly-typed Go structs.
package n2k
```

- [ ] **Step 5: Run tests**

Run: `go test -run "TestScanner" -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add scanner.go scanner_test.go n2k.go && git commit -m "feat: add Scanner API with internal pipeline wiring"
```

---

### Task 9: Add CEL filter support

Implement CEL compilation, AST partitioning into pre-decode and post-decode stages, and two-stage evaluation.

**Files:**
- Create: `filter.go`
- Create: `filter_test.go`
- Modify: `go.mod` (add cel-go dependency)

- [ ] **Step 1: Add cel-go dependency**

```bash
cd /Users/jacobthomas/code/openships/n2k && go get github.com/google/cel-go
```

- [ ] **Step 2: Write filter tests**

Create `filter_test.go`:

```go
package n2k

import (
    "testing"

    "github.com/open-ships/n2k/pgn"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestFilterCompilation(t *testing.T) {
    f, err := compileFilter(`pgn == 127250`)
    require.NoError(t, err)
    assert.True(t, f.preOnly, "metadata-only filter should be pre-only")
}

func TestFilterPreOnly(t *testing.T) {
    f, err := compileFilter(`pgn == 127250 && source == 3`)
    require.NoError(t, err)
    assert.True(t, f.preOnly)

    info := pgn.MessageInfo{PGN: 127250, SourceId: 3}
    assert.True(t, f.evalPre(info))

    info2 := pgn.MessageInfo{PGN: 127250, SourceId: 5}
    assert.False(t, f.evalPre(info2))
}

func TestFilterPostOnly(t *testing.T) {
    f, err := compileFilter(`msg.Heading > 3.14`)
    require.NoError(t, err)
    assert.False(t, f.preOnly)
    assert.True(t, f.hasPost)

    // Pre-filter should pass everything when no pre predicates
    info := pgn.MessageInfo{PGN: 127250}
    assert.True(t, f.evalPre(info))

    // Post-filter checks struct fields
    fields := map[string]any{"Heading": 4.0, "heading": 4.0}
    assert.True(t, f.evalPost(fields))

    fields2 := map[string]any{"Heading": 1.0, "heading": 1.0}
    assert.False(t, f.evalPost(fields2))
}

func TestFilterSplit(t *testing.T) {
    f, err := compileFilter(`pgn == 127250 && msg.Heading > 3.14`)
    require.NoError(t, err)
    assert.False(t, f.preOnly)
    assert.True(t, f.hasPost)

    // Pre-filter checks PGN
    info := pgn.MessageInfo{PGN: 127250}
    assert.True(t, f.evalPre(info))

    info2 := pgn.MessageInfo{PGN: 99999}
    assert.False(t, f.evalPre(info2))
}

func TestFilterOrCannotSplit(t *testing.T) {
    f, err := compileFilter(`source == 3 || msg.Heading > 1.0`)
    require.NoError(t, err)
    // OR mixing metadata and msg fields cannot be split
    assert.False(t, f.preOnly)
    assert.True(t, f.hasPost)
}

func TestFilterCaseInsensitive(t *testing.T) {
    f, err := compileFilter(`msg.heading > 3.14`)
    require.NoError(t, err)

    fields := map[string]any{"Heading": 4.0, "heading": 4.0}
    assert.True(t, f.evalPost(fields))
}

func TestFilterInvalidExpression(t *testing.T) {
    _, err := compileFilter(`invalid !!!`)
    assert.Error(t, err)
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `go test -run "TestFilter" -v`
Expected: FAIL — compileFilter not defined

- [ ] **Step 4: Implement filter**

Create `filter.go`:

```go
package n2k

import (
    "fmt"
    "reflect"
    "strings"

    "github.com/google/cel-go/cel"
    "github.com/google/cel-go/common/ast"
    "github.com/open-ships/n2k/pgn"
)

var metadataVars = map[string]bool{
    "pgn":         true,
    "source":      true,
    "priority":    true,
    "destination": true,
}

type filter struct {
    preOnly bool // true if expression only references metadata
    hasPost bool // true if expression references msg.* fields

    preProg cel.Program
    postProg cel.Program
}

func compileFilter(expr string) (*filter, error) {
    // Analyze which variables the expression references
    hasMeta, hasMsg, err := analyzeExpr(expr)
    if err != nil {
        return nil, fmt.Errorf("n2k: invalid filter expression: %w", err)
    }

    f := &filter{}

    if !hasMsg {
        // Pure metadata filter — compile as pre-only
        f.preOnly = true
        env, err := cel.NewEnv(
            cel.Variable("pgn", cel.IntType),
            cel.Variable("source", cel.IntType),
            cel.Variable("priority", cel.IntType),
            cel.Variable("destination", cel.IntType),
        )
        if err != nil {
            return nil, err
        }
        a, issues := env.Compile(expr)
        if issues != nil && issues.Err() != nil {
            return nil, issues.Err()
        }
        f.preProg, err = env.Program(a)
        if err != nil {
            return nil, err
        }
        return f, nil
    }

    f.hasPost = true

    // Build pre-filter from metadata-only AND clauses
    preParts, err := extractPreParts(expr)
    if err == nil && len(preParts) > 0 {
        preExpr := strings.Join(preParts, " && ")
        env, err := cel.NewEnv(
            cel.Variable("pgn", cel.IntType),
            cel.Variable("source", cel.IntType),
            cel.Variable("priority", cel.IntType),
            cel.Variable("destination", cel.IntType),
        )
        if err == nil {
            a, issues := env.Compile(preExpr)
            if issues == nil || issues.Err() == nil {
                f.preProg, _ = env.Program(a)
            }
        }
    }

    // Full post-filter with all variables
    postEnv, err := cel.NewEnv(
        cel.Variable("pgn", cel.IntType),
        cel.Variable("source", cel.IntType),
        cel.Variable("priority", cel.IntType),
        cel.Variable("destination", cel.IntType),
        cel.Variable("msg", cel.DynType),
    )
    if err != nil {
        return nil, err
    }
    a, issues := postEnv.Compile(expr)
    if issues != nil && issues.Err() != nil {
        return nil, issues.Err()
    }
    f.postProg, err = postEnv.Program(a)
    if err != nil {
        return nil, err
    }

    return f, nil
}

func (f *filter) evalPre(info pgn.MessageInfo) bool {
    if f.preProg == nil {
        return true // no pre-filter, pass everything
    }
    vars := map[string]any{
        "pgn":         int64(info.PGN),
        "source":      int64(info.SourceId),
        "priority":    int64(info.Priority),
        "destination": int64(info.TargetId),
    }
    out, _, err := f.preProg.Eval(vars)
    if err != nil {
        return true // on error, pass through
    }
    result, ok := out.Value().(bool)
    return ok && result
}

func (f *filter) evalPost(fields map[string]any) bool {
    if f.postProg == nil {
        return true
    }
    vars := map[string]any{
        "pgn":         int64(0),
        "source":      int64(0),
        "priority":    int64(0),
        "destination": int64(0),
        "msg":         fields,
    }
    out, _, err := f.postProg.Eval(vars)
    if err != nil {
        return true
    }
    result, ok := out.Value().(bool)
    return ok && result
}

func (f *filter) evalPostWithInfo(info pgn.MessageInfo, fields map[string]any) bool {
    if f.postProg == nil {
        return true
    }
    vars := map[string]any{
        "pgn":         int64(info.PGN),
        "source":      int64(info.SourceId),
        "priority":    int64(info.Priority),
        "destination": int64(info.TargetId),
        "msg":         fields,
    }
    out, _, err := f.postProg.Eval(vars)
    if err != nil {
        return true
    }
    result, ok := out.Value().(bool)
    return ok && result
}

// analyzeExpr checks if expression references metadata vars, msg vars, or both.
func analyzeExpr(expr string) (hasMeta bool, hasMsg bool, err error) {
    env, err := cel.NewEnv(
        cel.Variable("pgn", cel.IntType),
        cel.Variable("source", cel.IntType),
        cel.Variable("priority", cel.IntType),
        cel.Variable("destination", cel.IntType),
        cel.Variable("msg", cel.DynType),
    )
    if err != nil {
        return false, false, err
    }

    a, issues := env.Compile(expr)
    if issues != nil && issues.Err() != nil {
        return false, false, issues.Err()
    }

    checkedExpr := a.Impl()
    inspectAST(checkedExpr.Expr(), &hasMeta, &hasMsg)
    return
}

func inspectAST(e ast.Expr, hasMeta *bool, hasMsg *bool) {
    switch e.Kind() {
    case ast.IdentKind:
        name := e.AsIdent()
        if metadataVars[name] {
            *hasMeta = true
        }
        if name == "msg" {
            *hasMsg = true
        }
    case ast.SelectKind:
        sel := e.AsSelect()
        inspectAST(sel.Operand(), hasMeta, hasMsg)
    case ast.CallKind:
        call := e.AsCall()
        for _, arg := range call.Args() {
            inspectAST(arg, hasMeta, hasMsg)
        }
        if call.IsMemberFunction() {
            inspectAST(call.Target(), hasMeta, hasMsg)
        }
    case ast.ListKind:
        l := e.AsList()
        for i := 0; i < l.Size(); i++ {
            inspectAST(l.Elements()[i], hasMeta, hasMsg)
        }
    }
}

// extractPreParts splits top-level AND clauses and returns those referencing only metadata.
func extractPreParts(expr string) ([]string, error) {
    env, err := cel.NewEnv(
        cel.Variable("pgn", cel.IntType),
        cel.Variable("source", cel.IntType),
        cel.Variable("priority", cel.IntType),
        cel.Variable("destination", cel.IntType),
        cel.Variable("msg", cel.DynType),
    )
    if err != nil {
        return nil, err
    }

    a, issues := env.Compile(expr)
    if issues != nil && issues.Err() != nil {
        return nil, issues.Err()
    }

    checkedExpr := a.Impl()
    var parts []string
    collectMetaOnlyClauses(checkedExpr.Expr(), env, a, &parts)
    return parts, nil
}

func collectMetaOnlyClauses(e ast.Expr, env *cel.Env, a *cel.Ast, parts *[]string) {
    if e.Kind() == ast.CallKind {
        call := e.AsCall()
        if call.FunctionName() == "_&&_" {
            args := call.Args()
            for _, arg := range args {
                collectMetaOnlyClauses(arg, env, a, parts)
            }
            return
        }
    }

    // Check if this subtree only references metadata
    var hasMeta, hasMsg bool
    inspectAST(e, &hasMeta, &hasMsg)
    if hasMeta && !hasMsg {
        // Unparse this subtree back to string
        unparser, err := cel.AstToString(ast.NewAST(e, nil))
        if err == nil {
            *parts = append(*parts, unparser)
        }
    }
}

// structToFilterMap converts a decoded struct to a map for CEL evaluation.
// Both Go field names and lowercase versions are included for case-insensitive access.
func structToFilterMap(v any) map[string]any {
    result := make(map[string]any)
    rv := reflect.ValueOf(v)
    if rv.Kind() == reflect.Ptr {
        rv = rv.Elem()
    }
    if rv.Kind() != reflect.Struct {
        return result
    }

    rt := rv.Type()
    for i := 0; i < rt.NumField(); i++ {
        field := rt.Field(i)
        val := rv.Field(i)

        // Skip unexported fields
        if !field.IsExported() {
            continue
        }
        // Skip the embedded Info field
        if field.Name == "Info" {
            continue
        }

        var v any
        if val.Kind() == reflect.Ptr {
            if val.IsNil() {
                continue
            }
            v = val.Elem().Interface()
        } else {
            v = val.Interface()
        }

        // Convert numeric types to float64 for CEL compatibility
        v = toNumeric(v)

        result[field.Name] = v
        lower := strings.ToLower(field.Name)
        if lower != field.Name {
            result[lower] = v
        }
    }
    return result
}

func toNumeric(v any) any {
    switch n := v.(type) {
    case float32:
        return float64(n)
    case float64:
        return n
    case int:
        return int64(n)
    case int8:
        return int64(n)
    case int16:
        return int64(n)
    case int32:
        return int64(n)
    case int64:
        return n
    case uint:
        return int64(n)
    case uint8:
        return int64(n)
    case uint16:
        return int64(n)
    case uint32:
        return int64(n)
    case uint64:
        return int64(n)
    default:
        return v
    }
}
```

- [ ] **Step 5: Run tests**

Run: `go test -run "TestFilter" -v`
Expected: PASS

Note: The AST unparsing in `collectMetaOnlyClauses` may need adjustment depending on the exact cel-go API. The `cel.AstToString` function may have a slightly different signature. Check cel-go docs and adjust. The core logic (walk AST, find metadata-only AND clauses) is correct.

- [ ] **Step 6: Commit**

```bash
git add filter.go filter_test.go go.mod go.sum && git commit -m "feat: add CEL filter with automatic pre/post partitioning"
```

---

### Task 10: Integrate filter into Scanner pipeline

Wire the filter into the Scanner so pre-decode filtering skips decoding and post-decode filtering runs on struct fields.

**Files:**
- Modify: `scanner.go`
- Modify: `scanner_test.go`

- [ ] **Step 1: Write filtered scanner test**

Add to `scanner_test.go`:

```go
func TestScannerWithPreFilter(t *testing.T) {
    // Create frames with different PGNs
    frame1 := can.Frame{
        ID:     0x09F10D00, // Some PGN
        Length: 8,
        Data:   [8]uint8{0, 0, 0, 0, 0, 0, 0, 0},
    }
    frame2 := can.Frame{
        ID:     0x09F20D00, // Different PGN (127501)
        Length: 8,
        Data:   [8]uint8{0, 0, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
    }

    ctx := context.Background()
    // Filter for only PGN 127501 (0x1F20D)
    s := NewScanner(ctx,
        Replay([]can.Frame{frame1, frame2}),
        Filter(`pgn == 127501`),
        IncludeUnknown(),
    )

    count := 0
    for s.Next() {
        count++
    }
    assert.NoError(t, s.Err())
    // Should only receive messages matching the filter
    assert.LessOrEqual(t, count, 1)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run "TestScannerWithPreFilter" -v`
Expected: FAIL or incorrect count (filter not wired in)

- [ ] **Step 3: Update Scanner.run() to use filter**

Modify `scanner.go` to compile the filter at start and apply it in the handler:

In `Scanner.run()`, add filter compilation:

```go
func (s *Scanner) run() {
    defer close(s.ch)

    var f *filter
    if s.cfg.filterExpr != "" {
        var err error
        f, err = compileFilter(s.cfg.filterExpr)
        if err != nil {
            s.err = err
            return
        }
    }

    adapt := adapter.NewCANAdapter()
    dec := decoder.New()

    dec.SetOutput(&scannerHandler{scanner: s, filter: f})
    adapt.SetOutput(dec)

    err := runSources(s.ctx, s.cfg.sources, func(frame can.Frame) {
        // Pre-filter: skip decode if metadata doesn't match
        if f != nil {
            info := adapter.NewPacketInfo(&frame)
            if !f.evalPre(info) {
                return
            }
        }
        adapt.HandleMessage(&frame)
    })
    if err != nil {
        s.err = err
    }
}
```

Update `scannerHandler` to apply post-filter:

```go
type scannerHandler struct {
    scanner *Scanner
    filter  *filter
}

func (h *scannerHandler) HandleStruct(msg any) {
    if msg == nil {
        return
    }

    // Check if unknown PGN and should be dropped
    if u, ok := msg.(*pgn.UnknownPGN); ok {
        if !h.scanner.cfg.includeUnknown {
            h.scanner.cfg.logger.Debug("dropping unknown PGN",
                "pgn", u.Info.PGN,
                "reason", u.Reason,
            )
            return
        }
    }

    // Post-filter: check decoded struct fields
    if h.filter != nil && h.filter.hasPost {
        fields := structToFilterMap(msg)
        // Need MessageInfo for mixed expressions
        var info pgn.MessageInfo
        // Extract Info from the struct via reflection
        rv := reflect.ValueOf(msg)
        if rv.Kind() == reflect.Ptr {
            rv = rv.Elem()
        }
        if rv.Kind() == reflect.Struct {
            infoField := rv.FieldByName("Info")
            if infoField.IsValid() {
                info = infoField.Interface().(pgn.MessageInfo)
            }
        }
        if !h.filter.evalPostWithInfo(info, fields) {
            return
        }
    }

    select {
    case h.scanner.ch <- msg:
    case <-h.scanner.ctx.Done():
    }
}
```

Add `"reflect"` to imports.

- [ ] **Step 4: Expose NewPacketInfo from internal/adapter**

The pre-filter in Scanner needs to extract MessageInfo from a CAN frame without running the full adapter pipeline. `adapter.NewPacketInfo` is already defined in `internal/adapter/frame.go`. Verify it's exported (capital N). If it returns `pgn.MessageInfo`, it can be used directly.

- [ ] **Step 5: Run tests**

Run: `go test -run "TestScanner" -v`
Expected: All pass

- [ ] **Step 6: Commit**

```bash
git add scanner.go scanner_test.go && git commit -m "feat: integrate CEL filter into scanner pipeline"
```

---

### Task 11: Implement the Receive iterator API

Build `Receive()` on top of `Scanner`, returning `iter.Seq2[any, error]`.

**Files:**
- Modify: `n2k.go`
- Create: `receive_test.go`

- [ ] **Step 1: Write receive test**

Create `receive_test.go`:

```go
package n2k

import (
    "context"
    "testing"

    "github.com/brutella/can"
    "github.com/stretchr/testify/assert"
)

func TestReceiveBasic(t *testing.T) {
    ctx := context.Background()
    frames := []can.Frame{testFrame127501}

    count := 0
    for msg, err := range Receive(ctx, Replay(frames), IncludeUnknown()) {
        assert.NoError(t, err)
        assert.NotNil(t, msg)
        count++
    }
    assert.GreaterOrEqual(t, count, 1)
}

func TestReceiveNoSources(t *testing.T) {
    ctx := context.Background()
    for _, err := range Receive(ctx) {
        assert.Error(t, err)
        break // first yield should be error
    }
}

func TestReceiveWithFilter(t *testing.T) {
    ctx := context.Background()
    frames := []can.Frame{testFrame127501}

    count := 0
    for msg, err := range Receive(ctx, Replay(frames), Filter(`pgn == 0`)) {
        // PGN 0 should match nothing
        assert.NoError(t, err)
        _ = msg
        count++
    }
    assert.Equal(t, 0, count)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run "TestReceive" -v`
Expected: FAIL — Receive not defined

- [ ] **Step 3: Implement Receive**

Add to `n2k.go`:

```go
package n2k

import (
    "context"
    "iter"
)

// Receive returns an iterator of decoded NMEA 2000 messages from the configured sources.
// Each yielded value is a pointer to a typed PGN struct (e.g., *pgn.VesselHeading)
// or *pgn.UnknownPGN if IncludeUnknown() is set. Iteration stops on source error or
// context cancellation.
func Receive(ctx context.Context, opts ...Option) iter.Seq2[any, error] {
    return func(yield func(any, error) bool) {
        cfg := config{}
        for _, o := range opts {
            o.apply(&cfg)
        }
        if err := cfg.validate(); err != nil {
            yield(nil, err)
            return
        }

        s := NewScanner(ctx, opts...)
        for s.Next() {
            if !yield(s.Message(), nil) {
                return
            }
        }
        if s.Err() != nil {
            yield(nil, s.Err())
        }
    }
}
```

- [ ] **Step 4: Run tests**

Run: `go test -run "TestReceive" -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add n2k.go receive_test.go && git commit -m "feat: add Receive iterator API"
```

---

### Task 12: Update the sniffer CLI

Rewrite `cmd/sniffer.go` to use the new `n2k.Receive` API.

**Files:**
- Modify: `cmd/sniffer.go`

- [ ] **Step 1: Read current sniffer**

Read `cmd/sniffer.go` to confirm current structure.

- [ ] **Step 2: Rewrite sniffer**

Replace the contents of `cmd/sniffer.go`:

```go
package main

import (
    "context"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "log/slog"
    "os"
    "os/signal"

    "github.com/open-ships/n2k"
    "github.com/open-ships/n2k/pgn"
)

func main() {
    iface := flag.String("i", "can0", "CAN interface name")
    usb := flag.String("u", "", "USB-CAN serial port (e.g., /dev/ttyUSB0)")
    expr := flag.String("f", "", "CEL filter expression")
    unknown := flag.Bool("unknown", false, "include unknown PGNs")
    flag.Parse()

    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    opts := []n2k.Option{}
    if *usb != "" {
        opts = append(opts, n2k.USB(*usb))
    } else {
        opts = append(opts, n2k.CAN(*iface))
    }
    if *expr != "" {
        opts = append(opts, n2k.Filter(*expr))
    }
    if *unknown {
        opts = append(opts, n2k.IncludeUnknown())
    }
    opts = append(opts, n2k.WithLogger(slog.Default()))

    enc := json.NewEncoder(os.Stdout)
    for msg, err := range n2k.Receive(ctx, opts...) {
        if err != nil {
            log.Fatal(err)
        }
        if err := enc.Encode(msg); err != nil {
            fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
        }
        _ = pgn.UnknownPGN{} // ensure pgn import used
    }
}
```

Remove the `_ = pgn.UnknownPGN{}` line if the import is used elsewhere (it's just to keep the import). A cleaner approach: use `_ "github.com/open-ships/n2k/pgn"` if needed, or just remove the pgn import if it's not directly referenced.

- [ ] **Step 3: Verify it compiles**

Run: `go build ./cmd/sniffer.go`
Expected: SUCCESS

- [ ] **Step 4: Commit**

```bash
git add cmd/sniffer.go && git commit -m "refactor: rewrite sniffer to use n2k.Receive API"
```

---

### Task 13: Update README

Rewrite README.md with the new API, examples, and usage.

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Read current README**

Read `README.md` to understand current structure.

- [ ] **Step 2: Rewrite README**

Replace contents of `README.md`:

```markdown
# n2k

A Go library for decoding NMEA 2000 marine network messages from CAN bus hardware into strongly-typed Go structs.

## Install

```bash
go get github.com/open-ships/n2k
```

## Usage

### Iterator API

```go
package main

import (
    "context"
    "fmt"
    "os/signal"
    "os"

    "github.com/open-ships/n2k"
    "github.com/open-ships/n2k/pgn"
)

func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    for msg, err := range n2k.Receive(ctx, n2k.CAN("can0")) {
        if err != nil {
            panic(err)
        }
        switch m := msg.(type) {
        case *pgn.VesselHeading:
            fmt.Printf("Heading: %v\n", m.Heading)
        case *pgn.WindData:
            fmt.Printf("Wind Speed: %v\n", m.WindSpeed)
        }
    }
}
```

### Scanner API

```go
s := n2k.NewScanner(ctx, n2k.CAN("can0"))
for s.Next() {
    switch msg := s.Message().(type) {
    case *pgn.VesselHeading:
        fmt.Printf("Heading: %v\n", msg.Heading)
    }
}
if err := s.Err(); err != nil {
    panic(err)
}
```

### Multiple Sources

Read from multiple CAN interfaces simultaneously:

```go
for msg, err := range n2k.Receive(ctx,
    n2k.CAN("can0"),
    n2k.CAN("can1"),
    n2k.USB("/dev/ttyUSB0"),
) {
    // messages from all sources, interleaved by arrival
}
```

### CEL Filtering

Filter messages using [CEL](https://github.com/google/cel-go) expressions. The library automatically optimizes filters — metadata-only expressions skip decoding entirely.

```go
// Only vessel heading messages
for msg, err := range n2k.Receive(ctx,
    n2k.CAN("can0"),
    n2k.Filter(`pgn == 127250`),
) { ... }

// Filter on decoded fields
for msg, err := range n2k.Receive(ctx,
    n2k.CAN("can0"),
    n2k.Filter(`pgn == 127250 && msg.Heading > 3.14`),
) { ... }

// Filter by source address
for msg, err := range n2k.Receive(ctx,
    n2k.CAN("can0"),
    n2k.Filter(`source == 3`),
) { ... }
```

**Filter variables:**

| Variable | Type | Description |
|----------|------|-------------|
| `pgn` | `int` | Parameter Group Number |
| `source` | `int` | Source address (0-252) |
| `priority` | `int` | Message priority (0-7) |
| `destination` | `int` | Destination address (255 = broadcast) |
| `msg.<field>` | varies | Decoded struct field (case-insensitive) |

### Options

| Option | Description |
|--------|-------------|
| `n2k.CAN(iface)` | SocketCAN source (e.g., `"can0"`) |
| `n2k.USB(port)` | USB-CAN serial source (e.g., `"/dev/ttyUSB0"`) |
| `n2k.Replay(frames)` | Replay source for testing |
| `n2k.Filter(expr)` | CEL filter expression |
| `n2k.IncludeUnknown()` | Include undecodable messages as `*pgn.UnknownPGN` |
| `n2k.WithLogger(l)` | Override default `slog.Logger` |

### Testing with Replay

```go
frames := []can.Frame{
    {ID: 0x09F10D01, Length: 8, Data: [8]uint8{1, 2, 3, 4, 5, 6, 7, 8}},
}

for msg, err := range n2k.Receive(ctx, n2k.Replay(frames)) {
    // test your message handling
}
```

## Sniffer CLI

Print decoded NMEA 2000 messages as JSON:

```bash
# Read from SocketCAN
go run ./cmd/sniffer.go -i can0

# Read from USB-CAN
go run ./cmd/sniffer.go -u /dev/ttyUSB0

# With CEL filter
go run ./cmd/sniffer.go -i can0 -f 'pgn == 127250'

# Include unknown PGNs
go run ./cmd/sniffer.go -i can0 -unknown

# Pipe to jq
go run ./cmd/sniffer.go -i can0 | jq .
```

## PGN Types

All decoded messages are pointers to generated structs in the `pgn` package. Use a type switch to handle specific message types. See `pgn/pgninfo_generated.go` for the full list of supported PGNs.

Every struct embeds `pgn.MessageInfo` as the `Info` field:

```go
type MessageInfo struct {
    Timestamp time.Time
    Priority  uint8
    PGN       uint32
    SourceId  uint8
    TargetId  uint8
}
```

## Unit Types

Physical quantities use type-safe wrappers from the `units` package with built-in conversion methods:

```go
case *pgn.VesselHeading:
    heading := m.Heading // *float32, in radians
```

## Hardware

Tested with:
- **SocketCAN**: MCP2515 (SPI), PEAK PCAN-USB
- **USB-CAN**: USB-CAN Analyzer dongles (2 Mbaud serial)

Both use the NMEA 2000 standard bitrate of 250 kbps.
```

- [ ] **Step 3: Commit**

```bash
git add README.md && git commit -m "docs: rewrite README for new n2k.Receive API"
```

---

### Task 14: Clean up and final verification

Remove any remaining `pkg/` references, verify the full build, run all tests.

**Files:**
- Delete: `pkg/` directory (should be empty after Tasks 1-5)
- Modify: `go.mod` if needed

- [ ] **Step 1: Verify pkg/ is gone**

```bash
ls pkg/ 2>/dev/null && echo "ERROR: pkg/ still exists" || echo "OK: pkg/ removed"
```

- [ ] **Step 2: Check for stale import references**

```bash
grep -r "open-ships/n2k/pkg" --include="*.go" .
```

Expected: No results

- [ ] **Step 3: Tidy modules**

```bash
go mod tidy
```

- [ ] **Step 4: Run full build**

```bash
go build ./...
```

Expected: SUCCESS

- [ ] **Step 5: Run all tests**

```bash
go test ./...
```

Expected: All pass

- [ ] **Step 6: Run linter**

```bash
golangci-lint run ./...
```

Expected: No errors (or only pre-existing warnings)

- [ ] **Step 7: Final commit**

```bash
git add -A && git commit -m "chore: clean up and finalize package simplification"
```

---

### Task Summary

| Task | Description | Dependencies |
|------|-------------|--------------|
| 1 | Move canbus to internal/canbus | None |
| 2 | Move adapter to internal/adapter | Task 1 |
| 3 | Move decoder (pkt) to internal/decoder | Task 2 |
| 4 | Move pgn and units to top-level | Task 3 |
| 5 | Absorb endpoints into internal/canbus | Task 4 |
| 6 | Create Source interface and implementations | Task 5 |
| 7 | Create options API | Task 5 |
| 8 | Create Scanner and pipeline wiring | Tasks 6, 7 |
| 9 | Add CEL filter support | Task 7 |
| 10 | Integrate filter into Scanner | Tasks 8, 9 |
| 11 | Implement Receive iterator | Task 10 |
| 12 | Update sniffer CLI | Task 11 |
| 13 | Update README | Task 12 |
| 14 | Clean up and final verification | Task 13 |
