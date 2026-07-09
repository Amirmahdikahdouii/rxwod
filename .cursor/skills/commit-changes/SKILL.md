---
name: commit-changes
description: >-
  Commit uncommitted work using plan-split Conventional Commits: one commit per
  completed plan task or logical layer. Use when the user runs the commit-changes
  skill, asks to commit changes, or wants commits after finishing plan tasks.
disable-model-invocation: true
---

# Commit Changes

Commit current work with **plan-split** behavior: one commit per completed plan task or logical layer. Never batch unrelated work into a single commit unless the user explicitly asks for one commit.

## Safety Rules

- NEVER update git config
- NEVER use `--no-verify`, `--no-gpg-sign`, or other hook-skipping flags
- NEVER amend unless the user explicitly asks AND the amend conditions in the user rules are met
- NEVER push unless the user explicitly asks
- NEVER commit secrets (`.env`, credentials, private keys)
- Do NOT stage unrelated noise (vendor churn, `.env.example` unless that is the task, unrelated features)
- If there is nothing to commit, say so and stop

## Workflow

Run these in parallel first:

1. `git status`
2. `git diff` and `git diff --staged`
3. `git log -8 --oneline` (match message style)
4. If a plan is in context, note which todos/tasks are done and which files each task touched

Then:

### 1. Group changes

Split uncommitted files into commit groups:

| Situation | How to split |
|-----------|--------------|
| Multi-step plan just implemented | One commit per completed plan todo/task (or per logical layer if the plan groups that way) |
| Mixed unrelated edits | Separate commits by concern; leave unrelated files unstaged |
| User says "one commit" | Single commit for the relevant changes only |

Prefer layer/task boundaries such as:

- domain → application → infrastructure → delivery/HTTP → tests/docs
- one feature slice per commit when not following a numbered plan

### 2. Draft messages

Use Conventional Commits, focused on **why**:

```text
feat: ...
fix: ...
chore: ...
docs: ...
test: ...
```

Examples:

```text
feat: add WOD Archive and Delete HTTP handlers
feat: register WOD delete and archive routes with wod:delete auth
fix: map ErrCannotDeletePublished to 400
test: cover authz permission matrix edge cases
```

Keep the subject to one short line. Optional body only if needed for context.

### 3. Commit each group sequentially

For each group:

1. Stage **only** that group's files (`git add <paths>`)
2. Commit with a HEREDOC message:

```bash
git commit -m "$(cat <<'EOF'
feat: short description here

EOF
)"
```

3. If a pre-commit hook fails: fix the issue, then create a **new** commit (do not amend unless amend rules allow it)
4. After all groups: `git status` and report the new commit hashes/subjects

### 4. What to leave out

Unless the user asks to include them:

- `vendor/` / lockfile-only churn unrelated to the task
- `.env`, secrets, credentials
- Unrelated WIP in other directories

Warn if the user asks to commit secret-looking files.

## Output

After finishing, briefly list:

- Each commit hash + subject
- Any intentionally unstaged files
