<script setup lang="ts">
import { queryNetworks, updateContainer } from '@/api/container'
import { error, success } from '@/components/Message/message'
import type { Container, Network } from '@/interface/container'
import {
  containerToForm,
  restartPolicyOptions,
  serializeContainerUpdate,
} from '@/views/container/container/containerForm'
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
  title?: string
  container: Container
  update?: () => void | Promise<void>
}>()

const emit = defineEmits<{
  'update:visible': [visible: boolean]
  'close': []
}>()

const drawerVisible = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
    if (!visible)
      emit('close')
  },
})

const { t } = useI18n()
const formRef = ref<FormInstance>()
const form = reactive(containerToForm(props.container))
const networks = ref<Network[]>([])
const loading = shallowRef(false)

const rules: FormRules = {
  containerName: [{ required: true, message: '请输入容器名称', trigger: 'blur' }],
  imageName: [{ required: true, message: '请输入镜像名称', trigger: 'blur' }],
}

async function submit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) {
    error('请检查表单')
    return
  }

  loading.value = true
  try {
    await updateContainer(serializeContainerUpdate(props.container.id, form))
    success('容器更新成功')
    drawerVisible.value = false
    await props.update?.()
  }
  catch (cause) {
    error(cause instanceof Error ? cause.message : String(cause))
  }
  finally {
    loading.value = false
  }
}

onMounted(async () => {
  const { data } = await queryNetworks({ page: 1, size: 100 })
  networks.value = data.data ?? []
})
</script>

<template>
    <el-drawer v-model="drawerVisible" :title="t(props.title ?? 'container.editContainer')" size="480px" destroy-on-close>
        <el-form ref="formRef" class="container-form" :model="form" :rules="rules" label-position="top">
            <el-form-item :label="t('container.containerName')" prop="containerName">
                <el-input v-model="form.containerName" />
            </el-form-item>
            <el-form-item :label="t('container.imageName')" prop="imageName">
                <el-input v-model="form.imageName" />
            </el-form-item>
            <el-form-item :label="t('container.networkName')">
                <el-select v-model="form.networkName" filterable clearable :placeholder="t('container.keepCurrentNetwork')">
                    <el-option v-for="item in networks" :key="item.id" :label="item.name" :value="item.name" />
                </el-select>
            </el-form-item>
            <el-form-item :label="t('container.containerPort')">
                <el-input v-model="form.ports" type="textarea" :rows="2" placeholder="8080:80（每行一条）" />
            </el-form-item>
            <el-form-item :label="t('container.restartPolicy')">
                <el-select v-model="form.restartPolicy">
                    <el-option v-for="policy in restartPolicyOptions" :key="policy" :label="policy" :value="policy" />
                </el-select>
            </el-form-item>
            <el-form-item :label="t('container.environment')">
                <el-input v-model="form.environments" type="textarea" :rows="3" placeholder="KEY=value（每行一条）" />
            </el-form-item>
            <el-form-item :label="t('container.volume')">
                <el-input v-model="form.volumes" type="textarea" :rows="3" placeholder="/host:/container（每行一条）" />
            </el-form-item>
            <el-form-item :label="t('container.tag')">
                <el-input v-model="form.labels" type="textarea" :rows="3" placeholder="key=value（每行一条）" />
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="drawerVisible = false">
                {{ t('container.cancel') }}
            </el-button>
            <el-button :loading="loading" type="primary" @click="submit">
                {{ t('container.saveChanges') }}
            </el-button>
        </template>
    </el-drawer>
</template>

<style scoped lang="scss">
.container-form {
  padding-bottom: var(--am-spacing-lg);
}
</style>
