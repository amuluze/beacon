<script setup lang="ts">
import { createRole, queryResource } from '@/api/account'
import DialogHeader from '@/components/DialogHeader/index.vue'
import { success } from '@/components/Message/message'
import type { FormInstance, FormRules } from 'element-plus'
import type { Resource } from '@/interface/account.ts'
import { useI18n } from 'vue-i18n'
import IconShieldPlus from '~icons/lucide/shield-plus'

const props = defineProps<{
  visible: boolean
  title?: string
  update?: () => void
}>()

const emits = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'close'): void
}>()

const { t } = useI18n()

const dialogVisible = computed<boolean>({
  get() {
    return props.visible
  },
  set(visible: boolean) {
    emits('update:visible', visible)
    if (!visible)
      emits('close')
  },
})

interface RuleForm {
  name: string
  remark: string
  status: string
  resource_ids: string[]
}

const formRef = ref<FormInstance>()
const loading = ref(false)
const resources = ref<Resource[]>([])
const form = ref<RuleForm>({ name: '', remark: '', status: '1', resource_ids: [] })

const rules = reactive<FormRules<RuleForm>>({
  name: [{ required: true, message: t('account.roleNameRequired'), trigger: 'blur' }],
  resource_ids: [{ required: true, message: t('account.rolePermissionRequired'), trigger: 'change' }],
})

// 接口按 /api/v1/<group>/<action> 的首段分组；分组名直接用该段，避免为分组维护额外元数据。
function groupOf(path: string) {
  return path.split('/').filter(Boolean)[2] || 'other'
}

const groups = computed(() => {
  const map = new Map<string, Resource[]>()
  resources.value.forEach((item) => {
    const key = groupOf(item.path)
    const list = map.get(key) || []
    list.push(item)
    map.set(key, list)
  })
  return [...map.entries()].map(([name, items]) => ({ name, items }))
})

const selectedCount = computed(() => form.value.resource_ids.length)
const allSelected = computed(() => resources.value.length > 0 && selectedCount.value === resources.value.length)

function toggleAll() {
  form.value.resource_ids = allSelected.value ? [] : resources.value.map(item => item.id)
}

function clearAll() {
  form.value.resource_ids = []
}

async function loadResources() {
  loading.value = true
  try {
    const { data } = await queryResource()
    resources.value = data.data || []
  }
  finally {
    loading.value = false
  }
}

function resetForm() {
  form.value = { name: '', remark: '', status: '1', resource_ids: [] }
  formRef.value?.clearValidate()
}

async function confirmCreate(formEl: FormInstance | undefined) {
  if (!formEl)
    return
  await formEl.validate((valid) => {
    if (!valid)
      return
    loading.value = true
    createRole({
      name: form.value.name,
      remark: form.value.remark,
      status: Number(form.value.status),
      resource_ids: form.value.resource_ids,
    }).then(() => {
      success(t('account.roleCreateSuccess'))
      dialogVisible.value = false
      resetForm()
      props.update && props.update()
    }).finally(() => {
      loading.value = false
    })
  })
}

onMounted(() => {
  void loadResources()
})
</script>

