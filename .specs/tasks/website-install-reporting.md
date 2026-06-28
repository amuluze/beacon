# Website Install Reporting

## Goal

Amprobe containers started from the website-provided Compose installation flow must report installation startup metadata back to the website service for aggregate installation statistics.

## Requirements

- The website exposes an unauthenticated install reporting endpoint for Amprobe container startup reports.
- The website stores each startup report with installation identifier, image/version signals, exposed ports, public base URL, install directory, container hostname, client IP, user agent, and timestamp metadata.
- Amprobe creates and reuses a stable installation identifier in its mounted data directory.
- Install reporting defaults are owned by Amprobe internal configuration, not user-facing `.env` or Compose environment variables.
- Amprobe startup reporting must be best-effort and must not block container startup when the website endpoint is unavailable.

## Non-goals

- No user-facing analytics dashboard is required in this task.
- No authentication or license enforcement is added to the install reporting endpoint.
