/**
 * useConfirmCommand — wraps the ConfirmDialog component as a command
 * (imperative call site) like the existing useCommandComponent.
 *
 * The caller configures a title/confirm-message/i18n prefix and an async
 * confirm action. Calling `open(id?)` mounts the dialog and the confirm
 * button invokes the action. Success closes the dialog; failure keeps it
 * open while `onResolved` refreshes the caller's data to reflect partial work.
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
import { createVNode, getCurrentInstance, render } from 'vue'
import type { Component, VNode } from 'vue'

export interface ConfirmCommandOptions {
    title?: string
    message?: string
    // Optional interpolation params for the message key.
    params?: Record<string, string | number>
    i18nPrefix?: string
    // Optional header badge icon; omitted for plain confirmations (draft: 删除确认无图标块).
    badge?: Component
    // Optional width override; defaults to the 480 rung of DESIGN.md §11.1.
    width?: string
    action: (id?: string) => Promise<unknown> | unknown
    // Optional refresh callback, triggered after the action settles (success or failure).
    onResolved?: () => void
}

interface LastInstance {
    vnode: VNode
    container: HTMLElement
}

let lastInstance: LastInstance | null = null

function cleanup(instance: LastInstance | null) {
    if (!instance)
        return
    render(null, instance.container)
    instance.container.remove()
    if (lastInstance === instance)
        lastInstance = null
}

export function useConfirmCommand(options: ConfirmCommandOptions) {
    const appContext = getCurrentInstance()?.appContext ?? null

    return function open(id?: string) {
        cleanup(lastInstance)

        const container = document.createElement('div')
        document.body.appendChild(container)
        let closed = false
        const closeCurrent = () => {
            if (closed)
                return
            closed = true
            cleanup({ vnode, container })
        }

        const vnode = createVNode(ConfirmDialog, {
            'visible': true,
            'title': options.title,
            'badge': options.badge,
            'message': options.message,
            'params': options.params,
            'i18nPrefix': options.i18nPrefix || 'common',
            'confirmation': true,
            'width': options.width ?? '480px',
            'confirmAction': () => options.action(id),
            id,
            'onUpdate:visible': (visible: boolean) => {
                if (!visible)
                    closeCurrent()
            },
            'onResolved': () => {
                options.onResolved?.()
                closeCurrent()
            },
            'onFailed': () => {
                options.onResolved?.()
            },
        })
        vnode.appContext = appContext

        lastInstance = { vnode, container }
        render(vnode, container)
    }
}
