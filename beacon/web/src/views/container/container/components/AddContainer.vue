<script setup lang="ts">
import { createContainer, queryImages, queryNetworks } from '@/api/container'
import { error, success } from '@/components/Message/message'
import type { CreateContainerArgs } from '@/interface/container'
import useStore from '@/store'
import type { FormInstance, FormRules } from 'element-plus'
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

// 容器创建
const store = useStore()
const containerCreateRef = ref<FormInstance>()
interface Port {
  hostPort: string
  containerPort: string
}

interface Volume {
  hostPath: string
  containerPath: string
}

interface Environment {
  key: string
  value: string
}

interface RuleForm {
  containerName: string
  imageName: string
  networkName: string
  ports: Port[]
  volumes: Volume[]
  environments: Environment[]
  labels: string
}

const containerCreateMode = reactive<RuleForm>({
  containerName: '',
  imageName: '',
  networkName: '',
  ports: [],
  volumes: [],
  environments: [],
  labels: '',
})

const rules = reactive<FormRules<RuleForm>>({
  containerName: [
    {
      required: true,
      message: '请输入容器名称',
      trigger: 'blur',
    },
  ],
  imageName: [
    {
      required: true,
      message: '请选择镜像',
      trigger: 'blur',
    },
  ],
  networkName: [
    {
      required: true,
      message: '请选择网络',
      trigger: 'blur',
    },
  ],
})

const imageNameOptions = ref<{
  value: string
  label: string
}[]>([])

async function initImageOptions() {
  const params = {
    page: 1,
    size: 100,
  }
  const { data } = await queryImages(params)
  store.app.setImages(data.data)
  imageNameOptions.value = data.data.map((item) => {
    return {
      value: `${item.name}:${item.tag}`,
      label: `${item.name}:${item.tag}`,
    }
  })
}

const networkNameOptions = ref<{
  value: string
  label: string
}[]>([])

async function initNetworkOptions() {
  const params = {
    page: 1,
    size: 100,
  }
  const { data } = await queryNetworks(params)
  store.app.setNetworks(data.data)
  networkNameOptions.value = data.data.map((item) => {
    return {
      value: item.name,
      label: item.name,
    }
  })
}

function addPort() {
  containerCreateMode.ports.push({
    hostPort: '',
    containerPort: '',
  })
}

function deletePort(index: number) {
  containerCreateMode.ports.splice(index, 1)
}

function addVolume() {
  containerCreateMode.volumes.push({
    hostPath: '',
    containerPath: '',
  })
}

function deleteVolume(index: number) {
  containerCreateMode.volumes.splice(index, 1)
}

function addEnvironment() {
  containerCreateMode.environments.push({
    key: '',
    value: '',
  })
}

function deleteEnvironment(index: number) {
  containerCreateMode.environments.splice(index, 1)
}

const createContainerLoading = ref(false)
async function confirmCreateContainer(formEl: FormInstance | undefined) {
  if (!formEl)
    return
  await formEl?.validate((valid, fields) => {
    if (!valid) {
      console.log('error submit!', fields)
      error('请检查表单')
    }
    else {
      createContainerLoading.value = true
      const ps: string[] = []
      const vs: string[] = []
      const es: string[] = []
      const ls: Map<string, string> = new Map()
      let network_mode = ''
      let network_id = ''
      for (let i = 0; i < containerCreateMode.ports.length; i++) {
        ps.push(`${containerCreateMode.ports[i].hostPort}:${containerCreateMode.ports[i].containerPort}`)
      }
      for (let i = 0; i < containerCreateMode.volumes.length; i++) {
        vs.push(`${containerCreateMode.volumes[i].hostPath}:${containerCreateMode.volumes[i].containerPath}`)
      }
      for (let i = 0; i < containerCreateMode.environments.length; i++) {
        es.push(`${containerCreateMode.environments[i].key}=${containerCreateMode.environments[i].value}`)
      }

      if (containerCreateMode.labels) {
        const labelArr = containerCreateMode.labels.split('\n')
        for (let i = 0; i < labelArr.length; i++) {
          const label = labelArr[i].split(':')
          ls.set(label[0], label[1])
        }
      }
      store.app.networks
        .filter(item => item.name === containerCreateMode.networkName)
        .forEach((item) => {
          if (item.name === containerCreateMode.networkName) {
            network_mode = item.driver
            network_id = item.id
          }
        })
      const params: CreateContainerArgs = {
        container_name: containerCreateMode.containerName,
        image_name: containerCreateMode.imageName,
        network_name: containerCreateMode.networkName,
        network_mode,
        network_id,
        ports: ps,
        volumes: vs,
        environments: es,
        labels: ls,
      }
      createContainer(params)
        .then((res) => {
          const data = res.data
          console.log('container id: ', data.container_id)
          createContainerLoading.value = false
          drawerVisible.value = false
          success('容器创建成功')
          props.update && props.update()
        })
        .catch((err) => {
          error(err)
          createContainerLoading.value = false
        })
    }
  })
}

