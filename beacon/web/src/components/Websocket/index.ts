/**
 * @Author     : Amu
 * @Date       : 2024/4/6 11:29
 * @Description:
 */
/**
 * Websocket 封装
 * @ url： 请求地址       类型： string     默认： ''      备注： 'web/msg'
 */
import { useUserStore } from '@/store/modules/user'
import { useAgentStore } from '@/store/modules/agent'

interface WebsocketOptions {
    agentScoped?: boolean
    agentID?: string
    query?: Record<string, string | undefined>
}

export class Websocket {
    url: string
    ws: WebSocket
    close: Function
    send: Function
    constructor(
        url: string,
        onOpen: ((ws: Websocket, ev: Event) => any) | null = null,
        onMessage: ((ws: Websocket, ev: MessageEvent) => any) | null = null,
        onError: ((ws: Websocket, ev: Event) => any) | null = null,
        onClose: ((ws: Websocket, ev: Event) => any) | null = null,
        options: WebsocketOptions = {},
    ) {
        const location: Location = window.location
        url = `${location.host}/${url}`
        this.url = /https/.test(location.protocol) ? `wss://${url}` : `ws://${url}`
        // 浏览器 ws 握手无法携带 Authorization header，统一把 access token 写入 query。
        const params = new URLSearchParams()
        const token = useUserStore().token
        if (token) {
            params.set('token', token)
        }
        const agentScoped = options.agentScoped ?? true
        const agentID = options.agentID || useAgentStore().selectedAgentID
        if (agentScoped && agentID) {
            params.set('agent_id', agentID)
        }
        Object.entries(options.query || {}).forEach(([key, value]) => {
            if (value) {
                params.set(key, value)
            }
        })
        const query = params.toString()
        if (query) {
            this.url += `?${query}`
        }
        this.ws = new WebSocket(this.url)
        this.close = (): void => {
            this.ws.close()
        }
        this.send = (msg: string): void => {
            this.ws.send(msg)
        }

        this.ws.onopen = (ev: Event): any => {
            if (onOpen !== null) {
                onOpen(this, ev)
            }
        }

        this.ws.onmessage = (ev: MessageEvent): any => {
            if (onMessage !== null) {
                onMessage(this, ev)
            }
        }

        this.ws.onerror = (ev: Event): any => {
            if (onError != null) {
                onError(this, ev)
            }
        }

        this.ws.onclose = (ev: Event): any => {
            if (onClose !== null) {
                onClose(this, ev)
            }
        }
    }
}
