# CLI Reference

This page documents the commands, subcommands, and flag options supported by the `flowsync` command-line utility.

---

## Global Usage
```bash
flowsync <command> [options]
```

### Global Flags
*   `-h`, `--help`, `help`: Display help text and available commands.
*   `-v`, `--version`, `version`: Print the FlowSync version string.

---

##  Commands

### 1. `activate`
Prepares a repository or user directory to accept FlowSync hooks and starts the daemon if needed.
```bash
flowsync activate [--scope project|user] [--quiet]
```
*   `--scope`: Specifies whether hooks should be local to the current directory (`project`) or home folder (`user`). Default is `project`.
*   `--quiet`: Silences all output logs except errors.

---

### 2. `init`
Generates and writes adapter hook files for specific client harnesses.
```bash
flowsync init <claude-code|codex>
```
*   `claude-code`: Writes hook files compatible with Claude Code.
*   `codex`: Writes hook files compatible with Codex.

---

### 3. `daemon`
Controls the life cycle of the long-running background daemon process.
```bash
flowsync daemon <start|stop|restart|status>
```
*   `start`: Starts the daemon process in the background.
*   `stop`: Sends a termination signal to the running daemon process.
*   `restart`: Safely stops and restarts the daemon.
*   `status`: Returns whether the daemon socket is alive and details its active process ID.

---

### 4. `status`
Displays current path configurations and details for the active workspace.
```bash
flowsync status [--working-dir .]
```
*   Prints the working directory, the resolved global root directory, the calculated project hash, and directory locations of the journal and checkpoint logs.

---

### 5. `journal`
Retrieves and prints recent entries written to the project's journal.
```bash
flowsync journal [--last N] [--working-dir .]
```
*   `--last`: Specifies the number of recent journal entries to print. Default is `3`.
*   `--working-dir`: Path to the active repository. Default is `.`.

---

### 6. `enable` / `disable`
Overrides or implements project-level opt-out flags.
```bash
flowsync disable  # Creates a opt-out file in .flowsync to block logging
flowsync enable   # Removes the opt-out file to re-enable logging
```

---

### 7. `hook`
Internal hook bridge endpoint. This is utilized by adapter scripts and should not typically be called manually by users.
```bash
flowsync hook <harness> [--example]
```

---

### 8. `purge`
Deletes all files in the current repository's `.flowsync` folder. **Caution: This action is destructive and cannot be undone.**
```bash
flowsync purge
```

---

## Daemon CLI Configuration (`flowsyncd`)

When launching the daemon directly (or when configured via config templates), you can pass the following options:

*   `--root <path>`: Override the default global configuration directory (default: `~/.flowsync`).
*   `--socket <path>`: Override the default IPC unix socket path (default: `~/.flowsync/daemon.sock`).
*   `--reflector-model <model>`: Specify which local LLM model to query via Ollama (default: `qwen2.5-coder:7b`).
*   `--reflector-timeout <duration>`: Maximum execution timeout for reflection LLM calls (default: `2m`).
*   `--tick-interval <duration>`: Interval at which the daemon checks for idle sessions (default: `30s`).
*   `--idle-trigger-after <duration>`: Duration of inactivity before a checkpoint is automatically generated (default: `5m`).
