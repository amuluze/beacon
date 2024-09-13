<template>
    <div class="am-account-header">
        <div class="am-account-header__menu">
            <span @click="$router.push('/account/user')">用户</span>
            <span @click="$router.push('/account/role')">角色</span>
            <span @click="$router.push('/account/api')">接口</span>
        </div>
        <div>
            <el-button type="primary">新增用户</el-button>
        </div>
    </div>
</template>
<script setup lang="ts">

</script>

<style scoped lang="scss">
@include b(account-header) {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    height: 48px;
    width: 100%;
    background-color: #fff;

    border-radius: 4px;
    margin-bottom: 8px;
    padding: 0 16px;

    @include e(menu) {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-start;

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
}
</style>
