/**
 * @Author     : Amu
 * @Date       : 2024/06/30 15:23:29
 * @Description: https://blog.csdn.net/weixin_45128657/article/details/136388537
 */

import dayjs from 'dayjs'

export function formatDate(date: Date, format: string = 'YYYY-MM-DD HH:mm:ss'): string {
    return dayjs(date).format(format)
}
