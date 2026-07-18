# Developer Guide

This guide is for developers and contributors looking to compile, test, debug, and modify FlowSync.

---

## 🛠️ Building the Project

Ensure you have Go 1.20+ installed. To compile all binaries within the workspace:

```bash
# Build all sub-packages and binaries
go build ./...
```

To compile specific entrypoint commands:
```bash
# Build the client CLI
go build -o bin/flowsync ./cmd/flowsync

# Build the daemon server
go build -o bin/flowsyncd ./cmd/flowsyncd
```

---

## 🧪 Running Tests

FlowSync features unit tests for core debouncing logic, socket IPC, configuration parsing, and journal writing.

To run all unit tests in the repository:
```bash
# Run tests with verbose output
go test -v ./...
```

To run tests with code coverage:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 📂 Source Code Layout

*   `cmd/flowsync/`: Entrypoint code for the `flowsync` client command-line tool.
*   `cmd/flowsyncd/`: Entrypoint code for the `flowsyncd` long-running daemon.
*   `adapters/`: Script and configuration files used to hook into external client harnesses.
*   `internal/core/`: Event schemas, session threading logic, trigger classifications, and idle debouncers.
*   `internal/daemon/`: Connection handler loop, job queues, and background ticker coordination.
*   `internal/ipc/`: Unix domain socket Server and Client abstractions.
*   `internal/journal/`: Safe thread-locked storage engine for saving, reading, and appending markdown journals.
*   `internal/reflector/`: Client wrappers for calling Ollama and compiling completions.

---

## 🐛 Local Debugging & Testing

When developing modifications to the daemon or client, it is useful to run binaries directly:

1.  **Stop any running global daemon**:
    ```bash
    flowsync daemon stop
    ```
2.  **Launch the daemon locally in verbose foreground mode**:
    ```bash
    go run ./cmd/flowsyncd/main.go --tick-interval=10s
    ```
3.  **Run commands against the active local debug daemon**:
    In another shell window, you can run CLI commands:
    ```bash
    go run ./cmd/flowsync/main.go status
    ```
4.  **Send signals directly**:
    You can trigger shutdown paths by sending `SIGINT` (Ctrl+C) or `SIGTERM` to the daemon process to check connection cleanup and file lock releases.
