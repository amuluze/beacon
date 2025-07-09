<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const version = ref('v2.0.0')
const lastBuildTime = ref('2025-01-24 11:57:44')
const { t } = useI18n()
onMounted(() => {
    lastBuildTime.value = new Date().toLocaleString()
})
</script>

<template>
    <div class="am-about-container">
        <content-wrap :title="t('about.introduction')" :message="t('about.introduction')">
            <div class="am-about-info">
                <h4>{{ t('about.info') }}</h4>
                <div class="am-about-item">
                    <svg-icon icon-class="done" class="am-about-icon" />
                    <p>{{ t('about.infoFirst') }}</p>
                </div>
                <div class="am-about-item">
                    <svg-icon icon-class="done" class="am-about-icon" />
                    <p>{{ t('about.infoSecond') }}</p>
                </div>
                <div class="am-about-item">
                    <svg-icon icon-class="done" class="am-about-icon" />
                    <p>{{ t('about.infoThird') }}</p>
                </div>
            </div>
        </content-wrap>

        <content-wrap :title="t('about.projectInfo')" :message="t('about.projectInfo')">
            <el-descriptions :column="2" border class="am-about-descriptions">
                <el-descriptions-item :label="t('about.version')" label-align="left">
                    <el-tag effect="plain">{{ version }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item :label="t('about.releaseTime')" label-align="left">
                    <el-tag effect="plain">{{ lastBuildTime }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="Github" label-align="left">
                    <el-link type="primary" href="https://github.com/amuluze/amprobe" target="_blank"
                        class="am-about-link">
                        <svg-icon icon-class="github" class="am-about-link-icon" />
                        Github
                    </el-link>
                </el-descriptions-item>
                <el-descriptions-item :label="t('about.documentation')" label-align="left">
                    <el-link type="primary" href="https://amuluze.com/probe/" target="_blank" class="am-about-link">
                        <svg-icon icon-class="document" class="am-about-link-icon" />
                        {{ t('about.documentation') }}
                    </el-link>
                </el-descriptions-item>
            </el-descriptions>
        </content-wrap>
    </div>
</template>

<style scoped lang="scss">
@include b(about-container) {
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: 16px;
    background-color: var(--el-bg-color);
    overflow-y: auto;

    &::-webkit-scrollbar {
        width: 4px;
    }

    &::-webkit-scrollbar-thumb {
        background: var(--el-border-color-darker);
        border-radius: 2px;
    }

    &::-webkit-scrollbar-track {
        background: var(--el-border-color-light);
        border-radius: 2px;
    }
}

@include b(about-info) {
    h4 {
        font-size: 14px;
        font-weight: 500;
        color: var(--el-text-color-primary);
        margin-top: 0;
        margin-bottom: 16px;
    }
}

@include b(about-item) {
    display: flex;
    align-items: center;
    margin-bottom: 12px;

    p {
        margin: 0;
        font-size: 14px;
        color: var(--el-text-color-regular);
        line-height: 1.5;
    }
}

@include b(about-icon) {
    font-size: 16px;
    color: var(--el-color-success);
    margin-right: 8px;
}

@include b(about-descriptions) {
    :deep(.el-descriptions__label) {
        font-size: 14px;
        color: var(--el-text-color-secondary);
    }

    :deep(.el-descriptions__content) {
        display: flex;
        align-items: center;
    }

    .el-tag {
        border-radius: 4px;
        font-size: 14px;
        padding: 0 8px;
        height: 24px;
        line-height: 24px;
        transition: all 0.3s;
    }
}

@include b(about-link) {
    display: flex;
    align-items: center;
    transition: all 0.3s;

    &:hover {
        transform: translateY(-2px);
    }
}

@include b(about-link-icon) {
    margin-right: 4px;
    font-size: 16px;
}
</style>
