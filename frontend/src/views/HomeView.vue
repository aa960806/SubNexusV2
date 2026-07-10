<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="homeContent" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContentUrl"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Default Home Page -->
  <div v-else class="ikun-home">
    <div ref="ghostBgRef" class="ghost-background">
      <div class="ghost-container">
        <img src="/cxk.gif" alt="" />
      </div>
    </div>

    <header class="home-header">
      <router-link to="/home" class="brand-mark" aria-label="Home">
        <img :src="siteLogo || '/logo.png'" alt="Logo" />
      </router-link>

      <div class="header-actions">
        <LocaleSwitcher />
        <a
          v-if="docUrl"
          class="icon-button"
          :href="docUrl"
          target="_blank"
          rel="noopener noreferrer"
          :title="t('home.viewDocs')"
        >
          <Icon name="book" size="sm" />
        </a>
        <button class="icon-button" type="button" :title="themeTitle" @click="toggleTheme">
          <Icon v-if="isDark" name="sun" size="sm" />
          <Icon v-else name="moon" size="sm" />
        </button>
        <router-link class="console-pill" :to="isAuthenticated ? dashboardPath : '/login'">
          <span class="avatar-dot">{{ userInitial || 'A' }}</span>
          <span>{{ isAuthenticated ? '控制台' : '登录' }}</span>
          <Icon name="arrowRight" size="xs" />
        </router-link>
      </div>
    </header>

    <main class="main-container">
      <section class="content-left" aria-labelledby="home-title">
        <p class="system-label">Enterprise AI Gateway</p>

        <h1 id="home-title" class="main-title">{{ heroTitle }}</h1>

        <p class="subtitle">{{ siteSubtitle || 'Next-Generation AI Infrastructure Platform' }}</p>

        <p class="description">
          统一 API 协议，全模型接入。企业级性能与稳定性保障。
        </p>

        <div class="cta-row">
          <router-link class="ghost-button" :to="isAuthenticated ? dashboardPath : '/login'">
            立即开始
            <Icon name="arrowRight" size="sm" />
          </router-link>
          <router-link v-if="!isAuthenticated" class="register-link" to="/register">
            创建账号
          </router-link>
        </div>
      </section>
    </main>

    <section class="features-grid" aria-label="Platform features">
      <article class="feature-item">
        <div class="feature-icon">⚡</div>
        <h3>极速响应</h3>
        <p>毫秒级 API 调用延迟，智能负载均衡</p>
      </article>

      <article class="feature-item">
        <div class="feature-icon">🛡️</div>
        <h3>企业级稳定</h3>
        <p>99.99% SLA 保障，全链路监控</p>
      </article>

      <article class="feature-item">
        <div class="feature-icon">🌐</div>
        <h3>多模型统一</h3>
        <p>OpenAI 标准协议，全厂商兼容</p>
      </article>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()

const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || '')
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || '')
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const heroTitle = computed(() => siteName.value || '只因API')

const homeContentUrl = computed(() => sanitizeUrl(homeContent.value))
const isHomeContentUrl = computed(() => Boolean(homeContentUrl.value))

const isDark = ref(document.documentElement.classList.contains('dark'))
const themeTitle = computed(() => (isDark.value ? '切换到浅色模式' : '切换到深色模式'))

const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => (isAdmin.value ? '/admin/dashboard' : '/dashboard'))
const userInitial = computed(() => {
  const email = authStore.user?.email || ''
  return email ? email.charAt(0).toUpperCase() : ''
})

const ghostBgRef = ref<HTMLElement | null>(null)
let mouseX = 0
let mouseY = 0
let currentX = 0
let currentY = 0
let rafId = 0
let pulseTimer: number | undefined
let visualEffectsActive = false

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  } else {
    isDark.value = false
    document.documentElement.classList.remove('dark')
  }
}

function handleMouseMove(event: MouseEvent) {
  mouseX = (event.clientX / window.innerWidth - 0.5) * 2
  mouseY = (event.clientY / window.innerHeight - 0.5) * 2
}

