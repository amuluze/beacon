# Website Install Reporting Plan

## Steps

1. Add website schema/model/repository/service/API support for install startup reports.
2. Register `POST /api/v1/install/report` in the website router.
3. Keep reporting internals out of generated `.env` and Compose environment variables.
4. Add Amprobe startup reporter using standard library HTTP, called asynchronously from service initialization, with config-owned defaults and a data-volume install id file.
5. Verify website server tests, Amprobe tests/build, and Nuxt build.
