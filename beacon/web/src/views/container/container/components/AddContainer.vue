<script setup lang="ts">
import { createContainer, queryImages, queryNetworks } from '@/api/container'
import DialogHeader from '@/components/DialogHeader/index.vue'
import { error, success } from '@/components/Message/message'
import type { Image, Network } from '@/interface/container'
import {
  emptyContainerForm,
  restartPolicyOptions,
  serializeContainerCreate,
} from '@/views/container/container/containerForm'
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import IconPackagePlus from '~icons/lucide/package-plus'

const props = defineProps<{
  visible: boolean
  title?: string
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
const form = reactive(emptyContainerForm())
const images = ref<Image[]>([])
const networks = ref<Network[]>([])
const loading = shallowRef(false)

const rules: FormRules = {
  containerName: [{ required: true, message: '请输入容器名称', trigger: 'blur' }],
  imageName: [{ required: true, message: '请选择镜像', trigger: 'change' }],
  networkName: [{ required: true, message: '请选择网络', trigger: 'change' }],
}

async function loadOptions() {
  const [imageReply, networkReply] = await Promise.all([
    queryImages({ page: 1, size: 100 }),
    queryNetworks({ page: 1, size: 100 }),
  ])
  images.value = imageReply.data.data ?? []
  networks.value = networkReply.data.data ?? []
}

async function submit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) {
    error('请检查表单')
    return
  }

  loading.value = true
  try {
    const network = networks.value.find(item => item.name === form.networkName)
    await createContainer(serializeContainerCreate(form, network && { id: network.id, driver: network.driver }))
    success('容器创建成功')
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

onMounted(() => {
  void loadOptions()
})
</script>

<template>
    <el-drawer v-model="drawerVisible" :title="t(props.title ?? 'container.addContainer')" size="480px" destroy-on-close>
        <template #header>
            <DialogHeader :title="t(props.title ?? 'container.addContainer')" :icon="IconPackagePlus" />
        </template>
        <el-form ref="formRef" class="container-form" :model="form" :rules="rules" label-position="top">
            <el-form-item :label="t('container.containerName')" prop="containerName">
                <el-input v-model="form.containerName" :placeholder="t('container.containerCreatePlaceholder')" />
            </el-form-item>
            <el-form-item :label="t('container.imageName')" prop="imageName">
                <el-select v-model="form.imageName" filterable :placeholder="t('container.containerCreateImagePlaceholder')">
                    <el-option
                        v-for="item in images"
                        :key="item.id"
                        :label="`${item.name}:${item.tag}`"
                        :value="`${item.name}:${item.tag}`"
                    />
                </el-select>
            </el-form-item>
            <el-form-item :label="t('container.networkName')" prop="networkName">
                <el-select v-model="form.networkName" :placeholder="t('container.containerCreateNetworkPlaceholder')">
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
                {{ t('container.confirm') }}
            </el-button>
        </template>
    </el-drawer>
</template>

<style scoped lang="scss">
.container-form {
  padding-bottom: var(--am-spacing-lg);
}
</style>
