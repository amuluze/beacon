export type ToastType = 'info' | 'success' | 'error'

export interface ToastMessage {
    id: number
    message: string
    type: ToastType
}

const toastState = shallowRef<ToastMessage | null>(null)
let toastID = 0
let closeTimer: ReturnType<typeof setTimeout> | undefined

export function dismissToast() {
    if (closeTimer !== undefined) {
        clearTimeout(closeTimer)
        closeTimer = undefined
    }
    toastState.value = null
}

export function showToast(message: string, type: ToastType = 'info', duration = 3200) {
    dismissToast()
    toastID += 1
    toastState.value = { id: toastID, message, type }
    if (duration > 0)
        closeTimer = setTimeout(dismissToast, duration)
}

export function useToast() {
    return {
        toast: readonly(toastState),
        dismissToast,
    }
}
