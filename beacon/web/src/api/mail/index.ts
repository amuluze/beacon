/**
 * @Author     : Amu
 * @Date       : 2024/10/17 16:31:02
 * @Description:
 */
import request from '@/api'
import type { Mail, MailCreateArgs, MailTestArgs, MailUpdateArgs } from '@/interface/mail'

export async function queryMail() {
    return request.get<Mail>('/api/v1/mail/mail_query', {})
}

export async function createMail(params: MailCreateArgs) {
    return request.post('/api/v1/mail/mail_create', params)
}

export async function updateMail(params: MailUpdateArgs) {
    return request.post('/api/v1/mail/mail_update', params)
}

export async function testMail(params: MailTestArgs) {
    return request.post('/api/v1/mail/mail_test', params)
}
