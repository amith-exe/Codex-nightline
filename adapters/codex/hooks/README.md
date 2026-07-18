# Codex hook shims

These bash shims are installed into `.codex/hooks.json` or `~/.codex/hooks.json` so Codex calls FlowSync on `SessionStart`, `UserPromptSubmit`, `PostToolUse`, `Stop`, `PreCompact`, and `PostCompact`. Each shim reads the hook payload from stdin and forwards it to the FlowSync daemon at `~/.flowsync/daemon.sock` over a unix-domain-socket connection. The shim exits quickly; the daemon does the work.

See [HOW_IT_WORKS.md](../../../HOW_IT_WORKS.md) for the shim+daemon rationale.

## v0 contents

- `flowsync-hook.sh` is the generic shim used by all hook events. It finds `flowsync` via `THREADMARK_BIN`, `PATH`, or the repo-local `./flowsync` / `./bin/flowsync` fallback, then execs `flowsync hook codex`. Set `THREADMARK_ROOT` or `THREADMARK_SOCKET` to make the shim forward to a non-default root/socket during tests.
- `session-start.sh`, `user-prompt-submit.sh`, `post-tool-use.sh`, `stop.sh`, `pre-compact.sh`, and `post-compact.sh` are thin event-named wrappers for installs that prefer one script per hook. `pre-compact.sh` is the high-priority immediate-fire path for context compaction; `post-compact.sh` records continuity metadata only and does not fire a separate checkpoint.
- `../hooks.fragment.json` is a template fragment. Replace `__THREADMARK_HOOK_COMMAND__` with an absolute path to `flowsync-hook.sh`, or run:

```sh
flowsync init codex --print
```

To install into project hooks:

```sh
flowsync init codex --scope project
```

Codex requires one-time review for project-local command hooks. After installing, open Codex in the repo and run `/hooks` to trust the FlowSync hooks.

For first live testing, launch Codex with `THREADMARK_NO_JOURNAL=true` in the environment so an auto-spawned `flowsyncd` runs with journal writes suppressed.
