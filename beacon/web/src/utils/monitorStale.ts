/**
 * 监控数据过期态的展示格式化（DESIGN.md §7.5）。
 * 单调的秒数 → 中文时长；时间戳 → 稿面「最后采样」的两种精度。
 */
import { dayjs } from 'element-plus'

/** 把「已中断多少秒」格式化为中文时长：<60 秒 → 秒；<60 分钟 → 分钟；否则小时 + 分钟。 */
export function formatMonitorAge(seconds: number): string {
    if (!Number.isFinite(seconds) || seconds <= 0) {
        return '0 分钟'
    }
    const total = Math.floor(seconds)
    if (total < 60) {
        return `${total} 秒`
    }
    const minutes = Math.floor(total / 60)
    if (minutes < 60) {
        return `${minutes} 分钟`
    }
    const hours = Math.floor(minutes / 60)
    const rest = minutes % 60
    return rest > 0 ? `${hours} 小时 ${rest} 分钟` : `${hours} 小时`
}

/** 完整采样时刻：YYYY-MM-DD HH:mm:ss，零值回落 `-`。 */
export function formatSampleTime(epochSeconds: number): string {
    if (!epochSeconds) {
        return '-'
    }
    return dayjs.unix(epochSeconds).format('YYYY-MM-DD HH:mm:ss')
}

/** 指标值旁的短时刻：HH:mm:ss，零值回落 `-`。 */
export function formatSampleClock(epochSeconds: number): string {
    if (!epochSeconds) {
        return '-'
    }
    return dayjs.unix(epochSeconds).format('HH:mm:ss')
}
