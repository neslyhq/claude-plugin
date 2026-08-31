# nesly for Claude Code

Deploy applications to [nesly](https://nesly.net) without leaving
Claude Code. One install brings:

- **MCP tools** — deploy, diagnose, logs, env, addons, domains,
  rollback, usage — backed by the `nesly` CLI's credentials.
- **A deploy skill** — Claude knows the platform's contract
  (`$PORT`, `0.0.0.0`), the build detection, and the
  diagnose-first loop for failed deploys.

## Install

```
# prerequisite: the CLI
brew install neslyhq/tap/nesly
nesly login

# the plugin (inside Claude Code)
/plugin marketplace add neslyhq/claude-plugin
/plugin install nesly@nesly
```

Headless (CI, a server): `nesly token create agent` once, then set
`NESLY_TOKEN=nesly_pat_...` in the environment.

## Docs

- Agent contract: https://nesly.net/docs/agents/
- One-page briefing: https://nesly.net/llms.txt

This repository is published from the platform monorepo; open issues
here, changes land via the platform release process.
