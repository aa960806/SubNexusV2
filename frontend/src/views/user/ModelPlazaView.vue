<template>
  <AppLayout>
    <div class="mx-auto w-full max-w-6xl px-5 py-14 sm:px-8 sm:py-20">
      <!-- Hero -->
      <div class="mb-12 text-center">
        <h1 class="text-4xl font-bold tracking-tight text-gray-900 dark:text-gray-50 sm:text-5xl">
          {{ t('modelPlaza.title') }}
        </h1>
        <p class="mx-auto mt-4 max-w-2xl text-base text-gray-400 dark:text-gray-500 sm:text-lg">{{ t('modelPlaza.subtitle') }}</p>
        <div class="mt-6 inline-flex items-center gap-2 rounded-full bg-primary-50 px-4 py-2 text-sm font-medium text-primary-600 dark:bg-primary-500/10 dark:text-primary-400">
          <svg class="h-4 w-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.9">
            <path stroke-linecap="round" stroke-linejoin="round" d="M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z" />
          </svg>
          <span>{{ t('modelPlaza.pricingNote') }}</span>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex flex-col items-center justify-center py-24 text-gray-400">
        <span class="h-6 w-6 animate-spin rounded-full border-2 border-current border-t-transparent"></span>
      </div>

      <!-- Empty (no category has any group) -->
      <div
        v-else-if="!visibleCategories.length"
        class="rounded-3xl border border-gray-100 bg-white py-20 text-center dark:border-dark-700/50 dark:bg-dark-800/40"
      >
        <p class="text-[0.9375rem] text-gray-400 dark:text-gray-500">{{ t('modelPlaza.empty') }}</p>
      </div>

      <template v-else>
        <!-- Level 1: model type tabs (minimal underline) -->
        <div class="mb-8 flex flex-wrap justify-center gap-1 border-b border-gray-100 dark:border-dark-700/60">
          <button
            v-for="c in visibleCategories"
            :key="c.key"
            type="button"
            class="relative inline-flex items-center gap-2 px-5 pb-4 pt-3 text-base font-semibold transition-colors"
            :class="c.key === activeCat
              ? 'text-primary-600 dark:text-primary-400'
              : 'text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300'"
            @click="selectCat(c.key)"
          >
            <span class="h-1.5 w-1.5 rounded-full bg-current" :class="c.key === activeCat ? 'opacity-100' : 'opacity-30'"></span>
            <span>{{ categoryLabel(c.key) }}</span>
            <span
              v-if="c.key === activeCat"
              class="absolute inset-x-2 -bottom-px h-0.5 rounded-sm bg-primary-500 dark:bg-primary-400"
            ></span>
          </button>
        </div>

        <!-- Level 2: group pills (horizontal scroll when many) -->
        <div class="relative mb-8">
          <div
            ref="groupsBar"
            class="flex gap-2 overflow-x-auto py-0.5 scrollbar-hide"
            :class="groupsCentered ? 'justify-center' : 'justify-start'"
            @scroll="updateEdges"
          >
            <button
              v-for="(g, gi) in activeGroups"
              :key="gi"
              type="button"
              class="inline-flex flex-shrink-0 items-center gap-2 rounded-full border px-3.5 py-1.5 text-[0.8125rem] font-medium transition-all"
              :class="gi === activeGroupIdx
                ? 'border-primary-500 bg-primary-500 text-white'
                : 'border-gray-200 text-gray-500 hover:border-gray-300 hover:text-gray-700 dark:border-dark-700 dark:text-gray-400 dark:hover:border-dark-600 dark:hover:text-gray-200'"
              @click="activeGroupIdx = gi"
            >
              <span>{{ g.name || t('modelPlaza.untitledGroup') }}</span>
              <span class="text-[0.6875rem] font-semibold" :class="gi === activeGroupIdx ? 'text-white/70' : 'text-gray-400 dark:text-gray-500'">{{ g.models?.length || 0 }}</span>
            </button>
          </div>
          <!-- fade edges -->
          <div v-show="edgeL" class="pointer-events-none absolute inset-y-0 left-0 w-10 bg-gradient-to-r from-gray-50 to-transparent dark:from-dark-950"></div>
          <div v-show="edgeR" class="pointer-events-none absolute inset-y-0 right-0 w-10 bg-gradient-to-l from-gray-50 to-transparent dark:from-dark-950"></div>
        </div>

        <!-- Price card -->
        <div
          v-if="activeGroup"
          :key="activeCat + ':' + activeGroupIdx"
          class="animate-fade-in overflow-hidden rounded-[1.75rem] border border-gray-100 bg-white shadow-card dark:border-dark-700/50 dark:bg-dark-800/50"
        >
          <!-- header -->
          <div class="border-b border-gray-50 px-8 pb-6 pt-8 dark:border-dark-700/50 sm:px-10">
            <div class="flex items-start justify-between gap-4">
              <div>
                <h2 class="text-2xl font-bold tracking-tight text-gray-900 dark:text-gray-50">
                  {{ activeGroup.name || t('modelPlaza.untitledGroup') }}
                </h2>
                <p v-if="activeGroup.description" class="mt-2 max-w-2xl text-[0.9375rem] leading-relaxed text-gray-400 dark:text-gray-500">
                  {{ activeGroup.description }}
                </p>
              </div>
              <span
                v-if="activeGroup.currency"
                class="inline-flex flex-shrink-0 items-center gap-1 rounded-full bg-gray-100 px-3 py-1.5 text-xs font-semibold text-gray-500 dark:bg-dark-700 dark:text-gray-300"
              >
                {{ t('modelPlaza.currencyLabel', { currency: activeGroup.currency }) }}
              </span>
            </div>
          </div>

          <!-- model rows: name + stacked label/value price blocks -->
          <div class="divide-y divide-gray-50 dark:divide-dark-700/40">
            <div
              v-for="(m, mi) in activeGroup.models"
              :key="mi"
              class="flex flex-col gap-5 px-8 py-6 transition-colors hover:bg-gray-50/60 dark:hover:bg-dark-700/20 sm:px-10 lg:flex-row lg:items-center lg:gap-8"
            >
              <!-- model name -->
              <div class="lg:w-56 lg:flex-shrink-0">
                <span class="inline-flex flex-wrap items-center gap-2 text-base font-semibold text-gray-900 dark:text-gray-100">
                  {{ m.name || '-' }}
                  <span v-if="m.note" class="rounded-full bg-primary-50 px-2 py-0.5 text-[0.6875rem] font-semibold text-primary-600 dark:bg-primary-500/10 dark:text-primary-400">{{ m.note }}</span>
                </span>
              </div>
              <!-- price blocks: label on top, amount directly below -->
              <div class="grid flex-1 grid-cols-2 gap-x-6 gap-y-5 sm:grid-cols-4">
                <div
                  v-for="col in priceColumns"
                  :key="col.key"
                  class="flex flex-col gap-1.5"
                >
                  <div class="text-xs font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500">{{ t(col.label) }}</div>
                  <div class="text-lg" v-html="priceCell(activeGroup.currency, m.price?.[col.key])"></div>
                </div>
              </div>
            </div>
            <div v-if="!activeGroup.models?.length" class="py-14 text-center text-sm text-gray-400">{{ t('modelPlaza.noModels') }}</div>
          </div>

          <!-- footer -->
          <div class="border-t border-gray-50 px-8 py-4 text-right text-xs text-gray-300 dark:border-dark-700/50 dark:text-gray-600 sm:px-10">{{ t('modelPlaza.unitHint') }}</div>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import modelPlazaAPI, {
  MODEL_PLAZA_CATEGORY_KEYS,
  type ModelPlazaCategory,
  type ModelPlazaCategoryKey,
  type ModelPlazaGroup,
  type ModelPlazaModel,
} from '@/api/modelPlaza'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const appStore = useAppStore()

