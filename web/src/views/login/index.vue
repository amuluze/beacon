<!-- @format -->
<template>
    <div class="am-login">
        <div class="am-login-left">
            <!-- <img src="@/assets/images/login_left.png" alt="" /> -->
            <div class="am-login-left__title">
                <img class="am-login-left__img" src="@/assets/images/amprobe.png" alt="" />
                <h1>Amprobe</h1>
            </div>
            <h3>Amprobe 是一款轻量级主机及 Docker 容器监控工具，它可以轻松的帮助我们完成以下几方面的工作:</h3>
            <div class="am-login-left__item">
                <span>· 监控主机的 CPU、内存、磁盘 IO、网络 IO情况</span>
                <span>· 监控部署于主机上 Docker 容器的运行状态、CPU、内存使用情况</span>
                <span>· 实时查看 Docker 容器的日志，并支持日志下载</span>
            </div>
        </div>

        <div class="am-login-right">
            <div class="am-login-right__form">
                <el-form :model="loginForm" :rules="loginFormRules">
                    <div class="title">
                        <span>登录</span>
                    </div>

                    <el-form-item prop="username">
                        <el-input
                            ref="username"
                            v-model="loginForm.username"
                            size="large"
                            class="username-input"
                            :placeholder="'请输入用户名'"
                            name="username"
                            autocomplete="on"
                            :prefix-icon="User"
                        />
                    </el-form-item>
                    <el-form-item prop="password">
                        <el-input
                            v-model="loginForm.password"
                            size="large"
                            type="password"
                            class="password-input"
                            :placeholder="'请输入密码'"
                            name="password"
                            :prefix-icon="Lock"
                            :show-password="true"
                        >
                        </el-input>
                    </el-form-item>
                    <el-button class="btn" size="large" type="primary" @click.prevent="handleLogin"> 登录 </el-button>
                </el-form>
            </div>
        </div>
    </div>
</template>
<script setup lang="ts">
import { login } from '@/api/auth'
import { loginFormData } from '@/interface/auth.ts'
import useStore from '@/store'
import { getTimeState } from '@/utils/state.ts'
import { Lock, User } from '@element-plus/icons-vue'

const loginForm = reactive<loginFormData>({
    username: 'amprobe',
    password: '123456'
})

const loginFormRules = {
    username: [{ required: true, trigger: 'blur' }],
    password: [{ required: true, trigger: 'blur', validator: passwordValidator }]
}

/**
 * 密码校验器
 */
function passwordValidator(_: any, value: string, callback: any) {
    if (value === '') {
        callback(new Error('password is required'))
    } else if (value.length < 6) {
        callback(new Error('The password can not be less than 6 digits'))
    } else {
        callback()
    }
}

/**
 * 登录
 */
const store = useStore()
const router = useRouter()
const handleLogin = async () => {
    try {
        const { data } = await login({ ...loginForm })
        store.user.setToken(data.access_token, data.refresh_token)
        // 跳转首页
        await router.replace('/')
        ElNotification({
            title: '欢迎来到 Amprobe',
            message: getTimeState(),
            type: 'success'
        })
    } catch (error) {
        console.log(error)
        if (error instanceof Error) ElMessage.error(error.message)
    }
}
</script>
<style scoped lang="scss">
@include b(login) {
    @include flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    height: 100%;
    min-height: 100%;
    background-color: #eee;
    background-image: url('@/assets/images/login_bg.svg');
    background-size: auto;
}
@include b(login-left) {
    @include flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background-image: url('@/assets/images/bg-bd839ea3.svg');
    width: 48%;
    height: 100%;
    min-height: 100%;
    background-color: #22325b;
    color: #fff;
    padding: 20px;

    @include e(title) {
        @include flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;
        margin-bottom: 30px;
    }

    @include e(img) {
        height: 36px;
        width: 36px;
        margin-right: 6px;
    }

    @include e(item) {
        @include flex;
        flex-direction: column;
        align-items: flex-start;
        justify-content: center;
    }
}

@include b(login-right) {
    @include flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    width: 60%;
    height: 100%;

    @include e(form) {
        width: 400px;
        padding: 50px 40px 45px;
        background-color: var(--el-bg-color);
        border-radius: 10px;
        box-shadow: rgb(0 0 0 / 10%) 0 2px 10px 2px;

        .title {
            padding-bottom: 1rem;
            margin: 0 auto;
            text-align: center;

            span {
                font-size: 24px;
                font-weight: bold;
                color: #34495e;
                white-space: nowrap;
            }
        }

        .btn {
            width: 100%;
        }
    }
}
</style>
