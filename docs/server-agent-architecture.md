# Server-Agent Architecture

Amprobe is the Server side. It owns the web UI, HTTP API, authentication, authorization, audit records, alarm configuration, and request orchestration.

Collia is the Agent side. It is installed on target machines and owns command execution, Docker operations, file operations, system operations, and metric collection.

## Interaction Model

The Server never executes host commands and never collects host metrics directly. Every host-side capability is exposed by Collia through RPC.

Browser requests target an Agent by passing one of:

- `X-Agent-ID: <agent-id>`
- `?agent_id=<agent-id>`

If no Agent ID is provided, Amprobe uses `Rpc.DefaultAgentID`.

## Amprobe Configuration

Single local Agent compatibility:

```toml
[Rpc]
Network = "unix"
Address = "/app/collia.sock"
DefaultAgentID = "default"
```

Multiple remote Agents:

```toml
[Rpc]
DefaultAgentID = "host-a"

[Rpc.TLS]
Enable = true
CertDir = "/app/certs/amprobe"
ServerName = "collia-host-a"

[[Rpc.Agents]]
ID = "host-a"
Network = "tcp"
Address = "10.0.0.11:18080"
TLS.ServerName = "collia-host-a"

[[Rpc.Agents]]
ID = "host-b"
Network = "tcp"
Address = "10.0.0.12:18080"
TLS.ServerName = "collia-host-b"
```

## Collia Configuration

Local Unix socket mode:

```yaml
rpc:
  network: unix
  address: /data/amprobe/resources/collia/socks/collia.sock
```

Remote TCP mode:

```yaml
rpc:
  network: tcp
  address: 0.0.0.0:18080
  tls:
    enable: true
    cert_dir: /etc/collia/certs
    client_names:
      - amprobe
```

## Transport Security

Remote TCP mode should use mTLS. Amprobe verifies each Collia Agent certificate
with the configured `TLS.ServerName`; Collia requires an Amprobe client
certificate signed by the same CA. `agent_id` is only an Amprobe routing key and
must not be treated as a security identity.

Each certificate directory must contain:

- `ca.pem`
- `tls.crt`
- `tls.key`

Suggested identity mapping:

| Component | Certificate identity | Configuration |
| --- | --- | --- |
| Amprobe Server | `amprobe` | `collia.rpc.tls.client_names: ["amprobe"]` |
| Collia host-a | `collia-host-a` | `amprobe.Rpc.Agents[].TLS.ServerName = "collia-host-a"` |
| Collia host-b | `collia-host-b` | `amprobe.Rpc.Agents[].TLS.ServerName = "collia-host-b"` |

## Compose Deployment

Amprobe is expected to run as the panel service through the repository root
`compose.yaml`:

```bash
docker compose up -d --build
```

The compose service exposes the panel on host port `1443` and mounts
`./deploy/downloads` into `/app/downloads` for Collia Agent distribution.

Expected Collia package layout:

```text
deploy/downloads/collia/
  linux/
    amd64/
      collia.install
    arm64/
      collia.install
  certs/
    1.tar.gz
    host-a.tar.gz
```

Each cert tarball should extract `ca.pem`, `tls.crt`, and `tls.key` directly
into `/etc/collia/certs`.

## Agent Bootstrap

Install a Collia Agent from Amprobe:

```bash
curl -kfsSL 'http://<amprobe-host>:1443/api/v1/host/install?node=1&os_type=linux' | sudo bash -s -- --token=<install-token>
```

`/api/v1/host/install` only returns a generic bootstrap script. The script then
uses `--token` as `X-Install-Token` to download:

- `collia.install`
- generated `/etc/collia/config.yml`
- the node certificate tarball when `AgentInstall.TLSEnable = true`

Amprobe configuration:

```toml
[AgentInstall]
Enable = true
Token = "change-me"
PublicBaseURL = ""
PackageDir = "/app/downloads/collia"
RPCPort = 18080
TLSEnable = true
CertDir = "/etc/collia/certs"
```

When `PublicBaseURL` is empty, Amprobe builds the script download base URL from
the incoming request host and scheme.

## Current Boundary

Amprobe currently proxies these capabilities to Collia:

- host metrics and host info
- Docker/container/image/network operations
- file create/delete/upload/download/search
- system DNS/time/timezone/reboot/shutdown operations
- container logs
- alarm checks that query Agent-side metrics

Interactive terminal sessions are intentionally disabled on the Server side. If web terminal support is required, it should be implemented as an Agent-side session protocol instead of direct Server SSH.
