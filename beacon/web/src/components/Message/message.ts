export function message(msg: string) {
    ElMessage({
        showClose: true,
        dangerouslyUseHTMLString: true,
        message: msg,
    })
}

/**
 * 成功提示
 * @param msg 提示信息
 */
export function success(msg: string) {
    ElMessage({
        showClose: true,
        dangerouslyUseHTMLString: true,
        message: msg,
        type: 'success',
    })
}

/**
 * 消息提示
 * @param msg 提示信息
 */
export function info(msg: string) {
    ElMessage({
        showClose: true,
        dangerouslyUseHTMLString: true,
        message: msg,
        type: 'info',
    })
}

/**
 * 警告提示
 * @param msg 提示信息
 */
export function warning(msg: string) {
    ElMessage({
        showClose: true,
        dangerouslyUseHTMLString: true,
        message: msg,
        type: 'warning',
    })
}

/**
 * 错误提示
 * @param msg 提示信息
 */
export function error(msg: string) {
    ElMessage({
        showClose: true,
        dangerouslyUseHTMLString: true,
        message: msg,
        type: 'error',
    })
}
