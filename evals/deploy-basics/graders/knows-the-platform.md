Score 1.0 only if the answer gets all three of these right; subtract
about a third for each miss:

1. **The deploy command is `nesly deploy`** (from the project
   directory; `--json` optional). Inventing other commands or a web
   upload flow is a miss.
2. **The application contract is named**: the app must listen on the
   `$PORT` environment variable and bind `0.0.0.0`. A FastAPI app
   needs a start command (e.g. uvicorn with `--host 0.0.0.0 --port
   $PORT`) — mentioning `nesly.toml [build] start = ...` or
   `nesly build set --start`, or previewing with `nesly build show`,
   counts.
3. **Failure path starts with `nesly diagnose`** (or the
   diagnose_deployment MCP tool) — NOT with reading raw logs; raw logs
   are the fallback when the diagnosis is not enough.

Wrong or fabricated flags and commands cap the score at 0.4.
