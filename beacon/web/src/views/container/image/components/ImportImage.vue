<script setup lang="ts">
import { success } from '@/components/Message/message.ts'
import useStore from '@/store'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
  id: string
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

// 导入镜像
const store = useStore()
const imageImportLoading = ref(false)

function onSuccess() {
  success('导入成功')
  imageImportLoading.value = false
  props.update && props.update()
}
const { t } = useI18n()
</script>

<template>
    <el-dialog v-model="dialogVisible" :title="t(props.title as string)" width="50%">
        <el-upload
            drag
            action="/api/v1/container/image_import"
            :headers="{
                Authorization: `Bearer ${store.user.token}`,
            }"
            :limit="1"
            :loading="imageImportLoading"
            multiple
            :on-success="onSuccess"
        >
            <svg-icon icon-class="upload" />
            <div class="el-upload__text">
                Drop file here or <em>click to upload</em>
            </div>
            <template #tip>
                <div class="el-upload__tip">
                    tar/tar.gz files
                </div>
            </template>
        </el-upload>
    </el-dialog>
</template>

<style scoped lang="scss">

</style>
