# Sign-in and credentials (shared by setup, deploy, debug, rollback)

Contents: whoami · the `reason` values · logging in · MCP 401 · CI and
headless · what never to do.

## whoami

Call it before any MCP mutate or project lookup in the turn: MCP `whoami` or
`nesly whoami --json`. Fields: `signed_in` (bool), `reason` (only when false),
`retryable` (bool), the account and email when signed in.

Do not end a turn on a bare `signed_in: false`: most reasons are not a lost
session.

## The `reason` values

| reason | what it means | what you do |
|---|---|---|
| `api-unreachable` (`retryable: true`) | a session is stored, the API went quiet | wait a few seconds, call whoami again, up to three times; never ask for a login |
| `no-credential` | nothing stored on this machine | run `nesly login` yourself |
| `rejected` | the API refused the stored credential | run `nesly login` yourself |
| `credential-unreadable` | the stored credential cannot be read | run `nesly login` yourself |
| `account-mismatch` | session fine, wrong account selected | `nesly account ls`, then `nesly account use <name>` |

## Logging in

`nesly login` opens the browser; the person only approves there. The command
bounds its own wait (10 minutes, or the code's expiry): do not wrap it in
`timeout`, which stock macOS does not even have. Say in one line that a browser
window is waiting for them, then call whoami again and carry on.

With a password in the environment there is a no-browser path:
`NESLY_PASSWORD=… nesly login --email <account>` (plus `--otp` when MFA asks).

Only if the CLI is absent, or login returns an error, tell the person to run
`nesly login`, then `nesly mcp setup`, and restart the agent tool.

## MCP returns 401 / not signed in

Stop retrying MCP. The CLI shares the same credentials: run the same operation
as `nesly … --json` in the same turn. Do not hunt for secrets or other login
paths for MCP. A Cursor remote-OAuth MCP session is a different session from
the CLI's: on a machine with the CLI prefer `nesly mcp setup` (stdio) so one
login is enough.

## CI and headless machines

`nesly token create agent` once (in a signed-in shell), then set
`NESLY_TOKEN=nesly_pat_…` in the environment of the job. The token goes in the
CI secret store, never in the chat, a commit or a log.

## Never

- Invent Vercel, Netlify, Fly, Railway or another cloud.
- Paste a token or password into the chat, or ask the person to.
- Search the home directory, `~/.secrets`, other projects or stray `.env` files
  for credentials to give the app or the CLI.
