<template>
  <AppLayout>
    <div class="mx-auto w-full max-w-5xl px-4 py-6 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="mb-6 flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ t('admin.modelPlaza.title') }}</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('admin.modelPlaza.description') }}</p>
        </div>
        <button class="btn btn-primary" :disabled="saving || loading" @click="save">
          <Icon v-if="saving" name="refresh" size="sm" class="mr-1 animate-spin" />
          {{ t('common.save') }}
        </button>
      </div>

      <div v-if="loading" class="flex items-center justify-center py-24 text-gray-400">
        <Icon name="refresh" size="xl" class="animate-spin" />
      </div>

      <template v-else>
        <!-- Master switch -->
        <div class="mb-6 flex items-center justify-between rounded-xl border border-gray-200 bg-white p-5 dark:border-dark-700 dark:bg-dark-800">
          <div>
            <div class="font-medium text-gray-900 dark:text-gray-100">{{ t('admin.modelPlaza.enableLabel') }}</div>
            <div class="mt-0.5 text-sm text-gray-500 dark:text-gray-400">{{ t('admin.modelPlaza.enableHint') }}</div>
          </div>
          <button
            type="button"
            role="switch"
            :aria-checked="enabled"
            class="relative inline-flex h-6 w-11 flex-shrink-0 items-center rounded-full transition"
            :class="enabled ? 'bg-primary-500' : 'bg-gray-300 dark:bg-dark-600'"
            @click="enabled = !enabled"
          >
            <span class="inline-block h-4 w-4 transform rounded-full bg-white transition" :class="enabled ? 'translate-x-6' : 'translate-x-1'" />
          </button>
        </div>

        <!-- Category selector (fixed 5) -->
        <div class="mb-5 flex flex-wrap gap-2">
          <button
            v-for="c in categories"
            :key="c.key"
            type="button"
            class="inline-flex items-center gap-2 rounded-xl border px-4 py-2 text-sm font-semibold transition-all"
            :class="c.key === activeCat
              ? 'border-primary-500 bg-primary-50 text-primary-600 dark:border-primary-600 dark:bg-primary-500/10 dark:text-primary-400'
              : 'border-gray-200 text-gray-600 hover:border-gray-300 dark:border-dark-700 dark:text-gray-300 dark:hover:border-dark-600'"
            @click="activeCat = c.key"
          >
            <span>{{ categoryLabel(c.key) }}</span>
            <span class="rounded-full bg-gray-100 px-1.5 text-[0.6875rem] font-semibold text-gray-500 dark:bg-dark-700 dark:text-gray-400">{{ c.groups.length }}</span>
          </button>
        </div>

        <!-- Active category's groups -->
        <div v-if="activeCategory" class="space-y-5">
          <div
            v-for="(group, gi) in activeCategory.groups"
            :key="gi"
            class="rounded-xl border border-gray-200 bg-white p-5 dark:border-dark-700 dark:bg-dark-800"
          >
            <!-- group fields -->
            <div class="mb-4 flex flex-wrap items-end gap-3">
              <div class="min-w-[180px] flex-1">
                <label class="mb-1 block text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.modelPlaza.groupName') }}</label>
                <input v-model="group.name" class="input" :placeholder="t('admin.modelPlaza.groupNamePlaceholder')" />
              </div>
              <div class="w-28">
                <label class="mb-1 block text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.modelPlaza.currency') }}</label>
                <input v-model="group.currency" class="input" placeholder="$" />
              </div>
              <button class="btn btn-danger" :title="t('admin.modelPlaza.removeGroup')" @click="removeGroup(gi)">
                <Icon name="trash" size="sm" />
              </button>
            </div>
            <div class="mb-4">
              <label class="mb-1 block text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.modelPlaza.groupDescription') }}</label>
              <input v-model="group.description" class="input" :placeholder="t('admin.modelPlaza.groupDescriptionPlaceholder')" />
            </div>

            <!-- models table -->
            <div class="overflow-x-auto rounded-lg border border-gray-100 dark:border-dark-700">
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b border-gray-100 bg-gray-50 text-left text-xs text-gray-500 dark:border-dark-700 dark:bg-dark-700/40 dark:text-gray-400">
                    <th class="px-3 py-2 font-medium">{{ t('modelPlaza.columns.model') }}</th>
                    <th class="px-3 py-2 font-medium">{{ t('modelPlaza.columns.input') }}</th>
                    <th class="px-3 py-2 font-medium">{{ t('modelPlaza.columns.output') }}</th>
                    <th class="px-3 py-2 font-medium">{{ t('modelPlaza.columns.cacheRead') }}</th>
                    <th class="px-3 py-2 font-medium">{{ t('modelPlaza.columns.cacheWrite') }}</th>
                    <th class="px-3 py-2 font-medium">{{ t('admin.modelPlaza.note') }}</th>
                    <th class="w-10 px-3 py-2"></th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(model, mi) in group.models" :key="mi" class="border-b border-gray-50 dark:border-dark-700/60">
                    <td class="px-2 py-1.5"><input v-model="model.name" class="input !py-1 text-sm" /></td>
                    <td class="px-2 py-1.5"><input v-model="model.price.input" class="input !py-1 text-sm" /></td>
                    <td class="px-2 py-1.5"><input v-model="model.price.output" class="input !py-1 text-sm" /></td>
                    <td class="px-2 py-1.5"><input v-model="model.price.cache_read" class="input !py-1 text-sm" /></td>
                    <td class="px-2 py-1.5"><input v-model="model.price.cache_write" class="input !py-1 text-sm" /></td>
                    <td class="px-2 py-1.5"><input v-model="model.note" class="input !py-1 text-sm" /></td>
                    <td class="px-2 py-1.5 text-center">
                      <button class="text-gray-400 hover:text-red-500" :title="t('admin.modelPlaza.removeModel')" @click="removeModel(group, mi)">
                        <Icon name="trash" size="sm" />
                      </button>
                    </td>
                  </tr>
                  <tr v-if="!group.models.length">
                    <td colspan="7" class="px-3 py-4 text-center text-xs text-gray-400">{{ t('admin.modelPlaza.noModels') }}</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <button class="btn btn-secondary mt-3 btn-sm" @click="addModel(group)">
              <Icon name="plus" size="sm" class="mr-1" />{{ t('admin.modelPlaza.addModel') }}
            </button>
          </div>

          <div v-if="!activeCategory.groups.length" class="rounded-xl border border-dashed border-gray-200 py-12 text-center text-sm text-gray-400 dark:border-dark-700">
            {{ t('admin.modelPlaza.categoryEmpty', { name: categoryLabel(activeCat) }) }}
          </div>

          <button class="btn btn-secondary" @click="addGroup">
            <Icon name="plus" size="sm" class="mr-1" />{{ t('admin.modelPlaza.addGroup') }}
          </button>
        </div>

        <p class="mt-4 text-xs text-gray-400">{{ t('modelPlaza.unitHint') }}</p>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import {
  MODEL_PLAZA_CATEGORY_KEYS,
  getModelPlazaAdminConfig,
  type ModelPlazaCategory,
  type ModelPlazaCategoryKey,
  type ModelPlazaGroup,
  type ModelPlazaModel,
  updateModelPlazaAdminConfig,
} from '@/api/modelPlaza'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const appStore = useAppStore()