function animateParallax() {
  currentX += (mouseX - currentX) * 0.05
  currentY += (mouseY - currentY) * 0.05

  const offsetX = currentX * 20
  const offsetY = currentY * 20
  if (ghostBgRef.value) {
    ghostBgRef.value.style.transform = `translate(calc(-50% + ${offsetX}px), calc(-50% + ${offsetY}px))`
  }

  rafId = window.requestAnimationFrame(animateParallax)
}

function pulseGhost() {
  const el = ghostBgRef.value
  if (!el || Math.random() <= 0.7) return
  el.style.transition = 'opacity 0.3s ease, transform 0.3s ease-out'
  el.style.opacity = '0.15'

  window.setTimeout(() => {
    if (!ghostBgRef.value) return
    ghostBgRef.value.style.transition = 'opacity 1.5s ease, transform 0.3s ease-out'
    ghostBgRef.value.style.opacity = ''
  }, 300)
}

function startVisualEffects() {
  if (
    visualEffectsActive ||
    homeContent.value.trim() ||
    window.matchMedia('(prefers-reduced-motion: reduce)').matches
  ) {
    return
  }
  visualEffectsActive = true
  window.addEventListener('mousemove', handleMouseMove)
  rafId = window.requestAnimationFrame(animateParallax)
  pulseTimer = window.setInterval(pulseGhost, 8000)
}

function stopVisualEffects() {
  window.removeEventListener('mousemove', handleMouseMove)
  if (rafId) {
    window.cancelAnimationFrame(rafId)
    rafId = 0
  }
  if (pulseTimer) {
    window.clearInterval(pulseTimer)
    pulseTimer = undefined
  }
  visualEffectsActive = false
}

onMounted(() => {
  initTheme()
  authStore.checkAuth()

  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }

  void nextTick(startVisualEffects)
})

watch(homeContent, (content) => {
  if (content.trim()) {
    stopVisualEffects()
    return
  }
  void nextTick(startVisualEffects)
})

onUnmounted(stopVisualEffects)
</script>

<style scoped>
.ikun-home {
  position: relative;
  min-height: 100vh;
  overflow: hidden;
  background:
    radial-gradient(circle at 85% 10%, rgba(20, 184, 166, 0.15), transparent 35%),
    radial-gradient(circle at 35% 70%, rgba(14, 165, 233, 0.1), transparent 38%),
    radial-gradient(circle at 50% 50%, rgba(30, 41, 59, 0.5), transparent 60%),
    #020617;
  color: #ffffff;
  font-family:
    -apple-system,
    BlinkMacSystemFont,
    Inter,
    "SF Pro Display",
    "Segoe UI",
    sans-serif;
}

.ikun-home::before {
  position: absolute;
  inset: 0;
  pointer-events: none;
  content: "";
  background-image:
    linear-gradient(rgba(148, 163, 184, 0.045) 1px, transparent 1px),
    linear-gradient(90deg, rgba(148, 163, 184, 0.045) 1px, transparent 1px);
  background-size: 64px 64px;
  mask-image: linear-gradient(to bottom, #000 0%, transparent 88%);
}

.ghost-background {
  position: fixed;
  top: 50%;
  left: 50%;
  z-index: 0;
  width: 100vw;
  height: 100vh;
  pointer-events: none;
  opacity: 0.045;
  transform: translate(-50%, -50%);
  animation: ghostBreathing 12s ease-in-out infinite;
  transition: transform 0.3s ease-out;
}

.ghost-container {
  position: absolute;
  top: 50%;
  left: 58%;
  width: min(50vw, 760px);
  height: min(72vh, 760px);
  min-width: 360px;
  min-height: 460px;
  transform: translate(-50%, -50%);
}

.ghost-container img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  filter: invert(1) grayscale(1) contrast(150%) brightness(120%) blur(40px);
  mix-blend-mode: screen;
  opacity: 0.5;
}

