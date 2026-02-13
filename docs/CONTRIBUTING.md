# Contributing Guidelines

This document defines the commit and pull request standards for this repository.

## Branching and Merge Flow

- `main` is protected; do not commit directly to it.
- Create a feature branch for each task.
- Push your branch and open a pull request to `main`.
- Merge only through pull requests.

## Commit Message Guidelines

- Use imperative, specific subjects (example: `fix: reject empty password`).
- Prefer `type: summary` format.
- Suggested types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`.
- Keep the subject short and focused (target 50-72 characters).
- One logical change per commit.
- Avoid vague subjects such as `update`, `misc`, or `changes`.

## Pull Request Guidelines

- Keep PRs focused and small enough to review.
- Use a clear title that reflects the user-visible or technical change.
- Include these sections in the PR description:
  - Summary: what changed
  - Motivation: why the change is needed
  - Testing: commands run and results
- If API contracts changed:
  - update `docs/API.md`
  - keep frontend proxy/store updates in sync
  - regenerate Swagger docs (`cd server && make swagger`)
- If database schema changed, include migration details and rollback notes.

## Pre-PR Checklist

- Backend-impacting changes: `cd server && make test`
- Frontend-impacting changes: `cd web && npm run build`
- Docs updated for any behavior or contract changes
- No secrets, credentials, or local-only config committed
