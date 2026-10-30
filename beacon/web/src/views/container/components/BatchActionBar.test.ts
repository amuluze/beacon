import { mount } from '@vue/test-utils'
import { defineComponent } from 'vue'
import { describe, expect, it } from 'vitest'

import BatchActionBar from './BatchActionBar.vue'

const actions = [
    { key: 'start', label: 'container.batchStart' },
    { key: 'delete', label: 'container.batchDelete' },
]

const ElDropdownStub = defineComponent({
    name: 'ElDropdownStub',
    props: { disabled: Boolean },
    emits: ['command'],
    template: '<div class="el-dropdown-stub" :data-disabled="String(disabled)"><slot /><slot name="dropdown" /></div>',
})

const ElDropdownMenuStub = defineComponent({
    name: 'ElDropdownMenuStub',
    template: '<div class="el-dropdown-menu-stub"><slot /></div>',
})

const ElDropdownItemStub = defineComponent({
    name: 'ElDropdownItemStub',
    props: { command: String, disabled: Boolean },
    template: '<span class="el-dropdown-item-stub" :data-command="command" :data-disabled="String(disabled)"><slot /></span>',
})

const ElButtonStub = defineComponent({
    name: 'ElButtonStub',
    props: { disabled: Boolean },
    template: '<button :disabled="disabled"><slot /></button>',
})

function mountBar(count: number) {
    return mount(BatchActionBar, {
        props: { count, actions },
        global: {
            mocks: {
                $t: (key: string, params?: Record<string, unknown>) => (params ? `${key}:${params.count}` : key),
            },
            stubs: {
                'el-dropdown': ElDropdownStub,
                'el-dropdown-menu': ElDropdownMenuStub,
                'el-dropdown-item': ElDropdownItemStub,
                'el-button': ElButtonStub,
            },
        },
    })
}

describe('batch action bar', () => {
    it('shows the selected count and disables the menu when nothing is selected', () => {
        const wrapper = mountBar(0)

        expect(wrapper.text()).toContain('container.selectedCount:0')
        expect(wrapper.find('.el-dropdown-stub').attributes('data-disabled')).toBe('true')
        expect(wrapper.find('.am-table-operator__batch').classes()).toContain('is-empty')
    })

    it('enables the menu and renders every provided action once rows are selected', () => {
        const wrapper = mountBar(3)

        expect(wrapper.text()).toContain('container.selectedCount:3')
        expect(wrapper.find('.el-dropdown-stub').attributes('data-disabled')).toBe('false')
        expect(wrapper.findAll('.el-dropdown-item-stub').map(item => item.attributes('data-command'))).toEqual([
            'start',
            'delete',
        ])
    })

    it('dispatches the selected action key to the host view', () => {
        const wrapper = mountBar(2)

        wrapper.findComponent(ElDropdownStub).vm.$emit('command', 'delete')

        expect(wrapper.emitted('select')).toEqual([['delete']])
    })
})
