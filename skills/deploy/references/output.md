# Output contract: JSON fields, exit codes, reasons

Contents: JSON · deployment fields · exit codes · approval · reasons.

## JSON

Every CLI command accepts `--json` and then prints exactly one JSON object on
stdout; progress and logs go to stderr (so never hide stderr). MCP tools return
the same fields.

## Deployment fields

`deploy` and `wait` print: `deployment_id`, `project_id`, `status` (`queued`,
`building`, `deploying`, `ready`, `failed`, `superseded`, `cancelled`), `tier` (static or container),
`url`, `environment` (`preview` / `production`), `failure_reason` when it
failed, `build{framework,root,install,command,output,start}`, and `next` — a
one-line hint for the follow-up.

`ready_note{path,status}` is set when the deployment is `ready` but the
platform's own check got a 4xx at that path (below 500 counts as ready, as on
Vercel, whose default is no HTTP check at all). `say` then carries a sentence
about it after the URL line: relay it, and if the site should serve that path,
check the build output and routes. A health path in the config makes the check
strict (opt-in, like Vercel's Deployment Checks). `diagnose` returns it as a
finding with `class: ready-4xx` and `advice: true` (advice, not a failure).

`deploy --detach` prints `deployment_id`, `project`, `state`, `status` and
`next` (the `nesly wait …` command); `url` is there for production only (a
preview's address exists once it is ready). `--temporary` adds `claim_token`.

`diagnose` prints `status`, `missing_env`, and `findings[]` with `cause`,
`evidence` (the log lines) and a fix command; an `advice: true` finding is not a failure.

## Pull request with the platform's repair

After a repair that changed code (`build.fix_diff`) on a deployment built from a
git commit, `say` ends with an offer: "Want the same change as a pull request?
`nesly fix pr`". Relay it once, in one sentence. Run `nesly fix pr
[--deployment id]` (MCP `call_api` `GitService.OpenFixPullRequest`) only when the
person says yes; never accept on their behalf. It opens a new branch
`nesly/fix-<id>` from the deployment's commit, aimed at the production branch,
and never writes to the production branch. If the answer says the GitHub App
needs Contents and Pull requests write access, an owner of the GitHub account
approves it in the installation settings: tell the person that.

## Exit codes

| code | meaning | do |
|---|---|---|
| 0 | ok | continue |
| 1 | error: failed, cancelled, superseded, wait timeout | `diagnose`; on `WAIT_TIMEOUT` wait again |
| 2 | usage: wrong flag or argument; also a container refused for `--temporary` | read the message, correct the command |
| 3 | auth | whoami, then the setup skill |
| 4 | not found | wrong project or id: `list_projects`, `nesly ls` |
| 5 | conflict | name taken, rollout in progress: read the message |
| 6 | approval pending | below |

## Approval (exit 6)

Money, deletes, quota and customer data wait for the person's phone. The
command exits 6 and prints a two-digit approval number: show it as is ("07"),
end the turn, and on their yes the platform completes the action (or repeat the
command with `--approval <number>`). A repeat while the phone has not answered
returns the same number and sends no second push.

## Reasons

A refusal carries a reason code in `reason`; the CLI prints the next step and
MCP names the tool. Frequent ones:

| reason | next |
|---|---|
| `WAIT_TIMEOUT` | the deployment is still running: `wait` again |
| `START_MISSING_ENV` | `diagnose` lists the variables; addon, a self-generated secret (the platform makes it: deploy again) or the person's own secret (deploy skill) |
| `ADDON_MISSING` | `addon add <type>` |
| `PROJECT_NOT_FOUND` | `list_projects`; `deploy` with `create=true` |
| `NO_SOURCE` | `dir` set to the source directory, or `repository` set to a git URL |
| `APPROVAL_PENDING` | the person answers the approval, then repeat the call |
| `ROLLBACK_IS_PREVIEW` | `promote` with the preview's id |
| `PROMOTE_NEEDS_READY` | the deployment is not ready: `wait`, or roll back |
| `ROLLOUT_IN_PROGRESS` | `call_api` `DeploymentService.PromoteRollout` or `AbortRollout` |
| `DOMAIN_NOT_DELEGATED` | `attach_domain` with `delegate=true` |
| `SPEND_CAP_REACHED` | the person raises the cap (`call_api` `UsageService.SetUsageBudget`) |
| `PLAN_UPGRADE_REQUIRED` | no tool changes a plan: tell the account owner |
| `START_CRASHLOOP` | `get_logs`, then `diagnose` |
| `BUILD_*`, `START_*`, `ROLLOUT_*`, `IMAGE_PULL`, `PORT_MISMATCH`, … | a deploy failure: `diagnose` names the cause |

All codes with meaning and next command: https://nesly.net/docs/errors/.
