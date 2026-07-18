# Installation Guide

This page covers the prerequisites, build steps, and OS-specific setup notes needed to get FlowSync running on your system.

---

## 📋 Prerequisites

Before installing FlowSync, ensure your system has the following tools installed:

1.  **Go (version 1.20 or newer)**: Required for building the CLI and daemon from source.
    *   To verify, run: `go version`
2.  **Ollama** (optional, but highly recommended): Used for local LLM reflection.
    *   Download from [ollama.com](https://ollama.com).
    *   Ensure the `ollama` command is available in your terminal's `PATH`.
3.  **direnv** (optional): Recommended for automating project-level activation.
    *   Available through most package managers (`brew`, `apt`, `pacman`).

---

## 🛠️ Build and Installation

To install the latest version of FlowSync directly from the source code, run the following:

```bash
# Navigate to your clone of the Codex-nightline/flowsync repo
cd /path/to/flowsync

# Build and install the client CLI
go install ./cmd/flowsync

# Build and install the background daemon
go install ./cmd/flowsyncd
```

This will place the `flowsync` and `flowsyncd` binaries into your `$GOPATH/bin` directory (typically `~/go/bin` on Unix/macOS or `%USERPROFILE%\go\bin` on Windows).

---

## ⚙️ Running the Daemon

Once installed, you can start the background process using the CLI:

```bash
# Start the daemon
flowsync daemon start
```

### Verification
Confirm that the daemon is running and check its socket connection:

```bash
flowsync daemon status
```

You should see output indicating that the socket at `~/.flowsync/daemon.sock` is active and that the daemon is listening for events.

---

## 🖥️ Operating System Specific Notes

### Windows & WSL (Windows Subsystem for Linux)
Because FlowSync uses Unix Domain Sockets for high-performance IPC between client hooks and the daemon, running under WSL is highly recommended for Windows developers.

*   **WSL Setup**: Build and run both `flowsync` and `flowsyncd` entirely within your WSL distribution (e.g., Ubuntu).
*   **PowerShell / cmd.exe**: If you are using Windows-native IDE tools but a WSL shell, you can bridge commands by running them through `wsl`:
    ```powershell
    wsl flowsync daemon status
    ```

### Linux & macOS
On Unix-based operating systems, FlowSync works out of the box with standard file permissions:
*   The daemon creates directories with strict `0700` and `0600` octal permissions to protect sensitive CLI/AI transcripts and reflections from unauthorized local access.

---

## 🤖 Configuring the Local LLM (Ollama)

FlowSync uses a local LLM to summarize development sessions. By default, it looks for an Ollama server running on `localhost:11434` and requests the `qwen2.5-coder:7b` model.

1.  Start the Ollama server:
    ```bash
    ollama serve
    ```
2.  Pull the default model (or another model of your choice):
    ```bash
    ollama pull qwen2.5-coder:7b
    ```
3.  Ensure your model name is passed to the daemon at startup if not using the default:
    ```bash
    flowsyncd --reflector-model="llama3.1:8b"
    ```
