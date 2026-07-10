<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- Header -->
      <div class="ac-header relative overflow-hidden rounded-2xl px-6 py-6 sm:px-8 sm:py-7">
        <div class="pointer-events-none absolute inset-0 ac-header__glow"></div>
        <div class="relative flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div class="flex items-center gap-3">
            <span class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-white/15 text-white ring-1 ring-white/25">
              <Icon name="sparkles" size="md" />
            </span>
            <div>
              <h1 class="text-xl font-bold text-white sm:text-2xl">活动中心</h1>
              <p class="mt-0.5 text-sm text-white/75">
                {{ enabled && items.length ? `当前有 ${items.length} 个活动可参与` : '查看当前可参与的活动' }}
              </p>
            </div>
          </div>
          <button class="ac-refresh" :disabled="loading" @click="loadData">
            <Icon name="refresh" size="sm" :class="loading ? 'animate-spin' : ''" />
            <span class="ml-1.5 text-sm font-medium">刷新</span>
          </button>
        </div>
      </div>

      <!-- Loading: skeleton cards -->
      <div v-if="loading" class="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
        <div v-for="n in 3" :key="`skel-${n}`" class="card overflow-hidden p-0">
          <div class="ac-skel h-28 w-full rounded-none"></div>
          <div class="space-y-3 p-5">
            <div class="ac-skel h-4 w-2/3"></div>
            <div class="ac-skel h-3 w-1/2"></div>
            <div class="ac-skel h-3 w-full"></div>
            <div class="ac-skel mt-3 h-9 w-full"></div>
          </div>
        </div>
      </div>

      <!-- Disabled / empty -->
      <div v-else-if="!enabled || items.length === 0" class="card flex flex-col items-center justify-center px-5 py-16 text-center">
        <span class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-gray-100 text-gray-400 dark:bg-dark-800 dark:text-dark-500">
          <Icon name="gift" size="xl" />
        </span>
        <p class="text-sm font-medium text-gray-600 dark:text-dark-300">{{ !enabled ? '活动中心暂未开启' : '暂无可参与活动' }}</p>
        <p class="mt-1 text-xs text-gray-400 dark:text-dark-500">{{ !enabled ? '敬请期待精彩活动' : '新的活动即将上线，敬请期待' }}</p>
      </div>

      <!-- Activity cards -->
      <div v-else class="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
        <article
          v-for="(item, index) in items"
          :key="item.id"
          class="ac-card group card overflow-hidden p-0"
          :style="{ animationDelay: `${Math.min(index, 8) * 60}ms` }"
        >
          <!-- Banner: cover image or gradient -->
          <div class="ac-banner" :style="bannerStyle(index, item)">
            <img v-if="safeCoverUrl(item)" :src="safeCoverUrl(item)" :alt="item.title" class="ac-cover" />
            <span v-else class="ac-banner__deco"><Icon name="gift" size="xl" /></span>
            <span class="ac-banner__shade"></span>
            <span class="ac-live">
              <span class="ac-live__dot"></span>进行中
            </span>
          </div>

          <!-- Floating icon badge -->
          <span class="ac-badge" :style="badgeStyle(index)">
            <Icon name="gift" size="md" />
          </span>

          <!-- Body -->
          <div class="ac-body">
            <h2 class="truncate text-base font-semibold text-gray-900 dark:text-white">{{ item.title }}</h2>
            <p v-if="item.subtitle" class="mt-1 truncate text-sm font-medium" :style="{ color: accent(index)[0] }">{{ item.subtitle }}</p>
            <p v-if="item.description" class="ac-desc mt-3 text-sm text-gray-600 dark:text-dark-300">{{ item.description }}</p>
            <button class="ac-cta" :style="ctaStyle(index)" @click="openActivity(item)">
              <span>{{ item.action_label || '查看' }}</span>
              <Icon name="arrowRight" size="sm" />
            </button>
          </div>
        </article>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore } from '@/stores'
