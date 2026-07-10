<template>
  <AppLayout>
    <div class="space-y-6">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">活动入口管理</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">统一维护用户活动中心里的活动入口</p>
        </div>
        <label class="inline-flex items-center gap-3 text-sm text-gray-700 dark:text-dark-200">
          <input v-model="config.enabled" type="checkbox" class="rounded border-gray-300 text-primary-600 focus:ring-primary-500" @change="saveConfig" />
          开启用户活动中心入口
        </label>
      </div>

      <section class="grid gap-6 xl:grid-cols-[380px_1fr]">
        <form class="card space-y-4 p-5" @submit.prevent="saveItem">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ editingId ? '编辑活动' : '新增活动' }}</h2>
          <div>
            <label class="input-label">唯一标识</label>
            <input v-model="form.slug" class="input" placeholder="leaderboard-week" />
          </div>
          <div>
            <label class="input-label">活动标题</label>
            <input v-model="form.title" class="input" placeholder="周榜活动" />
          </div>
          <div>
            <label class="input-label">副标题</label>
            <input v-model="form.subtitle" class="input" placeholder="冲榜赢奖励" />
          </div>
          <div>
            <label class="input-label">描述</label>
            <textarea v-model="form.description" class="input min-h-24 resize-y" placeholder="活动说明"></textarea>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="input-label">图标</label>
              <input v-model="form.icon" class="input" placeholder="gift" />
            </div>
            <div>
              <label class="input-label">排序</label>
              <input v-model.number="form.sort_order" class="input" type="number" step="1" />
            </div>
          </div>
          <div>
            <label class="input-label">站内路径</label>
            <input v-model="form.route_path" class="input" placeholder="/activities/summer" />
          </div>
          <div>
            <label class="input-label">外部链接</label>
            <input v-model="form.external_url" class="input" placeholder="https://..." />
          </div>
          <div>
            <label class="input-label">封面图</label>
            <input v-model="form.cover_url" class="input" placeholder="https://..." />
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="input-label">开始时间</label>
              <input v-model="form.start_at" class="input" type="datetime-local" />
            </div>
            <div>
              <label class="input-label">结束时间</label>
              <input v-model="form.end_at" class="input" type="datetime-local" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="input-label">按钮文字</label>
              <input v-model="form.action_label" class="input" placeholder="查看" />
            </div>
            <label class="mt-7 inline-flex items-center gap-2 text-sm text-gray-700 dark:text-dark-200">
              <input v-model="form.enabled" type="checkbox" class="rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
              上架
            </label>
          </div>
          <div class="flex gap-3">
            <button class="btn btn-primary flex-1" :disabled="saving || !form.slug.trim() || !form.title.trim()">
              {{ saving ? '保存中...' : '保存活动' }}
            </button>
            <button v-if="editingId" type="button" class="btn btn-secondary" @click="resetForm">取消</button>
          </div>
        </form>

        <section class="card overflow-hidden">
          <div class="flex items-center justify-between border-b border-gray-200 px-5 py-4 dark:border-dark-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">活动列表</h2>
            <button class="btn btn-secondary" :disabled="loading" @click="loadData">
              <Icon name="refresh" size="sm" :class="loading ? 'animate-spin' : ''" />
            </button>
          </div>
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200 dark:divide-dark-700">
              <thead class="bg-gray-50 dark:bg-dark-800/80">
                <tr>
                  <th class="table-th">活动</th>
                  <th class="table-th">路径</th>
                  <th class="table-th">状态</th>
                  <th class="table-th text-right">排序</th>
                  <th class="table-th text-right">操作</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
                <tr v-if="loading">
                  <td colspan="5" class="px-5 py-12 text-center text-sm text-gray-500 dark:text-dark-400">加载中...</td>
                </tr>
                <tr v-else-if="items.length === 0">
                  <td colspan="5" class="px-5 py-12 text-center text-sm text-gray-500 dark:text-dark-400">暂无活动</td>
                </tr>
                <tr v-for="item in items" :key="item.id" class="hover:bg-gray-50 dark:hover:bg-dark-800/60">
                  <td class="table-td">
                    <div class="font-medium text-gray-900 dark:text-white">{{ item.title }}</div>
                    <div class="text-xs text-gray-500 dark:text-dark-400">{{ item.slug }}</div>
                  </td>
                  <td class="table-td">{{ item.route_path || item.external_url || '-' }}</td>
                  <td class="table-td">
                    <span :class="item.enabled ? 'badge badge-success' : 'badge badge-gray'">{{ item.enabled ? '上架' : '下架' }}</span>
                  </td>
                  <td class="table-td text-right">{{ item.sort_order }}</td>
                  <td class="table-td text-right">
                    <button class="btn-icon mr-2" @click="editItem(item)"><Icon name="edit" size="sm" /></button>
                    <button class="btn-icon text-red-500" @click="removeItem(item.id)"><Icon name="trash" size="sm" /></button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </section>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore } from '@/stores'
