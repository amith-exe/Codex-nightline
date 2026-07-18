# Quickstart Guide

This quickstart guide gets you up and running with the FlowSync daemon, shows you how to activate a project, and walks you through simulating a hook to generate your first journal entry.

---

##  Step-by-Step Walkthrough

### Step 1: Build and Install Binaries
Ensure you are in the directory containing the FlowSync source code, then run:

```bash
# Compile and install both the CLI and daemon binaries
go install ./cmd/flowsync
go install ./cmd/flowsyncd
```
This compiles the tools and places them into your `GOBIN` directory (e.g., `~/go/bin`).

---

### Step 2: Start the Daemon
Start the FlowSync background service. It will initialize its default storage under your user profile (`~/.flowsync/`):

```bash
flowsync daemon start
```

To verify the daemon started successfully and is listening on the Unix socket, check the status:
```bash
flowsync daemon status
```

---

### Step 3: Activate the Project
Change directories to the root of a project repository where you want to capture your coding sessions, and run:

```bash
cd /path/to/your/project
flowsync activate
```

> [!TIP]
> Running `flowsync activate` creates a local `.flowsync` folder inside the project. When active, all session excerpts and journals are saved directly inside this directory instead of the global `~/.flowsync` directory, keeping project logs localized.

---

### Step 4: Simulate a Hook Event
To test that your hook shims, socket IPC, and background logging are functional, run the built-in mock hook:

```bash
# Simulates a coding event from a Codex harness
flowsync hook codex --example
```
This command sends a mock tool call event to the daemon, triggers a checkpoint, and completes a reflection cycle.

---

### Step 5: Verify the Output Files
Check the output folder to verify the daemon recorded your event:

1.  **Project Journal**: 
    Look at the generated journal file to see the formatted summary:
    ```bash
    cat .flowsync/projects/<projectHash>/journal.md
    ```
2.  **Checkpoint Excerpts**:
    Inspect the raw transcripts collected during the session:
    ```bash
    ls .flowsync/checkpoint-logs/
    ```

---

##  What's Next?
*   Configure **Ollama** locally to automatically generate summaries for real-world coding turns (see [Installation](installation.md)).
*   Integrate hook shims into your active IDE extension or terminal environment (see [Project Activation](activate.md)).
