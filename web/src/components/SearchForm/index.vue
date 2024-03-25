<template>
    <div class="container">
        <el-form ref="searchFormRef" v-bind="_options" :model="searchForm" inline @submit.prevent>
            <template v-for="(item, _index) in items" :key="_index">
                <!-- 文本输入框 -->
                <el-form-item v-if="item.type === 'text'" :label="item.label" :rules="item.rules" :prop="item.prop">
                    <el-input
                        v-model="searchForm[item.prop]"
                        :readonly="item.readonly"
                        :type="item.type"
                        :placeholder="item.placeholder || item.label"
                        :disabled="item.disabled"
                        :clearable="item.clearable"
                        @keyup.enter="handleKeyUp()"
                    />
                </el-form-item>

                <!-- 下拉框 -->
                <el-form-item
                    v-else-if="item.type === 'select'"
                    :rules="item.rules"
                    :label="item.label"
                    :prop="[item.prop]"
                >
                    <el-select
                        v-model="searchForm[item.prop]"
                        :placeholder="item.options?.placeholder || '请选择'"
                        :clearable="item.clearable"
                    >
                        <el-option
                            v-for="option in item.options?.data"
                            :key="option[item.options?.value || 'value']"
                            :label="option[item.options?.label || 'label']"
                            :value="option[item.options?.value || 'value']"
                        />
                    </el-select>
                </el-form-item>

                <!-- 日期选择器 -->
                <!-- 日期时间选择器 -->
                <el-form-item v-else-if="item.type === 'datetimepicker'" :label="item.label" :prop="item.prop">
                    <el-date-picker
                        v-model="searchForm[item.prop]"
                        type="datetime"
                        placeholder="Pick a Date Time"
                        format="YYYY-MM-DD HH:mm:ss"
                        date-format="MMM DD, YYYY"
                        time-format="HH:mm:ss"
                    />
                </el-form-item>
            </template>

            <!--     按钮插槽       -->
            <el-form-item>
                <slot name="buttons" :model="searchForm" :form-ref="searchFormRef">
                    <el-button type="primary" @click="onSubmit(searchFormRef)">
                        <i-ep-search />{{ _options.submitButtonText }}
                    </el-button>
                    <el-button v-if="_options.showResetButton" type="info" @click="resetForm(searchFormRef)">
                        <i-ep-refresh />{{ _options.resetButtonText }}
                    </el-button>
                </slot>
            </el-form-item>
        </el-form>
    </div>
</template>
<script setup lang="ts">
import { ComputedRef } from 'vue'
import { FormInstance } from 'element-plus'

// 接收父组件传值

interface Props {
    items: Form.FormItem[] // 需要渲染的表达组件
    model?: Record<string, any>
    options?: Form.ButtonOptions
}

// 表单数据对象
const searchFormRef = ref<FormInstance>()
const searchForm = ref<Record<string, any>>({})
const props = defineProps<Props>()
// 根据 items 初始化 formModel
// 如果formModel有传值就用传递的model数据模型，否则就给上面声明的formModel设置相应的(key,value) [item.prop]， item.value是表单的默认值（选填）
watch(
    () => props.model,
    () => {
        props.items.map((item: Form.FormItem) => {
            // 通过 props.model 设置默认值
            props.model ? (searchForm.value = props.model) : (searchForm.value[item.prop] = item.value)
        })
    },
    { immediate: true }
)
const _options: ComputedRef<Form.ButtonOptions> = computed(() => {
    const option = {
        labelPosition: 'right',
        disabled: false,
        submitButtonText: '搜索',
        resetButtonText: '重置'
    }
    return Object.assign(option, props?.options)
})

interface EmitEvent {
    (e: 'submit', params: any): void
    (e: 'reset'): void
}
const emit = defineEmits<EmitEvent>()
defineExpose({
    searchFormRef
})

// 提交
const onSubmit = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.validate((valid) => {
        if (valid) {
            emit('submit', searchForm.value)
        } else {
            return false
        }
    })
}

// 重置
const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
}

const handleKeyUp = () => {}
</script>

<style scoped></style>
