/**
 * @Author     : Amu
 * @Date       : 2025/02/12 16:59:48
 * @Description:
 */
export interface Statistics {
    id: number
    times: number
}

export interface StatisticsQueryReply {
    id: number
    times: number
}

export interface StatisticsUpdateParams {
    id: number
}

export interface StatisticsUpdateReply {}
