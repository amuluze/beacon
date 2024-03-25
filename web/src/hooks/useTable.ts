import { success, warning } from '@/components/Message/message.ts'

export function useTable(loadDataFunc: Function, queryForm: object, deleteDataFunc?: Function) {
    const loading = ref(true)
    const tableData = ref()
    const total = ref(0)

    const pagination = reactive({
        page: 1,
        size: 20
    })

    const getData = async () => {
        loading.value = true
        const { data } = await loadDataFunc({ ...queryForm, ...pagination })
        console.log('res', data.data)
        data.data.map((item: any) => {
            item.status = item.status === 1
        })
        tableData.value = data.data
        total.value = data.total
        loading.value = false
    }

    onMounted(() => {
        getData()
    })

    // 搜索
    const handleSearch = async (param: Record<string, any>) => {
        loading.value = true
        console.log('search param: ', param)
        const { data } = await loadDataFunc({ ...param, ...pagination })
        console.log('search data: ', data)
        data.data.map((item: any) => {
            item.status = item.status === 1
        })
        tableData.value = data.data
        total.value = data.total
        loading.value = false
    }

    // 切换页码
    const handlePageChange = (val: number) => {
        pagination.page = val
        getData()
    }

    const multipleSelection = ref([])
    const handleSelectionChange = (val: []) => {
        multipleSelection.value = val
    }

    // 单个删除、批量删除
    const handleDelete = (id?: string) => {
        let ids: string[] | null = null
        if (id != undefined) {
            ids = [id]
        } else {
            ids = multipleSelection.value
        }

        if (!ids || ids.length === 0) {
            warning('请选择要删除的数据')
            return
        }
        ElMessageBox.confirm('确定删除？此操作不可恢复', '提示', {
            type: 'warning'
        })
            .then(async () => {
                if (deleteDataFunc === undefined) {
                    return
                }
                await deleteDataFunc({ ids })
                success('删除成功')
                await getData()
            })
            .catch(() => {})
    }

    // 导出表格
    const handleExport = () => {
        ElMessageBox.confirm('确定导出所选数据？', '提示')
            .then(() => {
                ElNotification({
                    title: '数据导出',
                    message: '正在为您导出数据，请稍后',
                    position: 'bottom-right'
                })
            })
            .catch(() => {})
    }
    console.log('containers: ', tableData.value)
    return {
        loading,
        tableData,
        total,
        pagination,
        getData,
        handleSearch,
        handleExport,
        handleSelectionChange,
        handlePageChange,
        handleDelete
    }
}
