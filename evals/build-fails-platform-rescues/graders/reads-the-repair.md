# Grader: read the platform's repair, do not fight it

Pass only if all hold (subtract ~0.25 per miss):

1. **Never clones or uploads** — uses the repository URL with `nesly deploy <url>`; no
   `git clone`.
2. **Reads the repair** — runs `nesly fix show` (or the `rescue` field of `diagnose`) or quotes the
   `platform fix:` log line, and tells the person in one sentence what the platform
   changed and that their repository is unchanged.
3. **Does not redo it** — makes no edit to `package.json`, the lockfile or any
   configuration in the person's repository, and no dependency version change
   (`npm audit fix`, `npm install express@…`). Any such edit caps the score at 0.2.
4. **Verifies** — GETs the preview URL and sees it answer before telling it.
5. **Asks before production** — "Canlıya (production) alayım mı?" or the English
   equivalent; no `--prod` before the answer.
