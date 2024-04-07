<template>
    <div class="am-term-operator">
        <el-card>
            <span>主机地址:</span> <el-input v-model="ip" style="width: 160px" placeholder="127.0.0.1" />
            <span>端口: </span> <el-input v-model="port" style="width: 80px" placeholder="22" /> <span>用户名:</span>
            <el-input v-model="username" style="width: 120px" placeholder="root" /> <span>密码:</span>
            <el-input v-model="password" type="password" style="width: 120px" placeholder="123456" show-password />
            <el-button type="primary" @click="connect">连接</el-button>
            <el-button type="warning" @click="closeConn">断开连接</el-button>
        </el-card>
    </div>
    <el-card>
        <div id="xterm-container" ref="xtermRef" style="width: 100%; height: 800px" />
    </el-card>
</template>

<script setup lang="ts">
import { Terminal } from '@xterm/xterm'

const xtermRef = ref<HTMLElement>()
const ip = ref('localhost')
const port = ref('22')
const username = ref('root')
const password = ref('123456')

let term: Terminal
let socket: WebSocket

onMounted(() => init())
onBeforeUnmount(() => {
    term.dispose()
    if (socket) {
        socket.close()
    }
})

const init = () => {
    // 初始化 终端
    console.log('初始化终端')
    term = new Terminal({
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

    // 3.websocket集成的插件,这里要注意,网上写了很多websocket相关代码.xterm4版本没必要.
    // const attachAddon = new AttachAddon(socket)
    // const fitAddon = new FitAddon() // 全屏插件
    // // term.loadAddon(attachAddon)
    // term.loadAddon(fitAddon) // 自适应尺寸
    // term.open(xtermRef.value!)
    // fitAddon.fit()
    // term.focus()
}

const writeDefaultInfo = () => {
    let defaultInfo = [
        '┌\x1b[1m terminals \x1b[0m─────────────────────────────────────────────────────────────────┐ ',
        '│                                                                            │ ',
        '│  \x1b[1;34m 欢迎使用Web Docker SSH \x1b[0m                                                  │ ',
        '│                                                                            │ ',
        '└────────────────────────────────────────────────────────────────────────────┘ '
    ]
    // 测试颜色区间
    // let arr = Array.from({length:100},(v,i)=>v = i)
    // console.log(arr)
    // arr.map((item,i) => {
    //     defaultInfo.push(`Hello from \x1B[1;3;${i}m ${i} \x1B[0m  \u2764\ufe0f   ${i}`)
    // })
    term.write('\r\n')
    term.write(defaultInfo.join('\n\n\r'))
    // writeOfColor('我是加粗斜体红色的字呀', '1;3;', '31m')
    term.write('\r\n')
}

const writeOfColor = (txt: string | ArrayBuffer | null, fontCss = '', bgColor = '') => {
    // 在Linux脚本中以 \x1B[ 开始，中间前部分是样式+内容，以 \x1B[0m 结尾
    // 示例 \x1B[1;3;31m 内容 \x1B[0m
    // fontCss
    // 0;-4;字体样式（0;正常 1;加粗 2;变细 3;斜体 4;下划线）
    // bgColor
    // 30m-37m字体颜色（30m:黑色 31m:红色 32m:绿色 33m:棕色字 34m:蓝色 35m:洋红色/紫色 36m:蓝绿色/浅蓝色 37m:白色）
    // 40m-47m背景颜色（40m:黑色 41m:红色 42m:绿色 43m:棕色字 44m:蓝色 45m:洋红色/紫色 46m:蓝绿色/浅蓝色 47m:白色）
    term.write(`\x1B[${fontCss}${bgColor}${txt}\x1B[0m`)
}

// 监听输入
let command = ref('')
const userWrite = () => {
    term.onData((e) => {
        switch (e) {
            case '\u0003': // Ctrl+C
                term.write('^C ')
                term.write('\r\n$ ')
                break
            case '\r': // Enter
                socket.send(command.value)
                command.value = ''
                term.write('\r\n')
                break
            case '\u007F': // Backspace (DEL)
                term.write('\b \b')
                if (command.value.length > 0) {
                    command.value = command.value.slice(0, command.value.length - 1)
                }
                break
            default: // Print all other characters for demo
                if ((e >= String.fromCharCode(0x20) && e <= String.fromCharCode(0x7e)) || e >= '\u00a0') {
                    command.value += e
                    writeOfColor(e, '2;3;', '33m')
                    console.log('用户输入command', command)
                }
        }
    })
}
const connect = () => {
    // 初始化 websocket
    socket = new WebSocket('ws://localhost:8000/ws')
    nextTick(() => {
        userWrite()
    })

    socket.onopen = (ev: Event) => {
        console.log(ev, '建立连接')
        socket.send(
            JSON.stringify({
                ip: ip.value,
                port: port.value,
                username: username.value,
                password: password.value
            })
        )
        socket.send('\n')
        term.open(xtermRef.value!)
        term.focus()
        writeDefaultInfo()
        console.log(ev, '建立连接')
    }
    socket.onmessage = (ev: MessageEvent) => {
        console.log('接收信息：', ev)
        //将字符串转换成 Blob对象
        const blob = new Blob([ev.data], {
            type: 'text/plain'
        })
        //将Blob 对象转换成字符串
        const reader = new FileReader()
        reader.readAsText(blob, 'utf-8')
        reader.onload = (e) => {
            console.log(e)
            // 可以根据返回值判断使用何种颜色或者字体，不过返回值自带了一些字体颜色
            writeOfColor(reader.result, '0;', '37m')
            term.write('\r\n$ ')
        }
    }
}

const closeConn = () => {
    socket.close()
    term.dispose()
}
</script>

<style scoped lang="scss">
@include b(term-operator) {
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
            .el-input {
                margin-left: 8px;
                margin-right: 16px;
            }
            .el-button {
                margin-left: 8px;
            }
        }
    }
}
:deep(.el-card__body) {
    padding: 0;
}
</style>
