# Architecture

This document describes the high-level architecture of FlowSync, including components, data flow, and filesystem layout.

---

## FlowSync Components

FlowSync is divided into several focused components:

*   **`cmd/flowsync` (CLI)**: A client-facing CLI tool that developers interact with. It manages daemon control, repo activation, and hooks mapping.
*   **`cmd/flowsyncd` (Daemon)**: A background server running per user. It listens on an IPC unix domain socket, ingests hook events, manages session threads, runs debouncers, and writes journals.
*   **`internal/reflector`**: Integrates with local LLMs (typically Ollama). It parses raw transcripts and outputs structured, concise journal entries.
*   **`internal/journal`**: Manages the low-level storage layout, directory lookup based on project paths, file lock concurrency, and journal appends.
*   **`adapters/*`**: Contains custom shell and python hook shims designed for different AI harnesses (e.g., Claude Code, Codex). They intercept command inputs/outputs and forward them to FlowSync via socket events.

---

## Data Flow Pipeline

Below is the chronological path of a coding action recorded in FlowSync:

1.  **Event Generation**: An AI assistant (e.g., Claude Code) executes a tool call or edits a file.
2.  **Hook Execution**: The harness's post-turn hook calls `flowsync hook <harness>`, formatting the turn transcript as a JSON payload.
3.  **IPC Event Ingestion**: The CLI sends this JSON payload over the Unix domain socket to `flowsyncd`.
4.  **Debouncing & Threading**:
    *   The daemon checks if this belongs to an existing session thread.
    *   If a long gap occurs, a new Thread ID is assigned.
    *   High-frequency events are debounced to prevent generating too many checkpoints.
5.  **Checkpoint Reaching**: Once the activity settles (idle timeout or user command finishes):
    *   The raw transcript excerpt is written to the `checkpoint-logs/` directory.
    *   If journaling is enabled, the daemon schedules a task to reflect on the transcript.
6.  **Reflecting & Summarizing**: The Reflector invokes the configured model (e.g., `qwen2.5-coder:7b` via Ollama) with a special prompt template to generate a markdown summary.
7.  **Journal Appending**: The generated summary is formatted with YAML frontmatter (detailing files touched, harness used, timestamp) and safely appended to the project's `journal.md`.

---

## System Architecture Diagram

```mermaid
flowchart LR
	subgraph Local ["Client Workspace"]
		User[User / Harness]
		HookScript[Hook Script\nadapters/*]
		RepoFlowsync[<repo>/.flowsync]
	end

	subgraph DaemonRoot ["Daemon Home (~/.flowsync)"]
		Daemon[flowsyncd]\n(ingest, checkpoint, journal)
		JournalStore[(projects/<hash>/journal.md)]
		CheckpointLogs[(checkpoint-logs/*)]
		Reflector[internal/reflector\n(ollama or other)]
	end

	User --> HookScript -->|call| Daemon
	HookScript -->|optionally runs| RepoFlowsync
	Daemon --> JournalStore
	Daemon --> CheckpointLogs
	Daemon --> Reflector
	Reflector -->|writes summary| JournalStore
	RepoFlowsync -->|preferred storage| JournalStore

	style RepoFlowsync fill:#f9f,stroke:#333,stroke-width:1px
	style Daemon fill:#efe,stroke:#333,stroke-width:1px
	style Reflector fill:#eef,stroke:#333,stroke-width:1px
```

---

## Security and Filesystem Layout

*   **Local Privacy**: No data is uploaded to remote servers. All communication between the shims and the daemon happens over local loopback Unix sockets.
*   **Access Control**: Files and directories are created with strict permission masks (`0700` directory access, `0600` file access) to ensure only the operating system user running the daemon can read the captured transcripts and journals.
*   **Redirected Storage**: Local `.flowsync` folders under an active repository root are automatically detected and preferred, isolating logs to the project folder itself.
