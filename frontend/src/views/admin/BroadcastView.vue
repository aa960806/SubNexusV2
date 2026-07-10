<template>
  <AppLayout>
    <div class="space-y-6">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">跑马灯管理</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">管理全局滚动公告及展示时间</p>
        </div>
        <label class="inline-flex items-center gap-3 text-sm text-gray-700 dark:text-dark-200">
          <input v-model="config.enabled" type="checkbox" class="rounded border-gray-300 text-primary-600 focus:ring-primary-500" @change="saveConfig" />
          开启跑马灯
        </label>
      </div>

      <section class="grid gap-6 xl:grid-cols-[380px_1fr]">
        <form class="card space-y-4 p-5" @submit.prevent="saveBroadcast">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ editingId ? '编辑公告' : '新增公告' }}</h2>
          <div>
            <label class="input-label">标题</label>
            <input v-model="form.title" class="input" maxlength="120" placeholder="系统公告" />
          </div>
          <div>
            <label class="input-label">内容</label>
            <textarea v-model="form.content" class="input min-h-28 resize-y" maxlength="2000" required placeholder="请输入滚动公告内容"></textarea>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="input-label">优先级</label>
              <input v-model.number="form.priority" class="input" min="0" max="1000" type="number" />
            </div>
            <label class="mt-7 inline-flex items-center gap-2 text-sm text-gray-700 dark:text-dark-200">
              <input v-model="form.enabled" type="checkbox" class="rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
              启用
            </label>
          </div>
          <div>
            <label class="input-label">开始时间</label>
            <input v-model="form.start_at" class="input" type="datetime-local" />
          </div>
          <div>
            <label class="input-label">结束时间</label>
            <input v-model="form.end_at" class="input" type="datetime-local" />
          </div>
          <div class="flex gap-3">
            <button class="btn btn-primary flex-1" :disabled="saving || !form.content.trim()">
              {{ saving ? '保存中...' : editingId ? '更新' : '创建' }}
            </button>
            <button v-if="editingId" type="button" class="btn btn-secondary" @click="resetForm">取消</button>
          </div>
        </form>

        <div class="card overflow-hidden p-0">
          <div class="flex items-center justify-between border-b border-gray-200 px-5 py-4 dark:border-dark-700">
            <h2 class="font-semibold text-gray-900 dark:text-white">公告列表</h2>
            <div class="flex gap-2">
              <button class="btn btn-secondary btn-sm" :disabled="cleaning" @click="cleanup">清理过期系统消息</button>
              <button class="btn btn-secondary btn-sm" :disabled="loading" @click="load">刷新</button>
            </div>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full min-w-[760px] divide-y divide-gray-200 dark:divide-dark-700">
              <thead class="bg-gray-50 dark:bg-dark-800">
                <tr>
                  <th class="table-th">标题 / 内容</th>
                  <th class="table-th">来源</th>
                  <th class="table-th">优先级</th>
                  <th class="table-th">状态</th>
                  <th class="table-th text-right">操作</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
                <tr v-if="loading"><td colspan="5" class="px-5 py-10 text-center text-sm text-gray-400">加载中...</td></tr>
                <tr v-else-if="items.length === 0"><td colspan="5" class="px-5 py-10 text-center text-sm text-gray-400">暂无公告</td></tr>
                <template v-else>
                  <tr v-for="item in items" :key="item.id">
                    <td class="table-td max-w-md">
                      <p class="font-medium text-gray-900 dark:text-white">{{ item.title || '无标题' }}</p>
                      <p class="mt-1 truncate text-xs text-gray-500">{{ item.content }}</p>
                    </td>
                    <td class="table-td">{{ item.source === 'admin' ? '管理员' : '历史系统消息' }}</td>
                    <td class="table-td">{{ item.priority }}</td>
                    <td class="table-td"><span :class="item.enabled ? 'badge badge-success' : 'badge badge-secondary'">{{ item.enabled ? '启用' : '停用' }}</span></td>
                    <td class="table-td text-right">
                      <button class="btn-icon mr-2" @click="edit(item)"><Icon name="edit" size="sm" /></button>
                      <button class="btn-icon text-red-500" @click="remove(item.id)"><Icon name="trash" size="sm" /></button>
                    </td>
                  </tr>
                </template>
              </tbody>
            </table>
          </div>
        </div>
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
  cleanupExpiredBroadcasts,
  createBroadcast,
  deleteBroadcast,
  getBroadcastConfig,
  listBroadcasts,
  updateBroadcast,
  updateBroadcastConfig,
  type ActivityBroadcast,
} from '@/api/broadcast'

const appStore = useAppStore()
const config = reactive({ enabled: true })
const items = ref<ActivityBroadcast[]>([])
const editingId = ref<number | null>(null)
const loading = ref(false)
const saving = ref(false)
const cleaning = ref(false)
const form = reactive({ title: '', content: '', enabled: true, priority: 10, start_at: '', end_at: '' })

function toAPI(value: string) {
  return value ? new Date(value).toISOString() : null
}

function fromAPI(value?: string | null) {
  if (!value) return ''
  const date = new Date(value)
  return new Date(date.getTime() - date.getTimezoneOffset() * 60000).toISOString().slice(0, 16)
}

async function load() {
  loading.value = true
  try {
    const [savedConfig, broadcasts] = await Promise.all([getBroadcastConfig(), listBroadcasts()])
    config.enabled = savedConfig.enabled
    items.value = broadcasts
  } catch (error: any) {
    appStore.showError(error?.message || '加载跑马灯配置失败')
  } finally {
    loading.value = false
  }
}

async function saveConfig() {
  try {
    await updateBroadcastConfig({ enabled: config.enabled })
    appStore.showSuccess('跑马灯开关已保存')
  } catch (error: any) {
    appStore.showError(error?.message || '保存失败')
    await load()
  }
}

async function saveBroadcast() {
  saving.value = true
  try {
    const payload = { title: form.title, content: form.content, enabled: form.enabled, priority: form.priority, start_at: toAPI(form.start_at), end_at: toAPI(form.end_at) }
    if (editingId.value) await updateBroadcast(editingId.value, payload)
    else await createBroadcast(payload)
    appStore.showSuccess(editingId.value ? '公告已更新' : '公告已创建')
    resetForm()
    await load()
  } catch (error: any) {
    appStore.showError(error?.message || '保存公告失败')
  } finally {
    saving.value = false
  }
}

function edit(item: ActivityBroadcast) {
  editingId.value = item.id
  form.title = item.title
  form.content = item.content
  form.enabled = item.enabled
  form.priority = item.priority
  form.start_at = fromAPI(item.start_at)
  form.end_at = fromAPI(item.end_at)
}

function resetForm() {
  editingId.value = null
  Object.assign(form, { title: '', content: '', enabled: true, priority: 10, start_at: '', end_at: '' })
}

async function remove(id: number) {
  if (!window.confirm('确定删除这条公告吗？')) return
  try {
    await deleteBroadcast(id)
    await load()
  } catch (error: any) {
    appStore.showError(error?.message || '删除失败')
  }
}

async function cleanup() {
  cleaning.value = true
  try {
    const deleted = await cleanupExpiredBroadcasts(7)
    appStore.showSuccess(`已清理 ${deleted} 条过期系统消息`)
    await load()
  } catch (error: any) {
    appStore.showError(error?.message || '清理失败')
  } finally {
    cleaning.value = false
  }
}

onMounted(load)
</script>
