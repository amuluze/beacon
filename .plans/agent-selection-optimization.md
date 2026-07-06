# Agent Selection Optimization Plan

## Steps

1. Add shared frontend Agent state and API client.
2. Add a global navbar Agent selector.
3. Inject selected Agent into Axios requests and WebSocket URLs.
4. Reuse shared Agent state in monitoring pages and refresh container management on Agent change.
5. Add JSON tags to the Agent model.
6. Add Agent/time composite indexes to monitoring tables.
7. Adjust frontend embed/test baseline so Go tests can run before frontend build.
8. Run backend and frontend validation.
