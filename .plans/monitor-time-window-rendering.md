# Monitor Time Window Rendering Plan

## Steps

1. 用 Repository 测试证明主机和容器查询当前未限制 `end_time`。
2. 用 Vue 组件测试证明当前图表未把请求边界映射到 `xAxis.min/max`，并丢失 timestamp。
3. 将全部趋势查询收紧到闭区间 `[start_time, end_time]`。
4. 将主机和容器图表切换为时间轴，以 `[timestamp_ms, value]` 渲染，并固定本次刷新窗口。
5. 补齐容器近 5 分钟选项，验证空数据、竞态和 tooltip 兼容性。
6. 运行 Go 测试、Vitest、类型检查、定向 ESLint、生产构建和差异检查。