.home-header {
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: min(100% - 48px, 1152px);
  margin: 0 auto;
  padding: 12px 0;
}

.brand-mark {
  display: inline-flex;
  width: 38px;
  height: 38px;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.92);
  border-radius: 11px;
  box-shadow: 0 10px 28px rgba(20, 184, 166, 0.18);
}

.brand-mark img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.icon-button,
.console-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: rgba(226, 232, 240, 0.86);
  transition:
    color 0.25s ease,
    background 0.25s ease,
    border-color 0.25s ease;
}

.icon-button {
  width: 34px;
  height: 34px;
  border: 0;
  background: transparent;
  cursor: pointer;
}

.icon-button:hover {
  color: #ffffff;
}

.console-pill {
  gap: 6px;
  min-height: 32px;
  padding: 5px 10px 5px 5px;
  font-size: 12px;
  font-weight: 700;
  text-decoration: none;
  background: rgba(30, 41, 59, 0.72);
  border: 1px solid rgba(148, 163, 184, 0.18);
  border-radius: 999px;
}

.console-pill:hover {
  color: #ffffff;
  background: rgba(51, 65, 85, 0.82);
  border-color: rgba(20, 184, 166, 0.4);
}

.avatar-dot {
  display: inline-flex;
  width: 20px;
  height: 20px;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  font-size: 10px;
  font-weight: 800;
  background: linear-gradient(135deg, #2dd4bf, #0f766e);
  border-radius: 50%;
}

.main-container {
  position: relative;
  z-index: 1;
  display: flex;
  min-height: calc(100vh - 62px);
  align-items: center;
  justify-content: flex-start;
  width: min(100% - 48px, 1152px);
  margin: 0 auto;
  padding: 48px 0 180px;
}

.content-left {
  max-width: 600px;
}

.system-label {
  margin-bottom: 40px;
  color: rgba(20, 184, 166, 0.75);
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 2.5px;
  text-transform: uppercase;
  text-shadow: 0 0 20px rgba(20, 184, 166, 0.3);
}

.main-title {
  margin-bottom: 20px;
  color: #ffffff;
  font-size: clamp(48px, 8vw, 72px);
  font-weight: 700;
  letter-spacing: -1px;
  line-height: 1.08;
}

.subtitle {
  margin-bottom: 54px;
  color: #d1d5db;
  font-size: clamp(16px, 2vw, 20px);
  font-weight: 500;
  letter-spacing: -0.2px;
  line-height: 1.6;
}

.description {
  max-width: 500px;
  margin-bottom: 54px;
  color: rgba(226, 232, 240, 0.72);
  font-size: 16px;
  font-weight: 400;
  line-height: 1.8;
}

.cta-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 18px;
}

.ghost-button {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 12px;
  min-height: 52px;
  padding: 0 32px;
  overflow: hidden;
  color: #ffffff;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: 0.5px;
  text-decoration: none;
  background: linear-gradient(135deg, rgba(20, 184, 166, 0.15), rgba(14, 165, 233, 0.15));
  border: 1px solid rgba(20, 184, 166, 0.4);
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(20, 184, 166, 0.1);
  transition:
    all 0.35s cubic-bezier(0.4, 0, 0.2, 1);
}

