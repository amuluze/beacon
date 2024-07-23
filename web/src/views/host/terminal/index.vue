<template>
    <div class="am-host-header">
        <span @click="$router.push('/host/monitor')">监控</span>
        <span @click="$router.push('/host/file')">文件</span>
        <span @click="$router.push('/host/terminal')">终端</span>
        <span @click="$router.push('/host/settings')">设置</span>
    </div>
    <div class="am-host-operator">
        <el-card shadow="never">
            <el-row :gutter="1">
                <el-col :span="4">
                    <span>主机: </span>
                    <el-input type="text" v-model="termModel.host" style="width: 140px" placeholder="192.168.1.1" />
                </el-col>
                <el-col :span="4">
                    <span>端口: </span>
                    <el-input type="number" v-model="termModel.port" style="width: 140px" placeholder="22" />
                </el-col>
                <el-col :span="4">
                    <span>用户: </span>
                    <el-input type="text" v-model="termModel.username" style="width: 140px" placeholder="root" />
                </el-col>
                <el-col :span="4">
                    <span>密码: </span>
                    <el-input type="password" v-model="termModel.password" style="width: 140px" placeholder="123456" />
                </el-col>
                <el-col :span="4">
                    <el-button type="primary" plain @click="createConnection">建立连接</el-button>
                    <el-button type="primary" plain @click="closeConnection">断开连接</el-button>
                </el-col>
            </el-row>
        </el-card>
    </div>
    <el-card shadow="never" class="am-host-terminal">
        <el-empty v-if="(termModel.connection = false)" description="暂无终端连接" />
        <div v-else ref="terminal" style="width: 100%; height: 100%" v-loading="termModel.loading" />
    </el-card>
</template>
<script setup lang="ts">
import { Websocket } from '@/components/Websocket'
import { FitAddon } from '@xterm/addon-fit'
import { Terminal } from '@xterm/xterm'

let term = ref<Terminal>()
const terminal = ref()
const fitAddon = new FitAddon()

let ws: Websocket

const initTerminal = () => {
    // 初始化 终端
    term.value = new Terminal({
        convertEol: true, // 启用时，光标将设置为下一行的开头
        disableStdin: false, // 是否应禁用输入
        cursorStyle: 'block', // 光标样式
        cursorBlink: true, // 光标闪烁
        theme: {
            foreground: 'yellow', // 字体
            background: 'black', // 背景
            cursor: 'help' // 设置光标
        }
    })
    term.value.open(terminal.value) //挂载dom窗口
    term.value.loadAddon(fitAddon) //自适应尺寸
    // 不能初始化的时候fit,需要等terminal准备就绪,可以设置延时操作
    setTimeout(() => {
        fitAddon.fit()
    }, 5)

    const onOpen = (ws: Websocket, ev: Event) => {
        console.log(
            '连接成功',
            JSON.stringify({
                host: termModel.host,
                port: termModel.port,
                user: termModel.username,
                password: termModel.password
            }),
            ev
        )
        ws.send(
            JSON.stringify({
                host: termModel.host,
                port: termModel.port,
                user: termModel.username,
                password: termModel.password
            })
        )
    }

    const onMessage = (ws: Websocket, message: MessageEvent) => {
        console.log(ws)
        terminal.value.write(message.data.toString())
    }

    const onError = (ws: Websocket, error: Event) => {
        console.error('连接失败', error)
        ws.close()
    }

    const onClose = (ws: Websocket, close: Event) => {
        console.log('连接关闭', close)
        ws.close()
        terminal.value.write('\r\nWebSSH quit!')
    }

    ws = new Websocket('ws', onOpen, onMessage, onError, onClose)
}

const termModel = reactive({
    host: '',
    port: 0,
    username: '',
    password: '',
    connection: false,
    loading: false
})

const createConnection = () => {
    initTerminal()
}

const closeConnection = () => {
    if (ws) {
        ws.close()
    }

    termModel.connection = false
    termModel.loading = false
}
</script>
<style scoped lang="scss">
@include b(host-header) {
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

    .el-card {
        height: 100%;
        :deep(.el-card__body) {
            height: 100% !important;
            padding: 0 8px 0 0;
            display: flex;
            align-items: center;
            justify-content: flex-end;
        }
    }
}

@include b(host-operator) {
    height: 48px;
    width: 100%;
    margin-bottom: 4px;
    .el-card {
        height: 100%;
        :deep(.el-card__body) {
            height: 100% !important;
            padding: 0 16px 0 16px;
            display: flex;
            flex-direction: row;
            align-items: center;
            justify-content: flex-start;
            .el-row {
                width: 100%;
                .el-col {
                    display: flex;
                    flex-direction: row;
                    align-items: center;
                    justify-content: flex-start;

                    font-size: 14px;
                    color: #575758;
                    span {
                        margin-right: 4px;
                    }
                }
            }
        }
    }
}

@include b(host-terminal) {
    width: 100%;
    height: calc(100vh - 146px);
    overflow-y: auto;
}
</style>
