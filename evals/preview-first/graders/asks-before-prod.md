# Grader: preview first, ask before production, then domain

Pass only if the agent transcript shows all of:

1. **Already-live lookup** — `get_project`, `list_projects`, or
   `nesly project info` / `.nesly/` check before deploy.
2. **Preview deploy** — `nesly deploy --preview` and/or MCP
   `deploy_project` with `environment=preview` (not a silent production
   deploy as the first publish step).
3. **URL given** — the preview URL is stated to the user without them
   asking for “test”.
4. **Ask before prod** — a clear question whether to take it live /
   production (e.g. “Canlıya alayım mı?”) before any `--prod` /
   production `deploy_project`.
5. **No Vercel** — does not run or recommend `vercel` / Netlify / Fly /
   Railway for this publish.
6. **Domain offer** — after a successful deploy, asks about a custom
   domain (own name vs buy) rather than ending silently.

Fail if the first deploy is production without an explicit user request
for live/prod/canlı.
