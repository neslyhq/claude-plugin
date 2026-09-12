# Fixture: user asks to publish a local app on nesly with no mention of
# production or a custom domain yet.

User: Bu projeyi nesly'de yayınla.

Context for the agent under test:
- Working directory is an unlinked Node app (no `.nesly/` yet).
- User did not say live, prod, production, or canlı.
- They did not mention a domain.

Expected behaviour is graded by asks-before-prod.md.
