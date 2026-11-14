export interface WeComSetting {
    id: number
    enabled: boolean
    webhook_masked: string
    webhook_configured: boolean
    at_all: boolean
}

export interface WeComUpdateArgs {
    enabled: boolean
    webhook: string
    at_all: boolean
}
