/**
 * @Author     : Amu
 * @Date       : 2024/11/17 16:21
 * @Description:
 */

export interface Pageable {
    page: number
    size: number
    total: number
    options: number[]
}

export interface Table {
    tableData: any[]
    pageable: Pageable
    loading: boolean
    totalParam: {
        [key: string]: any
    }
}
