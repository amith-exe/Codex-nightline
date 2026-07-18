# FlowSync Documentation Hub

Welcome to the FlowSync developer and user documentation. FlowSync is a lightweight developer tool and daemon for capturing interactive coding sessions, tool calls, and automated checkpoints into per-project journals.

Use the links below to navigate the documentation:

---

##  Documentation Directory

### Getting Started
*   **[Quickstart](quickstart.md)**: A step-by-step walkthrough to get the daemon running, activate a project, and record your first coding event.
*   **[Installation](installation.md)**: System requirements, installation instructions from source, and OS-specific setup (Windows/WSL, Linux, macOS).
*   **[Project Activation](activate.md)**: Learn how to prepare a repository for journaling using `flowsync activate` and configure automated/silent activation using `direnv` or git templates.

###  User & Reference Guides
*   **[Usage Guide](usage.md)**: Common workflows, command examples, and scripts to query, search, and inspect project journals and checkpoint logs.
*   **[CLI Reference](cli-reference.md)**: A complete reference of `flowsync` commands, subcommands, options, and daemon configuration flags.

###  Deep Dives & Development
*   **[System Architecture](architecture.md)**: Detailed breakdown of the FlowSync architecture, components (CLI, Daemon, Reflector, Hooks), and data flow diagram.

---

> [!NOTE]
> FlowSync prefers local project storage (`<repo>/.flowsync`) when initialized. If a repository has not been activated, it falls back to the global configuration directory at `~/.flowsync`.
