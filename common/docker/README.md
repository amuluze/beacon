# docker
docker manager

## 集成测试

测试会连接真实 Docker daemon，并使用唯一命名的临时容器、网络和镜像标签；所有创建的资源都会通过 `t.Cleanup` 回收，不执行全局 prune。

Docker SDK 读取 `DOCKER_HOST`，但不会自动读取 Docker CLI 的当前 context。在 macOS 使用 Colima 或 Docker Desktop 时，可按当前 context 运行：

```bash
DOCKER_HOST="$(docker context inspect --format '{{.Endpoints.docker.Host}}')" \
BEACON_DOCKER_TEST_IMAGE=alpine:3.20 \
go test -count=1 ./...
```

`BEACON_DOCKER_TEST_IMAGE` 可指向任意本机已有、包含 `sh` 和 `sleep` 的 Linux 镜像；未设置时默认使用 `alpine:3.20`，本机不存在则由测试拉取。
