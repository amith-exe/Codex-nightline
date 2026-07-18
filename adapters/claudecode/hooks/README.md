# Claude Code hook shims

These bash shims are installed into a user's `~/.claude/settings.json` so that Claude Code calls them on every `SessionStart`, `UserPromptSubmit`, `PostToolUse`, `Stop`, `PreCompact`, and `PostCompact` event. Each shim reads the hook payload from stdin/env and forwards it to the FlowSync daemon at `~/.flowsync/daemon.sock` over a unix-domain-socket connection. The shim exits immediately; the daemon does the work.

See [HOW_IT_WORKS.md](../../../HOW_IT_WORKS.md) for the shim+daemon rationale.

## v0 contents

- `flowsync-hook.sh` is the generic shim used by the installed hook events. It finds `flowsync` via `THREADMARK_BIN`, `PATH`, or the repo-local `./flowsync` / `./bin/flowsync` fallback, then execs `flowsync hook claude-code`. Set `THREADMARK_ROOT` or `THREADMARK_SOCKET` to make the shim forward to a non-default root/socket during tests.
- `session-start.sh`, `user-prompt-submit.sh`, `post-tool-use.sh`, `stop.sh`, `pre-compact.sh`, and `post-compact.sh` are thin event-named wrappers for installs that prefer one script per hook. `pre-compact.sh` is the high-priority immediate-fire path for context compaction; `post-compact.sh` records continuity metadata only and does not fire a separate checkpoint.
- `../settings.fragment.json` is a template fragment. Replace `__THREADMARK_HOOK_SCRIPT__` with an absolute path to `flowsync-hook.sh`, or run:

```sh
flowsync init claude-code --print
```

To install into project settings:

```sh
flowsync init claude-code --scope project
```

For first live testing, launch Claude Code with `THREADMARK_NO_JOURNAL=true` in the environment so an auto-spawned `flowsyncd` runs with journal writes suppressed.
