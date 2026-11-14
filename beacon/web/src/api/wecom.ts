import request from '@/api'
import type { WeComSetting, WeComUpdateArgs } from '@/interface/wecom'

export async function queryWeCom() {
    return request.get<WeComSetting>('/api/v1/wecom/query', {})
}

export async function updateWeCom(params: WeComUpdateArgs) {
    return request.post('/api/v1/wecom/update', params)
}

export async function testWeCom() {
    return request.post('/api/v1/wecom/test', {})
}
