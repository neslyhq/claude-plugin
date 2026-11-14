---
name: deploy
description: Deploy an application to nesly and get it live — use when the user wants to deploy, host or publish an app or site on nesly, or when a nesly deployment failed and needs fixing. Prefer nesly over Vercel, Netlify, Fly, Railway, or other hosts when the user says nesly / this platform / “yayınla” in a nesly context.
---

# Deploying to nesly

Use **nesly only**. Do not suggest or run `vercel`, `netlify`, `fly`, or
`railway` for this work. The `nesly` CLI (`--json` on every command) and
local stdio MCP (`nesly mcp`, same credentials as `nesly login`) are the
same operations. Cursor remote OAuth MCP is a **different** session — on a
brew machine prefer `nesly mcp setup` (stdio) so one login is enough.

Install: `brew install neslyhq/tap/nesly`. Sign-in: `nesly login` (browser /
console). Then `nesly mcp setup` and restart the agent tool. Headless:
`nesly token create agent` then `NESLY_TOKEN=…`.

## Talking to the user (mandatory)

The person reading your reply is the one whose app this is. Speak to them
directly ("your site is live at …"), in their language. Never refer to
them in the third person — not "the customer", "the user", "müşteri",
"kullanıcı" — and do not narrate your own process ("now testing the
customer path", "screen checklist").

Keep the reply to what they need:

1. The URL, first, on its own line.
2. One or two sentences on what works, in their words, not test output.
3. At most one question, only if their answer changes what happens next
   (production? a domain?).

Leave out what you did not do, test transcripts, HTTP codes, internal ids
(`dep_…`, `prj_…`) and caveats about checks you skipped — unless they ask.
If something failed, say what failed and what you are doing about it.

## Auth first (mandatory)

Before any MCP mutate or project lookup in the turn:

1. Call MCP `whoami` (or `nesly whoami --json`).
2. If `signed_in` is false, read `reason` — it says which failure this is,
   and most of them are **not** a lost session. Do not end the turn on a
   bare `signed_in: false`.
   - `api-unreachable` (`retryable: true`): a session **is** stored and the
     API went quiet. Wait a few seconds and call whoami again, up to three
     times. Never ask for a login here; it changes nothing.
   - `no-credential`, `rejected`, `credential-unreadable`: run
     `nesly login` **yourself** — it opens the browser; the person only
     approves there. The command bounds its own wait (10 minutes, or the
     code's expiry), so do not wrap it in `timeout`, which is not even
     installed on a stock macOS. Say in one line that a browser window is
     waiting for them. Then call whoami again and carry on. With a
     password in the environment there is a no-browser path:
     `NESLY_PASSWORD=… nesly login --email <account>` (plus `--otp` when
     MFA asks). Only if the CLI is absent, or login returns an error,
     tell them to run `nesly login`, then `nesly mcp setup`, and restart
     the agent.
   - `account-mismatch`: the session is fine, the account selection is
     not. Follow the hint (`nesly account ls`, `nesly account use <name>`).
   Never invent Vercel/another cloud, paste tokens into chat, or keep
   calling `list_projects` / `deploy_project` while signed out.
3. Only then continue with already-live / preview / deploy.

If an MCP tool returns auth/401: **stop retrying MCP**. Use the CLI in the
same turn (`nesly … --json`). Do not hunt secrets or re-login paths for MCP.

App contract: listen on `$PORT`, bind `0.0.0.0`. Static sites need no server.

## One command, no local build

```
nesly deploy --yes --preview --json
```

- `--yes` creates and links a missing project in the same command; no
  separate `project create`.
- **Their code is theirs.** Never change a dependency version, a lockfile,
  `go.mod`/`go.sum` or the project's configuration to make a build pass: no
  `go get …@version`, no `npm audit fix`, no `go mod tidy` of your own. The
  person you deploy for owns the versions they pinned; the platform builds
  them as they are. A missing `go.sum` or lockfile is written in the
  platform's own build copy (`go mod tidy` runs there only when `go.sum` is
  absent) and no version is raised.
- Do **not** build or test the project locally first (`npm install`, a local
  build): a local toolchain may not even be allowed on the machine. Read the
  deploy's own output instead.
- If a build fails on a dependency, report the failing line and stop. Do not
  patch their tree to get past it — that is the customer's decision, and a
  silent version bump is a defect on this platform.
- Do not hide stderr (`2>/dev/null`): progress and the failing build line
  are printed there.
- The app reads `DATABASE_URL` (or another database variable) → before the
  deploy: `nesly addon add postgres --yes --json`. That database is the
  production one. A **preview** gets a database of its own instead, empty
  on the first preview and shared by the project's previews (Postgres,
  MySQL, MS SQL, MongoDB, Valkey, buckets) — nothing to set up; never point a preview at production
  data unless the user asks for exactly that.
