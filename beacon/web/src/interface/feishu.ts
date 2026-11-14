export interface FeishuSetting {
    id: number
    enabled: boolean
    webhook_masked: string
    webhook_configured: boolean
    secret_configured: boolean
}

export interface FeishuUpdateArgs {
    enabled: boolean
    webhook: string
    secret: string
    clear_secret: boolean
}