import { listActivityCenter, type ActivityCenterItem } from '@/api/activityCenter'
import { sanitizeUrl } from '@/utils/url'

const appStore = useAppStore()
const router = useRouter()
const loading = ref(false)
const enabled = ref(false)
const items = ref<ActivityCenterItem[]>([])

// 每张卡片按顺序循环取一组渐变配色，让多活动展示更鲜活（纯展示，不影响数据/跳转逻辑）。
const palette: [string, string][] = [
  ['#6366f1', '#8b5cf6'],
  ['#ec4899', '#f43f5e'],
  ['#f59e0b', '#f97316'],
  ['#10b981', '#14b8a6'],
  ['#06b6d4', '#3b82f6'],
  ['#8b5cf6', '#d946ef'],
]

function accent(index: number): [string, string] {
  return palette[index % palette.length]
}

function bannerStyle(index: number, item: ActivityCenterItem) {
  if (safeCoverUrl(item)) return {}
  const [a, b] = accent(index)
  return { background: `linear-gradient(135deg, ${a}, ${b})` }
}

function safeCoverUrl(item: ActivityCenterItem) {
  return sanitizeUrl(item.cover_url || '')
}

function badgeStyle(index: number) {
  const [a, b] = accent(index)
  return { background: `linear-gradient(135deg, ${a}, ${b})` }
}

function ctaStyle(index: number) {
  const [a, b] = accent(index)
  return { backgroundImage: `linear-gradient(135deg, ${a}, ${b})` }
}

async function loadData() {
  loading.value = true
  try {
    const result = await listActivityCenter()
    enabled.value = result.enabled
    items.value = result.items || []
  } catch (error: any) {
    appStore.showError(error?.message || '加载活动中心失败')
  } finally {
    loading.value = false
  }
}

function openActivity(item: ActivityCenterItem) {
  const externalURL = sanitizeUrl(item.external_url || '')
  if (externalURL) {
    window.open(externalURL, '_blank', 'noopener,noreferrer')
    return
  }
  const routePath = sanitizeUrl(item.route_path || '', { allowRelative: true })
  if (routePath.startsWith('/')) {
    router.push(routePath)
  }
}

onMounted(loadData)
</script>

