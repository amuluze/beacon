import { mount } from '@vue/test-utils'
import { defineComponent, h, Teleport } from 'vue'
import { afterEach, describe, expect, it, vi } from 'vitest'

import useCommandComponent from '@/hooks/useCommandComponent'

describe('useCommandComponent', () => {
    const closeCommands: Array<() => void> = []

    afterEach(() => {
        closeCommands.splice(0).forEach(close => close())
    })

    it('mounts a fresh command component for repeated invocations', () => {
        const mounted = vi.fn()
        let openCommand: ReturnType<typeof useCommandComponent>

        const Command = defineComponent({
            props: {
                id: String,
                visible: Boolean,
            },
            setup(props) {
                mounted(props.id)
                return () => h('div', props.id)
            },
        })
        const Host = defineComponent({
            setup() {
                openCommand = useCommandComponent(Command)
                closeCommands.push(openCommand.close)
                return () => h('div')
            },
        })

        const wrapper = mount(Host)
        openCommand!({ id: 'container-1' })
        openCommand!({ id: 'container-2' })

        expect(mounted.mock.calls).toEqual([
            ['container-1'],
            ['container-2'],
        ])
        wrapper.unmount()
    })

    it('removes an appended command component when its host unmounts', () => {
        let openCommand: ReturnType<typeof useCommandComponent>
        const Command = defineComponent({
            props: {
                visible: Boolean,
            },
            emits: ['close'],
            setup() {
                return () => h(Teleport, { to: 'body' }, [
                    h('div', { 'data-testid': 'command-overlay' }, 'command'),
                ])
            },
        })
        const Host = defineComponent({
            setup() {
                openCommand = useCommandComponent(Command)
                closeCommands.push(openCommand.close)
                return () => h('div')
            },
        })

        const wrapper = mount(Host)
        openCommand!({})
        expect(document.body.querySelector('[data-testid="command-overlay"]')).not.toBeNull()

        wrapper.unmount()

        expect(document.body.querySelector('[data-testid="command-overlay"]')).toBeNull()
    })
})