const { t } = useI18n()
onMounted(() => {
  initImageOptions()
  initNetworkOptions()
})
</script>

<template>
    <el-drawer v-model="drawerVisible" :destroy-on-close="true" :title="t(props.title as string)" size="50%">
        <el-form ref="containerCreateRef" :model="containerCreateMode" :rules="rules" label-width="180px" label-position="left">
            <el-form-item :label="t('container.containerName')" prop="containerName">
                <el-input v-model="containerCreateMode.containerName" :placeholder="t('container.containerCreatePlaceholder')" />
            </el-form-item>
            <el-form-item :label="t('container.imageName')" prop="imageName">
                <el-select v-model="containerCreateMode.imageName" style="width: 320px" :placeholder="t('container.containerCreateImagePlaceholder')">
                    <el-option v-for="item in imageNameOptions" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
            </el-form-item>
            <el-form-item :label="t('container.networkName')" prop="networkName">
                <el-select v-model="containerCreateMode.networkName" style="width: 320px" :placeholder="t('container.containerCreateNetworkPlaceholder')">
                    <el-option v-for="item in networkNameOptions" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
            </el-form-item>
            <el-form-item :label="t('container.containerPort')" prop="ports">
                <span v-if="containerCreateMode.ports.length === 0">{{ t('container.containerPortTips') }}</span>
                <el-form-item v-for="(port, index) in containerCreateMode.ports" :key="index" style="margin-bottom: 4px">
                    <el-input v-model="port.hostPort" style="width: 160px" placeholder="t('container.serverPortPlaceholder')" />
                    <el-input v-model="port.containerPort" style="width: 160px" placeholder="t('container.containerPortPlaceholder')" />
                    <el-button type="danger" text @click="deletePort(index)">
                        <svg-icon icon-class="close" />
                    </el-button>
                </el-form-item>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" text @click="addPort">
                    {{ t('container.addPort') }}<svg-icon icon-class="plus" />
                </el-button>
            </el-form-item>
            <el-form-item :label="t('container.environment')" prop="environments">
                <span v-if="containerCreateMode.environments.length === 0">{{ t('container.environmentTips') }}</span>
                <el-form-item v-for="(environment, index) in containerCreateMode.environments" :key="index" style="margin-bottom: 4px">
                    <el-input v-model="environment.key" style="width: 160px" :placeholder="t('container.varName')" />
                    <el-input v-model="environment.value" style="width: 160px" :placeholder="t('container.varValue')" />
                    <el-button type="danger" text @click="deleteEnvironment(index)">
                        <svg-icon icon-class="close" />
                    </el-button>
                </el-form-item>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" text @click="addEnvironment">
                    {{ t('container.addEnvironment') }}
                    <svg-icon icon-class="plus" />
                </el-button>
            </el-form-item>
            <el-form-item :label="t('container.volume')" prop="volumes">
                <span v-if="containerCreateMode.volumes.length === 0">{{ t('container.volumeTips') }}</span>
                <el-form-item v-for="(volume, index) in containerCreateMode.volumes" :key="index" style="margin-bottom: 4px">
                    <el-input v-model="volume.hostPath" style="width: 200px" :placeholder="t('container.serverVolumn')" />
                    <el-input v-model="volume.containerPath" style="width: 200px" :placeholder="t('container.containerVolumn')" />
                    <el-button type="danger" text @click="deleteVolume(index)">
                        <svg-icon icon-class="close" />
                    </el-button>
                </el-form-item>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" text @click="addVolume">
                    {{ t('container.addVolume') }}<svg-icon icon-class="plus" />
                </el-button>
            </el-form-item>
            <el-form-item :label="t('container.tag')" prop="labels">
                <el-input v-model="containerCreateMode.labels" type="textarea" :rows="4" :placeholder="t('container.tagPlaceholder')" />
            </el-form-item>
        </el-form>
        <div class="am-container-create__operator">
            <el-button type="default" size="default" plain @click="drawerVisible = false">
                {{ t('container.cancel') }}
            </el-button>
            <el-button v-loading="createContainerLoading" type="primary" size="default" plain @click="confirmCreateContainer(containerCreateRef)">
                {{ t('container.confirm') }}
            </el-button>
        </div>
    </el-drawer>
</template>

<style scoped lang="scss">

</style>
