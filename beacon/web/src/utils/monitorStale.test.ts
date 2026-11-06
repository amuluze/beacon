import { describe, expect, it } from 'vitest'

import { formatMonitorAge, formatSampleClock, formatSampleTime } from './monitorStale'

describe('monitor stale formatting', () => {
    it('formats interruption durations across the second/minute/hour boundaries', () => {
        expect(formatMonitorAge(0)).toBe('0 分钟')
        expect(formatMonitorAge(-5)).toBe('0 分钟')
        expect(formatMonitorAge(59)).toBe('59 秒')
        expect(formatMonitorAge(60)).toBe('1 分钟')
        expect(formatMonitorAge(12 * 60)).toBe('12 分钟')
        expect(formatMonitorAge(60 * 60)).toBe('1 小时')
        expect(formatMonitorAge(2 * 60 * 60 + 5 * 60)).toBe('2 小时 5 分钟')
    })

    it('formats sample timestamps at full and clock precision', () => {
        const epoch = 1789000000

        expect(formatSampleTime(epoch)).toMatch(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/)
        expect(formatSampleClock(epoch)).toMatch(/^\d{2}:\d{2}:\d{2}$/)
        expect(formatSampleTime(0)).toBe('-')
        expect(formatSampleClock(0)).toBe('-')
    })
})
