/**
 * useConfirmCommand — wraps the ConfirmDialog component as a command
 * (imperative call site) like the existing useCommandComponent.
 *
 * The caller configures a title/confirm-message/i18n prefix and an async
 * confirm action. Calling `open(id?)` mounts the dialog and the confirm
 * button invokes the action; on success/failure (with .finally) the
 * dialog auto-closes.
 *
 * Usage:
 *   const confirmDelete = useConfirmCommand({
 *       title: 'container.deleteContainer',
 *       message: 'container.confirmDelete',
 *       i18nPrefix: 'container',
 *       action: (id: string) => removeContainer({ container_id: id }),
 *   })
 *   // ... button click:
 *   confirmDelete(row.id)
 */
import ConfirmDialog from '@/components/ConfirmDialog/index.vue'
import { createVNode, ref, render } from 'vue'

export interface ConfirmCommandOptions {
    title?: string
    message?: string
    i18nPrefix?: string
    action: (id?: string) => Promise<unknown> | unknown
    // Optional progress callback. Triggered after action resolves.
    onResolved?: () => void
}

interface LastInstance {
    vnode: ReturnType<typeof createVNode>
    container: HTMLElement
    dialogRef: any
}

let lastInstance: LastInstance | null = null

function ensureContainer(): HTMLElement {
    let el = document.getElementById('__confirm_dialog_host__')
    if (!el) {
        el = document.createElement('div')
        el.id = '__confirm_dialog_host__'
        document.body.appendChild(el)
    }
    return el
}

export function useConfirmCommand(options: ConfirmCommandOptions) {
    return function open(id?: string) {
        const container = ensureContainer()

        const dialogRef = ref<{ confirm: () => void; close: () => void } | null>(null)

        const vnode = createVNode(ConfirmDialog, {
            visible: true,
            title: options.title,
            message: options.message,
            i18nPrefix: options.i18nPrefix || 'common',
            confirmation: true,
            width: '500px',
            id,
            'onUpdate:visible': (visible: boolean) => {
                if (!visible) cleanup()
            },
            onConfirm: async () => {
                try {
                    await options.action(id)
                }
                finally {
                    dialogRef.value?.close()
                    options.onResolved?.()
                    cleanup()
                }
            },
            ref: dialogRef,
        })

        render(vnode, container)
        lastInstance = { vnode, container, dialogRef }
    }
}

function cleanup() {
    if (lastInstance) {
        render(null, lastInstance.container)
        lastInstance = null
    }
}