const categories = ref<ModelPlazaCategory[]>([])
const loading = ref(false)
const activeCat = ref<ModelPlazaCategoryKey | ''>('')
const activeGroupIdx = ref(0)
const groupsBar = ref<HTMLElement | null>(null)
const groupsCentered = ref(true)
const edgeL = ref(false)
const edgeR = ref(false)
let controller: AbortController | null = null

// Only categories that actually have at least one group are shown to users.
const visibleCategories = computed(() =>
  MODEL_PLAZA_CATEGORY_KEYS
    .map((key) => categories.value.find((c) => c.key === key))
    .filter((c): c is ModelPlazaCategory => !!c && (c.groups?.length ?? 0) > 0),
)

const activeGroups = computed<ModelPlazaGroup[]>(
  () => visibleCategories.value.find((c) => c.key === activeCat.value)?.groups ?? [],
)
const activeGroup = computed<ModelPlazaGroup | undefined>(() => activeGroups.value[activeGroupIdx.value])

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

/** Price columns rendered for each model (label key + price field). */
const priceColumns: { key: keyof ModelPlazaModel['price']; label: string }[] = [
  { key: 'input', label: 'modelPlaza.columns.input' },
  { key: 'output', label: 'modelPlaza.columns.output' },
  { key: 'cache_read', label: 'modelPlaza.columns.cacheRead' },
  { key: 'cache_write', label: 'modelPlaza.columns.cacheWrite' },
]

/** Render a price value as HTML; empty -> muted em dash, otherwise mono + currency prefix. */
function priceCell(currency: string | undefined, value: string | undefined): string {
  const v = (value ?? '').toString().trim()
  if (!v) return '<span class="text-gray-300 dark:text-gray-600">—</span>'
  const text = `${currency ? esc(currency) : ''}${esc(v)}`
  return `<span class="font-mono font-semibold tabular-nums text-gray-900 dark:text-gray-100">${text}</span>`
}
function esc(s: string): string {
  return s.replace(/[&<>"]/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' })[c] as string)
}

function selectCat(key: ModelPlazaCategoryKey) {
  activeCat.value = key
  activeGroupIdx.value = 0
}

function updateEdges() {
  const el = groupsBar.value
  if (!el) return
  const overflow = el.scrollWidth > el.clientWidth + 2
  groupsCentered.value = !overflow
  edgeL.value = overflow && el.scrollLeft > 2
  edgeR.value = overflow && el.scrollLeft + el.clientWidth < el.scrollWidth - 2
}

watch([activeCat, activeGroups], () => nextTick(updateEdges))

async function load() {
  loading.value = true
  controller?.abort()
  controller = new AbortController()
  try {
    const cfg = await modelPlazaAPI.getModelPlaza({ signal: controller.signal })
    categories.value = cfg.categories ?? []
    // default to first visible category
    activeCat.value = visibleCategories.value[0]?.key ?? ''
    activeGroupIdx.value = 0
    await nextTick()
    updateEdges()
  } catch (error) {
    const err = error as { code?: string; name?: string }
    if (err?.name === 'CanceledError' || err?.code === 'ERR_CANCELED') return
    appStore.showError(t('modelPlaza.loadFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  load()
  window.addEventListener('resize', updateEdges)
})
onUnmounted(() => {
  controller?.abort()
  window.removeEventListener('resize', updateEdges)
})
</script>
