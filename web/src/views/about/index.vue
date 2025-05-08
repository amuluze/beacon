<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const version = ref('v1.3.10')
const lastBuildTime = ref('2025-01-24 11:57:44')
const { t } = useI18n()
</script>

<template>
    <content-wrap :title="t('about.introduction')" :message="t('about.introduction')">
        <h4>{{ t('about.info') }}</h4>
        <p><svg-icon icon-class="done" /> {{ t('about.infoFirst') }}</p>
        <p><svg-icon icon-class="done" /> {{ t('about.infoSecond') }}</p>
        <p><svg-icon icon-class="done" /> {{ t('about.infoThird') }}</p>
    </content-wrap>

    <content-wrap :title="t('about.projectInfo')" :message="t('about.projectInfo')">
        <el-descriptions :column="2" border>
            <el-descriptions-item :label="t('about.version')" label-align="left">
                <el-tag>{{ version }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="t('about.releaseTime')" label-align="left">
                <el-tag>{{ lastBuildTime }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="Github" label-align="left">
                <el-link type="primary" href="https://github.com/amuluze/amprobe" target="_blank">
                    Github
                </el-link>
            </el-descriptions-item>
            <el-descriptions-item :label="t('about.documentation')" label-align="left">
                <el-link type="primary" href="https://amprobe.amuluze.com/doc/" target="_blank">
                    {{ t('about.documentation') }}
                </el-link>
            </el-descriptions-item>
        </el-descriptions>
    </content-wrap>
</template>

<style lang="scss" scoped>
</style>
