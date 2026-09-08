# Infra PR — `pr-infra`

**When:** AWS, pipeline, or infrastructure work owned by your infrastructure package.

**Rule:** `layer:infra` work is `chore(infra):` and follows the chore path. Security-sensitive changes get a `security` review before merge.

## The rules

| | |
|---|---|
| Branch | `docs/<short-description>` or `feature/<issue>-short-description` if tied to an issue |
| Title | `chore(infra):` — no `—` suffix; all context goes in the body |
| Issue link | `Closes #N` if tied to an issue, otherwise none |
| Proof | `your validation command`; CloudFormation diff reviewed manually |
| Labels | `layer:infra` |

## Template body fields

`Closes #<N>` if applicable.
