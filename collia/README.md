# collia

`collia` 是 Amprobe Server-Agent 架构中的 Agent 端，负责主机与 Docker 指标采集、本机执行能力，并通过本地 rpcx Unix Socket 提供给 `amprobe` Server 调用。

安装：

```bash
collia install
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
