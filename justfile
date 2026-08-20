golangci_lint_version := "v2.12.0"
secure_go_toolchain := "go1.26.6"
govulncheck_version := "v1.5.0"
gosec_version := "v2.27.1"

# list available recipes
default:
    @just --list

# ensure dev tools are installed at the right versions
setup:
    @if command -v golangci-lint >/dev/null && golangci-lint --version 2>&1 | grep -q "{{trim_start_match(golangci_lint_version, "v")}}"; then \
        echo "golangci-lint {{golangci_lint_version}} already installed"; \
    else \
        echo "Installing golangci-lint {{golangci_lint_version}}..."; \
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b "$(go env GOPATH)/bin" {{golangci_lint_version}}; \
    fi
    @if command -v govulncheck >/dev/null \
        && govulncheck -version 2>&1 | grep -q "govulncheck@{{govulncheck_version}}" \
        && go version "$(command -v govulncheck)" 2>&1 | grep -q "{{secure_go_toolchain}}"; then \
        echo "govulncheck {{govulncheck_version}} already installed"; \
    else \
        echo "Installing govulncheck {{govulncheck_version}}..."; \
        GOTOOLCHAIN={{secure_go_toolchain}} go install golang.org/x/vuln/cmd/govulncheck@{{govulncheck_version}}; \
    fi
    @if command -v gosec >/dev/null \
        && go version -m "$(command -v gosec)" 2>&1 | grep -q "github.com/securego/gosec/v2[[:space:]]*{{gosec_version}}" \
        && go version "$(command -v gosec)" 2>&1 | grep -q "{{secure_go_toolchain}}"; then \
        echo "gosec {{gosec_version}} already installed"; \
    else \
        echo "Installing gosec {{gosec_version}}..."; \
        GOTOOLCHAIN={{secure_go_toolchain}} go install github.com/securego/gosec/v2/cmd/gosec@{{gosec_version}}; \
    fi

# run all tests
test:
    go test ./...

# require every tracked Go file to match gofmt
format-check:
    @unformatted=$(git ls-files -z '*.go' | xargs -0 gofmt -l); \
        if [ -n "$unformatted" ]; then \
            echo "tracked Go files require gofmt:" >&2; \
            echo "$unformatted" >&2; \
            exit 1; \
        fi

# verify VERSION, changelog, and published-tag consistency
release-check:
    git fetch origin --tags --force
    bash scripts/release-check.sh

# run the public, reproducible protocol evidence suite (not NMEA certification)
conformance-local:
    go test -count=1 -v . -run '^(TestConformance|TestHeartbeat_|TestISORequest_|TestGroupFunction_|TestProtocolTransmission|TestRequiredProtocol|TestAdvisoryProtocol|TestTCPClientReconnect|TestObservationHub|TestReplayObservation)'
    go test -count=1 ./internal/actisense ./internal/adapter ./internal/claiming ./internal/ebl ./internal/framer ./internal/gateway ./internal/transport

# compare runtime PGN support against the current upstream schema
upstream-parity:
    UPSTREAM_PARITY=1 go test ./pgn -run TestUpstreamParity -count=1 -v

# regenerate source metadata from the pinned upstream schema revision
generate-pgn:
    go run ./cmd/pgngen

# run tests with verbose output
test-v:
    go test -v ./...

# run tests with race detector
test-race:
    go test -race ./...

# short, deterministic smoke run of every fuzz harness
fuzz-smoke:
    go test ./internal/actisense -run=^$ -fuzz=FuzzParser -fuzztime=2s
    go test ./internal/adapter -run=^$ -fuzz=FuzzCANAdapter -fuzztime=2s
    go test ./internal/canbus -run=^$ -fuzz=FuzzUSBCANParseFrames -fuzztime=2s
    go test ./internal/gateway -run=^$ -fuzz=FuzzActisenseReader -fuzztime=2s
    go test ./pgn -run=^$ -fuzz=FuzzDecodeEncodeMessage -fuzztime=2s

# run tests with coverage
test-cover:
    go test -coverpkg=./... -coverprofile=c.out ./...

# format all Go source files
fmt:
    gofmt -w .

# regenerate pgn/manifest.json from the initialized runtime PGN metadata
pgn-manifest:
    UPDATE_PGN_MANIFEST=1 go test ./pgn -run TestPGNManifestMatchesMetadata -count=1

# sync the pinned upstream source metadata and runtime PGN manifest
pgn-sync:
    go run ./cmd/pgngen
    UPDATE_PGN_MANIFEST=1 go test ./pgn -run TestPGNManifestMatchesMetadata -count=1

# check generated metadata and the manifest against the pinned schema revision
pgn-sync-check:
    go run ./cmd/pgngen --check
    go test ./pgn -run TestPGNManifestMatchesMetadata -count=1

# run go vet
vet:
    go vet ./...

# run golangci-lint
lint:
    golangci-lint run ./...

# run security review (vulnerability scan + static analysis)
secure:
    GOTOOLCHAIN={{secure_go_toolchain}} govulncheck ./...
    GOTOOLCHAIN={{secure_go_toolchain}} gosec -exclude-dir=.claude -exclude-generated -exclude=G115 ./...

# tidy go modules
tidy:
    go mod tidy
