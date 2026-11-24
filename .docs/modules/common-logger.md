# Common Logger 模块

## 模块职责

`common-logger` 是独立的 Go module（`github.com/amuluze/amutool/logger`），提供基于 `go.uber.org/zap` 的日志抽象，供 `beacon` Server 与 `collia` Agent 复用。承担以下职责：

- 顶层日志函数：`Debug`、`Debugf`、`DebugFields`、`Info`、`Warn`、`Error` 等全局入口，业务侧无需关心 zap 实例。
- Logger 配置：`Config` 结构 + `Option` 函数式选项，支持设置日志级别、格式（JSON/text）、输出位置、文件名后缀、文件最大保留时间、文件最大年龄。
- 文件轮转：依赖 `lestrrat-go/file-rotatelogs v2.4.0`，按时间间隔（`WithRotationTime`）轮转日志文件，避免单文件过大；不提供按大小轮转。
- 默认与 JSON 文件 logger：`Default()` 返回默认 zap 实例；`NewJsonFileLogger` 创建 JSON 格式 + 文件输出的 logger。
- 标准 slog 桥接：`NewJSONSlogHandler` 构造输出到 stdout + 按天切割文件 sink 的 slog JSON Handler（含 AddSource），供宿主进程 `slog.SetDefault` 接入，恢复旧版文件日志语义；rotator 初始化失败时降级为仅 stdout 并返回错误。
- 全局初始化：`SetLogLevel`、`SetLogFormat`、`SetLogOutput`、`SetLogFile` 等顶层 setter，由 `config` 解析后调用。

`common-logger` **不**负责：

- 业务事件追踪、APM、链路追踪：仅提供基础日志原语。
- 结构化字段语义约束：约定优于配置，业务侧自行定义字段名。
- 日志聚合与查询：仅本地输出，不内置远程上报。

相关源码包：`common/logger`。

## 职责边界

| 职责 | 说明 | 是否负责 |
|------|------|----------|
| zap 日志抽象 | `Debug/Info/Warn/Error` 全局入口 | yes |
| 函数式配置 | `Config` + `Option` | yes |
| 文件轮转 | 基于 `lestrrat-go/file-rotatelogs` | yes |
| 默认与 JSON logger 工厂 | `Default`、`NewJsonFileLogger` | yes |
| 业务事件追踪 | 不实现 | no |
| 远程日志聚合 | 不实现 | no |
| HTTP 入口 | 纯 Go 库，不暴露路由 | no |
| 跨模块协调 | 仅作为 Go library 被引用 | conditional |
| 其他模块内部状态 | 不直接读写其他模块私有实现 | no |

## 模块事实

- 模块路径：`github.com/amuluze/amutool/logger`（独立 Go module，通过 monorepo 路径 `common/logger/` 维护）
- Go 版本：`1.25.0`
- Go 源码文件：7
- Go 测试文件：1
- 包名：`logger`
- 推断角色：基于 zap 的日志抽象与文件轮转工具，为 beacon 与 collia 提供统一日志入口

## 接口与符号信号

本节记录源码扫描得到的导出符号或前端入口信号；它们是候选接口线索，不等同于跨模块公开 API。跨模块调用仍以 facade、模块文档和 Domain Spec 约束为准。

- `func SetLogFile (common/logger/config.go)`
- `func SetLogFileMaxAge (common/logger/config.go)`
- `func SetLogFileRotationTime (common/logger/config.go)`
- `func SetLogFileSuffix (common/logger/config.go)`
- `func SetLogFormat (common/logger/config.go)`
- `func SetLogLevel (common/logger/config.go)`
- `func SetLogOutput (common/logger/config.go)`
- `func SetName (common/logger/config.go)`
- `struct Config (common/logger/config.go)`
- `type Option (common/logger/config.go)`
- `func NewJsonFileLogger (common/logger/jsonfile.go)`
- `func NewJSONSlogHandler (common/logger/slog.go)`
- `func NewLogger (common/logger/default.go)`
- `func Debug (common/logger/default.go)`
- `func DebugFields (common/logger/default.go)`
- `func Debugf (common/logger/default.go)`
- `func Default (common/logger/default.go)`
- 另有 17 个导出符号未展开；请用 `rg`、codegraph 或 IDE 按需查看完整清单。

## 依赖关系

- `github.com/lestrrat-go/file-rotatelogs` `v2.4.0+incompatible`
- `go.uber.org/zap` `v1.28.0`
- 间接依赖 4 个，完整列表以模块 `go.mod` 为准。

## 配置与入口信号

- 配置：`Config` in `common/logger/config.go`
- 配置：`Option` in `common/logger/config.go`

## 数据模型与状态

- 数据模型以源码结构、导出类型信号和测试断言为准；本节只记录静态发现结果。
- `struct Logger (common/logger/logger.go)`
- 另有 2 个导出模型/接口符号未展开；完整清单以源码为准。

## 使用注意事项

- 不跨越模块边界直接依赖内部路径或未导出符号。
- 修改跨模块公开 API、导出符号信号、配置或副作用边界时，同步更新本文件、根入口文档和相关 Domain Spec。
- 配置、凭据和外部副作用只记录语义边界，不记录真实敏感值。

## 验证命令

```bash
# common/logger 子模块无独立 Taskfile.yml，验证通过 Go 标准工具链完成
cd common/logger

# 编译检查
go build ./...

# 单元测试（带 race 检测）
go test -race ./...

# 静态分析
go vet ./...
```

## Taskfile Commands

- `task beacon`
- `task beacon-web`
- `task collia`
- `task default`
