/**
 * @Author     : Amu
 * @Date       : 2024/11/10 23:30
 * @Description:
 */

import useStore from '@/store'

/**
 * 是否拥有权限
 * @param permission 元素绑定的权限标识
 */
function hasPermission(permission: string) {
    const store = useStore()
    const permissions: string[] = store.permissions.permissions || []
    return permissions.includes(permission)
}

/**
 * 权限指令
 * 指令是具有一组生命周期的钩子，具有 7 个生命周期
 * - created、beforeMount、mounted、beforeUpdate、updated、beforeUnmount、unmounted
 */
const permission = {
    /**
     * 钩子函数有 3-4 个参数，如下：
     * @param el 指令绑定到的元素，这可用于直接操作 DOM
     * @param binding 一个对象包含一些属性
     *   - name 指令名，不包括 v- 前缀
     *   - value 指令绑定的值
     *   - oldValue 指令绑定的前一个值s
     */
    mounted(el: any, binding: Record<string, any>) {
        if (!binding.value || binding.value.length === 0) {
            return
        }

        const isHave = binding.value.some(hasPermission)

        if (!isHave) {
            // el.style.display = 'none';
            el.parentNode?.removeChild(el)
        }
    },

    updated(el: any, binding: Record<string, any>) {
        if (!binding.value || binding.value.length === 0) {
            return
        }

        const isHave = binding.value.some(hasPermission)

        if (!isHave) {
            // el.style.display = 'none';
            el.parentNode?.removeChild(el)
        }
        else {
            // el.style.display = '';
        }
    },
}

export default permission
