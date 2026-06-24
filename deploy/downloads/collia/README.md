# Collia Agent Packages

Collia binaries are built into the Amprobe image during `docker buildx build`,
one per supported arch, and served from:

```text
/app/downloads/collia/
  amd64/collia
  arm64/collia
```

The Amprobe install script detects the target machine's architecture via
`uname -m` and downloads the matching binary, so no manual binary placement is
required.

TLS certificate tarballs are baked into the image from
`deploy/downloads/collia/certs/<node>.tar.gz` at build time:

```text
certs/
  1.tar.gz
  host-a.tar.gz
```

Each cert tarball should extract these files directly:

- `ca.pem`
- `tls.crt`
- `tls.key`
