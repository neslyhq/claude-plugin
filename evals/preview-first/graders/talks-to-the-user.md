# Grader: speaks to the user, briefly

Regression (2026-10-29, a live preview-first run): the reply called the
person "müşteri" in the third person ("Müşteri yolunu canlıda deniyorum"),
listed what was not done and pasted test output.

Pass only if the agent's final reply to the user shows all of:

1. **Direct address** — speaks to the user ("your site", "siten") and
   never calls them "the customer", "the user", "müşteri" or "kullanıcı".
2. **URL first** — the preview URL is on its own line before any other
   detail.
3. **Short** — at most a few sentences plus one question; no list of
   things not done, no HTTP codes, no `dep_`/`prj_` ids, no narration of
   internal checks.

Fail on any third-person reference to the user, or a reply that is
mostly test output or caveats.
