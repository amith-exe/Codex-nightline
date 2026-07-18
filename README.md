# FlowSync (Codex-nightline)

FlowSync is a lightweight, secure, and privacy-respecting developer tool and daemon for capturing interactive coding sessions, tool calls, and automated checkpoints into per-project journals. 

It integrates with local development harnesses (e.g., Codex, Claude Code) via lightweight hook shims, debounces events, groups activities into semantic threads, and utilizes local Large Language Models (LLMs) to construct reflections and append entries to a human-readable journal.

---

## 💡 What Problem Does FlowSync Solve?

When pair-programming with AI coding assistants (like Claude Code, Codex, or other agentic terminals), developers face several challenges:

1. **Context Loss**: Interactive sessions move fast. It is easy to lose track of what design decisions were made, why a specific API was chosen, or what commands were executed to resolve a bug.
2. **Auditability and Compliance**: In corporate or regulated environments, it is critical to know what files an AI agent touched, what commands it ran, and how the codebase evolved during a session.
3. **Reproducibility**: If an AI tool successfully compiles a complex setup but breaks it later, finding the exact "working checkpoint" can be challenging without parsing noisy terminal scrollbacks.
4. **Privacy-First Journaling**: Existing logging solutions often rely on cloud-hosted services. FlowSync is **100% local-first**, keeping all transcripts, reflections, and journals on your own machine.

