# beacon official

Beacon 官网（Nuxt 4 SSR/预渲染前端 + Go Fiber API/发布物后端）。

## 开发

```bash
task web:dev      # 安装依赖并启动前端 dev server
task web:build    # 构建前端
```

后端开发（根 `go.work` 未纳入 `website/server`，构建需 `GOWORK=off`）：

```bash
cd server && GOWORK=off go run . run -c configs/config.dev.toml
```

## 部署

根目录 `task amd64` / `task arm64` 构建镜像，详见根 `Taskfile.yml`。
