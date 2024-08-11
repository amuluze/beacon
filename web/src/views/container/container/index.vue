<template>
    <div class="am-container-title">
        <span @click="$router.push('/docker/container')">容器</span>
        <span @click="$router.push('/docker/image')">镜像</span>
        <span @click="$router.push('/docker/network')">网络</span>
        <span @click="$router.push('/docker/settings')">设置</span>
    </div>
    <div class="am-container-operator">
        <el-card shadow="never">
            <el-button type="primary" plain @click="drawer = true">创建容器</el-button>
            <!-- <el-button type="warning">清理容器</el-button> -->
        </el-card>
    </div>
    <!-- 表格主体 -->
    <el-card shadow="never">
        <!-- https://blog.csdn.net/qq_24950043/article/details/114292940 -->
        <div class="am-table">
            <el-table :data="data" :key="containerKey" stripe ref="multipleTable" v-loading="loading">
                <el-table-column prop="id" label="容器 ID" min-width="100" align="center" fixed sortable />
                <el-table-column
                    prop="name"
                    label="容器名称"
                    align="center"
                    min-width="120"
                    show-overflow-tooltip
                    fixed
                    sortable
                />
                <el-table-column
                    prop="image"
                    label="镜像名称"
                    align="center"
                    min-width="120"
                    show-overflow-tooltip
                    sortable
                />
                <el-table-column
                    prop="ip"
                    label="容器 IP"
                    align="center"
                    min-width="100"
                    show-overflow-tooltip
                    sortable
                />
                <el-table-column prop="ports" label="容器端口" align="center" min-width="100" show-overflow-tooltip />
                <el-table-column prop="state" label="运行状态" align="center" min-width="120" show-overflow-tooltip>
                    <template #default="scope">
                        <el-tag :type="scope.row.state === 'running' ? 'success' : 'danger'">{{
                            scope.row.state
                        }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column
                    prop="uptime"
                    label="启动时间"
                    align="center"
                    min-width="120"
                    show-overflow-tooltip
                    sortable
                />
                <el-table-column
                    prop="cpu_percent"
                    label="CPU使用率"
                    align="center"
                    min-width="130"
                    show-overflow-tooltip
                    sortable
                />
                <el-table-column
                    prop="memory_percent"
                    label="内存使用率"
                    align="center"
                    min-width="130"
                    show-overflow-tooltip
                    sortable
                />
                <el-table-column
                    prop="memory_usage"
                    label="内存使用量"
                    align="center"
                    min-width="130"
                    show-overflow-tooltip
                    sortable
                />
                <el-table-column
                    prop="memory_limit"
                    label="内存限制"
                    align="center"
                    min-width="120"
                    show-overflow-tooltip
                />
                <el-table-column label="操作" width="200" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="primary" size="small" text @click="viewLog(scope.row.id)"> 日志 </el-button>
                        <el-button type="primary" size="small" text @click="startContainerByID(scope.row.id)">
                            启动
                        </el-button>
                        <el-dropdown>
                            <el-button type="primary" size="small" text>更多</el-button>
                            <template #dropdown>
                                <el-dropdown-menu>
                                    <el-dropdown-item>
                                        <el-button
                                            type="warning"
                                            size="small"
                                            text
                                            @click="stopContainerByID(scope.row.id)"
                                        >
                                            停止
                                        </el-button>
                                    </el-dropdown-item>
                                    <el-dropdown-item>
                                        <el-button
                                            type="warning"
                                            size="small"
                                            text
                                            @click="restartContainerByID(scope.row.id)"
                                        >
                                            重启
                                        </el-button>
                                    </el-dropdown-item>
                                    <el-dropdown-item>
                                        <el-button
                                            type="danger"
                                            size="small"
                                            text
                                            @click="deleteContainerByID(scope.row.id)"
                                        >
                                            删除
                                        </el-button>
                                    </el-dropdown-item>
                                </el-dropdown-menu>
                            </template>
                        </el-dropdown>
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
                    @size-change="(size: number) => pagination.onSizeChange(size, params)"
                    @current-change="(page: number) => pagination.onPageChange(page, params)"
                />
            </el-config-provider>
        </div>
    </el-card>

    <!--  查看日志弹窗  -->
    <el-dialog v-model="dialogVisible" title="日志" width="50%" :before-close="handleClose">
        <el-input v-model="logData" rows="20" type="textarea" />
        <template #footer>
            <div class="dialog-footer">
                <el-button size="small" type="primary" plain @click="downloadLog">下载</el-button>
                <el-button size="small" type="info" plain @click="stopLogView">停止</el-button>
                <el-button size="small" type="success" plain @click="handleClose">关闭</el-button>
            </div>
        </template>
    </el-dialog>
    <!-- 创建容器 -->
    <div class="am-container-create">
        <el-drawer v-model="drawer" size="50%" title="创建容器">
            <el-form
                ref="containerCreateRef"
                :model="containerCreateMode"
                :rules="rules"
                label-width="120px"
                label-position="left"
            >
                <el-form-item label="名称" prop="containerName">
                    <el-input v-model="containerCreateMode.containerName" placeholder="请输入名称" />
                </el-form-item>
                <el-form-item label="镜像" prop="imageName">
                    <el-select v-model="containerCreateMode.imageName" style="width: 320px" placeholder="请选择镜像">
                        <el-option
                            v-for="item in imageNameOptions"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                        />
                    </el-select>
                </el-form-item>
                <el-form-item label="网络" prop="networkName">
                    <el-select v-model="containerCreateMode.networkName" style="width: 320px" placeholder="请选择镜像">
                        <el-option
                            v-for="item in networkNameOptions"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                        />
                    </el-select>
                </el-form-item>
                <el-form-item label="端口" prop="ports">
                    <span v-if="containerCreateMode.ports.length === 0">可以添加端口绑定</span>
                    <el-form-item
                        v-for="(port, index) in containerCreateMode.ports"
                        :key="index"
                        style="margin-bottom: 4px"
                    >
                        <el-input v-model="port.hostPort" style="width: 160px" placeholder="服务器端口" />
                        <el-input v-model="port.containerPort" style="width: 160px" placeholder="容器端口" />
                        <el-button type="danger" text @click="deletePort(index)">
                            <svg-icon icon-class="close" />
                        </el-button>
                    </el-form-item>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" text @click="addPort">添加端口<svg-icon icon-class="plus" /></el-button>
                </el-form-item>
                <el-form-item label="环境变量" prop="environments">
                    <span v-if="containerCreateMode.environments.length === 0">可以添加环境变量</span>
                    <el-form-item
                        v-for="(environment, index) in containerCreateMode.environments"
                        :key="index"
                        style="margin-bottom: 4px"
                    >
                        <el-input v-model="environment.key" style="width: 160px" placeholder="变量名" />
                        <el-input v-model="environment.value" style="width: 160px" placeholder="变量值" />
                        <el-button type="danger" text @click="deleteEnvironment(index)">
                            <svg-icon icon-class="close" />
                        </el-button>
                    </el-form-item>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" text @click="addEnvironment">
                        添加环境变量
                        <svg-icon icon-class="plus" />
                    </el-button>
                </el-form-item>
                <el-form-item label="目录/文件挂载" prop="volumes">
                    <span v-if="containerCreateMode.volumes.length === 0">可以添加目录/文件挂载</span>
                    <el-form-item
                        v-for="(volume, index) in containerCreateMode.volumes"
                        :key="index"
                        style="margin-bottom: 4px"
                    >
                        <el-input v-model="volume.hostPath" style="width: 200px" placeholder="服务器目录/文件" />
                        <el-input v-model="volume.containerPath" style="width: 200px" placeholder="容器目录/文件" />
                        <el-button type="danger" text @click="deleteVolume(index)">
                            <svg-icon icon-class="close" />
                        </el-button>
                    </el-form-item>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" text @click="addVolume">
                        添加目录/文件挂载<svg-icon icon-class="plus" />
                    </el-button>
                </el-form-item>
                <el-form-item label="标签(可选)" prop="labels">
                    <el-input
                        v-model="containerCreateMode.labels"
                        type="textarea"
                        :rows="4"
                        placeholder="一行一个，例如:&#10;key1=value1&#10;key2=value2"
                    />
                </el-form-item>
            </el-form>
            <div class="am-container-create__operator">
                <el-button type="default" size="default" plain @click="drawer = false">取消</el-button>
                <el-button
                    type="primary"
                    size="default"
                    plain
                    @click="confirmCreateContainer(containerCreateRef)"
                    v-loading="createContainerLoading"
                >
                    确定
                </el-button>
            </div>
        </el-drawer>
    </div>
</template>

<script setup lang="ts">
import {
    createContainer,
    queryContainers,
    queryImages,
    queryNetworks,
    removeContainer,
    restartContainer,
    startContainer,
    stopContainer
} from '@/api/container'
import { error, success } from '@/components/Message/message'
import { Websocket } from '@/components/Websocket'
import { useTable } from '@/hooks/useTable'
import {
    CreateContainerArgs,
    RemoveContainerArgs,
    RestartContainerArgs,
    StartContainerArgs,
    StopContainerArgs
} from '@/interface/container.ts'
import useStore from '@/store'
import { FormInstance, FormRules } from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
const store = useStore()
onMounted(() => {
    refresh()
    initNetworkOptions()
    initImageOptions()
})

const params = {}
const { data, refresh, loading, pagination } = useTable(queryContainers, {}, {})

const dialogVisible = ref(false)
const containerKey = ref(0)
const logData = ref('')

let ws: Websocket

const viewLog = (container_id: string) => {
    logData.value = ''
    dialogVisible.value = true
    console.log('container_id', container_id)
    console.log('host', location.host)
    console.log('port', location.port)

    const onOpen = (ws: Websocket, ev: Event) => {
        console.log(ev)
        ws.send(container_id)
    }

    const onMessage = (ws: Websocket, ev: MessageEvent) => {
        console.log(ws)
        logData.value = logData.value + '\n' + ev.data
    }

    ws = new Websocket('ws/' + container_id, onOpen, onMessage)
}

const startContainerByID = (container_id: string) => {
    console.log('start container', container_id)
    const params: StartContainerArgs = {
        container_id: container_id
    }
    startContainer(params).finally(() => {
        refresh()
        containerKey.value += 1
    })
}

const stopContainerByID = (container_id: string) => {
    console.log('start container', container_id)
    const params: StopContainerArgs = {
        container_id: container_id
    }
    stopContainer(params).finally(() => {
        refresh()
        containerKey.value += 1
    })
}

const restartContainerByID = (container_id: string) => {
    console.log('start container', container_id)
    const params: RestartContainerArgs = {
        container_id: container_id
    }
    restartContainer(params).finally(() => {
        refresh()
        containerKey.value += 1
    })
}

const deleteContainerByID = (container_id: string) => {
    ElMessageBox.confirm('该操作会强制删除该容器. 继续吗?', '警告', {
        confirmButtonText: 'OK',
        cancelButtonText: 'Cancel',
        type: 'warning',
        beforeClose: (action, instance, done) => {
            if (action === 'confirm') {
                instance.confirmButtonLoading = true
                const params: RemoveContainerArgs = {
                    container_id: container_id
                }
                removeContainer(params)
                    .then(() => {
                        ElMessage({
                            type: 'success',
                            message: '删除完成'
                        })
                    })
                    .finally(() => {
                        instance.confirmButtonLoading = false
                        done()
                        refresh()
                        containerKey.value += 1
                    })
            } else {
                done()
            }
        }
    })
        .then(() => {})
        .catch(() => {
            ElMessage({
                type: 'info',
                message: '删除取消'
            })
        })
}

const handleClose = () => {
    dialogVisible.value = false
    ws.close()
}

const stopLogView = () => {
    ws.close()
}

const downloadLog = () => {
    console.log(`${logData.value}`)
    const a = document.createElement('a')
    a.setAttribute('href', `data:text/plain;charset=utf-8,${encodeURIComponent(logData.value)}`)
    a.setAttribute('download', 'log.txt')
    a.style.display = 'none'
    a.click()
}

const drawer = ref(false)
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
    labels: ''
})

