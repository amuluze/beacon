import request from '@/api'
import type { DingTalkSetting, DingTalkUpdateArgs } from '@/interface/dingtalk'

export async function queryDingTalk() {
    return request.get<DingTalkSetting>('/api/v1/dingtalk/query', {})
}

export async function updateDingTalk(params: DingTalkUpdateArgs) {
    return request.post('/api/v1/dingtalk/update', params)
}

export async function testDingTalk() {
    return request.post('/api/v1/dingtalk/test', {})
}
