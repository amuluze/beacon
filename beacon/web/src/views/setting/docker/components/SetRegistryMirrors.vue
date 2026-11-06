<script setup lang="ts">
import type { SetDockerRegistryMirrorsArgs } from '@/interface/container.ts'
import { SetDockerRegistryMirrors } from '@/api/container'
import DialogHeader from '@/components/DialogHeader/index.vue'
import { useI18n } from 'vue-i18n'
import IconServer from '~icons/lucide/server'

const props = defineProps<{
  visible: boolean
  title?: string
  registryMirrors: string
  update?: () => void
}>()

const emits = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'close'): void
}>()

const drawerVisible = computed<boolean>({
  get() {
    return props.visible
  },
  set(visible: boolean) {
    emits('update:visible', visible)
    if (!visible) {
      emits('close')
    }
  },
})

// 设置docker镜像仓库
const textarea = ref('')
const loading = shallowRef(false)

onMounted(() => {
  textarea.value = props.registryMirrors
})

async function confirmEditDockerRegistryMirrors() {
  const params: SetDockerRegistryMirrorsArgs = {
    registry_mirrors: textarea.value.split('\n').map(item => item.trim()),
  }
  loading.value = true
  try {
    await SetDockerRegistryMirrors(params)
    drawerVisible.value = false
    await props.update?.()
  }
  finally {
    loading.value = false
  }
}
const { t } = useI18n()
</script>

<template>
    <el-drawer v-model="drawerVisible" size="540px" :title="t(props.title as string)">
        <template #header>
            <DialogHeader :title="t(props.title as string)" :icon="IconServer" />
        </template>
        <p class="am-registry-mirrors__hint">
            {{ t('setting.mirrorRegistryDescription') }}
        </p>
        <el-form label-position="top">
            <el-form-item>
                <el-input v-model="textarea" :rows="6" type="textarea" />
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="drawerVisible = false">
                {{ t('setting.cancel') }}
            </el-button>
            <el-button :loading="loading" type="primary" @click="confirmEditDockerRegistryMirrors">
                {{ t('setting.confirm') }}
            </el-button>
        </template>
    </el-drawer>
</template>

<style scoped lang="scss">
.am-registry-mirrors__hint {
  margin: 0 0 var(--am-spacing-md);
  color: var(--am-foreground-muted);
  font-size: var(--am-font-sm);
  line-height: 1.55;
}
</style>