const rules = reactive<FormRules<RuleForm>>({
    containerName: [
        {
            required: true,
            message: '请输入容器名称',
            trigger: 'blur'
        }
    ],
    imageName: [
        {
            required: true,
            message: '请选择镜像',
            trigger: 'blur'
        }
    ],
    networkName: [
        {
            required: true,
            message: '请选择网络',
            trigger: 'blur'
        }
    ]
})

const imageNameOptions = ref<
    {
        value: string
        label: string
    }[]
>([])

const initImageOptions = async () => {
    const params = {
        page: 1,
        size: 100
    }
    const { data } = await queryImages(params)
    store.app.setImages(data.data)
    imageNameOptions.value = data.data.map((item) => {
        return {
            value: item.name + ':' + item.tag,
            label: item.name + ':' + item.tag
        }
    })
}

const networkNameOptions = ref<
    {
        value: string
        label: string
    }[]
>([])

const initNetworkOptions = async () => {
    const params = {
        page: 1,
        size: 100
    }
    const { data } = await queryNetworks(params)
    store.app.setNetworks(data.data)
    networkNameOptions.value = data.data.map((item) => {
        return {
            value: item.name,
            label: item.name
        }
    })
}

const addPort = () => {
    containerCreateMode.ports.push({
        hostPort: '',
        containerPort: ''
    })
}