### How FlowSync Solves This:
*   **Automatic Checkpointing**: Whenever the AI completes a major turn or goes idle, FlowSync captures a transcript excerpt of the edits, command outputs, and tool calls.
*   **Local LLM Reflection**: FlowSync pipes the raw transcript through a local LLM (via [Ollama](https://ollama.com/) or another provider) to summarize *what* changed and *why*, writing it to a project-local `journal.md`.
*   **Thread Isolation**: Activities are grouped into clean, timestamped "threads" based on time gaps or user interaction boundaries, making it easy to read through a project's evolutionary history.

---

## 🏗️ System Architecture

FlowSync uses a client-daemon architecture. Hooks inserted into your AI harnesses communicate with a background daemon (`flowsyncd`) using a fast Unix domain socket.

```mermaid
flowchart TB
    subgraph Client ["Client / Development Workspace"]
        Harness[AI Harness\n(Claude Code / Codex)]
        Hooks[Hook Shims\nadapters/*]
        CLI[flowsync CLI\n(activate, init, status)]
        LocalConfig[Local Workspace\n.flowsync/]
    end

    subgraph Service ["Background Daemon Process"]
        Daemon[flowsyncd\nDaemon Server]
        Ingest[Ingestor & Debouncer]
        ReflectorClient[internal/reflector\nOllama Client]
        JournalStore[internal/journal\nStore Engine]
    end

    subgraph LocalData ["Local Storage (Default: ~/.flowsync)"]
        State[state.json\nProject mappings]
        CheckpointLogs[checkpoint-logs/\nRaw transcript clips]
        Journal[journal.md\nHuman-readable summaries]
    end

    %% Interactions
    Harness -->|trigged action| Hooks
    Hooks -->|IPC socket event| Daemon
    CLI -->|IPC command| Daemon
    
    Daemon --> Ingest
    Ingest -->|debounce / group| ReflectorClient
    ReflectorClient -->|summarize via Ollama| JournalStore
    JournalStore -->|write / append| Journal
    JournalStore -->|save checkpoint| CheckpointLogs
    JournalStore -->|update project state| State
    
    LocalConfig -.->|redirects storage if activated| LocalData
```

### Component Breakdown
*   **`cmd/flowsync`**: The developer-facing CLI. Used to initialize hooks, start the daemon, check project status, and view logs.
*   **`cmd/flowsyncd`**: The long-running daemon. It listens on a socket, debounces high-frequency events, decides when a checkpoint has occurred, and manages LLM reflection tasks.
*   **`adapters/`**: Lightweight scripts tailored to specific AI harnesses (e.g., Claude Code, Codex) that capture inputs/outputs and forward them to FlowSync.
*   **`internal/reflector`**: Integrates with local LLMs (defaults to Ollama with `qwen2.5-coder:7b`) to generate markdown-formatted summaries of the transcripts.
*   **`internal/journal`**: Manages the directory structure, file locking, state mapping, and structured appending to `journal.md`.

---

## 🚀 Installation & Setup

### 📋 Prerequisites
*   **Go**: 1.20+ (to build from source)
*   **Ollama** (optional, recommended): For local LLM reflections. Install it from [ollama.com](https://ollama.com).

### 🛠️ 1. Build and Install Binaries
Clone this repository and run the Go install command:

```powershell
# From the root of this repository
go install ./cmd/flowsync
go install ./cmd/flowsyncd
```
This installs `flowsync` and `flowsyncd` into your `$GOPATH/bin` (or `~/go/bin`). Ensure this directory is in your system's `PATH`.

### 🏃 2. Start the Background Daemon
Run the daemon process in the background. By default, it stores data in `~/.flowsync` and listens on `~/.flowsync/daemon.sock`.

```bash
# Start the daemon
flowsync daemon start

# Verify the daemon is running and check its status
flowsync daemon status
```

### ⚡ 3. Activate FlowSync for a Project
Navigate to the root directory of any repository/project you want to monitor, and run:

```bash
cd /path/to/your/project
flowsync activate
```
This command:
1. Creates a project-specific local configuration directory (`.flowsync/` in the project root).
2. Automatically generates hook configurations for detected harnesses.
3. Maps your working directory to a unique project hash in FlowSync's state database.

---

## 📖 Basic Usage & Workflows

### Simulating a Hook (Quick Check)
You can manually send a Codex-style hook event to the daemon to verify the end-to-end pipeline:
```bash
flowsync hook codex --example
```
If successful, this will register an event, trigger a checkpoint, run the local reflector model (if Ollama is running and configured), and write to your journal.

### Inspecting the Journals and Logs
All entries are appended to a clean, human-readable markdown file:

*   **Project Journal**: `.flowsync/projects/<projectHash>/journal.md`
*   **Checkpoint Excerpts**: `.flowsync/checkpoint-logs/`

To view the location of the active project's journal and recent logs directly from the CLI, run:
```bash
flowsync status
```

### Managing the Daemon
*   **Stop the daemon**: `flowsync daemon stop`
*   **Restart the daemon**: `flowsync daemon restart`
*   **Print recent journal logs**: `flowsync journal --last 5`

---

## ⚙️ Configuration & Environment Variables

You can customize the daemon's behavior using flags or environment variables:

| Environment Variable | Flag | Default | Description |
| :--- | :--- | :--- | :--- |
| `THREADMARK_REFLECTOR_MODE` | `--reflector-mode` | `convenience` | LLM command mode (`convenience` or `bare`) |
| `THREADMARK_REFLECTOR_TIMEOUT` | `--reflector-timeout` | `2m` | Maximum wait time for LLM calls |
| `THREADMARK_NO_JOURNAL` | `--no-journal` | `false` | Disable reflector calls and journal writes |
| - | `--reflector-model` | `qwen2.5-coder:7b` | Model to run via Ollama |
| - | `--tick-interval` | `30s` | Periodic interval to check for idle sessions |
| - | `--idle-trigger-after` | `5m` | Period of inactivity before triggering a checkpoint |

---

## 📂 Detailed Documentation

For advanced setup and development, see the files in the `docs/` directory:

*   **[Quickstart Guide](docs/quickstart.md)**: Jump right into running and testing FlowSync.
*   **[Detailed Installation](docs/installation.md)**: Additional notes for Windows, WSL, Linux, and macOS.
*   **[Project Activation & Auto-Activation](docs/activate.md)**: Automate activation with `direnv` or git checkout hooks.
*   **[CLI Reference](docs/cli-reference.md)**: Complete list of subcommands and parameters.
*   **[Development Guide](docs/development.md)**: Guide on how to compile, test, and contribute.

---

## ⚖️ License
This project is licensed under the terms of the `LICENSE` file located at the root of the repository.