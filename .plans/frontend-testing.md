# Frontend Testing Plan

## Steps

1. **Tooling setup** (P1)
   - Add `vitest`, `jsdom`, `@vue/test-utils` to `amprobe/web/package.json` devDependencies.
   - Create `vitest.config.ts` with Vue plugin, jsdom environment, and `@` alias.
   - Add `test` and `test:watch` scripts to `package.json`.

2. **Agent store tests** (P1)
   - Create `src/store/modules/agent.test.ts` covering `setAgents`, `setCurrentAgent`, and empty-list behavior.

3. **Axios interceptor tests** (P2)
   - Create `src/api/index.test.ts` mocking Axios and verifying `X-Agent-ID` and `Authorization` headers.

4. **CI integration** (P2)
   - Add `task amprobe-web:test` target to `amprobe/web/Taskfile.yml`.
   - Ensure `pnpm test` runs before `pnpm build` in release pipelines.

5. **Validation**
   - `pnpm test` passes in `amprobe/web`.
   - `pnpm build` still succeeds after adding test files.
