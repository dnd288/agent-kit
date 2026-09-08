---
description: Run review and security audit gates before shipping a change
argument-hint: [change-id-or-pr]
---

You are running the ship gate for a change.

Target: `$ARGUMENTS`

Run these gates in parallel:

1. `reviewer` agent
   - Review the target for correctness, architecture, security/privacy, tests, and maintainability.
   - If a change/spec folder exists, review against it before the diff.
   - Require concrete failure scenarios and evidence.

2. `security-auditor` agent
   - Audit authentication, authorization, sessions, cookies, uploads, environment variables, dependencies, and CI/CD workflows.
   - Report confirmed vulnerabilities separately from hardening recommendations.

After both gates return:

- Summarize blockers first.
- Summarize non-blocking follow-ups second.
- List validation commands already run and their result.
- State one of:
  - `SHIP: approved` when both gates pass and required validation has evidence.
  - `SHIP: blocked` when any correctness, security, spec, or missing-proof issue remains.

Do not approve a change on intent alone. A claim without a test, guard, manual verification note, or other concrete evidence is not ready to ship.