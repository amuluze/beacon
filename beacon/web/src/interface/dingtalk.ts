export interface DingTalkSetting {
    id: number
    enabled: boolean
    webhook_masked: string
    webhook_configured: boolean
    secret_configured: boolean
    at_all: boolean
}

export interface DingTalkUpdateArgs {
    enabled: boolean
    webhook: string
    secret: string
    clear_secret: boolean
    at_all: boolean
}
