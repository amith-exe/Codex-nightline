# Project Activation

To record developer activities and generate journal summaries, a repository must be "activated" within FlowSync. This guide details what happens during activation, how to activate manually, and how to automate the process.

---

## 🛠️ What Does Activation Do?

When you run `flowsync activate` in a project directory:
1.  **Installs Harness Hooks**: Creates or updates configuration files (e.g., `.claude` or `.codex` settings) in the target directory to load FlowSync's hook shim scripts.
2.  **Launches the Daemon**: Automatically spawns a background `flowsyncd` instance if one is not already active.
3.  **Localizes Storage**: When scoped to a project (the default), it creates a `.flowsync/` folder in the root of the repository. This ensures that session transcripts, project state, and journals are stored inside the project rather than in the global home directory.

---

## 💻 Activation Commands

### Manual Project-Level Activation (Recommended)
Prepares the current directory and installs local-scoped configurations:
```bash
cd /path/to/your/project
flowsync activate
```

### Global User-Level Activation
Installs hooks into global settings files in your home directory:
```bash
flowsync activate --scope user
```

---

## 🤖 Automating Activation

To avoid typing `flowsync activate` every time you navigate to a project, choose one of these auto-activation methods:

### Option A: Using `direnv` (Highly Recommended)
If you use `direnv` for managing environment variables, you can automatically activate FlowSync when entering a project directory.

1. Create or edit an `.envrc` file in the root of your project:
   ```bash
   # Add this line
   flowsync activate --quiet
   ```
2. Approve the environment change:
   ```bash
   direnv allow
   ```

Now, every time you `cd` into the directory, `direnv` will silently start the daemon and update the hooks if necessary.

### Option B: Using Git Hooks (Global)
You can configure Git to trigger activation whenever a repository is checked out or created.

1.  Add a `post-checkout` hook into your global git template directory (often `~/.git-templates/hooks/post-checkout`):
    ```bash
    #!/bin/sh
    if command -v flowsync >/dev/null 2>&1; then
        flowsync activate --quiet --scope project &
    fi
    ```
2.  Make sure the script is executable:
    ```bash
    chmod +x ~/.git-templates/hooks/post-checkout
    ```

---

## 📁 Local `.flowsync` Directory Layout

Once active, the `.flowsync` directory in your repository will contain:

*   **`projects/<hash>/state.json`**: Tracks the thread count, last activity timestamp, and metadata configuration for the repository.
*   **`projects/<hash>/journal.md`**: The primary markdown file containing your chronological development reflections.
*   **`checkpoint-logs/`**: Raw logs/excerpts of terminal commands and file edits captured during each AI turn.
