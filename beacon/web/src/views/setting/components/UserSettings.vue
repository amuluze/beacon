<script setup lang="ts">
import ApiManager from '@/views/account/api/index.vue'
import RoleManager from '@/views/account/role/index.vue'
import UserManager from '@/views/account/user/index.vue'

type SettingsView = 'users' | 'roles' | 'apis'

const tabs: { name: SettingsView, label: string }[] = [
  { name: 'users', label: 'menu.userManager' },
  { name: 'roles', label: 'menu.roleManager' },
  { name: 'apis', label: 'menu.apiManager' },
]

const route = useRoute()
const router = useRouter()

function resolveInitialView(): SettingsView {
  const requested = route?.query.view
  return tabs.some(tab => tab.name === requested) ? requested as SettingsView : 'users'
}

const activeView = ref<SettingsView>(resolveInitialView())
// 与设计稿一致：三个视图首次激活后保持挂载，切换不重复请求。
const visited = reactive(new Set<SettingsView>(['users']))
// 主操作按钮在工具栏右侧，通过票据驱动用户列表打开新建弹层。
const createTicket = ref(0)

watch(activeView, (view) => {
  visited.add(view)
  if (!router)
    return
  void router.replace({ query: { ...route.query, view } })
})

function openCreate() {
  createTicket.value += 1
}
</script>

<template>
    <div class="user-settings">
        <div class="user-settings__toolbar">
            <div class="user-settings__tabs" role="tablist">
                <button
                    v-for="tab in tabs"
                    :key="tab.name"
                    type="button"
                    role="tab"
                    class="user-settings__tab"
                    :class="{ 'user-settings__tab--active': activeView === tab.name }"
                    :aria-selected="activeView === tab.name"
                    @click="activeView = tab.name"
                >
                    {{ $t(tab.label) }}
                </button>
            </div>
            <el-button
                v-if="activeView === 'users'"
                type="primary"
                class="user-settings__primary"
                data-testid="create-user"
                @click="openCreate"
            >
                {{ $t('account.newAccount') }}
            </el-button>
        </div>

        <div class="user-settings__content">
            <UserManager
                v-if="visited.has('users')"
                v-show="activeView === 'users'"
                :create-ticket="createTicket"
            />
            <RoleManager v-if="visited.has('roles')" v-show="activeView === 'roles'" />
            <ApiManager v-if="visited.has('apis')" v-show="activeView === 'apis'" />
        </div>
    </div>
</template>

<style scoped lang="scss">
.user-settings {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-md);
  min-width: 0;
}

.user-settings__toolbar {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: var(--am-spacing-md);
  border-bottom: 1px solid var(--am-border-subtle);
}

.user-settings__tabs {
  display: flex;
  align-items: center;
  gap: var(--am-spacing-lg);
}

.user-settings__tab {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 4px;
  margin-bottom: -1px;
  color: var(--am-foreground-secondary);
  font-family: inherit;
  font-size: var(--am-font-sm);
  font-weight: 500;
  background: transparent;
  border: 0;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  transition:
    color 0.15s ease,
    border-color 0.15s ease;
}

.user-settings__tab:hover:not(.user-settings__tab--active) {
  color: var(--am-foreground-primary);
}

.user-settings__tab--active {
  color: var(--am-accent-primary);
  font-weight: 600;
  border-bottom-color: var(--am-accent-primary);
}

.user-settings__primary {
  margin-bottom: 8px;
}

.user-settings__content {
  height: 560px;
  min-height: 0;

  :deep(.am-container) {
    padding: var(--am-spacing-md) 0 0;
    background: transparent;
  }
}

@media (max-width: 900px) {
  .user-settings__content {
    height: 520px;
  }
}
</style>
