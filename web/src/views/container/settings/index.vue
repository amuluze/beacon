<template>
    <div class="am-settings-title">
        <span @click="$router.push('/docker/container')">容器</span>
        <span @click="$router.push('/docker/image')">镜像</span>
        <span @click="$router.push('/docker/network')">网络</span>
        <span @click="$router.push('/docker/settings')">设置</span>
    </div>
    <el-card shadow="never" class="am-settings-image">
        <h4>镜像加速设置</h4>
        <el-divider />
        <el-input v-model="textarea" style="width: 240px" :rows="6" type="textarea" placeholder="https://docker.nju.edu.cn&#10;https://mirror.baidubce.com" />
        <p>优先使用加速 URL 执行操作，设置为空则取消镜像加速。</p>
        <el-button type="primary" plain @click="drawer = true"><svg-icon icon-class="settings" /> 设置</el-button>
    </el-card>
    <div class="am-settings-drawer">
        <el-drawer v-model="drawer" size="540" title="镜像加速">
            <div class="am-settings-drawer__input">
                <el-input v-model="textarea" :rows="6" type="textarea" />
            </div>
            <div class="am-settings-drawer__operator">
                <el-button type="default" size="default" plain @click="cancelEditDockerRegistryMirrors">取消</el-button>
                <el-button type="primary" size="default" plain @click="confirmEditDockerRegistryMirrors"> 确定 </el-button>
            </div>
        </el-drawer>
    </div>
</template>
<script setup lang="ts">
import { SetDockerRegistryMirrors, getDockerRegistryMirrors } from '@/api/container'
import { SetDockerRegistryMirrorsArgs } from '@/interface/container'

const textarea = ref('')
const drawer = ref(false)

const queryDockerRegistryMirrors = async () => {
    const { data } = await getDockerRegistryMirrors()
    textarea.value = data.registry_mirrors.join('\n')
}

const cancelEditDockerRegistryMirrors = () => {
    drawer.value = false
}

const confirmEditDockerRegistryMirrors = async () => {
    let params: SetDockerRegistryMirrorsArgs = {
        registry_mirrors: textarea.value.split('\n').map((item) => item.trim())
    }
    console.log('set mirrors params: ', params)
    await SetDockerRegistryMirrors(params)
    queryDockerRegistryMirrors()
    drawer.value = false
}

onMounted(() => {
    queryDockerRegistryMirrors()
})
</script>
<style scoped lang="scss">
@include b(settings-title) {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: flex-start;
    height: 48px;
    width: 100%;
    background-color: #ffffff;
    // box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);

    border-radius: 4px;
    margin-bottom: 8px;
    padding: 0 16px;
    span {
        display: flex;
        align-items: center;
        font-size: 16px;
        font-weight: 600;
        color: #424244;
        margin-left: 16px;
        margin-right: 16px;
        cursor: pointer;
        &:nth-child(4) {
            color: #2f7bff;
            &::before {
                content: '';
                position: absolute;
                left: 212px;
                width: 4px;
                height: 16px;
                text-align: center;
                background-color: #2f7bff;
                border-radius: 2px;
            }
        }
    }
}

@include b(settings-image) {
    width: 50%;

    :deep(.el-textarea) {
        min-width: 100% !important;
    }

    p {
        color: #bbb6b6;
        font-size: 14px;
    }
}

@include b(settings-drawer) {
    :deep(.el-drawer__header) {
        margin-bottom: 0;
    }

    @include e(input) {
        :deep(.el-textarea) {
            width: 500px !important;
        }
    }

    @include e(operator) {
        margin-top: 16px;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-end;
    }
}
</style>
