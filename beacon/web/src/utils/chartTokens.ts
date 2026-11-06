/**
 * @Description: 图表序列色与图例文字色统一从设计 token 解析，避免写死色值。
 * 与 DESIGN.md §5.1（chart-1..chart-6 图表序列色）保持单一来源。
 */

export const CHART_SERIES_TOKENS = [
    '--am-chart-1',
    '--am-chart-2',
    '--am-chart-3',
    '--am-chart-4',
    '--am-chart-5',
    '--am-chart-6',
]

const CHART_LEGEND_TEXT_TOKEN = '--am-foreground-secondary'

/** 读取单个设计 token 的当前主题取值，浏览器外返回空串。 */
export function resolveChartTokenValue(token: string): string {
    if (typeof window === 'undefined')
        return ''
    return getComputedStyle(document.documentElement).getPropertyValue(token).trim()
}

/** 图表序列色板，顺序与 chart-1..chart-6 一致。 */
export function resolveChartPalette(): string[] {
    if (typeof window === 'undefined')
        return CHART_SERIES_TOKENS.map(() => '')
    const styles = getComputedStyle(document.documentElement)
    return CHART_SERIES_TOKENS.map(token => styles.getPropertyValue(token).trim())
}

/** 图例文字色，取次级前景色，暗色主题下同样可读。 */
export function resolveChartLegendTextColor(): string {
    return resolveChartTokenValue(CHART_LEGEND_TEXT_TOKEN)
}
