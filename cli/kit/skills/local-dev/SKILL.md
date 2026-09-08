---
name: local-dev
description: "Run local development stacks end to end: select the smallest scenario for the goal, bootstrap a fresh clone, start required services, seed local data, debug startup failures, and verify that each running surface is healthy. Use this for local setup, install, dev server startup, container/native service choices, local data seeding, and environment gotchas; use a test-specific skill for test stacks."
---

# Local development stacks, end to end

This skill owns running a project locally: choosing the smallest stack that proves the goal, performing a one-pass bootstrap, and diagnosing startup sequences that look correct but fail.

It does not own test-runner topology, browser installation, CI-only services, or test data isolation. Route those to the project's testing instructions.

## Select a scenario

Ask what the goal needs, then start the smallest stack that proves it.

| Goal | Scenario |
|---|---|
| Smoke the app shell, auth, ordinary UI, and sessions | **minimal**: database + required local dependencies + API + client |
| Run with host services and no containers | **no-container**: native database, optional external dependencies disabled |
| Exercise a third-party integration without the real service | **mock-service**: app plus a standalone or in-process mock |
| Match deployment routing, path prefixes, proxy headers, or built images | **container-parity**: gateway/proxy + production-like containers |
| Use a real upstream service | **real-service override**: personal or approved environment only |
| Exercise async/background processing | **worker-included**: app plus the local queue/storage path and worker |
| Replace cloud-like services with local in-process plumbing | **in-process**: local filesystem/database queues instead of emulators |

Do not mix scenarios. If a cloud-emulator path and an in-process path both exist, choose one and remove the other path's environment variables before starting.

## Know the surfaces

List every process the stack may need and what starts it. Typical surfaces are:

| Surface | Started by | Needs |
|---|---|---|
| Database | container compose or native service | schema migrations and connection string |
| Object storage / queue emulator | local compose or emulator command | buckets, CORS, queues/topics primed |
| API server | workspace dev command | database, secrets, env loaded |
| Client server | workspace dev command | API URL, route prefix if any |
| Mock upstream service | mock command or in-process flag | signing keys or fixtures if applicable |
| Gateway/proxy | container stack | API/client services and path-prefix config |
| Background worker | worker command | queue/storage configuration and provider keys |

Document which surfaces are deliberately absent. A missing optional surface should fail at first use with a clear message, not be described as passing.

## First-time bootstrap

Prefer the project's one-command local bootstrap if it exists. By hand, the order is usually:

```sh
nvm use || true
corepack enable
pnpm install
cp .env.example .env          # fill required secrets; keep multiline secrets quoted
pnpm local:up                 # database and local dependencies, if containers are used
pnpm db:migrate
pnpm db:seed                  # only what the app needs to let an operator sign in
set -a; . ./.env; set +a      # export env for child processes
pnpm dev --env-mode=loose     # or the project-specific equivalent
```

Adjust package manager and commands to the repository. Preserve the ordering: dependencies, environment, infrastructure, migrations, seed/bootstrap, then long-running processes.

## Seed and bootstrap data

Separate bootstrap data from demo data.

- Bootstrap data creates the minimum account, tenant, or configuration needed to enter the app.
- Demo/reference data should be run through the same operator or maintenance path used outside local development.
- Re-running a seed or maintenance routine should be idempotent.
- Do not claim a seeded flow works until the login or entry point has been exercised.

If the project intentionally avoids a one-command demo seed, do not invent one. Use the operator path so local development exercises the same control plane as production.

## Scenario recipes

### No-container

Use native host services. Ensure the database exists, unset emulator-only values, run migrations and bootstrap seed, then start API and client processes.

State which legs are unavailable without containers: uploads, queues, workers, webhooks, or provider callbacks. Their clear first-use failures are expected, not broken boot.

### Minimal

Start only the surfaces required for ordinary app use. This is the default dev loop. It should prove sign-in, session persistence, navigation, and common read/write paths.

### Mock-service

Start the app and the mock service, or enable the in-process mock. Ensure both sides agree on base URLs, signing secrets, public keys, tenant IDs, and route prefixes. Use this for integration behavior when the real upstream is unavailable or inappropriate.

### Container-parity

Use this when routing, proxy headers, cookie attributes, path prefixes, image builds, or production server behavior matter. Remember that production-mode containers may enforce secure cookies, stricter origins, or built-asset behavior that dev servers do not.

### Real-service override

Use only in an approved personal or shared environment. Verify the real service and local app agree on base URLs, callback URLs, secrets, and signing keys. Keep mock credentials and real credentials distinct.

### Worker-included

Start the worker in a separate process after the API and dependencies are ready. Confirm queue/storage configuration and provider keys are visible to the worker. If a worker is absent, queued work may look like a hang.

### In-process path

Use local filesystem storage, database queues, or inline adapters when the project supports them. Remove emulator variables so the API and worker do not disagree about where work is written and read.

## Environment gotchas

- **Environment loading:** task runners may strip variables unless run in a loose or explicit env mode. Export `.env` before launching children when required.
- **Conflicted env files:** check `.env` and `.env.example` for `<<<<<<<`, `=======`, `>>>>>>>` before booting.
- **Quoted secrets:** multiline values, PEMs, JSON blobs, and secrets with spaces must stay quoted when sourced by a shell.
- **Port squatting:** list every port the scenario uses. Two stacks publishing the same port produce misleading failures.
- **Cookie security:** production-mode local containers may mark cookies `Secure` and refuse plain HTTP sessions.
- **Path prefixes:** configured prefixes affect cookies, redirects, proxy rewrites, callbacks, and CSRF routes. Never hard-code them.
- **Provider keys:** generation, payments, email, or other provider legs need real keys unless a mock path is selected.
- **Built package staleness:** workers and scripts often consume built packages or generated clients. Rebuild/regenerate after schema or package changes.
- **Migrations:** local drift can trigger prompts or checksum failures. Prefer deploy-style migrations for bootstrapping; reset only disposable local databases.

## Red flags

| Symptom | Likely cause |
|---|---|
| Server says a required secret is missing though `.env` has it | env not exported, env-mode stripped it, or wrong process loaded it |
| Sign-in returns success but no cookie appears | secure-cookie/plain-HTTP mismatch, cookie path mismatch, or proxy headers wrong |
| Migrations cannot connect | database not running, wrong connection string, or container/host networking mismatch |
| Demo screens are empty after seed | demo/reference maintenance routine did not run |
| Host works but gateway fails | path prefix, forwarded host/proto, cookie path, or proxy rewrite mismatch |
| Uploads or async jobs never complete | storage/queue path mismatch, missing emulator, or worker not running |
| Mock integration fails signature checks | unmatched HMAC secret, public key, base URL, or altered request body |
| API returns CSRF/forbidden on direct POST | CSRF bootstrap step missing; fetch first to mint token/cookie |
| A route fails only during restart | dev watcher is mid-rebuild; wait for health check and retry |
| Missing enum/export in worker or seed | generated client or built package is stale |

## Verification

- [ ] The scenario matches the goal and excludes unnecessary surfaces.
- [ ] Required environment variables are loaded by every process that needs them.
- [ ] Mutually exclusive local paths are not mixed.
- [ ] Ports are free or intentionally shared through one stack.
- [ ] Migrations and bootstrap seed ran successfully and are idempotent.
- [ ] Any demo/reference data path was run through the project's normal operator mechanism.
- [ ] Authentication/session flow was exercised through the configured host and prefix.
- [ ] Each running component answered its health check.
- [ ] Optional legs not run are named with the reason.
- [ ] No result is claimed for a provider, worker, upload, or integration leg that was not actually started.