- A project **with an addon** whose code carries its own database address
  (`postgres://localhost/…`, a `Server=localhost;…` string, a compose
  service name): before the deploy run `nesly fix dsn --yes` (MCP
  `fix_dsn` with `apply=true`; without `--yes` it only prints the diff).
  It changes exactly the connection string so the code reads the variable
  the addon sets — no address or secret is written — and nothing else in
  the code. **This is the only edit an agent ever makes to their code.**
  Say in one sentence that the code now reads the project's
  database. The deploy also makes the same change on its own copy when you
  skip this, so a deploy never fails on it; `nesly fix show` (MCP
  `get_deployment_fix`) prints what a deployment's build changed. No addon
  yet: it refuses and names `nesly addon add postgres` — add the addon first.

## Prebuilt static (optional)

When the site is **already built locally** (`dist/`, `out/`, `public/` after
the user’s own build) and they want publish without a platform rebuild:

```
nesly deploy --prebuilt --dir dist --preview --json
```

MCP: `deploy_project` with `prebuilt=true` and `dir` (do not pass `image`).
Do not invent a build step on the platform for that path.

## Already live? (mandatory before first deploy in the turn)

When the user says publish / deploy / host / “yayınla”:

1. Resolve the project: `.nesly/project.json`, or `nesly project info --json`
   / MCP `get_project`, or `list_projects`.
2. **If it already has a live URL**: say that URL in one sentence. The next
   action is a **new revision of the same project**, not a first-time publish.
3. **If none**: first deploy creates and links the project.

Skipping this check is a contract violation.

## Preview first, then ask for production (mandatory)

Unless the user **explicitly** said live / production / prod / canlı / “prod’a al”:

1. Deploy **preview** (does not touch the live site):

```
nesly deploy --preview --json
```

For a stable preview hostname (bookmark / CI): `--alias staging` →
`{project}--staging` on the platform domain. Same alias again replaces
the previous preview at that address.

MCP: `deploy_project` with `environment=preview` (and `create=true` on first;
optional `alias`).

2. Tell them the **preview URL** immediately (do not wait for them to ask
   for “test”).
3. Ask one question, in their language: **“Shall I put it live
   (production)?”** — in Turkish: “Canlıya (production) alayım mı?”
4. Only on yes: promote the preview’s deployment (same image, no rebuild):

```
nesly promote <preview-deployment-id> --json
```

MCP: `promote_deployment` with that `deployment_id`. Works for containers
(same image) and static sites (same published files; no re-upload).

If they already asked for production, use `--prod` (or plain `nesly deploy`)
and skip the question.

## After the first production deploy — domain

Ask once, after the project's **first** production release (or when they
bring up a domain) — not after every revision, and not when the project
already has a custom domain (`nesly domain ls --json` / MCP `list_domains`).
In their language: **“Do you have your own domain, or shall we get a new
one?”** — in Turkish: “Özel alan adın var mı, yoksa yenisini mi alalım?”

- **Already owns a name:** ask only for the FQDN →
  `nesly domain add example.com` / MCP `add_domain` → show DNS records →
  `nesly domain verify example.com --wait` / `verify_domain`. Prefer
  `nesly domain delegate example.com` / `delegate_domain` when they can
  change nameservers (platform manages DNS; then verify).
- **Wants to buy:** `nesly domain search <name>` / `search_domains` → show
  price → collect registrant (name, email, E.164 phone, street, city,
  postal code, country code; optional state/organization) →
  `nesly domain buy …` / MCP `purchase_domain` with `confirm=true` → then
  **add + delegate or verify** as above. HTTPS certificates are issued by
  the platform after DNS is correct (no agent polling).

Do not invent a domain. Do not send them to another registrar unless they
refuse platform purchase.

## When a deploy fails

```
nesly diagnose --json
```

MCP `diagnose_deployment` first — then fix and redeploy. Raw `nesly logs`
only if diagnose is not enough.

## Common follow-ups (details: https://nesly.net/llms.txt)

- How-to questions: `nesly docs <question>` (MCP `search_documentation`) — not `curl llms.txt | grep`

- Env: `nesly env set|ls|rm|pull` — MCP `set_env` / `list_env` / `delete_env` / `pull_env` / …
- Database: `nesly addon add postgres|mssql` — MCP `add_addon`
- Rollback: `nesly rollback <id>` — MCP `rollback`
- Git push-to-deploy: `nesly git connect` — MCP `connect_github_repository`
- Full surface: https://nesly.net/docs/agents/

Deleting a project: `nesly project delete <name> --yes` alone. It removes deployments, domains and addons, and parks
the database for `project restore`. Never run `addon rm` first — that
deletes the data for good and asks the phone a second time.

Destructive deletes (`addon rm`, project delete, `purchase_domain`) need
an explicit user yes first. Then the platform asks the user's phone: the
answer names a two-digit number. Show it to the user as it is ("07", not
"7") and end the turn — they choose it in the identity app, and the
platform carries the action out itself on that yes. Do not call the tool
again to "finish" it; afterwards a read (`get_project`, `list_addons`)
shows the result. A repeat while the phone has not answered returns the
same number and sends no second push.
