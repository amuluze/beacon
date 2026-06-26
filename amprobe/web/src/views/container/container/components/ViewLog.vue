<script setup lang="ts">
import { Websocket } from '@/components/Websocket';
import { useI18n } from 'vue-i18n';

import Codemirror from "codemirror-editor-vue3";
import "codemirror/mode/javascript/javascript.js";

const cmOptions = {
  mode: "log",
  theme: "default",

}

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
  set(visible: boolean) {
    emits('update:visible', visible)
    if (!visible) {
      emits('close')
    }
  },
})

// 查看日志
let ws: Websocket
const logData = ref('')

function viewLog(container_id: string): void {
  logData.value = ''
  dialogVisible.value = true

  const onOpen = (_ws: Websocket, _ev: Event) => {
    ws.send(container_id)
  }

  const onMessage = (_ws: Websocket, ev: MessageEvent) => {
    logData.value = `${logData.value}\n${ev.data}`
  }

  ws = new Websocket(`ws/${container_id}`, onOpen, onMessage)
}

onMounted(() => {
  viewLog(props.id)
})

onUnmounted(() => {
  ws.close()
})

function downloadLog() {
  const a = document.createElement('a')
  a.setAttribute('href', `data:text/plain;charset=utf-8,${encodeURIComponent(logData.value)}`)
  a.setAttribute('download', 'log.txt')
  a.style.display = 'none'
  a.click()
}

function stopLogView() {
  ws.close()
}

function handleClose() {
  dialogVisible.value = false
  ws.close()
}

const { t } = useI18n()
</script>

<template>
    <!--  查看日志弹窗  -->
    <el-dialog v-model="dialogVisible" :title="t('container.log')" width="50%" :destroy-on-close="true">
        <Codemirror
            v-model:value="logData"
            :options="cmOptions"
            border
            height="100%"
            width="100%"
        />
        <template #footer>
            <div class="dialog-footer">
                <el-button size="small" type="primary" plain @click="downloadLog">
                    {{ t('container.download') }}
                </el-button>
                <el-button size="small" type="info" plain @click="stopLogView">
                    {{ t('container.stop') }}
                </el-button>
                <el-button size="small" type="success" plain @click="handleClose">
                    {{ t('container.close') }}
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>

<style scoped lang="scss">

</style>
