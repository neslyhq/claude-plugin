# Custom domains

Ask once, after the project's first production release (or when the person
brings up a domain), not after every revision, and not when
`nesly domain ls --json` (MCP `get_project`) already shows one. In their
language: "Do you have your own domain, or shall we get a new one?" — "Özel
alan adın var mı, yoksa yenisini mi alalım?"

## Own domain

Ask only for the FQDN, then:

1. `nesly domain add example.com` (MCP `attach_domain`) prints the DNS records.
2. If the person can change nameservers, prefer `nesly domain delegate
   example.com` (MCP `attach_domain` with `delegate=true`): the platform then manages the DNS.
3. `nesly domain verify example.com --wait` (MCP `attach_domain` again).

HTTPS certificates are issued by the platform once DNS is right; do not poll.

## Buy a new one

1. `nesly domain search <name>` (MCP `call_api` `DomainService.SearchDomains`); show the price.
2. Collect the registrant: name, email, E.164 phone, street, city, postal code,
   country code (optional state, organization).
3. `nesly domain buy …` (MCP `call_api` `DomainService.PurchaseDomain`, `confirm=true`). Buying costs
   money: explicit yes first, then the phone approval (exit 6).
4. Then add + delegate or verify, as above.

## Pin a domain to a deployment

`nesly domain pin <fqdn> <deployment-id>` and `nesly domain unpin <fqdn>`.

## Never

Invent a domain, or send the person to another registrar unless they refuse
the platform purchase.
