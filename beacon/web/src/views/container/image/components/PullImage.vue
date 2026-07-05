<script setup lang="ts">
import { pullImage } from '@/api/container'
import type { PullImageArgs } from '@/interface/container.ts'
import { useI18n } from 'vue-i18n'

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
  const params: PullImageArgs = {
    image_name: imageNameForPull.value,
  }
  const { data } = await pullImage(params)
  console.log(data)
  imagePullLoading.value = false
  dialogVisible.value = false
  props.update && props.update()
}

const { t } = useI18n()
</script>

<template>
    <el-dialog v-model="dialogVisible" :title="t(props.title as string)" width="50%">
        <div class="am-image-pull">
            <el-input v-model="imageNameForPull" style="margin-right: 8px" :placeholder="t('image.pullImagePlaceholder')" />
            <el-button v-loading="imagePullLoading" size="default" type="primary" plain @click="confirmImagePull">
                {{ t('image.confirm') }}
            </el-button>
        </div>
    </el-dialog>
</template>

<style scoped lang="scss">
@include b(image-pull) {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: flex-start;
}
</style>
