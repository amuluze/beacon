# collia

`collia` 是 Beacon Server-Agent 架构中的 Agent 端，负责主机与 Docker 指标采集、本机执行能力，并通过反向 `gRPC tunnel` 主动连接 `beacon` Server，按 `agent_id` 让 Server 反向调用本机资源。

安装（默认读取 `/etc/collia/config.yml`，可用 `-c` 指定其他路径）：

```bash
collia install
# collia install -c /path/to/config.yml
```

启动

```bash
collia start
```

停止

```bash
collia stop
```

移除

```bash
collia remove
```

查看状态

```bash
collia status
```

前台运行（调试）/ 仅执行环境预检 / 查看版本

```bash
collia run
collia probe
collia version
```

说明：

- 默认文件位置：二进制 `/usr/sbin/collia`、配置 `/etc/collia/config.yml`、日志 `/data/beacon/logs/collia/collia.log`、本地数据 `/data/beacon/collia/storage/`。
- 本机 Agent 由 Beacon Server 经反向隧道触发升级；是否自动升级由 Server 控制台开关决定（默认开启），Agent 侧校验 SHA-256 后原子替换二进制。
- `remove`（别名 `uninstall`）先停止服务再清理二进制、配置、本地数据、日志与系统服务注册；未被 systemd 跟踪的残留 daemon 进程也会被终止。
