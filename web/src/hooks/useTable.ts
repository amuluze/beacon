/**
 * @Author     : Amu
 * @Date       : 2024/03/27 10:27:31
 * @Description:
 */

import type { Table } from '@/interface/table.ts'

/**
 * @description table 页面操作方法封装
 * @param {Function} api 获取表格数据的 api 方法 (必传)
 * @param {object} params 获取数据的初始化参数 (非必传，默认为{})
 * @param {boolean} isPageable 是否有分页(非必传，默认为 true)
 * @param {Function} dataCallBack 对后台数据进行处理的方法(非必传)
 */

export function useTable(
    api?: (params: any) => Promise<any>,
    params: object = {},
    isPageable: boolean = true,
    dataCallBack?: (data: any) => any,
) {
    const state = reactive<Table>({
        tableData: [],
        pageable: {
            page: 1,
            size: 10,
            total: 0,
            options: [10, 20, 50, 100, 200],
        },
        loading: false,
        totalParam: {},
    })

    /**
     * @description 分页查询参数(只包括分页和表格字段排序，其他排序方式可自行配置)
     */
    const pageParam = computed({
        get: () => {
            return {
                page: state.pageable.page,
                size: state.pageable.size,
            }
        },
        set: () => {},
    })

    /**
     * @description 获取表格数据
     * @return void
     */
    const getList = async () => {
        if (!api)
            return
        try {
            state.loading = true
            // 先把初始化参数和分页参数放到总参数里面
            Object.assign(state.totalParam, params, isPageable ? pageParam.value : {})
            console.log('total param: ', state.totalParam)
            let { data } = await api({ ...state.totalParam })
            dataCallBack && (data = dataCallBack(data))
            console.log(data)
            state.tableData = data.data
            // 解构后台返回的分页数据（如果有分页更新分页信息）
            if (isPageable) {
                state.pageable.total = data.total
            }
        }
        finally {
            state.loading = false
        }
    }

    /**
     * @description 更新查询参数
     * @return void
     */
    const updateTotalParam = async (searchParams: { [key: string]: any }) => {
        console.log('search params:', searchParams)
        Object.assign(state.totalParam, searchParams)
        await getList()
    }

    /**
     * @description 表格数据查询
     * @return void
     */
    const search = async () => {
        state.pageable.page = 1
        await getList()
    }

    /**
     * @description 每页条数改变
     * @param {number} val 当前条数
     * @return void
     */
    const handleSizeChange = async (val: number) => {
        state.pageable.page = 1
        state.pageable.size = val
        await getList()
    }

    /**
     * @description 当前页改变
     * @param {number} val 当前页
     * @return void
     */
    const handleCurrentChange = async (val: number) => {
        console.log('current change:', val)
        state.pageable.page = val
        await getList()
    }

    return {
        ...toRefs(state),
        search,
        handleSizeChange,
        handleCurrentChange,
        updateTotalParam,
    }
}