import {
  createActivityCenterItem,
  deleteActivityCenterItem,
  getActivityCenterConfig,
  listAdminActivityCenterItems,
  updateActivityCenterConfig,
  updateActivityCenterItem,
  type ActivityCenterItem,
  type ActivityCenterItemInput,
} from '@/api/activityCenter'

const appStore = useAppStore()
const loading = ref(false)
const saving = ref(false)
const editingId = ref<number | null>(null)
const items = ref<ActivityCenterItem[]>([])
const config = reactive({ enabled: false })

const form = reactive<ActivityCenterItemInput>({
  slug: '',
  title: '',
  subtitle: '',
  description: '',
  icon: 'gift',
  cover_url: '',
  route_path: '',
  external_url: '',
  action_label: '查看',
  activity_type: 'custom',
  enabled: true,
  sort_order: 0,
  start_at: '',
  end_at: '',
  metadata: {},
})

async function loadData() {
  loading.value = true
  try {
    Object.assign(config, await getActivityCenterConfig())
    items.value = await listAdminActivityCenterItems()
  } catch (error: any) {
    appStore.showError(error?.message || '加载活动入口失败')
  } finally {
    loading.value = false
  }
}

async function saveConfig() {
  try {
    Object.assign(config, await updateActivityCenterConfig({ enabled: config.enabled }))
    window.dispatchEvent(new CustomEvent('activity-center-config-changed'))
    appStore.showSuccess('活动中心开关已保存')
  } catch (error: any) {
    appStore.showError(error?.message || '保存开关失败')
  }
}

async function saveItem() {
  saving.value = true
  try {
    const payload = buildPayload()
    if (editingId.value) {
      await updateActivityCenterItem(editingId.value, payload)
      appStore.showSuccess('活动已更新')
    } else {
      await createActivityCenterItem(payload)
      appStore.showSuccess('活动已创建')
    }
    resetForm()
    await loadData()
  } catch (error: any) {
    appStore.showError(error?.message || '保存活动失败')
  } finally {
    saving.value = false
  }
}

function editItem(item: ActivityCenterItem) {
  editingId.value = item.id
  form.slug = item.slug
  form.title = item.title
  form.subtitle = item.subtitle || ''
  form.description = item.description || ''
  form.icon = item.icon || 'gift'
  form.cover_url = item.cover_url || ''
  form.route_path = item.route_path || ''
  form.external_url = item.external_url || ''
  form.action_label = item.action_label || '查看'
  form.activity_type = item.activity_type || 'custom'
  form.enabled = item.enabled
  form.sort_order = item.sort_order || 0
  form.start_at = fromApiDateTime(item.start_at)
  form.end_at = fromApiDateTime(item.end_at)
  form.metadata = item.metadata || {}
}

async function removeItem(id: number) {
  try {
    await deleteActivityCenterItem(id)
    appStore.showSuccess('活动已删除')
    await loadData()
  } catch (error: any) {
    appStore.showError(error?.message || '删除活动失败')
  }
}

function resetForm() {
  editingId.value = null
  form.slug = ''
  form.title = ''
  form.subtitle = ''
  form.description = ''
  form.icon = 'gift'
  form.cover_url = ''
  form.route_path = ''
  form.external_url = ''
  form.action_label = '查看'
  form.activity_type = 'custom'
  form.enabled = true
  form.sort_order = 0
  form.start_at = ''
  form.end_at = ''
  form.metadata = {}
}

function buildPayload(): ActivityCenterItemInput {
  return {
    ...form,
    slug: form.slug.trim(),
    title: form.title.trim(),
    start_at: toApiDateTime(form.start_at || ''),
    end_at: toApiDateTime(form.end_at || ''),
    metadata: form.metadata || {},
  }
}

function toApiDateTime(value: string) {
  return value ? new Date(value).toISOString() : null
}

function fromApiDateTime(value?: string) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return new Date(date.getTime() - date.getTimezoneOffset() * 60000).toISOString().slice(0, 16)
}

onMounted(loadData)
</script>

<style scoped>
.table-th {
  @apply whitespace-nowrap px-5 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400;
}
.table-td {
  @apply whitespace-nowrap px-5 py-4 text-sm text-gray-600 dark:text-dark-300;
}
</style>
