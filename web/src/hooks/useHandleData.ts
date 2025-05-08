/**
 * @Author     : Amu
 * @Date       : 2024/11/17 16:10
 * @Description: 操作单条数据信息
 */

/**
 * @description 操作单条数据（二次确认【删除、禁用、启用、重置密码】）
 * @param {Function} api 操作数据接口的 api（必传）
 * @param {object} params 携带的操作数据参数 {id,params} （必传）
 * @param {string} message 提示信息（必传）
 * @returns {Promise}
 */
export async function useHandleData(
    api: (params: any) => Promise<any>,
    params: any = {},
    message: string,
) {
    return new Promise((resolve, reject) => {
        ElMessageBox.confirm(`是否${message}?`, '温馨提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            draggable: true,
        }).then(async () => {
            const res = await api(params)
            if (!res)
                return reject(res.message)
            ElMessage({
                message: `${message} 成功`,
                type: 'success',
            })
            resolve(true)
        }).catch(() => {})
    })
}
