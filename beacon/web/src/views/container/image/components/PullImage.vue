<script setup lang="ts">
import { pullImage } from '@/api/container'
import DialogHeader from '@/components/DialogHeader/index.vue'
import type { PullImageArgs } from '@/interface/container.ts'
import { useI18n } from 'vue-i18n'
import IconDownload from '~icons/lucide/download'

const props = defineProps<{
  visible: boolean
  title?: string
  update?: () => void
}>()

const emits = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'close'): void
}>()

const dialogVisible = computed<boolean>({
  get() {
    return props.visible
  },
  set(visible) {
    emits('update:visible', visible)
    if (!visible) {
      emits('close')
    }
  },
})

// 拉取镜像
const imagePullLoading = ref(false)
const imageNameForPull = ref('')
async function confirmImagePull() {
  imagePullLoading.value = true
  try {
    const params: PullImageArgs = {
      image_name: imageNameForPull.value,
    }
    await pullImage(params)
    dialogVisible.value = false
    await props.update?.()
  }
  finally {
    imagePullLoading.value = false
  }
}

const { t } = useI18n()
</script>

<template>
    <el-dialog v-model="dialogVisible" :title="t(props.title as string)" width="480px">
        <template #header>
            <DialogHeader :title="t(props.title as string)" :icon="IconDownload" />
        </template>
        <div class="am-image-pull">
            <el-form label-position="top">
                <el-form-item :label="t('image.imageName')">
                    <el-input v-model="imageNameForPull" :placeholder="t('image.pullImagePlaceholder')" />
                </el-form-item>
            </el-form>
            <p class="am-image-pull__hint">
                {{ t('image.pullImageHint') }}
            </p>
        </div>
        <template #footer>
            <el-button @click="dialogVisible = false">
                {{ t('common.cancel') }}
            </el-button>
            <el-button type="primary" :loading="imagePullLoading" @click="confirmImagePull">
                {{ t('image.pullImage') }}
            </el-button>
        </template>
    </el-dialog>
</template>

<style scoped lang="scss">
@include b(image-pull) {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-sm);

  @include e(hint) {
    margin: 0;
    color: var(--am-foreground-muted);
    font-size: var(--am-font-xs);
    line-height: 1.55;
  }
}
</style>
