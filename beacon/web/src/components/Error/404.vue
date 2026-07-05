<script setup lang="ts">
const router = useRouter()

function goHome() {
  router.replace('/')
}
</script>

<template>
    <div class="am-error-container">
        <img src="@/assets/images/404.png" alt="404" />
        <div class="am-error-container__msg">
            <p>404!</p>
            <span>抱歉，您访问的页面不存在~🤷‍♂️🤷‍♀️</span>
            <el-button type="primary" @click.prevent="goHome">
                返回首页
            </el-button>
        </div>
    </div>
</template>

<style scoped lang="scss">
@include b(error-container) {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;

  @include e(msg) {
    display: flex;
    flex-direction: column;
    margin-left: 30px;
    p {
      font-size: 4rem;
      color: #2f7bff;
      font-weight: 700;
    }
    span {
      font-size: 1.5rem;
      color: #2f7bff;
      margin-bottom: 30px;
    }
  }
}
</style>