const deletePort = (index: number) => {
    containerCreateMode.ports.splice(index, 1)
}

const addVolume = () => {
    containerCreateMode.volumes.push({
        hostPath: '',
        containerPath: ''
    })
}

const deleteVolume = (index: number) => {
    containerCreateMode.volumes.splice(index, 1)
}

const addEnvironment = () => {
    containerCreateMode.environments.push({
        key: '',
        value: ''
    })
}

const deleteEnvironment = (index: number) => {
    containerCreateMode.environments.splice(index, 1)
}

const createContainerLoading = ref(false)
const confirmCreateContainer = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl?.validate((valid, fields) => {
        if (!valid) {
            console.log('error submit!', fields)
            error('请检查表单')
            return
        } else {
            createContainerLoading.value = true
            let ps: string[] = []
            let vs: string[] = []
            let es: string[] = []
            let ls: Map<string, string> = new Map()
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
                .filter((item) => item.name === containerCreateMode.networkName)
                .forEach((item) => {
                    if (item.name == containerCreateMode.networkName) {
                        network_mode = item.driver
                        network_id = item.id
                    }
                })
            const params: CreateContainerArgs = {
                container_name: containerCreateMode.containerName,
                image_name: containerCreateMode.imageName,
                network_name: containerCreateMode.networkName,
                network_mode: network_mode,
                network_id: network_id,
                ports: ps,
                volumes: vs,
                environments: es,
                labels: ls
            }
            createContainer(params)
                .then((res) => {
                    const data = res.data
                    console.log('container id: ', data.container_id)
                    createContainerLoading.value = false
                    drawer.value = false
                    success('容器创建成功')
                    refresh()
                })
                .catch((err) => {
                    error(err)
                    createContainerLoading.value = false
                })
        }
    })
}
</script>

<style scoped lang="scss">
@include b(container-title) {
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
        flex-direction: row;
        align-items: center;
        justify-content: flex-start;
        font-size: 16px;
        font-weight: 600;
        margin-left: 16px;
        margin-right: 16px;
        color: #424244;
        cursor: pointer;
        &:first-child {
            color: #2f7bff;
            &::before {
                content: '';
                position: absolute;
                left: 20px;
                width: 4px;
                height: 16px;
                text-align: center;
                background-color: #2f7bff;
                border-radius: 2px;
            }
        }
    }
}
@include b(container-operator) {
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

@include b(container-create) {
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
    height: 100%;
    width: 100%;
}
</style>