const enabled = ref(false)
const categories = ref<ModelPlazaCategory[]>([])
const activeCat = ref<ModelPlazaCategoryKey>('claude')
const loading = ref(false)
const saving = ref(false)

const activeCategory = computed(() => categories.value.find((c) => c.key === activeCat.value))

const CATEGORY_LABEL_KEYS: Record<ModelPlazaCategoryKey, string> = {
  claude: 'modelPlaza.categories.claude',
  openai: 'modelPlaza.categories.openai',
  gemini: 'modelPlaza.categories.gemini',
  domestic: 'modelPlaza.categories.domestic',
  image: 'modelPlaza.categories.image',
}
function categoryLabel(key: ModelPlazaCategoryKey): string {
  return t(CATEGORY_LABEL_KEYS[key])
}

function emptyModel(): ModelPlazaModel {
  return { name: '', price: { input: '', output: '', cache_read: '', cache_write: '' }, note: '' }
}
function addGroup() {
  activeCategory.value?.groups.push({ name: '', description: '', currency: '$', models: [emptyModel()] })
}
function removeGroup(index: number) {
  activeCategory.value?.groups.splice(index, 1)
}
function addModel(group: ModelPlazaGroup) {
  group.models.push(emptyModel())
}
function removeModel(group: ModelPlazaGroup, index: number) {
  group.models.splice(index, 1)
}

/** Build the fixed 5 categories in order, carrying over loaded groups (fully normalized). */
function normalize(loaded: ModelPlazaCategory[]): ModelPlazaCategory[] {
  const byKey = new Map(loaded.map((c) => [c.key, c]))
  return MODEL_PLAZA_CATEGORY_KEYS.map((key) => {
    const src = byKey.get(key)
    return {
      key,
      groups: (src?.groups ?? []).map((g) => ({
        name: g.name ?? '',
        description: g.description ?? '',
        currency: g.currency ?? '',
        models: (g.models ?? []).map((m) => ({
          name: m.name ?? '',
          note: m.note ?? '',
          price: {
            input: m.price?.input ?? '',
            output: m.price?.output ?? '',
            cache_read: m.price?.cache_read ?? '',
            cache_write: m.price?.cache_write ?? '',
          },
        })),
      })),
    }
  })
}

async function load() {
  loading.value = true
  try {
    const res = await getModelPlazaAdminConfig()
    enabled.value = res.enabled
    categories.value = normalize(res.config?.categories ?? [])
  } catch {
    appStore.showError(t('admin.modelPlaza.loadFailed'))
  } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  try {
    const res = await updateModelPlazaAdminConfig({
      enabled: enabled.value,
      config: { categories: categories.value },
    })
    enabled.value = res.enabled
    categories.value = normalize(res.config?.categories ?? [])
    appStore.showSuccess(t('admin.modelPlaza.saved'))
    window.dispatchEvent(new CustomEvent('model-plaza-config-changed'))
  } catch (error) {
    const err = error as { message?: string }
    appStore.showError(err?.message || t('admin.modelPlaza.saveFailed'))
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>
