<script setup lang="ts">
import { queryAgentList, queryAgentUpdateSetting, updateAgentUpdateSetting } from '@/api/agent'
import { querySystemVersion } from '@/api/system'
import { success } from '@/components/Message/message.ts'
import type { AgentInfo } from '@/interface/agent.ts'
import useCommandComponent from '@/hooks/useCommandComponent.ts'
import UpgradeAgentDialog from '@/components/Agent/UpgradeAgentDialog.vue'
import useStore from '@/store'
import i18n from '@/languages'

const store = useStore()

// 升级 Collia 是管理员专属资源，普通用户只看列表，不暴露升级入口。
const isAdmin = computed(() => store.user.userInfo.isAdmin === 1)

type StatusFilter = 'all' | 'online' | 'offline'

const agents = ref<AgentInfo[]>([])
const loading = ref(false)
const keyword = ref('')
const statusFilter = ref<StatusFilter>('all')
// 最新可用版本来自 /api/v1/system/version，用于在版本列标记「可升级」。
const latestVersion = ref('')

const onlineCount = computed(() => agents.value.filter(item => item.status === 'online').length)
const offlineCount = computed(() => agents.value.length - onlineCount.value)

const filterCounts = computed(() => ({
  all: agents.value.length,
  online: onlineCount.value,
  offline: offlineCount.value,
}))

// 搜索与状态筛选都在前端完成：Agent 数量级小，没有必要为此增加后端查询参数。
const filteredAgents = computed(() => {
  const term = keyword.value.trim().toLowerCase()
  return agents.value.filter((item) => {
    if (statusFilter.value !== 'all' && item.status !== statusFilter.value)
      return false
    if (!term)
      return true
    return (item.hostname || '').toLowerCase().includes(term) || item.agent_id.toLowerCase().includes(term)
  })
})

function isUpgradable(row: AgentInfo) {
  return Boolean(latestVersion.value) && row.version !== latestVersion.value
}

async function loadAgents() {
  loading.value = true
  try {
    const [agentResult, versionResult] = await Promise.all([
      queryAgentList(),
      // 版本信息只用于标记列，取不到时退化为不标注，不影响列表可用。
      querySystemVersion().catch(() => null),
    ])
    agents.value = agentResult.data || []
    const version = versionResult?.data
    latestVersion.value = version?.update_available && version.latest_version
      ? version.latest_version
      : version?.version || ''
  }
  finally {
    loading.value = false
  }
}

// 自动更新开关：记录缺失时后端按默认值（开启）返回；部署级门禁关闭时开关不可操作。
const autoUpdate = ref(true)
const deploymentEnabled = ref(true)
const settingLoading = ref(false)
const settingSaving = ref(false)

async function loadAutoUpdateSetting() {
  if (!isAdmin.value)
    return
  settingLoading.value = true
  try {
    const result = await queryAgentUpdateSetting()
    autoUpdate.value = result.data?.auto_update ?? true
    deploymentEnabled.value = result.data?.deployment_enabled ?? true
  }
  catch {
    // 读取失败保留默认开启，管理员仍可重新保存开关。
  }
  finally {
    settingLoading.value = false
  }
}

async function changeAutoUpdate(value: boolean | string | number) {
  const next = Boolean(value)
  settingSaving.value = true
  try {
    const result = await updateAgentUpdateSetting(next)
    autoUpdate.value = result.data?.auto_update ?? next
    success(i18n.global.t('probe.autoUpdateSaved'))
  }
  catch {
    // 失败时回滚开关；错误提示由全局响应拦截器负责。
    autoUpdate.value = !next
  }
  finally {
    settingSaving.value = false
  }
}

onMounted(() => {
  void loadAgents()
  void loadAutoUpdateSetting()
})

const upgradeAgentDialog = useCommandComponent(UpgradeAgentDialog)

function openUpgrade(row: AgentInfo) {
  upgradeAgentDialog({
    title: 'agent.upgradeTitle',
    agentID: row.agent_id,
    hostname: row.hostname,
    version: row.version,
    onSuccess: () => {
      // 升级后 Agent 版本会变化，重新拉取列表让版本列同步。
      void loadAgents()
    },
  })
}
</script>

