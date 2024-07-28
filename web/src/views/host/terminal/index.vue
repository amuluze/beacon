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
        <div v-else id="terminal" style="width: 100%; height: 100%" />
    </el-card>
</template>
<script setup lang="ts">
import { warning } from '@/components/Message/message'

// let term: Terminal
// let ws: WebSocket
// const fitAddon = new FitAddon()

// onBeforeUnmount(() => {
//     if (term) {
//         term.dispose()
//     }
//     if (ws) {
//         ws.close()
//     }
// })

// const initWebsocket = () => {
//     const onOpen = (ws: Websocket, ev: Event) => {
//         console.log(
//             '请求连接',
//             JSON.stringify({
//                 host: termModel.host,
//                 port: termModel.port as Number,
//                 user: termModel.username,
//                 password: termModel.password
//             }),
//             ev
//         )
//         ws.send(
//             JSON.stringify({
//                 host: termModel.host,
//                 port: termModel.port as Number,
//                 user: termModel.username,
//                 password: termModel.password
//             })
//         )
//     }

//     const onMessage = (ws: Websocket, message: MessageEvent) => {
//         console.log(ws, message.data)
//         console.log('接收信息：', message)
//         //将字符串转换成 Blob对象
//         const blob = new Blob([message.data], {
//             type: 'text/plain'
//         })
//         //将Blob 对象转换成字符串
//         const reader = new FileReader()
//         reader.readAsText(blob, 'utf-8')
//         reader.onload = (e) => {
//             console.log('ddd', e)
//             // 可以根据返回值判断使用何种颜色或者字体，不过返回值自带了一些字体颜色
//             writeOfColor(reader.result, '0;', '37m')
//             term.write('\r\n$ ')
//         }
//     }

//     const onError = (ws: Websocket, error: Event) => {
//         console.error('连接失败', error)
//         ws.close()
//     }

//     const onClose = (ws: Websocket, ev: Event) => {
//         console.log('连接关闭', ev)
//         ws.close()
//         termModel.connection = false
//         term.write('\r\nWebSSH quit!')
//     }

//     ws = new Websocket('ws', onOpen, onMessage, onError, onClose)
// }
// const writeDefaultInfo = () => {
//     let defaultInfo = [
//         '┌\x1b[1m terminals \x1b[0m─────────────────────────────────────────────────────────────────┐ ',
//         '│                         \x1b[1;34m 欢迎使用 Amprobe SSH \x1b[0m                             │ ',
//         '└────────────────────────────────────────────────────────────────────────────┘ '
//     ]
//     // 测试颜色区间
//     // let arr = Array.from({length:100},(v,i)=>v = i)
//     // console.log(arr)
//     // arr.map((item,i) => {
//     //     defaultInfo.push(`Hello from \x1B[1;3;${i}m ${i} \x1B[0m  \u2764\ufe0f   ${i}`)
//     // })
//     term.write('\r\n')
//     term.write(defaultInfo.join('\n\n\r'))
//     term.write('\r\n')
// }

// const writeOfColor = (txt: string | ArrayBuffer | null, fontCss = '', bgColor = '') => {
//     // 在Linux脚本中以 \x1B[ 开始，中间前部分是样式+内容，以 \x1B[0m 结尾
//     // 示例 \x1B[1;3;31m 内容 \x1B[0m
//     // fontCss
//     // 0;-4;字体样式（0;正常 1;加粗 2;变细 3;斜体 4;下划线）
//     // bgColor
//     // 30m-37m字体颜色（30m:黑色 31m:红色 32m:绿色 33m:棕色字 34m:蓝色 35m:洋红色/紫色 36m:蓝绿色/浅蓝色 37m:白色）
//     // 40m-47m背景颜色（40m:黑色 41m:红色 42m:绿色 43m:棕色字 44m:蓝色 45m:洋红色/紫色 46m:蓝绿色/浅蓝色 47m:白色）
//     term.write(`\x1B[${fontCss}${bgColor}${txt}\x1B[0m`)
// }

// const initTerminal = () => {
//     // 初始化 终端
//     term = new Terminal({
//         convertEol: true, // 启用时，光标将设置为下一行的开头
//         disableStdin: false, // 是否应禁用输入
//         cursorStyle: 'block', // 光标样式
//         cursorBlink: true, // 光标闪烁
//         theme: {
//             foreground: 'yellow', // 字体
//             background: 'black', // 背景
//             cursor: 'help' // 设置光标
//         }
//     })
//     setTimeout(() => {
//         fitAddon.fit()
//     }, 5)
//     term.loadAddon(fitAddon) //自适应尺寸
//     // 不能初始化的时候fit,需要等terminal准备就绪,可以设置延时操作
//     termModel.connection = true
//     term.open(document.getElementById('terminal')!)
//     term.focus()
//     writeDefaultInfo()

//     termData()
// }

// const termData = () => {
//     term?.onData((data) => {
//         ws.send(data)
//     })
//     term?.onResize((size) => {
//         ws.send(JSON.stringify({ cols: size.cols, rows: size.rows }))
//     })
// }

const termModel = reactive({
    host: '',
    port: 0,
    username: '',
    password: '',
    connection: false,
    loading: false
})

const createConnection = () => {
    warning('该功能尚未开放')
    // ws = new WebSocket('ws://101.42.246.113:8000/ws')
    // ws.onopen = (message) => {
    //     console.log('连接成功', message)
    //     termModel.loading = true
    //     termModel.connection = true
    //     term = new Terminal({
    //         convertEol: true, // 启用时，光标将设置为下一行的开头
    //         disableStdin: false, // 是否应禁用输入
    //         cursorStyle: 'block', // 光标样式
    //         cursorBlink: true, // 光标闪烁
    //         theme: {
    //             foreground: 'yellow', // 字体
    //             background: 'black', // 背景
    //             cursor: 'help' // 设置光标
    //         }
    //     })
    //     setTimeout(() => {
    //         fitAddon.fit()
    //     }, 5)
    //     term.loadAddon(fitAddon) //自适应尺寸
    //     // 不能初始化的时候fit,需要等terminal准备就绪,可以设置延时操作
    //     term.open(document.getElementById('terminal')!)
    //     term.focus()
    //     writeDefaultInfo()
    //     ws.send(
    //         JSON.stringify({
    //             host: termModel.host,
    //             port: termModel.port as Number,
    //             user: termModel.username,
    //             password: termModel.password
    //         })
    //     )
    // }
    // ws.onmessage = (ev: MessageEvent) => {
    //     console.log('接收信息：', ev)
    //     //将字符串转换成 Blob对象
    //     const blob = new Blob([ev.data], {
    //         type: 'text/plain'
    //     })
    //     //将Blob 对象转换成字符串
    //     const reader = new FileReader()
    //     reader.readAsText(blob, 'utf-8')
    //     reader.onload = (e) => {
    //         console.log('ddd', e)
    //         // 可以根据返回值判断使用何种颜色或者字体，不过返回值自带了一些字体颜色
    //         writeOfColor(reader.result, '0;', '37m')
    //         term.write('\r\n$ ')
    //     }
    // }
    // ws.onclose = (ev: Event) => {
    //     termModel.connection = false
    //     console.log('关闭连接', ev, termModel.connection)
    // }
}

const closeConnection = () => {
    warning('该功能尚未开放')
    // if (ws) {
    //     ws.close()
    // }
    // termModel.connection = false
    // console.log('关闭', termModel.connection)
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
