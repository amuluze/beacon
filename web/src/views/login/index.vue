<!-- @format -->
<template>
    <div class="am-login">
        <div class="am-login__img">
            <img src="@/assets/images/login_left.png" alt="" />
        </div>

        <div class="am-login__form">
            <el-form :model="loginForm" :rules="loginFormRules">
                <div class="title">
                    <span>Amprobe</span>
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
</template>
<script setup lang="ts">
import { login } from '@/api/auth'
import { loginFormData } from '@/interface/auth.ts'
import useStore from '@/store'
import { getTimeState } from '@/utils/state.ts'
import { Lock, User } from '@element-plus/icons-vue'

const loginForm = reactive<loginFormData>({
    username: 'admin',
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
    const { data } = await login({ ...loginForm })

    // 过期时间为 2h 后，是与后端的约定值
    store.user.setToken(data.token, data.refresh)
    // 跳转首页
    await router.replace('/')
    ElNotification({
        title: '欢迎来到 Amprobe',
        message: getTimeState(),
        type: 'success'
    })
}
</script>
<style scoped lang="scss">
@include b(login) {
    @include flex;
    align-items: center;
    justify-content: center;
    min-height: 100%;
    background-color: #eee;
    background-image: url('@/assets/images/login_bg.svg');
    background-size: auto;

    @include e(img) {
        width: 800px;
        margin-right: 10px;

        img {
            width: 100%;
            height: 100%;
        }
    }

    @include e(form) {
        width: 520px;
        padding: 50px 40px 45px;
        background-color: var(--el-bg-color);
        border-radius: 10px;
        box-shadow: rgb(0 0 0 / 10%) 0 2px 10px 2px;

        .title {
            padding-bottom: 1rem;
            margin: 0 auto;
            text-align: center;

            span {
                font-size: 40px;
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
