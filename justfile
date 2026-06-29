golangci_lint_version := "v2.12.0"
secure_go_toolchain := "go1.25.11"

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
    @if command -v govulncheck >/dev/null; then \
        echo "govulncheck already installed"; \
    else \
        echo "Installing govulncheck..."; \
        go install golang.org/x/vuln/cmd/govulncheck@latest; \
    fi
    @if command -v gosec >/dev/null; then \
        echo "gosec already installed"; \
    else \
        echo "Installing gosec..."; \
        go install github.com/securego/gosec/v2/cmd/gosec@latest; \
    fi

# run all tests
test:
    go test ./...

# compare runtime PGN support against the current upstream schema
upstream-parity:
    UPSTREAM_PARITY=1 go test ./pgn -run TestUpstreamParity -count=1 -v

# regenerate upstream source-derived runtime metadata from upstream master
generate-pgn:
    go run ./cmd/pgngen

# run tests with verbose output
test-v:
    go test -v ./...

# run tests with race detector
test-race:
    go test -race ./...

# run tests with coverage
test-cover:
    go test -coverpkg=./... -coverprofile=c.out ./...

# format all Go source files
fmt:
    gofmt -w .

# regenerate pgn/manifest.json from the initialized runtime PGN metadata
pgn-manifest:
    UPDATE_PGN_MANIFEST=1 go test ./pgn -run TestPGNManifestMatchesMetadata -count=1

# sync generated upstream source metadata and regenerate the runtime PGN manifest
pgn-sync:
    go run ./cmd/pgngen
    UPDATE_PGN_MANIFEST=1 go test ./pgn -run TestPGNManifestMatchesMetadata -count=1

# check whether generated upstream source metadata and runtime PGN manifest are current
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
    GOTOOLCHAIN={{secure_go_toolchain}} gosec -exclude-generated -exclude=G115 ./...

# tidy go modules
tidy:
    go mod tidy
