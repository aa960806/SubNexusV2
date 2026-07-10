import { existsSync, readFileSync, statSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { describe, expect, it } from 'vitest'

const dir = dirname(fileURLToPath(import.meta.url))
const homeViewPath = resolve(dir, '../HomeView.vue')
const visualAssetPath = resolve(dir, '../../../public/cxk.gif')
const source = readFileSync(homeViewPath, 'utf8')

describe('customized home view', () => {
  it('keeps the forked visual identity and asset', () => {
    expect(source).toContain('class="ikun-home"')
    expect(source).toContain('class="ghost-background"')
    expect(source).toContain('src="/cxk.gif"')
    expect(source).toContain('class="features-grid"')
    expect(existsSync(visualAssetPath)).toBe(true)
    expect(statSync(visualAssetPath).size).toBeGreaterThan(0)
  })

  it('preserves custom home content behavior', () => {
    expect(source).toContain('v-if="homeContent"')
    expect(source).toContain('v-if="isHomeContentUrl"')
    expect(source).toContain('v-html="homeContent"')
  })

  it('retains URL sanitization for configurable resources', () => {
    expect(source).toContain("import { sanitizeUrl } from '@/utils/url'")
    expect(source).toContain('sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo')
    expect(source).toContain('sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl')
    expect(source).toContain('const homeContentUrl = computed(() => sanitizeUrl(homeContent.value))')
  })

  it('cleans up animation work and respects reduced motion', () => {
    expect(source).toContain("window.matchMedia('(prefers-reduced-motion: reduce)').matches")
    expect(source).toContain('onUnmounted(stopVisualEffects)')
    expect(source).toContain('@media (prefers-reduced-motion: reduce)')
  })
})
