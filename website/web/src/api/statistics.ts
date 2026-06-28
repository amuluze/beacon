/**
 * @Author     : Amu
 * @Date       : 2025/02/12 16:58:24
 * @Description:
 */

import type { StatisticsQueryReply, StatisticsUpdateParams, StatisticsUpdateReply } from '~/interface/statistics'

import { useHttp } from '~/composables/useHttp'
import { API } from '~/config/api'

export async function statisticQuery() {
    return useHttp.get<StatisticsQueryReply>(API.statistics_query, {})
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
