<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import type { NetworkCreateArgs } from '@/interface/container.ts'
import { createNetwork } from '@/api/container'
import { error, success } from '@/components/Message/message.ts'
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

// 创建网络
const networkCreateRef = ref<FormInstance>()

const networkCreateMode = reactive<RuleForm>({
  networkName: '',
  networkMode: 'bridge',
  networkSubnet: '',
  networkGateway: '',
  networkLabels: '',
})

interface RuleForm {
  networkName: string
  networkMode: string
  networkSubnet: string
  networkGateway: string
  networkLabels: string
}

const rules = reactive<FormRules<RuleForm>>({
  networkName: [{ required: true, message: 'Please input network name', trigger: 'blur' }],
  networkMode: [{ required: true, message: 'Please select network mode', trigger: 'blur' }],
  networkSubnet: [{ required: true, message: 'Please input network subnet', trigger: 'blur' }],
  networkGateway: [{ required: true, message: 'Please input network gateway', trigger: 'blur' }],
})

const networkOptions = [
  {
    value: 'bridge',
    label: 'bridge',
  },
  {
    value: 'host',
    label: 'host',
  },
]

const createNetworkLoading = ref(false)
async function confirmCreateNetwork(formEl: FormInstance | undefined) {
  if (!formEl)
    return
  await formEl?.validate((valid, fields) => {
    if (!valid) {
      console.log('error fields!', fields)
      error('请检查输入')
    }
    else {
      createNetworkLoading.value = true
      const ls: Map<string, string> = new Map()
      const labelsArr = networkCreateMode.networkLabels.split('\n')
      labelsArr.forEach((label) => {
        const [key, value] = label.split('=')
        ls.set(key, value)
      })
      const params: NetworkCreateArgs = {
        name: networkCreateMode.networkName,
        driver: networkCreateMode.networkMode,
        network_sebnet: networkCreateMode.networkSubnet,
        network_gateway: networkCreateMode.networkGateway,
        labels: ls,
      }
      console.log(params)
      createNetwork(params)
        .then((res) => {
          const { data } = res
          console.log(data.network_id)
          createNetworkLoading.value = false
          drawerVisible.value = false
          success('创建成功')
          props.update && props.update()
        })
        .catch((err) => {
          error(err)
          createNetworkLoading.value = false
        })
    }
  })
}

const { t } = useI18n()
</script>

<template>
    <el-drawer v-model="drawerVisible" :title="t(props.title as string)" size="50%">
        <div class="am-network-create__content">
            <el-form ref="networkCreateRef" :model="networkCreateMode" :rules="rules" label-width="160px" label-position="left">
                <el-form-item :label="t('network.networkName')" prop="networkName">
                    <el-input v-model="networkCreateMode.networkName" :placeholder="t('network.networkPlaceholder')" />
                </el-form-item>
                <el-form-item :label="t('network.networkMode')" prop="networkMode">
                    <el-select v-model="networkCreateMode.networkMode" style="width: 240px" :placeholder="t('network.networkModePlaceholder')">
                        <el-option v-for="item in networkOptions" :key="item.value" :label="item.label" :value="item.value" />
                    </el-select>
                </el-form-item>
                <el-form-item :label="t('network.subNetwork')" prop="networkSubnet">
                    <el-input v-model="networkCreateMode.networkSubnet" placeholder="172.16.10.0/24" />
                </el-form-item>
                <el-form-item :label="t('network.networkGateway')" prop="networkGateway">
                    <el-input v-model="networkCreateMode.networkGateway" placeholder="172.16.10.1" />
                </el-form-item>
                <el-form-item :label="t('network.networkTag')" prop="networkLabels">
                    <el-input v-model="networkCreateMode.networkLabels" type="textarea" :rows="4" :placeholder="t('network.networkLabelPlaceholder')" />
                </el-form-item>
            </el-form>
        </div>
        <div class="am-network-create__operator">
            <el-button type="default" size="default" plain @click="drawerVisible = false">
                {{ t('network.cancel') }}
            </el-button>
            <el-button v-loading="createNetworkLoading" type="primary" size="default" plain @click="confirmCreateNetwork(networkCreateRef)">
                {{ t('network.confirm') }}
            </el-button>
        </div>
    </el-drawer>
</template>

<style scoped lang="scss">

</style>
