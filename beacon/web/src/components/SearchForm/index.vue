<script setup lang="ts">
import type { FormInstance } from 'element-plus'
import { useI18n } from 'vue-i18n'

/**
 * conditionList   组件外部传入的查询条件的配置列表数据
 */
interface Props {
  items: Form.Item[]
  model: Record<string, any>
  search: (params: any) => void
}
const props = withDefaults(defineProps<Props>(), {
  items: () => [],
})

// 查询条件组件的实例
const searchFormRef = ref<FormInstance>()
// 查询条件表单数据
const searchForm = ref<Record<string, any>>({})

function reset() {
  searchFormRef.value?.resetFields()
}

// 切换展开和收起查询条件
const collapsed = ref(true)
function toggleCollapsed() {
  collapsed.value = !collapsed.value
}
// 初始化折叠查询条件的断点，从第几个查询条件开始（默认是从第3个，因为默认配置的span值是6）
const initConditionLen = 3

// 展示右侧按钮组（折叠||收起按钮）
const showConBtn = computed(() => {
  return props.items.length > initConditionLen
})

onMounted(() => {
  props.items.forEach((item) => {
    searchForm.value[item.prop] = props.model[item.prop]
    // 如果类型为checkbox，默认值需要设置一个空数组
    const value = item.type === 'checkbox' ? [] : ''
    // 通过 prop.mode 可以传入一些默认值
    props.model ? (searchForm.value = props.model) : (searchForm.value[item.prop] = item.value || value)
  })
})
const { t } = useI18n()
</script>

<template>
    <div class="am-search">
        <el-form ref="searchFormRef" :inline="true" :model="searchForm" label-width="auto">
            <el-row :class="{ 'condition-col': !collapsed }" align="middle" :gutter="5">
                <el-col v-for="(item, index) in props.items.slice(0, collapsed ? initConditionLen : props.items.length)" :key="index" class="custom-col" :xl="6" :lg="6" :md="8" :sm="12" :xs="24">
                    <el-form-item :label="t(item.label)" :prop="item.prop" style="width: 90%">
                        <!-- 输入框 -->
                        <el-input
                            v-if="item.type === 'input'"
                            v-model="searchForm[item.prop]"
                            :type="item.type ?? 'text'"
                            :placeholder="t(item.placeholder as string)"
                            :clearable="item.clearable"
                            :disabled="item.disabled"
                        />

                        <!-- 下拉框 -->
                        <el-select
                            v-if="item.type === 'select'"
                            v-model="searchForm[item.prop]"
                            :placeholder="t(item.placeholder as string)"
                            :clearable="item.clearable"
                            :multiple="item.multiple"
                        >
                            <el-option
                                v-for="option in item.options?.data"
                                :key="option[item.options?.value || 'value']"
                                :label="option[item.options?.label || 'label']"
                                :value="option[item.options?.value || 'value']"
                            />
                        </el-select>

                        <!-- 日期选择器 -->
                        <el-date-picker
                            v-if="item.type === 'datepicker'"
                            v-model="searchForm[item.prop]"
                        />

                        <!-- 时间范围选择器 -->
                        <el-date-picker
                            v-if="item.type === 'datetimerange'"
                            v-model="searchForm[item.prop]"
                            :type="item.type"
                            range-separator="To"
                            start-placeholder="Start date"
                            end-placeholder="End date"
                        />
                    </el-form-item>
                </el-col>
                <el-col :xl="6" :lg="6" :md="8" :sm="12" :xs="24" class="custom-col">
                    <el-form-item class="btn-group-item flex-end">
                        <el-button type="primary" plain @click="search(searchForm)">
                            <svg-icon icon-class="search" style="margin-right: 4px" />
                            {{ t('search.search') }}
                        </el-button>
                        <el-button type="danger" plain @click="reset">
                            <svg-icon icon-class="update" style="margin-right: 4px" />
                            {{ t('search.reset') }}
                        </el-button>
                        <el-button v-show="showConBtn" type="primary" link @click="toggleCollapsed">
                            <span v-if="collapsed">
                                {{ t('search.unfold') }}
                            </span>
                            <span v-else>{{ t('search.fold') }}</span>
                            <!-- {{ collapsed ? "展开" : "收起" }} -->
                            <svg-icon v-if="collapsed" icon-class="down" />
                            <svg-icon v-else icon-class="up" />
                        </el-button>
                    </el-form-item>
                </el-col>
            </el-row>
        </el-form>
    </div>
</template>

<style scoped lang="scss">
@include b(search) {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 16px;
  margin-bottom: 8px;
  border-radius: 4px;
  border: 1px solid var(--el-header-border-color);
  color: var(--el-header-text-color);
  background-color: var(--el-header-bg-color);
  .el-form {
    width: 100%;
    height: 100%;
    .el-form-item {
      margin: 0;
    }
  }
}

.condition-col {
  .el-col {
    margin-bottom: 8px;
  }
}

.custom-col {
  /*默认值*/
  margin-bottom: 0;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
}
/* xl */
@media (min-width: 1920px) {
  .custom-col {
    margin-bottom: 0;
  }
}
/* lg */
@media (max-width: 1919px) {
  .custom-col {
    margin-bottom: 0;
  }
}
/* md */
@media (max-width: 1199px) {
  .custom-col {
    margin-bottom: 8px;
  }
}
/* sm */
@media (max-width: 991px) {
  .custom-col {
    margin-bottom: 8px;
  }
}
/* xs */
@media (max-width: 767px) {
  .custom-col {
    margin-bottom: 8px;
  }
}
</style>
