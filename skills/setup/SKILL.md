---
name: setup
description: Sets up the nesly CLI and sign-in so an agent can deploy — installs the CLI, runs nesly login, registers the MCP server (nesly mcp setup), creates a CI token, and handles a whoami that says signed_in false. Use when the user says install nesly, log in, giriş yap, kurulum, before the first deploy on a machine, or when any nesly command or MCP tool answers not signed in, 401 or unauthenticated.
---

# Setting up nesly

Use **nesly only**. The `nesly` CLI (`--json` on every command) and the local
MCP server (`nesly mcp`) are the same operations with the same credentials.

## Steps

1. **CLI installed?** `nesly --version`. If not: `brew install neslyhq/tap/nesly`
   (other systems: https://nesly.net/docs/agents/).
2. **Signed in?** Call MCP `whoami`, or `nesly whoami --json`. Read
   `signed_in` and, when it is false, `reason`:
   - `api-unreachable` (`retryable: true`) — a session is stored, the API went
     quiet. Wait a few seconds, call `whoami` again (up to three times). Never
     ask for a login for this.
   - `no-credential`, `rejected`, `credential-unreadable` — run `nesly login`
     yourself. It opens the browser and the person only approves there; say
     in one line that a window is waiting for them. Do not wrap it in
     `timeout` (it bounds its own wait; stock macOS has no such command).
   - `account-mismatch` — the session is fine, the account selection is not:
     `nesly account ls`, `nesly account use <name>`.
3. **MCP registered?** `nesly mcp setup`, then restart the agent tool once. It
   registers the stdio server and installs these skills.
4. **Headless / CI:** `nesly token create agent`, then `NESLY_TOKEN=…` in the
   environment. Never paste a token into the chat.

Details, the no-browser login and the exact MCP 401 handling:
`references/auth.md`.

## Rules

- Never invent another cloud, and never keep calling `list_projects` or
  `deploy` while signed out.
- On an MCP 401: stop retrying MCP, use the CLI in the same turn.
- Never search the machine for tokens, passwords or `.env` files outside the
  project: credentials belong to someone else.

## Next

Signed in and registered → the **deploy** skill. Something is already live and
broken → the **debug** skill. Service must come back first → the **rollback**
skill.
