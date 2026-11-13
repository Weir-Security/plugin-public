# Weir for Claude Code

Weir checks each action Claude Code proposes before it runs. It blocks or asks about the
risky ones (secrets, data leaving the company, force pushes, changes to shared systems),
saves a checkpoint before a delete so it can be undone, and logs labels: never a
command, a path or a file's content. The one exception is what the person asked, kept
word for word so it is clear whether the agent did it, with passwords, keys, tokens and
card numbers replaced by `[secret]`.

## Install

```
/plugin marketplace add Weir-Security/plugin-public
/plugin install weir@weir-security
```

Restart Claude Code afterwards.

For a whole company, an admin installs it through managed settings instead
(`extraKnownMarketplaces`, `enabledPlugins`, and `WEIR_SERVER` and `WEIR_INSTALL_KEY` in
`env`), so nobody installs anything by hand.

## Weir comes from your company's server

This plugin does not hold Weir itself, only a small launcher. The first time a hook runs,
the launcher downloads Weir from the company's Weir server (`WEIR_SERVER`) with the
company's install key (`WEIR_INSTALL_KEY`), or with the token of a computer enrolled by
hand (`weir enroll`). Without one of them nothing is downloaded. The download is checked
against the SHA-256 the server sends and kept in Weir's folder (`bin/`); at the start of a
session, at most every 6 hours, the launcher asks the server for a newer one.

## What it runs

`hooks/hooks.json` registers Weir for seven events: SessionStart, UserPromptSubmit,
PreToolUse, PostToolUse, PostToolUseFailure, PermissionDenied and SessionEnd. Each is an
exec-form hook (no shell) with a 10 second timeout (60 for SessionStart, which may
download Weir). PreToolUse has `"onFailure": "block"`: a Weir that can't start, times out
or crashes blocks that tool call instead of letting it through (Claude Code 2.1.295 and
later; an older one ignores the field and lets it through). The other hooks let the work
go on when they fail. A launcher that has not got Weir yet is not a failure: it lets the
work go on and says so once, at the session's start.

The launchers live in `runtime/`, not `bin/`: Claude Code puts a plugin's `bin/` on the
Bash tool's PATH, and Weir's commands are for a person, not for the agent.

- macOS and Linux: `runtime/weir` is a short sh script that runs
  `runtime/weir-<os>-<arch>`, the launcher for that system.
- Windows: the hook's command resolves to `runtime/weir.exe`. To be confirmed on a clean
  Windows machine (WEIR-28).

Weir keeps its log, checkpoints and settings in its own folder in the user's home
directory, never inside the plugin, so an update does not lose them.

## Update

Weir updates itself from the company's server. Each Weir release also rebuilds the
launchers in `runtime/` (from `cmd/weir-launcher` in Weir-Security/Spotter) and raises the
version in `.claude-plugin/plugin.json`; `/plugin update weir@weir-security` picks it up.

## Uninstall

```
/plugin uninstall weir@weir-security
```

The log and checkpoints stay in Weir's folder until you delete them.
