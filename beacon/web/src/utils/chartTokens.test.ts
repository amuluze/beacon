import { afterEach, describe, expect, it } from 'vitest'
import { CHART_TOKEN_VALUES, stubChartTokenValues } from '@/__tests__/helpers/chartTestUtils'
import { CHART_SERIES_TOKENS, resolveChartLegendTextColor, resolveChartPalette, resolveChartTokenValue } from './chartTokens'

const restorers: Array<() => void> = []

function useTokenValues(values?: Record<string, string>) {
    const spy = stubChartTokenValues(values)
    restorers.push(() => spy.mockRestore())
    return spy
}

afterEach(() => {
    restorers.splice(0).forEach(restore => restore())
})

describe('chart series tokens', () => {
    it('lists the six design tokens in chart order without duplicates', () => {
        expect(CHART_SERIES_TOKENS).toEqual([
            '--am-chart-1',
            '--am-chart-2',
            '--am-chart-3',
            '--am-chart-4',
            '--am-chart-5',
            '--am-chart-6',
        ])
        expect(new Set(CHART_SERIES_TOKENS).size).toBe(CHART_SERIES_TOKENS.length)
    })
})

describe('resolveChartTokenValue', () => {
    it('returns the trimmed computed value of the token', () => {
        useTokenValues({ '--am-chart-1': '  #0A7C79  ' })

        expect(resolveChartTokenValue('--am-chart-1')).toBe('#0A7C79')
    })

    it('falls back to an empty string for unknown tokens', () => {
        useTokenValues({})

        expect(resolveChartTokenValue('--am-chart-unknown')).toBe('')
    })
})

describe('resolveChartPalette', () => {
    it('returns the palette in chart-1..chart-6 order', () => {
        useTokenValues()

        expect(resolveChartPalette()).toEqual([
            CHART_TOKEN_VALUES['--am-chart-1'],
            CHART_TOKEN_VALUES['--am-chart-2'],
            CHART_TOKEN_VALUES['--am-chart-3'],
            CHART_TOKEN_VALUES['--am-chart-4'],
            CHART_TOKEN_VALUES['--am-chart-5'],
            CHART_TOKEN_VALUES['--am-chart-6'],
        ])
    })

    it('keeps empty positions when some tokens are undefined', () => {
        useTokenValues({ '--am-chart-2': '#4FB8B4' })

        expect(resolveChartPalette()).toEqual(['', '#4FB8B4', '', '', '', ''])
    })
})

describe('resolveChartLegendTextColor', () => {
    it('reads the secondary foreground token', () => {
        useTokenValues()

        expect(resolveChartLegendTextColor()).toBe('#475569')
    })
})
