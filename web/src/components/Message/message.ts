export const message = (msg: string) => {
    ElMessage({
        showClose: true,
        dangerouslyUseHTMLString: true,
        message: msg
    })
}

/**
 * 成功提示
 * @param msg 提示信息
 */
export const success = (msg: string) => {
    ElMessage({
        showClose: true,
        dangerouslyUseHTMLString: true,
        message: msg,
        type: 'success'
    })
}

/**
 * 消息提示
 * @param msg 提示信息
 */
export const info = (msg: string) => {
    ElMessage({
        showClose: true,
        dangerouslyUseHTMLString: true,
        message: msg,
        type: 'info'
    })
}

/**
 * 警告提示
 * @param msg 提示信息
 */
export const warning = (msg: string) => {
    ElMessage({
        showClose: true,
        dangerouslyUseHTMLString: true,
        message: msg,
        type: 'warning'
    })
}

/**
 * 错误提示
 * @param msg 提示信息
 */
export const error = (msg: string) => {
    ElMessage({
        showClose: true,
        dangerouslyUseHTMLString: true,
        message: msg,
        type: 'error'
    })
}

/**
 * confirm 提示框
 * @param title 标题
 * @param msg 信息
 * @param ok ok函数
 * @param okText ok按钮文字
 * @param cancel 取消函数
 * @param cText cancel按钮文字
 */
export const confirm = (title: string, msg: string, ok: any, okText: string, cancel: any, cText: string) => {
    ElMessageBox.confirm(msg, title ? title : '提示', {
        confirmButtonText: okText ? okText : '确定',
        cancelButtonText: cText ? cText : '取消',
        draggable: true
    })
        .then(ok ? ok : () => {})
        .catch(cancel ? cancel : () => {})
}
