# Guard pattern

Guards are repository-specific static analysis checks that run beside tests. They protect architectural invariants that a general-purpose linter cannot know: package boundaries, generated-file rules, API contracts, migration policy, security-sensitive patterns, or design-system constraints.

A guard is not a formatter and not a style preference. It is a small executable rule that answers: "Does this repository still obey an invariant we have chosen to protect?"

## Guards versus linters

Linters usually enforce language and style conventions: unused variables, unreachable code, import order, hooks rules, formatting, and common bug patterns.

Guards enforce project policy:

- A frontend package must not import a backend package.
- Route handlers must declare response schemas.
- Business services must not accept framework request objects.
- UI components must use theme tokens rather than raw colours.
- Generated files must not be edited by hand.
- Security-sensitive APIs must pass through approved wrappers.

Keep guards narrow. A guard should produce a finding a maintainer can understand and fix without debating the rule from scratch.

## Check interface

A check is a small unit of policy.

```ts
export type Severity = "error" | "warning";

export interface Check {
  id: string;
  group: "invariants" | "security" | string;
  title: string;
  severity: Severity;
  run(files: SourceFile[], context: GuardContext): Violation[] | Promise<Violation[]>;
}
```

Recommended fields:

- `id`: stable machine-readable name, for example `ui/no-raw-hex-colours`.
- `group`: broad category used for reporting and selective runs.
- `title`: human-readable title.
- `severity`: `error` blocks, `warning` advises.
- `run(files, context)`: returns violations. It should not write files.

Use stable IDs because CI annotations, ignore files, documentation, and dashboards may refer to them.

## Violation structure

Every violation must include a fix. A guard without a fix trains people to resent the guard.

```ts
export interface Violation {
  check: string;
  path: string;
  line: number;
  message: string;
  fix: string;
}
```

Fields:

- `check`: the check ID.
- `path`: normalized repository-relative path.
- `line`: 1-based line number. Use `1` for file-level findings.
- `message`: what is wrong, in concrete terms.
- `fix`: what to do next. Mandatory.

Good message:

```text
Raw colour "#ffcc00" appears in a component.
```

Good fix:

```text
Move the value to the theme token file, regenerate tokens if required, and reference the token here.
```

Weak fix:

```text
Fix this.
```

## Runner pattern

The runner coordinates checks. Keep it boring.

1. **Collect files** from the repository root.
2. **Build context** once: package graph, config files, workspace metadata, generated-file manifests, ignore rules, or parsed schemas.
3. **Select checks** by group, path, changed files, severity, or CLI flags.
4. **Run checks** with the file list and context.
5. **Report** violations as terminal output, CI annotations, JSON, or SARIF.
6. **Exit non-zero** when any `error` violation exists.

Pseudocode:

```ts
async function main() {
  const root = findRepoRoot(process.cwd());
  const files = await walkFiles(root, defaultWalkOptions);
  const context = await buildContext(root, files);
  const checks = selectChecks(allChecks, parseArgs(process.argv));

  const violations = (
    await Promise.all(checks.map((check) => check.run(files, context)))
  ).flat();

  report(violations);

  if (violations.some((violation) => severityFor(violation.check) === "error")) {
    process.exitCode = 1;
  }
}
```

Avoid checks that shell out repeatedly per file. Build shared context once, then run fast pure checks over it.

## File walker

The file walker decides what source text exists.

Requirements:

- Walk recursively from the repository root.
- Ignore dependency, build, cache, coverage, and VCS directories.
- Scan only extensions relevant to checks.
- Normalize paths to POSIX-style repository-relative paths.
- Keep file contents and line starts available for checks.
- Follow symlink policy deliberately: either skip symlinks or only follow known safe ones.

Typical ignore directories:

```ts
const ignoredDirs = new Set([
  ".git",
  "node_modules",
  "dist",
  "build",
  "coverage",
  ".cache",
  ".next",
  ".turbo",
]);
```

Typical source extensions:

```ts
const sourceExtensions = new Set([
  ".ts",
  ".tsx",
  ".js",
  ".jsx",
  ".mjs",
  ".cjs",
  ".json",
  ".md",
  ".css",
  ".scss",
  ".html",
  ".yml",
  ".yaml",
]);
```

Normalize paths:

```ts
function normalizePath(path: string): string {
  return path.split("\\").join("/");
}
```

Line lookup helper:

```ts
function lineOf(content: string, index: number): number {
  let line = 1;
  for (let i = 0; i < index; i += 1) {
    if (content.charCodeAt(i) === 10) line += 1;
  }
  return line;
}
```