<template>
    <div class="am-container">
        <div class="am-probe-toolbar">
            <div class="am-probe-toolbar__left">
                <el-input
                    v-model="keyword"
                    class="am-probe-toolbar__search"
                    size="small"
                    clearable
                    :placeholder="$t('probe.searchPlaceholder')"
                >
                    <template #prefix>
                        <svg-icon icon-class="search" />
                    </template>
                </el-input>
                <el-radio-group v-model="statusFilter" size="small">
                    <el-radio-button value="all">
                        {{ $t('probe.filterAll') }} {{ filterCounts.all }}
                    </el-radio-button>
                    <el-radio-button value="online">
                        {{ $t('probe.online') }} {{ filterCounts.online }}
                    </el-radio-button>
                    <el-radio-button value="offline">
                        {{ $t('probe.offline') }} {{ filterCounts.offline }}
                    </el-radio-button>
                </el-radio-group>
            </div>
            <div class="am-probe-toolbar__right">
                <el-tooltip
                    v-if="isAdmin"
                    placement="top"
                    :disabled="deploymentEnabled"
                    :content="$t('probe.autoUpdateDeploymentDisabled')"
                >
                    <div class="am-probe-toolbar__switch">
                        <span class="am-probe-toolbar__switch-label">{{ $t('probe.autoUpdate') }}</span>
                        <el-switch
                            v-model="autoUpdate"
                            size="small"
                            :loading="settingLoading"
                            :disabled="settingLoading || settingSaving || !deploymentEnabled"
                            @change="changeAutoUpdate"
                        />
                    </div>
                </el-tooltip>
                <span class="am-probe-toolbar__summary">
                    {{ $t('probe.summary', { total: agents.length, online: onlineCount }) }}
                </span>
                <el-button size="small" :loading="loading" @click="loadAgents">
                    <svg-icon icon-class="refresh" />
                    {{ $t('probe.refresh') }}
                </el-button>
            </div>
        </div>
        <div class="am-table">
            <el-table
                v-loading="loading"
                :data="filteredAgents"
                :header-cell-style="{ height: '44px', fontSize: '13px', color: 'var(--am-foreground-primary)', background: 'var(--am-surface-secondary)' }"
                height="100%"
                border
            >
                <el-table-column prop="hostname" :label="$t('probe.columnName')" min-width="240">
                    <template #default="scope">
                        <div class="am-probe__name">
                            {{ scope.row.hostname || scope.row.agent_id }}
                        </div>
                        <div class="am-probe__meta">
                            {{ scope.row.os || '-' }}/{{ scope.row.arch || '-' }}
                        </div>
                    </template>
                </el-table-column>
                <el-table-column prop="status" :label="$t('probe.columnStatus')" width="120">
                    <template #default="scope">
                        <el-tag v-if="scope.row.status === 'online'" type="success">
                            {{ $t('probe.online') }}
                        </el-tag>
                        <el-tag v-else type="danger">
                            {{ $t('probe.offline') }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="version" :label="$t('probe.columnVersion')" width="220">
                    <template #default="scope">
                        <div class="am-probe__version-cell">
                            <span class="am-probe__version">{{ scope.row.version || '-' }}</span>
                            <span v-if="isUpgradable(scope.row as AgentInfo)" class="am-probe__badge am-probe__badge--upgradable">
                                {{ $t('probe.upgradable') }}
                            </span>
                            <span v-else-if="latestVersion && scope.row.version" class="am-probe__badge">
                                {{ $t('probe.latest') }}
                            </span>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column :label="$t('probe.columnAction')" width="160" fixed="right">
                    <template #default="scope">
                        <el-button
                            v-if="isAdmin"
                            type="primary"
                            size="small"
                            text
                            :disabled="scope.row.status !== 'online'"
                            @click="openUpgrade(scope.row as AgentInfo)"
                        >
                            <svg-icon icon-class="update" />
                            {{ $t('probe.upgrade') }}
                        </el-button>
                        <span v-else class="am-probe__placeholder">—</span>
                    </template>
                </el-table-column>
            </el-table>
        </div>
    </div>
</template>

<style scoped lang="scss">
.am-probe-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--am-spacing-md);
  padding-bottom: var(--am-spacing-sm);

  &__left,
  &__right {
    display: flex;
    align-items: center;
    gap: var(--am-spacing-md);
  }

  &__switch {
    display: flex;
    align-items: center;
    gap: var(--am-spacing-xs);
  }

  &__switch-label {
    font-size: var(--am-font-xs);
    color: var(--am-foreground-secondary);
    white-space: nowrap;
  }

  &__search {
    width: 240px;
  }

  &__summary {
    color: var(--am-foreground-secondary);
    font-size: var(--am-font-xs);
  }
}

.am-probe__name {
  color: var(--am-foreground-primary);
  font-size: var(--am-font-sm);
  font-weight: 500;
}

.am-probe__meta {
  color: var(--am-foreground-muted);
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  font-size: var(--am-font-xs);
}

.am-probe__version-cell {
  display: flex;
  align-items: center;
  gap: var(--am-spacing-sm);
}

.am-probe__version {
  min-width: 64px;
  color: var(--am-foreground-secondary);
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  font-size: var(--am-font-xs);
  font-weight: 500;
}

.am-probe__badge {
  padding: 0 6px;
  border-radius: var(--am-radius-xs, 4px);
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
  line-height: 18px;

  &--upgradable {
    background: var(--am-accent-primary-soft);
    color: var(--am-accent-primary);
    font-weight: 600;
  }
}

.am-probe__placeholder {
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
}
</style>
