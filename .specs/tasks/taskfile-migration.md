# Taskfile Migration

## Goal

Replace Makefile-based development task management with Taskfile-based task definitions across the workspace.

## Requirements

- Existing Makefile targets must have equivalent Taskfile tasks.
- Developers must be able to run module-local tasks with `task <target>` from the module directory.
- The root workspace must expose included module task namespaces.
- Generated outputs, dependency directories and local environment files remain excluded from task management changes.
- AI entry documents and implementation docs must reference Taskfile commands instead of Makefile commands.

## Non-goals

- No change to application runtime behavior.
- No change to Docker image tags or build outputs beyond task runner migration.
