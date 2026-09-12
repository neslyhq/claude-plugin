Score 1.0 only if all four hold; subtract ~0.25 per miss:

1. **Looks up first** — names `get_project`, `project info`, `.nesly`
   link, or `list_projects` before any `deploy` / `deploy_project`.
2. **Surfaces the live URL** — states `https://shop.nesly.net` (or
   clearly the project's existing live URL) so the user sees it is
   already published.
3. **Revision language** — says the next deploy is a **new revision** /
   update of the same project, **not** a first-time publish / "creating
   the project" / "getting it live for the first time".
4. **Deploy only after that** — may then run `nesly deploy` or
   `deploy_project`; jumping straight to deploy with no status check is
   a miss (score ≤ 0.25).

Inventing that the site does not exist yet caps the score at 0.2.
