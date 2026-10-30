<!--
    TablePagination — 容器工作区三个表格共用的分页栏。

    容器 / 镜像 / 网络三张表原本各自复制了一份 el-config-provider + el-pagination
    与 locale 计算逻辑；这里收敛为单一组件，避免分页样式与语言包再次分叉。
-->
<script setup lang="ts">
import type { Pageable } from '@/interface/table.ts'

import en from 'element-plus/es/locale/lang/en'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

import useStore from '@/store'
import { getBrowserLanguage } from '@/utils'

const props = defineProps<{ pageable: Pageable }>()

const emit = defineEmits<{
  sizeChange: [number]
  currentChange: [number]
}>()

const store = useStore()
const locale = computed(() => {
  if (store.app.language === 'zh')
    return zhCn
  if (store.app.language === 'en')
    return en
  return getBrowserLanguage() === 'zh' ? zhCn : en
})
</script>

<template>
    <div class="am-pagination">
        <el-config-provider :locale="locale">
            <el-pagination
                :current-page="props.pageable.page"
                :page-size="props.pageable.size"
                :page-sizes="props.pageable.options"
                :total="props.pageable.total"
                layout="total, sizes, prev, pager, next, jumper"
                @size-change="emit('sizeChange', $event)"
                @current-change="emit('currentChange', $event)"
            />
        </el-config-provider>
    </div>
</template>
