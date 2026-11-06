import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import MonitorStalePanel from './MonitorStalePanel.vue'

const now = Math.floor(Date.now() / 1000)

function mountPanel() {
    return mount(MonitorStalePanel, {
        props: {
            lastSampleAt: now - 12 * 60,
            kpis: [
                { key: 'cpu', label: 'CPU 使用率', value: '23.5%', note: '12 核 · 最近一次采样来自 10:24:08', lastSampleAt: now - 12 * 60 },
                { key: 'memory', label: '内存使用率', value: '72.1%', note: '3.2 GiB / 7.8 GiB · 采样停在上报中断前', lastSampleAt: now - 12 * 60 },
            ],
            rows: [
                { key: 'cpu', label: 'CPU 使用率', value: '23.5%', lastSampleAt: now - 12 * 60 },
                { key: 'memory', label: '内存使用率', value: '72.1%', lastSampleAt: now - 12 * 60 },
                { key: 'disk_usage', label: '磁盘使用率', value: '/dev/vda1 68.4%', lastSampleAt: now - 12 * 60 },
            ],
        },
    })
}

describe('monitor stale panel', () => {
    it('states the interruption window and that data resumes automatically', () => {
        const wrapper = mountPanel()

        const banner = wrapper.get('.am-stale__banner').text()
        expect(banner).toContain('Agent 上报已中断 12 分钟')
        expect(banner).toContain('数据恢复后自动续传，无需手动刷新')
        expect(banner).toContain('最后一次采样')
    })

    it('marks every KPI as expired and keeps the last sampled value', () => {
        const wrapper = mountPanel()
        const kpis = wrapper.findAll('.am-stale__kpi')

        expect(kpis).toHaveLength(2)
        expect(kpis[0].text()).toContain('CPU 使用率')
        expect(kpis[0].text()).toContain('23.5%')
        expect(kpis[0].text()).toContain('数据过期')
        expect(kpis[0].text()).toContain('12 核 · 最近一次采样来自 10:24:08')
    })

    it('lists one detail row per metric with its expiry duration and waiting status', () => {
        const wrapper = mountPanel()
        const rows = wrapper.findAll('.am-stale__row:not(.am-stale__row--header)')

        expect(wrapper.get('.am-stale__detail-summary').text()).toBe('共 3 项指标 · 全部过期')
        expect(rows).toHaveLength(3)
        expect(rows[0].text()).toContain('CPU 使用率')
        expect(rows[0].text()).toContain('23.5%')
        expect(rows[0].text()).toContain('12 分钟')
        expect(rows[0].text()).toContain('等待 Agent 续传')
    })
})
