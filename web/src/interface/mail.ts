/**
 * @Author     : Amu
 * @Date       : 2024/10/17 16:32:05
 * @Description:
 */

export interface Mail {
    id: number
    server: string
    port: number
    sender: string
    receiver: string
}

export interface MailCreateArgs {
    server: string
    port: number
    sender: string
    password: string
    receiver: string
}

export interface MailUpdateArgs {
    id: number
    server: string
    port: number
    sender: string
    password: string
    receiver: string
}

export interface MailDeleteArgs {
    id: number
}

export interface MailTestArgs {
    receiver: string
}
