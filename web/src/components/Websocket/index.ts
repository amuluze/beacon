/**
 * @Author     : Amu
 * @Date       : 2024/4/6 11:29
 * @Description:
 */
/**
 * Websocket 封装
 * @ url： 请求地址       类型： string     默认： ''      备注： 'web/msg'
 */

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
        onClose: ((ws: Websocket, ev: Event) => any) | null = null
    ) {
        let location: Location = window.location
        url = location.host + '/' + url
        // this.url = /https/.test(location.protocol) ? 'wss://' + url : 'ws://' + url
        this.url = 'ws://101.42.246.113:8000/ws'
        console.log('----->', this.url)
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
