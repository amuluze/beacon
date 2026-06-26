# Plan: Agent 远程更新推送与自更新

> Task Spec: [.specs/tasks/agent-remote-upgrade.md](../.specs/tasks/agent-remote-upgrade.md)
> Domain Spec: [.specs/domain/agent-lifecycle-update.md](../.specs/domain/agent-lifecycle-update.md)

## 目标

Server 通过反向 gRPC tunnel 主动推送更新指令给指定 Agent，Agent 安全完成自更新（下载→校验→原子替换→重启→重注册）。

## 技术上下文

- 反向 tunnel 已实现（`common/rpc/tunnel`），支持 `RegisterUnary[A,R]` 注册模式
- Agent 注册帧已携带版本信息（`agent-version-reporting` done）
- 安装包端点已有 `/api/v1/host/install/package`（支持 Install Token）
- Agent 进程路径通过 `-prefix` flag 配置（默认 `/data/amprobe`）
- `collia/service/rpc/dispatcher.go` 提供方法注册表

## 层级分解

### L2-后端（Server + Agent）

| 任务 | 文件 | 说明 |
|------|------|------|
| T1 | `common/rpc/schema/system.go` | 新增 `UpgradeAgentArgs{DownloadURL, SHA256, Version, InstallToken}`、`UpgradeAgentReply{Success, Version, Error}` |
| T2 | `collia/service/rpc/upgrade.go` | 实现 `UpgradeAgent` 方法：HTTP下载→SHA256校验→原子替换（两步rename）→重启 |
| T3 | `collia/service/rpc/dispatcher.go` | 注册 `"UpgradeAgent"` handler，调用 `RegisterUnary` |
| T4 | `amprobe/service/agent/api.go` | 新增 `Upgrade` API 端点，接收 `agent_id` + `version`，构建 UpgradeAgentArgs，通过 tunnel Call |
| T5 | `amprobe/service/router.go` | 注册 `POST /api/v1/agent/upgrade` 路由 |

### L3-前端

| 任务 | 文件 | 说明 |
|------|------|------|
| T6 | `amprobe/web/src/views/monitor/host/index.vue` | Agent 列表添加"升级"按钮，输入目标版本号 |
| T7 | `amprobe/web/src/api/host/index.ts` | 新增 `upgradeAgent` API 方法 |

### L4-验证

| 任务 | 说明 |
|------|------|
| T8 | `common/rpc/schema/system_test.go` — UpgradeAgentArgs/Reply 序列化测试 |
| T9 | `collia/service/rpc/upgrade_test.go` — SHA256 校验、并发锁、回退策略测试 |
| T10 | 全量构建验证 |

## 依赖图

```
T1 → T2 → T3
T1 → T4 → T5
T4 → T6 → T7
T2,T3 → T8,T9 → T10
```

## 关键实现细节

### T2: UpgradeAgent 方法（最核心）

```go
func (s *Service) UpgradeAgent(ctx context.Context, args rpcSchema.UpgradeAgentArgs, reply *rpcSchema.UpgradeAgentReply) error {
    // 1. 获取并发锁（防止并发更新）
    if !s.upgradeMu.TryLock() {
        reply.Error = "upgrade in progress"
        return nil
    }
    defer s.upgradeMu.Unlock()

    // 2. HTTP GET 下载新二进制到临时路径
    tmpPath := prefix + "/resources/collia/collia_new"
    if err := downloadFile(args.DownloadURL, args.InstallToken, tmpPath); err != nil {
        reply.Error = fmt.Sprintf("download failed: %v", err)
        return nil
    }

    // 3. SHA256 校验
    if err := verifySHA256(tmpPath, args.SHA256); err != nil {
        os.Remove(tmpPath)
        reply.Error = fmt.Sprintf("checksum mismatch: %v", err)
        return nil
    }

    // 4. 两步原子替换
    currentPath := prefix + "/resources/collia/collia"
    backupPath := prefix + "/resources/collia/collia.bak"

    os.Rename(currentPath, backupPath)    // 旧 → .bak
    os.Rename(tmpPath, currentPath)       // new → current（原子）

    // 5. 重启服务
    // Agent 进程自行退出，systemd/daemon 自动重启为新版本
    // 重启后以新版本重新注册到 Server

    reply.Success = true
    reply.Version = args.Version
    return nil
}
```

### T4: Server API 端点

```go
func (a *API) Upgrade(ctx *fiber.Ctx) error {
    agentID := contextx.FromAgentID(ctx.UserContext())
    version := ctx.Query("version")

    // 构建下载 URL
    baseURL := a.config.AgentInstall.PublicBaseURL
    downloadURL := fmt.Sprintf("%s/api/v1/host/install/package?version=%s", baseURL, version)

    args := rpcSchema.UpgradeAgentArgs{
        DownloadURL:  downloadURL,
        SHA256:      "", // TODO: 需从包元数据获取
        Version:     version,
        InstallToken: a.config.AgentInstall.Token,
    }

    var reply rpcSchema.UpgradeAgentReply
    if err := a.svc.tunnel.Call(ctx.UserContext(), agentID, "UpgradeAgent", args, &reply); err != nil {
        return fiberx.Failure(ctx, fiberx.ServiceError(err))
    }
    return fiberx.Success(ctx, reply)
}
```

## 风险与回退

| 风险 | 等级 | 缓解措施 |
|------|------|----------|
| 原子替换失败（权限） | 高 | `collia.bak` 保留，可手动回退 |
| 下载超时（大文件） | 中 | HTTP 下载设置 5 分钟超时 |
| 并发更新冲突 | 中 | `sync.Mutex` + `TryLock` 拒绝 |
| SHA256 计算时间 | 低 | 现代 CPU < 1s for 50MB |

## 提交策略

1. `feat(schema): add UpgradeAgentArgs/Reply for remote upgrade`
2. `feat(collia): implement UpgradeAgent with download→verify→atomic replace`
3. `feat(collia): register UpgradeAgent handler in dispatcher`
4. `feat(amprobe): add Agent upgrade API endpoint`
5. `feat(amprobe): register /api/v1/agent/upgrade route`
6. `test: add upgrade agent schema and handler tests`
