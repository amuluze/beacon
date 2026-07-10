<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import useStore from '@/store'

const languageList = [
  { label: '简体中文', value: 'zh' },
  { label: 'English', value: 'en' },
]

const i18n = useI18n()
const store = useStore()
const language = computed(() => store.app.language)

function changeLanguage(lang: string) {
  i18n.locale.value = lang
  store.app.setLanguage(lang)
}
</script>

<template>
    <el-dropdown trigger="click" @command="changeLanguage">
        <button class="am-topbar-action" type="button" aria-label="切换语言">
            <svg-icon size="14px" icon-class="translate" />
        </button>
        <template #dropdown>
            <el-dropdown-menu>
                <el-dropdown-item
                    v-for="item in languageList"
                    :key="item.value"
                    :command="item.value"
                    :disabled="language === item.value"
                >
                    {{ item.label }}
                </el-dropdown-item>
            </el-dropdown-menu>
        </template>
    </el-dropdown>
</template>

<style scoped lang="scss">

</style>
