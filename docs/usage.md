# Usage and Examples

This page provides practical workflows, terminal commands, and scripts to help you query, manage, and inspect FlowSync journals.

---

## 🔍 Locating Your Journals and Logs

FlowSync stores project data based on the activation scope. The directory layout is identical for both global and local scopes:

*   **Global Project Storage**: `~/.flowsync/projects/<projectHash>/`
*   **Local Project Storage**: `<repo>/.flowsync/projects/<projectHash>/`
*   **Checkpoint Excerpts**: `<storageRoot>/checkpoint-logs/`

---

## 🐍 Script: Find the Journal for a Repository Path

Since project directories are hashed to prevent path collisions and handle nested structures, finding the correct journal for a specific directory can be done programmatically.

Here is a Python helper script to resolve a directory path to its corresponding `journal.md`:

```bash
# Run this from your terminal to print the active journal.md path for a repo
python - <<'PY'
import json
import glob
import sys
import os

target_dir = os.path.abspath(sys.argv[1] if len(sys.argv) > 1 else ".")
search_path = os.path.expanduser("~/.flowsync/projects/*/state.json")

found = False
for path in glob.glob(search_path):
    try:
        with open(path) as f:
            state = json.load(f)
            if os.path.abspath(state.get("working_dir", "")) == target_dir:
                print(path.replace("state.json", "journal.md"))
                found = True
                break
    except Exception:
        pass

if not found:
    # Check project-local storage if global directory didn't match
    local_path = os.path.join(target_dir, ".flowsync", "projects")
    for path in glob.glob(os.path.join(local_path, "*", "state.json")):
        try:
            with open(path) as f:
                state = json.load(f)
                if os.path.abspath(state.get("working_dir", "")) == target_dir:
                    print(path.replace("state.json", "journal.md"))
                    found = True
                    break
        except Exception:
            pass

if not found:
    print(f"No FlowSync journal found active for: {target_dir}", file=sys.stderr)
PY /abs/path/to/your/repo
```

---

## 🛠️ Common Workflows

### Session Capture
Just write code and interact with your AI assistant (Codex/Claude Code) as you normally would. The hooks run in the background on every turn, streaming transcripts to the daemon. 
If you pause or finish a task:
*   The daemon notices the inactivity (defaults to 5 minutes of idle, or based on command triggers).
*   A checkpoint is automatically created.
*   The reflector processes the transcript and appends a markdown summary to `journal.md`.

### Viewing Recent Logs and State
To see the active configuration paths and how many entries have been logged, run:
```bash
flowsync status
```

To tail the journal entries directly in your console:
```bash
# Print the last 3 journal entries
flowsync journal --last 3
```

---

## ⚙️ Daemon Management

To ensure everything runs smoothly, use these commands to query and manage the background daemon process:

*   **Check Daemon Status**:
    ```bash
    flowsync daemon status
    ```
*   **Stop the Daemon**:
    ```bash
    flowsync daemon stop
    ```
*   **Restart the Daemon**:
    ```bash
    flowsync daemon restart
    ```
*   **Wipe Local Storage** (Deletes local caches and history, use with caution!):
    ```bash
    flowsync purge
    ```
