# Collia Agent Packages

Place Collia installer packages here for Amprobe to distribute:

```text
linux/
  amd64/
    collia.install
  arm64/
    collia.install
certs/
  1.tar.gz
  host-a.tar.gz
```

Each cert tarball should extract these files directly:

- `ca.pem`
- `tls.crt`
- `tls.key`