<template>
    <el-dialog v-model="dialogVisible" :title="t('account.roleCreateTitle')" width="560px">
        <template #header>
            <DialogHeader :title="t('account.roleCreateTitle')" :icon="IconShieldPlus" />
        </template>
        <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
            <p class="am-role-create__tip">
                {{ t('account.roleCreateTip') }}
            </p>
            <el-form-item :label="t('account.roleNameLabel')" prop="name">
                <el-input v-model="form.name" :placeholder="t('account.roleNamePlaceholder')" />
            </el-form-item>
            <el-form-item :label="t('account.roleRemarkLabel')" prop="remark">
                <el-input v-model="form.remark" :placeholder="t('account.roleRemarkPlaceholder')" />
            </el-form-item>
            <el-form-item :label="t('account.rolePermissionLabel')" prop="resource_ids">
                <div class="am-role-create__perm">
                    <div class="am-role-create__perm-head">
                        <span class="am-role-create__perm-count">
                            {{ t('account.rolePermissionSelected', { selected: selectedCount, total: resources.length }) }}
                        </span>
                        <button type="button" class="am-role-create__link" @click="toggleAll">
                            {{ t('account.rolePermissionSelectAll') }}
                        </button>
                        <button type="button" class="am-role-create__link am-role-create__link--muted" @click="clearAll">
                            {{ t('account.rolePermissionClear') }}
                        </button>
                    </div>
                    <div v-if="groups.length" v-loading="loading" class="am-role-create__perm-list">
                        <el-checkbox-group v-model="form.resource_ids">
                            <template v-for="group in groups" :key="group.name">
                                <div class="am-role-create__group">
                                    <span class="am-role-create__group-name">{{ group.name }}</span>
                                    <span class="am-role-create__group-count">{{ group.items.length }}</span>
                                </div>
                                <label v-for="item in group.items" :key="item.id" class="am-role-create__option">
                                    <el-checkbox :value="item.id">
                                        <span class="am-role-create__option-name">{{ item.name }}</span>
                                        <span class="am-role-create__option-path">{{ item.path }}</span>
                                    </el-checkbox>
                                    <span class="am-role-create__method" :class="`am-role-create__method--${item.method.toLowerCase()}`">
                                        {{ item.method }}
                                    </span>
                                </label>
                            </template>
                        </el-checkbox-group>
                    </div>
                    <p v-else-if="!loading" class="am-role-create__empty">
                        {{ t('account.rolePermissionEmpty') }}
                    </p>
                </div>
            </el-form-item>
            <el-form-item :label="t('account.roleStatusLabel')" prop="status">
                <div class="am-role-create__status">
                    <el-switch v-model="form.status" active-value="1" inactive-value="2" />
                    <span class="am-role-create__status-hint">{{ t('account.roleStatusHint') }}</span>
                </div>
            </el-form-item>
        </el-form>

        <template #footer>
            <span>
                <el-button @click="dialogVisible = false">{{ t('account.cancel') }}</el-button>
                <el-button type="primary" :loading="loading" @click="confirmCreate(formRef)">
                    {{ t('account.newRole') }}
                </el-button>
            </span>
        </template>
    </el-dialog>
</template>

<style scoped lang="scss">
.am-role-create {
  &__tip {
    margin: 0 0 var(--am-spacing-md);
    color: var(--am-foreground-secondary);
    font-size: var(--am-font-xs);
    line-height: 1.6;
  }

  &__perm {
    width: 100%;
  }

  &__perm-head {
    display: flex;
    align-items: center;
    gap: var(--am-spacing-md);
    margin-bottom: var(--am-spacing-sm);
  }

  &__perm-count {
    color: var(--am-foreground-muted);
    font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
    font-size: var(--am-font-xs);
  }

  &__link {
    padding: 0;
    border: none;
    background: none;
    color: var(--am-accent-primary);
    font-size: var(--am-font-xs);
    cursor: pointer;

    &--muted {
      color: var(--am-foreground-secondary);
    }
  }

  &__perm-list {
    max-height: 268px;
    overflow-y: auto;
    border: 1px solid var(--am-border-subtle);
    border-radius: var(--am-radius-sm);
    background: var(--am-surface-card);
  }

  &__group {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 6px 12px;
    background: var(--am-surface-secondary);
    color: var(--am-foreground-secondary);
    font-size: var(--am-font-xs);
    font-weight: 600;
  }

  &__group-name {
    font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  }

  &__group-count {
    color: var(--am-foreground-muted);
    font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
    font-weight: 400;
  }

  &__option {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--am-spacing-md);
    padding: 6px 12px;
    cursor: pointer;
  }

  &__option-name {
    color: var(--am-foreground-primary);
    font-size: var(--am-font-sm);
  }

  &__option-path {
    margin-left: var(--am-spacing-sm);
    color: var(--am-foreground-muted);
    font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
    font-size: var(--am-font-xs);
  }

  &__method {
    padding: 1px 6px;
    border-radius: var(--am-radius-xs, 4px);
    font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
    font-size: var(--am-font-xs);
    font-weight: 600;

    &--get {
      background: var(--am-accent-success-soft);
      color: var(--am-accent-success);
    }

    &--post {
      background: var(--am-accent-warning-soft);
      color: var(--am-accent-warning);
    }
  }

  &__empty {
    margin: 0;
    padding: var(--am-spacing-md);
    color: var(--am-foreground-muted);
    font-size: var(--am-font-xs);
  }

  &__status {
    display: flex;
    align-items: center;
    gap: var(--am-spacing-md);
  }

  &__status-hint {
    color: var(--am-foreground-muted);
    font-size: var(--am-font-xs);
  }
}
</style>
