/**
 * @Author     : Amu
 * @Date       : 2025/02/12 16:58:24
 * @Description:
 */

import type { StatisticsQueryReply, StatisticsUpdateParams, StatisticsUpdateReply } from '~/interface/statistics'

import { useHttp } from '~/composables/useHttp'
import { API } from '~/config/api'

export async function statisticQuery() {
    const reply = await useHttp.get<StatisticsQueryReply>(API.statistics_query, {}, { silent: true })
    const statistic = reply?.data
    if (
        !statistic
        || !Number.isInteger(statistic.id)
        || statistic.id <= 0
        || !Number.isInteger(statistic.times)
        || statistic.times < 0
    ) {
        throw new Error('统计响应格式错误')
    }
    return statistic
}

export async function statisticUpdate(params: StatisticsUpdateParams) {
    try {
        return await useHttp.post<StatisticsUpdateReply>(API.statistics_update, params)
    }
    catch (e) {
        console.error('[API Error] statisticUpdate:', e)
        throw new Error('数据更新失败')
    }
}
