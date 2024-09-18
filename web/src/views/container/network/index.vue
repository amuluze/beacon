<template>
    <div class="am-network-title">
        <span @click="$router.push('/docker/container')">容器</span>
        <span @click="$router.push('/docker/image')">镜像</span>
        <span @click="$router.push('/docker/network')">网络</span>
        <span @click="$router.push('/docker/settings')">设置</span>
    </div>
    <div class="am-network-operator">
        <el-card shadow="never">
            <el-button type="primary" plain @click="drawer = true">创建网络</el-button>
            <!-- <el-button type="warning">清理网络</el-button> -->
        </el-card>
    </div>
    <el-card shadow="never">
        <div class="am-table">
            <el-table :data="data" :key="imageKey" highlight-current-row height="100%" stripe v-loading="loading">
                <el-table-column prop="name" label="名称" align="center" width="120" fixed />
                <el-table-column prop="driver" label="模式" align="center" min-width="120" />
                <el-table-column prop="subnet" label="子网" align="center" show-overflow-tooltip width="120" />
                <el-table-column prop="gateway" label="网关" align="center" show-overflow-tooltip width="120" />
                <el-table-column prop="created" label="创建时间" align="center" width="200" sortable />
                <el-table-column label="操作" width="160" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="danger" plain size="small" @click="deleteNetworkByID(scope.row.id)"> 删除 </el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>

        <div class="am-pagination">
            <el-config-provider :locale="zhCn">
                <el-pagination
                    v-model:current-page="pagination.page"
                    :page-size="pagination.size"
                    layout="total, sizes, prev, pager, next, jumper"
                    :page-sizes="pagination.sizeOption"
                    :total="pagination.total"
                    @size-change="pagination.onSizeChange"
                    @current-change="pagination.onPageChange"
                />
            </el-config-provider>
        </div>
    </el-card>
    <div class="am-network-create">
        <el-drawer v-model="drawer" size="540" title="创建网络">
            <div class="am-network-create__content">
                <el-form ref="networkCreateRef" :model="networkCreateMode" :rules="rules" label-width="100px" label-position="left">
                    <el-form-item label="名称" prop="networkName">
                        <el-input v-model="networkCreateMode.networkName" placeholder="请输入名称" />
                    </el-form-item>
                    <el-form-item label="模式" prop="networkMode">
                        <el-select v-model="networkCreateMode.networkMode" style="width: 240px" placeholder="请选择模式">
                            <el-option v-for="item in networkOptions" :key="item.value" :label="item.label" :value="item.value" />
                        </el-select>
                    </el-form-item>
                    <el-form-item label="子网" prop="networkSubnet">
                        <el-input v-model="networkCreateMode.networkSubnet" placeholder="172.16.10.0/24" />
                    </el-form-item>
                    <el-form-item label="网关" prop="networkGateway">
                        <el-input v-model="networkCreateMode.networkGateway" placeholder="172.16.10.1" />
                    </el-form-item>
                    <el-form-item label="标签" prop="networkLabels">
                        <el-input v-model="networkCreateMode.networkLabels" type="textarea" :rows="4" placeholder="一行一个，例如:&#10;key1=value1&#10;key2=value2" />
                    </el-form-item>
                </el-form>
            </div>
            <div class="am-network-create__operator">
                <el-button type="default" size="default" plain @click="drawer = false">取消</el-button>
                <el-button type="primary" size="default" plain @click="confirmCreateNetwork(networkCreateRef)" v-loading="createNetworkLoading"> 确定 </el-button>
            </div>
        </el-drawer>
    </div>
</template>
<script setup lang="ts">
import { createNetwork, deleteNetwork, queryNetworks } from '@/api/container'
import { error, success } from '@/components/Message/message'
import { useTable } from '@/hooks/useTable'
import { NetworkCreateArgs, NetworkDeleteArgs } from '@/interface/container'
import { FormInstance, type FormRules } from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

onMounted(() => {
    refresh()
})

const imageKey = ref(0)
const { data, refresh, loading, pagination } = useTable(queryNetworks, {}, {})
interface RuleForm {
    networkName: string
    networkMode: string
    networkSubnet: string
    networkGateway: string
    networkLabels: string
}
const networkCreateMode = reactive<RuleForm>({
    networkName: '',
    networkMode: 'bridge',
    networkSubnet: '',
    networkGateway: '',
    networkLabels: ''
})

const drawer = ref(false)
const networkCreateRef = ref<FormInstance>()
const rules = reactive<FormRules<RuleForm>>({
    networkName: [{ required: true, message: 'Please input network name', trigger: 'blur' }],
    networkMode: [{ required: true, message: 'Please select network mode', trigger: 'blur' }],
    networkSubnet: [{ required: true, message: 'Please input network subnet', trigger: 'blur' }],
    networkGateway: [{ required: true, message: 'Please input network gateway', trigger: 'blur' }]
})

const networkOptions = [
    {
        value: 'bridge',
        label: 'bridge'
    },
    {
        value: 'host',
        label: 'host'
    }
]

const createNetworkLoading = ref(false)
const confirmCreateNetwork = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl?.validate((valid, fields) => {
        if (!valid) {
            console.log('error fields!', fields)
            error('请检查输入')
            return
        } else {
            createNetworkLoading.value = true
            const ls: Map<string, string> = new Map()
            const labelsArr = networkCreateMode.networkLabels.split('\n')
            labelsArr.forEach((label) => {
                const [key, value] = label.split('=')
                ls.set(key, value)
            })
            let params: NetworkCreateArgs = {
                name: networkCreateMode.networkName,
                driver: networkCreateMode.networkMode,
                network_sebnet: networkCreateMode.networkSubnet,
                network_gateway: networkCreateMode.networkGateway,
                labels: ls
            }
            console.log(params)
            createNetwork(params)
                .then((res) => {
                    const { data } = res
                    console.log(data.network_id)
                    createNetworkLoading.value = false
                    drawer.value = false
                    success('创建成功')
                    refresh()
                })
                .catch((err) => {
                    error(err)
                    createNetworkLoading.value = false
                })
        }
    })
}

const deleteNetworkByID = async (id: string) => {
    console.log(id)
    let params: NetworkDeleteArgs = {
        network_id: id
    }
    const { data } = await deleteNetwork(params)
    console.log(data)
    success('删除成功')
    refresh()
}
</script>
<style scoped lang="scss">
@include b(network-title) {
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
        &:nth-child(3) {
            color: #2f7bff;
            &::before {
                content: '';
                position: absolute;
                left: 148px;
                width: 4px;
                height: 16px;
                text-align: center;
                background-color: #2f7bff;
                border-radius: 2px;
            }
        }
    }
}
@include b(network-operator) {
    height: 48px;
    width: 100%;
    margin-bottom: 4px;
    .el-card {
        height: 100%;
        :deep(.el-card__body) {
            height: 100% !important;
            padding: 0 0 0 16px;
            display: flex;
            align-items: center;
            justify-content: flex-start;
        }
    }
}

@include b(pagination) {
    margin-top: 10px;
    display: flex;
    justify-content: flex-end;
}

@include b(table) {
    width: 100%;
    height: calc(100vh - 230px);
    overflow-y: auto;
}
</style>