For large repositories, precompute line starts instead of scanning from the beginning for every match.

## How to write a check

This TypeScript example flags raw hex colours in component files. The implementation is intentionally simple; production checks may parse ASTs when syntax matters.

```ts
type SourceFile = {
  path: string;
  content: string;
};

type GuardContext = {
  root: string;
};

type Severity = "error" | "warning";

type Violation = {
  check: string;
  path: string;
  line: number;
  message: string;
  fix: string;
};

type Check = {
  id: string;
  group: "invariants" | "security" | string;
  title: string;
  severity: Severity;
  run(files: SourceFile[], context: GuardContext): Violation[];
};

function lineOf(content: string, index: number): number {
  return content.slice(0, index).split("\n").length;
}

const hexColourPattern = /#[0-9a-fA-F]{3}(?:[0-9a-fA-F]{3})?\b/g;

export const noRawHexColours: Check = {
  id: "ui/no-raw-hex-colours",
  group: "invariants",
  title: "Components use theme tokens instead of raw hex colours",
  severity: "error",
  run(files) {
    const componentFiles = files.filter(
      (file) =>
        file.path.startsWith("src/components/") &&
        (file.path.endsWith(".tsx") || file.path.endsWith(".css")),
    );

    const violations: Violation[] = [];

    for (const file of componentFiles) {
      for (const match of file.content.matchAll(hexColourPattern)) {
        const colour = match[0];
        const index = match.index ?? 0;

        violations.push({
          check: "ui/no-raw-hex-colours",
          path: file.path,
          line: lineOf(file.content, index),
          message: `Raw colour "${colour}" appears in a component.`,
          fix: "Move the value to the theme token file and reference the token here.",
        });
      }
    }

    return violations;
  },
};
```

When regex is not enough, parse the file. Use an AST for import rules, call signatures, JSX props, route definitions, or anything where comments and strings would create false positives.

## Integration

### CI pipeline

Run guards in CI on every pull request. The guard command should:

- Print concise terminal output.
- Annotate files and lines when the CI system supports it.
- Upload JSON or SARIF if the project tracks trends.
- Fail the job on `error` findings.
- Allow warnings without blocking merge, unless policy says otherwise.

### Git hooks

Run a fast guard subset in pre-commit or pre-push hooks.

- Pre-commit: changed-file checks that complete quickly.
- Pre-push: static checks that do not need services.
- Never make local hooks the only enforcement. CI is the source of truth.

### IDE

Optional IDE integration can run guards on save or expose JSON output through diagnostics. Keep IDE integration advisory. Developers must be able to reproduce every result with the CLI.

## Severity model

Use two severities by default.

### Error

Blocks the build. Use for crisp invariants with low false-positive risk.

Examples:

- Forbidden dependency direction.
- Missing required response schema.
- Secret-looking file committed under source control.
- Generated file edited by hand.

### Warning

Advises without blocking. Use for emerging rules, migration periods, or checks that are useful but not precise enough to block.

Examples:

- A file is approaching complexity limits.
- A deprecated API is still used during migration.
- A public module lacks documentation.

A warning that everyone ignores should either become precise enough to block or be removed.

## Trust principle

If a guard fails, fix the code rather than the guard. A guard that produces false positives is a higher-priority bug than one that misses a case.

This principle keeps the system trusted. A false negative is a missed opportunity. A false positive teaches people the guard is noise. Once a guard is noise, real findings are ignored too.

When a guard is wrong:

1. Confirm the intended invariant.
2. Add a regression fixture for the false positive.
3. Fix the guard.
4. Keep the rule if it still protects real cases.
5. Remove or demote the rule if it cannot be made trustworthy.

## Groups

Start with two groups.

### Invariants

Architecture, layering, generated files, design-system constraints, API contracts, migration rules, and repository layout.

### Security

Secrets, unsafe dependencies, untrusted input, authorization bypasses, insecure transport, dangerous APIs, and missing hardening around uploads or webhooks.

Add groups only when reporting or selective execution needs them. Possible future groups: `performance`, `accessibility`, `docs`, `migrations`, `dependencies`, `observability`.

## Guard design checklist

Before adding a guard, answer:

- What invariant does this protect?
- Where is the rule documented?
- Can the check identify violations with low false-positive risk?
- What exact fix will it print?
- Should it block (`error`) or advise (`warning`)?
- Which files should it scan?
- Which cases are explicitly allowed?
- What test fixtures prove allowed and disallowed cases?

A guard earns its place by making the correct path cheaper than the wrong one.