.ghost-button::before {
  position: absolute;
  inset: 0;
  z-index: -1;
  content: "";
  background: linear-gradient(135deg, #14b8a6, #0ea5e9);
  opacity: 0;
  transition: opacity 0.35s cubic-bezier(0.4, 0, 0.2, 1);
}

.ghost-button:hover {
  color: #ffffff;
  border-color: rgba(20, 184, 166, 0.8);
  box-shadow: 0 8px 24px rgba(20, 184, 166, 0.25);
  transform: translateY(-2px);
}

.ghost-button:hover::before {
  opacity: 1;
}

.register-link {
  color: rgba(226, 232, 240, 0.82);
  font-size: 14px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s ease;
  position: relative;
}

.register-link::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 2px;
  background: linear-gradient(90deg, #14b8a6, #0ea5e9);
  transition: width 0.3s ease;
}

.register-link:hover {
  color: #ffffff;
}

.register-link:hover::after {
  width: 100%;
}

.features-grid {
  position: fixed;
  right: 50%;
  bottom: 50px;
  z-index: 2;
  display: grid;
  grid-template-columns: repeat(3, minmax(150px, 1fr));
  width: min(820px, calc(100% - 48px));
  gap: 1px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
  transform: translateX(50%);
}

.feature-item {
  position: relative;
  min-height: 110px;
  padding: 24px 32px;
  background: rgba(2, 6, 23, 0.9);
  backdrop-filter: blur(12px);
  transition:
    background 0.3s ease,
    transform 0.3s ease,
    border-color 0.3s ease;
  border-bottom: 2px solid transparent;
}

.feature-item:hover {
  background: rgba(15, 23, 42, 0.92);
  transform: translateY(-2px);
  border-bottom-color: rgba(20, 184, 166, 0.5);
}

.feature-icon {
  font-size: 24px;
  margin-bottom: 12px;
  opacity: 0.85;
  transition: transform 0.3s ease;
}

.feature-item:hover .feature-icon {
  transform: scale(1.1);
}

.feature-item h3 {
  margin-bottom: 8px;
  color: rgba(255, 255, 255, 0.9);
  font-size: 14px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.feature-item p {
  color: rgba(226, 232, 240, 0.58);
  font-size: 12px;
  font-weight: 400;
  line-height: 1.6;
}

.system-label,
.main-title,
.subtitle,
.description,
.cta-row {
  opacity: 0;
  animation: fadeInUp 0.8s cubic-bezier(0.4, 0, 0.2, 1) forwards;
}

.system-label {
  animation-delay: 0.1s;
}

.main-title {
  animation-delay: 0.2s;
}

.subtitle {
  animation-delay: 0.3s;
}

.description {
  animation-delay: 0.4s;
}

.cta-row {
  animation-delay: 0.5s;
}

@keyframes ghostBreathing {
  0%,
  100% {
    opacity: 0.045;
  }

  50% {
    opacity: 0.07;
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 900px) {
  .ikun-home {
    min-height: 100svh;
    overflow-x: hidden;
    overflow-y: auto;
  }

  .home-header,
  .main-container {
    width: calc(100% - 40px);
  }

  .main-container {
    align-items: flex-start;
    min-height: auto;
    padding-top: 72px;
    padding-bottom: 36px;
  }

  .system-label {
    margin-bottom: 28px;
  }

  .subtitle,
  .description {
    margin-bottom: 36px;
  }

  .ghost-container {
    left: 62%;
    width: 82vw;
    height: 56vh;
    min-width: 280px;
    min-height: 360px;
  }

  .features-grid {
    position: relative;
    right: auto;
    bottom: auto;
    grid-template-columns: 1fr;
    width: calc(100% - 40px);
    max-width: 560px;
    margin: 0 auto 40px;
    transform: none;
  }
}

@media (max-width: 560px) {
  .home-header {
    width: calc(100% - 32px);
  }

  .header-actions {
    gap: 8px;
  }

  .console-pill {
    padding-right: 8px;
  }

  .console-pill span:nth-child(2) {
    display: none;
  }

  .main-container {
    width: calc(100% - 32px);
    padding-top: 54px;
    padding-bottom: 28px;
  }

  .description {
    font-size: 15px;
  }

  .features-grid {
    width: calc(100% - 32px);
    margin-bottom: 28px;
  }

  .feature-item {
    min-height: auto;
    padding: 18px 20px;
  }

  .feature-icon {
    font-size: 20px;
    margin-bottom: 8px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .ghost-background,
  .feature-item,
  .system-label,
  .main-title,
  .subtitle,
  .description,
  .cta-row {
    animation: none !important;
    transition: none !important;
  }

  .system-label,
  .main-title,
  .subtitle,
  .description,
  .cta-row {
    opacity: 1;
    transform: none;
  }
}
</style>
