# Contributing

This template describes a generic contribution workflow for projects scaffolded with agent-kit. Replace placeholder commands and paths with the commands your repository actually supports.

## Before you start

Install the project prerequisites and confirm the local toolchain works.

- Runtime: `<runtime and version>`
- Package manager: `<package manager and version>`
- Database or service dependencies: `<local services>`
- Browser or mobile tooling: `<optional tooling>`
- Secrets: copy `<example env file>` to `<local env file>` and fill only local development values

Run the bootstrap command once:

```sh
<install command>
<setup command>
<validate command>
```

If setup fails, fix setup before writing feature code. A broken local loop makes every later result suspect.

## Workflow

Use the same shape for human and agent-driven changes.

1. **Branch** from the default branch.
2. **Specify** the behaviour when the change is user-visible, cross-workspace, or ambiguous.
3. **Build** the smallest vertical slice that proves the behaviour.
4. **Validate** with the right test modes and repository checks.
5. **Open a pull request** that links the issue, explains the change, and lists the evidence.

Do not mix unrelated layers in one pull request unless the specification says the integration is the change. Prefer small PRs that retire one uncertainty at a time.

## Day-to-day development commands

Replace these placeholders with real commands.

```sh
<install command>          # install dependencies
<dev command>              # run the local development server
<format command>           # format files
<lint command>             # lint files
<typecheck command>        # typecheck
<unit test command>        # run unit tests
<component test command>   # run component/browser tests
<e2e command>              # run end-to-end tests
<validate command>         # run the common ready check
```

Keep commands deterministic. If a command needs external services, document the service and the expected local URL.

## Frontend development

Frontend changes should keep rendering, data ownership, and copy ownership separate.

- Presentational components receive data through props.
- Application routes or screens own data loading and mutation wiring.
- User-facing copy comes from the catalogue when the project has i18n.
- Theme values come from tokens rather than raw colours or spacing values.
- Behaviour changes need tests at the layer that can observe them.

Run the component or browser test mode for changed interactions, responsive behaviour, accessibility states, or visual states.

## Backend development

Backend changes should keep transport and business logic separate.

- Route handlers validate input, call services, and return declared response shapes.
- Services hold business rules and should not depend on framework request/response objects.
- Database migrations must be backward compatible with the previous release unless the deployment plan says otherwise.
- External calls need typed failures, timeouts, and test doubles.
- Auth, authorization, webhooks, uploads, and secrets require a security pass.

Run unit tests for services, contract tests for API boundaries, and integration tests for persistence or external-system behaviour.

## Fixing a bug

A bug fix starts with a diagnosis.

1. Reproduce or observe the failure.
2. Identify the smallest code path that explains it.
3. Check what else relies on the behaviour you are about to change.
4. Fix the cause, not only the symptom.
5. Add or update the test that would have caught it.
6. Validate the affected slice.

If the cause is not obvious, write a short bugfix specification or issue note before changing code: observation, hypothesis, impact, fix, and verification.

## Testing

Use the cheapest mode that can prove the claim during development, then run the required ready checks before PR.

- **Static checks**: formatting, linting, type checking, repository guards.
- **Unit tests**: pure functions, services, reducers, validators, and utilities.
- **Component tests**: component states, events, accessibility contracts, and rendering branches.
- **Integration tests**: persistence, API contracts, queues, external clients, and framework wiring.
- **End-to-end tests**: full user journeys, cross-service flows, auth/session behaviour, and deployment-sensitive paths.
- **Security tests**: permissions, secret handling, dependency risk, injection, upload handling, and untrusted input.

Do not say a mode passed unless you ran it. If a mode is unavailable locally, record why and how CI or a reviewer should run it.

## Ready checklist

Before marking work ready:

- [ ] The change has an issue, task, or written reason.
- [ ] User-visible or cross-workspace behaviour has a specification.
- [ ] Non-obvious technical choices have an ADR.
- [ ] Tests cover the claim at the right layer.
- [ ] Static checks and guards pass.
- [ ] Documentation no longer describes shipped work as planned.
- [ ] Migrations and rollout steps are documented.
- [ ] Security-sensitive changes had a security review.
- [ ] The PR description lists commands run and results.

## Working with a coding agent

Keep project knowledge portable across tools.

- Put durable rules in `AGENTS.md`, workspace `AGENTS.md` files, docs, specs, and ADRs.
- Put task procedures in skills.
- Do not rely on a single vendor's memory feature for rules another agent must obey.
- Skills should route to owning documents instead of duplicating them.
- If an agent changes a rule in docs, it should check the skills that cite that rule in the same change.
- Ask agents to report checks they ran and checks they could not run.

## Commits

Use conventional commits:

```text
feat(scope): add new behaviour
fix(scope): correct broken behaviour
docs(scope): update documentation
test(scope): add or repair tests
refactor(scope): change structure without behaviour change
chore(scope): maintain tooling or dependencies
```

The scope should usually be the workspace, package, service, or layer. Keep commits coherent: one reason per commit.

## Pull requests

A pull request should include:

- Linked issue or task
- Summary of the change
- Test plan with commands and results
- Screenshots or recordings for UI changes when useful
- Migration, rollout, and rollback notes when applicable
- Labels that trigger optional CI suites when the repository uses label-driven checks
- Known follow-up work, with owners or issue links

Do not hide risk in the body text. Call out compatibility breaks, data migrations, new secrets, permission changes, and operational dependencies near the top.

## Decisions

Use ADRs for decisions that future maintainers would otherwise re-litigate.

File name:

```text
docs/adr/NNNN-short-title.md
```

Recommended sections:

- Title
- Status
- Context
- Decision
- Consequences
- What would reverse this

Every ADR needs an exit condition. If the decision is later reversed, update the old ADR's status and link the replacement.

## Secrets

Never commit credentials, tokens, private keys, production data, or customer data.

- Keep local secrets in ignored files.
- Provide example environment files with placeholder values only.
- Run secret scanners in CI and, where practical, pre-commit or pre-push hooks.
- If a secret is committed, treat it as compromised: rotate it, remove it from history according to project policy, and document the incident.
- Do not paste secrets into issues, PRs, logs, screenshots, or agent prompts.
