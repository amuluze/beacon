# Amprobe

![MIT License](https://img.shields.io/badge/License-MIT-green.svg)
![Go Reference](https://pkg.go.dev/badge/github.com/shirou/gopsutil/v3.svg)
[![GitHub stars](https://img.shields.io/github/stars/amuluze/amprobe)](https://github.com/amuluze/amprobe/stargazers)

English | [中文](./README.md)

## Introduction

Amprobe is a lightweight host and Docker container monitoring tool, it can easily help us to complete the following aspects of work:

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

Official website Address:[Website | Amprobe (official.amprobe.amuluze.com)](https://official.amprobe.amuluze.com/)

> **docker version required：>= 20.10.9**

## Technology Stack

Amprobe adopts a front-end and back-end separated technical architecture.

**Frontend Technology Stack:**

- `Vue3`
- `TypeScript`
- `Element+`
- `Vue-router`
- `Pinia`

**Backend Technology Stack:**

- `Golang`
- `Fiber`
- `Sqlite`

## Download & Installation

### Method 1: Download Binary Files Directly

Download the latest binary file for your system from the [Releases](https://github.com/amuluze/amprobe/releases) page.

### Method 2: Compile from Source

1. Clone the repository
```bash
git clone https://github.com/amuluze/amprobe.git
cd amprobe
```

2. Build the frontend
```bash
cd web
make install  # Install dependencies
make build    # Build frontend
cd ..
```

3. Build the backend
```bash
make assets   # Package frontend resources into backend
make wire     # Generate dependency injection code
make bin      # Compile binary file
```

## Usage

Install and use instructions can be found in the [Manual](https://official.amprobe.amuluze.com/document)

## License

Amprobe is available under the MIT license
