# Beacon

![MIT License](https://img.shields.io/badge/License-MIT-green.svg)
![Go Reference](https://pkg.go.dev/badge/github.com/shirou/gopsutil/v3.svg)
[![GitHub stars](https://img.shields.io/github/stars/amuluze/amprobe)](https://github.com/amuluze/amprobe/stargazers)

English | [中文](./README.md)

## Introduction

Beacon is a lightweight host and Docker container monitoring tool. The project uses a Server-Agent architecture:

- `beacon`: the Server side, responsible for the Web UI, HTTP API, authentication, authorization, auditing, alarm configuration, and task orchestration.
- `collia`: the Agent side, responsible for host and Docker metric collection, local execution capabilities, and exposing them to the Server through a reverse gRPC tunnel (the Agent dials the Server, the Server calls back per `agent_id`).

It can help us complete the following tasks:

Architecture details: [System Architecture](./.docs/architecture.md)

### Container Manager

- View Docker version information
- Create, start, stop, restart, delete, and view container logs
- Import, export, and delete images, and clean up suspended images
- Create, delete, and view network status

### Host Monitor

- View the host name, startup time, release version, kernel version, and system type
- View host CPU, memory, disk IO, network IO monitoring

### Audit

- View user login, logout, and operation records

Official website Address:[Website | Beacon (beacon.amuluze.com)](https://beacon.amuluze.com/)

> **docker version required：>= 20.10.9**

## Technology Stack

Golang + Vue3

## License

Beacon is available under the MIT license

## Thanks

> [GoLand](https://www.jetbrains.com/go/?from=gopay) A Go IDE with extended support for JavaScript, TypeScript, and Databases。

Special thanks to [JetBrains](https://www.jetbrains.com/?from=gopay) for the Goland license for the open source project
