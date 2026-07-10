<template>
  <div v-if="currentMessage" class="pointer-events-none fixed left-1/2 top-3 z-[100000] w-[min(820px,calc(100vw-24px))] -translate-x-1/2">
    <div class="pointer-events-auto flex items-center gap-3 overflow-hidden rounded-lg border border-amber-200 bg-white/95 px-3 py-2 shadow-lg shadow-amber-900/10 backdrop-blur dark:border-amber-400/30 dark:bg-dark-900/95">
      <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-amber-100 text-amber-700 dark:bg-amber-400/15 dark:text-amber-300">
        <Icon name="bell" size="sm" />
      </div>
      <div class="min-w-0 flex-1 overflow-hidden">
        <div
          :key="currentMessageKey"
          class="marquee-track whitespace-nowrap text-sm font-medium text-gray-800 dark:text-dark-100"
          :style="{ animationDuration: `${animationDuration}s` }"
          @animationend="finishCurrentRun"
        >
          <span class="inline-flex items-center gap-2">
            <span v-if="currentMessage.title" class="font-semibold text-amber-700 dark:text-amber-300">{{ currentMessage.title }}</span>
            <span>{{ currentMessage.content }}</span>
          </span>
        </div>
      </div>
      <button class="btn-icon shrink-0 text-gray-400 hover:text-gray-700 dark:text-dark-400 dark:hover:text-dark-100" @click="dismiss">
        <Icon name="x" size="sm" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import Icon from '@/components/icons/Icon.vue'
import { listActiveBroadcasts, type ActivityBroadcast } from '@/api/broadcast'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const messages = ref<ActivityBroadcast[]>([])
let timer: number | undefined

const storageKey = computed(() => `sub2api:broadcast-marquee:seen:${authStore.user?.id || 'guest'}`)
const visibleMessages = computed(() => messages.value.filter((item) => !seenKeys.value.has(messageKey(item))))
const currentIndex = ref(0)
const cycleCount = ref(0)
const currentMessage = computed(() => {
  const list = visibleMessages.value
  if (!list.length) return null
  return list[currentIndex.value % list.length] || list[0]
})
// Re-keying on every cycle forces the marquee animation to replay, so each
// active announcement keeps looping (and rotating) until the user dismisses it.
const currentMessageKey = computed(() => currentMessage.value ? `${messageKey(currentMessage.value)}:${cycleCount.value}` : '')
const seenKeys = ref<Set<string>>(new Set())
const animationDuration = computed(() => {
  const item = currentMessage.value
  const textLength = item ? (item.title?.length || 0) + item.content.length : 0
  return Math.min(42, Math.max(16, Math.ceil(textLength / 6)))
})

function messageKey(item: ActivityBroadcast) {
  return `${item.id}:${item.updated_at || item.created_at || ''}`
}

function loadSeenKeys() {
  try {
    const raw = window.localStorage.getItem(storageKey.value)
    const parsed = raw ? JSON.parse(raw) : []
    seenKeys.value = new Set(Array.isArray(parsed) ? parsed.map(String) : [])
  } catch {
    seenKeys.value = new Set()
  }
}

function saveSeenKeys() {
  try {
    const values = Array.from(seenKeys.value).slice(-200)
    window.localStorage.setItem(storageKey.value, JSON.stringify(values))
    seenKeys.value = new Set(values)
  } catch {
    // Ignore storage failures; the marquee still works for the current session.
  }
}

async function loadBroadcasts() {
  if (!authStore.isAuthenticated) {
    messages.value = []
    return
  }
  try {
    const items = await listActiveBroadcasts(12)
    messages.value = items.filter((item) => !seenKeys.value.has(messageKey(item)))
  } catch {
    messages.value = []
  }
}

function startPolling() {
  stopPolling()
  loadBroadcasts()
  timer = window.setInterval(loadBroadcasts, 30000)
}

function stopPolling() {
  if (timer) {
    window.clearInterval(timer)
    timer = undefined
  }
}

function dismiss() {
  markVisibleAsSeen()
  messages.value = []
}

function finishCurrentRun() {
  const list = visibleMessages.value
  if (!list.length) return
  // Advance to the next active announcement and loop. Announcements are NOT
  // marked as seen here, so they keep cycling for repeated exposure; the user
  // dismisses them explicitly via the close button.
  currentIndex.value = (currentIndex.value + 1) % list.length
  cycleCount.value += 1
}

function markVisibleAsSeen() {
  if (!visibleMessages.value.length) return
  const next = new Set(seenKeys.value)
  for (const item of visibleMessages.value) {
    next.add(messageKey(item))
  }
  seenKeys.value = next
  saveSeenKeys()
}

watch(
  () => storageKey.value,
  () => {
    loadSeenKeys()
    messages.value = []
  },
  { immediate: true },
)

watch(
  () => authStore.isAuthenticated,
  (authenticated) => {
    if (authenticated) startPolling()
    else {
      stopPolling()
      messages.value = []
    }
  },
  { immediate: true },
)

onBeforeUnmount(stopPolling)
</script>

<style scoped>
.marquee-track {
  display: inline-block;
  min-width: 100%;
  animation-name: broadcast-marquee;
  animation-timing-function: linear;
  animation-iteration-count: 1;
  animation-fill-mode: forwards;
}

.marquee-track:hover {
  animation-play-state: paused;
}

@keyframes broadcast-marquee {
  0% {
    transform: translateX(100%);
  }
  100% {
    transform: translateX(-100%);
  }
}
</style>
