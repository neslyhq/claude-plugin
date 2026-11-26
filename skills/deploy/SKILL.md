---
name: deploy
description: Deploys an application to nesly and verifies it is live — from a repository URL, a directory, prebuilt static files or a temporary no-account upload — and drives a failed deploy to green (wait, diagnose, repair, redeploy). Use when the user wants to deploy, publish, host, ship, go live, yayınla, canlıya al, or when a nesly deployment failed, is stuck building, or the app needs env vars, a database, mail or a domain.
---

# Deploying to nesly

Do not suggest or run `vercel`, `netlify`, `fly` or `railway` for this work:
nesly only. The CLI (`--json` on every command) and the MCP tools are the same
operations. App contract: listen on `$PORT`, bind `0.0.0.0`; static sites need
no server.

## Auth first (mandatory)

Before any project lookup or deploy: `whoami` (MCP or `nesly whoami --json`).
If `signed_in` is false, read `reason`; most are not a lost session. Steps and
the `reason` table: the **setup** skill (`references/auth.md`). Never invent
another cloud or paste tokens into the chat.

## Choose the source

- **A repository URL** ("deploy github.com/o/r"): `nesly deploy
  https://github.com/o/r --yes --json` (MCP `deploy` with `repository`). The
  platform clones it: **never `git clone` it, never upload it.** Private repo →
  `nesly git status`, then `nesly git connect o/r`.
- **This directory:** `nesly deploy --yes --json` (`--yes` creates and links a
  missing project). Do not build or test locally first.
- **Built static files** (`dist/`): `nesly deploy --prebuilt --dir dist --preview --json`.
- **Not signed in, static only:** `nesly deploy --temporary --dir <out> --json`
  (`url`, `claim_token`; gone in an hour unless `nesly claim <claim_token>`).

## Read the model first (mandatory)

With the CLI, before the first deploy in a turn run `nesly app show --json`: the
one answer to "what will run" (processes, port, health, release, env names by
source, generated secrets, memory and plan ceiling, needed addons, refused
files). An older CLI says `unknown command`: `nesly update` (Homebrew: `brew
upgrade nesly`) and retry; still unknown, `nesly build show --json`. Then:

1. Set only what is MISSING: `services.needed` → run its command; `env.hints`
   → only the person has the value (`nesly env request`). Never invent or copy
   a value (.env.example, compose, app.json); NEVER search this machine for one.
2. Compose is refused on purpose (`refused[]`): never translate it by hand.
   A worker the app's own docs name: `nesly build set --worker name="cmd"`.
3. Memory is the plan's ceiling unless something set less. On exit 137 the
   message names the limit, its source and the ceiling: remove a lower
   `nesly build set --memory`; tell the person when only a larger plan helps.
4. Never set or rotate generated secrets. 5. Destructive actions: last section.

## Already live? (mandatory before the first deploy in the turn)

Resolve the project (`.nesly/project.json`, `nesly project info --json`, MCP
`get_project` / `list_projects`). If it has a live URL, say it in one sentence:
the next action is a **new revision** of the same project, not a first
publish. Skipping this check is a contract violation.

## Preview first, then ask for production (mandatory)

Unless the person said live / production / canlı / "prod'a al", deploy a
preview (`--alias staging` for a stable hostname), tell them its URL at once,
and ask one question: "Shall I put it live (production)?" — in Turkish:
"Canlıya (production) alayım mı?" Only on yes: `nesly promote <preview-id>
--json` (MCP `promote`; same image, no rebuild). If they already asked for
production: `--prod` and skip the question.

## The loop (mandatory)

1. **Start:** `nesly deploy --preview --detach --yes --json` (MCP `deploy`:
   past `wait_seconds` it answers `state=building` with `next_tool`).
2. **Wait:** `nesly wait <id> --timeout 30m --json` (MCP `wait_deployment`).
   `reason: "WAIT_TIMEOUT"` stops only the waiting: wait again, never start a
   second deploy, never poll `inspect`.
3. **Ready → verify.** GET the URL (`curl -sI <url>`) and see it answer; only
   then give the URL and ask the production question. Not answering: step 4.
4. **Failed or not answering:** `nesly diagnose --json` (MCP `diagnose`):
   - **The platform already tried to repair it** in its own build copy: the
     log says `platform fix: …`, `nesly fix show` prints the diff. Do not
     fight it; if it could not, state the root cause it gives.
   - **Missing env — a backing service** (`env_findings` kind `addon`):
     `nesly addon add <type> --yes --json`. Catalog: `references/addons.md`.
   - **Missing env — a value the app makes up itself** (kind `self_generated`:
     `JWT_SECRET`, `APP_KEY`, …): the platform generates it itself on the
     deploy. You generate nothing and ask nobody: deploy again.
   - **Missing env — a third-party secret** (kind `third_party`: `STRIPE_*`,
     `OPENAI_API_KEY`): `nesly env request KEY…` (MCP `request_secret`) gives
     a console URL to type it into. Never ask them to paste it into the chat.
     Never search the machine for it.
   - **A code bug:** change their code only when they ask. Never a dependency
     version, a lockfile, `go.mod`/`go.sum` or a config file to get a build through.
5. **Redeploy → wait → verify the URL again** until it answers or the cause
   is one only the person can resolve.

A hardcoded database address with an addon present: `nesly fix dsn --yes` (MCP
`fix_dsn`, `apply=true`), the only edit an agent makes to their code. A scan
note never stops a deploy; `nesly scan` lists it; you only report the bump.

## After the first production deploy

Ask once: "Do you have your own domain, or shall we get a new one?" in Turkish:
"Özel alan adın var mı, yoksa yenisini mi alalım?" Own: `nesly domain add`.
New: `nesly domain search`, `nesly domain buy` (MCP `call_api`
`DomainService.PurchaseDomain`, with `confirm=true`). Steps: `references/domains.md`.

## Talking to the user (mandatory)

Answers carry `say`: the sentence for the person, URL first, in their
language; say it. `verified:false`: open the address yourself first. A repair's
`say` offers a pull request (`nesly fix pr`): relay it once, run it only on yes.
`rescue.repo_unchanged`: their repository was not touched. Speak to the person
("your site is live at …"), never "the user" or "müşteri": URL first on its own
line, one or two sentences, at most one question, no test output or ids.

## Destructive actions

Project delete, deployment delete, `addon rm`, backup delete, domain remove and
a domain purchase need the person's explicit yes first; then the platform asks
their phone: exit 6 with a two-digit number. Show it as is ("07"), end the turn,
do not call again; it completes on their yes. No flag (`--yes`, `--json`) skips
the phone. Delete a project with `nesly project delete <name> --yes` alone,
never `addon rm` first. JSON fields, exit codes: `references/output.md`.