<style scoped>
/* Header banner */
.ac-header {
  background: linear-gradient(120deg, #4f46e5 0%, #7c3aed 50%, #c026d3 100%);
  box-shadow: 0 18px 40px -24px rgba(124, 58, 237, 0.8);
}
.ac-header__glow {
  background:
    radial-gradient(circle at 12% 20%, rgba(255, 255, 255, 0.22), transparent 45%),
    radial-gradient(circle at 90% 90%, rgba(255, 255, 255, 0.12), transparent 50%);
}
.ac-refresh {
  display: inline-flex;
  align-items: center;
  border-radius: 9999px;
  padding: 7px 16px;
  color: #fff;
  background: rgba(255, 255, 255, 0.16);
  border: 1px solid rgba(255, 255, 255, 0.28);
  transition: background 0.2s ease, transform 0.2s ease;
}
.ac-refresh:hover:not(:disabled) { background: rgba(255, 255, 255, 0.26); }
.ac-refresh:disabled { opacity: 0.7; cursor: not-allowed; }

/* Card */
.ac-card {
  position: relative;
  display: flex;
  flex-direction: column;
  border-radius: 18px;
  transition: transform 0.25s ease, box-shadow 0.25s ease;
  animation: ac-in 0.5s cubic-bezier(0.16, 1, 0.3, 1) both;
}
.ac-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 22px 40px -22px rgba(31, 41, 55, 0.45);
}
@keyframes ac-in {
  from { opacity: 0; transform: translateY(14px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Banner */
.ac-banner {
  position: relative;
  height: 112px;
  overflow: hidden;
}
.ac-cover {
  height: 100%;
  width: 100%;
  object-fit: cover;
  transition: transform 0.4s ease;
}
.ac-card:hover .ac-cover { transform: scale(1.06); }
.ac-banner__deco {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.32);
  transform: scale(2.4) rotate(-8deg);
}
.ac-banner__shade {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, transparent 40%, rgba(0, 0, 0, 0.18) 100%);
}
.ac-live {
  position: absolute;
  top: 12px;
  right: 12px;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  border-radius: 9999px;
  padding: 3px 10px;
  font-size: 11px;
  font-weight: 600;
  color: #fff;
  background: rgba(0, 0, 0, 0.28);
  backdrop-filter: blur(4px);
}
.ac-live__dot {
  height: 6px;
  width: 6px;
  border-radius: 9999px;
  background: #34d399;
  box-shadow: 0 0 0 0 rgba(52, 211, 153, 0.7);
  animation: ac-pulse 1.6s ease-in-out infinite;
}
@keyframes ac-pulse {
  0%, 100% { box-shadow: 0 0 0 0 rgba(52, 211, 153, 0.6); }
  50% { box-shadow: 0 0 0 5px rgba(52, 211, 153, 0); }
}
/* 可参与状态：把「进行中」绿色胶囊变为红色「可参与」，更强的引导 */
.ac-live--hot { background: rgba(239, 68, 68, 0.92); }
.ac-live--hot .ac-live__dot {
  background: #fff;
  box-shadow: 0 0 0 0 rgba(255, 255, 255, 0.7);
  animation: ac-pulse-hot 1.4s ease-in-out infinite;
}
@keyframes ac-pulse-hot {
  0%, 100% { box-shadow: 0 0 0 0 rgba(255, 255, 255, 0.7); }
  50% { box-shadow: 0 0 0 5px rgba(255, 255, 255, 0); }
}

/* Floating badge */
.ac-badge {
  position: absolute;
  top: 88px;
  left: 18px;
  z-index: 2;
  display: flex;
  height: 48px;
  width: 48px;
  align-items: center;
  justify-content: center;
  border-radius: 14px;
  color: #fff;
  box-shadow: 0 8px 18px -6px rgba(0, 0, 0, 0.4);
}
:global(.dark) .ac-badge { box-shadow: 0 8px 18px -6px rgba(0, 0, 0, 0.6), 0 0 0 3px rgba(17, 24, 39, 0.6); }
/* 图标徽标右上角红点：可参与活动的醒目引导 */
.ac-badge__dot {
  position: absolute;
  top: -3px;
  right: -3px;
  height: 12px;
  width: 12px;
  border-radius: 9999px;
  background: #ef4444;
  border: 2px solid #fff;
  box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.7);
  animation: ac-dot-pulse 1.5s ease-in-out infinite;
}
:global(.dark) .ac-badge__dot { border-color: #1f2937; }
@keyframes ac-dot-pulse {
  0%, 100% { box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.6); }
  50% { box-shadow: 0 0 0 5px rgba(239, 68, 68, 0); }
}

/* Body */
.ac-body {
  display: flex;
  flex: 1;
  flex-direction: column;
  padding: 28px 18px 18px;
}
.ac-desc {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.ac-cta {
  margin-top: auto;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  border-radius: 12px;
  padding: 10px 16px;
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  box-shadow: 0 8px 18px -8px rgba(79, 70, 229, 0.7);
  transition: transform 0.15s ease, filter 0.15s ease, box-shadow 0.15s ease;
}
.ac-cta:hover { filter: brightness(1.06); transform: translateY(-1px); }
.ac-cta:active { transform: translateY(0); }

/* Skeleton */
.ac-skel {
  border-radius: 8px;
  background: linear-gradient(90deg, rgba(148, 163, 184, 0.18) 25%, rgba(148, 163, 184, 0.32) 37%, rgba(148, 163, 184, 0.18) 63%);
  background-size: 400% 100%;
  animation: ac-shimmer 1.4s ease infinite;
}
@keyframes ac-shimmer {
  0% { background-position: 100% 0; }
  100% { background-position: 0 0; }
}
</style>
