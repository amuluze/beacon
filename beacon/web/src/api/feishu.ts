import request from '@/api'
import type { FeishuSetting, FeishuUpdateArgs } from '@/interface/feishu'

export async function queryFeishu() {
    return request.get<FeishuSetting>('/api/v1/feishu/query', {})
}

export async function updateFeishu(params: FeishuUpdateArgs) {
    return request.post('/api/v1/feishu/update', params)
}

export async function testFeishu() {
    return request.post('/api/v1/feishu/test', {})
}
