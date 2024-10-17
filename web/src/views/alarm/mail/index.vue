<template>
    <div class="am-alarm-title">
        <span @click="$router.push('/alarm/mail')">告警配置</span>
        <span @click="$router.push('/alarm/threshold')">告警阈值</span>
    </div>
    <el-card shadow="never">
        <el-descriptions title="发信邮箱服务配置" :column=1>
            <el-descriptions-item label="告警邮箱">
                {{ mailSettingModel.sender }}
                <svg-icon icon-class="edit" style="cursor: pointer" @click="editMailSetting" />
            </el-descriptions-item>
            <div class="am-alarm-mail-test">
                <el-descriptions-item label="发信测试">
                        <el-input v-model="testReciever" style="width: 240px" size="small" placeholder="请输入收信邮箱" />
                        <el-button style="margin-left: 8px;" size="small" plain type="primary" @click="mailTest">测试</el-button>
                </el-descriptions-item>
            </div>
        </el-descriptions>
    </el-card>
    <div class="am-alamr-edit-mail">
        <el-drawer v-model="mailDrawer" size="540" title="发信邮箱服务设置">
            <el-form ref="mailSettingRef" :model="mailSettingModel" :rules="rules" label-width="120px" label-position="left">
              <el-form-item prop="server" label="邮箱服务器">
                <el-input v-model="mailSettingModel.server" placeholder="请输入邮箱服务器地址"></el-input>
              </el-form-item>
              <el-form-item prop="port" label="邮箱服务端口">
                <el-input v-model="mailSettingModel.port" placeholder="请输入邮箱服务器端口"></el-input>
              </el-form-item>
              <el-form-item prop="sender" label="发信邮箱账号">
                <el-input v-model="mailSettingModel.sender" placeholder="请输入发信邮箱账号"></el-input>
              </el-form-item>
              <el-form-item prop="password" label="发信邮箱密码">
                <el-input v-model="mailSettingModel.password" type="password" placeholder="请输入发信邮箱密码" show-password></el-input>
              </el-form-item>
              <el-form-item prop="receiver" label="收信邮箱">
                <el-input v-model="mailSettingModel.receiver" placeholder="请输入收信邮箱(多个用,分割)"></el-input>
              </el-form-item>
            </el-form>
            <el-button type="default" size="default" plain @click="cancelMaileEdit">取消</el-button>
            <el-button type="primary" size="default" plain @click="confirmMailEdit">确定</el-button>
        </el-drawer>
    </div>
</template>
<script setup lang="ts">
import { createMail, queryMail, testMail, updateMail } from '@/api/mail';
import { info } from '@/components/Message/message';
import { MailCreateArgs, MailTestArgs, MailUpdateArgs } from '@/interface/mail';
import { FormInstance, FormRules } from 'element-plus';

onMounted(async() => {
    const { data } = await queryMail()
    mailSettingModel.id = data.id
    mailSettingModel.server = data.server
    mailSettingModel.port = data.port
    mailSettingModel.sender = data.sender
    mailSettingModel.receiver = data.receiver
})

const mailDrawer = ref(false)
const editMailSetting = () => mailDrawer.value = true
const mailSettingRef = ref<FormInstance>()
interface RuleForm {
    id: number
    server: string
    port: number
    sender: string
    password: string
    receiver: string
}
const mailSettingModel = reactive<RuleForm>({
    id: 0,
    server: '',
    port: 465,
    sender: '',
    password: '******',
    receiver: ''
})
const rules = reactive<FormRules<RuleForm>>({
    server: [
        {
            required: true,
            message: '请输入邮箱服务器地址',
            trigger: 'blur'
        }
    ],
    port: [
        {
            required: true,
            message: '请输入邮箱服务器端口',
            trigger: 'blur'
        }
    ],
    sender: [
        {
            required: true,
            message: '请输入发信邮箱账号',
            trigger: 'blur'
        }
    ],
    password: [
        {
            required: true,
            message: '请输入发信邮箱密码',
            trigger: 'blur'
        }
    ]
})
const cancelMaileEdit = () => {
    mailDrawer.value = false
}
const confirmMailEdit = () => {
    if (mailSettingModel.id == 0) {
        let params: MailCreateArgs = {
            server: mailSettingModel.server,
            port: Number(mailSettingModel.port),
            sender: mailSettingModel.sender,
            password: mailSettingModel.password,
            receiver: mailSettingModel.receiver
        }
        createMail(params).finally(() => {
            info("邮件设置创建成功")
            mailDrawer.value = false
        })
    } else {
        let params: MailUpdateArgs = {
            id: mailSettingModel.id,
            server: mailSettingModel.server,
            port: Number(mailSettingModel.port),
            sender: mailSettingModel.sender,
            password: mailSettingModel.password,
            receiver: mailSettingModel.receiver
        }
        updateMail(params).finally(() => {
            info("邮件设置更新成功")
            mailDrawer.value = false
        })
    }

}

const testReciever = ref('')
const mailTest = () => {
    let Params: MailTestArgs = {
        receiver: testReciever.value
    }
    testMail(Params).finally(() => {
        info("邮件发送成功")
    })
}
</script>
<style scoped lang="scss">
@include b(alarm-title) {
    display: flex;
    align-items: center;
    justify-content: space-start;
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
@include b(alarm-mail-test) {
    display: flex;
    align-items: center;
    justify-content: space-start;
}
</style>
