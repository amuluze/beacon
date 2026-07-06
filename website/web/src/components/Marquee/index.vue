<script setup lang="ts">
interface Props {
    text: string
    time?: number
}
const props = withDefaults(defineProps<Props>(), {
    time: 0.5,
})

const rootRef = ref<HTMLElement>()
const textArray = ref<string[]>(props.text.split(''))
const displayedText = ref<string[]>([])
let currentIndex: number = 0
// 修改前
let typingTimer: NodeJS.Timeout | null = null

onBeforeUnmount(() => {
    if (typingTimer) {
        clearTimeout(typingTimer)
    }
})

function typeText() {
    if (currentIndex < props.text.length) {
        displayedText.value.push(textArray.value[currentIndex])
        currentIndex++
        typingTimer = setTimeout(typeText, props.time)
    }
    else {
        currentIndex = 0
        displayedText.value = []
        typingTimer = setTimeout(typeText, props.time)
    }
}

onMounted(() => {
    console.log('mounted', props.time)
    typingTimer = setTimeout(typeText, props.time)
})
onUnmounted(() => {
    if (typingTimer) {
        clearTimeout(typingTimer)
    }
})
</script>

<template>
    <!--  文字打字机 -->
    <div ref="rootRef" class="am-marquee">
        <span v-for="(char, index) in displayedText" :key="index" class="am-marquee__text">{{ char }}</span>
    </div>
</template>

<style scoped lang="scss">
@include b('marquee') {
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;

  @include e('text') {
    font-size: 24px;
    font-weight: bold;
    color: #6c7280;
    line-height: 24px;
  }
}
</style>
