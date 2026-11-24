# Common Timex 模块

## 模块职责

`common-timex` 是独立的 Go module（`github.com/amuluze/amutool/timex`），提供时间工具集，供 `beacon` Server 与 `collia` Agent 复用，避免各模块重复实现时间戳、日期与定时器行为差异。承担以下职责：

- 当前时间获取：`CurrentTime()` 返回 HH:MM:SS 字符串；`CurrentDate(format)` 按传入格式返回日期字符串（空格式默认 `2006-01-02`）；`CurrentTimestamp()` 返回秒级 Unix 时间戳。
- 时间戳转换：`Int64ToTime(int64)` 把 Unix 秒级时间戳转换为 `time.Time`（`time.Unix(timestamp, 0)`，纯秒级，不处理毫秒）。
- Ticker 抽象：`Ticker` 接口 + `NewTicker(duration)` 工厂，封装 `time.Ticker` 与停止语义，简化跨模块定时任务写法。

`common-timex` **不**负责：

- cron 表达式解析：不实现 cron 调度。
- 时区业务逻辑：当前日期和时间采用进程本地时区，不统一转换为 UTC；业务侧按需转换。
- 业务定时任务调度：仅提供原语。

相关源码包：`common/timex`。

## 职责边界

| 职责 | 说明 | 是否负责 |
|------|------|----------|
| 当前时间/日期/时间戳 | `CurrentTime`、`CurrentDate(format)`、`CurrentTimestamp` | yes |
| 时间戳转换 | `Int64ToTime` | yes |
| Ticker 抽象 | `Ticker` 接口 + `NewTicker` 工厂 | yes |
| cron 解析 | 不实现 | no |
| 时区业务逻辑 | 不实现 | no |
| 业务定时调度 | 不实现 | no |
| HTTP 入口 | 纯 Go 库，不暴露路由 | no |
| 跨模块协调 | 仅作为 Go library 被引用 | conditional |
| 其他模块内部状态 | 不直接读写其他模块私有实现 | no |

## 模块事实

- 模块路径：`github.com/amuluze/amutool/timex`（独立 Go module，通过 monorepo 路径 `common/timex/` 维护）
- Go 版本：`1.21.4`
- Go 源码文件：2
- Go 测试文件：2
- 包名：`timex`
- 推断角色：时间工具与 Ticker 抽象，避免各模块重复实现时间处理

## 接口与符号信号

本节记录源码扫描得到的导出符号或前端入口信号；它们是候选接口线索，不等同于跨模块公开 API。跨模块调用仍以 facade、模块文档和 Domain Spec 约束为准。

- `interface Ticker (common/timex/ticker.go)`
- `func NewTicker (common/timex/ticker.go)`
- `func CurrentDate (common/timex/timex.go)`
- `func CurrentTime (common/timex/timex.go)`
- `func CurrentTimestamp (common/timex/timex.go)`
- `func Int64ToTime (common/timex/timex.go)`

## 依赖关系

- 未检测到第三方 Go require；依赖边界以 workspace 内模块和 Domain Spec 为准。

## 配置与入口信号

- 未检测到明显的配置结构或路由注册信号。

## 数据模型与状态

- 数据模型以源码结构、导出类型信号和测试断言为准；本节只记录静态发现结果。
- `interface Ticker (common/timex/ticker.go)`

## 使用注意事项

- 不跨越模块边界直接依赖内部路径或未导出符号。
- 修改跨模块公开 API、导出符号信号、配置或副作用边界时，同步更新本文件、根入口文档和相关 Domain Spec。
- 配置、凭据和外部副作用只记录语义边界，不记录真实敏感值。

## 验证命令

```bash
# common/timex 子模块无独立 Taskfile.yml，验证通过 Go 标准工具链完成
cd common/timex

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
