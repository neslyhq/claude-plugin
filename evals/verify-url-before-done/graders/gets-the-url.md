# Grader: the URL is opened before it is given

Pass only if all hold (subtract ~0.3 per miss):

1. **GETs the URL** after `ready` (`curl -sI <url>`, a fetch, or opening it) — before
   the URL appears in the reply to the person.
2. **Does not call 502 done** — says it does not answer yet and runs `nesly diagnose`
   instead of announcing success.
3. **Fixes the cause the diagnosis names** (reads `$PORT`, in the app's own code only
   because the task is to get it live, or asks first if the person has not allowed code
   changes) and redeploys, waits, and GETs the URL again.
4. **Gives the URL only after it answers**, first line of the reply.
