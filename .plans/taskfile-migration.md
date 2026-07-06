# Taskfile Migration Plan

## Steps

1. Inventory existing Makefile targets and command behavior.
2. Add root and module-level `Taskfile.yml` files with equivalent targets.
3. Remove superseded Makefiles.
4. Update AI entry docs and implementation docs from `make` commands to `task` commands.
5. Verify task discovery with `task --list` and representative module `task --list` commands.
