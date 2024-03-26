// export function useTable(loadDataFunc: Function, queryForm: object, deleteDataFunc?: Function) {
//     const loading = ref(true)
//     const tableData = ref()
//     const total = ref(0)

//     const pagination = reactive({
//         page: 1,
//         size: 20
//     })

//     const getData = async () => {
//         loading.value = true
//         const { data } = await loadDataFunc({ ...queryForm, ...pagination })
//         console.log('res', data.data)
//         data.data.map((item: any) => {
//             item.status = item.status === 1
//         })
//         tableData.value = data.data
//         total.value = data.total
//         loading.value = false
//     }

//     onMounted(() => {
//         getData()
//     })

//     // 搜索
//     const handleSearch = async (param: Record<string, any>) => {
//         loading.value = true
//         console.log('search param: ', param)
//         const { data } = await loadDataFunc({ ...param, ...pagination })
//         console.log('search data: ', data)
//         data.data.map((item: any) => {
//             item.status = item.status === 1
//         })
//         tableData.value = data.data
//         total.value = data.total
//         loading.value = false
//     }

//     // 切换页码
//     const handlePageChange = (val: number) => {
//         pagination.page = val
//         getData()
//     }

//     const multipleSelection = ref([])
//     const handleSelectionChange = (val: []) => {
//         multipleSelection.value = val
//     }

//     // 单个删除、批量删除
//     const handleDelete = (id?: string) => {
//         let ids: string[] | null = null
//         if (id != undefined) {
//             ids = [id]
//         } else {
//             ids = multipleSelection.value
//         }

//         if (!ids || ids.length === 0) {
//             warning('请选择要删除的数据')
//             return
//         }
//         ElMessageBox.confirm('确定删除？此操作不可恢复', '提示', {
//             type: 'warning'
//         })
//             .then(async () => {
//                 if (deleteDataFunc === undefined) {
//                     return
//                 }
//                 await deleteDataFunc({ ids })
//                 success('删除成功')
//                 await getData()
//             })
//             .catch(() => {})
//     }

//     // 导出表格
//     const handleExport = () => {
//         ElMessageBox.confirm('确定导出所选数据？', '提示')
//             .then(() => {
//                 ElNotification({
//                     title: '数据导出',
//                     message: '正在为您导出数据，请稍后',
//                     position: 'bottom-right'
//                 })
//             })
//             .catch(() => {})
//     }
//     console.log('containers: ', tableData.value)
//     return {
//         loading,
//         tableData,
//         total,
//         pagination,
//         getData,
//         handleSearch,
//         handleExport,
//         handleSelectionChange,
//         handlePageChange,
//         handleDelete
//     }
// }
import { warning } from '@/components/Message/message'
import { defaults, get, has } from 'lodash-es'

type keyPath = Array<string> | string

export function useTable(
    api: Function,
    params?: object | (() => object),
    options?: {
        path?: { data?: keyPath; total?: keyPath; page?: string; size?: string }
        immediate?: boolean
    }
) {
    // 参数处理
    defaults(options, {
        path: { data: 'data', total: 'total', page: 'page', size: 'size' },
        immediate: false
    })
    console.log('options 1 : ', options)
    // refresh 调用是允许传入参数
    const { pagination } = usePagination((extraData?: object) => (extraData ? refresh(extraData) : refresh()))
    const data = ref([])
    const loading = ref(false)
    console.log('pagination: ', pagination.page, pagination.size)

    const refresh = (extraData?: object) => {
        console.log('extra data: ', extraData)
        console.log('options: ', options)
        const requestData = {
            [options?.path?.page as string]: pagination.page,
            [options?.path?.size as string]: pagination.size
        }
        console.log('---------0', requestData)
        if (params) {
            if (typeof params === 'function') {
                Object.assign(requestData, params())
            } else {
                Object.assign(requestData, params)
            }
        }

        console.log('---------1', requestData)
        if (extraData) {
            if (typeof extraData === 'function') {
                Object.assign(requestData, extraData())
            } else {
                Object.assign(requestData, extraData)
            }
        }

        console.log('---------2', requestData)
        loading.value = true
        return api(requestData)
            .then((res: any) => {
                console.log('---------3', res['data'])
                data.value = get(res['data'], options?.path?.data as string, [])
                pagination.setTotal(get(res['data'], options?.path?.total as string, 0))
                if (
                    !has(res['data'], options?.path?.data as string) ||
                    !has(res['data'], options?.path?.total as string)
                ) {
                    warning('返回数据格式错误')
                }
            })
            .finally(() => {
                loading.value = false
            })
    }
    return { data, refresh, loading, pagination }
}

export function usePagination(cb: Function, sizeOption: Array<number> = [10, 20, 50, 100, 200]) {
    const pagination = reactive({
        page: 1,
        total: 0,
        size: sizeOption[0],
        sizeOption,

        // 维护 page 和 size （一般主动触发）
        onPageChange: (page: number, extraData?: object) => {
            pagination.page = page
            return extraData ? cb(extraData) : cb()
        },

        onSizeChange: (size: number, extraData?: object) => {
            pagination.size = size
            pagination.page = 1
            return extraData ? cb(extraData) : cb()
        },

        // 一般调用 cb 后，还会修改 total，一般是被动触发
        setTotal: (total: number) => {
            pagination.total = total
        },

        reset() {
            pagination.page = 1
            pagination.total = 0
            pagination.size = pagination.sizeOption[0]
        }
    })

    return { pagination }
}
