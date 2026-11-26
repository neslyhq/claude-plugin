# Backing services (addons)

Contents: the rule · catalog · SQL naming · previews · mail vs mailapi ·
sms · whatsapp · voice · model · removing.

## The rule

If the app needs a database, cache, mail, storage or login, ask the platform
first: `nesly addon add <type> --yes --json` (MCP `add_addon`). The variables
arrive on the next deploy. An app whose setup screen, README or code names
`DATABASE_URL`, `SMTP_URL`, `REDIS_URL` or `MONGODB_URI` means an addon
exists. Never search the machine for those values.

When credentials of the person's own provider are needed (sms, mailapi,
whatsapp, voice), ask the person for them and pass them on standard input
(`--credentials-file -`); never look for them, never put them in the chat
history, a file in the repository, or a command line.

## Catalog

| type | injects |
|---|---|
| `postgres`, `mysql`, `mssql`, `mongo` | `DATABASE_URL` / `MONGODB_URI`, spelled for the driver the build found |
| `valkey` | `REDIS_URL` / `VALKEY_URL`; keys start with `$REDIS_KEY_PREFIX` when the client can (required until the server is Valkey 9.1; from 9.1 any key name works in the project's own database) |
| `bucket` | `AWS_*` (S3-compatible) |
| `oidc` | `OIDC_*` (issuer URL, client id and secret) — login for the app |
| `mail` | `SMTP_URL`, `SMTP_HOST/PORT/USER/PASSWORD/FROM`, `MAIL_*` |
| `mailapi`, `sms`, `voice`, `whatsapp` | the person's own provider account, below |
| `sip` | a SIP account on the platform's border controller: `SIP_URI`, `SIP_USERNAME`, `SIP_PASSWORD`, `SIP_DOMAIN`, `SIP_PROXY`, `SIP_TRANSPORT=tls`; register over TLS, the app opens no UDP port |
| `model` | an OpenAI-compatible base URL, key and model name |

## SQL naming

A SQL database also arrives under the names the framework reads: Laravel gets
`DB_URL` and `DB_CONNECTION/DB_HOST/DB_PORT/DB_DATABASE/DB_USERNAME/
DB_PASSWORD`; .NET gets `ConnectionStrings__Default`. A value the project sets
itself wins.

## Production and previews

The addon's database is the production one. A **preview** gets a database of
its own, empty on the first preview and shared by the project's previews
(Postgres, MySQL, MS SQL, MongoDB, Valkey, buckets): nothing to set up. Never
point a preview at production data unless the person asks for exactly that.

With an addon present and the code carrying its own address
(`postgres://localhost/…`, `Server=localhost;…`, a compose service name):
`nesly fix dsn --yes` (MCP `fix_dsn`, `apply=true`; without `--yes` it only
prints the diff). It changes exactly the connection string and nothing else.
No addon yet: it refuses and names `nesly addon add postgres`.

## mail and mailapi

`mail` is the platform relay. `mailapi` is outgoing mail through the person's
own account at resend, sendgrid, postmark, mailgun or ses over the provider's
HTTPS API (provider SMTP ports are closed from apps): `--provider`,
`--sender <verified From>`, credentials on stdin. It injects the provider's SDK
names plus `MAIL_PROVIDER`, `MAIL_API_URL`, `MAIL_FROM`. A project holds one of
the two.

## sms

The person's own account at twilio, plivo, telnyx, netgsm, iletimerkezi or
verimor: `--provider`, `--sender`, credentials as `NAME=VALUE` lines on stdin.
It injects the provider's SDK names plus `SMS_PROVIDER`, `SMS_API_URL`,
`SMS_SENDER`, `SMS_USER/PASSWORD/API_KEY/API_SECRET`. The platform sends no
SMS. Run again with new credentials to replace them.

## whatsapp

The person's own WhatsApp Business Cloud API, credentials on stdin:
`WHATSAPP_TOKEN`, `WHATSAPP_PHONE_NUMBER_ID`, `WHATSAPP_BUSINESS_ACCOUNT_ID`,
`WHATSAPP_VERIFY_TOKEN` (generated if none), `META_APP_SECRET`,
`WHATSAPP_API_URL`. The webhook is the app's own HTTPS URL. The platform sends
nothing; App Review and templates are the person's.

## voice

The person's own account at twilio, verimor (Bulut Santral), plivo or netgsm:
`--provider`, `--caller-id`, credentials on stdin. It injects the provider's
SDK names plus `VOICE_PROVIDER`, `VOICE_API_URL`, `VOICE_CALLER_ID`,
`VOICE_WEBHOOK_SECRET`, `VOICE_API_KEY/API_SECRET`. HTTPS and webhooks only:
the webhook is the app's own HTTPS URL, media stays at the provider (a SIP
stack is the `sip` addon).

## model

`nesly addon add model <id>` takes an id from `nesly models list` (Pro plan);
the app gets `OPENAI_BASE_URL` and `OPENAI_API_KEY` too unless the project sets
its own.

## Removing

`nesly addon rm <type>` deletes the data for good and needs an explicit yes plus
the phone approval (exit 6). To delete a whole project use
`nesly project delete <name> --yes` alone: it parks the database for
`project restore`.
