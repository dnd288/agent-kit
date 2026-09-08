import { existsSync, readFileSync } from 'node:fs';
import path from 'node:path';

export const meta = {
  name: 'feature-wf',
  description: 'Multi-agent feature build from a change specification',
  phases: [
    { title: 'Preflight', detail: 'Read change folder, discover facts' },
    { title: 'Build', detail: 'Execute tasks in dependency order' },
    { title: 'Integration', detail: 'Wire up, run validation' },
    { title: 'Verify', detail: 'Adversarial verification against spec' },
  ],
};

const DEFAULT_CHANGE_ROOTS = ['changes', 'specs/changes', 'openspec/changes', 'docs/changes'];
const DEFAULT_TASK_FILES = ['tasks.md', 'tasks.txt', 'plan.md'];
const DEFAULT_VALIDATION_COMMAND = 'npm run validate --if-present';

export async function run({ args = [], step, agent, shell, log, workspace = process.cwd() }) {
  const changeId = readChangeId(args);
  const changeDir = findChangeDir(workspace, changeId);

  if (!changeDir) {
    throw new Error(
      `Change folder not found for "${changeId}". Looked in: ${DEFAULT_CHANGE_ROOTS.join(', ')}`,
    );
  }

  await step('Preflight', async () => {
    log(`Change: ${changeId}`);
    log(`Change folder: ${changeDir}`);
  });

  const manifest = readBuildManifest(changeDir);
  const taskFile = findTaskFile(changeDir);
  const tasks = normalizeTasks(manifest, taskFile);

  if (tasks.length === 0) {
    throw new Error(
      `No tasks found. Add build.json or one of: ${DEFAULT_TASK_FILES.join(', ')} in ${changeDir}`,
    );
  }

  await step('Build', async () => {
    for (const group of dependencyGroups(tasks)) {
      const independent = group.filter((task) => task.parallel !== false);
      const serial = group.filter((task) => task.parallel === false);

      if (independent.length > 0) {
        await Promise.all(independent.map((task) => runWorkerTask({ task, changeId, changeDir, agent, log })));
      }

      for (const task of serial) {
        await runWorkerTask({ task, changeId, changeDir, agent, log });
      }
    }
  });

  await step('Integration', async () => {
    const commands = manifest?.validation?.commands ?? manifest?.commands?.validation ?? [DEFAULT_VALIDATION_COMMAND];
    for (const command of commands) {
      log(`Running: ${command}`);
      await shell(command, { cwd: workspace });
    }
  });

  await step('Verify', async () => {
    await agent({
      name: 'feature-verifier',
      description: 'Verify feature against specification',
      prompt: [
        'Verify this completed change against its specification as a serial release gate.',
        `Change ID: ${changeId}`,
        `Change folder: ${changeDir}`,
        'Read the spec, tasks, implementation diff, and tests.',
        'Report pass/fail with concrete evidence. A claim without proof fails verification.',
      ].join('\n'),
    });
  });
}

function readChangeId(args) {
  const value = Array.isArray(args) ? args.find(Boolean) : String(args || '').trim();
  if (!value) {
    throw new Error('Usage: feature-wf <change-id>');
  }
  return value.replace(/^--change=/, '').trim();
}

function findChangeDir(workspace, changeId) {
  if (path.isAbsolute(changeId) && existsSync(changeId)) {
    return changeId;
  }

  for (const root of DEFAULT_CHANGE_ROOTS) {
    const candidate = path.join(workspace, root, changeId);
    if (existsSync(candidate)) {
      return candidate;
    }
  }

  const direct = path.join(workspace, changeId);
  return existsSync(direct) ? direct : null;
}

function readBuildManifest(changeDir) {
  const manifestPath = path.join(changeDir, 'build.json');
  if (!existsSync(manifestPath)) {
    return null;
  }

  try {
    return JSON.parse(readFileSync(manifestPath, 'utf8'));
  } catch (error) {
    throw new Error(`Invalid build manifest at ${manifestPath}: ${error.message}`);
  }
}

function findTaskFile(changeDir) {
  return DEFAULT_TASK_FILES.map((name) => path.join(changeDir, name)).find((candidate) => existsSync(candidate));
}

function normalizeTasks(manifest, taskFile) {
  if (Array.isArray(manifest?.tasks)) {
    return manifest.tasks.map((task, index) => ({
      id: String(task.id ?? index + 1),
      title: task.title ?? task.name ?? `Task ${index + 1}`,
      description: task.description ?? task.detail ?? '',
      dependsOn: task.dependsOn ?? task.depends_on ?? [],
      parallel: task.parallel,
      agent: task.agent,
      files: task.files ?? [],
    }));
  }

  if (!taskFile) {
    return [];
  }

  return readFileSync(taskFile, 'utf8')
    .split('\n')
    .map((line) => line.trim())
    .filter((line) => /^[-*]\s+\[[ xX]\]\s+/.test(line))
    .map((line, index) => ({
      id: String(index + 1),
      title: line.replace(/^[-*]\s+\[[ xX]\]\s+/, ''),
      description: '',
      dependsOn: [],
      parallel: true,
      files: [],
    }));
}

function dependencyGroups(tasks) {
  const remaining = new Map(tasks.map((task) => [task.id, task]));
  const completed = new Set();
  const groups = [];

  while (remaining.size > 0) {
    const ready = [...remaining.values()].filter((task) =>
      asArray(task.dependsOn).every((dependency) => completed.has(String(dependency))),
    );

    if (ready.length === 0) {
      throw new Error(`Circular or missing task dependencies: ${[...remaining.keys()].join(', ')}`);
    }

    groups.push(ready);
    for (const task of ready) {
      remaining.delete(task.id);
      completed.add(task.id);
    }
  }

  return groups;
}

async function runWorkerTask({ task, changeId, changeDir, agent, log }) {
  log(`Task ${task.id}: ${task.title}`);
  await agent({
    name: `feature-task-${task.id}`,
    description: task.title,
    subagent_type: task.agent,
    prompt: [
      'Build one task from a change specification. Do not work outside this task unless required to keep the build coherent.',
      `Change ID: ${changeId}`,
      `Change folder: ${changeDir}`,
      `Task ID: ${task.id}`,
      `Task title: ${task.title}`,
      task.description && `Task detail: ${task.description}`,
      asArray(task.files).length > 0 && `Expected files: ${asArray(task.files).join(', ')}`,
      'Read the relevant spec and repository conventions before editing.',
      'When done, report files changed, tests run, and any unresolved risks.',
    ]
      .filter(Boolean)
      .join('\n'),
  });
}

function asArray(value) {
  if (value == null) return [];
  return Array.isArray(value) ? value : [value];
}
